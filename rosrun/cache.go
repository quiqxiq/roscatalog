package rosrun

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	ros "example.com/roscatalog"
)

// ErrCacheFull dikembalikan bila jumlah feed aktif mencapai Cache.MaxEntries
// dan tidak ada entry tanpa subscriber yang bisa dibuang.
var ErrCacheFull = errors.New("rosrun: cache full")

// Subscription adalah konsumen shared feed. Banyak Subscription dengan menu
// dan perilaku yang sama berbagi satu upstream ke router. Masing-masing
// menyaring hasil ke daftar field-nya sendiri.
type Subscription struct {
	hub   *feedCache
	entry *cacheEntry

	ch chan Event

	// fields adalah proplist yang diminta subscription ini; kosong berarti
	// semua field. key selalu dipertahankan agar diff/snapshot tetap bekerja.
	fields []string
	key    string

	closeOnce   sync.Once
	releaseOnce sync.Once
	chOnce      sync.Once
	doneOnce    sync.Once
	done        chan struct{}

	// onClose dipanggil sekali saat subscription ditutup; dipakai Manager
	// untuk melepas refcount sesi.
	onClose func()
}

// fireClose menjalankan onClose sekali, baik subscription ditutup pemakai
// maupun upstream berakhir sendiri.
func (s *Subscription) fireClose() {
	s.releaseOnce.Do(func() {
		if s.onClose != nil {
			s.onClose()
		}
	})
}

// Events mengembalikan channel event. Channel ditutup saat Subscription.Close,
// sesi ditutup, atau upstream berakhir.
func (s *Subscription) Events() <-chan Event { return s.ch }

// Latest mengembalikan event data terakhir yang masih di-cache, jika ada.
func (s *Subscription) Latest() (Event, bool) {
	s.hub.mu.Lock()
	defer s.hub.mu.Unlock()
	if s.entry == nil || !s.entry.hasLatest {
		return Event{}, false
	}
	return trimEvent(s.entry.latest, s), true
}

// Close melepas subscription. Upstream dibatalkan setelah subscriber terakhir
// pergi; event terakhir tetap disimpan sampai TTL untuk subscriber berikutnya.
func (s *Subscription) Close() {
	s.closeOnce.Do(func() {
		s.signalDone()
		s.hub.unsubscribe(s)
		s.fireClose()
	})
}

func (s *Subscription) send(ev Event) {
	select {
	case s.ch <- ev:
	default:
		// Konsumen lambat: event terbaru tetap tersedia lewat Latest dan
		// upstream tidak boleh terblokir oleh satu dashboard.
	}
}

func (s *Subscription) closeChannel() {
	s.chOnce.Do(func() { close(s.ch) })
}

func (s *Subscription) signalDone() {
	s.doneOnce.Do(func() { close(s.done) })
}

type cacheEntry struct {
	req ros.Request

	// baseSentence/Fallback adalah kalimat asli tanpa perubahan proplist.
	// Union direkonstruksi dari keduanya setiap kali subscriber berubah.
	baseSentence []string
	baseFallback []string

	stop  context.CancelFunc
	runID uint64

	subs map[*Subscription]struct{}

	// fields adalah union proplist semua subscriber; allFields berarti salah
	// satu subscriber meminta semua kolom sehingga proplist tidak dikirim.
	fields    map[string]bool
	allFields bool

	latest    Event
	latestAt  time.Time
	hasLatest bool

	lastUsed time.Time

	// poll adalah keadaan diff untuk scheduler; hanya dipakai DeliveryPoll.
	poll *pollState
}

func newCacheEntry(req ros.Request, now time.Time) *cacheEntry {
	return &cacheEntry{
		req:          req,
		baseSentence: append([]string(nil), req.Sentence...),
		baseFallback: append([]string(nil), req.Fallback...),
		subs:         map[*Subscription]struct{}{},
		lastUsed:     now,
	}
}

type feedCache struct {
	m  *MikroTik
	mu sync.Mutex

	cfg     CacheConfig
	entries map[string]*cacheEntry
	closed  bool
	sched   *feedScheduler
}

func newFeedCache(m *MikroTik) *feedCache {
	c := &feedCache{
		m:       m,
		cfg:     m.cfg.Cache,
		entries: map[string]*cacheEntry{},
	}
	c.sched = newFeedScheduler(c)
	c.sched.start()
	return c
}

// Subscribe membuka atau bergabung ke shared feed untuk req. Request harus
// memakai Delivery stream atau poll.
func (m *MikroTik) Subscribe(ctx context.Context, req ros.Request) (*Subscription, error) {
	if err := m.ready(m.feed); err != nil {
		return nil, err
	}
	if req.Delivery != ros.DeliveryStream && req.Delivery != ros.DeliveryPoll {
		return nil, fmt.Errorf("rosrun.Subscribe: Delivery %q tidak didukung", req.Delivery)
	}
	if m.cache == nil {
		return nil, ErrClosed
	}
	return m.cache.subscribe(ctx, req)
}

func (c *feedCache) subscribe(ctx context.Context, req ros.Request) (*Subscription, error) {
	key := menuKey(req)
	now := time.Now()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return nil, ErrClosed
	}

	var entry *cacheEntry
	if !c.cfg.Disable {
		c.pruneLocked(now)
		entry = c.entries[key]
		if entry == nil {
			if len(c.entries) >= c.cfg.MaxEntries {
				if !c.evictLocked() {
					return nil, fmt.Errorf("%w: %d feed", ErrCacheFull, c.cfg.MaxEntries)
				}
			}
			entry = newCacheEntry(req, now)
			c.entries[key] = entry
		}
	} else {
		entry = newCacheEntry(req, now)
	}

	sub := &Subscription{
		hub:    c,
		entry:  entry,
		ch:     make(chan Event, 256),
		done:   make(chan struct{}),
		fields: append([]string(nil), req.Fields...),
		key:    req.Key,
	}
	entry.subs[sub] = struct{}{}
	entry.lastUsed = now
	c.recomputeLocked(entry)

	// Subscriber baru langsung mendapat event terakhir bila masih segar.
	if entry.hasLatest && now.Sub(entry.latestAt) <= c.cfg.TTL {
		sub.send(trimEvent(entry.latest, sub))
	} else if entry.hasLatest {
		entry.hasLatest = false
		entry.latest = Event{}
	}

	c.ensureRunningLocked(entry)

	// ctx pemanggil hanya mengikat subscription ini, bukan upstream bersama.
	go func() {
		select {
		case <-ctx.Done():
			sub.Close()
		case <-sub.done:
		}
	}()

	return sub, nil
}

// ensureRunningLocked memulai upstream sesuai delivery: stream membuka channel,
// poll didaftarkan ke scheduler per sesi.
func (c *feedCache) ensureRunningLocked(entry *cacheEntry) {
	switch entry.req.Delivery {
	case ros.DeliveryPoll:
		if entry.poll == nil {
			entry.poll = newPollState()
		}
		c.sched.add(entry)
	case ros.DeliveryStream:
		c.startStreamLocked(entry)
	}
}

func (c *feedCache) startStreamLocked(entry *cacheEntry) {
	if entry.stop != nil {
		return
	}
	entry.runID++
	runID := entry.runID
	runCtx, cancel := context.WithCancel(c.m.runCtx)
	entry.stop = cancel
	go c.runStream(entry, runCtx, runID)
}

func (c *feedCache) restartStreamLocked(entry *cacheEntry) {
	if entry.stop != nil {
		entry.stop()
		entry.stop = nil
	}
	c.startStreamLocked(entry)
}

func (c *feedCache) runStream(entry *cacheEntry, ctx context.Context, runID uint64) {
	ch, err := c.m.Watch(ctx, entry.req)
	if err != nil {
		c.publish(entry, Event{Kind: KindError, At: time.Now(), Err: err})
		c.finish(entry, runID)
		return
	}
	for ev := range ch {
		c.publish(entry, ev)
	}
	c.finish(entry, runID)
}

func (c *feedCache) publish(entry *cacheEntry, ev Event) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return
	}
	now := time.Now()
	entry.lastUsed = now
	switch ev.Kind {
	case KindError, KindReconnecting:
	default:
		entry.latest = ev
		entry.latestAt = now
		entry.hasLatest = true
	}
	for sub := range entry.subs {
		sub.send(trimEvent(ev, sub))
	}
}

// finish menandai upstream stream selesai. Bila masih ada subscriber, channel
// mereka ditutup agar pemanggil tahu stream berakhir dan bisa berlangganan ulang.
func (c *feedCache) finish(entry *cacheEntry, runID uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if entry.runID != runID {
		// Upstream generasi lama; jangan sentuh subscriber baru.
		return
	}
	entry.stop = nil
	if c.closed {
		return
	}
	for sub := range entry.subs {
		sub.closeChannel()
		sub.signalDone()
		sub.fireClose()
		delete(entry.subs, sub)
	}
	entry.lastUsed = time.Now()
}

func (c *feedCache) unsubscribe(s *Subscription) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry := s.entry
	s.closeChannel()
	if entry == nil {
		return
	}
	delete(entry.subs, s)
	if len(entry.subs) == 0 {
		if entry.stop != nil {
			entry.stop()
			entry.stop = nil
		}
		if entry.req.Delivery == ros.DeliveryPoll {
			c.sched.remove(entry)
			entry.poll = nil
		}
		entry.lastUsed = time.Now()
		return
	}
	c.recomputeLocked(entry)
}

// recomputeLocked menghitung ulang union proplist setelah subscriber berubah.
// Union yang menyempit juga boleh: subscriber datang dan pergi bersama halaman.
func (c *feedCache) recomputeLocked(entry *cacheEntry) {
	entry.allFields = false
	entry.fields = map[string]bool{}
	for sub := range entry.subs {
		if len(sub.fields) == 0 {
			entry.allFields = true
			continue
		}
		for _, f := range sub.fields {
			entry.fields[f] = true
		}
	}

	// Listen TIDAK PERNAH mengirim .proplist: di RouterOS 7.24.6 proplist
	// menekan event .dead=true. Kolomnya disaring di klien lewat trimEvent.
	if entry.req.Mode == ros.ModeListen {
		if !equalWords(entry.req.Sentence, entry.baseSentence) {
			entry.req.Sentence = append([]string(nil), entry.baseSentence...)
		}
		if len(entry.baseFallback) > 0 {
			entry.req.Fallback = append([]string(nil), entry.baseFallback...)
		}
		return
	}

	union := entry.unionLocked()
	sent := withProplist(entry.baseSentence, union, entry.allFields)
	if equalWords(sent, entry.req.Sentence) {
		return
	}
	entry.req.Sentence = sent
	if len(entry.baseFallback) > 0 {
		entry.req.Fallback = withProplist(entry.baseFallback, union, entry.allFields)
	}
	if entry.req.Delivery == ros.DeliveryStream && entry.stop != nil {
		// Field baru tidak ada di channel yang sedang berjalan; buka ulang.
		c.restartStreamLocked(entry)
	}
}

func (e *cacheEntry) unionLocked() []string {
	if e.allFields {
		return nil
	}
	set := map[string]bool{}
	if e.req.Key != "" {
		set[e.req.Key] = true
	}
	for f := range e.fields {
		set[f] = true
	}
	out := make([]string, 0, len(set))
	for f := range set {
		out = append(out, f)
	}
	sort.Strings(out)
	return out
}

// menuKey mengidentifikasi upstream bersama: command, mode, flag, dan query
// yang sama boleh berbagi satu channel; daftar proplist tidak ikut menjadi key
// karena union-lah yang dikirim.
func menuKey(req ros.Request) string {
	parts := []string{
		string(req.Delivery), string(req.Mode), string(req.Shape), req.Key,
		strconv.FormatBool(req.CountOnly), strconv.FormatBool(req.Diff),
		strconv.FormatBool(req.Baseline), strconv.FormatBool(req.Seed),
	}
	parts = append(parts, stripProplist(req.Sentence)...)
	parts = append(parts, "\x1e")
	parts = append(parts, stripProplist(req.Fallback)...)
	return strings.Join(parts, "\x1f")
}

func stripProplist(words []string) []string {
	out := make([]string, 0, len(words))
	for _, w := range words {
		if strings.HasPrefix(w, "=.proplist=") {
			continue
		}
		out = append(out, w)
	}
	return out
}

// withProplist menulis ulang kalimat dengan daftar field union. all=true
// membuang proplist sama sekali (perilaku raw).
func withProplist(sentence []string, fields []string, all bool) []string {
	out := make([]string, 0, len(sentence)+1)
	inserted := false
	for _, w := range sentence {
		if strings.HasPrefix(w, "=.proplist=") {
			continue
		}
		if !inserted && (strings.HasPrefix(w, ".tag=") || strings.HasPrefix(w, "?")) {
			if !all && len(fields) > 0 {
				out = append(out, "=.proplist="+strings.Join(fields, ","))
			}
			inserted = true
		}
		out = append(out, w)
	}
	if !inserted && !all && len(fields) > 0 {
		out = append(out, "=.proplist="+strings.Join(fields, ","))
	}
	return out
}

func equalWords(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// trimEvent menyaring satu event ke proplist subscription. Subscription dengan
// fields kosong menerima semua kolom.
func trimEvent(ev Event, s *Subscription) Event {
	if s == nil || len(s.fields) == 0 {
		return ev
	}
	trim := func(row map[string]string) map[string]string {
		out := make(map[string]string, len(s.fields)+1)
		if s.key != "" {
			if v, ok := row[s.key]; ok {
				out[s.key] = v
			}
		}
		for _, f := range s.fields {
			if v, ok := row[f]; ok {
				out[f] = v
			}
		}
		return out
	}
	switch ev.Kind {
	case KindSnapshot:
		rows := make([]map[string]string, len(ev.Rows))
		for i, row := range ev.Rows {
			rows[i] = trim(row)
		}
		ev.Rows = rows
	case KindRow, KindAdded, KindChanged, KindRemoved:
		if ev.Row != nil {
			ev.Row = trim(ev.Row)
		}
	}
	return ev
}

// pruneLocked membuang entry tanpa subscriber yang sudah melewati TTL.
func (c *feedCache) pruneLocked(now time.Time) {
	for key, entry := range c.entries {
		if len(entry.subs) == 0 && now.Sub(entry.lastUsed) > c.cfg.TTL {
			delete(c.entries, key)
			c.sched.remove(entry)
		}
	}
}

// evictLocked membuang entry tanpa subscriber yang paling lama tidak dipakai.
func (c *feedCache) evictLocked() bool {
	var oldestKey string
	var oldest time.Time
	for key, entry := range c.entries {
		if len(entry.subs) != 0 {
			continue
		}
		if oldestKey == "" || entry.lastUsed.Before(oldest) {
			oldestKey, oldest = key, entry.lastUsed
		}
	}
	if oldestKey == "" {
		return false
	}
	c.sched.remove(c.entries[oldestKey])
	delete(c.entries, oldestKey)
	return true
}

func (c *feedCache) close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return
	}
	c.closed = true
	if c.sched != nil {
		c.sched.stopAll()
	}
	for _, entry := range c.entries {
		if entry.stop != nil {
			entry.stop()
			entry.stop = nil
		}
		for sub := range entry.subs {
			sub.closeChannel()
			sub.signalDone()
			delete(entry.subs, sub)
		}
	}
	c.entries = make(map[string]*cacheEntry)
}

// pollEntry menjalankan satu siklus poll untuk scheduler.
func (c *feedCache) pollEntry(entry *cacheEntry) {
	// Salin di bawah lock: recompute bisa mengganti Sentence saat subscriber
	// lain datang/menunggu.
	c.mu.Lock()
	req := entry.req
	c.mu.Unlock()
	ctx, cancel := context.WithTimeout(c.m.runCtx, pollTimeout(req.Every))
	defer cancel()
	emit := func(ev Event) bool {
		c.publish(entry, ev)
		return true
	}
	if err := c.m.feed.run().pollOnce(ctx, req, entry.poll, emit); err != nil {
		c.m.logf("rosrun: poll %s gagal: %v", req.Sentence[0], err)
		if ctx.Err() == nil && !errors.Is(err, context.Canceled) {
			c.publish(entry, Event{Kind: KindError, At: time.Now(), Err: err})
		}
	}
}

func pollTimeout(every time.Duration) time.Duration {
	d := 2 * every
	if d < 5*time.Second {
		d = 5 * time.Second
	}
	return d
}
