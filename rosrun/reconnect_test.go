package rosrun

import (
	"context"
	"os"
	"runtime"
	"testing"
	"time"

	ros "github.com/quiqxiq/roscatalog"
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
		return !st.Down && st.Reconnects >= 1
	}, "connection to reconnect")

	if got := m.acceptedCount(); got < 2 {
		t.Fatalf("accepted connections = %d, want at least 2", got)
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

	req := fastInterval(t, "hotspot-active", "live", 300*time.Millisecond, ros.Params{})

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

	req := fastInterval(t, "hotspot-active", "live", 300*time.Millisecond, ros.Params{})

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

func countOpenFDs(t *testing.T) int {
	t.Helper()
	entries, err := os.ReadDir("/proc/self/fd")
	if err != nil {
		return -1
	}
	return len(entries)
}

func TestReconnectNoGoroutineFDLeak(t *testing.T) {
	m := newMock(t)
	s := openMock(t, m, nil)
	defer s.Close()

	// Warm-up goroutine & fd baseline.
	time.Sleep(50 * time.Millisecond)
	baseG := runtime.NumGoroutine()
	baseFD := countOpenFDs(t)

	for i := 0; i < 5; i++ {
		want := int64(i + 1)
		m.dropAll()
		waitFor(t, 5*time.Second, func() bool {
			st := s.Stats()
			return !st.Down && st.Reconnects >= want
		}, "reconnect cycle")
	}

	// Beri waktu goroutine async lama selesai setelah cancel()+Close().
	waitFor(t, 2*time.Second, func() bool {
		runtime.GC()
		return runtime.NumGoroutine() <= baseG+2
	}, "goroutine count to settle after reconnects")

	if baseFD >= 0 {
		waitFor(t, 2*time.Second, func() bool {
			return countOpenFDs(t) <= baseFD+2
		}, "fd count to settle after reconnects")
	}
}

func TestOpenCloseNoGoroutineLeak(t *testing.T) {
	m := newMock(t)
	// Warm-up.
	s0 := openMock(t, m, nil)
	_ = s0.Close()
	m.dropAll()
	time.Sleep(50 * time.Millisecond)
	baseG := runtime.NumGoroutine()
	baseFD := countOpenFDs(t)

	for i := 0; i < 5; i++ {
		s := openMock(t, m, nil)
		_ = s.Close()
	}
	m.dropAll()

	waitFor(t, 2*time.Second, func() bool {
		runtime.GC()
		return runtime.NumGoroutine() <= baseG+1
	}, "goroutine count after Open/Close")

	if baseFD >= 0 {
		waitFor(t, 2*time.Second, func() bool {
			return countOpenFDs(t) <= baseFD+1
		}, "fd count after Open/Close")
	}
}
