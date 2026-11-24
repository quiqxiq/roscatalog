package roscatalog

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Requestable adalah antarmuka untuk builder atau objek yang dapat
// menghasilkan struktur ros.Request standar.
type Requestable interface {
	Request() Request
	Build() (Request, error)
}

// PathBuilder adalah pintu masuk pembuatan command MikroTik berbasis path/menu.
type PathBuilder struct {
	path string
}

// Path membuat PathBuilder untuk menu atau submenu MikroTik.
// Pemisah spasi atau slash dinormalkan menjadi path API biner.
// Contoh: Path("ip firewall filter") atau Path("/interface").
func Path(path string) *PathBuilder {
	return &PathBuilder{path: normalizePath(path)}
}

// Menu adalah alias dari Path.
func Menu(path string) *PathBuilder {
	return Path(path)
}

func normalizePath(path string) string {
	path = strings.TrimSpace(path)
	path = strings.ReplaceAll(path, " ", "/")
	parts := strings.Split(path, "/")
	clean := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			clean = append(clean, p)
		}
	}
	if len(clean) == 0 {
		return "/"
	}
	return "/" + strings.Join(clean, "/")
}

// Print membuat builder untuk action print pada menu ini.
// Menerima switches awal (misal "stats", "detail", "without-paging").
func (p *PathBuilder) Print(switches ...string) *PrintCmdBuilder {
	b := &PrintCmdBuilder{
		path:     p.path,
		rawArgs:  make(map[string]string),
		delivery: DeliveryOnce,
		mode:     ModeOneShot,
	}
	if len(switches) > 0 {
		b.Switch(switches...)
	}
	return b
}

// Set membuat builder action set dengan selector (nomor item, nama, atau ID internal).
// Data parameter dapat dioper sekaligus lewat map, typed struct (fields), atau variadic key-value.
func (p *PathBuilder) Set(selector any, data ...any) *CmdBuilder {
	b := &CmdBuilder{
		path:     p.path,
		action:   "set",
		selector: formatSelector(selector),
		args:     make(map[string]string),
		delivery: DeliveryOnce,
		mode:     ModeOneShot,
	}
	if len(data) > 0 {
		b.Args(data...)
	}
	return b
}

// Add membuat builder action add untuk menambah item baru.
// Parameter dapat dioper lewat map, typed struct (fields), atau variadic key-value.
func (p *PathBuilder) Add(data ...any) *CmdBuilder {
	b := &CmdBuilder{
		path:     p.path,
		action:   "add",
		args:     make(map[string]string),
		delivery: DeliveryOnce,
		mode:     ModeOneShot,
	}
	if len(data) > 0 {
		b.Args(data...)
	}
	return b
}

// Remove membuat builder action remove untuk menghapus item.
func (p *PathBuilder) Remove(selector any) *CmdBuilder {
	return &CmdBuilder{
		path:     p.path,
		action:   "remove",
		selector: formatSelector(selector),
		args:     make(map[string]string),
		delivery: DeliveryOnce,
		mode:     ModeOneShot,
	}
}

// Enable membuat builder action enable untuk mengaktifkan item.
func (p *PathBuilder) Enable(selector any) *CmdBuilder {
	return &CmdBuilder{
		path:     p.path,
		action:   "enable",
		selector: formatSelector(selector),
		args:     make(map[string]string),
		delivery: DeliveryOnce,
		mode:     ModeOneShot,
	}
}

// Disable membuat builder action disable untuk menonaktifkan item.
func (p *PathBuilder) Disable(selector any) *CmdBuilder {
	return &CmdBuilder{
		path:     p.path,
		action:   "disable",
		selector: formatSelector(selector),
		args:     make(map[string]string),
		delivery: DeliveryOnce,
		mode:     ModeOneShot,
	}
}

// Comment membuat builder action comment untuk memberi komentar pada item.
func (p *PathBuilder) Comment(selector any, text string) *CmdBuilder {
	b := &CmdBuilder{
		path:     p.path,
		action:   "comment",
		selector: formatSelector(selector),
		args:     map[string]string{"comment": text},
		delivery: DeliveryOnce,
		mode:     ModeOneShot,
	}
	return b
}

// Reset membuat builder action reset untuk mengembalikan item ke nilai bawaan.
func (p *PathBuilder) Reset(selector any) *CmdBuilder {
	return &CmdBuilder{
		path:     p.path,
		action:   "reset",
		selector: formatSelector(selector),
		args:     make(map[string]string),
		delivery: DeliveryOnce,
		mode:     ModeOneShot,
	}
}

// Do membuat builder untuk action atau child command apa pun di menu ini
// (misal "monitor-traffic", "check-for-updates", "run", "ping").
func (p *PathBuilder) Do(action string, uparamAndArgs ...any) *CmdBuilder {
	b := &CmdBuilder{
		path:     p.path,
		action:   action,
		args:     make(map[string]string),
		delivery: DeliveryOnce,
		mode:     ModeOneShot,
	}
	if len(uparamAndArgs) > 0 {
		first := uparamAndArgs[0]
		// Jika parameter pertama adalah string skalar (bukan map/struct), jadikan uparam/selector
		if s, ok := first.(string); ok && len(uparamAndArgs) > 1 {
			b.selector = s
			b.Args(uparamAndArgs[1:]...)
		} else if isDataInput(first) {
			b.Args(uparamAndArgs...)
		} else {
			b.selector = formatSelector(first)
			if len(uparamAndArgs) > 1 {
				b.Args(uparamAndArgs[1:]...)
			}
		}
	}
	return b
}

func isDataInput(v any) bool {
	if v == nil {
		return false
	}
	switch v.(type) {
	case map[string]string, map[string]any:
		return true
	default:
		val := reflect.ValueOf(v)
		for val.Kind() == reflect.Pointer {
			if val.IsNil() {
				return false
			}
			val = val.Elem()
		}
		return val.Kind() == reflect.Struct
	}
}

// ========================================================================
// PrintCmdBuilder: Builder khusus untuk print query & streaming
// ========================================================================

type PrintCmdBuilder struct {
	path       string
	switches   map[string]bool
	proplist   []string
	interval   time.Duration
	follow     bool
	followOnly bool
	from       string
	file       string
	queries    []Query
	rawArgs    map[string]string
	delivery   Delivery
	mode       Mode
	err        error
}

// Switch menambahkan flag/switch tanpa nilai (misal: stats, detail, without-paging, terse, brief).
func (b *PrintCmdBuilder) Switch(switches ...string) *PrintCmdBuilder {
	if b.switches == nil {
		b.switches = make(map[string]bool)
	}
	for _, s := range switches {
		s = strings.TrimSpace(s)
		if s != "" {
			b.switches[s] = true
		}
	}
	return b
}

// Proplist memilih kolom properti yang dikembalikan (.proplist).
func (b *PrintCmdBuilder) Proplist(props ...string) *PrintCmdBuilder {
	for _, p := range props {
		for _, part := range strings.Split(p, ",") {
			part = strings.TrimSpace(part)
			if part != "" {
				b.proplist = append(b.proplist, part)
			}
		}
	}
	return b
}

// Interval mengatur interval refresh output stream (misal 1s..5s).
// Secara otomatis menyetel Delivery menjadi DeliveryStream dan Mode menjadi ModeInterval.
func (b *PrintCmdBuilder) Interval(d time.Duration) *PrintCmdBuilder {
	b.interval = d
	if d > 0 {
		b.delivery = DeliveryStream
		b.mode = ModeInterval
	}
	return b
}

// Follow mengatur print follow untuk streaming perubahan realtime.
func (b *PrintCmdBuilder) Follow(followOnly ...bool) *PrintCmdBuilder {
	b.follow = true
	b.delivery = DeliveryStream
	b.mode = ModeFollow
	if len(followOnly) > 0 && followOnly[0] {
		b.followOnly = true
		b.Switch("follow-only")
	} else {
		b.Switch("follow")
	}
	return b
}

// From membatasi item berdasarkan nomor/urutan (misal "0,1,2").
func (b *PrintCmdBuilder) From(from any) *PrintCmdBuilder {
	b.from = formatSelector(from)
	return b
}

// File menulis output print ke berkas di router.
func (b *PrintCmdBuilder) File(filename string, appendMode ...bool) *PrintCmdBuilder {
	b.file = filename
	if len(appendMode) > 0 && appendMode[0] {
		b.Switch("append")
	}
	return b
}

// Where menambahkan filter sisi server (?prop=val, ?<prop=val, ?>prop=val).
func (b *PrintCmdBuilder) Where(prop, op, val string) *PrintCmdBuilder {
	var operator Op
	switch op {
	case "=", "==":
		operator = OpEq
	case "<":
		operator = OpLt
	case ">":
		operator = OpGt
	case "has", "?":
		operator = OpHas
	case "!has", "not-has", "-":
		operator = OpNotHas
	default:
		operator = OpEq
	}
	b.queries = append(b.queries, Query{Prop: prop, Op: operator, Value: val})
	return b
}

// WhereEq menambahkan filter equality (?prop=val).
func (b *PrintCmdBuilder) WhereEq(prop, val string) *PrintCmdBuilder {
	return b.Where(prop, "=", val)
}

// WhereHas menambahkan filter keberadaan properti (?prop).
func (b *PrintCmdBuilder) WhereHas(prop string) *PrintCmdBuilder {
	return b.Where(prop, "has", "")
}

// WhereNotHas menambahkan filter ketidakberadaan properti (?-prop).
func (b *PrintCmdBuilder) WhereNotHas(prop string) *PrintCmdBuilder {
	return b.Where(prop, "!has", "")
}

// Query menambahkan objek Query secara langsung.
func (b *PrintCmdBuilder) Query(queries ...Query) *PrintCmdBuilder {
	b.queries = append(b.queries, queries...)
	return b
}

// Arg menambahkan satu parameter bernama ke print.
func (b *PrintCmdBuilder) Arg(key string, val any) *PrintCmdBuilder {
	if b.rawArgs == nil {
		b.rawArgs = make(map[string]string)
	}
	b.rawArgs[key] = formatValue(val)
	return b
}

// Args menambahkan banyak parameter bernama sekaligus (map, struct, atau variadic key-value).
func (b *PrintCmdBuilder) Args(data ...any) *PrintCmdBuilder {
	m, err := normalizeArgs(data...)
	if err != nil {
		b.err = err
		return b
	}
	if b.rawArgs == nil {
		b.rawArgs = make(map[string]string)
	}
	for k, v := range m {
		b.rawArgs[k] = v
	}
	return b
}

// Sentence merakit kalimat kata-kata API biner MikroTik.
func (b *PrintCmdBuilder) Sentence() []string {
	words := []string{joinCommandWord(b.path, "print")}

	// 1. Interval
	if b.interval > 0 {
		words = append(words, "=interval="+FormatTime(b.interval))
	}

	// 2. Switches (diurutkan agar deterministik)
	if len(b.switches) > 0 {
		swList := make([]string, 0, len(b.switches))
		for sw := range b.switches {
			swList = append(swList, sw)
		}
		sort.Strings(swList)
		for _, sw := range swList {
			words = append(words, "="+sw+"=")
		}
	}

	// 3. Proplist (.proplist diawali titik)
	if len(b.proplist) > 0 {
		words = append(words, "=.proplist="+strings.Join(b.proplist, ","))
	}

	// 4. From
	if b.from != "" {
		words = append(words, "=from="+b.from)
	}

	// 5. File
	if b.file != "" {
		words = append(words, "=file="+b.file)
	}

	// 6. Named arguments
	if len(b.rawArgs) > 0 {
		keys := make([]string, 0, len(b.rawArgs))
		for k := range b.rawArgs {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			words = append(words, "="+k+"="+b.rawArgs[k])
		}
	}

	// 7. Queries
	for _, q := range b.queries {
		w, err := q.word()
		if err == nil {
			words = append(words, w)
		}
	}

	return words
}

// Build menyusun Request yang siap dieksekusi oleh rosrun.
func (b *PrintCmdBuilder) Build() (Request, error) {
	if b.err != nil {
		return Request{}, b.err
	}

	sent := b.Sentence()
	var fallback []string
	if len(b.proplist) > 0 {
		// Fallback kalimat tanpa .proplist jika ditolak router lama
		for _, w := range sent {
			if !strings.HasPrefix(w, "=.proplist=") {
				fallback = append(fallback, w)
			}
		}
	}

	shape := ShapeSnapshot
	if b.follow {
		shape = ShapeRows
	}

	delivery := b.delivery
	if b.switches["count-only"] {
		delivery = DeliveryOnce
	}

	return Request{
		ResourceID:     "path:" + strings.TrimPrefix(b.path, "/"),
		View:           "print",
		Delivery:       delivery,
		Mode:           b.mode,
		Sentence:       sent,
		Fallback:       fallback,
		Every:          b.interval,
		Shape:          shape,
		Fields:         b.proplist,
	}, nil
}

// Request mengembalikan Request. Panic jika terjadi error pada konfigurasi builder.
func (b *PrintCmdBuilder) Request() Request {
	req, err := b.Build()
	if err != nil {
		panic(err)
	}
	return req
}

// ========================================================================
// CmdBuilder: Builder untuk command aksi (Set, Add, Remove, Do, Tools)
// ========================================================================

type CmdBuilder struct {
	path     string
	action   string
	selector string
	switches map[string]bool
	args     map[string]string
	delivery Delivery
	mode     Mode
	interval time.Duration
	shape    Shape
	err      error
}

// Switch menambahkan flag/switch tanpa nilai (misal: "once", "stats", "compact", "show-sensitive").
func (b *CmdBuilder) Switch(switches ...string) *CmdBuilder {
	if b.switches == nil {
		b.switches = make(map[string]bool)
	}
	for _, s := range switches {
		s = strings.TrimSpace(s)
		if s != "" {
			b.switches[s] = true
		}
	}
	return b
}

// Uparam mengatur selector / positional parameter (misal "=numbers=").
func (b *CmdBuilder) Uparam(selector any) *CmdBuilder {
	b.selector = formatSelector(selector)
	return b
}

// Arg menambahkan satu parameter bernama (nama=nilai).
func (b *CmdBuilder) Arg(key string, val any) *CmdBuilder {
	if b.args == nil {
		b.args = make(map[string]string)
	}
	b.args[key] = formatValue(val)
	return b
}

// Args menambahkan banyak parameter bernama sekaligus dalam satu panggilan:
// Mendukung map[string]any, map[string]string, typed struct, atau variadic key-value.
func (b *CmdBuilder) Args(data ...any) *CmdBuilder {
	m, err := normalizeArgs(data...)
	if err != nil {
		b.err = err
		return b
	}
	if b.args == nil {
		b.args = make(map[string]string)
	}
	for k, v := range m {
		b.args[k] = v
	}
	return b
}

// Interval mengatur durasi interval streaming (untuk child command seperti monitor-traffic).
func (b *CmdBuilder) Interval(d time.Duration) *CmdBuilder {
	b.interval = d
	if d > 0 {
		b.delivery = DeliveryStream
		b.mode = ModeInterval
	}
	return b
}

// Delivery mengatur strategi delivery (Once, Stream, Events, Poll).
func (b *CmdBuilder) Delivery(d Delivery) *CmdBuilder {
	b.delivery = d
	return b
}

// Sentence merakit kalimat kata-kata API biner MikroTik.
func (b *CmdBuilder) Sentence() []string {
	cmdWord := joinCommandWord(b.path, b.action)
	words := []string{cmdWord}

	// 1. Selector / numbers
	if b.selector != "" {
		words = append(words, "=numbers="+b.selector)
	}

	// 2. Switches
	if len(b.switches) > 0 {
		swList := make([]string, 0, len(b.switches))
		for sw := range b.switches {
			swList = append(swList, sw)
		}
		sort.Strings(swList)
		for _, sw := range swList {
			words = append(words, "="+sw+"=")
		}
	}

	// 3. Interval jika disetel
	if b.interval > 0 {
		words = append(words, "=interval="+FormatTime(b.interval))
	}

	// 4. Named arguments
	if len(b.args) > 0 {
		keys := make([]string, 0, len(b.args))
		for k := range b.args {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			words = append(words, "="+k+"="+b.args[k])
		}
	}

	return words
}

// Build menyusun Request yang siap dieksekusi oleh rosrun.
func (b *CmdBuilder) Build() (Request, error) {
	if b.err != nil {
		return Request{}, b.err
	}

	shape := b.shape
	if shape == "" {
		shape = ShapeSnapshot
	}

	delivery := b.delivery
	if delivery == "" {
		delivery = DeliveryOnce
	}

	mode := b.mode
	if mode == "" {
		mode = ModeOneShot
	}

	actionName := b.action
	if actionName == "" {
		actionName = "do"
	}

	cmdWord := joinCommandWord(b.path, b.action)
	return Request{
		ResourceID: "path:" + strings.TrimPrefix(cmdWord, "/"),
		View:       actionName,
		Delivery:   delivery,
		Mode:       mode,
		Sentence:   b.Sentence(),
		Every:      b.interval,
		Shape:      shape,
	}, nil
}

// Request mengembalikan Request. Panic jika terjadi error pada konfigurasi builder.
func (b *CmdBuilder) Request() Request {
	req, err := b.Build()
	if err != nil {
		panic(err)
	}
	return req
}

// ========================================================================
// Helper format data & normalisasi
// ========================================================================

func joinCommandWord(path, action string) string {
	path = strings.TrimSuffix(path, "/")
	action = strings.TrimPrefix(action, "/")
	if path == "" {
		if action == "" {
			return "/"
		}
		return "/" + action
	}
	if action == "" {
		return path
	}
	return path + "/" + action
}

func formatSelector(selector any) string {
	if selector == nil {
		return ""
	}
	switch s := selector.(type) {
	case string:
		return s
	case int:
		return strconv.Itoa(s)
	case int64:
		return strconv.FormatInt(s, 10)
	case []string:
		return strings.Join(s, ",")
	case []int:
		parts := make([]string, len(s))
		for i, n := range s {
			parts[i] = strconv.Itoa(n)
		}
		return strings.Join(parts, ",")
	case fmt.Stringer:
		return s.String()
	default:
		return fmt.Sprint(s)
	}
}

func formatValue(val any) string {
	if val == nil {
		return ""
	}
	switch v := val.(type) {
	case string:
		return v
	case bool:
		if v {
			return "yes"
		}
		return "no"
	case time.Duration:
		return FormatTime(v)
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprint(v)
	}
}

func normalizeArgs(inputs ...any) (map[string]string, error) {
	out := make(map[string]string)
	if len(inputs) == 0 {
		return out, nil
	}
	if len(inputs) == 1 {
		if inputs[0] == nil {
			return out, nil
		}
		switch v := inputs[0].(type) {
		case map[string]string:
			for k, val := range v {
				out[k] = val
			}
			return out, nil
		case map[string]any:
			for k, val := range v {
				out[k] = formatValue(val)
			}
			return out, nil
		default:
			val := reflect.ValueOf(v)
			for val.Kind() == reflect.Pointer {
				if val.IsNil() {
					return out, nil
				}
				val = val.Elem()
			}
			if val.Kind() == reflect.Struct {
				return ArgsOf(v)
			}
			return nil, fmt.Errorf("argumen harus berupa map, struct, atau pasangan key-value genap, dapat %T", inputs[0])
		}
	}
	if len(inputs)%2 != 0 {
		return nil, errors.New("argumen pasangan key-value harus berjumlah genap")
	}
	for i := 0; i < len(inputs); i += 2 {
		key := fmt.Sprint(inputs[i])
		val := formatValue(inputs[i+1])
		out[key] = val
	}
	return out, nil
}
