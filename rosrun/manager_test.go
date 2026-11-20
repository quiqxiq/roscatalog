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
		if got := m.acceptedCount(); got != 2 {
			t.Errorf("accepted = %d, want 2", got)
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
			Every:              time.Second,
			Timeout:            5 * time.Second,
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
	cfg2.Stream = false

	mgr, err := OpenAll(context.Background(), []Config{cfg1, cfg2})
	if err != nil {
		t.Fatalf("OpenAll: %v", err)
	}
	t.Cleanup(func() { _ = mgr.Close() })

	res, _ := ros.Lookup("hotspot-active")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch, err := mgr.WatchAutoAll(ctx, res, "live", ros.AnyVersion, ros.Plan{}, ros.Params{})
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
