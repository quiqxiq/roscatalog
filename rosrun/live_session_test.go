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
		PollEvery: 5 * time.Second,
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
	if !st.Open || st.Down {
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
	req, err := ifs.Build("list", ros.Once(), ros.Params{})
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
	if got := sess.PlanFor(mt, "run").Delivery; got != ros.DeliveryStream {
		t.Fatalf("default delivery = %q, want stream", got)
	}
	mreq, err := mt.Build("run", ros.StreamEvery(time.Second), ros.Params{
		Args: map[string]string{"interface": name},
	})
	if err != nil {
		t.Fatalf("build monitor: %v", err)
	}

	ch, err := sess.Watch(ctx, mreq)
	if err != nil {
		t.Fatalf("watch monitor: %v", err)
	}

	// Command harus tetap jalan saat stream terbuka di satu koneksi async.
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
	req, err := active.Build("list", ros.Poll(5*time.Second), ros.Params{})
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

func TestLiveSessionUnionProplist(t *testing.T) {
	sess, ctx := openLive(t)

	res, ok := ros.Lookup("interface")
	if !ok {
		t.Fatal("resource interface tidak ada")
	}
	build := func(only ...string) ros.Request {
		t.Helper()
		req, err := res.Build("list", ros.Poll(5*time.Second), ros.Params{Only: only})
		if err != nil {
			t.Fatalf("build %v: %v", only, err)
		}
		return req
	}

	subA, err := sess.Subscribe(ctx, build("name"))
	if err != nil {
		t.Fatalf("subscribe A: %v", err)
	}
	defer subA.Close()
	subB, err := sess.Subscribe(ctx, build("type"))
	if err != nil {
		t.Fatalf("subscribe B: %v", err)
	}
	defer subB.Close()

	wait := func(sub *rosrun.Subscription, label string) {
		t.Helper()
		for {
			select {
			case ev, ok := <-sub.Events():
				if !ok {
					t.Fatalf("%s: channel closed", label)
				}
				if ev.Kind != rosrun.KindSnapshot || len(ev.Rows) == 0 {
					continue
				}
				for _, row := range ev.Rows {
					if label == "A" {
						if _, has := row["type"]; has {
							t.Fatalf("A menerima type: %v", row)
						}
					} else {
						if _, has := row["name"]; has {
							t.Fatalf("B menerima name: %v", row)
						}
					}
				}
				return
			case <-time.After(8 * time.Second):
				t.Fatalf("%s: tidak ada snapshot", label)
			}
		}
	}
	wait(subA, "A")
	wait(subB, "B")

	st := sess.Stats()
	if st.Commands == 0 {
		t.Fatal("statistik command tidak tercatat")
	}
	t.Logf("union ok: commands=%d inflight=%d", st.Commands, st.InFlight)
}

func TestLiveSessionListen(t *testing.T) {
	sess, ctx := openLive(t)

	al, ok := ros.Lookup("fw-address-list")
	if !ok {
		t.Fatal("resource fw-address-list tidak ada")
	}
	req, err := al.Build("events", ros.Events(), ros.Params{})
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
