package rosrun

import (
	"context"
	"testing"
	"time"

	ros "example.com/roscatalog"
)

func waitFor(t *testing.T, d time.Duration, cond func() bool, msg string) {
	t.Helper()
	deadline := time.Now().Add(d)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("timeout waiting for %s", msg)
}

func TestReconnectCommand(t *testing.T) {
	m := newMock(t)
	s := openMock(t, m, nil)

	m.dropAll()

	waitFor(t, 8*time.Second, func() bool {
		st := s.Stats()
		return !st.CmdDown && !st.FeedDown && st.Reconnects >= 2
	}, "both connections to reconnect")

	if got := m.acceptedCount(); got < 4 {
		t.Fatalf("accepted connections = %d, want at least 4", got)
	}

	id, err := s.Exec(testContext(t, 3*time.Second), []string{"/ip/hotspot/user/add", "=name=vc1"})
	if err != nil {
		t.Fatalf("exec after reconnect: %v", err)
	}
	if id != "*A1" {
		t.Fatalf("ret = %q, want *A1", id)
	}
}

func TestReconnectResumesWatch(t *testing.T) {
	m := newMock(t)
	s := openMock(t, m, nil)

	req := fastInterval(t, "hotspot-active", "live", ros.AnyVersion, 300*time.Millisecond, ros.Params{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, err := s.Watch(ctx, req)
	if err != nil {
		t.Fatalf("watch: %v", err)
	}

	if ev := waitEvent(t, ch, 3*time.Second); ev.Kind != KindSnapshot {
		t.Fatalf("first event = %s, want snapshot", ev.Kind)
	}

	m.dropAll()

	sawReconnect := false
	deadline := time.NewTimer(10 * time.Second)
	defer deadline.Stop()
	for {
		select {
		case ev, ok := <-ch:
			if !ok {
				t.Fatal("watch channel closed before resuming")
			}
			switch ev.Kind {
			case KindReconnecting:
				sawReconnect = true
			case KindSnapshot:
				if sawReconnect {
					return
				}
			case KindError:
				t.Fatalf("unexpected error event: %v", ev.Err)
			}
		case <-deadline.C:
			t.Fatalf("watch did not resume (sawReconnect=%v, stats=%+v)", sawReconnect, s.Stats())
		}
	}
}

func TestCloseEndsWatch(t *testing.T) {
	m := newMock(t)
	s := openMock(t, m, nil)

	req := fastInterval(t, "hotspot-active", "live", ros.AnyVersion, 300*time.Millisecond, ros.Params{})

	ch, err := s.Watch(context.Background(), req)
	if err != nil {
		t.Fatalf("watch: %v", err)
	}
	if ev := waitEvent(t, ch, 3*time.Second); ev.Kind != KindSnapshot {
		t.Fatalf("first event = %s, want snapshot", ev.Kind)
	}

	if err := s.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	select {
	case _, ok := <-ch:
		if ok {
			// Buang sisa event yang sudah ter-buffer, lalu tunggu tutup.
			for range ch {
			}
		}
	case <-time.After(3 * time.Second):
		t.Fatal("watch channel not closed after session close")
	}
}
