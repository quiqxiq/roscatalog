package rosrun

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	ros "example.com/roscatalog"
)

// StatusState adalah keadaan ringkas satu router.
type StatusState string

const (
	// StatusConnected: kedua koneksi hidup dan ping status baru saja berhasil.
	StatusConnected StatusState = "connected"
	// StatusWarning: router masih terhubung, tetapi ada kondisi yang perlu
	// diperhatikan (TLS mati, packet loss, satu koneksi down, reconnecting).
	StatusWarning StatusState = "warning"
	// StatusError: sesi tertutup, kedua koneksi down, atau ping status sudah
	// melewati batas waktu.
	StatusError StatusState = "error"
)

var (
	// ErrStatusTimeout menandai ping status tidak lagi menerima balasan.
	ErrStatusTimeout = errors.New("rosrun: status ping timeout")
	// ErrBothConnectionsDown menandai kedua koneksi sesi tidak tersedia.
	ErrBothConnectionsDown = errors.New("rosrun: both connections down")
)

// PingStatus adalah data kumulatif probe /ping terakhir.
type PingStatus struct {
	Seq        int
	Sent       int
	Received   int
	PacketLoss float64 // persen, 0..100
	MinRTT     time.Duration
	AvgRTT     time.Duration
	MaxRTT     time.Duration
	LastRTT    time.Duration
	LastStatus string // "success" atau status router, mis. "timeout"
	Updated    time.Time
}

// Status adalah snapshot status realtime satu router.
type Status struct {
	ID      string
	Address string

	// State adalah connected, warning, atau error.
	State StatusState
	// Connected berarti kedua koneksi hidup dan ping status masih segar.
	// Nilainya bisa true walaupun State=warning (mis. TLS mati).
	Connected bool
	// Warnings berisi alasan state warning.
	Warnings []string
	// Err terisi saat State=error.
	Err error

	// Target adalah alamat ping status (bawaan 8.8.8.8).
	Target string
	Ping   PingStatus

	CmdDown    bool
	FeedDown   bool
	Epoch      uint64
	Reconnects int64
	Uptime     time.Duration
	Updated    time.Time
}

// statusTracker menyimpan data ping dan subscriber status. Semua field
// dilindungi mu; Status dibangun dari data ini plus Stats koneksi.
type statusTracker struct {
	mu sync.Mutex

	started time.Time
	ping    PingStatus
	lastOK  time.Time
	pingErr error

	current Status
	subs    map[chan Status]struct{}
}

// startStatus menginisialisasi monitor status. Monitor selalu memakai native
// stream /ping, bukan poll, dan berjalan terlepas dari Config.Stream.
func (m *MikroTik) startStatus() {
	m.status.mu.Lock()
	m.status.started = time.Now()
	if m.status.subs == nil {
		m.status.subs = make(map[chan Status]struct{})
	}
	m.status.mu.Unlock()

	m.publishStatus()
	if m.cfg.Status.Disable {
		return
	}

	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		m.statusLoop()
	}()
}

// statusRequest menyusun /ping native tanpa count/duration. Status sengaja
// tidak memakai .proplist agar tetap kompatibel dengan RouterOS 6.
func (m *MikroTik) statusRequest() ros.Request {
	every := m.cfg.Status.Every
	return ros.Request{
		ResourceID: "tool-ping",
		View:       "run",
		Delivery:   ros.DeliveryStream,
		Mode:       ros.ModeNative,
		Every:      every,
		Shape:      ros.ShapeRows,
		Sentence: []string{
			"/ping",
			"=address=" + m.cfg.Status.Target,
			"=interval=" + ros.FormatTime(every),
		},
		Fields: []string{
			"seq", "host", "size", "ttl", "time", "status", "sent", "received",
			"packet-loss", "min-rtt", "avg-rtt", "max-rtt",
		},
	}
}

// statusLoop menjaga stream ping tetap hidup. Bila feed putus, Watch sudah
// melanjutkan request yang sama; loop ini hanya menangani kegagalan awal atau
// penolakan permanen dari router.
func (m *MikroTik) statusLoop() {
	req := m.statusRequest()
	backoff := time.Second

	for {
		if m.runCtx.Err() != nil {
			return
		}

		ch, err := m.Watch(m.runCtx, req)
		if err != nil {
			m.setPingProblem(err)
			m.publishStatus()
			if !sleepContext(m.runCtx, backoff) {
				return
			}
			if backoff < 10*time.Second {
				backoff *= 2
			}
			continue
		}
		backoff = time.Second

		ticker := time.NewTicker(time.Second)
	consume:
		for {
			select {
			case ev, ok := <-ch:
				if !ok {
					break consume
				}
				switch ev.Kind {
				case KindRow:
					m.observePing(ev.Row)
					m.publishStatus()
				case KindError:
					if m.runCtx.Err() == nil {
						m.setPingProblem(ev.Err)
						m.publishStatus()
					}
				}
			case <-ticker.C:
				// Tetap terbitkan status walau tidak ada baris, supaya
				// timeout dan perubahan koneksi terlihat realtime.
				m.publishStatus()
			case <-m.runCtx.Done():
				ticker.Stop()
				return
			}
		}
		ticker.Stop()

		if m.runCtx.Err() != nil {
			return
		}
		if !sleepContext(m.runCtx, time.Second) {
			return
		}
	}
}

func sleepContext(ctx context.Context, d time.Duration) bool {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-t.C:
		return true
	}
}

// observePing memperbarui data ping dari satu baris /ping.
func (m *MikroTik) observePing(row map[string]string) {
	p := parsePingRow(row)
	now := time.Now()

	m.status.mu.Lock()
	prev := m.status.ping
	// Field kumulatif tidak ada di baris ini: pertahankan nilai sebelumnya
	// daripada melaporkan 0%.
	if p.Sent == 0 && prev.Sent > 0 {
		p.Sent = prev.Sent
		p.Received = prev.Received
		p.PacketLoss = prev.PacketLoss
		p.MinRTT, p.AvgRTT, p.MaxRTT = prev.MinRTT, prev.AvgRTT, prev.MaxRTT
	}
	m.status.ping = p
	m.status.pingErr = nil
	if pingSucceeded(p) {
		m.status.lastOK = now
	}
	m.status.mu.Unlock()
}

func pingSucceeded(p PingStatus) bool {
	s := strings.ToLower(strings.TrimSpace(p.LastStatus))
	return s == "" || s == "success" || s == "ok"
}

// parsePingRow mengubah satu baris /ping menjadi PingStatus. Field yang tidak
// ada dibiarkan nol agar kompatibel dengan RouterOS 6/7.
func parsePingRow(row map[string]string) PingStatus {
	var p PingStatus
	p.Seq = atoi(row["seq"])
	p.Sent = atoi(row["sent"])
	p.Received = atoi(row["received"])
	if v, ok := row["packet-loss"]; ok {
		p.PacketLoss, _ = strconv.ParseFloat(strings.TrimSuffix(v, "%"), 64)
	}
	p.LastRTT, _ = time.ParseDuration(row["time"])
	p.MinRTT, _ = time.ParseDuration(row["min-rtt"])
	p.AvgRTT, _ = time.ParseDuration(row["avg-rtt"])
	p.MaxRTT, _ = time.ParseDuration(row["max-rtt"])
	p.LastStatus = strings.ToLower(strings.TrimSpace(row["status"]))
	if p.LastStatus == "" {
		p.LastStatus = "success"
	}
	p.Updated = time.Now()
	return p
}

func atoi(s string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	return n
}

// setPingProblem mencatat kegagalan stream ping. Nil berarti pulih.
func (m *MikroTik) setPingProblem(err error) {
	m.status.mu.Lock()
	m.status.pingErr = err
	m.status.mu.Unlock()
}

// Status mengembalikan snapshot status terbaru. Aman dipanggil kapan saja,
// termasuk sebelum monitor ping menerima baris pertama.
func (m *MikroTik) Status() Status {
	m.status.mu.Lock()
	defer m.status.mu.Unlock()
	return m.buildStatusLocked(time.Now())
}

// WatchStatus mengirim snapshot status saat ini lalu pembaruan realtime.
// Channel ditutup saat ctx selesai atau sesi ditutup. Update yang datang lebih
// cepat dari konsumen di-coalesce: yang dijamin adalah snapshot terbaru.
func (m *MikroTik) WatchStatus(ctx context.Context) <-chan Status {
	ch := make(chan Status, 1)

	m.status.mu.Lock()
	if m.status.subs == nil {
		m.status.subs = make(map[chan Status]struct{})
	}
	m.status.subs[ch] = struct{}{}
	sendLatest(ch, m.buildStatusLocked(time.Now()))
	m.status.mu.Unlock()

	go func() {
		select {
		case <-ctx.Done():
		case <-m.runCtx.Done():
		}
		m.status.mu.Lock()
		if _, ok := m.status.subs[ch]; ok {
			delete(m.status.subs, ch)
			close(ch)
		}
		m.status.mu.Unlock()
	}()

	return ch
}

// publishStatus menghitung ulang status dan meneruskannya ke subscriber.
// Snapshot yang tidak berubah (mis. ticker timeout tanpa baris ping baru)
// tidak dikirim ulang agar subscriber tidak dibanjiri duplikat.
func (m *MikroTik) publishStatus() {
	now := time.Now()
	m.status.mu.Lock()
	st := m.buildStatusLocked(now)
	prev := m.status.current
	m.status.current = st
	if !sameStatus(prev, st) {
		for ch := range m.status.subs {
			sendLatest(ch, st)
		}
	}
	m.status.mu.Unlock()
}

func sameStatus(a, b Status) bool {
	if a.ID != b.ID || a.Address != b.Address || a.State != b.State ||
		a.Connected != b.Connected || a.Target != b.Target ||
		a.CmdDown != b.CmdDown || a.FeedDown != b.FeedDown ||
		a.Epoch != b.Epoch || a.Reconnects != b.Reconnects ||
		!slices.Equal(a.Warnings, b.Warnings) || !sameErr(a.Err, b.Err) {
		return false
	}
	pa, pb := a.Ping, b.Ping
	return pa.Seq == pb.Seq && pa.Sent == pb.Sent && pa.Received == pb.Received &&
		pa.PacketLoss == pb.PacketLoss && pa.MinRTT == pb.MinRTT &&
		pa.AvgRTT == pb.AvgRTT && pa.MaxRTT == pb.MaxRTT &&
		pa.LastRTT == pb.LastRTT && pa.LastStatus == pb.LastStatus
}

func sameErr(a, b error) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return a.Error() == b.Error()
}

// sendLatest mengisi channel berkapasitas satu dengan nilai terbaru.
func sendLatest(ch chan Status, st Status) {
	select {
	case ch <- st:
		return
	default:
	}
	select {
	case <-ch:
	default:
	}
	select {
	case ch <- st:
	default:
	}
}

// buildStatusLocked membangun Status dari tracker dan Stats koneksi. Pemanggil
// wajib memegang status.mu.
func (m *MikroTik) buildStatusLocked(now time.Time) Status {
	cfg := m.cfg
	open := !m.closed.Load()
	cmdDown, feedDown := m.cmd.isDown(), m.feed.isDown()
	lastErr := m.feed.err()
	if lastErr == nil {
		lastErr = m.cmd.err()
	}

	st := Status{
		ID:         cfg.ID,
		Address:    cfg.addr(),
		State:      StatusConnected,
		Target:     cfg.Status.Target,
		Ping:       m.status.ping,
		CmdDown:    cmdDown,
		FeedDown:   feedDown,
		Epoch:      m.epochs.current(),
		Reconnects: m.reconnects.Load(),
		Updated:    now,
	}
	if open {
		st.Uptime = now.Sub(m.openedAt)
	}

	st.Connected = open && !cmdDown && !feedDown
	if !cfg.Status.Disable {
		timeout := cfg.Status.Timeout
		switch {
		case !m.status.lastOK.IsZero() && now.Sub(m.status.lastOK) > timeout:
			st.Connected = false
			st.State = StatusError
			st.Err = fmt.Errorf("%w: no reply from %s for %s", ErrStatusTimeout,
				cfg.Status.Target, now.Sub(m.status.lastOK).Round(time.Second))
		case m.status.lastOK.IsZero() && now.Sub(m.status.started) > timeout:
			st.Connected = false
			st.State = StatusError
			st.Err = fmt.Errorf("%w: no reply from %s since %s", ErrStatusTimeout,
				cfg.Status.Target, now.Sub(m.status.started).Round(time.Second))
		}
	}

	if !open {
		st.State = StatusError
		st.Err = ErrClosed
	} else if cmdDown && feedDown {
		st.State = StatusError
		if lastErr != nil {
			st.Err = fmt.Errorf("%w: %v", ErrBothConnectionsDown, lastErr)
		} else {
			st.Err = ErrBothConnectionsDown
		}
	}

	// Warning dikumpulkan terlepas dari error agar UI bisa menampilkan
	// seluruh alasan yang diketahui.
	if cmdDown && !feedDown {
		st.Warnings = append(st.Warnings, "command connection down")
	}
	if feedDown && !cmdDown {
		st.Warnings = append(st.Warnings, "feed connection down")
	}
	if lastErr != nil && (cmdDown || feedDown) && st.State != StatusError {
		st.Warnings = append(st.Warnings, lastErr.Error())
	}
	if err := m.status.pingErr; err != nil {
		st.Warnings = append(st.Warnings, "status ping: "+err.Error())
	}
	if !cfg.TLS && !cfg.Status.SuppressTLSWarning {
		st.Warnings = append(st.Warnings, "TLS disabled: API traffic is plaintext")
	}
	if cfg.Status.WarnPacketLoss > 0 && st.Ping.PacketLoss > cfg.Status.WarnPacketLoss {
		st.Warnings = append(st.Warnings, fmt.Sprintf("packet loss %.1f%%", st.Ping.PacketLoss))
	}
	if !cfg.Status.Disable && st.State != StatusError && m.status.lastOK.IsZero() {
		st.Warnings = append(st.Warnings, "waiting for first status ping reply")
		if st.Connected {
			st.State = StatusWarning
		}
	}
	if st.State != StatusError && (len(st.Warnings) > 0 || !st.Connected) {
		st.State = StatusWarning
	}
	return st
}

// Statuses mengembalikan status semua router yang dikelola Manager.
func (m *Manager) Statuses() []Status {
	if m == nil {
		return nil
	}
	m.mu.RLock()
	order := append([]string(nil), m.order...)
	routers := make([]*MikroTik, 0, len(order))
	for _, id := range order {
		routers = append(routers, m.routers[id])
	}
	m.mu.RUnlock()

	out := make([]Status, 0, len(routers))
	for _, r := range routers {
		out = append(out, r.Status())
	}
	return out
}
