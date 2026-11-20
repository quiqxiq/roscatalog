package roscatalog

import (
	"errors"
	"testing"
	"time"
)

func TestCatalogFind(t *testing.T) {
	c := Builtin()
	if c.Len() == 0 {
		t.Fatal("builtin katalog kosong")
	}
	r, err := c.Find("ip-route")
	if err != nil || r == nil || r.ID != "ip-route" {
		t.Fatalf("Find(ip-route) = %v, %v", r, err)
	}
	if _, err := c.Find("tidak-ada"); !errors.Is(err, ErrResourceNotFound) {
		t.Errorf("Find(tidak-ada) error = %v, mau %v", err, ErrResourceNotFound)
	}
	if r, ok := c.Lookup("ip-route"); !ok || r == nil {
		t.Error("Lookup(ip-route) tidak ditemukan")
	}
	if _, ok := c.Lookup("tidak-ada"); ok {
		t.Error("Lookup(tidak-ada) seharusnya tidak ditemukan")
	}
}

func TestCatalogAddFindValidate(t *testing.T) {
	c := NewCatalog(Builtin().All()...)
	res := Table("queue-custom", GroupQueue, "Queue custom", "/queue/custom").
		Describe("Resource definisi pemakai.").
		Filter("name").
		List("name", "target", "rate", "comment").
		Count().
		Live(time.Second, true, "name", "rate").
		Add(P("name"), P("target", "rate", "comment")).
		Set("target").
		Remove("hapus")
	if err := c.Add(res); err != nil {
		t.Fatalf("Add: %v", err)
	}
	got, err := c.Find("queue-custom")
	if err != nil {
		t.Fatalf("Find: %v", err)
	}
	if got != res {
		t.Error("Find tidak mengembalikan instance yang sama")
	}
	if _, ok := c.Lookup("queue-custom"); !ok {
		t.Error("Lookup(queue-custom) tidak ditemukan")
	}
	found := false
	for _, r := range c.ByGroup(GroupQueue) {
		if r.ID == "queue-custom" {
			found = true
		}
	}
	if !found {
		t.Error("ByGroup(queue) tidak mengandung queue-custom")
	}
	if err := c.Validate(); err != nil {
		t.Errorf("Validate builtin+custom: %v", err)
	}
}

func TestCatalogAddDuplicateAndInvalid(t *testing.T) {
	c := NewCatalog(Builtin().All()...)
	if err := c.Add(Table("ip-route", GroupIP, "Duplikat", "/ip/route").List("name")); !errors.Is(err, ErrDuplicateResource) {
		t.Errorf("Add(dup) = %v, mau %v", err, ErrDuplicateResource)
	}
	if err := c.Add(Table("bad-path", GroupQueue, "Path salah", "queue/custom").List("name")); !errors.Is(err, ErrInvalidResource) {
		t.Errorf("Add(bad path) = %v, mau %v", err, ErrInvalidResource)
	}
	empty := NewCatalog()
	if err := empty.Validate(); err != nil {
		t.Errorf("Validate(katalog kosong) = %v, mau nil", err)
	}
}

func TestCatalogReplaceAndExtend(t *testing.T) {
	a := Table("res-a", GroupIP, "A", "/ip/a").List("name")
	b := Table("res-b", GroupIP, "B", "/ip/b").List("name")
	c := NewCatalog(a)
	repl := Table("res-a", GroupIP, "A2", "/ip/a").List("name")
	if err := c.Replace(repl); err != nil {
		t.Fatalf("Replace: %v", err)
	}
	if r, _ := c.Find("res-a"); r.Title != "A2" {
		t.Errorf("Replace tidak berubah: %+v", r)
	}
	if err := c.Replace(b); !errors.Is(err, ErrResourceNotFound) {
		t.Errorf("Replace(res-b) = %v, mau %v", err, ErrResourceNotFound)
	}
	if err := c.Extend(NewCatalog(b)); err != nil {
		t.Fatalf("Extend: %v", err)
	}
	if _, ok := c.Lookup("res-b"); !ok {
		t.Error("Extend tidak menambah res-b")
	}
	if err := c.Extend(NewCatalog(b)); !errors.Is(err, ErrDuplicateResource) {
		t.Errorf("Extend(dup) = %v, mau %v", err, ErrDuplicateResource)
	}
}

func TestBuildSentinelErrors(t *testing.T) {
	r, ok := Lookup("ip-route")
	if !ok {
		t.Fatal("ip-route tidak ada")
	}
	if _, err := r.Build("tidak-ada", Once(), Params{}); !errors.Is(err, ErrViewNotFound) {
		t.Errorf("Build(tidak-ada) = %v, mau %v", err, ErrViewNotFound)
	}
	if _, err := r.Build("list", Once(), Params{Only: []string{"bogus"}}); !errors.Is(err, ErrUnknownColumn) {
		t.Errorf("Build(bogus column) = %v, mau %v", err, ErrUnknownColumn)
	}
	add, ok := r.Action("add")
	if !ok {
		t.Fatal("ip-route/add tidak ada")
	}
	if _, err := add.Sentence("", nil, ""); !errors.Is(err, ErrMissingArgument) {
		t.Errorf("add tanpa arg = %v, mau %v", err, ErrMissingArgument)
	}
	set, ok := r.Action("set")
	if !ok {
		t.Fatal("ip-route/set tidak ada")
	}
	if _, err := set.Sentence("", map[string]string{"dst-address": "10.0.0.0/24", "gateway": "10.0.0.1"}, ""); !errors.Is(err, ErrSelectorRequired) {
		t.Errorf("set tanpa selector = %v, mau %v", err, ErrSelectorRequired)
	}
	if _, err := r.Build("list", Once(), Params{Flags: []string{"tidak-ada"}}); !errors.Is(err, ErrUnknownFlag) {
		t.Errorf("Build(flag tidak ada) = %v, mau %v", err, ErrUnknownFlag)
	}
	// Error router: Find tidak pernah menemukan.
	if _, ok := Lookup("tidak-ada"); ok {
		t.Error("Lookup(tidak-ada) seharusnya tidak ditemukan")
	}
}

func TestCustomResourceCRUDSentences(t *testing.T) {
	res := Table("queue-custom", GroupQueue, "Queue custom", "/queue/custom").
		Filter("name").
		List("name", "target", "rate", "comment").
		Add(P("name"), P("target", "rate", "comment")).
		Set("target").
		Remove("hapus")

	req, err := res.Build("list", Once(), Params{})
	if err != nil {
		t.Fatalf("Build list: %v", err)
	}
	if req.Sentence[0] != "/queue/custom/print" {
		t.Errorf("list sentence = %q", req.Sentence[0])
	}
	if add, ok := res.Action("add"); !ok {
		t.Error("Action(add) tidak ada")
	} else if sent, err := add.Sentence("", map[string]string{"name": "q1", "target": "1.2.3.4", "rate": "10M"}, ""); err != nil || sent[0] != "/queue/custom/add" {
		t.Errorf("add sentence = %q err=%v", sent, err)
	}
	if set, ok := res.Action("set"); !ok {
		t.Error("Action(set) tidak ada")
	} else if sent, err := set.Sentence("q1", map[string]string{"target": "9.9.9.9"}, ""); err != nil || sent[0] != "/queue/custom/set" {
		t.Errorf("set sentence = %q err=%v", sent, err)
	}
	if rm, ok := res.Action("remove"); !ok {
		t.Error("Action(remove) tidak ada")
	} else if sent, err := rm.Sentence("*1", nil, ""); err != nil || sent[0] != "/queue/custom/remove" {
		t.Errorf("remove sentence = %q err=%v", sent, err)
	}
}
