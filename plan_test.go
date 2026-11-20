package roscatalog

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func mustBuild(t *testing.T, id, view string, ver Version, plan Plan, p Params) Request {
	t.Helper()
	r, ok := Lookup(id)
	if !ok {
		t.Fatalf("resource %s tidak ada", id)
	}
	req, err := r.Build(view, ver, plan, p)
	if err != nil {
		t.Fatalf("%s/%s %+v: %v", id, view, plan, err)
	}
	return req
}

func TestBuildInterval(t *testing.T) {
	const props = "=.proplist=.id,user,address,uptime,session-time-left,idle-time,bytes-in,bytes-out"

	// Stream bawaan: interval dari katalog.
	req := mustBuild(t, "hotspot-active", "live", AnyVersion, Stream(), Params{})
	want := []string{"/ip/hotspot/active/print", "=interval=2s", "=stats=", props}
	if !reflect.DeepEqual(req.Sentence, want) || req.Shape != ShapeSnapshot || req.Delivery != DeliveryStream || req.Mode != ModeInterval {
		t.Errorf("stream bawaan: %+v", req)
	}
	if req.Every != 2*time.Second || req.Key != ".id" {
		t.Errorf("every/key: %v %q", req.Every, req.Key)
	}
	// Fallback tanpa proplist.
	if !reflect.DeepEqual(req.Fallback, []string{"/ip/hotspot/active/print", "=interval=2s", "=stats="}) {
		t.Errorf("fallback: %q", req.Fallback)
	}
	if !reflect.DeepEqual(req.Fields, []string{"user", "address", "uptime", "session-time-left", "idle-time", "bytes-in", "bytes-out"}) {
		t.Errorf("fields: %q", req.Fields)
	}

	// Stream dengan interval 2 detik.
	req = mustBuild(t, "hotspot-active", "live", AnyVersion, StreamEvery(2*time.Second), Params{})
	if req.Sentence[1] != "=interval=2s" || req.Every != 2*time.Second {
		t.Errorf("interval 2s: %q", req.Sentence)
	}

	// Poll: one-shot tanpa interval, diulang Go.
	req = mustBuild(t, "hotspot-active", "live", AnyVersion, Poll(5*time.Second), Params{})
	want = []string{"/ip/hotspot/active/print", "=stats=", props}
	if !reflect.DeepEqual(req.Sentence, want) || req.Delivery != DeliveryPoll || req.Mode != ModeOneShot ||
		req.Every != 5*time.Second || req.Diff || req.Shape != ShapeSnapshot {
		t.Errorf("poll: %+v", req)
	}
	// Poll + diff.
	req = mustBuild(t, "hotspot-active", "live", AnyVersion, Poll(5*time.Second).WithDiff(), Params{})
	if !req.Diff {
		t.Error("diff tidak terbawa")
	}

	// stats ikut terbawa ke poll.
	req = mustBuild(t, "interface", "live", AnyVersion, Poll(5*time.Second), Params{})
	if !contains(req.Sentence, "=stats=") || contains(req.Sentence, "=interval=") {
		t.Errorf("poll stats: %q", req.Sentence)
	}
}

func TestBuildNative(t *testing.T) {
	// Stream native apa adanya.
	req := mustBuild(t, "interface-monitor-traffic", "run", AnyVersion, Stream(),
		Params{Args: map[string]string{"interface": "ether1"}})
	if req.Shape != ShapeRows || req.Mode != ModeNative || contains(req.Sentence, "=once=") {
		t.Errorf("native stream: %+v", req)
	}
	if req.Key != "name" {
		t.Errorf("key monitor-traffic = %q", req.Key)
	}

	// Poll native = command yang sama dengan =once=.
	req = mustBuild(t, "interface-monitor-traffic", "run", AnyVersion, Poll(5*time.Second),
		Params{Args: map[string]string{"interface": "ether1,ether2"}})
	if !contains(req.Sentence, "=once=") || req.Delivery != DeliveryPoll || req.Mode != ModeNative {
		t.Errorf("native poll: %+v", req)
	}

	// Once native.
	req = mustBuild(t, "interface-ethernet", "monitor", AnyVersion, Once(),
		Params{Args: map[string]string{"numbers": "sfp1"}})
	if !contains(req.Sentence, "=once=") || req.Delivery != DeliveryOnce {
		t.Errorf("native once: %+v", req)
	}

	// Stream native dengan interval (monitor menerima interval).
	req = mustBuild(t, "interface-ethernet", "monitor", AnyVersion, StreamEvery(2*time.Second),
		Params{Args: map[string]string{"numbers": "sfp1"}})
	if !contains(req.Sentence, "=interval=2s") {
		t.Errorf("native interval: %q", req.Sentence)
	}

	// Argumen yang tidak dimutasi: Params.Args milik pemanggil tetap utuh.
	args := map[string]string{"interface": "ether1"}
	mustBuild(t, "interface-monitor-traffic", "run", AnyVersion, Poll(5*time.Second), Params{Args: args})
	if len(args) != 1 {
		t.Errorf("Params.Args termutasi: %v", args)
	}

	// ping tidak punya once: poll ditolak, stream boleh.
	r, _ := Lookup("tool-ping")
	if _, err := r.Build("run", AnyVersion, Poll(5*time.Second), Params{Args: map[string]string{"address": "1.1.1.1"}}); err == nil {
		t.Error("poll ping seharusnya ditolak")
	}
	req = mustBuild(t, "tool-ping", "run", AnyVersion, Stream(), Params{Args: map[string]string{"address": "1.1.1.1", "count": "5"}})
	if req.Shape != ShapeRows {
		t.Errorf("ping: %+v", req)
	}
}

func TestBuildListen(t *testing.T) {
	// Stream listen: tanpa .proplist (di RouterOS 7.24.6 proplist menekan
	// event .dead=true); kolom tetap disaring di klien dari view list.
	// Stream() pada view listen dinormalkan menjadi Events.
	req := mustBuild(t, "hotspot-active", "events", AnyVersion, Stream(), Params{})
	if req.Delivery != DeliveryEvents || req.Shape != ShapeChanges || req.Mode != ModeListen ||
		req.Sentence[0] != "/ip/hotspot/active/listen" {
		t.Errorf("listen: %+v", req)
	}

	req = mustBuild(t, "hotspot-active", "events", AnyVersion, Events(), Params{})
	if req.Delivery != DeliveryEvents || req.Shape != ShapeChanges {
		t.Errorf("events eksplisit: %+v", req)
	}
	for _, w := range req.Sentence {
		if strings.HasPrefix(w, "=.proplist=") {
			t.Errorf("listen tidak boleh mengirim .proplist: %q", req.Sentence)
		}
	}
	if req.Fallback != nil {
		t.Errorf("listen tidak perlu fallback proplist: %q", req.Fallback)
	}
	if len(req.Fields) == 0 || req.Fields[0] != "user" {
		t.Errorf("fields listen: %v", req.Fields)
	}
	if req.Seed {
		t.Error("seed tidak diminta")
	}

	// Seed: snapshot awal lewat print list.
	req = mustBuild(t, "hotspot-active", "events", AnyVersion, Stream().WithSeed(), Params{})
	if !req.Seed || req.SeedSentence[0] != "/ip/hotspot/active/print" {
		t.Errorf("seed: %+v", req)
	}

	// Poll listen = list berulang + diff.
	req = mustBuild(t, "hotspot-active", "events", AnyVersion, Poll(5*time.Second), Params{})
	if req.Sentence[0] != "/ip/hotspot/active/print" || !req.Diff || req.Mode != ModeOneShot {
		t.Errorf("poll listen: %+v", req)
	}

	// View ber-versi wajib menyebut versi.
	r, _ := Lookup("ip-route")
	if _, err := r.Build("events", AnyVersion, Stream(), Params{}); err == nil {
		t.Error("ip-route tanpa versi seharusnya ditolak")
	}
	req = mustBuild(t, "ip-route", "events", V7, Events(), Params{})
	if req.Sentence[0] != "/ip/route/listen" || !contains(req.Fields, "routing-table") {
		t.Errorf("ip-route v7: sentence=%q fields=%v", req.Sentence, req.Fields)
	}

	// Events hanya untuk view listen.
	if ar, ok := Lookup("hotspot-active"); ok {
		if _, err := ar.Build("list", AnyVersion, Events(), Params{}); err == nil {
			t.Error("Events pada view list seharusnya ditolak")
		}
	}
	if contains(req.Fields, "routing-mark") {
		t.Errorf("ip-route v7 memakai field v6: %v", req.Fields)
	}
}

func TestBuildFollow(t *testing.T) {
	req := mustBuild(t, "log", "tail-only", AnyVersion, Stream(), Params{})
	want := []string{"/log/print", "=follow-only=", "=.proplist=.id,time,topics,message"}
	if !reflect.DeepEqual(req.Sentence, want) || req.Shape != ShapeRows {
		t.Errorf("follow-only stream: %+v", req)
	}

	req = mustBuild(t, "log", "tail-only", AnyVersion, Poll(5*time.Second), Params{})
	if contains(req.Sentence, "=follow-only=") || !req.Baseline || !req.Diff {
		t.Errorf("follow-only poll: %+v", req)
	}
	req = mustBuild(t, "log", "tail", AnyVersion, Poll(5*time.Second), Params{})
	if req.Baseline || !req.Diff {
		t.Errorf("follow poll: %+v", req)
	}
}

func TestBuildOneShotAndErrors(t *testing.T) {
	req := mustBuild(t, "hotspot-user", "list", AnyVersion, Plan{}, Params{})
	if req.Delivery != DeliveryOnce || req.Shape != ShapeSnapshot {
		t.Errorf("oneshot bawaan: %+v", req)
	}
	req = mustBuild(t, "hotspot-user", "count", AnyVersion, Once(), Params{Query: []Query{{Prop: "profile", Value: "1jam"}}})
	if !req.CountOnly || req.Sentence[len(req.Sentence)-1] != "?profile=1jam" {
		t.Errorf("count: %+v", req)
	}
	req = mustBuild(t, "hotspot-user", "list", AnyVersion, Poll(5*time.Second), Params{})
	if req.Delivery != DeliveryPoll {
		t.Errorf("poll oneshot: %+v", req)
	}

	r, _ := Lookup("hotspot-user")
	if _, err := r.Build("list", AnyVersion, Stream(), Params{}); err == nil {
		t.Error("stream pada view oneshot seharusnya ditolak")
	}
	if _, err := r.Build("list", AnyVersion, Poll(10*time.Millisecond), Params{}); err == nil {
		t.Error("poll 10ms seharusnya ditolak")
	}
	if _, err := r.Build("list", AnyVersion, Poll(31*time.Second), Params{}); err == nil {
		t.Error("poll 31s seharusnya ditolak")
	}
	if _, err := r.Build("tidak-ada", AnyVersion, Once(), Params{}); err == nil {
		t.Error("view tidak ada seharusnya ditolak")
	}
	live, _ := Lookup("hotspot-active")
	if _, err := live.Build("live", AnyVersion, StreamEvery(100*time.Millisecond), Params{}); err == nil {
		t.Error("interval 100ms seharusnya ditolak")
	}
	if _, err := live.Build("live", AnyVersion, StreamEvery(6*time.Second), Params{}); err == nil {
		t.Error("interval 6s seharusnya ditolak")
	}
}

func TestRequestFingerprint(t *testing.T) {
	a := mustBuild(t, "hotspot-active", "live", AnyVersion, Stream(), Params{})
	b := mustBuild(t, "hotspot-active", "live", AnyVersion, Stream(), Params{})
	if a.Fingerprint() != b.Fingerprint() {
		t.Error("fingerprint tidak stabil untuk Request yang sama")
	}
	c := mustBuild(t, "hotspot-active", "live", AnyVersion, Stream(),
		Params{Query: []Query{{Prop: "server", Value: "hs1"}}})
	if a.Fingerprint() == c.Fingerprint() {
		t.Error("query berbeda berbagi fingerprint")
	}
	d := mustBuild(t, "hotspot-active", "live", AnyVersion, Poll(5*time.Second), Params{})
	if a.Fingerprint() == d.Fingerprint() {
		t.Error("stream dan poll berbagi fingerprint")
	}
}

func TestDeliveriesTable(t *testing.T) {
	cases := map[string][]Delivery{
		"hotspot-user/list":             {DeliveryOnce, DeliveryPoll},
		"hotspot-active/live":           {DeliveryStream, DeliveryPoll},
		"hotspot-active/events":         {DeliveryEvents, DeliveryPoll},
		"log/tail":                      {DeliveryStream, DeliveryPoll},
		"interface-monitor-traffic/run": {DeliveryStream, DeliveryOnce, DeliveryPoll},
		"tool-ping/run":                 {DeliveryStream},
		"interface-ethernet/monitor":    {DeliveryStream, DeliveryOnce, DeliveryPoll},
	}
	for k, want := range cases {
		parts := strings.SplitN(k, "/", 2)
		r, _ := Lookup(parts[0])
		if got := r.Deliveries(parts[1], AnyVersion); !reflect.DeepEqual(got, want) {
			t.Errorf("%s: got %v want %v", k, got, want)
		}
	}
}
