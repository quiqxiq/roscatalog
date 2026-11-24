package rosrun

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	ros "example.com/roscatalog"
)

func openMock(t *testing.T, m *mock, mutate func(*Config)) *MikroTik {
	t.Helper()

	cfg := Config{
		ID:       "test",
		Address:  m.addr(),
		Username: "admin",
		Timeout:  2 * time.Second,
		backoff:  25 * time.Millisecond,
		Status:   StatusConfig{Disable: true, SuppressTLSWarning: true},
	}
	if mutate != nil {
		mutate(&cfg)
	}

	s, err := Open(context.Background(), cfg)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

func testContext(t *testing.T, d time.Duration) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), d)
	t.Cleanup(cancel)
	return ctx
}

func waitEvent(t *testing.T, ch <-chan Event, d time.Duration) Event {
	t.Helper()
	select {
	case ev, ok := <-ch:
		if !ok {
			t.Fatalf("channel closed before event")
		}
		return ev
	case <-time.After(d):
		t.Fatalf("no event within %s", d)
	}
	return Event{}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		ok   bool
	}{
		{"minimal", Config{Address: "192.168.1.1", Username: "admin"}, true},
		{"host port", Config{Address: "192.168.1.1:8729", Username: "admin"}, true},
		{"no address", Config{Username: "admin"}, false},
		{"no username", Config{Address: "192.168.1.1"}, false},
		{"port too big", Config{Address: "h", Username: "u", Port: 70000}, false},
		{"negative timeout", Config{Address: "h", Username: "u", Timeout: -time.Second}, false},
		{"poll too fast", Config{Address: "h", Username: "u", PollEvery: 10 * time.Millisecond}, false},
		{"poll too slow", Config{Address: "h", Username: "u", PollEvery: 31 * time.Second}, false},
		{"poll valid", Config{Address: "h", Username: "u", PollEvery: 5 * time.Second}, true},
		{"stream too fast", Config{Address: "h", Username: "u", StreamEvery: 500 * time.Millisecond}, false},
		{"stream too slow", Config{Address: "h", Username: "u", StreamEvery: 6 * time.Second}, false},
		{"status too fast", Config{Address: "h", Username: "u", Status: StatusConfig{Every: 500 * time.Millisecond}}, false},
		{"status too slow", Config{Address: "h", Username: "u", Status: StatusConfig{Every: 6 * time.Second}}, false},
		{"status timeout below interval", Config{Address: "h", Username: "u", Status: StatusConfig{Every: 2 * time.Second, Timeout: time.Second}}, false},
		{"negative cache ttl", Config{Address: "h", Username: "u", Cache: CacheConfig{TTL: -time.Second}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.ok && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.ok && err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestConfigAddr(t *testing.T) {
	tests := []struct {
		cfg  Config
		want string
	}{
		{Config{Address: "192.168.1.1"}, "192.168.1.1:8728"},
		{Config{Address: "192.168.1.1", Port: 9000}, "192.168.1.1:9000"},
		{Config{Address: "192.168.1.1:8729"}, "192.168.1.1:8729"},
		{Config{Address: "192.168.1.1", TLS: true}, "192.168.1.1:8729"},
		{Config{Address: "::1"}, "[::1]:8728"},
	}
	for _, tt := range tests {
		if got := tt.cfg.addr(); got != tt.want {
			t.Errorf("addr(%+v) = %q, want %q", tt.cfg, got, tt.want)
		}
	}
}

func TestOpenUsesOneConnection(t *testing.T) {
	m := newMock(t)
	s := openMock(t, m, nil)

	if got := m.acceptedCount(); got != 1 {
		t.Fatalf("accepted connections = %d, want 1", got)
	}
	st := s.Stats()
	if !st.Open {
		t.Fatal("session not open")
	}
	if st.Down {
		t.Fatalf("connection down at open: %+v", st)
	}
	if st.Epoch != 0 {
		t.Fatalf("epoch = %d, want 0", st.Epoch)
	}
}

func TestExecReturnsRet(t *testing.T) {
	m := newMock(t)
	s := openMock(t, m, nil)

	id, err := s.Exec(testContext(t, 2*time.Second), []string{"/ip/hotspot/user/add", "=name=vc1"})
	if err != nil {
		t.Fatalf("exec: %v", err)
	}
	if id != "*A1" {
		t.Fatalf("ret = %q, want *A1", id)
	}
}

func TestOnceUsesSessionConnection(t *testing.T) {
	m := newMock(t)
	s := openMock(t, m, nil)

	res, ok := ros.Lookup("hotspot-active")
	if !ok {
		t.Fatal("resource hotspot-active tidak ada")
	}
	req, err := res.Build("list", ros.Once(), ros.Params{})
	if err != nil {
		t.Fatalf("build: %v", err)
	}

	rows, err := s.Once(testContext(t, 2*time.Second), req)
	if err != nil {
		t.Fatalf("once: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("rows = %d, want 2", len(rows))
	}
}

func TestRequestableCmds(t *testing.T) {
	m := newMock(t)
	s := openMock(t, m, nil)

	// Test OnceCmd dengan ros.Path
	rows, err := s.OnceCmd(testContext(t, 2*time.Second), ros.Path("/ip/hotspot/active").Print())
	if err != nil {
		t.Fatalf("OnceCmd: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("rows = %d, want 2", len(rows))
	}

	// Test ExecCmd dengan ros.Path
	ret, err := s.ExecCmd(testContext(t, 2*time.Second), ros.Path("/ip/hotspot/user").Add(map[string]any{
		"name": "vc1",
	}))
	if err != nil {
		t.Fatalf("ExecCmd: %v", err)
	}
	if ret != "*A1" {
		t.Fatalf("ret = %q, want *A1", ret)
	}
}

func TestWatchPoll(t *testing.T) {
	m := newMock(t)
	s := openMock(t, m, nil)

	req := fastPoll(t, "hotspot-active", "list", 60*time.Millisecond, ros.Params{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, err := s.Watch(ctx, req)
	if err != nil {
		t.Fatalf("watch: %v", err)
	}

	ev := waitEvent(t, ch, 2*time.Second)
	if ev.Kind != KindSnapshot {
		t.Fatalf("kind = %s, want snapshot", ev.Kind)
	}
	if len(ev.Rows) == 0 {
		t.Fatal("snapshot kosong")
	}
}

func TestCloseStopsSession(t *testing.T) {
	m := newMock(t)
	s := openMock(t, m, nil)

	if err := s.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	if err := s.Close(); err != nil {
		t.Fatalf("second close: %v", err)
	}

	if _, err := s.Exec(context.Background(), []string{"/ip/hotspot/user/add"}); !errors.Is(err, ErrClosed) {
		t.Fatalf("exec after close = %v, want ErrClosed", err)
	}
}

func TestGateCalledPerCommand(t *testing.T) {
	m := newMock(t)
	var calls atomic.Int32

	s := openMock(t, m, func(c *Config) {
		c.Gate = func(ctx context.Context, path string) (func(), error) {
			calls.Add(1)
			return func() {}, nil
		}
	})

	ctx := testContext(t, 2*time.Second)

	if _, err := s.Exec(ctx, []string{"/ip/hotspot/user/add", "=name=x"}); err != nil {
		t.Fatalf("exec: %v", err)
	}

	res, _ := ros.Lookup("hotspot-active")
	req, _ := res.Build("list", ros.Once(), ros.Params{})
	if _, err := s.Once(ctx, req); err != nil {
		t.Fatalf("once: %v", err)
	}

	if got := calls.Load(); got != 2 {
		t.Fatalf("gate calls = %d, want 2", got)
	}
}

func TestConfigDefaults(t *testing.T) {
	c := Config{Address: "h", Username: "u"}.withDefaults()
	if c.PollEvery != 5*time.Second {
		t.Errorf("PollEvery = %s", c.PollEvery)
	}
	if c.Status.Target != "8.8.8.8" || c.Status.Every != time.Second || c.Status.Timeout != 15*time.Second {
		t.Errorf("status defaults = %+v", c.Status)
	}
	if c.Cache.TTL != 5*time.Second || c.Cache.MaxEntries != 64 {
		t.Errorf("cache defaults = %+v", c.Cache)
	}
	if c.Status.WarnPacketLoss != 5 {
		t.Errorf("WarnPacketLoss = %v", c.Status.WarnPacketLoss)
	}
}

func TestPlanFor(t *testing.T) {
	m := newMock(t)
	stream := openMock(t, m, func(c *Config) {
		c.Stream = true
		c.StreamEvery = 2 * time.Second
		c.PollEvery = 5 * time.Second
	})
	poll := openMock(t, m, func(c *Config) { c.Poll = true; c.PollEvery = 5 * time.Second })

	live, _ := ros.Lookup("hotspot-active")

	tests := []struct {
		name  string
		sess  *MikroTik
		view  string
		want  ros.Delivery
		every time.Duration
	}{
		{"interval stream", stream, "live", ros.DeliveryStream, 2 * time.Second},
		{"interval poll", poll, "live", ros.DeliveryPoll, 5 * time.Second},
		{"oneshot stream falls back", stream, "list", ros.DeliveryOnce, 0},
		{"oneshot poll", poll, "list", ros.DeliveryPoll, 5 * time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan := tt.sess.PlanFor(live, tt.view)
			if plan.Delivery != tt.want {
				t.Fatalf("delivery = %s, want %s", plan.Delivery, tt.want)
			}
			if plan.Every != tt.every {
				t.Fatalf("every = %s, want %s", plan.Every, tt.every)
			}
		})
	}
}
