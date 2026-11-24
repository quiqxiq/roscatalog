package rosrun

import (
	"context"
	"strings"
	"testing"
	"time"

	ros "example.com/roscatalog"
)

func TestPlanForDefaultsToStream(t *testing.T) {
	m := newMock(t)
	s := openMock(t, m, nil)
	res, ok := ros.Lookup("hotspot-active")
	if !ok {
		t.Fatal("resource hotspot-active tidak ada")
	}
	if got := s.PlanFor(res, "live").Delivery; got != ros.DeliveryStream {
		t.Fatalf("PlanFor default = %q, want stream", got)
	}
	if got := s.PlanFor(res, "events").Delivery; got != ros.DeliveryEvents {
		t.Fatalf("PlanFor events = %q, want events", got)
	}

	m2 := newMock(t)
	s2 := openMock(t, m2, func(c *Config) { c.Poll = true })
	if got := s2.PlanFor(res, "live").Delivery; got != ros.DeliveryPoll {
		t.Fatalf("PlanFor Poll=true = %q, want poll", got)
	}
}

func TestSingleConnectionServesStreamAndExec(t *testing.T) {
	m := newMock(t)
	s := openMock(t, m, nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := fastInterval(t, "hotspot-active", "live", 50*time.Millisecond, ros.Params{})
	ch, err := s.Watch(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	if ev, ok := next(t, ch, 2*time.Second); !ok || ev.Kind != KindSnapshot {
		t.Fatalf("event stream pertama: %+v (ok=%v)", ev, ok)
	}

	if _, err := s.Exec(testContext(t, 2*time.Second), []string{"/ip/hotspot/user/add", "=name=single"}); err != nil {
		t.Fatalf("exec saat stream terbuka: %v", err)
	}
	if got := m.acceptedCount(); got != 1 {
		t.Fatalf("accepted connections = %d, want 1", got)
	}
}

func TestLimiterCapAndStats(t *testing.T) {
	l := newLimiter(1)
	release := l.acquire()
	st := l.stats()
	if st.InFlight != 1 {
		t.Fatalf("InFlight = %d, want 1", st.InFlight)
	}

	done := make(chan struct{})
	go func() {
		r := l.acquire()
		close(done)
		r()
	}()
	select {
	case <-done:
		t.Fatal("acquire kedua tidak boleh langsung lolos")
	case <-time.After(30 * time.Millisecond):
	}
	release()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("acquire kedua tidak pernah lolos setelah release")
	}
	if st := l.stats(); st.InFlight != 0 {
		t.Fatalf("InFlight = %d, want 0", st.InFlight)
	}

	l.note("/ip/address/print")
	closeStream := l.streamOpened()
	if st := l.stats(); st.OpenStreams != 1 || st.Commands != 1 {
		t.Fatalf("statistik salah: %+v", st)
	}
	closeStream()
	l.streamsGone()
	if st := l.stats(); st.OpenStreams != 0 {
		t.Fatalf("OpenStreams = %d, want 0", st.OpenStreams)
	}
}

func TestUnionProplistSharesOneUpstream(t *testing.T) {
	m := newMock(t)
	s := openMock(t, m, func(c *Config) { c.Logf = t.Logf })
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reqA := fastPoll(t, "hotspot-active", "list", 10*time.Millisecond,
		ros.Params{Only: []string{"address"}})
	reqB := fastPoll(t, "hotspot-active", "list", 10*time.Millisecond,
		ros.Params{Only: []string{"user"}})

	subA, err := s.Subscribe(ctx, reqA)
	if err != nil {
		t.Fatal(err)
	}
	defer subA.Close()
	subB, err := s.Subscribe(ctx, reqB)
	if err != nil {
		t.Fatal(err)
	}
	defer subB.Close()

	waitFor(t, 2*time.Second, func() bool {
		for _, c := range m.callsTo("/ip/hotspot/active/print") {
			pl := c.attrs[".proplist"]
			if strings.Contains(pl, "address") && strings.Contains(pl, "user") {
				return true
			}
		}
		return false
	}, "proplist union address+user")

	if got := len(m.callsTo("/ip/hotspot/active/print")); got != 1 {
		t.Fatalf("upstream print = %d, want 1 (dua subscriber berbagi satu menu)", got)
	}

	var ev Event
	var ok bool
	select {
	case ev, ok = <-subA.Events():
	case <-time.After(2 * time.Second):
	}
	if !ok {
		s.cache.mu.Lock()
		t.Logf("entries=%d", len(s.cache.entries))
		for k, e := range s.cache.entries {
			t.Logf("entry %q subs=%d stop=%v poll=%v req=%v", k, len(e.subs), e.stop != nil, e.poll != nil, e.req.Sentence)
		}
		s.cache.mu.Unlock()
		t.Fatalf("event pertama subA tidak datang")
	}
	if ev.Kind != KindSnapshot || len(ev.Rows) == 0 {
		t.Fatalf("event pertama subA: %+v", ev)
	}
	for _, row := range ev.Rows {
		if _, has := row["user"]; has {
			t.Fatalf("subA tidak boleh menerima user: %v", row)
		}
		if row[".id"] == "" || row["address"] == "" {
			t.Fatalf("key .id dan address harus ada: %v", row)
		}
	}
}

func TestListenNeverSendsProplistWithUnion(t *testing.T) {
	m := newMock(t)
	s := openMock(t, m, nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reqA := build(t, "hotspot-active", "events", ros.Stream(),
		ros.Params{Only: []string{"user"}})
	reqB := build(t, "hotspot-active", "events", ros.Stream(),
		ros.Params{Only: []string{"address"}})

	subA, err := s.Subscribe(ctx, reqA)
	if err != nil {
		t.Fatal(err)
	}
	defer subA.Close()
	subB, err := s.Subscribe(ctx, reqB)
	if err != nil {
		t.Fatal(err)
	}
	defer subB.Close()

	waitFor(t, 2*time.Second, func() bool {
		return len(m.callsTo("/ip/hotspot/active/listen")) > 0
	}, "listen dibuka")
	for _, c := range m.callsTo("/ip/hotspot/active/listen") {
		if c.has(".proplist") {
			t.Fatalf("listen tidak boleh mengirim .proplist: %+v", c.attrs)
		}
	}
}

func TestListenIsNotReopenedOnSilence(t *testing.T) {
	m := newMock(t)
	s := openMock(t, m, func(c *Config) {
		c.StreamCheck = 50 * time.Millisecond
		c.StreamStale = 100 * time.Millisecond
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := build(t, "hotspot-active", "events", ros.Events(), ros.Params{})
	ch, err := s.Watch(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for range ch {
		}
	}()

	time.Sleep(400 * time.Millisecond)
	if got := len(m.callsTo("/ip/hotspot/active/listen")); got != 1 {
		t.Fatalf("listen dibuka %d kali; diam tidak boleh memicu watchdog", got)
	}
}

func TestManagerLifecycleHoldAndGrace(t *testing.T) {
	m := newMock(t)
	cfg := managerConfig("r1", m)
	cfg.Grace = 60 * time.Millisecond
	mgr, err := OpenAll(context.Background(), []Config{cfg})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = mgr.Close() })

	if _, err := mgr.Acquire(context.Background(), "r1"); err != nil {
		t.Fatal(err)
	}
	if err := mgr.Retain("r1", "alerts"); err != nil {
		t.Fatal(err)
	}
	mgr.Release("r1")

	time.Sleep(120 * time.Millisecond)
	if _, ok := mgr.Get("r1"); !ok {
		t.Fatal("sesi ditutup padahal masih di-hold")
	}

	mgr.Drop("r1", "alerts")
	waitFor(t, 2*time.Second, func() bool {
		_, ok := mgr.Get("r1")
		return !ok
	}, "sesi ditutup setelah grace")
}

func TestSessionStatsIncludeLimiter(t *testing.T) {
	m := newMock(t)
	s := openMock(t, m, nil)
	if _, err := s.Exec(testContext(t, 2*time.Second), []string{"/ip/hotspot/user/add", "=name=vc1"}); err != nil {
		t.Fatal(err)
	}
	st := s.Stats()
	if st.Commands < 1 {
		t.Fatalf("Commands = %d, want >= 1", st.Commands)
	}
	if st.InFlight != 0 {
		t.Fatalf("InFlight = %d, want 0", st.InFlight)
	}
}
