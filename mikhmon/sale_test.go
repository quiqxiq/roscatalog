package mikhmon

import "testing"

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
