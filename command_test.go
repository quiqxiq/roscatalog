package roscatalog

import (
	"reflect"
	"strings"
	"testing"
)

func TestActionSentenceWithFlags(t *testing.T) {
	a := Action{
		Name: "print", Command: "/ip/address/print", Read: true,
		Args: []Arg{
			{Name: "brief", Kind: ArgFlag},
			{Name: "detail", Kind: ArgFlag},
			{Name: "stats", Kind: ArgFlag},
			{Name: "proplist"},
		},
	}
	got, err := a.SentenceWith("", []string{"brief"}, map[string]string{"proplist": "address"}, "")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"/ip/address/print", "=brief=", "=proplist=address"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}

	if _, err := a.SentenceWith("", []string{"bogus"}, nil, ""); err == nil {
		t.Fatal("flag asing harus ditolak")
	}
	if _, err := a.SentenceWith("", []string{"brief", "detail"}, nil, ""); err == nil {
		t.Fatal("brief + detail harus ditolak")
	}
}

func TestViewSentenceWithParamsFlags(t *testing.T) {
	v := View{
		Name: "list", Mode: ModeOneShot, Command: "/interface/bridge/print",
		AllowedFlags: []string{"brief", "detail", "stats", "proplist"},
	}
	got, err := v.Sentence(Params{Flags: []string{"detail"}})
	if err != nil {
		t.Fatal(err)
	}
	if !containsWord(got, "=detail=") {
		t.Fatalf("flag detail tidak masuk: %v", got)
	}

	if _, err := v.Sentence(Params{Flags: []string{"terse"}}); err == nil {
		t.Fatal("flag yang tidak diizinkan harus ditolak")
	}

	v.Flags = []string{"stats"}
	if _, err := v.Sentence(Params{Flags: []string{"count-only"}}); err == nil {
		t.Fatal("stats + count-only harus ditolak")
	}
}

func TestCatalogEnrichment(t *testing.T) {
	r, ok := Lookup("ip-address")
	if !ok {
		t.Fatal("resource ip-address tidak ada")
	}
	a, ok := r.Action("print")
	if !ok || !a.Read {
		t.Fatal("action print harus ada dan bertanda Read")
	}
	if !hasFlag(a.FlagNames(), "stats") && !hasFlag(a.FlagNames(), "detail") {
		t.Fatalf("flag print tidak terisi: %v", a.FlagNames())
	}
	v, ok := r.View("list", AnyVersion)
	if !ok {
		t.Fatal("view list tidak ada")
	}
	if len(v.AllowedFlags) == 0 {
		t.Fatal("AllowedFlags tidak terisi oleh enrich")
	}
	if _, ok := r.Action("comment"); !ok {
		t.Fatal("aksi comment belum ditambahkan")
	}
}

func TestBuildRaw(t *testing.T) {
	r, ok := Lookup("ip-address")
	if !ok {
		t.Fatal("resource ip-address tidak ada")
	}
	req, err := r.Build("list", AnyVersion, Once(), Params{Raw: true})
	if err != nil {
		t.Fatal(err)
	}
	for _, w := range req.Sentence {
		if strings.HasPrefix(w, "=.proplist=") {
			t.Fatalf("Raw tidak boleh mengirim proplist: %v", req.Sentence)
		}
	}
	if len(req.Fields) != 0 {
		t.Fatalf("Raw harus mengembalikan semua kolom: %v", req.Fields)
	}
	if len(req.Fallback) != 0 {
		t.Fatalf("Raw tidak butuh fallback: %v", req.Fallback)
	}
}

func TestBuildOnly(t *testing.T) {
	r, ok := Lookup("ip-address")
	if !ok {
		t.Fatal("resource ip-address tidak ada")
	}
	req, err := r.Build("list", AnyVersion, Once(), Params{Only: []string{"address", "comment"}})
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, w := range req.Sentence {
		if w == "=.proplist=.id,address,comment" {
			found = true
		}
	}
	if !found {
		t.Fatalf("proplist tidak terbatasi: %v", req.Sentence)
	}
	if !reflect.DeepEqual(req.Fields, []string{"address", "comment"}) {
		t.Fatalf("Fields salah: %v", req.Fields)
	}
	if _, err := r.Build("list", AnyVersion, Once(), Params{Only: []string{"bogus"}}); err == nil {
		t.Fatal("kolom asing harus ditolak")
	}
}

func TestPrintPresets(t *testing.T) {
	r := Table("x-test", GroupTools, "x", "/x").Stats("rx", "tx")
	v, ok := r.View("stats", AnyVersion)
	if !ok {
		t.Fatal("view stats tidak ada")
	}
	if v.Mode != ModeOneShot || !v.hasFlag("stats") {
		t.Fatalf("preset stats salah: %+v", v)
	}
	req, err := r.Build("stats", AnyVersion, Once(), Params{})
	if err != nil {
		t.Fatal(err)
	}
	if !containsWord(req.Sentence, "=stats=") {
		t.Fatalf("flag stats tidak masuk: %v", req.Sentence)
	}
}

func hasFlag(flags []string, name string) bool {
	for _, f := range flags {
		if f == name {
			return true
		}
	}
	return false
}

func containsWord(words []string, want string) bool {
	for _, w := range words {
		if w == want || strings.HasPrefix(w, want+"=") {
			return true
		}
	}
	return false
}
