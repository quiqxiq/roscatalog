package rosrun

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"

	ros "example.com/roscatalog"
)

func build(t *testing.T, id, view string, ver ros.Version, plan ros.Plan, p ros.Params) ros.Request {
	t.Helper()
	res, ok := ros.Lookup(id)
	if !ok {
		t.Fatalf("resource %s tidak ada", id)
	}
	req, err := res.Build(view, ver, plan, p)
	if err != nil {
		t.Fatalf("Build %s/%s: %v", id, view, err)
	}
	return req
}

// fastPoll membangun Request poll yang valid (5s) lalu mempercepat Every untuk
// test. Batas 5s tetap berlaku untuk pemakai API; ini hanya jalur test.
func fastPoll(t *testing.T, id, view string, ver ros.Version, every time.Duration, p ros.Params) ros.Request {
	t.Helper()
	req := build(t, id, view, ver, ros.Poll(ros.MinPollEvery), p)
	req.Every = every
	return req
}

// fastInterval membangun Request interval yang valid (1s) lalu mengganti
// =interval= dan Every agar test tidak menunggu satu detik per siklus.
func fastInterval(t *testing.T, id, view string, ver ros.Version, every time.Duration, p ros.Params) ros.Request {
	t.Helper()
	req := build(t, id, view, ver, ros.StreamEvery(ros.MinPrintInterval), p)
	req.Every = every
	rewriteInterval := func(words []string) {
		for i, w := range words {
			if strings.HasPrefix(w, "=interval=") {
				words[i] = "=interval=" + ros.FormatTime(every)
			}
		}
	}
	rewriteInterval(req.Sentence)
	rewriteInterval(req.Fallback)
	return req
}

func fastStream(t *testing.T, req ros.Request, every time.Duration) ros.Request {
	t.Helper()
	req.Every = every
	for i, w := range req.Sentence {
		if strings.HasPrefix(w, "=interval=") {
			req.Sentence[i] = "=interval=" + ros.FormatTime(every)
		}
	}
	for i, w := range req.Fallback {
		if strings.HasPrefix(w, "=interval=") {
			req.Fallback[i] = "=interval=" + ros.FormatTime(every)
		}
	}
	return req
}

// next membaca satu event atau gagal bila lewat batas waktu.
func next(t *testing.T, ch <-chan Event, d time.Duration) (Event, bool) {
	t.Helper()
	select {
	case ev, ok := <-ch:
		return ev, ok
	case <-time.After(d):
		t.Fatalf("timeout %s menunggu event", d)
	}
	return Event{}, false
}

func take(t *testing.T, ch <-chan Event, n int) []Event {
	t.Helper()
	var out []Event
	for len(out) < n {
		ev, ok := next(t, ch, 5*time.Second)
		if !ok {
			t.Fatalf("channel tertutup setelah %d event, butuh %d", len(out), n)
		}
		if ev.Kind == KindError {
			t.Fatalf("error di stream: %v", ev.Err)
		}
		out = append(out, ev)
	}
	return out
}

func drain(t *testing.T, ch <-chan Event) []Event {
	t.Helper()
	var out []Event
	deadline := time.After(5 * time.Second)
	for {
		select {
		case ev, ok := <-ch:
			if !ok {
				return out
			}
			out = append(out, ev)
		case <-deadline:
			t.Fatal("channel tidak tertutup")
		}
	}
}

func ids(rows []map[string]string) []string {
	var out []string
	for _, r := range rows {
		out = append(out, r[".id"])
	}
	return out
}

func TestOnceCountAndExec(t *testing.T) {
	m := newMock(t)
	_, r := m.client()
	ctx := context.Background()

	n, err := r.Count(ctx, build(t, "hotspot-user", "count", ros.AnyVersion, ros.Once(),
		ros.Params{Query: []ros.Query{{Prop: "profile", Value: "1jam"}}}))
	if err != nil || n != 7 {
		t.Fatalf("Count = %d, %v", n, err)
	}
	if q := m.callsTo("/ip/hotspot/user/print")[0].query; !reflect.DeepEqual(q, []string{"?profile=1jam"}) {
		t.Errorf("query terkirim: %v", q)
	}

	add, _ := mustAction(t, "hotspot-user", "add")
	sent, err := add.Sentence("", map[string]string{"name": "v1", "password": "x", "profile": "1jam"}, "")
	if err != nil {
		t.Fatal(err)
	}
	ret, err := r.Exec(ctx, sent)
	if err != nil || ret != "*A1" {
		t.Fatalf("Exec = %q, %v", ret, err)
	}
}

func mustAction(t *testing.T, id, name string) (ros.Action, *ros.Resource) {
	t.Helper()
	res, _ := ros.Lookup(id)
	a, ok := res.Action(name)
	if !ok {
		t.Fatalf("aksi %s/%s tidak ada", id, name)
	}
	return a, res
}

func TestOnceTrimsToFields(t *testing.T) {
	m := newMock(t)
	_, r := m.client()
	// print one-shot hotspot-active; router mengirim "comment" yang tidak diminta
	req := build(t, "hotspot-active", "list", ros.AnyVersion, ros.Once(), ros.Params{})
	rows, err := r.Once(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) == 0 {
		t.Fatal("tidak ada baris")
	}
	if _, has := rows[0]["comment"]; has {
		t.Errorf("kolom di luar Fields lolos: %v", rows[0])
	}
	if rows[0][".id"] != "*1" || rows[0]["user"] != "alice" {
		t.Errorf("baris: %v", rows[0])
	}
}

func TestPollDiff(t *testing.T) {
	m := newMock(t)
	_, r := m.client()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := fastPoll(t, "hotspot-active", "live", ros.AnyVersion, 60*time.Millisecond, ros.Params{})
	req.Diff = true
	ch, err := r.Watch(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	// tick1: snapshot (*1,*2). tick2: changed *1. tick3: changed *1 + removed *2. tick4: changed *1 + added *3.
	evs := take(t, ch, 6)
	if evs[0].Kind != KindSnapshot || !reflect.DeepEqual(ids(evs[0].Rows), []string{"*1", "*2"}) {
		t.Fatalf("snapshot awal: %+v", evs[0])
	}
	count := map[Kind]int{}
	for _, e := range evs[1:] {
		count[e.Kind]++
		if e.Kind == KindRemoved && e.Row[".id"] != "*2" {
			t.Errorf("removed salah: %v", e.Row)
		}
		if e.Kind == KindAdded && e.Row["user"] != "carol" {
			t.Errorf("added salah: %v", e.Row)
		}
	}
	if count[KindChanged] != 3 || count[KindRemoved] != 1 || count[KindAdded] != 1 {
		t.Errorf("hitungan event: %v", count)
	}
	cancel()
	drain(t, ch)
}

func TestPollNativeOnce(t *testing.T) {
	m := newMock(t)
	_, r := m.client()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := fastPoll(t, "interface-monitor-traffic", "run", ros.AnyVersion, 50*time.Millisecond,
		ros.Params{Args: map[string]string{"interface": "ether1"}})
	ch, _ := r.Watch(ctx, req)
	evs := take(t, ch, 3)
	for _, e := range evs {
		if e.Kind != KindRow || e.Row["name"] != "ether1" || e.Row["rx-bits-per-second"] != "1000" {
			t.Errorf("row: %+v", e)
		}
		if _, has := e.Row["rx-drops-per-second"]; !has {
			t.Errorf("kolom Fields hilang: %v", e.Row)
		}
	}
	cancel()
	drain(t, ch)
	for _, c := range m.callsTo("/interface/monitor-traffic") {
		if !c.has("once") {
			t.Error("poll native harus mengirim =once=")
		}
	}
	if m.waitStreamsStopped(time.Second) == false {
		t.Error("ada stream tersisa di router")
	}
}

func TestStreamIntervalGroupsSnapshots(t *testing.T) {
	m := newMock(t)
	_, r := m.client()
	ctx, cancel := context.WithCancel(context.Background())

	req := fastInterval(t, "hotspot-active", "live", ros.AnyVersion, 200*time.Millisecond, ros.Params{})
	ch, _ := r.Watch(ctx, req)
	evs := take(t, ch, 3)
	for i, e := range evs {
		if e.Kind != KindSnapshot || !reflect.DeepEqual(ids(e.Rows), []string{"*1", "*2"}) {
			t.Fatalf("snapshot %d: %+v", i, e)
		}
		if _, has := e.Rows[0]["comment"]; has {
			t.Errorf("kolom di luar Fields: %v", e.Rows[0])
		}
	}
	// nilai berubah antar siklus (uptime naik)
	if evs[0].Rows[0]["uptime"] == evs[2].Rows[0]["uptime"] {
		t.Errorf("uptime tidak berubah antar siklus: %v vs %v", evs[0].Rows[0]["uptime"], evs[2].Rows[0]["uptime"])
	}
	calls := m.callsTo("/ip/hotspot/active/print")
	if calls[0].attrs["interval"] != "200ms" || !strings.HasPrefix(calls[0].attrs[".proplist"], ".id,user") {
		t.Errorf("sentence terkirim: %v", calls[0].attrs)
	}

	cancel()
	for _, e := range drain(t, ch) {
		if e.Kind == KindError {
			t.Errorf("cancel tidak boleh menghasilkan error: %v", e.Err)
		}
	}
	if !m.waitStreamsStopped(time.Second) {
		t.Error("router masih mengirim setelah cancel")
	}
	if len(m.callsTo("/cancel")) == 0 {
		t.Error("/cancel tidak terkirim")
	}
}

func TestStreamListenWithSeed(t *testing.T) {
	m := newMock(t)
	_, r := m.client()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := build(t, "hotspot-active", "events", ros.AnyVersion, ros.Stream().WithSeed(), ros.Params{})
	ch, _ := r.Watch(ctx, req)
	evs := take(t, ch, 4)

	if evs[0].Kind != KindSnapshot || !reflect.DeepEqual(ids(evs[0].Rows), []string{"*1", "*2"}) {
		t.Fatalf("seed: %+v", evs[0])
	}
	if evs[1].Kind != KindAdded || evs[1].Row["user"] != "carol" {
		t.Errorf("added: %+v", evs[1])
	}
	// perubahan parsial digabung dengan state seed: user=alice tetap ada
	if evs[2].Kind != KindChanged || evs[2].Row["uptime"] != "99s" || evs[2].Row["user"] != "alice" {
		t.Errorf("changed: %+v", evs[2])
	}
	if evs[3].Kind != KindRemoved || evs[3].Row["user"] != "bob" {
		t.Errorf("removed: %+v", evs[3])
	}

	// listen tidak boleh mengirim .proplist: di RouterOS 7.24.6 proplist
	// menekan event .dead=true sehingga penghapusan tidak terdeteksi.
	listen := m.callsTo("/ip/hotspot/active/listen")
	if len(listen) != 1 || listen[0].has(".proplist") {
		t.Errorf("listen harus sekali dan tanpa .proplist: %+v", listen)
	}
	cancel()
	drain(t, ch)
	if !m.waitStreamsStopped(time.Second) {
		t.Error("listen tidak dibatalkan")
	}
}

func TestStreamNativeCancel(t *testing.T) {
	m := newMock(t)
	_, r := m.client()
	ctx, cancel := context.WithCancel(context.Background())

	req := build(t, "interface-monitor-traffic", "run", ros.AnyVersion, ros.Stream(),
		ros.Params{Args: map[string]string{"interface": "ether1,ether2"}})
	ch, _ := r.Watch(ctx, req)
	for _, e := range take(t, ch, 3) {
		if e.Kind != KindRow || e.Row["tx-bits-per-second"] != "2000" {
			t.Errorf("row: %+v", e)
		}
	}
	cancel()
	for _, e := range drain(t, ch) {
		if e.Kind == KindError {
			t.Errorf("error saat cancel: %v", e.Err)
		}
	}
	if !m.waitStreamsStopped(time.Second) {
		t.Error("stream masih jalan di router")
	}
	if c := m.callsTo("/interface/monitor-traffic"); len(c) != 1 || c[0].has("once") {
		t.Errorf("native stream salah terkirim: %+v", c)
	}
}

func TestProplistFallback(t *testing.T) {
	m := newMock(t)
	m.rejectProplist["/log/print"] = true
	_, r := m.client()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := build(t, "log", "tail-only", ros.AnyVersion, ros.Stream(), ros.Params{})
	ch, _ := r.Watch(ctx, req)
	evs := take(t, ch, 3)
	if evs[0].Kind != KindRow || evs[0].Row["message"] != "baris 1" {
		t.Errorf("row: %+v", evs[0])
	}
	// hanya Fields yang tersisa walau router (tanpa proplist) mengirim semuanya
	calls := m.callsTo("/log/print")
	if len(calls) != 2 || !calls[0].has(".proplist") || calls[1].has(".proplist") {
		t.Fatalf("urutan percobaan: %+v", calls)
	}
	if !calls[1].has("follow-only") {
		t.Error("fallback harus tetap membawa follow-only")
	}
	cancel()
	drain(t, ch)
}

func TestSlowConsumerDoesNotBlockConnection(t *testing.T) {
	m := newMock(t)
	_, r := m.client()
	r.MaxBuffered = 1 << 20
	ctx, cancel := context.WithCancel(context.Background())

	// stream banjir yang TIDAK dibaca konsumen
	flood := ros.Request{
		Delivery: ros.DeliveryStream, Mode: ros.ModeNative, Shape: ros.ShapeRows,
		Sentence: []string{"/flood"},
	}
	ch, _ := r.Watch(ctx, flood)
	time.Sleep(300 * time.Millisecond)

	// command lain di koneksi yang sama harus tetap selesai
	done := make(chan error, 1)
	go func() {
		_, err := r.Exec(context.Background(), []string{"/ip/hotspot/user/add", "=name=v2"})
		done <- err
	}()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("koneksi terblokir oleh konsumen lambat")
	}
	cancel()
	drain(t, ch)
}

func TestQueueOverflowStopsStream(t *testing.T) {
	m := newMock(t)
	_, r := m.client()
	r.MaxBuffered = 5
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	flood := ros.Request{
		Delivery: ros.DeliveryStream, Mode: ros.ModeNative, Shape: ros.ShapeRows,
		Sentence: []string{"/flood"},
	}
	ch, _ := r.Watch(ctx, flood)
	time.Sleep(400 * time.Millisecond) // konsumen tidak membaca; antrean melampaui batas
	evs := drain(t, ch)
	last := evs[len(evs)-1]
	if last.Kind != KindError || !strings.Contains(last.Err.Error(), "antrean") {
		t.Fatalf("event terakhir: %+v", last)
	}
	if !m.waitStreamsStopped(time.Second) {
		t.Error("stream tidak dibatalkan setelah overflow")
	}
}

func TestWatchValidation(t *testing.T) {
	m := newMock(t)
	_, r := m.client()
	if _, err := r.Watch(context.Background(), ros.Request{}); err == nil {
		t.Error("Request kosong seharusnya ditolak")
	}
	if _, err := r.Watch(context.Background(), ros.Request{Sentence: []string{"/x"}, Delivery: ros.DeliveryPoll}); err == nil {
		t.Error("poll tanpa Every seharusnya ditolak")
	}
	if _, err := r.Watch(context.Background(), ros.Request{Sentence: []string{"/x"}, Delivery: ros.DeliveryPoll, Every: time.Second, Diff: true}); err == nil {
		t.Error("Diff tanpa Key seharusnya ditolak")
	}
}

func TestWatchAutoFallsBackToPoll(t *testing.T) {
	m := newMock(t)
	m.rejectPaths["/ip/hotspot/active/listen"] = true // listen tidak didukung
	_, r := m.client()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	res, _ := ros.Lookup("hotspot-active")
	ch, err := r.WatchAuto(ctx, res, "events", ros.AnyVersion, ros.Plan{Every: ros.MinPollEvery}, ros.Params{})
	if err != nil {
		t.Fatal(err)
	}
	first, ok := next(t, ch, 8*time.Second)
	if !ok {
		t.Fatal("channel tertutup sebelum snapshot")
	}
	if first.Kind != KindSnapshot {
		t.Errorf("event pertama setelah fallback harus snapshot: %+v", first)
	}
	second, ok := next(t, ch, 8*time.Second)
	if !ok {
		t.Fatal("channel tertutup sebelum event diff")
	}
	if second.Kind != KindChanged && second.Kind != KindAdded && second.Kind != KindRemoved {
		t.Errorf("event berikutnya harus hasil diff: %+v", second)
	}
	if len(m.callsTo("/ip/hotspot/active/listen")) < 1 || len(m.callsTo("/ip/hotspot/active/print")) < 2 {
		t.Errorf("urutan panggilan tidak sesuai: listen=%d print=%d",
			len(m.callsTo("/ip/hotspot/active/listen")), len(m.callsTo("/ip/hotspot/active/print")))
	}
	cancel()
	drain(t, ch)
}

func TestWatchAutoKeepsStreamWhenSupported(t *testing.T) {
	m := newMock(t)
	_, r := m.client()
	ctx, cancel := context.WithCancel(context.Background())
	res, _ := ros.Lookup("hotspot-active")
	ch, err := r.WatchAuto(ctx, res, "events", ros.AnyVersion, ros.Plan{Seed: true}, ros.Params{})
	if err != nil {
		t.Fatal(err)
	}
	evs := take(t, ch, 2)
	if evs[0].Kind != KindSnapshot || evs[1].Kind != KindAdded {
		t.Errorf("stream asli: %+v %+v", evs[0], evs[1])
	}
	cancel()
	drain(t, ch)
	if len(m.callsTo("/ip/hotspot/active/listen")) != 1 {
		t.Error("listen seharusnya dipanggil sekali")
	}
}
