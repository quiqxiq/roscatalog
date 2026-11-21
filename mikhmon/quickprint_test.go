package mikhmon

import (
	"testing"

	"example.com/roscatalog/fields"
)

func TestQuickPrintRoundTrip(t *testing.T) {
	p := QuickPrintPreset{
		Package: "1jam", Server: "all", UserMode: "vc", UserLength: "4",
		Prefix: "", Char: "mix", Profile: "1jam", TimeLimit: "1d",
		DataLimit: "0", Comment: "catatan", Validity: "1d",
		Price: "5000", SellingPrice: "7000", UserLock: "Disable",
	}
	got, err := ParseQuickPrintSource(p.Source())
	if err != nil {
		t.Fatal(err)
	}
	if got != p {
		t.Fatalf("roundtrip:\n got %+v\nwant %+v", got, p)
	}
}

func TestQuickPrintName(t *testing.T) {
	p := QuickPrintPreset{Package: "1 jam  hemat"}
	if got, want := p.Name(), "Quick_Print_1-jam-hemat"; got != want {
		t.Fatalf("Name=%q, want %q", got, want)
	}
	if got := (QuickPrintPreset{}).Name(); got != "Quick_Print_Preset" {
		t.Fatalf("default Name=%q", got)
	}
}

func TestQuickPrintSellingFallback(t *testing.T) {
	p := QuickPrintPreset{Package: "x", Price: "5000"}
	got, err := ParseQuickPrintSource(p.Source())
	if err != nil {
		t.Fatal(err)
	}
	if got.Price != "5000" || got.SellingPrice != "5000" {
		t.Fatalf("fallback harga salah: %+v", got)
	}
}

func TestQuickPrintParseInvalid(t *testing.T) {
	for _, source := range []string{"", "1#2", "#a#b"} {
		if _, err := ParseQuickPrintSource(source); err == nil {
			t.Errorf("source %q harus error", source)
		}
	}
}

func TestQuickPrintToFieldsAndFromRow(t *testing.T) {
	p := QuickPrintPreset{
		Package: "1jam", Server: "all", UserMode: "vc", UserLength: "4",
		Prefix: "", Char: "mix", Profile: "1jam", TimeLimit: "1d",
		DataLimit: "0", Comment: "catatan", Validity: "1d",
		Price: "5000", SellingPrice: "7000", UserLock: "Disable",
	}
	f := p.ToFields()
	if f.Name != p.Name() || f.Source != p.Source() || f.Comment != QuickPrintComment {
		t.Fatalf("ToFields salah: %+v", f)
	}

	row := fields.MikhmonQuickPrintRow{
		ID:      "*QP1",
		Name:    p.Name(),
		Source:  p.Source(),
		Comment: QuickPrintComment,
	}
	parsed, err := PresetFromRow(row)
	if err != nil {
		t.Fatalf("PresetFromRow failed: %v", err)
	}
	if parsed != p {
		t.Fatalf("PresetFromRow:\n got %+v\nwant %+v", parsed, p)
	}
}
