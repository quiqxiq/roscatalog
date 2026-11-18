// Perintah roscheck menguji satu sesi MikroTik: sambungan, dua koneksi,
// reconnect, sekali jalan, dan stream/poll.
//
//	# uji koneksi saja
//	go run ./cmd/roscheck -addr 192.168.122.189:8728 -user admin -pass secret
//
//	# stream live 1 detik
//	go run ./cmd/roscheck -addr 192.168.122.189:8728 -user admin -pass secret \
//	  -resource hotspot-active -view live -stream -every 1s
//
//	# poll daftar user
//	go run ./cmd/roscheck -addr 192.168.122.189:8728 -user admin -pass secret \
//	  -resource hotspot-user -view list -once
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
	"time"

	ros "example.com/roscatalog"
	"example.com/roscatalog/rosrun"
)

func main() {
	var (
		addr     = flag.String("addr", "192.168.88.1:8728", "alamat API router")
		user     = flag.String("user", "admin", "user API")
		pass     = flag.String("pass", "", "password API")
		useTLS   = flag.Bool("tls", false, "pakai API TLS (8729)")
		insecure = flag.Bool("insecure", true, "lewati verifikasi sertifikat TLS")

		resource = flag.String("resource", "", "ID resource katalog; kosong = uji koneksi")
		view     = flag.String("view", "", "nama view; kosong = pilih otomatis")
		ver      = flag.String("ver", "", "v6 atau v7")
		once     = flag.Bool("once", false, "jalankan sekali, bukan stream/poll")
		stream   = flag.Bool("stream", false, "preferensi stream (bawaan poll)")
		every    = flag.Duration("every", time.Second, "interval stream atau jeda poll")
		count    = flag.Int("n", 3, "jumlah event sebelum berhenti (0 = sampai Ctrl+C)")

		timeout = flag.Duration("timeout", 10*time.Second, "batas dial+login")
	)
	params := ros.Params{Args: map[string]string{}}
	flag.Func("arg", "argumen command native, mis. interface=ether1", func(s string) error {
		k, v, _ := strings.Cut(s, "=")
		params.Args[k] = v
		return nil
	})
	flag.Func("q", "query sisi server prop=nilai", func(s string) error {
		k, v, _ := strings.Cut(s, "=")
		params.Query = append(params.Query, ros.Query{Prop: k, Value: v})
		return nil
	})
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cfg := rosrun.Config{
		Address:   *addr,
		Username:  *user,
		Password:  *pass,
		TLS:       *useTLS,
		Insecure:  *insecure,
		Timeout:   *timeout,
		Stream:    *stream,
		PollEvery: *every,
		Logf:      log.Printf,
	}

	start := time.Now()
	sess, err := rosrun.Open(ctx, cfg)
	if err != nil {
		log.Fatalf("open: %v", err)
	}
	defer func() {
		if err := sess.Close(); err != nil {
			log.Printf("close: %v", err)
		}
	}()
	log.Printf("tersambung dalam %s ke %s (dua koneksi)", time.Since(start).Round(time.Millisecond), cfg.Address)

	if *resource == "" {
		checkConnect(ctx, sess)
		return
	}

	res, ok := ros.Lookup(*resource)
	if !ok {
		log.Fatalf("resource %q tidak ada", *resource)
	}
	viewName := *view
	if viewName == "" {
		viewName = defaultView(res)
	}

	if *once {
		req, err := res.Build(viewName, ros.Version(*ver), ros.Once(), params)
		if err != nil {
			log.Fatalf("build: %v", err)
		}
		runOnce(ctx, sess, res, req)
		return
	}

	plan := sess.PlanFor(res, viewName, ros.Version(*ver))
	plan.Every = *every
	plan.Diff = true
	req, err := res.Build(viewName, ros.Version(*ver), plan, params)
	if err != nil {
		log.Fatalf("build: %v", err)
	}
	log.Printf("watch %s/%s delivery=%s every=%s", res.ID, viewName, req.Delivery, plan.Every)

	events, err := sess.Watch(ctx, req)
	if err != nil {
		log.Fatalf("watch: %v", err)
	}

	n := 0
	for ev := range events {
		switch ev.Kind {
		case rosrun.KindSnapshot:
			log.Printf("snapshot %d baris", len(ev.Rows))
			for _, row := range ev.Rows {
				log.Printf("  %s", formatRow(row))
			}
		case rosrun.KindRow:
			log.Printf("row %s", formatRow(ev.Row))
		case rosrun.KindAdded:
			log.Printf("added %s", formatRow(ev.Row))
		case rosrun.KindChanged:
			log.Printf("changed %s", formatRow(ev.Row))
		case rosrun.KindRemoved:
			log.Printf("removed %s", formatRow(ev.Row))
		case rosrun.KindCount:
			log.Printf("count %d", ev.Count)
		case rosrun.KindReconnecting:
			log.Printf("koneksi putus; menyambung ulang...")
		case rosrun.KindError:
			log.Printf("error: %v", ev.Err)
		}

		if ev.Kind != rosrun.KindReconnecting {
			n++
		}
		if *count > 0 && n >= *count {
			return
		}
	}
}

func checkConnect(ctx context.Context, sess *rosrun.MikroTik) {
	for _, path := range []string{"/system/identity/print", "/system/resource/print"} {
		done := time.Now()
		out, err := sess.Exec(ctx, []string{path})
		if err != nil {
			log.Fatalf("%s: %v", path, err)
		}
		log.Printf("%s ok dalam %s (ret=%q)", path, time.Since(done).Round(time.Millisecond), out)
	}

	st := sess.Stats()
	log.Printf("stats: epoch=%d reconnects=%d cmdDown=%v feedDown=%v uptime=%s",
		st.Epoch, st.Reconnects, st.CmdDown, st.FeedDown, st.Uptime.Round(time.Millisecond))
}

func runOnce(ctx context.Context, sess *rosrun.MikroTik, res *ros.Resource, req ros.Request) {
	start := time.Now()
	rows, err := sess.Once(ctx, req)
	if err != nil {
		log.Fatalf("once: %v", err)
	}
	log.Printf("%s/%s: %d baris dalam %s", res.ID, req.View, len(rows), time.Since(start).Round(time.Millisecond))
	for _, row := range rows {
		log.Printf("  %s", formatRow(row))
	}
}

func defaultView(res *ros.Resource) string {
	for _, name := range []string{"list", "live", "run", "detail"} {
		if _, ok := res.View(name, ros.AnyVersion); ok {
			return name
		}
	}
	if len(res.Views) > 0 {
		return res.Views[0].Name
	}
	log.Fatalf("resource %s tidak punya view", res.ID)
	return ""
}

func formatRow(row map[string]string) string {
	keys := make([]string, 0, len(row))
	for k := range row {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	for i, k := range keys {
		if i > 0 {
			b.WriteString(" ")
		}
		fmt.Fprintf(&b, "%s=%s", k, row[k])
	}
	return b.String()
}
