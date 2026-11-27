// Package roscatalog adalah katalog command MikroTik RouterOS API yang sudah
// dikurasi untuk tooling ISP / hotspot (Mikhmon-style): tiap menu punya view
// baca (dengan mode, command, flag, proplist) dan aksi tulis (add, set, ...).
//
// Katalog ini hanya DATA + pembuat kalimat API (word list). Ia tidak membuka
// koneksi; hasil Sentence() bisa dikirim lewat library API apa pun
// (mis. goros) atau implementasi sendiri.
//
// Empat mode baca (lihat REFERENCE.md untuk penjelasan lengkap):
//
//	ModeOneShot  print biasa, selesai dengan !done
//	ModeNative   command yang memang streaming (monitor, monitor-traffic, ping, torch, ...)
//	ModeInterval print [stats] interval=N: snapshot penuh berulang dari server
//	ModeFollow   print follow / follow-only
//	ModeListen   /path/listen: hanya event perubahan item (API only)
package roscatalog

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// Mode adalah cara data dibaca dari router.
type Mode string

const (
	ModeOneShot  Mode = "oneshot"
	ModeNative   Mode = "native"
	ModeInterval Mode = "interval"
	ModeFollow   Mode = "follow"
	ModeListen   Mode = "listen"
)

// Batas interval command yang dikirim ke router: print interval, ping
// interval, dan command native lain yang memakai =interval=.
//
// Kebijakan proyek ini lebih ketat daripada batas per-command RouterOS supaya
// beban router terkendali dan bentuk stream tetap seragam di semua menu.
const (
	MinPrintInterval = time.Second
	MaxPrintInterval = 5 * time.Second
)

// Group mengelompokkan resource untuk dokumentasi dan UI.
type Group string

const (
	GroupSystem    Group = "system"
	GroupInterface Group = "interface"
	GroupIP        Group = "ip"
	GroupFirewall  Group = "firewall"
	GroupQueue     Group = "queue"
	GroupPPP       Group = "ppp-radius"
	GroupHotspot   Group = "hotspot"
	GroupTools     Group = "tools"
)

// GroupOrder adalah urutan tampil grup.
var GroupOrder = []Group{
	GroupHotspot, GroupPPP, GroupIP, GroupFirewall, GroupQueue,
	GroupInterface, GroupSystem, GroupTools,
}

// Selector menentukan nama argumen yang dipakai untuk memilih item.
type Selector string

const (
	SelNone    Selector = ""        // add, atau set pada singleton
	SelID      Selector = ".id"     // =.id=*1A  (set/remove/enable/disable)
	SelNumbers Selector = "numbers" // =numbers=*1A (move, reset-counters, make-static, ...)
	SelNumber  Selector = "number"  // =number=nama (system script run)
)

// ArgKind membedakan argumen bernilai (=name=value) dari flag (=name=).
// Zero value ArgValue menjaga kompatibilitas dengan Arg lama.
type ArgKind uint8

const (
	ArgValue ArgKind = iota
	ArgFlag
)

// Arg adalah argumen command, dipakai view native dan action (termasuk print).
type Arg struct {
	Name     string
	Kind     ArgKind
	Required bool
	Note     string
}

// View adalah satu cara membaca data dari sebuah resource.
type View struct {
	Name    string   // list, live, events, monitor, ...
	Mode    Mode     //
	Command string   // path command API lengkap, mis. /ip/hotspot/active/print
	Flags   []string // dikirim sebagai =flag= (stats, count-only, follow, ...)

	// AllowedFlags adalah flag yang benar-benar didukung command ini (dari
	// /console/inspect). Params.Flags divalidasi terhadap daftar ini.
	AllowedFlags []string
	Args         []Arg // argumen command (khusus ModeNative)

	// Proplist adalah kolom yang benar-benar dibutuhkan. Bila ServerProplist
	// true, dikirim sebagai =.proplist= (didukung print, follow, interval,
	// dan monitor); tetap saring ulang di klien (lihat Request.Fields).
	// Listen sengaja tidak mengirim .proplist: di RouterOS 7.24.6 proplist
	// menekan event .dead=true sehingga penghapusan tidak terdeteksi.
	Proplist       []string
	ServerProplist bool
	WithID         bool // tambahkan .id otomatis ke .proplist

	Key             string        // properti kunci baris (".id" untuk tabel); dipakai snapshot/diff
	Interval        time.Duration // default untuk ModeInterval
	SelfTerminating bool          // native: berhenti sendiri bila dibatasi count/duration/once
	Heavy           bool          // tabel besar: wajib pakai filter/proplist/count-only
	Note            string
}

// Action adalah satu perintah RouterOS. Read=true untuk perintah baca
// (print, get, find, export), termasuk seluruh varian print (stats/detail/...).
// Args memuat semua argumen yang didukung perintah ini; flag bertanda
// Kind=ArgFlag. Required/Optional dipertahankan untuk kompatibilitas dengan
// builder lama.
type Action struct {
	Name        string
	Command     string
	Selector    Selector
	Required    []string
	Optional    []string
	Destructive bool
	Read        bool
	Args        []Arg
	Note        string
}

// Resource adalah satu menu (atau satu tool) RouterOS.
type Resource struct {
	ID        string
	Group     Group
	Title     string
	Path      string
	Singleton bool     // menu tanpa baris (mis. /system/resource) atau tool
	Verified  bool     // dicocokkan dengan CLI reference RouterOS 7.26
	Filters   []string // properti yang aman dipakai sebagai query ?prop=value
	Views     []View
	Actions   []Action
	Notes     string
}

// ---------------------------------------------------------------------------
// Pencarian
// ---------------------------------------------------------------------------

// View mencari view menurut nama. Nama view unik di dalam satu resource.
func (r *Resource) View(name string) (View, bool) {
	for i := range r.Views {
		if r.Views[i].Name == name {
			return r.Views[i], true
		}
	}
	return View{}, false
}

// Action mencari aksi menurut nama.
func (r *Resource) Action(name string) (Action, bool) {
	for _, a := range r.Actions {
		if a.Name == name {
			return a, true
		}
	}
	return Action{}, false
}

// ---------------------------------------------------------------------------
// Pembuat kalimat API
// ---------------------------------------------------------------------------

// Op adalah operator query word API.
type Op string

const (
	OpEq     Op = "="       // ?prop=value
	OpLt     Op = "<"       // ?<prop=value
	OpGt     Op = ">"       // ?>prop=value
	OpHas    Op = "has"     // ?prop
	OpNotHas Op = "not-has" // ?-prop
)

// Query adalah satu query word. API tidak mendukung regex (~) dan where.
type Query struct {
	Prop  string
	Op    Op
	Value string
}

func (q Query) word() (string, error) {
	if q.Prop == "" {
		return "", fmt.Errorf("%w: empty property", ErrInvalidQuery)
	}
	switch q.Op {
	case "", OpEq:
		return "?" + q.Prop + "=" + q.Value, nil
	case OpLt:
		return "?<" + q.Prop + "=" + q.Value, nil
	case OpGt:
		return "?>" + q.Prop + "=" + q.Value, nil
	case OpHas:
		return "?" + q.Prop, nil
	case OpNotHas:
		return "?-" + q.Prop, nil
	}
	return "", fmt.Errorf("%w: unknown operator %q", ErrInvalidQuery, q.Op)
}

// Params adalah masukan untuk View.Sentence.
type Params struct {
	Tag      string            // .tag untuk multiplexing dan /cancel
	Query    []Query           // filter sisi server
	Interval time.Duration     // override interval (ModeInterval)
	Args     map[string]string // argumen view native; nilai kosong = flag (=once=)
	// Flags adalah flag print tambahan (brief, detail, stats, ...).
	// Divalidasi terhadap View.AllowedFlags; konflik antarflag ditolak.
	Flags []string
	// Proplist adalah daftar kolom yang dikirim apa adanya sebagai =.proplist=.
	// TIDAK divalidasi terhadap daftar field katalog — ini jalan keluar untuk
	// struct buatan pemakai yang memuat field di luar katalog. Raw menang atas
	// Proplist, dan Proplist menang atas Only.
	Proplist []string
	// Raw meminta semua kolom apa adanya: tidak mengirim =.proplist= dan
	//.Request.Fields dikosongkan supaya rosrun tidak menyaring kolom.
	Raw bool
	// Only membatasi kolom yang diambil dari daftar field view. Nama harus
	// tercantum di Proplist view, kalau tidak Build mengembalikan error.
	// ".id" dan View.Key selalu boleh diminta karena dipakai sebagai kunci
	// snapshot/diff; keduanya tidak ikut menjadi Request.Fields.
	Only []string
}

// FormatTime mengubah durasi Go ke format waktu RouterOS (1s, 500ms).
func FormatTime(d time.Duration) string {
	if d <= 0 {
		return "0s"
	}
	if d%time.Second == 0 {
		return fmt.Sprintf("%ds", int64(d/time.Second))
	}
	ms := int64(d / time.Millisecond)
	if ms < 1 {
		ms = 1
	}
	return fmt.Sprintf("%dms", ms)
}

// exclusiveFlagGroups berisi flag print yang tidak boleh dipakai bersamaan.
// Daftar ini sengaja konservatif: hanya pasangan yang jelas kontradiktif.
var exclusiveFlagGroups = [][]string{
	{"brief", "detail", "terse"},
	{"follow", "follow-only"},
	{"stats", "count-only"},
	{"stats", "stats-detail"},
}

// checkExclusive menolak kombinasi flag yang saling eksklusif.
func checkExclusive(flags []string) error {
	seen := map[string]bool{}
	for _, f := range flags {
		seen[f] = true
	}
	for _, group := range exclusiveFlagGroups {
		n := 0
		for _, g := range group {
			if seen[g] {
				n++
			}
		}
		if n > 1 {
			return fmt.Errorf("%w: %s are mutually exclusive", ErrInvalidFlag, strings.Join(group, "/"))
		}
	}
	return nil
}

// validateFlags menolak flag yang tidak didukung command dan kombinasi flag
// yang saling eksklusif.
func validateFlags(command string, allowed, flags []string) error {
	if len(flags) == 0 {
		return nil
	}
	ok := make(map[string]bool, len(allowed))
	for _, f := range allowed {
		ok[f] = true
	}
	for _, f := range flags {
		if !ok[f] {
			return fmt.Errorf("%w: %s: flag %q", ErrUnknownFlag, command, f)
		}
	}
	return checkExclusive(flags)
}

func sortedStrings(in []string) []string {
	out := append([]string(nil), in...)
	sort.Strings(out)
	return out
}

// ProplistWords mengembalikan isi =.proplist= (dengan .id bila perlu).
func (v View) ProplistWords() []string {
	if !v.WithID {
		return append([]string(nil), v.Proplist...)
	}
	for _, p := range v.Proplist {
		if p == ".id" {
			return append([]string(nil), v.Proplist...)
		}
	}
	return append([]string{".id"}, v.Proplist...)
}

// Sentence menyusun word list API untuk view ini.
func (v View) Sentence(p Params) ([]string, error) {
	if strings.TrimSpace(v.Command) == "" {
		return nil, fmt.Errorf("%w: empty view command", ErrInvalidResource)
	}
	w := []string{v.Command}

	switch v.Mode {
	case ModeNative:
		declared := make(map[string]bool, len(v.Args))
		for _, a := range v.Args {
			declared[a.Name] = true
			if _, ok := p.Args[a.Name]; a.Required && !ok {
				return nil, fmt.Errorf("%w: %s: %q", ErrMissingArgument, v.Command, a.Name)
			}
		}
		for _, k := range sortedKeys(p.Args) {
			if !declared[k] {
				return nil, fmt.Errorf("%w: %s: unknown argument %q", ErrInvalidQuery, v.Command, k)
			}
		}
		for _, a := range v.Args {
			if val, ok := p.Args[a.Name]; ok {
				if a.Name == "interval" && val != "" {
					if d, err := time.ParseDuration(val); err == nil {
						if d < MinPrintInterval || d > MaxPrintInterval {
							return nil, fmt.Errorf("%w: %s: interval %s out of range %s..%s",
								ErrInvalidInterval, v.Command, d, MinPrintInterval, MaxPrintInterval)
						}
					}
				}
				w = append(w, "="+a.Name+"="+val)
			}
		}
	case ModeInterval:
		d := p.Interval
		if d == 0 {
			d = v.Interval
		}
		if d < MinPrintInterval {
			return nil, fmt.Errorf("%w: %s: interval %s below minimum %s", ErrInvalidInterval, v.Command, d, MinPrintInterval)
		}
		if d > MaxPrintInterval {
			return nil, fmt.Errorf("%w: %s: interval %s above maximum %s", ErrInvalidInterval, v.Command, d, MaxPrintInterval)
		}
		w = append(w, "=interval="+FormatTime(d))
		fallthrough
	default:
		if len(p.Args) > 0 && v.Mode != ModeNative {
			return nil, fmt.Errorf("%w: %s: Args only supported on native views", ErrInvalidQuery, v.Command)
		}
	}

	for _, f := range v.Flags {
		w = append(w, "="+f+"=")
	}
	if len(p.Flags) > 0 {
		if err := validateFlags(v.Command, v.AllowedFlags, p.Flags); err != nil {
			return nil, err
		}
		if err := checkExclusive(append(append([]string(nil), v.Flags...), p.Flags...)); err != nil {
			return nil, fmt.Errorf("%s: %w", v.Command, err)
		}
		have := make(map[string]bool, len(v.Flags))
		for _, f := range v.Flags {
			have[f] = true
		}
		for _, f := range sortedStrings(p.Flags) {
			if !have[f] {
				w = append(w, "="+f+"=")
			}
		}
	}
	if v.ServerProplist && len(v.Proplist) > 0 {
		w = append(w, "=.proplist="+strings.Join(v.ProplistWords(), ","))
	}
	for _, q := range p.Query {
		qw, err := q.word()
		if err != nil {
			return nil, err
		}
		w = append(w, qw)
	}
	if p.Tag != "" {
		w = append(w, ".tag="+p.Tag)
	}
	return w, nil
}

// Sentence menyusun word list API untuk aksi tulis. selector diisi dengan
// *ID, nama, atau daftar dipisah koma sesuai Action.Selector.
func (a Action) Sentence(selector string, args map[string]string, tag string) ([]string, error) {
	return a.SentenceWith(selector, nil, args, tag)
}

// FlagNames mengembalikan nama flag yang dideklarasikan pada action.
func (a Action) FlagNames() []string {
	var out []string
	for _, arg := range a.Args {
		if arg.Kind == ArgFlag {
			out = append(out, arg.Name)
		}
	}
	return out
}

// SentenceWith menyusun word list untuk action apa pun, termasuk print.
// flags harus merupakan flag yang dideklarasikan di Action.Args (Kind=ArgFlag).
func (a Action) SentenceWith(selector string, flags []string, args map[string]string, tag string) ([]string, error) {
	w := []string{a.Command}
	if a.Selector != SelNone {
		if selector == "" {
			return nil, fmt.Errorf("%w: %s (%s)", ErrSelectorRequired, a.Command, a.Selector)
		}
		w = append(w, "="+string(a.Selector)+"="+selector)
	}
	if err := validateFlags(a.Command, a.FlagNames(), flags); err != nil {
		return nil, err
	}
	for _, f := range sortedStrings(flags) {
		w = append(w, "="+f+"=")
	}
	for _, r := range a.Required {
		val, ok := args[r]
		if !ok || strings.TrimSpace(val) == "" {
			return nil, fmt.Errorf("%w: %s: %q", ErrMissingArgument, a.Command, r)
		}
	}
	for _, k := range sortedKeys(args) {
		w = append(w, "="+k+"="+args[k])
	}
	if tag != "" {
		w = append(w, ".tag="+tag)
	}
	return w, nil
}

// Unknown mengembalikan nama argumen yang tidak terdaftar di Required/Optional.
// Hanya untuk logging; Sentence tidak menolaknya karena menu firewall punya
// ratusan properti.
func (a Action) Unknown(args map[string]string) []string {
	known := make(map[string]bool)
	for _, n := range a.Required {
		known[n] = true
	}
	for _, n := range a.Optional {
		known[n] = true
	}
	var out []string
	for _, k := range sortedKeys(args) {
		if !known[k] {
			out = append(out, k)
		}
	}
	return out
}

// InspectChild menyusun kalimat /console/inspect untuk membaca anak sebuah
// path, dipakai untuk mendeteksi kemampuan menu (mis. apakah print punya
// argumen follow/interval, apakah ada command monitor).
func InspectChild(path ...string) []string {
	return []string{"/console/inspect", "=request=child", "=path=" + strings.Join(path, ",")}
}

// NeverSeen mengembalikan properti yang diminta tetapi tidak muncul di satu
// pun baris. Ini hanya indikasi ketidakcocokan versi: properti kosong (mis.
// comment yang belum diisi) memang tidak dikirim router.
func NeverSeen(requested []string, rows []map[string]string) []string {
	seen := make(map[string]bool)
	for _, row := range rows {
		for k := range row {
			seen[k] = true
		}
	}
	var out []string
	for _, p := range requested {
		if p != ".id" && !seen[p] {
			out = append(out, p)
		}
	}
	return out
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
