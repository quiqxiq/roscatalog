package rosrun

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestParsePingRow(t *testing.T) {
	p := parsePingRow(map[string]string{
		"seq": "4", "sent": "5", "received": "4", "packet-loss": "20",
		"time": "30ms", "min-rtt": "10ms", "avg-rtt": "20ms", "max-rtt": "30ms",
	})
	if p.Seq != 4 || p.Sent != 5 || p.Received != 4 || p.PacketLoss != 20 {
		t.Errorf("counter: %+v", p)
	}
	if p.LastRTT != 30*time.Millisecond || p.MinRTT != 10*time.Millisecond ||
		p.AvgRTT != 20*time.Millisecond || p.MaxRTT != 30*time.Millisecond {
		t.Errorf("rtt: %+v", p)
	}
	if !pingSucceeded(p) || p.LastStatus != "success" {
		t.Errorf("status sukses: %+v", p)
	}

	f := parsePingRow(map[string]string{"sent": "6", "received": "4", "packet-loss": "33", "status": "timeout"})
	if pingSucceeded(f) || f.LastStatus != "timeout" {
		t.Errorf("status timeout: %+v", f)
	}
}

func TestStatusUsesNativePingRegardlessOfStream(t *testing.T) {
	m := newMock(t)
	s := openMock(t, m, func(c *Config) {
		c.Poll = true // status tetap harus native stream
		c.Status = StatusConfig{
			PingInterval:       time.Second,
			ReplyTimeout:       5 * time.Second,
			SuppressTLSWarning: true,
			WarnPacketLoss:     100,
		}
	})

	waitFor(t, 4*time.Second, func() bool {
		st := s.Status()
		return st.Connected && st.State == StatusConnected
	}, "status connected dari native ping")

	st := s.Status()
	if st.Target != "8.8.8.8" {
		t.Errorf("target = %q, want 8.8.8.8", st.Target)
	}
	if st.Ping.Received == 0 {
		t.Errorf("ping belum menerima balasan: %+v", st.Ping)
	}

	calls := m.callsTo("/ping")
	if len(calls) == 0 {
		t.Fatal("native /ping tidak dipanggil")
	}
	c := calls[0]
	if c.attrs["address"] != "8.8.8.8" || c.attrs["interval"] != "1s" {
		t.Errorf("sentence ping: %+v", c.attrs)
	}
	if c.has("count") || c.has("duration") || c.has(".proplist") {
		t.Errorf("status ping tidak boleh count/duration/proplist: %+v", c.attrs)
	}
}

func TestStatusRequestOmitsPingIntervalByDefault(t *testing.T) {
	m := newMock(t)
	s := openMock(t, m, nil) // status disabled; hanya memeriksa sentence

	req := s.statusRequest()
	for _, arg := range req.Sentence {
		if strings.HasPrefix(arg, "=interval=") {
			t.Fatalf("default status tidak boleh mengirim interval: %v", req.Sentence)
		}
	}
}

func TestStatusRequestHonorsPingInterval(t *testing.T) {
	m := newMock(t)
	s := openMock(t, m, func(c *Config) {
		c.Status.PingInterval = 2 * time.Second
		c.Status.ReplyTimeout = 5 * time.Second
	})

	req := s.statusRequest()
	for _, arg := range req.Sentence {
		if arg == "=interval=2s" {
			return
		}
	}
	t.Fatalf("interval ping tidak dikirim: %v", req.Sentence)
}

func TestStatusWarnsTLSByDefault(t *testing.T) {
	m := newMock(t)
	s := openMock(t, m, func(c *Config) {
		c.Status = StatusConfig{ReplyTimeout: 5 * time.Second}
	})

	waitFor(t, 4*time.Second, func() bool { return s.Status().Connected }, "status ping pertama")
	st := s.Status()
	if st.State != StatusWarning {
		t.Fatalf("state = %s, want warning (TLS mati)", st.State)
	}
	joined := strings.Join(st.Warnings, " | ")
	if !strings.Contains(strings.ToLower(joined), "tls") {
		t.Errorf("warning TLS tidak ada: %v", st.Warnings)
	}
}

func TestStatusTimeoutBecomesError(t *testing.T) {
	m := newMock(t)
	m.pingFail.Store(true)
	s := openMock(t, m, func(c *Config) {
		c.Status = StatusConfig{
			PingInterval:       time.Second,
			ReplyTimeout:       time.Second,
			SuppressTLSWarning: true,
		}
	})

	waitFor(t, 5*time.Second, func() bool {
		st := s.Status()
		return st.State == StatusError && errors.Is(st.Err, ErrStatusTimeout)
	}, "status error karena ping timeout")
}

func TestWatchStatusRealtime(t *testing.T) {
	m := newMock(t)
	s := openMock(t, m, func(c *Config) {
		c.Status = StatusConfig{
			PingInterval:       time.Second,
			ReplyTimeout:       5 * time.Second,
			SuppressTLSWarning: true,
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := s.WatchStatus(ctx)

	select {
	case st, ok := <-ch:
		if !ok {
			t.Fatal("channel status tertutup")
		}
		if st.ID != "test" {
			t.Errorf("status awal: %+v", st)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("tidak ada status awal")
	}

	deadline := time.After(4 * time.Second)
	for {
		select {
		case st, ok := <-ch:
			if !ok {
				t.Fatal("channel status tertutup sebelum connected")
			}
			if st.Connected && st.State == StatusConnected {
				return
			}
		case <-deadline:
			t.Fatalf("tidak ada update connected: %+v", s.Status())
		}
	}
}

func TestStatusConnectionDownUsesErrConnectionDown(t *testing.T) {
	m := newMock(t)
	s := openMock(t, m, nil) // status ping dimatikan; fokus ke koneksi

	m.close() // listener ikut tutup, jadi reconnect tidak bisa berhasil

	waitFor(t, 3*time.Second, func() bool {
		st := s.Status()
		return st.Down && !st.Connected && errors.Is(st.Err, ErrConnectionDown)
	}, "status error karena koneksi down")
}

func TestStatusDisabledStillReportsConnections(t *testing.T) {
	m := newMock(t)
	s := openMock(t, m, nil) // openMock menonaktifkan status

	st := s.Status()
	if !st.Connected || st.State != StatusConnected {
		t.Errorf("status tanpa ping: %+v", st)
	}
	if len(m.callsTo("/ping")) != 0 {
		t.Error("status disabled tidak boleh mengirim /ping")
	}
}
