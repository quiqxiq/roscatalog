// Package rosrun menjalankan roscatalog.Request memakai goros/v4
// (github.com/quiqxiq/goros/v4).
//
// Hal yang sengaja ditangani di sini karena jebakan di library:
//
//   - Library mengelola .tag sendiri; jangan kirim Params.Tag.
//   - Varian *Context di library memanggil c.r.Cancel() saat ctx selesai, yang
//     mematikan pembaca bersama untuk SEMUA command pada koneksi itu. Karena itu
//     runner memakai RunArgs/ListenArgs (tanpa ctx) dan membatalkan stream lewat
//     ListenReply.Cancel() (/cancel =tag=).
//   - Konsumen yang lambat memblokir loop baca bersama (channel stream penuh =
//     semua balasan berhenti, termasuk balasan command lain). Runner menguras
//     channel library ke antrean sendiri sebelum dibaca konsumen.
//
// Runner membungkus satu *routeros.Client. Untuk pemakaian aplikasi —
// reconnect otomatis, plan stream/poll bawaan, dan lifecycle per router —
// gunakan MikroTik (session.go). New tetap tersedia untuk pemakaian satu
// koneksi dan test.
package rosrun

import (
	"context"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/quiqxiq/goros/v4"
	"github.com/quiqxiq/goros/v4/proto"

	ros "example.com/roscatalog"
)

// Kind adalah jenis event.
type Kind string

const (
	KindSnapshot Kind = "snapshot" // Rows: seluruh baris satu siklus/poll
	KindRow      Kind = "row"      // Row: satu !re (native, follow)
	KindAdded    Kind = "added"    // Row: item baru
	KindChanged  Kind = "changed"  // Row: item berubah (nilai gabungan terbaru)
	KindRemoved  Kind = "removed"  // Row: item hilang (minimal Key)
	KindCount    Kind = "count"    // Count: hasil count-only
	KindError    Kind = "error"    // Err: stream berakhir karena error
	// KindReconnecting menandai jeda: koneksi putus dan Watch akan
	// melanjutkan request yang sama setelah tersambung lagi. Tidak ada data
	// yang hilang dari sisi pemanggil selain selama jeda ini.
	KindReconnecting Kind = "reconnecting"
)

// Event adalah satu kejadian dari Watch.
type Event struct {
	Kind  Kind
	At    time.Time
	Row   map[string]string
	Rows  []map[string]string
	Count int
	Err   error
}

// Runner menjalankan Request di atas satu *routeros.Client.
type Runner struct {
	c *routeros.Client

	// Logf opsional, untuk fallback proplist dan kejadian tak biasa.
	Logf func(format string, args ...any)

	// MaxBuffered membatasi antrean per stream (bawaan 65536 kalimat). Bila
	// terlampaui stream dibatalkan dan Watch mengirim KindError.
	MaxBuffered int

	// EmptySnapshots: pada stream interval, tabel kosong tidak mengirim apa pun.
	// Bila true, runner mengirim Snapshot kosong setelah 3x Every tanpa baris.
	// Ini inferensi, bukan konfirmasi router.
	EmptySnapshots bool

	// Gate opsional, dipanggil sekali per command yang menuju router (satu
	// kali per RunArgs, satu kali per stream). Cocok untuk limiter per router.
	Gate Gate

	// limiter membatasi command in-flight dan mencatat channel stream.
	// Nil berarti tanpa batas (mis. Runner yang dibuat manual lewat New).
	limiter *limiter

	// Watchdog stream: channel yang tidak mengirim apa pun selama streamStale
	// dibuka ulang. streamCheck adalah jarak pemeriksaan.
	streamStale time.Duration
	streamCheck time.Duration
	streamEvery time.Duration
}

// New membungkus client yang sudah login dan menyalakan mode async
// (wajib agar Run dan Listen bisa berjalan bersamaan).
func New(c *routeros.Client) *Runner {
	r := &Runner{c: c, streamStale: DefaultStreamStale, streamCheck: DefaultStreamCheck}
	if !c.IsAsync() {
		errC := c.Async()
		go func() {
			for err := range errC {
				r.logf("rosrun: loop async berhenti: %v", err)
			}
		}()
	}
	return r
}

// newRunner membangun Runner dari Config tanpa menyentuh mode async: koneksi
// sudah diaktifkan oleh dial, jadi Runner hanya mewarisi setelan.
func newRunner(c *routeros.Client, cfg Config, lim *limiter) *Runner {
	return &Runner{
		c:              c,
		Logf:           cfg.Logf,
		MaxBuffered:    cfg.MaxBuffered,
		EmptySnapshots: cfg.EmptySnapshots,
		Gate:           cfg.Gate,
		limiter:        lim,
		streamEvery:    cfg.StreamEvery,
		streamStale:    cfg.StreamStale,
		streamCheck:    cfg.StreamCheck,
	}
}

func (r *Runner) logf(format string, args ...any) {
	if r.Logf != nil {
		r.Logf(format, args...)
	}
}

func (r *Runner) maxBuffered() int {
	if r.MaxBuffered > 0 {
		return r.MaxBuffered
	}
	return 1 << 16
}

func (r *Runner) staleEvery() time.Duration {
	if r.streamStale > 0 {
		return r.streamStale
	}
	return DefaultStreamStale
}

func (r *Runner) checkEvery() time.Duration {
	if r.streamCheck > 0 {
		return r.streamCheck
	}
	return DefaultStreamCheck
}

// ---------------------------------------------------------------------------
// API publik
// ---------------------------------------------------------------------------

// Once menjalankan Request satu kali dan mengembalikan baris (sudah disaring
// ke Fields). Delivery Request harus Once.
func (r *Runner) Once(ctx context.Context, req ros.Request) ([]map[string]string, error) {
	if req.Delivery != ros.DeliveryOnce {
		return nil, fmt.Errorf("rosrun.Once: Delivery %q, butuh once", req.Delivery)
	}
	rep, err := r.run(ctx, req.Sentence)
	if err != nil {
		if e := r.retryNoProplist(ctx, req, err, &rep); e != nil {
			return nil, e
		}
	}
	return rows(rep, req), nil
}

// Count menjalankan view count dan mengembalikan jumlahnya.
func (r *Runner) Count(ctx context.Context, req ros.Request) (int, error) {
	if !req.CountOnly {
		return 0, errors.New("rosrun.Count: Request bukan count-only")
	}
	rep, err := r.run(ctx, req.Sentence)
	if err != nil {
		return 0, err
	}
	return countOf(rep)
}

// Exec menjalankan sentence aksi tulis (dari Action.Sentence) dan mengembalikan
// nilai =ret= (mis. *ID dari add).
func (r *Runner) Exec(ctx context.Context, sentence []string) (string, error) {
	rep, err := r.run(ctx, sentence)
	if err != nil {
		return "", err
	}
	if rep.Done != nil {
		return rep.Done.Map["ret"], nil
	}
	return "", nil
}

// Watch menjalankan Request sesuai Delivery-nya dan mengirim Event sampai ctx
// dibatalkan, stream berakhir, atau terjadi error (KindError lalu channel
// ditutup). Konsumen harus terus membaca channel atau membatalkan ctx.
func (r *Runner) Watch(ctx context.Context, req ros.Request) (<-chan Event, error) {
	if len(req.Sentence) == 0 {
		return nil, errors.New("rosrun.Watch: Request kosong")
	}
	if req.Diff && req.Key == "" {
		return nil, errors.New("rosrun.Watch: Diff butuh Key (tabel dengan .id atau KeyBy)")
	}
	if req.Delivery == ros.DeliveryPoll && req.Every <= 0 {
		return nil, errors.New("rosrun.Watch: poll butuh Every")
	}

	out := make(chan Event, 64)
	emit := func(ev Event) bool {
		ev.At = time.Now()
		select {
		case out <- ev:
			return true
		case <-ctx.Done():
			return false
		}
	}
	go func() {
		defer close(out)
		var err error
		switch req.Delivery {
		case ros.DeliveryOnce:
			err = r.watchOnce(ctx, req, emit)
		case ros.DeliveryPoll:
			err = r.poll(ctx, req, emit)
		case ros.DeliveryStream, ros.DeliveryEvents:
			err = r.stream(ctx, req, emit)
		default:
			err = fmt.Errorf("Delivery tidak dikenal: %q", req.Delivery)
		}
		if err != nil && ctx.Err() == nil {
			emit(Event{Kind: KindError, Err: err})
		}
	}()
	return out, nil
}

// ---------------------------------------------------------------------------
// Eksekusi dasar
// ---------------------------------------------------------------------------

// run menjalankan RunArgs di goroutine agar bisa ditinggalkan saat ctx selesai
// tanpa memakai varian Context milik library.
func (r *Runner) run(ctx context.Context, sentence []string) (*routeros.Reply, error) {
	if len(sentence) == 0 {
		return nil, errors.New("rosrun: empty sentence")
	}
	var releaseGate func()
	if r.Gate != nil {
		var err error
		releaseGate, err = r.Gate(ctx, sentence[0])
		if err != nil {
			return nil, err
		}
	}
	var releaseLim func()
	if r.limiter != nil {
		r.limiter.note(sentence[0])
		releaseLim = r.limiter.acquire()
	}

	type result struct {
		rep *routeros.Reply
		err error
	}
	ch := make(chan result, 1)
	go func() {
		rep, err := r.c.RunArgs(sentence)
		// Lepas slot setelah command benar-benar selesai di wire, bukan saat
		// pemanggil menyerah karena ctx.
		if releaseLim != nil {
			releaseLim()
		}
		if releaseGate != nil {
			releaseGate()
		}
		ch <- result{rep, err}
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case x := <-ch:
		return x.rep, x.err
	}
}

// retryNoProplist dipakai Once: bila router menolak, coba lagi tanpa .proplist.
func (r *Runner) retryNoProplist(ctx context.Context, req ros.Request, first error, out **routeros.Reply) error {
	var de *routeros.DeviceError
	if req.Fallback == nil || !errors.As(first, &de) {
		return first
	}
	r.logf("rosrun: %s ditolak (%v); mengulang tanpa .proplist", req.Sentence[0], first)
	rep, err := r.run(ctx, req.Fallback)
	if err != nil {
		return err
	}
	*out = rep
	return nil
}

func (r *Runner) watchOnce(ctx context.Context, req ros.Request, emit func(Event) bool) error {
	rep, err := r.run(ctx, req.Sentence)
	if err != nil {
		if e := r.retryNoProplist(ctx, req, err, &rep); e != nil {
			return e
		}
	}
	return deliverResult(req, rep, emit)
}

func deliverResult(req ros.Request, rep *routeros.Reply, emit func(Event) bool) error {
	if req.CountOnly {
		n, err := countOf(rep)
		if err != nil {
			return err
		}
		emit(Event{Kind: KindCount, Count: n})
		return nil
	}
	rs := rows(rep, req)
	if req.Shape == ros.ShapeRows {
		for _, row := range rs {
			if !emit(Event{Kind: KindRow, Row: row}) {
				return nil
			}
		}
		return nil
	}
	emit(Event{Kind: KindSnapshot, Rows: rs})
	return nil
}

func countOf(rep *routeros.Reply) (int, error) {
	if rep == nil || rep.Done == nil {
		return 0, errors.New("count-only tanpa !done")
	}
	ret, ok := rep.Done.Map["ret"]
	if !ok {
		return 0, errors.New("count-only: =ret= tidak ada")
	}
	return strconv.Atoi(ret)
}

func rows(rep *routeros.Reply, req ros.Request) []map[string]string {
	out := make([]map[string]string, 0, len(rep.Re))
	for _, s := range rep.Re {
		out = append(out, trim(s.Map, req))
	}
	return out
}

// trim menyisakan Key dan Fields (bila Fields kosong, salin semua).
func trim(row map[string]string, req ros.Request) map[string]string {
	if len(req.Fields) == 0 {
		out := make(map[string]string, len(row))
		for k, v := range row {
			out[k] = v
		}
		return out
	}
	out := make(map[string]string, len(req.Fields)+1)
	if req.Key != "" {
		if v, ok := row[req.Key]; ok {
			out[req.Key] = v
		}
	}
	for _, f := range req.Fields {
		if v, ok := row[f]; ok {
			out[f] = v
		}
	}
	return out
}

// ---------------------------------------------------------------------------
// Poll: Go mengulang one-shot
// ---------------------------------------------------------------------------

// pollState menampung keadaan diff antar-siklus poll. Disimpan per subscription
// agar satu scheduler bisa menjalankan banyak poll tanpa kehilangan baseline.
type pollState struct {
	prev         map[string]map[string]string
	first        bool
	fallbackUsed bool
}

func newPollState() *pollState { return &pollState{first: true} }

func (r *Runner) poll(ctx context.Context, req ros.Request, emit func(Event) bool) error {
	t := time.NewTicker(req.Every)
	defer t.Stop()
	st := newPollState()
	for {
		if err := r.pollOnce(ctx, req, st, emit); err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return nil
		case <-t.C:
		}
	}
}

// pollOnce menjalankan satu siklus poll. Scheduler memanggil ini ketika
// subscription sudah jatuh tempo.
func (r *Runner) pollOnce(ctx context.Context, req ros.Request, st *pollState, emit func(Event) bool) error {
	sentence := req.Sentence
	if st.fallbackUsed && req.Fallback != nil {
		sentence = req.Fallback
	}
	rep, err := r.run(ctx, sentence)
	if err != nil {
		var de *routeros.DeviceError
		if st.first && !st.fallbackUsed && req.Fallback != nil && errors.As(err, &de) {
			r.logf("rosrun: poll %s ditolak (%v); mengulang tanpa .proplist", sentence[0], err)
			st.fallbackUsed = true
			return nil
		}
		return err
	}

	switch {
	case req.CountOnly:
		if err := deliverResult(req, rep, emit); err != nil {
			return err
		}
	case req.Diff:
		cur := make(map[string]map[string]string, len(rep.Re))
		var order []string
		for _, s := range rep.Re {
			row := trim(s.Map, req)
			k := row[req.Key]
			if _, dup := cur[k]; !dup {
				order = append(order, k)
			}
			cur[k] = row
		}
		switch {
		case st.first && req.Baseline:
			// hanya baseline
		case st.first:
			if !emit(Event{Kind: KindSnapshot, Rows: orderedRows(cur, order)}) {
				return nil
			}
		default:
			for _, k := range order {
				row := cur[k]
				old, existed := st.prev[k]
				switch {
				case !existed:
					if !emit(Event{Kind: KindAdded, Row: row}) {
						return nil
					}
				case !reflect.DeepEqual(old, row):
					if !emit(Event{Kind: KindChanged, Row: row}) {
						return nil
					}
				}
			}
			for k, old := range st.prev {
				if _, still := cur[k]; !still {
					if !emit(Event{Kind: KindRemoved, Row: old}) {
						return nil
					}
				}
			}
		}
		st.prev = cur
	default:
		if err := deliverResult(req, rep, emit); err != nil {
			return err
		}
	}
	st.first = false
	return nil
}

func orderedRows(m map[string]map[string]string, order []string) []map[string]string {
	out := make([]map[string]string, 0, len(order))
	for _, k := range order {
		out = append(out, m[k])
	}
	return out
}

// ---------------------------------------------------------------------------
// Stream: router mengirim terus
// ---------------------------------------------------------------------------

// queue menampung kalimat dari library tanpa pernah memblokir loop baca bersama.
type queue struct {
	mu       sync.Mutex
	items    []*proto.Sentence
	closed   bool
	overflow bool
	max      int
	notify   chan struct{}
}

func newQueue(max int) *queue { return &queue{max: max, notify: make(chan struct{}, 1)} }

func (q *queue) push(s *proto.Sentence) {
	q.mu.Lock()
	if q.closed || q.overflow {
		q.mu.Unlock()
		return
	}
	if len(q.items) >= q.max {
		q.overflow = true
		q.items = nil
		q.mu.Unlock()
		q.signal()
		return
	}
	q.items = append(q.items, s)
	q.mu.Unlock()
	q.signal()
}

func (q *queue) close() {
	q.mu.Lock()
	q.closed = true
	q.mu.Unlock()
	q.signal()
}

func (q *queue) signal() {
	select {
	case q.notify <- struct{}{}:
	default:
	}
}

type popResult int

const (
	popGot popResult = iota
	popClosed
	popTimeout
	popCtx
	popOverflow
)

// pop menunggu kalimat berikut. timeout <= 0 berarti tanpa batas waktu.
func (q *queue) pop(ctx context.Context, timeout time.Duration) (*proto.Sentence, popResult) {
	var timer <-chan time.Time
	if timeout > 0 {
		tm := time.NewTimer(timeout)
		defer tm.Stop()
		timer = tm.C
	}
	for {
		q.mu.Lock()
		if q.overflow {
			q.mu.Unlock()
			return nil, popOverflow
		}
		if len(q.items) > 0 {
			s := q.items[0]
			q.items = q.items[1:]
			q.mu.Unlock()
			return s, popGot
		}
		if q.closed {
			q.mu.Unlock()
			return nil, popClosed
		}
		q.mu.Unlock()
		select {
		case <-q.notify:
		case <-timer:
			return nil, popTimeout
		case <-ctx.Done():
			return nil, popCtx
		}
	}
}

// errStreamStale menandai channel yang berhenti mengirim. stream membuka
// ulang channel itu; ini bukan error yang diteruskan ke konsumen.
var errStreamStale = errors.New("rosrun: stream stale")

func (r *Runner) stream(ctx context.Context, req ros.Request, emit func(Event) bool) error {
	sentence := req.Sentence
	fallbackUsed := false

	if len(sentence) == 0 {
		return errors.New("rosrun: empty sentence")
	}
	if r.Gate != nil {
		release, err := r.Gate(ctx, sentence[0])
		if err != nil {
			return err
		}
		if release != nil {
			defer release()
		}
	}
	if r.limiter != nil {
		r.limiter.note(sentence[0])
		defer r.limiter.streamOpened()()
	}

	for {
		l, err := r.c.ListenArgs(sentence)
		if err != nil {
			return err
		}
		q := newQueue(r.maxBuffered())
		go func() { // kuras channel library agar loop baca bersama tidak terblokir
			for s := range l.Chan() {
				q.push(s)
			}
			q.close()
		}()
		got, cerr := r.consume(ctx, req, q, emit)
		if ctx.Err() != nil || q.overflowed() || cerr != nil {
			r.cancel(l) // pembatalan eksplisit: pasti terkirim sebelum Watch menutup channel
		}
		if q.overflowed() {
			return fmt.Errorf("rosrun: antrean stream penuh (>%d kalimat); konsumen terlalu lambat", r.maxBuffered())
		}
		if errors.Is(cerr, errStreamStale) {
			// Watchdog: channel diam terlalu lama. Buka ulang channel yang sama.
			r.logf("rosrun: stream %s diam > %s; membuka ulang", sentence[0], r.staleEvery())
			continue
		}
		if cerr != nil {
			return cerr
		}
		if ctx.Err() != nil {
			return nil
		}
		// stream berakhir: baca alasannya
		if err := l.Err(); err != nil {
			var de *routeros.DeviceError
			if !got && !fallbackUsed && req.Fallback != nil && errors.As(err, &de) {
				r.logf("rosrun: stream %s ditolak (%v); mengulang tanpa .proplist", sentence[0], err)
				sentence, fallbackUsed = req.Fallback, true
				continue
			}
			return err
		}
		return nil
	}
}

// cancel mengirim /cancel =tag= dan menunggu balasannya (maksimal 2 detik).
func (r *Runner) cancel(l *routeros.ListenReply) {
	done := make(chan struct{})
	go func() {
		defer close(done)
		if _, err := l.Cancel(); err != nil && !errors.Is(err, io.EOF) &&
			!strings.Contains(err.Error(), "loop has ended") {
			r.logf("rosrun: cancel gagal: %v", err)
		}
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		r.logf("rosrun: cancel melewati batas 2s")
	}
}

func (q *queue) overflowed() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.overflow
}

// consume membaca antrean sesuai Shape. got = ada data yang sudah diterima.
func (r *Runner) consume(ctx context.Context, req ros.Request, q *queue, emit func(Event) bool) (got bool, err error) {
	last := time.Now()
	switch req.Shape {
	case ros.ShapeRows:
		for {
			s, res := q.pop(ctx, r.checkEvery())
			if res == popGot {
				last = time.Now()
			}
			if res == popTimeout {
				if time.Since(last) >= r.staleEvery() {
					return got, errStreamStale
				}
				continue
			}
			if res != popGot {
				return got, nil
			}
			got = true
			if !emit(Event{Kind: KindRow, Row: trim(s.Map, req)}) {
				return got, nil
			}
		}

	case ros.ShapeChanges:
		state := map[string]map[string]string{}
		if req.Seed {
			rep, err := r.run(ctx, req.SeedSentence)
			if err != nil {
				return got, fmt.Errorf("seed: %w", err)
			}
			var order []string
			for _, s := range rep.Re {
				k := s.Map[req.Key]
				if _, dup := state[k]; !dup {
					order = append(order, k)
				}
				state[k] = s.Map
			}
			snap := make([]map[string]string, 0, len(order))
			for _, k := range order {
				snap = append(snap, trim(state[k], req))
			}
			if !emit(Event{Kind: KindSnapshot, Rows: snap}) {
				return got, nil
			}
		}
		// Listen TIDAK memakai watchdog diam: tidak ada event adalah keadaan
		// normal, bukan channel rusak. Kesehatan koneksi dipantau lewat status.
		for {
			s, res := q.pop(ctx, 0)
			if res != popGot {
				return got, nil
			}
			got = true
			k := s.Map[req.Key]
			if dead := s.Map[".dead"]; dead == "true" || dead == "yes" {
				old, existed := state[k]
				delete(state, k)
				row := trim(old, req)
				if !existed {
					row = map[string]string{req.Key: k}
				}
				if !emit(Event{Kind: KindRemoved, Row: row}) {
					return got, nil
				}
				continue
			}
			old, existed := state[k]
			merged := make(map[string]string, len(old)+len(s.Map))
			for kk, v := range old {
				merged[kk] = v
			}
			for kk, v := range s.Map {
				merged[kk] = v
			}
			state[k] = merged
			kind := KindAdded
			if existed {
				kind = KindChanged
			}
			if !emit(Event{Kind: kind, Row: trim(merged, req)}) {
				return got, nil
			}
		}

	default: // ShapeSnapshot: siklus interval
		return r.consumeCycles(ctx, req, q, emit)
	}
}

// consumeCycles mengelompokkan !re stream interval menjadi snapshot. Router
// tidak menandai batas siklus, jadi dipakai tiga petunjuk: nilai .section
// berganti, Key yang sama muncul lagi, atau jeda tanpa baris (Every/5,
// dibatasi 20ms..300ms).
func (r *Runner) consumeCycles(ctx context.Context, req ros.Request, q *queue, emit func(Event) bool) (bool, error) {
	idle := req.Every / 5
	if idle < 20*time.Millisecond {
		idle = 20 * time.Millisecond
	}
	if idle > 300*time.Millisecond {
		idle = 300 * time.Millisecond
	}
	emptyAfter := 3 * req.Every
	if emptyAfter < time.Second {
		emptyAfter = time.Second
	}

	var (
		cur      []map[string]string
		seen     = map[string]struct{}{}
		section  string
		got      bool
		wasEmpty bool
		last     = time.Now()
	)
	flush := func() bool {
		if len(cur) == 0 {
			return true
		}
		snap := cur
		cur, seen, section = nil, map[string]struct{}{}, ""
		wasEmpty = false
		return emit(Event{Kind: KindSnapshot, Rows: snap})
	}

	for {
		timeout := time.Duration(0)
		switch {
		case len(cur) > 0:
			timeout = idle
		case r.EmptySnapshots && !wasEmpty:
			timeout = emptyAfter
		}
		if timeout == 0 && r.checkEvery() > 0 {
			// Tanpa baris yang menunggu flush, tetap bangun berkala supaya
			// channel yang diam bisa dideteksi.
			timeout = r.checkEvery()
		}
		s, res := q.pop(ctx, timeout)
		switch res {
		case popTimeout:
			if time.Since(last) >= r.staleEvery() {
				flush()
				return got, errStreamStale
			}
			if len(cur) > 0 {
				if !flush() {
					return got, nil
				}
			} else if !wasEmpty {
				wasEmpty = true
				if !emit(Event{Kind: KindSnapshot, Rows: []map[string]string{}}) {
					return got, nil
				}
			}
			continue
		case popGot:
			last = time.Now()
		default:
			flush()
			return got, nil
		}
		got = true

		raw := s.Map
		if sec, ok := raw[".section"]; ok {
			if len(cur) > 0 && sec != section {
				if !flush() {
					return got, nil
				}
			}
			section = sec
		} else if req.Key != "" {
			if k, ok := raw[req.Key]; ok {
				if _, dup := seen[k]; dup {
					if !flush() {
						return got, nil
					}
				}
				seen[k] = struct{}{}
			}
		}
		cur = append(cur, trim(raw, req))
	}
}
