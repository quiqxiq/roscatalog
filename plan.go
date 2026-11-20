package roscatalog

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Delivery menentukan siapa yang menggerakkan pengiriman data.
type Delivery string

const (
	// DeliveryOnce menjalankan command satu kali.
	DeliveryOnce Delivery = "once"
	// DeliveryStream: router yang mengirim terus-menerus (print interval=ms,
	// follow / follow-only, listen, atau command native seperti monitor).
	DeliveryStream Delivery = "stream"
	// DeliveryPoll: Go mengulang command one-shot (print, atau monitor =once=)
	// tiap Plan.Every. Router tidak menyimpan stream terbuka.
	DeliveryPoll Delivery = "poll"
)

// MinPollEvery dan MaxPollEvery membatasi polling dari sisi Go. Polling
// terlalu rapat membebani router; terlalu renggang membuat dashboard basi.
const (
	MinPollEvery = 5 * time.Second
	MaxPollEvery = 30 * time.Second
)

// Plan adalah pilihan pemakai: stream atau poll, seberapa sering, dan apakah
// perubahan perlu dihitung (diff) atau snapshot awal diambil (seed).
type Plan struct {
	Delivery Delivery
	// Every berarti:
	//   Stream pada view interval : nilai =interval= (MinPrintInterval..MaxPrintInterval)
	//   Stream pada view native   : =interval= bila command mendukungnya
	//   Poll                      : jeda antar one-shot dari sisi Go
	//                               (MinPollEvery..MaxPollEvery)
	Every time.Duration
	// Diff mengubah hasil poll menjadi event added/changed/removed (poll pada
	// view listen/follow selalu memakai diff).
	Diff bool
	// Seed pada stream listen: setelah listen aktif, ambil snapshot awal dengan
	// print agar pemakai tahu keadaan sebelum event pertama.
	Seed bool
}

// Once menjalankan satu kali.
func Once() Plan { return Plan{Delivery: DeliveryOnce} }

// Stream memakai mode asli view dengan interval bawaan.
func Stream() Plan { return Plan{Delivery: DeliveryStream} }

// StreamEvery seperti Stream dengan interval sendiri. Nilai harus berada di
// antara MinPrintInterval dan MaxPrintInterval.
func StreamEvery(d time.Duration) Plan { return Plan{Delivery: DeliveryStream, Every: d} }

// Poll mengulang one-shot dari Go tiap d. Nilai harus berada di antara
// MinPollEvery dan MaxPollEvery.
func Poll(d time.Duration) Plan { return Plan{Delivery: DeliveryPoll, Every: d} }

// WithDiff meminta event added/changed/removed.
func (p Plan) WithDiff() Plan { p.Diff = true; return p }

// WithSeed meminta snapshot awal untuk stream listen.
func (p Plan) WithSeed() Plan { p.Seed = true; return p }

// Shape menjelaskan bagaimana balasan router dikelompokkan oleh pemanggil.
type Shape string

const (
	// ShapeSnapshot: beberapa !re membentuk satu snapshot penuh (print, interval).
	ShapeSnapshot Shape = "snapshot"
	// ShapeRows: tiap !re adalah satu kejadian (native, follow).
	ShapeRows Shape = "rows"
	// ShapeChanges: tiap !re adalah perubahan satu item (listen; .dead=true = hilang).
	ShapeChanges Shape = "changes"
)

// Request adalah hasil Build: semua yang dibutuhkan runner untuk menjalankan
// satu view sesuai Plan. Struktur ini tidak bergantung pada library API.
type Request struct {
	ResourceID   string
	View         string
	Delivery     Delivery
	Mode         Mode     // mode di sisi router yang benar-benar dipakai (oneshot saat poll)
	Sentence     []string // kirim apa adanya; jangan tambah .tag bila library mengelolanya
	Fallback     []string // sama seperti Sentence tapi tanpa =.proplist= (retry bila ditolak)
	Every        time.Duration
	Fields       []string // kolom yang dipertahankan di klien (tanpa .id dan Key)
	Key          string   // properti kunci baris; kosong = tanpa kunci
	Shape        Shape
	CountOnly    bool // balasan berupa =ret=N pada !done
	Diff         bool
	Baseline     bool // poll follow-only: hasil pertama hanya jadi baseline
	Seed         bool // stream listen: ambil snapshot awal
	SeedSentence []string
}

// Fingerprint mengembalikan kunci stabil untuk cache/shared feed. Semua field
// yang memengaruhi kalimat ke router dan bentuk hasil ikut disertakan, sehingga
// dua Request dengan hasil berbeda tidak pernah berbagi stream.
func (r Request) Fingerprint() string {
	parts := []string{
		r.ResourceID, r.View, string(r.Delivery), string(r.Mode),
		r.Every.String(), string(r.Shape), r.Key,
		strconv.FormatBool(r.CountOnly), strconv.FormatBool(r.Diff),
		strconv.FormatBool(r.Baseline), strconv.FormatBool(r.Seed),
	}
	parts = append(parts, r.Sentence...)
	parts = append(parts, "\x1e")
	parts = append(parts, r.Fallback...)
	parts = append(parts, "\x1e")
	parts = append(parts, r.Fields...)
	parts = append(parts, "\x1e")
	parts = append(parts, r.SeedSentence...)
	return strings.Join(parts, "\x1f")
}

// Deliveries mengembalikan pengiriman yang valid untuk sebuah view; elemen
// pertama adalah bawaan bila Plan.Delivery kosong.
func (r *Resource) Deliveries(viewName string, ver Version) []Delivery {
	v, ok := r.View(viewName, ver)
	if !ok {
		return nil
	}
	switch v.Mode {
	case ModeOneShot:
		return []Delivery{DeliveryOnce, DeliveryPoll}
	case ModeInterval, ModeFollow:
		return []Delivery{DeliveryStream, DeliveryPoll}
	case ModeNative:
		out := []Delivery{DeliveryStream}
		if v.hasArg("once") {
			out = append(out, DeliveryOnce, DeliveryPoll)
		}
		return out
	case ModeListen:
		out := []Delivery{DeliveryStream}
		if _, ok := r.View("list", ver); ok {
			out = append(out, DeliveryPoll)
		}
		return out
	}
	return nil
}

func (v View) hasArg(name string) bool {
	for _, a := range v.Args {
		if a.Name == name {
			return true
		}
	}
	return false
}

func (v View) hasFlag(name string) bool {
	for _, f := range v.Flags {
		if f == name {
			return true
		}
	}
	return false
}

func fieldsOf(props []string) []string {
	var out []string
	for _, p := range props {
		if p != ".id" {
			out = append(out, p)
		}
	}
	return out
}

func cloneArgs(m map[string]string) map[string]string {
	out := make(map[string]string, len(m)+1)
	for k, v := range m {
		out[k] = v
	}
	return out
}

// Build menyusun Request untuk view dan Plan tertentu.
func (r *Resource) Build(viewName string, ver Version, plan Plan, p Params) (Request, error) {
	v, ok := r.View(viewName, ver)
	if !ok {
		return Request{}, fmt.Errorf("%s: view %q tidak ada", r.ID, viewName)
	}
	allowed := r.Deliveries(viewName, ver)
	d := plan.Delivery
	if d == "" {
		d = allowed[0]
	}
	valid := false
	for _, a := range allowed {
		if a == d {
			valid = true
		}
	}
	if !valid {
		names := make([]string, len(allowed))
		for i, a := range allowed {
			names[i] = string(a)
		}
		return Request{}, fmt.Errorf("%s/%s (%s): pengiriman %q tidak didukung; pilihan: %s",
			r.ID, v.Name, v.Mode, d, strings.Join(names, ", "))
	}
	if d == DeliveryStream && plan.Every > 0 {
		if plan.Every < MinPrintInterval || plan.Every > MaxPrintInterval {
			return Request{}, fmt.Errorf("%s/%s: interval %s di luar batas %s..%s",
				r.ID, v.Name, plan.Every, MinPrintInterval, MaxPrintInterval)
		}
	}
	if d == DeliveryPoll {
		if plan.Every < MinPollEvery {
			return Request{}, fmt.Errorf("%s/%s: jeda poll %s di bawah batas %s", r.ID, v.Name, plan.Every, MinPollEvery)
		}
		if plan.Every > MaxPollEvery {
			return Request{}, fmt.Errorf("%s/%s: jeda poll %s di atas batas %s", r.ID, v.Name, plan.Every, MaxPollEvery)
		}
	}

	req := Request{
		ResourceID: r.ID, View: v.Name, Delivery: d, Every: plan.Every,
		Key: v.Key, Diff: plan.Diff, CountOnly: v.hasFlag("count-only"),
	}

	// pv adalah view yang benar-benar dikirim; p2 parameter turunan.
	pv := v
	p2 := p
	p2.Args = cloneArgs(p.Args)
	if len(p.Args) == 0 {
		p2.Args = nil
	}
	fieldSrc := v.Proplist

	switch d {
	case DeliveryOnce, DeliveryPoll:
		switch v.Mode {
		case ModeOneShot:
			req.Shape = ShapeSnapshot
		case ModeInterval:
			pv.Mode, pv.Interval = ModeOneShot, 0
			req.Shape = ShapeSnapshot
		case ModeNative:
			if p2.Args == nil {
				p2.Args = map[string]string{}
			}
			p2.Args["once"] = ""
			req.Shape = ShapeRows
		case ModeFollow:
			pv.Mode, pv.Flags = ModeOneShot, nil
			req.Shape = ShapeSnapshot
			req.Diff = true
			req.Baseline = v.hasFlag("follow-only")
		case ModeListen:
			lv, _ := r.View("list", ver)
			pv, fieldSrc = lv, lv.Proplist
			req.Shape, req.Diff, req.Key = ShapeSnapshot, true, lv.Key
		}
		req.Mode = ModeOneShot
		if v.Mode == ModeNative {
			req.Mode = ModeNative
		}
	case DeliveryStream:
		req.Mode = v.Mode
		switch v.Mode {
		case ModeInterval:
			if plan.Every > 0 {
				if plan.Every < MinPrintInterval || plan.Every > MaxPrintInterval {
					return Request{}, fmt.Errorf("%s/%s: interval %s di luar batas %s..%s",
						r.ID, v.Name, plan.Every, MinPrintInterval, MaxPrintInterval)
				}
				p2.Interval = plan.Every
			}
			req.Every = p2.Interval
			if req.Every == 0 {
				req.Every = v.Interval
			}
			if req.Every < MinPrintInterval || req.Every > MaxPrintInterval {
				return Request{}, fmt.Errorf("%s/%s: interval %s di luar batas %s..%s",
					r.ID, v.Name, req.Every, MinPrintInterval, MaxPrintInterval)
			}
			req.Shape = ShapeSnapshot
		case ModeNative:
			if plan.Every > 0 && v.hasArg("interval") {
				if plan.Every < MinPrintInterval || plan.Every > MaxPrintInterval {
					return Request{}, fmt.Errorf("%s/%s: interval %s di luar batas %s..%s",
						r.ID, v.Name, plan.Every, MinPrintInterval, MaxPrintInterval)
				}
				if p2.Args == nil {
					p2.Args = map[string]string{}
				}
				if _, set := p2.Args["interval"]; !set {
					p2.Args["interval"] = FormatTime(plan.Every)
				}
			}
			req.Shape = ShapeRows
		case ModeFollow:
			req.Shape = ShapeRows
		case ModeListen:
			req.Shape = ShapeChanges
			if lv, ok := r.View("list", ver); ok {
				// Jangan kirim =.proplist= pada listen. Di RouterOS 7.24.6
				// proplist menekan event .dead=true sehingga penghapusan item
				// tidak terdeteksi. Kolom tetap disaring di klien lewat
				// Request.Fields yang diambil dari view list.
				fieldSrc = lv.Proplist
				if plan.Seed {
					seed, err := lv.Sentence(p)
					if err != nil {
						return Request{}, err
					}
					req.Seed, req.SeedSentence = true, seed
				}
			} else {
				for _, o := range r.Views {
					if o.Name == "list" {
						return Request{}, fmt.Errorf("%s: view list bergantung versi; sebutkan v6 atau v7", r.ID)
					}
				}
				if plan.Seed {
					return Request{}, fmt.Errorf("%s: Seed butuh view list", r.ID)
				}
			}
		}
	}

	sent, err := pv.Sentence(p2)
	if err != nil {
		return Request{}, err
	}
	req.Sentence = sent
	if pv.ServerProplist && len(pv.Proplist) > 0 {
		np := pv
		np.ServerProplist = false
		if fb, err := np.Sentence(p2); err == nil {
			req.Fallback = fb
		}
	}
	req.Fields = fieldsOf(fieldSrc)
	return req, nil
}
