package rosrun

import (
	"context"
	"sort"
	"sync"
)

// limiter membatasi jumlah command yang boleh in-flight ke SATU router.
//
// Batas ini adalah plafon, bukan throttle: burst normal jauh di bawahnya,
// sedangkan kasus tanpa batas (banyak subscription berjalan bersamaan) adalah
// yang ingin dicegah. Channel stream TIDAK mengambil slot command — channel
// yang ditahan bukan command yang belum dijawab — tetapi level-nya dicatat.
type limiter struct {
	max int
	sem chan struct{}

	mu       sync.Mutex
	inflight int
	streams  int
	commands int64
	menus    map[string]int64
}

// MenuCount adalah satu baris statistik menu.
type MenuCount struct {
	Menu  string
	Count int64
}

// LimiterStats adalah snapshot instrumentasi satu router.
type LimiterStats struct {
	InFlight    int
	OpenStreams int
	Commands    int64
	Menus       []MenuCount
}

func newLimiter(max int) *limiter {
	if max <= 0 {
		max = DefaultMaxInFlight
	}
	return &limiter{max: max, sem: make(chan struct{}, max), menus: map[string]int64{}}
}

func (l *limiter) acquire() func() {
	rel, _ := l.acquireCtx(context.Background())
	return rel
}

func (l *limiter) acquireCtx(ctx context.Context) (func(), error) {
	if l == nil {
		return func() {}, nil
	}
	select {
	case l.sem <- struct{}{}:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	l.mu.Lock()
	l.inflight++
	l.mu.Unlock()

	var once sync.Once
	return func() {
		once.Do(func() {
			<-l.sem
			l.mu.Lock()
			if l.inflight > 0 {
				l.inflight--
			}
			l.mu.Unlock()
		})
	}, nil
}

func (l *limiter) note(menu string) {
	if l == nil {
		return
	}
	l.mu.Lock()
	l.commands++
	if menu != "" {
		l.menus[menu]++
	}
	l.mu.Unlock()
}

// streamOpened mencatat satu channel terbuka dan mengembalikan release
// idempotent. Tidak memblokir: ini instrumen, bukan gerbang.
func (l *limiter) streamOpened() func() {
	if l == nil {
		return func() {}
	}
	l.mu.Lock()
	l.streams++
	l.mu.Unlock()

	var once sync.Once
	return func() {
		once.Do(func() {
			l.mu.Lock()
			if l.streams > 0 {
				l.streams--
			}
			l.mu.Unlock()
		})
	}
}

// streamsGone mereset level channel ketika koneksi putus: semua channel lama
// hilang bersama socket, terlepas dari apakah collector sempat melepasnya.
func (l *limiter) streamsGone() {
	if l == nil {
		return
	}
	l.mu.Lock()
	l.streams = 0
	l.mu.Unlock()
}

func (l *limiter) stats() LimiterStats {
	if l == nil {
		return LimiterStats{}
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	st := LimiterStats{
		InFlight:    len(l.sem),
		OpenStreams: l.streams,
		Commands:    l.commands,
		Menus:       make([]MenuCount, 0, len(l.menus)),
	}
	for menu, n := range l.menus {
		st.Menus = append(st.Menus, MenuCount{Menu: menu, Count: n})
	}
	sort.Slice(st.Menus, func(i, j int) bool {
		if st.Menus[i].Count != st.Menus[j].Count {
			return st.Menus[i].Count > st.Menus[j].Count
		}
		return st.Menus[i].Menu < st.Menus[j].Menu
	})
	return st
}
