// Perintah multicheck menjalankan beberapa MikroTik sekaligus: status native
// /ping, stream interval + stats, native monitor-traffic, follow-only, poll,
// dan command tulis/baca di tengah semua stream.
//
// Contoh:
//
//	go run ./cmd/multicheck \
//	  -router 192.168.122.60:8728 -router 192.168.122.26:8728 \
//	  -router 192.168.122.189:8728 \
//	  -user admin -pass secret -duration 20s
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sort"
	"strings"
	"sync"
	"time"

	ros "example.com/roscatalog"
	"example.com/roscatalog/rosrun"
)

type stringList []string

func (s *stringList) String() string { return strings.Join(*s, ",") }
func (s *stringList) Set(v string) error {
	*s = append(*s, v)
	return nil
}

func main() {
	var routers stringList
	flag.Var(&routers, "router", "alamat API host[:port]; boleh berulang")
	user := flag.String("user", "admin", "user API")
	pass := flag.String("pass", "", "password API")
	useTLS := flag.Bool("tls", false, "pakai API TLS (8729)")
	insecure := flag.Bool("insecure", true, "lewati verifikasi sertifikat TLS")
	duration := flag.Duration("duration", 20*time.Second, "lama pengujian")
	timeout := flag.Duration("timeout", 10*time.Second, "batas dial+login")
	streamEvery := flag.Duration("stream-every", time.Second, "interval stream (1s..5s)")
	pollEvery := flag.Duration("poll-every", 5*time.Second, "jeda poll (5s..30s)")
	statusEvery := flag.Duration("status-every", time.Second, "interval ping status (1s..5s)")
	noTLSWarning := flag.Bool("no-tls-warning", false, "jangan tampilkan warning TLS mati")
	iface := flag.String("interface", "ether1", "interface untuk monitor-traffic")
	flag.Parse()

	if len(routers) == 0 {
		log.Fatal("sebutkan minimal satu -router, mis. -router 192.168.122.189:8728")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	ctx, cancel := context.WithTimeout(ctx, *duration)
	defer cancel()

	cfgs := make([]rosrun.Config, 0, len(routers))
	for i, addr := range routers {
		cfgs = append(cfgs, rosrun.Config{
			ID:          addr,
			Address:     addr,
			Username:    *user,
			Password:    *pass,
			TLS:         *useTLS,
			Insecure:    *insecure,
			Timeout:     *timeout,
			Stream:      i < len(routers)-1, // satu router sengaja poll-preferred
			StreamEvery: *streamEvery,
			PollEvery:   *pollEvery,
			Status: rosrun.StatusConfig{
				Every:              *statusEvery,
				SuppressTLSWarning: *noTLSWarning,
			},
			Logf: log.Printf,
		})
	}

	start := time.Now()
	mgr, err := rosrun.OpenAll(ctx, cfgs)
	if err != nil {
		log.Printf("sebagian router gagal: %v", err)
	}
	if mgr == nil {
		log.Fatal("tidak ada router yang berhasil dibuka")
	}
	defer mgr.Close()
	log.Printf("manager siap dalam %s: %s", time.Since(start).Round(time.Millisecond), strings.Join(mgr.IDs(), ", "))

	counts := &counter{kinds: map[string]map[rosrun.Kind]int{}, errors: map[string]int{}}

	// Status realtime semua router.
	go consumeStatuses(ctx, mgr, counts)

	// Otomatis mengikuti Config.Stream tiap router.
	live, _ := ros.Lookup("interface")
	go consumeEvents("auto:"+live.ID, mustWatchAuto(ctx, mgr, live), counts)

	// Stream interval + stats di semua router.
	streamReq, err := live.Build("live", ros.AnyVersion, ros.StreamEvery(*streamEvery), ros.Params{})
	if err != nil {
		log.Fatalf("build stream: %v", err)
	}
	go consumeEvents("interval", mustWatchAll(ctx, mgr, streamReq), counts)

	// Poll terus-menerus di semua router, berjalan bersamaan dengan stream.
	pollReq, err := live.Build("live", ros.AnyVersion, ros.Poll(*pollEvery), ros.Params{})
	if err != nil {
		log.Fatalf("build poll: %v", err)
	}
	go consumeEvents("poll", mustWatchAll(ctx, mgr, pollReq), counts)

	// Native monitor-traffic.
	mt, _ := ros.Lookup("interface-monitor-traffic")
	mtReq, err := mt.Build("run", ros.AnyVersion, ros.Stream(), ros.Params{
		Args: map[string]string{"interface": *iface},
	})
	if err != nil {
		log.Fatalf("build monitor: %v", err)
	}
	go consumeEvents("monitor", mustWatchAll(ctx, mgr, mtReq), counts)

	// Follow-only dan follow (buffer lama + entri baru).
	lg, _ := ros.Lookup("log")
	lgReq, err := lg.Build("tail-only", ros.AnyVersion, ros.Stream(), ros.Params{})
	if err != nil {
		log.Fatalf("build follow-only: %v", err)
	}
	go consumeEvents("follow-only", mustWatchAll(ctx, mgr, lgReq), counts)

	tailReq, err := lg.Build("tail", ros.AnyVersion, ros.Stream(), ros.Params{})
	if err != nil {
		log.Fatalf("build follow: %v", err)
	}
	go consumeEvents("follow", mustWatchAll(ctx, mgr, tailReq), counts)

	// Command di koneksi cmd harus tetap jalan walau semua stream aktif.
	execTicker := time.NewTicker(2 * time.Second)
	defer execTicker.Stop()
	for {
		select {
		case <-ctx.Done():
			report(mgr, counts)
			return
		case <-execTicker.C:
			for _, sess := range mgr.Routers() {
				done := time.Now()
				_, err := sess.Exec(ctx, []string{"/system/identity/print"})
				if err != nil {
					counts.addError(idOf(sess))
					log.Printf("[%s] exec gagal: %v", idOf(sess), err)
					continue
				}
				counts.markExec(idOf(sess), time.Since(done))
			}
		}
	}
}

func idOf(s *rosrun.MikroTik) string {
	// Config.ID = alamat, jadi Stats.ID stabil untuk log.
	return s.Stats().ID
}

func mustWatchAll(ctx context.Context, mgr *rosrun.Manager, req ros.Request) <-chan rosrun.RouterEvent {
	ch, err := mgr.WatchAll(ctx, req)
	if err != nil {
		log.Fatalf("WatchAll: %v", err)
	}
	return ch
}

func mustWatchAuto(ctx context.Context, mgr *rosrun.Manager, res *ros.Resource) <-chan rosrun.RouterEvent {
	ch, err := mgr.WatchAutoAll(ctx, res, "live", ros.AnyVersion, ros.Plan{}, ros.Params{})
	if err != nil {
		log.Fatalf("WatchAutoAll: %v", err)
	}
	return ch
}

func consumeEvents(name string, ch <-chan rosrun.RouterEvent, c *counter) {
	for ev := range ch {
		c.add(ev)
		if ev.Kind == rosrun.KindError {
			log.Printf("[%s] %s error: %v", ev.RouterID, name, ev.Err)
		}
	}
}

func consumeStatuses(ctx context.Context, mgr *rosrun.Manager, c *counter) {
	ch, err := mgr.WatchStatuses(ctx)
	if err != nil {
		log.Printf("WatchStatuses: %v", err)
		return
	}
	for rs := range ch {
		st := rs.Status
		c.setStatus(rs.RouterID, st.State)
		warn := ""
		if len(st.Warnings) > 0 {
			warn = " warnings=" + strings.Join(st.Warnings, "; ")
		}
		log.Printf("[%s] status=%s connected=%v rtt=%s loss=%.1f%%%s%s",
			rs.RouterID, st.State, st.Connected, st.Ping.LastRTT,
			st.Ping.PacketLoss, warn, errText(st.Err))
	}
}

func errText(err error) string {
	if err == nil {
		return ""
	}
	return " err=" + err.Error()
}

type counter struct {
	mu     sync.Mutex
	kinds  map[string]map[rosrun.Kind]int
	errors map[string]int
	execs  map[string]int
	execMs map[string]time.Duration
	status map[string]rosrun.StatusState
}

func (c *counter) add(ev rosrun.RouterEvent) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.kinds[ev.RouterID] == nil {
		c.kinds[ev.RouterID] = map[rosrun.Kind]int{}
	}
	c.kinds[ev.RouterID][ev.Kind]++
	if ev.Kind == rosrun.KindError {
		c.errors[ev.RouterID]++
	}
}

func (c *counter) addError(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.errors == nil {
		c.errors = map[string]int{}
	}
	c.errors[id]++
}

func (c *counter) markExec(id string, d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.execs == nil {
		c.execs = map[string]int{}
		c.execMs = map[string]time.Duration{}
	}
	c.execs[id]++
	c.execMs[id] += d
}

func (c *counter) setStatus(id string, st rosrun.StatusState) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.status == nil {
		c.status = map[string]rosrun.StatusState{}
	}
	c.status[id] = st
}

func report(mgr *rosrun.Manager, c *counter) {
	c.mu.Lock()
	defer c.mu.Unlock()
	fmt.Println("\n=== ringkasan ===")
	ids := mgr.IDs()
	sort.Strings(ids)
	for _, id := range ids {
		exec := c.execs[id]
		avg := time.Duration(0)
		if exec > 0 {
			avg = c.execMs[id] / time.Duration(exec)
		}
		fmt.Printf("%s status=%s errors=%d exec=%d avgExec=%s\n",
			id, c.status[id], c.errors[id], exec, avg.Round(time.Millisecond))
		kinds := make([]string, 0, len(c.kinds[id]))
		for k, n := range c.kinds[id] {
			kinds = append(kinds, fmt.Sprintf("%s=%d", k, n))
		}
		sort.Strings(kinds)
		fmt.Printf("  event: %s\n", strings.Join(kinds, " "))
	}
}
