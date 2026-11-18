package rosrun_test

import (
	"context"
	"os"
	"testing"
	"time"

	ros "example.com/roscatalog"
	"example.com/roscatalog/rosrun"
)

func liveConfig(t *testing.T) rosrun.Config {
	t.Helper()

	addr := os.Getenv("ROUTEROS_LIVE_ADDR")
	if addr == "" {
		t.Skip("set ROUTEROS_LIVE_ADDR untuk uji ke router sungguhan")
	}
	user := os.Getenv("ROUTEROS_LIVE_USER")
	if user == "" {
		user = "admin"
	}

	return rosrun.Config{
		Address:   addr,
		Username:  user,
		Password:  os.Getenv("ROUTEROS_LIVE_PASS"),
		Timeout:   10 * time.Second,
		PollEvery: time.Second,
	}
}

func openLive(t *testing.T) (*rosrun.MikroTik, context.Context) {
	t.Helper()

	cfg := liveConfig(t)
	sess, err := rosrun.Open(context.Background(), cfg)
	if err != nil {
		t.Fatalf("open %s: %v", cfg.Address, err)
	}
	t.Cleanup(func() { _ = sess.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)
	return sess, ctx
}

func TestLiveSessionConnect(t *testing.T) {
	sess, ctx := openLive(t)

	for _, path := range []string{"/system/identity/print", "/system/resource/print"} {
		if _, err := sess.Exec(ctx, []string{path}); err != nil {
			t.Fatalf("%s: %v", path, err)
		}
	}

	st := sess.Stats()
	if !st.Open || st.CmdDown || st.FeedDown {
		t.Fatalf("unexpected stats: %+v", st)
	}
	t.Logf("connected: epoch=%d reconnects=%d uptime=%s", st.Epoch, st.Reconnects, st.Uptime.Round(time.Millisecond))
}

func TestLiveSessionStreamAndExec(t *testing.T) {
	sess, ctx := openLive(t)

	ifs, ok := ros.Lookup("interface")
	if !ok {
		t.Fatal("resource interface tidak ada")
	}
	req, err := ifs.Build("list", ros.AnyVersion, ros.Once(), ros.Params{})
	if err != nil {
		t.Fatalf("build interface: %v", err)
	}
	rows, err := sess.Once(ctx, req)
	if err != nil {
		t.Fatalf("interface list: %v", err)
	}

	name := ""
	for _, row := range rows {
		if row["name"] != "" && row["disabled"] != "true" {
			name = row["name"]
			break
		}
	}
	if name == "" {
		t.Skip("tidak ada interface untuk diuji")
	}

	mt, ok := ros.Lookup("interface-monitor-traffic")
	if !ok {
		t.Fatal("resource interface-monitor-traffic tidak ada")
	}
	mreq, err := mt.Build("run", ros.AnyVersion, ros.StreamEvery(time.Second), ros.Params{
		Args: map[string]string{"interface": name},
	})
	if err != nil {
		t.Fatalf("build monitor: %v", err)
	}

	ch, err := sess.Watch(ctx, mreq)
	if err != nil {
		t.Fatalf("watch monitor: %v", err)
	}

	// Command harus tetap jalan saat stream terbuka: bukti dua koneksi.
	for i := 0; i < 3; i++ {
		if _, err := sess.Exec(ctx, []string{"/system/identity/print"}); err != nil {
			t.Fatalf("exec during stream #%d: %v", i+1, err)
		}
	}

	select {
	case ev, ok := <-ch:
		if !ok {
			t.Fatal("stream channel closed")
		}
		if ev.Kind == rosrun.KindError {
			t.Fatalf("stream error: %v", ev.Err)
		}
		if ev.Kind != rosrun.KindRow && ev.Kind != rosrun.KindSnapshot {
			t.Fatalf("unexpected kind %s", ev.Kind)
		}
	case <-time.After(8 * time.Second):
		t.Fatalf("no stream event for interface %s", name)
	}
}

func TestLiveSessionPoll(t *testing.T) {
	sess, ctx := openLive(t)

	active, ok := ros.Lookup("hotspot-active")
	if !ok {
		t.Fatal("resource hotspot-active tidak ada")
	}
	req, err := active.Build("list", ros.AnyVersion, ros.Poll(time.Second), ros.Params{})
	if err != nil {
		t.Fatalf("build: %v", err)
	}

	ch, err := sess.Watch(ctx, req)
	if err != nil {
		t.Fatalf("watch: %v", err)
	}

	select {
	case ev, ok := <-ch:
		if !ok {
			t.Fatal("poll channel closed")
		}
		if ev.Kind != rosrun.KindSnapshot {
			t.Fatalf("kind = %s, want snapshot", ev.Kind)
		}
		t.Logf("poll snapshot: %d sesi aktif", len(ev.Rows))
	case <-time.After(8 * time.Second):
		t.Fatal("no snapshot from poll")
	}
}

func TestLiveSessionListen(t *testing.T) {
	sess, ctx := openLive(t)

	al, ok := ros.Lookup("fw-address-list")
	if !ok {
		t.Fatal("resource fw-address-list tidak ada")
	}
	req, err := al.Build("events", ros.AnyVersion, ros.Stream(), ros.Params{})
	if err != nil {
		t.Fatalf("build: %v", err)
	}

	ch, err := sess.Watch(ctx, req)
	if err != nil {
		t.Fatalf("watch: %v", err)
	}

	select {
	case ev, ok := <-ch:
		if !ok {
			t.Fatal("listen channel closed")
		}
		if ev.Kind == rosrun.KindError {
			t.Fatalf("listen error: %v", ev.Err)
		}
		t.Logf("listen event: %s", ev.Kind)
	case <-time.After(2 * time.Second):
		// Tidak ada perubahan address-list selama uji: bukan kegagalan.
		t.Log("listen aktif tanpa event (tidak ada perubahan address-list)")
	}
}
