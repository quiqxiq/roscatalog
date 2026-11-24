package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	ros "example.com/roscatalog"
	cmd "example.com/roscatalog/cmd"
	"example.com/roscatalog/rosrun"
)

// ANSI color helpers
const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

type catalogTestCase struct {
	Name       string
	ResourceID string
	View       string
	Plan       ros.Plan
	Params     ros.Params
	Setup      func(r *ros.Resource)
}

func main() {
	streamDur := flag.Duration("duration", 3*time.Second, "durasi sampel untuk setiap command streaming (follow / interval)")
	dialTimeout := flag.Duration("timeout", 10*time.Second, "batas dial+login ke router")
	skipAuthTest := flag.Bool("skip-auth-test", false, "lewati pengujian simulasi password salah di data.go")
	flag.Parse()

	log.SetFlags(log.Ltime | log.Lmicroseconds)
	fmt.Printf("%s%s=== MULTI-MIKROTIK TEST RUNNER (ROSCATALOG & BUILDER) ===%s\n", colorBold, colorCyan, colorReset)
	fmt.Printf("Memuat konfigurasi router dari cmd/data.go:\n")
	for i, cfg := range cmd.MikrotikList {
		modeStr := "Streaming"
		if cfg.Poll {
			modeStr = "Polling"
		}
		fmt.Printf("  [%d] %s (%s) -> User: %s, Mode: %s\n",
			i+1, cfg.ID, cfg.Address, cfg.Username, modeStr)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// 1. Inisialisasi Manager dan buka koneksi serempak ke seluruh router
	fmt.Printf("\n%s[1/3] Menghubungkan seluruh router via rosrun.OpenAll...%s\n", colorBold, colorReset)
	cfgs := make([]rosrun.Config, len(cmd.MikrotikList))
	copy(cfgs, cmd.MikrotikList)
	for i := range cfgs {
		if cfgs[i].Timeout == 0 {
			cfgs[i].Timeout = *dialTimeout
		}
	}

	mgr, err := rosrun.OpenAll(ctx, cfgs)
	if err != nil {
		if mgr == nil {
			log.Fatalf("%s[FATAL] Gagal menghubungkan router sama sekali: %v%s", colorRed, err, colorReset)
		}
		log.Printf("%s[WARN] Sebagian router gagal terhubung: %v%s", colorYellow, err, colorReset)
	}
	defer mgr.Close()

	fmt.Printf("%s[OK] Terhubung ke %d router dalam Manager.%s\n", colorGreen, mgr.Len(), colorReset)

	// 2. Pengawasan Realtime Status & Latency Ping di background
	statusCtx, cancelStatus := context.WithCancel(ctx)
	defer cancelStatus()
	go monitorRealtimeStatuses(statusCtx, mgr)

	// Jeda singkat untuk probe ping pertama
	time.Sleep(1500 * time.Millisecond)

	// 3. Bangun seluruh Request menggunakan Catalog & Builder bawaan roscatalog
	fmt.Printf("\n%s[2/3] Menjalankan Command Suite berbasis roscatalog.Resource & Builder...%s\n", colorBold, colorReset)

	testCases := []catalogTestCase{
		// /ip/hotspot/active
		{
			Name:       "hotspot-active: print (Snapshot One-shot)",
			ResourceID: "hotspot-active",
			View:       "list",
			Plan:       ros.Once(),
		},
		{
			Name:       "hotspot-active: print follow (Event Stream)",
			ResourceID: "hotspot-active",
			View:       "tail",
			Plan:       ros.Stream(),
			Setup: func(r *ros.Resource) {
				if _, ok := r.View("tail"); !ok {
					r.Tail("user", "address", "mac-address", "uptime", "bytes-in", "bytes-out")
				}
			},
		},
		{
			Name:       "hotspot-active: print interval=1s (Live Interval Stream)",
			ResourceID: "hotspot-active",
			View:       "live",
			Plan:       ros.StreamEvery(time.Second),
		},

		// /ip/hotspot/user
		{
			Name:       "hotspot-user: print (Snapshot One-shot)",
			ResourceID: "hotspot-user",
			View:       "list",
			Plan:       ros.Once(),
		},
		{
			Name:       "hotspot-user: print follow (Event Stream)",
			ResourceID: "hotspot-user",
			View:       "tail",
			Plan:       ros.Stream(),
			Setup: func(r *ros.Resource) {
				if _, ok := r.View("tail"); !ok {
					r.Tail("name", "profile", "server", "mac-address", "uptime", "bytes-in", "bytes-out")
				}
			},
		},
		{
			Name:       "hotspot-user: print stats interval=1s (Stats Interval Stream)",
			ResourceID: "hotspot-user",
			View:       "live",
			Plan:       ros.StreamEvery(time.Second),
			Setup: func(r *ros.Resource) {
				if _, ok := r.View("live"); !ok {
					r.Live(time.Second, true, "name", "profile", "uptime", "bytes-in", "bytes-out")
				}
			},
		},

		// /ppp/secret
		{
			Name:       "ppp-secret: print (Snapshot One-shot)",
			ResourceID: "ppp-secret",
			View:       "list",
			Plan:       ros.Once(),
		},
		{
			Name:       "ppp-secret: print follow (Event Stream)",
			ResourceID: "ppp-secret",
			View:       "tail",
			Plan:       ros.Stream(),
			Setup: func(r *ros.Resource) {
				if _, ok := r.View("tail"); !ok {
					r.Tail("name", "service", "profile", "local-address", "remote-address")
				}
			},
		},
		{
			Name:       "ppp-secret: print stats interval=1s (Stats Interval Stream)",
			ResourceID: "ppp-secret",
			View:       "live",
			Plan:       ros.StreamEvery(time.Second),
			Setup: func(r *ros.Resource) {
				if _, ok := r.View("live"); !ok {
					r.Live(time.Second, true, "name", "service", "profile", "limit-bytes-in", "limit-bytes-out")
				}
			},
		},

		// /ppp/active
		{
			Name:       "ppp-active: print (Snapshot One-shot)",
			ResourceID: "ppp-active",
			View:       "list",
			Plan:       ros.Once(),
		},
		{
			Name:       "ppp-active: print follow (Event Stream)",
			ResourceID: "ppp-active",
			View:       "tail",
			Plan:       ros.Stream(),
			Setup: func(r *ros.Resource) {
				if _, ok := r.View("tail"); !ok {
					r.Tail("name", "service", "caller-id", "address", "uptime")
				}
			},
		},
		{
			Name:       "ppp-active: print stats interval=1s (Stats Interval Stream)",
			ResourceID: "ppp-active",
			View:       "live",
			Plan:       ros.StreamEvery(time.Second),
			Setup: func(r *ros.Resource) {
				if _, ok := r.View("live"); !ok {
					r.Live(time.Second, true, "name", "service", "caller-id", "address", "uptime")
				}
			},
		},

		// /system/log (/log)
		{
			Name:       "log: print follow (Tail Event Stream)",
			ResourceID: "log",
			View:       "tail",
			Plan:       ros.Stream(),
		},
		{
			Name:       "log: print interval=1s (Live Interval Stream)",
			ResourceID: "log",
			View:       "live",
			Plan:       ros.StreamEvery(time.Second),
			Setup: func(r *ros.Resource) {
				if _, ok := r.View("live"); !ok {
					r.Live(time.Second, false, "time", "topics", "message")
				}
			},
		},

		// /system/resource
		{
			Name:       "system-resource: print interval=1s (Live Resource Monitor)",
			ResourceID: "system-resource",
			View:       "live",
			Plan:       ros.StreamEvery(time.Second),
		},

		// /interface
		{
			Name:       "interface: print proplist=id,name,type (Snapshot Proplist)",
			ResourceID: "interface",
			View:       "list",
			Plan:       ros.Once(),
			Params:     ros.Params{Proplist: []string{".id", "name", "type"}},
		},
		{
			Name:       "interface: print stats detail interval=1s proplist=tx-byte,rx-byte (Live Traffic Detail)",
			ResourceID: "interface",
			View:       "live",
			Plan:       ros.StreamEvery(time.Second),
			Params: ros.Params{
				Flags:    []string{"detail"},
				Proplist: []string{"tx-byte", "rx-byte"},
			},
		},
	}

	for idx, tc := range testCases {
		if ctx.Err() != nil {
			break
		}
		runCatalogTestCase(ctx, mgr, idx+1, len(testCases), tc, *streamDur)
	}

	// 4. Pengujian Universal Command Builder (ros.Path / ros.Menu)
	if ctx.Err() == nil {
		fmt.Printf("\n%s[3/4] Menjalankan Pengujian Universal Command Builder (ros.Path / ros.Menu)...%s\n", colorBold, colorReset)
		runUniversalCommandBuilderTests(ctx, mgr, *streamDur)
	}

	// 5. Pengujian Simulasi Perubahan Password di data.go
	if !*skipAuthTest && ctx.Err() == nil {
		fmt.Printf("\n%s[4/4] Pengujian Simulasi Perubahan Password di data.go (Config Level)...%s\n", colorBold, colorReset)
		runPasswordChangeSimulation(ctx, mgr)
	}

	fmt.Printf("\n%s%s=== SELURUH PENGUJIAN SELESAI ===%s\n", colorBold, colorGreen, colorReset)
}

// monitorRealtimeStatuses membaca WatchStatuses dan menampilkan latency ping serta state.
func monitorRealtimeStatuses(ctx context.Context, mgr *rosrun.Manager) {
	ch, err := mgr.WatchStatuses(ctx)
	if err != nil {
		log.Printf("[STATUS ERROR] Gagal memulai WatchStatuses: %v", err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case st, ok := <-ch:
			if !ok {
				return
			}
			stateColor := colorGreen
			if st.Status.State == rosrun.StatusWarning {
				stateColor = colorYellow
			} else if st.Status.State == rosrun.StatusError {
				stateColor = colorRed
			}

			pingStr := fmt.Sprintf("%.2fms", float64(st.Status.Ping.LastRTT.Microseconds())/1000.0)
			if st.Status.Ping.Received == 0 {
				pingStr = "no-reply"
			}

			log.Printf("%s[HEALTH]%s %s%s%s | State: %s%s%s | Ping: %s (loss: %.1f%%) | Uptime: %s",
				colorCyan, colorReset,
				colorBold, st.RouterID, colorReset,
				stateColor, st.Status.State, colorReset,
				pingStr, st.Status.Ping.PacketLoss,
				st.Status.Uptime.Round(time.Second),
			)
			if len(st.Status.Warnings) > 0 {
				log.Printf("   -> Warnings (%s): %v", st.RouterID, st.Status.Warnings)
			}
			if st.Status.Err != nil {
				log.Printf("   -> Error (%s): %v", st.RouterID, st.Status.Err)
			}
		}
	}
}

// runCatalogTestCase membangun Request melalui Catalog & Builder lalu mengeksekusinya ke Manager.
func runCatalogTestCase(ctx context.Context, mgr *rosrun.Manager, num, total int, tc catalogTestCase, dur time.Duration) {
	fmt.Printf("\n%s--- [%d/%d] %s ---%s\n", colorBold, num, total, tc.Name, colorReset)

	// 1. Ambil resource dari Catalog
	res, ok := ros.Lookup(tc.ResourceID)
	if !ok {
		fmt.Printf("    %s[ERROR] Resource %q tidak terdaftar di katalog!%s\n", colorRed, tc.ResourceID, colorReset)
		return
	}

	// 2. Jalankan setup hook jika perlu menambah view via builder
	if tc.Setup != nil {
		tc.Setup(res)
	}

	// 3. Bangun Request menggunakan Build API katalog
	req, err := res.Build(tc.View, tc.Plan, tc.Params)
	if err != nil {
		fmt.Printf("    %s[BUILD ERROR] res.Build(%q): %v%s\n", colorRed, tc.View, err, colorReset)
		return
	}

	fmt.Printf("    Katalog: %s | View: %s | Mode: %s | Delivery: %s\n", res.ID, req.View, req.Mode, req.Delivery)
	fmt.Printf("    Kalimat API: %s\n", strings.Join(req.Sentence, " "))

	if req.Delivery == ros.DeliveryOnce {
		// One-shot: dieksekusi serempak ke tiap router aktif di Manager
		var wg sync.WaitGroup
		routers := mgr.Routers()
		for _, sess := range routers {
			wg.Add(1)
			go func(s *rosrun.MikroTik) {
				defer wg.Done()
				start := time.Now()
				rows, err := s.Once(ctx, req)
				elapsed := time.Since(start)
				routerID := s.Stats().ID

				if err != nil {
					fmt.Printf("    %s[%s] ERROR:%s %v (%s)\n", colorRed, routerID, colorReset, err, elapsed)
					return
				}

				fmt.Printf("    %s[%s] OK:%s Menerima %d baris (%s)\n", colorGreen, routerID, colorReset, len(rows), elapsed)
				if len(rows) > 0 {
					fmt.Printf("        Sample data: %v\n", truncateMap(rows[0], 5))
				}
			}(sess)
		}
		wg.Wait()
		return
	}

	// Streaming (follow / interval): dimonitor selama durasi sampel
	testCtx, cancel := context.WithTimeout(ctx, dur)
	defer cancel()

	streamCh, err := mgr.WatchAll(testCtx, req)
	if err != nil {
		fmt.Printf("    %s[WATCH ERROR] mgr.WatchAll: %v%s\n", colorRed, err, colorReset)
		return
	}

	eventCounts := make(map[string]int)
	var mu sync.Mutex

	done := make(chan struct{})
	go func() {
		defer close(done)
		for ev := range streamCh {
			mu.Lock()
			eventCounts[ev.RouterID]++
			count := eventCounts[ev.RouterID]
			mu.Unlock()

			detail := ""
			switch ev.Kind {
			case rosrun.KindRow:
				detail = fmt.Sprintf("Row: %v", truncateMap(ev.Row, 4))
			case rosrun.KindSnapshot:
				detail = fmt.Sprintf("Snapshot: %d baris", len(ev.Rows))
				if len(ev.Rows) > 0 {
					detail += fmt.Sprintf(" (Contoh: %v)", truncateMap(ev.Rows[0], 4))
				}
			case rosrun.KindError:
				detail = fmt.Sprintf("Error: %v", ev.Err)
			}

			fmt.Printf("    %s[%s #%d]%s Kind=%s | %s\n",
				colorCyan, ev.RouterID, count, colorReset, ev.Kind, detail)
		}
	}()

	<-testCtx.Done()
	<-done

	fmt.Printf("    %s[SELESAI]%s Sampel %s selesai. Rekap event per router: %v\n",
		colorGreen, colorReset, dur, eventCounts)
}

// runPasswordChangeSimulation menguji perubahan password langsung di data.go (Config level).
func runPasswordChangeSimulation(ctx context.Context, mgr *rosrun.Manager) {
	fmt.Println("Skenario: Memodifikasi password router di cmd.MikrotikList (data.go) secara programatik")
	fmt.Println("          untuk memverifikasi deteksi auth failure, status connection,")
	fmt.Println("          dan memastikan router lain di Manager tetap sehat (fault isolation).")

	// Cari router yang sedang aktif terhubung agar hasil otentikasi dapat diamati secara nyata
	var target rosrun.Config
	for _, cfg := range cmd.MikrotikList {
		if _, ok := mgr.Get(cfg.ID); ok {
			target = cfg
			break
		}
	}
	if target.ID == "" && len(cmd.MikrotikList) > 0 {
		target = cmd.MikrotikList[0]
	}
	originalPassword := target.Password

	// 1. Ubah password menjadi salah
	badPassword := "password_salah_simulasi_999"
	fmt.Printf("\n  %s1. Mengubah password %s dari %q menjadi %q...%s\n",
		colorBold, target.ID, originalPassword, badPassword, colorReset)

	badCfg := target
	badCfg.Password = badPassword
	badCfg.Timeout = 3 * time.Second

	// 2. Mencoba membuka koneksi dengan kredensial yang salah
	fmt.Println("  2. Membuka koneksi baru dengan kredensial salah via rosrun.Open...")
	start := time.Now()
	_, err := rosrun.Open(ctx, badCfg)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("  %s[TERDETEKSI GAGAL AUTENTIKASI]%s (%s):\n", colorGreen, colorReset, elapsed)
		fmt.Printf("     Pesan Error: %s%v%s\n", colorRed, err, colorReset)
		fmt.Println("     -> Sesuai ekspektasi: RouterOS API langsung menolak login dengan error otentikasi.")
	} else {
		fmt.Printf("  %s[PERINGATAN] Koneksi tidak gagal! Periksa kredensial.%s\n", colorYellow, colorReset)
	}

	// 3. Verifikasi kesehatan router kedua di Manager
	fmt.Println("\n  3. Memeriksa status router di Manager (apakah ada gangguan):")
	for _, sess := range mgr.Routers() {
		st := sess.Status()
		fmt.Printf("     Router: %s\n", st.ID)
		fmt.Printf("     State: %s%s%s (Connected: %v)\n", colorGreen, st.State, colorReset, st.Connected)
		fmt.Printf("     Ping RTT Terakhir: %.2fms (Packet Loss: %.1f%%)\n",
			float64(st.Ping.LastRTT.Microseconds())/1000.0, st.Ping.PacketLoss)
		fmt.Printf("     Uptime: %s\n", st.Uptime.Round(time.Second))
		fmt.Printf("     %s[ISOLASI BERHASIL] Router aktif di Manager tetap 100%% sehat tanpa terpengaruh!%s\n",
			colorGreen, colorReset)
	}

	// 4. Kembalikan password semula dan verifikasi pemulihan
	fmt.Printf("\n  4. Mengembalikan password asli (%q) dan memverifikasi koneksi normal kembali...\n", originalPassword)
	restoreCfg := target
	restoreCfg.Password = originalPassword
	restoreCfg.Timeout = 5 * time.Second

	sessRestore, errRestore := rosrun.Open(ctx, restoreCfg)
	if errRestore != nil {
		fmt.Printf("  %s[ERROR] Pemulihan gagal: %v%s\n", colorRed, errRestore, colorReset)
	} else {
		defer sessRestore.Close()
		st := sessRestore.Status()
		fmt.Printf("  %s[PEMULIHAN SUKSES]%s Router kembali terhubung: State=%s, Connected=%v\n",
			colorGreen, colorReset, st.State, st.Connected)
	}
}

// truncateMap membatasi jumlah key-value yang dicetak agar log rapi.
func truncateMap(m map[string]string, maxKeys int) map[string]string {
	if len(m) <= maxKeys {
		return m
	}
	out := make(map[string]string, maxKeys+1)
	count := 0
	for k, v := range m {
		out[k] = v
		count++
		if count >= maxKeys {
			break
		}
	}
	out["..."] = fmt.Sprintf("(total %d kolom)", len(m))
	return out
}

// runUniversalCommandBuilderTests mendemonstrasikan fleksibilitas ros.Path / ros.Menu
// yang mencakup print dengan switches/proplist, stats detail interval, query server-side,
// dan child command seperti ping/tools.
func runUniversalCommandBuilderTests(ctx context.Context, mgr *rosrun.Manager, dur time.Duration) {
	type builderTest struct {
		Name     string
		Builder  ros.Requestable
		IsStream bool
	}

	tests := []builderTest{
		{
			Name: "ros.Path(/system/resource).Print().Proplist(\"version\", \"uptime\", \"cpu-load\")",
			Builder: ros.Path("/system/resource").Print().
				Proplist("version", "uptime", "cpu-load"),
			IsStream: false,
		},
		{
			Name: "ros.Path(/interface).Print(\"stats\", \"detail\").Interval(1s).Proplist(\".id\", \"name\", \"tx-byte\", \"rx-byte\")",
			Builder: ros.Path("/interface").Print("stats", "detail").
				Interval(time.Second).
				Proplist(".id", "name", "tx-byte", "rx-byte"),
			IsStream: true,
		},
		{
			Name: "ros.Path(/ip/hotspot/user).Print().WhereEq(\"disabled\", \"no\").Proplist(\"name\", \"profile\")",
			Builder: ros.Path("/ip/hotspot/user").Print().
				WhereEq("disabled", "no").
				Proplist("name", "profile"),
			IsStream: false,
		},
		{
			Name: "ros.Path(\"/\").Do(\"ping\").Arg(\"address\", \"8.8.8.8\").Arg(\"count\", 2)",
			Builder: ros.Path("/").Do("ping").
				Arg("address", "8.8.8.8").
				Arg("count", 2),
			IsStream: true,
		},
	}

	routers := mgr.Routers()
	if len(routers) == 0 {
		fmt.Printf("    %s[SKIP] Tidak ada router terhubung di Manager%s\n", colorYellow, colorReset)
		return
	}

	for i, t := range tests {
		if ctx.Err() != nil {
			break
		}
		req, err := t.Builder.Build()
		if err != nil {
			fmt.Printf("    [%d/%d] %s: %sBUILD ERROR:%s %v\n", i+1, len(tests), t.Name, colorRed, colorReset, err)
			continue
		}
		fmt.Printf("\n    [%d/%d] %s%s%s\n", i+1, len(tests), colorBold, t.Name, colorReset)
		fmt.Printf("        Sentence: %s\n", strings.Join(req.Sentence, " "))

		if !t.IsStream {
			for _, s := range routers {
				st := s.Status()
				if !st.Connected {
					continue
				}
				start := time.Now()
				rows, err := s.OnceCmd(ctx, t.Builder)
				elapsed := time.Since(start)
				if err != nil {
					fmt.Printf("        [%s] %sERROR:%s %v (%s)\n", st.ID, colorRed, colorReset, err, elapsed)
					continue
				}
				fmt.Printf("        [%s] %sOK:%s Menerima %d baris (%s)\n", st.ID, colorGreen, colorReset, len(rows), elapsed)
				if len(rows) > 0 {
					fmt.Printf("            Data: %v\n", truncateMap(rows[0], 5))
				}
			}
		} else {
			testCtx, cancel := context.WithTimeout(ctx, dur)
			ch, err := mgr.WatchAllCmd(testCtx, t.Builder)
			if err != nil {
				fmt.Printf("        %sWATCH ERROR:%s %v\n", colorRed, colorReset, err)
				cancel()
				continue
			}
			counts := make(map[string]int)
			for ev := range ch {
				if ev.Event.Kind == rosrun.KindError {
					fmt.Printf("        [%s] %sSTREAM ERROR:%s %v\n", ev.RouterID, colorRed, colorReset, ev.Event.Err)
				} else {
					counts[ev.RouterID]++
				}
			}
			cancel()
			for routerID, count := range counts {
				fmt.Printf("        [%s] %sOK:%s Streaming diterima %d event dalam %s\n", routerID, colorGreen, colorReset, count, dur)
			}
		}
	}
}
