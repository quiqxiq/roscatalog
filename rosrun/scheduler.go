package rosrun

import (
	"sync"
	"time"

	ros "github.com/quiqxiq/roscatalog"
)

// feedScheduler menjalankan semua feed DeliveryPoll untuk satu sesi dari SATU
// goroutine. Stream tidak masuk ke sini: channel-nya dikelola langsung oleh
// cache.
type pollTask struct {
	entry   *cacheEntry
	state   *pollState
	next    time.Time
	running bool
}

type feedScheduler struct {
	c    *feedCache
	tick time.Duration

	mu      sync.Mutex
	tasks   map[*cacheEntry]*pollTask
	stop    chan struct{}
	done    chan struct{}
	started bool
	wg      sync.WaitGroup
}

func newFeedScheduler(c *feedCache) *feedScheduler {
	return &feedScheduler{
		c:     c,
		tick:  250 * time.Millisecond,
		tasks: map[*cacheEntry]*pollTask{},
	}
}

func (s *feedScheduler) start() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.started {
		return
	}
	s.started = true
	s.stop = make(chan struct{})
	s.done = make(chan struct{})
	go s.loop()
}

func (s *feedScheduler) loop() {
	defer close(s.done)
	t := time.NewTicker(s.tick)
	defer t.Stop()
	for {
		select {
		case <-s.stop:
			return
		case now := <-t.C:
			s.run(now)
		}
	}
}

// add mendaftarkan entry poll. Idempotent: entry yang sudah ada tidak diubah,
// kecuali state-nya hilang.
func (s *feedScheduler) add(entry *cacheEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.tasks[entry] != nil {
		return
	}
	state := entry.poll
	if state == nil {
		state = newPollState()
		entry.poll = state
	}
	s.tasks[entry] = &pollTask{entry: entry, state: state, next: time.Now()}
}

func (s *feedScheduler) remove(entry *cacheEntry) {
	s.mu.Lock()
	delete(s.tasks, entry)
	s.mu.Unlock()
}

// stopAll menghentikan loop tanpa menunggu: pemanggil bisa memegang lock cache,
// dan pollEntry yang sedang berjalan akan berhenti sendiri saat publish melihat
// cache sudah ditutup.
func (s *feedScheduler) stopAll() {
	s.mu.Lock()
	if !s.started {
		s.mu.Unlock()
		return
	}
	s.started = false
	close(s.stop)
	s.tasks = map[*cacheEntry]*pollTask{}
	s.mu.Unlock()
}

func (s *feedScheduler) wait() {
	s.wg.Wait()
}

func (s *feedScheduler) run(now time.Time) {
	s.mu.Lock()
	due := make([]*pollTask, 0, len(s.tasks))
	for _, task := range s.tasks {
		if task.running || task.next.After(now) {
			continue
		}
		task.running = true
		every := task.entry.req.Every
		if every < ros.MinPollEvery {
			every = ros.MinPollEvery
		}
		task.next = now.Add(every)
		due = append(due, task)
	}
	s.mu.Unlock()

	for _, task := range due {
		s.wg.Add(1)
		go func(task *pollTask) {
			defer s.wg.Done()
			defer func() {
				s.mu.Lock()
				task.running = false
				s.mu.Unlock()
			}()
			s.c.pollEntry(task.entry)
		}(task)
	}
}
