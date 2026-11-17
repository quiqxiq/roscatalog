package roscatalog

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestStructure(t *testing.T) {
	groups := map[Group]bool{}
	for _, g := range GroupOrder {
		groups[g] = true
	}
	ids := map[string]bool{}
	for _, r := range All() {
		if r.ID == "" || ids[r.ID] {
			t.Errorf("ID kosong atau duplikat: %q", r.ID)
		}
		ids[r.ID] = true
		if !groups[r.Group] {
			t.Errorf("%s: grup %q tidak ada di GroupOrder", r.ID, r.Group)
		}
		if !strings.HasPrefix(r.Path, "/") || strings.ContainsAny(r.Path, " \t") {
			t.Errorf("%s: path tidak valid %q", r.ID, r.Path)
		}
		if len(r.Views) == 0 && len(r.Actions) == 0 {
			t.Errorf("%s: tanpa view maupun aksi", r.ID)
		}

		seenView := map[string]bool{}
		for _, v := range r.Views {
			key := v.Name + "|" + string(v.Version)
			if v.Name == "" || seenView[key] {
				t.Errorf("%s: nama view kosong/duplikat %q", r.ID, key)
			}
			seenView[key] = true
			if !strings.HasPrefix(v.Command, "/") {
				t.Errorf("%s/%s: command tidak valid %q", r.ID, v.Name, v.Command)
			}
			switch v.Mode {
			case ModeOneShot, ModeNative, ModeInterval, ModeFollow, ModeListen:
			default:
				t.Errorf("%s/%s: mode tidak dikenal %q", r.ID, v.Name, v.Mode)
			}
			if v.Mode == ModeInterval && v.Interval < MinPrintInterval {
				t.Errorf("%s/%s: interval %s < %s", r.ID, v.Name, v.Interval, MinPrintInterval)
			}
			if v.Mode == ModeListen {
				hasList := false
				for _, o := range r.Views {
					if o.Name == "list" {
						hasList = true
					}
				}
				if !hasList {
					t.Errorf("%s/%s: listen butuh view list untuk kolom dan seed", r.ID, v.Name)
				}
			}
			if v.Mode == ModeListen && v.ServerProplist {
				t.Errorf("%s/%s: listen tidak boleh memakai ServerProplist", r.ID, v.Name)
			}
			if v.ServerProplist && len(v.Proplist) == 0 {
				t.Errorf("%s/%s: ServerProplist tanpa Proplist", r.ID, v.Name)
			}
			dup := map[string]bool{}
			for _, p := range v.Proplist {
				if dup[p] {
					t.Errorf("%s/%s: properti ganda %q", r.ID, v.Name, p)
				}
				dup[p] = true
			}
		}

		seenAct := map[string]bool{}
		for _, a := range r.Actions {
			if a.Name == "" || seenAct[a.Name] {
				t.Errorf("%s: nama aksi kosong/duplikat %q", r.ID, a.Name)
			}
			seenAct[a.Name] = true
			if !strings.HasPrefix(a.Command, "/") {
				t.Errorf("%s/%s: command aksi tidak valid %q", r.ID, a.Name, a.Command)
			}
			if a.Name == "remove" && !a.Destructive {
				t.Errorf("%s: remove harus Destructive", r.ID)
			}
		}
	}
}

// View "list" dan "live" tidak boleh membawa secret atau isi besar.
func TestNoSensitiveInRoutineViews(t *testing.T) {
	banned := map[string]bool{
		"password": true, "secret": true, "otp-secret": true, "contents": true,
		"source": true, "on-login": true, "on-logout": true, "on-event": false,
	}
	for _, r := range All() {
		for _, v := range r.Views {
			if v.Name != "list" && v.Name != "live" {
				continue
			}
			for _, p := range v.Proplist {
				if banned[p] {
					t.Errorf("%s/%s membawa properti sensitif/berat %q", r.ID, v.Name, p)
				}
			}
		}
	}
}

func TestViewSentences(t *testing.T) {
	get := func(id, view string, ver Version) View {
		t.Helper()
		r, ok := Lookup(id)
		if !ok {
			t.Fatalf("resource %s tidak ada", id)
		}
		v, ok := r.View(view, ver)
		if !ok {
			t.Fatalf("%s: view %s tidak ada", id, view)
		}
		return v
	}

	cases := []struct {
		name string
		v    View
		p    Params
		want []string
	}{
		{
			"listen mentah (proplist diisi saat Build)",
			get("hotspot-active", "events", AnyVersion),
			Params{Tag: "t1"},
			[]string{"/ip/hotspot/active/listen", ".tag=t1"},
		},
		{
			"interval dengan query dan tag",
			get("hotspot-active", "live", AnyVersion),
			Params{Tag: "t2", Query: []Query{{Prop: "server", Value: "hs1"}}},
			[]string{
				"/ip/hotspot/active/print", "=interval=2s",
				"=.proplist=.id,user,address,uptime,session-time-left,idle-time,bytes-in,bytes-out",
				"?server=hs1", ".tag=t2",
			},
		},
		{
			"interval stats milidetik",
			get("interface", "live", AnyVersion),
			Params{Interval: 500 * time.Millisecond},
			[]string{
				"/interface/print", "=interval=500ms", "=stats=",
				"=.proplist=.id,name,running,rx-byte,tx-byte,rx-packet,tx-packet,rx-drop,tx-drop,tx-queue-drop,rx-error,tx-error",
			},
		},
		{
			"count-only dengan filter",
			get("hotspot-user", "count", AnyVersion),
			Params{Query: []Query{{Prop: "profile", Value: "1jam"}}},
			[]string{"/ip/hotspot/user/print", "=count-only=", "?profile=1jam"},
		},
		{
			"singleton tanpa .id",
			get("system-resource", "live", AnyVersion),
			Params{},
			[]string{
				"/system/resource/print", "=interval=1s",
				"=.proplist=uptime,cpu-load,free-memory,total-memory,free-hdd-space,total-hdd-space",
			},
		},
		{
			"native monitor-traffic dengan once",
			get("interface-monitor-traffic", "run", AnyVersion),
			Params{Args: map[string]string{"interface": "ether1", "once": ""}},
			[]string{
				"/interface/monitor-traffic", "=interface=ether1", "=once=",
				"=.proplist=name,rx-bits-per-second,tx-bits-per-second,rx-packets-per-second,tx-packets-per-second,rx-drops-per-second,tx-drops-per-second,rx-errors-per-second,tx-errors-per-second,tx-queue-drops-per-second",
			},
		},
		{
			"native monitor ethernet",
			get("interface-ethernet", "monitor", AnyVersion),
			Params{Args: map[string]string{"numbers": "sfp1", "interval": "1s"}},
			[]string{
				"/interface/ethernet/monitor", "=numbers=sfp1", "=interval=1s",
				"=.proplist=name,status,rate,full-duplex,auto-negotiation,tx-flow-control,rx-flow-control,sfp-module-present,sfp-type,sfp-vendor-name,sfp-wavelength,sfp-temperature,sfp-supply-voltage,sfp-tx-bias-current,sfp-tx-power,sfp-rx-power",
			},
		},
		{
			"follow-only",
			get("log", "tail-only", AnyVersion),
			Params{},
			[]string{"/log/print", "=follow-only=", "=.proplist=.id,time,topics,message"},
		},
	}
	for _, c := range cases {
		got, err := c.v.Sentence(c.p)
		if err != nil {
			t.Errorf("%s: %v", c.name, err)
			continue
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s:\n got  %q\n want %q", c.name, got, c.want)
		}
	}
}

func TestViewSentenceErrors(t *testing.T) {
	r, _ := Lookup("hotspot-active")
	live, _ := r.View("live", AnyVersion)
	if _, err := live.Sentence(Params{Interval: 100 * time.Millisecond}); err == nil {
		t.Error("interval 100ms seharusnya ditolak")
	}
	if _, err := live.Sentence(Params{Args: map[string]string{"x": "y"}}); err == nil {
		t.Error("Args pada view non-native seharusnya ditolak")
	}

	ping, _ := Lookup("tool-ping")
	pv, _ := ping.View("run", AnyVersion)
	if _, err := pv.Sentence(Params{}); err == nil {
		t.Error("ping tanpa address seharusnya ditolak")
	}
	if _, err := pv.Sentence(Params{Args: map[string]string{"address": "8.8.8.8", "typo": "1"}}); err == nil {
		t.Error("argumen tidak dikenal seharusnya ditolak")
	}
	got, err := pv.Sentence(Params{Args: map[string]string{"address": "8.8.8.8", "count": "4"}, Tag: "p"})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"/ping", "=address=8.8.8.8", "=count=4",
		"=.proplist=seq,host,size,ttl,time,status,sent,received,packet-loss,min-rtt,avg-rtt,max-rtt",
		".tag=p",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ping:\n got  %q\n want %q", got, want)
	}
}

func TestVersionedViews(t *testing.T) {
	r, _ := Lookup("system-health")
	v7, ok := r.View("list", V7)
	if !ok || v7.Proplist[0] != "name" {
		t.Errorf("view v7 salah: %+v", v7)
	}
	v6, ok := r.View("list", V6)
	if !ok || v6.Proplist[0] != "voltage" {
		t.Errorf("view v6 salah: %+v", v6)
	}
	route, _ := Lookup("ip-route")
	if v, _ := route.View("list", V7); !contains(v.Proplist, "routing-table") || contains(v.Proplist, "routing-mark") {
		t.Errorf("ip-route v7 salah: %v", v.Proplist)
	}
	if v, _ := route.View("list", V6); !contains(v.Proplist, "routing-mark") || contains(v.Proplist, "routing-table") {
		t.Errorf("ip-route v6 salah: %v", v.Proplist)
	}
}

func TestActionSentences(t *testing.T) {
	user, _ := Lookup("hotspot-user")
	set, _ := user.Action("set")
	got, err := set.Sentence("*1", map[string]string{"disabled": "yes", "comment": "x"}, "")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"/ip/hotspot/user/set", "=.id=*1", "=comment=x", "=disabled=yes"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("set:\n got  %q\n want %q", got, want)
	}
	if _, err := set.Sentence("", nil, ""); err == nil {
		t.Error("set tanpa selector seharusnya ditolak")
	}

	add, _ := user.Action("add")
	if _, err := add.Sentence("", map[string]string{"password": "x"}, ""); err == nil {
		t.Error("add tanpa name seharusnya ditolak")
	}
	if u := add.Unknown(map[string]string{"name": "a", "typo": "b"}); !reflect.DeepEqual(u, []string{"typo"}) {
		t.Errorf("Unknown = %v", u)
	}

	fw, _ := Lookup("fw-filter")
	mv, _ := fw.Action("move")
	got, err = mv.Sentence("*5", map[string]string{"destination": "*2"}, "")
	if err != nil {
		t.Fatal(err)
	}
	want = []string{"/ip/firewall/filter/move", "=numbers=*5", "=destination=*2"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("move:\n got  %q\n want %q", got, want)
	}

	sys, _ := Lookup("system-identity")
	idset, _ := sys.Action("set")
	if idset.Selector != SelNone {
		t.Error("set pada singleton tidak boleh memakai selector")
	}
	scr, _ := Lookup("system-script")
	run, _ := scr.Action("run")
	got, _ = run.Sentence("backup-harian", nil, "")
	if !reflect.DeepEqual(got, []string{"/system/script/run", "=number=backup-harian"}) {
		t.Errorf("run: %q", got)
	}
}

func TestInspectAndNeverSeen(t *testing.T) {
	got := InspectChild("ip", "hotspot", "active")
	want := []string{"/console/inspect", "=request=child", "=path=ip,hotspot,active"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("inspect = %q", got)
	}
	rows := []map[string]string{{".id": "*1", "name": "a"}, {".id": "*2", "name": "b", "comment": "c"}}
	if m := NeverSeen([]string{".id", "name", "comment", "uptime"}, rows); !reflect.DeepEqual(m, []string{"uptime"}) {
		t.Errorf("NeverSeen = %v", m)
	}
}

// Setiap resource yang ditandai Verified harus punya minimal satu view baca.
func TestVerifiedHaveViews(t *testing.T) {
	n := 0
	for _, r := range All() {
		if !r.Verified {
			continue
		}
		n++
		if len(r.Views) == 0 {
			t.Errorf("%s: Verified tanpa view", r.ID)
		}
	}
	if n == 0 {
		t.Error("tidak ada resource Verified")
	}
}

func contains(s []string, x string) bool {
	for _, v := range s {
		if v == x {
			return true
		}
	}
	return false
}
