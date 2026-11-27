package roscatalog

import "time"

func (r *Resource) keyField() string {
	if r.Singleton {
		return ""
	}
	return ".id"
}

// P adalah singkatan untuk daftar string.
func P(s ...string) []string { return s }

// Table membuat resource menu bertabel (punya .id per baris).
func Table(id string, g Group, title, path string) *Resource {
	return &Resource{ID: id, Group: g, Title: title, Path: path}
}

// Single membuat resource menu tanpa baris (settings/singleton).
func Single(id string, g Group, title, path string) *Resource {
	return &Resource{ID: id, Group: g, Title: title, Path: path, Singleton: true}
}

// Tool membuat resource yang hanya berupa satu command native (ping, torch, ...).
func Tool(id string, g Group, title, command string, args []Arg, fields []string, selfTerminating bool) *Resource {
	r := &Resource{ID: id, Group: g, Title: title, Path: command, Singleton: true}
	return r.add(View{
		Name: "run", Mode: ModeNative, Command: command,
		Args: args, Proplist: fields, ServerProplist: len(fields) > 0,
		SelfTerminating: selfTerminating,
	})
}

// ---------------------------------------------------------------------------
// Metadata
// ---------------------------------------------------------------------------

// Describe mengisi catatan resource.
func (r *Resource) Describe(s string) *Resource { r.Notes = s; return r }

// Check menandai resource sudah dicocokkan dengan CLI reference RouterOS 7.26.
func (r *Resource) Check() *Resource { r.Verified = true; return r }

// Filter mendaftarkan properti yang aman dipakai sebagai query ?prop=value.
func (r *Resource) Filter(props ...string) *Resource { r.Filters = props; return r }

// ---------------------------------------------------------------------------
// View
// ---------------------------------------------------------------------------

func (r *Resource) add(v View) *Resource {
	r.Views = append(r.Views, v)
	return r
}

func (r *Resource) last() *View {
	if len(r.Views) == 0 {
		return nil
	}
	return &r.Views[len(r.Views)-1]
}

// MarkHeavy menandai view terakhir sebagai tabel besar.
func (r *Resource) MarkHeavy() *Resource {
	if v := r.last(); v != nil {
		v.Heavy = true
	}
	return r
}

// UseServerProplist mengaktifkan =.proplist= pada view terakhir (untuk command
// native yang terbukti menerimanya, mis. monitor-traffic).
func (r *Resource) UseServerProplist() *Resource {
	if v := r.last(); v != nil {
		v.ServerProplist = true
	}
	return r
}

// WithoutServerProplist menonaktifkan =.proplist= pada view terakhir.
// Berguna untuk resource singleton atau menu tertentu yang bermasalah bila
// proplist dikirim ke router (mis. /system/resource pada RouterOS 6).
func (r *Resource) WithoutServerProplist() *Resource {
	if v := r.last(); v != nil {
		v.ServerProplist = false
	}
	return r
}

// Hint mengisi catatan pada view terakhir.
func (r *Resource) Hint(s string) *Resource {
	if v := r.last(); v != nil {
		v.Note = s
	}
	return r
}

// PrintFlags mendeklarasikan flag print yang didukung view terakhir.
// Sumbernya adalah /console/inspect; Params.Flags divalidasi terhadap daftar ini.
func (r *Resource) PrintFlags(flags ...string) *Resource {
	if v := r.last(); v != nil {
		v.AllowedFlags = append([]string(nil), flags...)
	}
	return r
}

// Print menambah view oneshot print dengan flag tertentu. AllowedFlags diisi
// dari flags yang diberikan; pemakai bisa menggantinya dengan PrintFlags.
func (r *Resource) Print(name string, props []string, flags ...string) *Resource {
	return r.add(View{
		Name: name, Mode: ModeOneShot, Command: r.Path + "/print",
		Flags:          append([]string(nil), flags...),
		AllowedFlags:   append([]string(nil), flags...),
		Proplist:       append([]string(nil), props...),
		ServerProplist: len(props) > 0,
		WithID:         !r.Singleton,
		Key:            r.keyField(),
	})
}

// Brief, Detail, Stats, dan StatsDetail adalah preset Print untuk varian print
// yang paling sering dipakai.
func (r *Resource) Brief(props ...string) *Resource {
	return r.Print("brief", props, "brief")
}

func (r *Resource) Detail(props ...string) *Resource {
	return r.Print("detail", props, "detail")
}

func (r *Resource) Stats(props ...string) *Resource {
	return r.Print("stats", props, "stats")
}

func (r *Resource) StatsDetail(props ...string) *Resource {
	return r.Print("stats-detail", props, "stats-detail")
}

// Named menambah view oneshot (print) dengan nama dan proplist tertentu.
func (r *Resource) Named(name string, props ...string) *Resource {
	return r.add(View{
		Name: name, Mode: ModeOneShot, Command: r.Path + "/print",
		Proplist: props, ServerProplist: true, WithID: !r.Singleton, Key: r.keyField(),
	})
}

// List adalah view oneshot standar bernama "list".
func (r *Resource) List(props ...string) *Resource { return r.Named("list", props...) }

// Count membaca jumlah baris saja (=count-only=; hasil di =ret=N pada !done).
func (r *Resource) Count() *Resource {
	return r.add(View{
		Name: "count", Mode: ModeOneShot, Command: r.Path + "/print",
		Flags: P("count-only"), Note: "jumlah ada di =ret=N pada !done",
	})
}

// Live membuat view interval: print [stats] interval=every.
func (r *Resource) Live(every time.Duration, stats bool, props ...string) *Resource {
	v := View{
		Name: "live", Mode: ModeInterval, Command: r.Path + "/print",
		Interval: every, Proplist: props, ServerProplist: true, WithID: !r.Singleton,
		Key: r.keyField(),
	}
	if stats {
		v.Flags = P("stats")
	}
	return r.add(v)
}

// Events membuat view listen (event perubahan item, tanpa snapshot awal).
// ServerProplist sengaja tidak dipakai: di RouterOS 7.24.6 proplist pada
// listen menekan event .dead=true. Kolom diambil dari view list saat Build
// dan disaring di klien.
func (r *Resource) Events() *Resource {
	return r.add(View{
		Name: "events", Mode: ModeListen, Command: r.Path + "/listen",
		Key:  r.keyField(),
		Note: "tanpa snapshot awal; .dead=true berarti item hilang; kolom diambil dari view list saat Build",
	})
}

// Tail membuat view follow: snapshot lalu entri baru (print follow).
func (r *Resource) Tail(props ...string) *Resource {
	return r.add(View{
		Name: "tail", Mode: ModeFollow, Command: r.Path + "/print",
		Flags: P("follow"), Proplist: props, ServerProplist: true, WithID: !r.Singleton,
		Key: r.keyField(), Note: "kirim entri yang ada dulu, lalu entri baru",
	})
}

// TailOnly membuat view follow-only: hanya entri baru.
func (r *Resource) TailOnly(props ...string) *Resource {
	return r.add(View{
		Name: "tail-only", Mode: ModeFollow, Command: r.Path + "/print",
		Flags: P("follow-only"), Proplist: props, ServerProplist: true, WithID: !r.Singleton,
		Key: r.keyField(), Note: "hanya entri baru setelah command dijalankan",
	})
}

// Monitor membuat view native /path/monitor dengan argumen pemilih.
func (r *Resource) Monitor(selector string, fields ...string) *Resource {
	return r.add(View{
		Name: "monitor", Mode: ModeNative, Command: r.Path + "/monitor",
		Args: []Arg{
			{Name: selector, Required: true, Note: "item yang dipantau"},
			{Name: "once", Note: "kirim satu kali lalu selesai"},
			{Name: "duration", Note: "batasi lama streaming, mis. 10s"},
			{Name: "interval", Note: "jeda antar update"},
		},
		Proplist: fields, ServerProplist: len(fields) > 0, SelfTerminating: true,
	})
}

// KeyBy menetapkan properti kunci baris pada view terakhir (mis. "name" untuk
// monitor-traffic dengan banyak interface), dipakai snapshot dan diff.
func (r *Resource) KeyBy(prop string) *Resource {
	if v := r.last(); v != nil {
		v.Key = prop
	}
	return r
}

// ---------------------------------------------------------------------------
// Aksi
// ---------------------------------------------------------------------------

func (r *Resource) act(a Action) *Resource {
	r.Actions = append(r.Actions, a)
	return r
}

// Do menambah aksi khusus.
func (r *Resource) Do(name string, sel Selector, req, opt []string, destructive bool, note string) *Resource {
	return r.act(Action{
		Name: name, Command: r.Path + "/" + name, Selector: sel,
		Required: req, Optional: opt, Destructive: destructive, Note: note,
	})
}

// DoAt menambah aksi dengan command eksplisit (mis. command root seperti /export atau /import).
func (r *Resource) DoAt(command, name string, sel Selector, req, opt []string, destructive bool, note string) *Resource {
	return r.act(Action{
		Name: name, Command: command, Selector: sel,
		Required: req, Optional: opt, Destructive: destructive, Note: note,
	})
}

// Command mengganti Command pada aksi terakhir (aman no-op bila belum ada aksi).
func (r *Resource) Command(command string) *Resource {
	if len(r.Actions) > 0 {
		r.Actions[len(r.Actions)-1].Command = command
	}
	return r
}

// ReadDo menambah perintah baca (print, get, find, ...) sebagai action.
// flags adalah flag yang didukung perintah itu; ArgFlag otomatis ditambahkan.
func (r *Resource) ReadDo(name string, flags, req, opt []string, note string) *Resource {
	args := make([]Arg, 0, len(flags)+len(req)+len(opt))
	for _, f := range flags {
		args = append(args, Arg{Name: f, Kind: ArgFlag})
	}
	for _, p := range req {
		args = append(args, Arg{Name: p})
	}
	for _, p := range opt {
		args = append(args, Arg{Name: p})
	}
	return r.act(Action{
		Name: name, Command: r.Path + "/" + name, Selector: SelNone,
		Required: req, Optional: opt, Read: true, Args: args, Note: note,
	})
}

// Add menambah aksi add.
func (r *Resource) Add(req, opt []string) *Resource {
	return r.Do("add", SelNone, req, opt, false, "balasan sukses berisi =ret=*ID")
}

// Set menambah aksi set. Pada singleton tanpa selector, pada tabel dengan .id.
func (r *Resource) Set(props ...string) *Resource {
	sel := SelID
	if r.Singleton {
		sel = SelNone
	}
	return r.Do("set", sel, nil, props, false, "")
}

// Remove menambah aksi remove.
func (r *Resource) Remove(note string) *Resource {
	return r.Do("remove", SelID, nil, nil, true, note)
}

// Toggle menambah enable dan disable.
func (r *Resource) Toggle() *Resource {
	return r.Do("enable", SelID, nil, nil, false, "").
		Do("disable", SelID, nil, nil, false, "")
}

// CRUD menambah add, set, remove, enable, disable.
func (r *Resource) CRUD(req, opt []string) *Resource {
	all := append(append([]string{}, req...), opt...)
	return r.Add(req, opt).Set(all...).Remove("").Toggle()
}

// Move menambah aksi move (urutan rule).
func (r *Resource) Move() *Resource {
	return r.Do("move", SelNumbers, P("destination"), nil, false,
		"destination berisi *ID item tujuan; item dipindah ke depannya")
}

// ResetCounters menambah reset-counters dan reset-counters-all.
func (r *Resource) ResetCounters() *Resource {
	return r.Do("reset-counters", SelNumbers, nil, nil, false, "").
		Do("reset-counters-all", SelNone, nil, nil, false, "semua item di menu")
}

// ---------------------------------------------------------------------------
// Registry
// ---------------------------------------------------------------------------

var registry []*Resource

func register(rs ...*Resource) { registry = append(registry, rs...) }

// All mengembalikan semua resource bawaan, urut menurut grup lalu ID.
func All() []*Resource {
	out := make([]*Resource, len(registry))
	for i, r := range registry {
		out[i] = r.Clone()
	}
	return sortedResources(out)
}

// Lookup mencari resource menurut ID (mengembalikan salinan aman).
func Lookup(id string) (*Resource, bool) {
	if r, ok := lookupRaw(id); ok {
		return r.Clone(), true
	}
	return nil, false
}

func lookupRaw(id string) (*Resource, bool) {
	for _, r := range registry {
		if r.ID == id {
			return r, true
		}
	}
	return nil, false
}

// ByGroup mengembalikan resource dalam satu grup.
func ByGroup(g Group) []*Resource {
	var out []*Resource
	for _, r := range All() {
		if r.Group == g {
			out = append(out, r)
		}
	}
	return out
}
