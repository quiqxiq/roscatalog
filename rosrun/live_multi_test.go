package rosrun_test

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	ros "example.com/roscatalog"
	"example.com/roscatalog/rosrun"
)

// TestLiveMultiRouter membuka beberapa router sekaligus dan memastikan status,
// stream, poll, native, dan cache berjalan bersamaan.
//
//	ROUTEROS_LIVE_ADDRS=192.168.122.60:8728,192.168.122.26:8728,192.168.122.189:8728 \
//	ROUTEROS_LIVE_USER=admin ROUTEROS_LIVE_PASS=secret \
//	go test ./rosrun -run TestLiveMultiRouter -v -count=1
func TestLiveMultiRouter(t *testing.T) {
	raw := os.Getenv("ROUTEROS_LIVE_ADDRS")
	if raw == "" {
		t.Skip("set ROUTEROS_LIVE_ADDRS untuk uji multi-router")
	}
	user := os.Getenv("ROUTEROS_LIVE_USER")
	if user == "" {
		user = "admin"
	}
	pass := os.Getenv("ROUTEROS_LIVE_PASS")

	addrs := strings.Split(raw, ",")
	cfgs := make([]rosrun.Config, 0, len(addrs))
	for i, addr := range addrs {
		addr = strings.TrimSpace(addr)
		if addr == "" {
			continue
		}
		cfgs = append(cfgs, rosrun.Config{
			ID:          addr,
			Address:     addr,
			Username:    user,
			Password:    pass,
			Timeout:     10 * time.Second,
			Stream:      i%2 == 0,
			StreamEvery: time.Second,
			PollEvery:   5 * time.Second,
			Status:      rosrun.StatusConfig{Every: time.Second, Timeout: 10 * time.Second},
			Logf:        t.Logf,
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mgr, err := rosrun.OpenAll(ctx, cfgs)
	if err != nil {
		if mgr == nil {
			t.Fatalf("OpenAll: %v", err)
		}
		t.Fatalf("sebagian router gagal: %v", err)
	}
	defer mgr.Close()
	if mgr.Len() != len(cfgs) {
		t.Fatalf("Len = %d, want %d", mgr.Len(), len(cfgs))
	}

	// Status: native ping di semua router, terlepas dari Config.Stream.
	deadline := time.Now().Add(15 * time.Second)
	for {
		all := true
		for _, st := range mgr.Statuses() {
			if !st.Connected || st.Ping.Received == 0 {
				all = false
			}
		}
		if all {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("status belum connected: %+v", mgr.Statuses())
		}
		time.Sleep(200 * time.Millisecond)
	}
	for _, st := range mgr.Statuses() {
		if st.Target != "8.8.8.8" || st.Ping.Received == 0 {
			t.Errorf("status %s: %+v", st.ID, st)
		}
	}

	// Stream/poll otomatis per router.
	res, _ := ros.Lookup("interface")
	auto, err := mgr.WatchAutoAll(ctx, res, "live", ros.AnyVersion, ros.Plan{}, ros.Params{})
	if err != nil {
		t.Fatalf("WatchAutoAll: %v", err)
	}
	collectRouters(t, auto, len(cfgs), 15*time.Second, "WatchAutoAll")

	// Stream interval + stats eksplisit di semua router.
	sreq, err := res.Build("live", ros.AnyVersion, ros.StreamEvery(time.Second), ros.Params{})
	if err != nil {
		t.Fatalf("build stream: %v", err)
	}
	streams, err := mgr.WatchAll(ctx, sreq)
	if err != nil {
		t.Fatalf("WatchAll stream: %v", err)
	}
	collectRouters(t, streams, len(cfgs), 15*time.Second, "stream")

	// Poll eksplisit di semua router, berjalan bersamaan dengan stream.
	preq, err := res.Build("live", ros.AnyVersion, ros.Poll(5*time.Second), ros.Params{})
	if err != nil {
		t.Fatalf("build poll: %v", err)
	}
	polls, err := mgr.WatchAll(ctx, preq)
	if err != nil {
		t.Fatalf("WatchAll poll: %v", err)
	}
	collectRouters(t, polls, len(cfgs), 15*time.Second, "poll")

	// Native monitor-traffic.
	mt, _ := ros.Lookup("interface-monitor-traffic")
	mreq, err := mt.Build("run", ros.AnyVersion, ros.Stream(), ros.Params{
		Args: map[string]string{"interface": "ether1"},
	})
	if err != nil {
		t.Fatalf("build monitor: %v", err)
	}
	mon, err := mgr.WatchAll(ctx, mreq)
	if err != nil {
		t.Fatalf("WatchAll monitor: %v", err)
	}
	collectRouters(t, mon, len(cfgs), 15*time.Second, "monitor")

	// Cache/shared feed: satu subscription per router.
	subs, err := mgr.SubscribeAll(ctx, mreq)
	if err != nil {
		t.Logf("SubscribeAll partial: %v", err)
	}
	for id, sub := range subs {
		select {
		case ev := <-sub.Events():
			if ev.Kind == rosrun.KindError {
				t.Errorf("%s cache error: %v", id, ev.Err)
			}
		case <-time.After(10 * time.Second):
			t.Errorf("%s tidak menerima event cache", id)
		}
		sub.Close()
	}

	// Command tetap responsif saat semua stream berjalan.
	for _, sess := range mgr.Routers() {
		start := time.Now()
		if _, err := sess.Exec(ctx, []string{"/system/identity/print"}); err != nil {
			t.Errorf("%s exec: %v", sess.Stats().ID, err)
		}
		if d := time.Since(start); d > 2*time.Second {
			t.Errorf("%s exec lambat: %s", sess.Stats().ID, d)
		}
	}
}

func collectRouters(t *testing.T, ch <-chan rosrun.RouterEvent, n int, d time.Duration, name string) {
	t.Helper()
	seen := map[string]bool{}
	deadline := time.After(d)
	for len(seen) < n {
		select {
		case ev, ok := <-ch:
			if !ok {
				t.Fatalf("%s: channel tertutup, diterima %v", name, seen)
			}
			if ev.Kind == rosrun.KindError {
				t.Fatalf("%s %s: %v", name, ev.RouterID, ev.Err)
			}
			seen[ev.RouterID] = true
		case <-deadline:
			t.Fatalf("%s: timeout, diterima %v", name, seen)
		}
	}
}
