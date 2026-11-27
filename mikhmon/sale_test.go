package mikhmon

import (
	"testing"

	"github.com/quiqxiq/roscatalog/fields"
)

func TestSaleRoundTrip(t *testing.T) {
	s := Sale{
		Date: "nov/18/2026", Time: "10:00:00", User: "vc001",
		Price: "5000", Address: "198.51.100.10", MAC: "AA:BB:CC:DD:EE:FF",
		Validity: "1d", Profile: "1jam", Comment: "vc-123-11.18.26",
	}
	if err := s.Validate(); err != nil {
		t.Fatal(err)
	}
	if got := s.Owner(); got != "nov2026" {
		t.Fatalf("Owner=%q", got)
	}
	got, err := ParseSaleName(s.Name())
	if err != nil {
		t.Fatal(err)
	}
	if got != s {
		t.Fatalf("roundtrip:\n got %+v\nwant %+v", got, s)
	}
}

func TestSaleOwnerV7(t *testing.T) {
	s := Sale{Date: "2026-11-18"}
	if got := s.Owner(); got != "nov2026" {
		t.Fatalf("Owner=%q", got)
	}
	if got := s.MonthText(); got != "nov2026" {
		t.Fatalf("MonthText=%q", got)
	}
	if got := s.MonthNumeric(); got != "112026" {
		t.Fatalf("MonthNumeric=%q", got)
	}
}

func TestSaleMatchesMonthAndDate(t *testing.T) {
	s7 := Sale{Date: "2024-03-18"}
	if !s7.MatchesMonth("mar2024") {
		t.Error("2024-03-18 harus cocok dengan mar2024")
	}
	if !s7.MatchesMonth("032024") {
		t.Error("2024-03-18 harus cocok dengan 032024")
	}
	if !s7.MatchesMonth("2024-03") {
		t.Error("2024-03-18 harus cocok dengan 2024-03")
	}
	if s7.MatchesMonth("apr2024") {
		t.Error("2024-03-18 tidak boleh cocok dengan apr2024")
	}
	if !s7.MatchesDate("2024-03-18") {
		t.Error("2024-03-18 harus cocok dengan 2024-03-18")
	}
	if !s7.MatchesDate("mar/18/2024") {
		t.Error("2024-03-18 harus cocok dengan mar/18/2024")
	}
	if !s7.MatchesDate("03/18/2024") {
		t.Error("2024-03-18 harus cocok dengan 03/18/2024")
	}
	if s7.MatchesDate("2024-03-19") {
		t.Error("2024-03-18 tidak boleh cocok dengan 2024-03-19")
	}

	s6 := Sale{Date: "mar/18/2024"}
	if !s6.MatchesMonth("mar2024") || !s6.MatchesMonth("032024") {
		t.Error("mar/18/2024 harus cocok dengan mar2024 dan 032024")
	}
	if !s6.MatchesDate("2024-03-18") {
		t.Error("mar/18/2024 harus cocok dengan 2024-03-18")
	}
}

func TestMonthBuckets(t *testing.T) {
	cases := []struct {
		input       string
		wantText    string
		wantNumeric string
		wantOK      bool
	}{
		{"2024-03-18", "mar2024", "032024", true},
		{"2024-03", "mar2024", "032024", true},
		{"mar2024", "mar2024", "032024", true},
		{"032024", "mar2024", "032024", true},
		{"mar/18/2024", "mar2024", "032024", true},
		{"03/18/2024", "mar2024", "032024", true},
		{"mar/2024", "mar2024", "032024", true},
		{"03/2024", "mar2024", "032024", true},
		{"invalid", "", "", false},
		{"", "", "", false},
	}
	for _, tc := range cases {
		text, num, ok := MonthBuckets(tc.input)
		if ok != tc.wantOK || text != tc.wantText || num != tc.wantNumeric {
			t.Errorf("MonthBuckets(%q) = (%q, %q, %v), want (%q, %q, %v)",
				tc.input, text, num, ok, tc.wantText, tc.wantNumeric, tc.wantOK)
		}
	}
}

func TestSaleOwnerInvalid(t *testing.T) {
	for _, date := range []string{"", "18/11/2026", "2026-13-01", "foo"} {
		if got := ownerFromDate(date); got != "" {
			t.Errorf("ownerFromDate(%q)=%q, want empty", date, got)
		}
	}
}

func TestParseSaleNameInvalid(t *testing.T) {
	if _, err := ParseSaleName("a-|-b-|-c"); err == nil {
		t.Fatal("nama kurang field harus error")
	}
}

func TestSaleValidate(t *testing.T) {
	s := Sale{Date: "nov/18/2026"}
	if err := s.Validate(); err == nil {
		t.Fatal("field kosong harus error")
	}
	s = Sale{
		Date: "nov/18/2026", Time: "10:00:00", User: "u",
		Price: "1", Validity: "1d", Profile: "p",
	}
	if err := s.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestSaleToFieldsAndFromRow(t *testing.T) {
	s := Sale{
		Date: "nov/18/2026", Time: "10:00:00", User: "vc001",
		Price: "5000", Address: "198.51.100.10", MAC: "AA:BB:CC:DD:EE:FF",
		Validity: "1d", Profile: "1jam", Comment: "vc-123",
	}
	f := s.ToFields()
	if f.Name != s.Name() || f.Source != s.Source() || f.Owner != "nov2026" || f.Comment != SalesComment {
		t.Fatalf("ToFields salah: %+v", f)
	}

	row := fields.MikhmonSalesRow{
		ID:      "*S1",
		Name:    s.Name(),
		Owner:   f.Owner,
		Source:  f.Source,
		Comment: f.Comment,
	}
	parsed, err := SaleFromRow(row)
	if err != nil {
		t.Fatalf("SaleFromRow failed: %v", err)
	}
	if parsed != s {
		t.Fatalf("SaleFromRow:\n got %+v\nwant %+v", parsed, s)
	}
}
