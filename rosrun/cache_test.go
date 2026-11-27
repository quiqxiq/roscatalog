package rosrun

import (
	"context"
	"testing"
	"time"

	ros "github.com/quiqxiq/roscatalog"
)

func monitorReq(iface string) ros.Request {
	return ros.Request{
		Delivery: ros.DeliveryStream,
		Mode:     ros.ModeNative,
		Shape:    ros.ShapeRows,
		Sentence: []string{"/interface/monitor-traffic", "=interface=" + iface},
	}
}

func nextSub(t *testing.T, sub *Subscription, d time.Duration) Event {
	t.Helper()
	select {
	case ev, ok := <-sub.Events():
		if !ok {
			t.Fatal("subscription tertutup")
		}
		return ev
	case <-time.After(d):
		t.Fatal("timeout menunggu event subscription")
	}
	return Event{}
}

func TestSubscribeSharesUpstream(t *testing.T) {
	m := newMock(t)
	s := openMock(t, m, nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := monitorReq("ether1")
	sub1, err := s.Subscribe(ctx, req)
	if err != nil {
		t.Fatalf("subscribe 1: %v", err)
	}
	sub2, err := s.Subscribe(ctx, req)
	if err != nil {
		t.Fatalf("subscribe 2: %v", err)
	}
	defer sub1.Close()
	defer sub2.Close()

	ev1 := nextSub(t, sub1, 2*time.Second)
	ev2 := nextSub(t, sub2, 2*time.Second)
	if ev1.Kind != KindRow || ev2.Kind != KindRow {
		t.Fatalf("event: %+v %+v", ev1, ev2)
	}
	if got := len(m.callsTo("/interface/monitor-traffic")); got != 1 {
		t.Fatalf("upstream dipanggil %d kali, want 1", got)
	}
	if _, ok := sub2.Latest(); !ok {
		t.Error("Latest kosong")
	}
}

func TestSubscribeLateGetsLatest(t *testing.T) {
	m := newMock(t)
	s := openMock(t, m, nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := monitorReq("ether1")
	sub1, err := s.Subscribe(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	defer sub1.Close()
	nextSub(t, sub1, 2*time.Second)

	sub2, err := s.Subscribe(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	defer sub2.Close()

	if ev, ok := sub2.Latest(); !ok || ev.Kind != KindRow {
		t.Fatalf("subscriber baru tidak mendapat latest: %+v ok=%v", ev, ok)
	}
	nextSub(t, sub2, 500*time.Millisecond)
	if got := len(m.callsTo("/interface/monitor-traffic")); got != 1 {
		t.Fatalf("late subscriber membuat upstream baru: %d", got)
	}
}

func TestSubscribeCachesAfterLastLeave(t *testing.T) {
	m := newMock(t)
	s := openMock(t, m, nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := monitorReq("ether1")
	sub1, err := s.Subscribe(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	nextSub(t, sub1, 2*time.Second)
	sub1.Close()
	if !m.waitStreamsStopped(2 * time.Second) {
		t.Fatal("upstream tidak berhenti setelah subscriber terakhir")
	}

	sub2, err := s.Subscribe(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	defer sub2.Close()
	if _, ok := sub2.Latest(); !ok {
		t.Error("snapshot cache tidak tersedia")
	}
	nextSub(t, sub2, 2*time.Second)
	waitFor(t, 2*time.Second, func() bool {
		return len(m.callsTo("/interface/monitor-traffic")) >= 2
	}, "upstream baru setelah cache hit")
}

func TestSubscribeIsolatesDifferentRequests(t *testing.T) {
	m := newMock(t)
	s := openMock(t, m, nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sub1, err := s.Subscribe(ctx, monitorReq("ether1"))
	if err != nil {
		t.Fatal(err)
	}
	defer sub1.Close()
	sub2, err := s.Subscribe(ctx, monitorReq("ether2"))
	if err != nil {
		t.Fatal(err)
	}
	defer sub2.Close()

	nextSub(t, sub1, 2*time.Second)
	nextSub(t, sub2, 2*time.Second)
	if got := len(m.callsTo("/interface/monitor-traffic")); got != 2 {
		t.Fatalf("request berbeda berbagi upstream: %d", got)
	}
}

func TestSessionCloseClosesSubscriptions(t *testing.T) {
	m := newMock(t)
	s := openMock(t, m, nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sub, err := s.Subscribe(ctx, monitorReq("ether1"))
	if err != nil {
		t.Fatal(err)
	}
	nextSub(t, sub, 2*time.Second)

	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	deadline := time.After(3 * time.Second)
	for {
		select {
		case _, ok := <-sub.Events():
			if !ok {
				return
			}
		case <-deadline:
			t.Fatal("channel subscription tidak tertutup setelah Close")
		}
	}
}

func TestSubscribeRejectsOnce(t *testing.T) {
	m := newMock(t)
	s := openMock(t, m, nil)
	if _, err := s.Subscribe(context.Background(), ros.Request{
		Delivery: ros.DeliveryOnce,
		Sentence: []string{"/x"},
	}); err == nil {
		t.Fatal("Subscribe dengan Delivery once seharusnya ditolak")
	}
}

func TestCacheRespectsNoServerProplist(t *testing.T) {
	m := newMock(t)
	s := openMock(t, m, nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	res, ok := ros.Lookup("system-resource")
	if !ok {
		t.Fatal("system-resource missing")
	}
	req, err := res.Build("list", ros.Poll(5*time.Second), ros.Params{})
	if err != nil {
		t.Fatal(err)
	}
	if !req.NoServerProplist {
		t.Fatal("expected NoServerProplist=true for system-resource/list")
	}
	sub, err := s.Subscribe(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Close()
	nextSub(t, sub, 2*time.Second)

	for _, c := range m.callsTo("/system/resource/print") {
		if c.has(".proplist") {
			t.Fatalf("cache injected .proplist into NoServerProplist request: %+v", c.attrs)
		}
	}
}

func TestCachePollIntervalKeySeparatesFeeds(t *testing.T) {
	req5s := ros.Request{
		ResourceID: "interface", View: "list",
		Delivery: ros.DeliveryPoll, Mode: ros.ModeOneShot,
		Every: 5 * time.Second, Shape: ros.ShapeSnapshot,
		Sentence: []string{"/interface/print"},
	}
	req30s := req5s
	req30s.Every = 30 * time.Second
	if menuKey(req5s) == menuKey(req30s) {
		t.Fatal("5s and 30s poll requests must not share the same menuKey")
	}
}
