package rosrun

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	ros "example.com/roscatalog"
)

// RouterEvent adalah Event yang diberi identitas router asal.
type RouterEvent struct {
	RouterID string
	Address  string
	Event
}

// RouterStatus adalah Status yang diberi identitas router asal.
type RouterStatus struct {
	RouterID string
	Address  string
	Status   Status
}

// Manager mengelola banyak sesi MikroTik sekaligus. Tiap router tetap memiliki
// satu koneksi sendiri, supervisor reconnect sendiri, dan monitor status
// sendiri; Manager hanya mengoordinasikan pembukaan dan fan-out event.
type Manager struct {
	mu      sync.RWMutex
	routers map[string]*MikroTik
	order   []string

	// Lifecycle: refcount pemakai, hold non-viewer, dan grace sebelum tutup.
	// cfgs/cfgOrder menyimpan registrasi (Register/Acquire) walau sesi sedang
	// ditutup sementara, sehingga Acquire berikutnya bisa membukanya lagi.
	cfgs     map[string]Config
	cfgOrder []string
	refs     map[string]int
	holds    map[string]map[string]bool
	linger   map[string]*time.Timer

	closed    atomic.Bool
	closeOnce sync.Once
	closeErr  error
}

// ensureMapsLocked menyiapkan map lifecycle untuk Manager yang dibuat manual.
// Caller memegang m.mu.
func (m *Manager) ensureMapsLocked() {
	if m.routers == nil {
		m.routers = map[string]*MikroTik{}
	}
	if m.cfgs == nil {
		m.cfgs = map[string]Config{}
	}
	if m.refs == nil {
		m.refs = map[string]int{}
	}
	if m.holds == nil {
		m.holds = map[string]map[string]bool{}
	}
	if m.linger == nil {
		m.linger = map[string]*time.Timer{}
	}
}

// NewManager membuat Manager kosong siap pakai.
func NewManager() *Manager {
	return &Manager{
		routers: make(map[string]*MikroTik),
		cfgs:    make(map[string]Config),
		refs:    make(map[string]int),
		holds:   make(map[string]map[string]bool),
		linger:  make(map[string]*time.Timer),
	}
}

// Register mendaftarkan atau memperbarui konfigurasi router ke manager tanpa
// langsung membuka koneksi. Koneksi dibuka saat pertama kali di-Acquire.
// Registrasi tetap tersimpan setelah sesi ditutup karena grace, sehingga
// Acquire berikutnya membuka sesi baru dengan Config yang sama.
func (m *Manager) Register(cfg Config) {
	if m == nil {
		return
	}
	if strings.TrimSpace(cfg.ID) == "" {
		cfg.ID = cfg.addr()
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ensureMapsLocked()
	m.registerLocked(cfg)
}

// registerLocked menyimpan Config dan urutan registrasi. Caller memegang m.mu.
func (m *Manager) registerLocked(cfg Config) {
	if _, exists := m.cfgs[cfg.ID]; !exists {
		m.cfgOrder = append(m.cfgOrder, cfg.ID)
	}
	m.cfgs[cfg.ID] = cfg
}

// OpenAll membuka semua router secara paralel. Bila sebagian gagal, router yang
// berhasil tetap dikembalikan bersama error gabungan (partial success).
// Bila tidak ada satu pun yang berhasil, hasilnya nil dan error.
func OpenAll(ctx context.Context, cfgs []Config) (*Manager, error) {
	m := &Manager{
		routers: make(map[string]*MikroTik),
		cfgs:    make(map[string]Config),
		refs:    make(map[string]int),
		holds:   make(map[string]map[string]bool),
		linger:  make(map[string]*time.Timer),
	}
	if len(cfgs) == 0 {
		return m, nil
	}

	normalized := make([]Config, len(cfgs))
	seen := make(map[string]bool, len(cfgs))
	var setupErrs []error
	var valid []int
	for i, cfg := range cfgs {
		if strings.TrimSpace(cfg.ID) == "" {
			cfg.ID = cfg.addr()
		}
		normalized[i] = cfg
		if seen[cfg.ID] {
			setupErrs = append(setupErrs, fmt.Errorf("rosrun: duplicate router id %q", cfg.ID))
			continue
		}
		seen[cfg.ID] = true
		valid = append(valid, i)
		// Simpan registrasi walau dial gagal: Acquire/WatchAll bisa mencoba
		// lagi tanpa perlu Register manual.
		m.registerLocked(cfg)
	}

	type result struct {
		sess *MikroTik
		err  error
	}
	results := make([]result, len(cfgs))
	var wg sync.WaitGroup
	for _, i := range valid {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sess, err := Open(ctx, normalized[i])
			results[i] = result{sess: sess, err: err}
		}(i)
	}
	wg.Wait()

	errs := append([]error(nil), setupErrs...)
	for _, i := range valid {
		cfg, res := normalized[i], results[i]
		if res.err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", cfg.ID, res.err))
			continue
		}
		m.routers[cfg.ID] = res.sess
		m.refs[cfg.ID] = 0
		m.order = append(m.order, cfg.ID)
	}

	if len(m.order) == 0 {
		return nil, errors.Join(errs...)
	}
	return m, errors.Join(errs...)
}

// Add membuka satu router lagi dan menambahkannya ke manager.
func (m *Manager) Add(ctx context.Context, cfg Config) (*MikroTik, error) {
	if m == nil {
		return nil, ErrClosed
	}
	if m.closed.Load() {
		return nil, ErrClosed
	}
	if strings.TrimSpace(cfg.ID) == "" {
		cfg.ID = cfg.addr()
	}

	m.mu.RLock()
	_, exists := m.routers[cfg.ID]
	m.mu.RUnlock()
	if exists {
		return nil, fmt.Errorf("rosrun: router id %q already exists", cfg.ID)
	}

	sess, err := Open(ctx, cfg)
	if err != nil {
		return nil, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.ensureMapsLocked()
	if m.closed.Load() {
		_ = sess.Close()
		return nil, ErrClosed
	}
	if _, exists := m.routers[cfg.ID]; exists {
		_ = sess.Close()
		return nil, fmt.Errorf("rosrun: router id %q already exists", cfg.ID)
	}
	m.registerLocked(cfg)
	m.routers[cfg.ID] = sess
	m.refs[cfg.ID] = 0
	m.order = append(m.order, cfg.ID)
	return sess, nil
}

// Get mencari sesi berdasarkan ID.
func (m *Manager) Get(id string) (*MikroTik, bool) {
	if m == nil {
		return nil, false
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.routers[id]
	return s, ok
}

// IDs mengembalikan ID router sesuai urutan pendaftaran.
func (m *Manager) IDs() []string {
	if m == nil {
		return nil
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]string(nil), m.order...)
}

// Routers mengembalikan sesi sesuai urutan pendaftaran.
func (m *Manager) Routers() []*MikroTik {
	if m == nil {
		return nil
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*MikroTik, 0, len(m.order))
	for _, id := range m.order {
		out = append(out, m.routers[id])
	}
	return out
}

// Len mengembalikan jumlah router aktif.
func (m *Manager) Len() int {
	if m == nil {
		return 0
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.order)
}

// Acquire menaikkan refcount sesi untuk router id, membuka sesi bila belum ada.
// Grace yang sedang berjalan dibatalkan.
func (m *Manager) Acquire(ctx context.Context, id string) (*MikroTik, error) {
	if m == nil {
		return nil, ErrClosed
	}
	m.mu.Lock()
	m.ensureMapsLocked()
	if m.closed.Load() {
		m.mu.Unlock()
		return nil, ErrClosed
	}
	if s := m.routers[id]; s != nil {
		m.refs[id]++
		if t := m.linger[id]; t != nil {
			t.Stop()
			delete(m.linger, id)
		}
		m.mu.Unlock()
		return s, nil
	}
	cfg, ok := m.cfgs[id]
	m.mu.Unlock()
	if !ok {
		return nil, fmt.Errorf("rosrun: unknown router %q", id)
	}

	sess, err := Open(ctx, cfg)
	if err != nil {
		return nil, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed.Load() {
		_ = sess.Close()
		return nil, ErrClosed
	}
	if s := m.routers[id]; s != nil {
		// Balapan dengan Acquire lain; sesi baru tidak dipakai.
		_ = sess.Close()
		m.refs[id]++
		return s, nil
	}
	m.routers[id] = sess
	m.refs[id] = 1
	m.order = append(m.order, id)
	return sess, nil
}

// Release menurunkan refcount. Bila mencapai nol dan tidak ada hold, grace
// dimulai; sesi baru benar-benar ditutup setelah grace lewat.
func (m *Manager) Release(id string) {
	if m == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.routers[id] == nil {
		return
	}
	if m.refs[id] > 0 {
		m.refs[id]--
	}
	if m.refs[id] == 0 {
		m.armLingerLocked(id)
	}
}

// Retain menahan sesi hidup tanpa viewer. Alasan yang sama boleh diulang.
func (m *Manager) Retain(id, reason string) error {
	if m == nil {
		return ErrClosed
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ensureMapsLocked()
	if m.routers[id] == nil {
		return fmt.Errorf("rosrun: unknown router %q", id)
	}
	if m.holds[id] == nil {
		m.holds[id] = map[string]bool{}
	}
	m.holds[id][reason] = true
	if t := m.linger[id]; t != nil {
		t.Stop()
		delete(m.linger, id)
	}
	return nil
}

// Drop melepas satu hold. Bila itu hold terakhir dan tidak ada viewer, grace
// dimulai.
func (m *Manager) Drop(id, reason string) {
	if m == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.holds[id] != nil {
		delete(m.holds[id], reason)
	}
	if m.refs[id] == 0 && len(m.holds[id]) == 0 {
		m.armLingerLocked(id)
	}
}

// Subscribe adalah Acquire + shared feed + Release saat subscription ditutup.
func (m *Manager) Subscribe(ctx context.Context, id string, req ros.Request) (*Subscription, error) {
	sess, err := m.Acquire(ctx, id)
	if err != nil {
		return nil, err
	}
	sub, err := sess.Subscribe(ctx, req)
	if err != nil {
		m.Release(id)
		return nil, err
	}
	sub.onClose = func() { m.Release(id) }
	return sub, nil
}

// armLingerLocked menjadwalkan teardown. Caller memegang m.mu.
func (m *Manager) armLingerLocked(id string) {
	if len(m.holds[id]) > 0 {
		return
	}
	if t := m.linger[id]; t != nil {
		t.Stop()
	}
	grace := m.cfgs[id].Grace
	if grace <= 0 {
		grace = DefaultGrace
	}
	m.linger[id] = time.AfterFunc(grace, func() { m.idleOut(id) })
}

func (m *Manager) idleOut(id string) {
	m.mu.Lock()
	if m.closed.Load() {
		m.mu.Unlock()
		return
	}
	if m.refs[id] > 0 || len(m.holds[id]) > 0 {
		m.mu.Unlock()
		return
	}
	sess := m.routers[id]
	delete(m.routers, id)
	// cfgs/cfgOrder sengaja dipertahankan: ini registrasi, bukan sesi.
	delete(m.refs, id)
	delete(m.holds, id)
	if t := m.linger[id]; t != nil {
		t.Stop()
		delete(m.linger, id)
	}
	for i, x := range m.order {
		if x == id {
			m.order = append(m.order[:i], m.order[i+1:]...)
			break
		}
	}
	m.mu.Unlock()
	if sess != nil {
		_ = sess.Close()
	}
}

// Remove memutuskan koneksi dan menghapus router dari manager.
func (m *Manager) Remove(id string) error {
	if m == nil {
		return ErrClosed
	}
	m.mu.Lock()
	if m.closed.Load() {
		m.mu.Unlock()
		return ErrClosed
	}
	sess, ok := m.routers[id]
	if !ok {
		m.mu.Unlock()
		return fmt.Errorf("rosrun: router id %q not found", id)
	}
	delete(m.routers, id)
	delete(m.cfgs, id)
	delete(m.refs, id)
	delete(m.holds, id)
	if t := m.linger[id]; t != nil {
		t.Stop()
		delete(m.linger, id)
	}
	for i, x := range m.order {
		if x == id {
			m.order = append(m.order[:i], m.order[i+1:]...)
			break
		}
	}
	for i, x := range m.cfgOrder {
		if x == id {
			m.cfgOrder = append(m.cfgOrder[:i], m.cfgOrder[i+1:]...)
			break
		}
	}
	m.mu.Unlock()
	if sess != nil {
		return sess.Close()
	}
	return nil
}

// Unregister adalah alias dari Remove.
func (m *Manager) Unregister(id string) error {
	return m.Remove(id)
}

// Config mengembalikan Config untuk router id jika ada.
func (m *Manager) Config(id string) (Config, bool) {
	if m == nil {
		return Config{}, false
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	cfg, ok := m.cfgs[id]
	return cfg, ok
}

// Has melaporkan apakah konfigurasi router id sudah didaftarkan.
func (m *Manager) Has(id string) bool {
	if m == nil {
		return false
	}
	_, ok := m.Config(id)
	return ok
}

// Configs mengembalikan seluruh Config router terdaftar (termasuk yang
// sesinya sedang ditutup karena grace) sesuai urutan pendaftaran.
func (m *Manager) Configs() []Config {
	if m == nil {
		return nil
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Config, 0, len(m.cfgOrder))
	for _, id := range m.cfgOrder {
		if cfg, ok := m.cfgs[id]; ok {
			out = append(out, cfg)
		}
	}
	return out
}

// WatchAll menjalankan Request yang sama di semua router terdaftar. Tiap
// router di-Acquire selama stream hidup, jadi sesi tidak bisa ditutup oleh
// grace saat masih dipakai. Router yang belum terbuka dibuka otomatis dari
// Config terdaftar. Kegagalan satu router menjadi RouterEvent KindError dan
// tidak menghentikan router lain. Channel ditutup setelah semua stream berhenti
// atau ctx selesai.
func (m *Manager) WatchAll(ctx context.Context, req ros.Request) (<-chan RouterEvent, error) {
	ids, err := m.watchIDs()
	if err != nil {
		return nil, err
	}
	out := make(chan RouterEvent, 64)
	if len(ids) == 0 {
		close(out)
		return out, nil
	}

	var wg sync.WaitGroup
	for _, id := range ids {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			addr := m.routerAddress(id)
			sess, aerr := m.Acquire(ctx, id)
			if aerr != nil {
				emitRouterEventID(ctx, out, id, addr, Event{Kind: KindError, At: time.Now(), Err: aerr})
				return
			}
			defer m.Release(id)
			ch, werr := sess.Watch(ctx, req)
			if werr != nil {
				emitRouterEventID(ctx, out, id, addr, Event{Kind: KindError, At: time.Now(), Err: werr})
				return
			}
			for ev := range ch {
				if !emitRouterEventID(ctx, out, id, addr, ev) {
					return
				}
			}
		}(id)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out, nil
}

// WatchAutoAll seperti WatchAll, tetapi tiap router memilih stream/poll dari
// Config.Poll-nya sendiri saat plan.Delivery kosong.
func (m *Manager) WatchAutoAll(ctx context.Context, res *ros.Resource, view string, plan ros.Plan, p ros.Params) (<-chan RouterEvent, error) {
	ids, err := m.watchIDs()
	if err != nil {
		return nil, err
	}
	out := make(chan RouterEvent, 64)
	if len(ids) == 0 {
		close(out)
		return out, nil
	}

	var wg sync.WaitGroup
	for _, id := range ids {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			addr := m.routerAddress(id)
			sess, aerr := m.Acquire(ctx, id)
			if aerr != nil {
				emitRouterEventID(ctx, out, id, addr, Event{Kind: KindError, At: time.Now(), Err: aerr})
				return
			}
			defer m.Release(id)
			rp := plan
			if rp.Delivery == "" {
				rp = sess.PlanFor(res, view)
			}
			ch, werr := sess.WatchAuto(ctx, res, view, rp, p)
			if werr != nil {
				emitRouterEventID(ctx, out, id, addr, Event{Kind: KindError, At: time.Now(), Err: werr})
				return
			}
			for ev := range ch {
				if !emitRouterEventID(ctx, out, id, addr, ev) {
					return
				}
			}
		}(id)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out, nil
}

// WatchStatuses menggabungkan WatchStatus semua router.
func (m *Manager) WatchStatuses(ctx context.Context) (<-chan RouterStatus, error) {
	ids, err := m.watchIDs()
	if err != nil {
		return nil, err
	}
	out := make(chan RouterStatus, 64)
	if len(ids) == 0 {
		close(out)
		return out, nil
	}

	var wg sync.WaitGroup
	for _, id := range ids {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			sess, aerr := m.Acquire(ctx, id)
			if aerr != nil {
				return
			}
			defer m.Release(id)
			ch := sess.WatchStatus(ctx)
			for st := range ch {
				ev := RouterStatus{RouterID: id, Address: m.routerAddress(id), Status: st}
				select {
				case out <- ev:
				case <-ctx.Done():
					return
				}
			}
		}(id)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out, nil
}

// SubscribeAll membuka shared feed yang sama di semua router. Router yang
// gagal tetap dilaporkan lewat error gabungan; subscription yang berhasil
// dikembalikan dalam map.
func (m *Manager) SubscribeAll(ctx context.Context, req ros.Request) (map[string]*Subscription, error) {
	ids, err := m.watchIDs()
	if err != nil {
		return nil, err
	}
	subs := make(map[string]*Subscription, len(ids))
	var errs []error
	for _, id := range ids {
		sub, serr := m.Subscribe(ctx, id, req)
		if serr != nil {
			errs = append(errs, fmt.Errorf("%s: %w", id, serr))
			continue
		}
		subs[id] = sub
	}
	return subs, errors.Join(errs...)
}

// watchIDs mengembalikan router aktif lalu router terdaftar yang belum terbuka.
// Mengembalikan error bila manager sudah ditutup.
func (m *Manager) watchIDs() ([]string, error) {
	if m == nil {
		return nil, ErrClosed
	}
	if m.closed.Load() {
		return nil, ErrClosed
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	ids := make([]string, 0, len(m.order)+len(m.cfgOrder))
	seen := make(map[string]bool, len(ids))
	for _, id := range m.order {
		if !seen[id] {
			seen[id] = true
			ids = append(ids, id)
		}
	}
	for _, id := range m.cfgOrder {
		if !seen[id] {
			seen[id] = true
			ids = append(ids, id)
		}
	}
	return ids, nil
}

func (m *Manager) routerAddress(id string) string {
	cfg, ok := m.Config(id)
	if !ok {
		return ""
	}
	return cfg.addr()
}

// Close menutup semua router dan menunggu stream/goroutine selesai.
func (m *Manager) Close() error {
	if m == nil {
		return nil
	}
	m.closeOnce.Do(func() {
		m.closed.Store(true)
		m.mu.RLock()
		routers := make([]*MikroTik, 0, len(m.order))
		for _, id := range m.order {
			routers = append(routers, m.routers[id])
		}
		m.mu.RUnlock()

		var wg sync.WaitGroup
		errs := make([]error, len(routers))
		for i, r := range routers {
			wg.Add(1)
			go func(i int, r *MikroTik) {
				defer wg.Done()
				errs[i] = r.Close()
			}(i, r)
		}
		wg.Wait()
		m.closeErr = errors.Join(errs...)
	})
	return m.closeErr
}

func emitRouterEventID(ctx context.Context, out chan<- RouterEvent, id, address string, ev Event) bool {
	re := RouterEvent{RouterID: id, Address: address, Event: ev}
	select {
	case out <- re:
		return true
	case <-ctx.Done():
		return false
	}
}
