package rosrun

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/quiqxiq/goros/v4"

	ros "example.com/roscatalog"
)

// Error yang bisa dibandingkan pemanggil. Tidak ada klasifikasi otomatis;
// pemanggil cukup memakai errors.Is bila perlu.
var (
	// ErrClosed dikembalikan setelah Close dipanggil.
	ErrClosed = errors.New("rosrun: session closed")
	// ErrReconnecting dikembalikan saat koneksi sedang dipulihkan.
	ErrReconnecting = errors.New("rosrun: reconnecting")
)

// MikroTik adalah satu sesi ke router: dua koneksi persisten (command dan
// feed stream), supervisor reconnect, dan pembungkus Runner.
//
// Koneksi dipisah supaya stream yang lambat tidak menahan command tulis.
// Keduanya memakai mode async masing-masing.
type MikroTik struct {
	cfg Config

	cmd  *managedConn
	feed *managedConn

	stopCh chan struct{}
	wg     sync.WaitGroup

	closeOnce sync.Once
	closeErr  error
	closed    atomic.Bool

	runCtx context.Context
	cancel context.CancelFunc

	epochs     *epochSignal
	reconnects atomic.Int64
	openedAt   time.Time

	status statusTracker

	cache   *feedCache
	limiter *limiter
}

// Open memvalidasi cfg lalu membuka dua koneksi ke router. Bila koneksi kedua
// gagal, koneksi pertama ditutup sebelum error dikembalikan.
func Open(ctx context.Context, cfg Config) (*MikroTik, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	cfg = cfg.withDefaults()

	lim := newLimiter(cfg.MaxInFlight)
	main, err := dialRetry(ctx, cfg, "main", lim)
	if err != nil {
		return nil, err
	}

	runCtx, cancel := context.WithCancel(context.Background())
	m := &MikroTik{
		cfg:      cfg,
		cmd:      main,
		feed:     main,
		limiter:  lim,
		stopCh:   make(chan struct{}),
		runCtx:   runCtx,
		cancel:   cancel,
		epochs:   newEpochSignal(),
		openedAt: time.Now(),
	}
	m.cache = newFeedCache(m)

	// Satu koneksi async adalah default (sama seperti MikroDash). Dua koneksi
	// hanya bila diminta eksplisit, dan opsi itu akan dihapus setelah satu rilis.
	m.wg.Add(1)
	go m.supervise(main)
	if cfg.DualConn {
		feed, err := dialRetry(ctx, cfg, "feed", lim)
		if err != nil {
			if c := main.get(); c != nil {
				_ = c.Close()
			}
			cancel()
			return nil, err
		}
		m.feed = feed
		m.wg.Add(1)
		go m.supervise(feed)
	}

	m.startStatus()

	return m, nil
}

func (m *MikroTik) logf(format string, args ...any) {
	if m.cfg.Logf != nil {
		m.cfg.Logf(format, args...)
	}
}

func (m *MikroTik) ready(conn *managedConn) error {
	if m.closed.Load() {
		return ErrClosed
	}
	if conn == nil {
		return ErrClosed
	}
	if conn.isDown() {
		return ErrReconnecting
	}
	return nil
}

// ---------------------------------------------------------------------------
// Command
// ---------------------------------------------------------------------------

// Exec menjalankan sentence aksi tulis dan mengembalikan nilai =ret=.
// Command tulis tidak diulang otomatis saat reconnect.
func (m *MikroTik) Exec(ctx context.Context, sentence []string) (string, error) {
	if err := m.ready(m.cmd); err != nil {
		return "", err
	}
	return m.cmd.run().Exec(ctx, sentence)
}

// Once menjalankan satu Request mode once. Bila koneksi command putus di
// tengah jalan, Once menunggu reconnect lalu mencoba sekali lagi; retry kedua
// memakai Fallback tanpa .proplist bila tersedia (beberapa versi RouterOS
// menutup koneksi untuk proplist pada menu tertentu).
func (m *MikroTik) Once(ctx context.Context, req ros.Request) ([]map[string]string, error) {
	if err := m.ready(m.cmd); err != nil {
		return nil, err
	}
	epoch := m.epochs.current()
	rows, err := m.cmd.run().Once(ctx, req)
	if err == nil || !isTransientDialErr(err) {
		return rows, err
	}
	if !m.epochs.wait(ctx, epoch, m.stopCh) || !m.waitCmdUp(ctx, epoch) {
		return nil, err
	}
	retry := req
	if len(req.Fallback) > 0 {
		retry.Sentence, retry.Fallback = req.Fallback, nil
	}
	return m.cmd.run().Once(ctx, retry)
}

// Count menjalankan Request count-only.
func (m *MikroTik) Count(ctx context.Context, req ros.Request) (int, error) {
	if err := m.ready(m.cmd); err != nil {
		return 0, err
	}
	epoch := m.epochs.current()
	n, err := m.cmd.run().Count(ctx, req)
	if err == nil || !isTransientDialErr(err) {
		return n, err
	}
	if !m.epochs.wait(ctx, epoch, m.stopCh) || !m.waitCmdUp(ctx, epoch) {
		return 0, err
	}
	return m.cmd.run().Count(ctx, req)
}

// waitCmdUp menunggu sampai koneksi command benar-benar up. Epoch dipakai
// karena reconnect feed juga menaikkannya.
func (m *MikroTik) waitCmdUp(ctx context.Context, epoch uint64) bool {
	for {
		if m.closed.Load() {
			return false
		}
		if !m.cmd.isDown() {
			return true
		}
		if !m.epochs.wait(ctx, epoch, m.stopCh) {
			return false
		}
		epoch = m.epochs.current()
	}
}

// ---------------------------------------------------------------------------
// Stream dan poll
// ---------------------------------------------------------------------------

// Watch menjalankan Request sesuai Delivery-nya di koneksi feed. Bila
// koneksi putus, Watch mengirim Event KindReconnecting lalu melanjutkan
// request yang sama setelah koneksi pulih. Command wajar yang berakhir sendiri
// (mis. ping dengan count) tidak diulang.
func (m *MikroTik) Watch(ctx context.Context, req ros.Request) (<-chan Event, error) {
	if err := m.ready(m.feed); err != nil {
		return nil, err
	}

	runner := m.feed.run()
	wctx, cancel := context.WithCancel(ctx)
	first, err := runner.Watch(wctx, req)
	if err != nil {
		cancel()
		return nil, err
	}

	out := make(chan Event)
	go m.forward(wctx, cancel, req, first, out)
	return out, nil
}

func (m *MikroTik) forward(ctx context.Context, cancel context.CancelFunc, req ros.Request, first <-chan Event, out chan<- Event) {
	defer close(out)
	defer cancel()

	ch := first
	for {
		epoch := m.epochs.current()
		restart := false
		terminal := false

		for ev := range ch {
			if ev.Kind == KindError {
				if shouldReconnect(ev.Err) {
					// Koneksi putus: jangan teruskan error mentah; pemakai
					// cukup melihat KindReconnecting.
					restart = true
					continue
				}
				// Penolakan router: teruskan, lalu berhenti.
				terminal = true
			}
			if !m.emit(ctx, out, ev) {
				return
			}
		}

		if ctx.Err() != nil || m.closed.Load() || terminal {
			return
		}

		// Poll dan stream non-native berjalan sampai ctx dibatalkan. Bila
		// channel tertutup tanpa KindError, itu tetap kehilangan koneksi,
		// bukan akhir yang wajar. Command native yang berhenti sendiri
		// (mis. ping count) tidak diulang.
		infinite := req.Delivery == ros.DeliveryPoll ||
			req.Delivery == ros.DeliveryEvents ||
			(req.Delivery == ros.DeliveryStream && req.Mode != ros.ModeNative)
		if !restart {
			if !infinite && !m.feed.isDown() && m.epochs.current() == epoch {
				return
			}
			restart = true
		}

		if !m.emit(ctx, out, Event{Kind: KindReconnecting, At: time.Now()}) {
			return
		}
		if !m.epochs.wait(ctx, epoch, m.stopCh) {
			return
		}

		// Epoch naik saat salah satu koneksi pulih. Stream hanya boleh
		// dimulai di feed yang benar-benar sudah up.
		var next <-chan Event
		for {
			if m.feed.isDown() {
				if !m.epochs.wait(ctx, m.epochs.current(), m.stopCh) {
					return
				}
				continue
			}
			n, err := m.feed.run().Watch(ctx, req)
			if err == nil {
				next = n
				break
			}
			if ctx.Err() != nil || m.closed.Load() {
				return
			}
			if !shouldReconnect(err) {
				_ = m.emit(ctx, out, Event{Kind: KindError, At: time.Now(), Err: err})
				return
			}
			// Feed masih transien; tunggu epoch berikutnya.
			if !m.epochs.wait(ctx, m.epochs.current(), m.stopCh) {
				return
			}
		}
		ch = next
	}
}

func (m *MikroTik) emit(ctx context.Context, out chan<- Event, ev Event) bool {
	select {
	case out <- ev:
		return true
	case <-ctx.Done():
		return false
	case <-m.runCtx.Done():
		return false
	}
}

// WatchAuto memakai plan bawaan Config bila plan kosong, lalu memilih stream
// atau poll. Bila router menolak stream sebelum ada data, ia beralih ke poll
// dengan jeda yang setara (lihat Runner.WatchAuto).
func (m *MikroTik) WatchAuto(ctx context.Context, res *ros.Resource, view string, ver ros.Version, plan ros.Plan, p ros.Params) (<-chan Event, error) {
	if plan.Delivery == "" {
		plan = m.PlanFor(res, view, ver)
	}

	req, err := res.Build(view, ver, plan, p)
	if err != nil {
		return nil, err
	}
	if req.Delivery != ros.DeliveryStream && req.Delivery != ros.DeliveryEvents {
		return m.Watch(ctx, req)
	}

	first, err := m.Watch(ctx, req)
	if err != nil {
		return nil, err
	}

	pollReq, pollErr := res.Build(view, ver, ros.Poll(m.pollEvery(res, view, ver, plan)), p)

	out := make(chan Event, 64)
	go func() {
		defer close(out)

		got := false
		stream := first
		for stream != nil {
			switched := false
			for ev := range stream {
				if ev.Kind == KindError && !got && pollErr == nil {
					var de *routeros.DeviceError
					if errors.As(ev.Err, &de) {
						next, werr := m.Watch(ctx, pollReq)
						if werr == nil {
							m.logf("rosrun: stream %s/%s ditolak router (%v); beralih ke poll", res.ID, view, ev.Err)
							stream, switched = next, true
							break
						}
					}
				}
				if ev.Kind != KindError {
					got = true
				}
				if !m.emit(ctx, out, ev) {
					return
				}
			}
			if !switched {
				return
			}
		}
	}()
	return out, nil
}

func (m *MikroTik) pollEvery(res *ros.Resource, view string, ver ros.Version, plan ros.Plan) time.Duration {
	every := plan.Every
	if v, ok := res.View(view, ver); ok && every < ros.MinPollEvery {
		every = v.Interval
	}
	if every < ros.MinPollEvery {
		every = m.cfg.PollEvery
	}
	if every < ros.MinPollEvery {
		every = ros.MinPollEvery
	}
	if every > ros.MaxPollEvery {
		every = ros.MaxPollEvery
	}
	return every
}

// PlanFor memilih plan bawaan untuk sebuah view. Default adalah STREAM
// (realtime). Poll hanya bila Config.Poll true, atau view tidak punya jalur
// stream. Plan eksplisit pemanggil tidak pernah ditimpa.
func (m *MikroTik) PlanFor(res *ros.Resource, view string, ver ros.Version) ros.Plan {
	wanted := ros.DeliveryStream
	if v, ok := res.View(view, ver); ok && v.Mode == ros.ModeListen {
		wanted = ros.DeliveryEvents
	}
	if m.cfg.Poll {
		wanted = ros.DeliveryPoll
	}

	allowed := res.Deliveries(view, ver)
	chosen := wanted
	found := false
	for _, d := range allowed {
		if d == wanted {
			found = true
			break
		}
	}
	if !found && len(allowed) > 0 {
		chosen = allowed[0]
	}

	plan := ros.Plan{Delivery: chosen}
	switch chosen {
	case ros.DeliveryPoll:
		plan.Every = m.cfg.PollEvery
	case ros.DeliveryStream:
		plan.Every = m.cfg.StreamEvery
	}
	return plan
}

// ---------------------------------------------------------------------------
// Status dan lifecycle
// ---------------------------------------------------------------------------

// Stats adalah snapshot status sesi.
type Stats struct {
	ID         string
	Address    string
	Open       bool
	Epoch      uint64
	Reconnects int64
	CmdDown    bool
	FeedDown   bool
	LastError  error
	Uptime     time.Duration

	// Instrumentasi per router.
	InFlight     int
	OpenStreams  int
	Commands     int64
	BusiestMenus []MenuCount
}

// Stats mengembalikan status sesi saat ini.
func (m *MikroTik) Stats() Stats {
	lastErr := m.feed.err()
	if lastErr == nil {
		lastErr = m.cmd.err()
	}

	st := Stats{
		ID:         m.cfg.ID,
		Address:    m.cfg.addr(),
		Open:       !m.closed.Load(),
		Epoch:      m.epochs.current(),
		Reconnects: m.reconnects.Load(),
		CmdDown:    m.cmd.isDown(),
		FeedDown:   m.feed.isDown(),
		LastError:  lastErr,
		Uptime:     time.Since(m.openedAt),
	}
	ls := m.limiter.stats()
	st.InFlight = ls.InFlight
	st.OpenStreams = ls.OpenStreams
	st.Commands = ls.Commands
	st.BusiestMenus = ls.Menus
	if !st.Open {
		st.Uptime = 0
	}
	return st
}

// Epoch naik setiap reconnect berhasil.
func (m *MikroTik) Epoch() uint64 { return m.epochs.current() }

// Close membatalkan stream yang berjalan dan menutup kedua koneksi. Aman
// dipanggil lebih dari sekali.
func (m *MikroTik) Close() error {
	m.closeOnce.Do(func() {
		m.closed.Store(true)
		// Terbitkan status error terakhir sebelum channel watcher ditutup.
		m.publishStatus()
		if m.cache != nil {
			m.cache.close()
		}
		m.cancel()
		close(m.stopCh)

		var errs []error
		conns := []*managedConn{m.cmd}
		if m.feed != m.cmd {
			conns = append(conns, m.feed)
		}
		for _, conn := range conns {
			if c := conn.get(); c != nil {
				if err := c.Close(); err != nil {
					errs = append(errs, fmt.Errorf("rosrun: close %s connection: %w", conn.name, err))
				}
			}
		}

		m.wg.Wait()
		m.closeErr = errors.Join(errs...)
	})
	return m.closeErr
}

// shouldReconnect membedakan koneksi putus dari penolakan command oleh
// router. Penolakan (DeviceError) diteruskan ke pemakai; selain itu dianggap
// gangguan koneksi.
func shouldReconnect(err error) bool {
	if err == nil {
		return true
	}
	var de *routeros.DeviceError
	return !errors.As(err, &de)
}

// ---------------------------------------------------------------------------
// Sinyal epoch
// ---------------------------------------------------------------------------

// epochSignal membroadcast perubahan epoch tanpa kehilangan event: setiap
// reconnect menutup channel lama dan membuat channel baru.
type epochSignal struct {
	mu sync.Mutex
	n  uint64
	ch chan struct{}
}

func newEpochSignal() *epochSignal {
	return &epochSignal{ch: make(chan struct{})}
}

func (e *epochSignal) current() uint64 {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.n
}

func (e *epochSignal) bump() {
	e.mu.Lock()
	e.n++
	close(e.ch)
	e.ch = make(chan struct{})
	e.mu.Unlock()
}

// wait menunggu sampai epoch berbeda dari seen.
func (e *epochSignal) wait(ctx context.Context, seen uint64, stop <-chan struct{}) bool {
	for {
		e.mu.Lock()
		if e.n != seen {
			e.mu.Unlock()
			return true
		}
		ch := e.ch
		e.mu.Unlock()

		select {
		case <-ch:
		case <-ctx.Done():
			return false
		case <-stop:
			return false
		}
	}
}
