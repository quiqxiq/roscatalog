package rosrun

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	ros "example.com/roscatalog"
)

// ErrCacheFull dikembalikan bila jumlah feed aktif mencapai Cache.MaxEntries
// dan tidak ada entry tanpa subscriber yang bisa dibuang.
var ErrCacheFull = errors.New("rosrun: cache full")

// Subscription adalah konsumen shared feed. Banyak Subscription dengan
// Request yang sama berbagi satu stream/poll upstream ke router.
type Subscription struct {
	hub   *feedCache
	entry *cacheEntry

	ch chan Event

	closeOnce sync.Once
	chOnce    sync.Once
	doneOnce  sync.Once
	done      chan struct{}
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
	return s.entry.latest, true
}

// Close melepas subscription. Upstream dibatalkan setelah subscriber terakhir
// pergi; event terakhir tetap disimpan sampai TTL untuk subscriber berikutnya.
func (s *Subscription) Close() {
	s.closeOnce.Do(func() {
		s.signalDone()
		s.hub.unsubscribe(s)
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

	stop  context.CancelFunc
	runID uint64

	subs map[*Subscription]struct{}

	latest    Event
	latestAt  time.Time
	hasLatest bool

	lastUsed time.Time
}

type feedCache struct {
	m  *MikroTik
	mu sync.Mutex

	cfg     CacheConfig
	entries map[string]*cacheEntry
	closed  bool
}

func newFeedCache(m *MikroTik) *feedCache {
	return &feedCache{
		m:       m,
		cfg:     m.cfg.Cache,
		entries: make(map[string]*cacheEntry),
	}
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
	key := req.Fingerprint()
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
			entry = &cacheEntry{
				req:      req,
				subs:     make(map[*Subscription]struct{}),
				lastUsed: now,
			}
			c.entries[key] = entry
		}
	} else {
		entry = &cacheEntry{
			req:      req,
			subs:     make(map[*Subscription]struct{}),
			lastUsed: now,
		}
	}

	sub := &Subscription{
		hub:   c,
		entry: entry,
		ch:    make(chan Event, 256),
		done:  make(chan struct{}),
	}
	entry.subs[sub] = struct{}{}
	entry.lastUsed = now

	// Subscriber baru langsung mendapat event terakhir bila masih segar.
	if entry.hasLatest && now.Sub(entry.latestAt) <= c.cfg.TTL {
		sub.send(entry.latest)
	} else if entry.hasLatest {
		entry.hasLatest = false
		entry.latest = Event{}
	}

	if entry.stop == nil {
		entry.runID++
		runID := entry.runID
		runCtx, cancel := context.WithCancel(c.m.runCtx)
		entry.stop = cancel
		go c.run(entry, runCtx, runID)
	}

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

func (c *feedCache) run(entry *cacheEntry, ctx context.Context, runID uint64) {
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
		sub.send(ev)
	}
}

// finish menandai upstream selesai. Bila masih ada subscriber, channel mereka
// ditutup agar pemanggil tahu stream berakhir dan bisa berlangganan ulang.
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
		entry.lastUsed = time.Now()
	}
}

// pruneLocked membuang entry tanpa subscriber yang sudah melewati TTL.
func (c *feedCache) pruneLocked(now time.Time) {
	for key, entry := range c.entries {
		if len(entry.subs) == 0 && now.Sub(entry.lastUsed) > c.cfg.TTL {
			delete(c.entries, key)
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
