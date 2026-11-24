package rosrun

import (
	"context"
	"strings"
	"testing"
	"time"

	ros "example.com/roscatalog"
)

func managerConfig(id string, m *mock) Config {
	return Config{
		ID:       id,
		Address:  m.addr(),
		Username: "admin",
		Timeout:  2 * time.Second,
		Status:   StatusConfig{Disable: true, SuppressTLSWarning: true},
		backoff:  25 * time.Millisecond,
	}
}

func TestOpenAllThreeRouters(t *testing.T) {
	m1, m2, m3 := newMock(t), newMock(t), newMock(t)
	mgr, err := OpenAll(context.Background(), []Config{
		managerConfig("r1", m1),
		managerConfig("r2", m2),
		managerConfig("r3", m3),
	})
	if err != nil {
		t.Fatalf("OpenAll: %v", err)
	}
	t.Cleanup(func() { _ = mgr.Close() })

	if mgr.Len() != 3 {
		t.Fatalf("Len = %d, want 3", mgr.Len())
	}
	if got := strings.Join(mgr.IDs(), ","); got != "r1,r2,r3" {
		t.Errorf("IDs = %q", got)
	}
	for _, m := range []*mock{m1, m2, m3} {
		if got := m.acceptedCount(); got != 1 {
			t.Errorf("accepted = %d, want 1", got)
		}
	}
	sts := mgr.Statuses()
	if len(sts) != 3 {
		t.Fatalf("Statuses = %d, want 3", len(sts))
	}
	for _, st := range sts {
		if !st.Connected || st.State != StatusConnected {
			t.Errorf("status %s: %+v", st.ID, st)
		}
	}
	if _, ok := mgr.Get("r2"); !ok {
		t.Error("Get r2 gagal")
	}
}

func TestOpenAllPartialFailure(t *testing.T) {
	good := newMock(t)
	mgr, err := OpenAll(context.Background(), []Config{
		managerConfig("good", good),
		{
			ID: "bad", Address: "127.0.0.1:1", Username: "admin",
			Timeout: 500 * time.Millisecond,
			Status:  StatusConfig{Disable: true},
		},
	})
	if err == nil {
		t.Fatal("partial failure harus mengembalikan error")
	}
	if mgr == nil {
		t.Fatal("manager sukses tidak boleh nil saat partial failure")
	}
	t.Cleanup(func() { _ = mgr.Close() })
	if mgr.Len() != 1 || mgr.IDs()[0] != "good" {
		t.Fatalf("manager = %v", mgr.IDs())
	}
}

func TestManagerWatchAll(t *testing.T) {
	m1, m2 := newMock(t), newMock(t)
	mgr, err := OpenAll(context.Background(), []Config{
		managerConfig("r1", m1),
		managerConfig("r2", m2),
	})
	if err != nil {
		t.Fatalf("OpenAll: %v", err)
	}
	t.Cleanup(func() { _ = mgr.Close() })

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := ros.Request{
		Delivery: ros.DeliveryPoll,
		Sentence: []string{"/ip/hotspot/active/print"},
		Every:    40 * time.Millisecond,
		Shape:    ros.ShapeSnapshot,
	}
	ch, err := mgr.WatchAll(ctx, req)
	if err != nil {
		t.Fatalf("WatchAll: %v", err)
	}

	seen := map[string]bool{}
	deadline := time.After(3 * time.Second)
	for len(seen) < 2 {
		select {
		case ev, ok := <-ch:
			if !ok {
				t.Fatalf("channel tertutup, seen=%v", seen)
			}
			if ev.Kind == KindError {
				t.Fatalf("error dari %s: %v", ev.RouterID, ev.Err)
			}
			seen[ev.RouterID] = true
		case <-deadline:
			t.Fatalf("tidak menerima event dari semua router: %v", seen)
		}
	}
	cancel()
	for range ch {
	}
}

func TestManagerWatchStatuses(t *testing.T) {
	m1, m2 := newMock(t), newMock(t)
	c1, c2 := managerConfig("r1", m1), managerConfig("r2", m2)
	for _, c := range []*Config{&c1, &c2} {
		c.Status = StatusConfig{
			PingInterval:       time.Second,
			ReplyTimeout:       5 * time.Second,
			SuppressTLSWarning: true,
		}
	}
	mgr, err := OpenAll(context.Background(), []Config{c1, c2})
	if err != nil {
		t.Fatalf("OpenAll: %v", err)
	}
	t.Cleanup(func() { _ = mgr.Close() })

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch, err := mgr.WatchStatuses(ctx)
	if err != nil {
		t.Fatalf("WatchStatuses: %v", err)
	}

	seen := map[string]bool{}
	deadline := time.After(5 * time.Second)
	for len(seen) < 2 {
		select {
		case rs, ok := <-ch:
			if !ok {
				t.Fatalf("channel status tertutup, seen=%v", seen)
			}
			if rs.Status.Connected && rs.Status.Ping.Received > 0 {
				seen[rs.RouterID] = true
			}
		case <-deadline:
			t.Fatalf("status belum connected: %v", seen)
		}
	}
	cancel()
	for range ch {
	}
}

func TestManagerWatchAutoAllUsesPerRouterStream(t *testing.T) {
	m1, m2 := newMock(t), newMock(t)
	cfg1 := managerConfig("stream", m1)
	cfg1.Stream = true
	cfg1.StreamEvery = time.Second
	cfg2 := managerConfig("poll", m2)
	cfg2.Poll = true

	mgr, err := OpenAll(context.Background(), []Config{cfg1, cfg2})
	if err != nil {
		t.Fatalf("OpenAll: %v", err)
	}
	t.Cleanup(func() { _ = mgr.Close() })

	res, _ := ros.Lookup("hotspot-active")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch, err := mgr.WatchAutoAll(ctx, res, "live", ros.Plan{}, ros.Params{})
	if err != nil {
		t.Fatalf("WatchAutoAll: %v", err)
	}

	seen := map[string]bool{}
	deadline := time.After(4 * time.Second)
	for len(seen) < 2 {
		select {
		case ev, ok := <-ch:
			if !ok {
				t.Fatalf("channel tertutup, seen=%v", seen)
			}
			if ev.Kind == KindError {
				t.Fatalf("error dari %s: %v", ev.RouterID, ev.Err)
			}
			seen[ev.RouterID] = true
		case <-deadline:
			t.Fatalf("tidak menerima event dari semua router: %v", seen)
		}
	}
	cancel()
	for range ch {
	}

	c1 := m1.callsTo("/ip/hotspot/active/print")
	c2 := m2.callsTo("/ip/hotspot/active/print")
	if len(c1) == 0 || !c1[0].has("interval") {
		t.Errorf("router stream tidak memakai interval: %+v", c1)
	}
	if len(c2) == 0 || c2[0].has("interval") {
		t.Errorf("router poll tidak boleh memakai interval: %+v", c2)
	}
}

func TestManagerIdleKeepsRegistration(t *testing.T) {
	m := newMock(t)
	mgr := NewManager()
	t.Cleanup(func() { _ = mgr.Close() })

	cfg := managerConfig("r-idle", m)
	cfg.Grace = 40 * time.Millisecond
	mgr.Register(cfg)

	if _, err := mgr.Acquire(context.Background(), "r-idle"); err != nil {
		t.Fatalf("Acquire pertama: %v", err)
	}
	mgr.Release("r-idle")

	waitFor(t, 2*time.Second, func() bool {
		_, open := mgr.Get("r-idle")
		return !open
	}, "sesi ditutup setelah grace")

	if !mgr.Has("r-idle") {
		t.Fatal("registrasi hilang setelah sesi idle")
	}
	if got := mgr.Configs(); len(got) != 1 || got[0].ID != "r-idle" {
		t.Fatalf("Configs = %+v", got)
	}

	sess, err := mgr.Acquire(context.Background(), "r-idle")
	if err != nil {
		t.Fatalf("Acquire setelah idle: %v", err)
	}
	if sess == nil {
		t.Fatal("sesi nil")
	}
	mgr.Release("r-idle")
}

func TestManagerGraceIsPerRouter(t *testing.T) {
	m1, m2 := newMock(t), newMock(t)
	c1 := managerConfig("fast", m1)
	c1.Grace = 50 * time.Millisecond
	c2 := managerConfig("slow", m2)
	c2.Grace = 2 * time.Second

	mgr, err := OpenAll(context.Background(), []Config{c1, c2})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = mgr.Close() })

	for _, id := range []string{"fast", "slow"} {
		if _, err := mgr.Acquire(context.Background(), id); err != nil {
			t.Fatalf("Acquire %s: %v", id, err)
		}
		mgr.Release(id)
	}

	waitFor(t, time.Second, func() bool {
		_, open := mgr.Get("fast")
		return !open
	}, "router fast ditutup sesuai Grace-nya")
	if _, open := mgr.Get("slow"); !open {
		t.Fatal("router slow ikut ditutup; Grace harus per-router")
	}
}

func TestManagerWatchAllOpensRegisteredRouter(t *testing.T) {
	m := newMock(t)
	mgr := NewManager()
	t.Cleanup(func() { _ = mgr.Close() })

	cfg := managerConfig("r-watch", m)
	cfg.Grace = 50 * time.Millisecond
	mgr.Register(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req := ros.Request{
		Delivery: ros.DeliveryPoll,
		Sentence: []string{"/ip/hotspot/active/print"},
		Every:    40 * time.Millisecond,
		Shape:    ros.ShapeSnapshot,
	}
	ch, err := mgr.WatchAll(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	select {
	case ev, ok := <-ch:
		if !ok {
			t.Fatal("channel tertutup")
		}
		if ev.Kind == KindError {
			t.Fatalf("error dari %s: %v", ev.RouterID, ev.Err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("tidak ada event dari router terdaftar")
	}

	// Watch memegang refcount, jadi grace tidak boleh menutup sesi.
	time.Sleep(120 * time.Millisecond)
	if _, open := mgr.Get("r-watch"); !open {
		t.Fatal("sesi ditutup saat masih di-watch")
	}
}

func TestManagerRegisterAndUnregister(t *testing.T) {
	m := newMock(t)
	mgr := NewManager()
	t.Cleanup(func() { _ = mgr.Close() })

	cfg := managerConfig("r-dyn", m)
	mgr.Register(cfg)

	// Acquire opens session on-demand from registered config
	sess, err := mgr.Acquire(context.Background(), "r-dyn")
	if err != nil {
		t.Fatalf("Acquire: %v", err)
	}
	if sess == nil {
		t.Fatal("sess is nil")
	}
	mgr.Release("r-dyn")

	// Unregister closes and removes
	if err := mgr.Unregister("r-dyn"); err != nil {
		t.Fatalf("Unregister: %v", err)
	}
	if _, ok := mgr.Get("r-dyn"); ok {
		t.Error("r-dyn still exists after Unregister")
	}
}
