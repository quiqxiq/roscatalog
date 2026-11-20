package roscatalog

import (
	"reflect"
	"strings"
	"testing"
)

type customRow struct {
	ID    string `ros:".id,key"`
	Addr  string `ros:"address"`
	Extra string `ros:"field-dari-router-lain"`
	Skip  string `ros:"-"`
}

type customAdd struct {
	Address   string `ros:"address,required"`
	Interface string `ros:"interface,required"`
	Comment   string `ros:"comment,omitempty"`
}

func TestFieldNamesFromTags(t *testing.T) {
	got := FieldNames[customRow]()
	want := []string{".id", "address", "field-dari-router-lain"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestFieldNamesFallbackKebab(t *testing.T) {
	type row struct {
		RateLimit string
	}
	if got := FieldNames[row](); !reflect.DeepEqual(got, []string{"rate-limit"}) {
		t.Fatalf("fallback kebab salah: %v", got)
	}
}

func TestDecodeAndDecodeAll(t *testing.T) {
	row := map[string]string{
		".id":                    "*1",
		"address":                "10.0.0.1/32",
		"field-dari-router-lain": "x",
		"comment":                "diabaikan",
	}
	got, err := Decode[customRow](row)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "*1" || got.Addr != "10.0.0.1/32" || got.Extra != "x" {
		t.Fatalf("decode salah: %+v", got)
	}

	all, err := DecodeAll[customRow]([]map[string]string{row, row})
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 2 || all[1].Extra != "x" {
		t.Fatalf("decodeAll salah: %+v", all)
	}
}

func TestArgsOfOmitEmptyAndRequired(t *testing.T) {
	got, err := ArgsOf(customAdd{Address: "10.0.0.2/32", Interface: "ether1"})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, map[string]string{
		"address": "10.0.0.2/32", "interface": "ether1",
	}) {
		t.Fatalf("argsOf salah: %v", got)
	}
	if _, err := ArgsOf(customAdd{Address: "10.0.0.2/32"}); err == nil {
		t.Fatal("required kosong harus error")
	}
}

func TestProplistOverrideWinsOverOnly(t *testing.T) {
	res, ok := Lookup("ip-address")
	if !ok {
		t.Fatal("resource tidak ada")
	}
	req, err := res.Build("list", Once(), Params{
		Only:     []string{"comment"},
		Proplist: FieldNames[customRow](),
	})
	if err != nil {
		t.Fatal(err)
	}
	if !containsWord(req.Sentence, "=.proplist=.id,address,field-dari-router-lain") {
		t.Fatalf("Proplist tidak dipakai: %v", req.Sentence)
	}
	want := []string{"address", "field-dari-router-lain"}
	if !reflect.DeepEqual(req.Fields, want) {
		t.Fatalf("Fields salah: %v", req.Fields)
	}
}

func TestRawWinsOverProplistAndOnly(t *testing.T) {
	res, ok := Lookup("ip-address")
	if !ok {
		t.Fatal("resource tidak ada")
	}
	req, err := res.Build("list", Once(), Params{
		Raw:      true,
		Only:     []string{"address"},
		Proplist: FieldNames[customRow](),
	})
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
}
