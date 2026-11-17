package mikhmon

import (
	"strings"
	"testing"

	ros "example.com/roscatalog"
)

func TestCatalogSubset(t *testing.T) {
	ids := IDs()
	if len(ids) == 0 {
		t.Fatal("subset kosong")
	}
	if got := len(Resources()); got != len(ids) {
		t.Fatalf("Resources()=%d, IDs()=%d (ada resource root yang hilang)", got, len(ids))
	}
	for _, id := range ids {
		if _, ok := ros.Lookup(id); !ok {
			t.Errorf("resource root %q tidak ada", id)
		}
		if !IsResource(id) {
			t.Errorf("IsResource(%q)=false", id)
		}
	}
	if _, ok := Lookup("hotspot-user"); !ok {
		t.Error("hotspot-user harus ada")
	}
	if _, ok := Lookup("user"); ok {
		t.Error("user router tidak boleh masuk subset Mikhmon")
	}
	if _, ok := Lookup("mikhmon-sales"); !ok {
		t.Error("mikhmon-sales harus ada")
	}
}

func TestIDsSalinan(t *testing.T) {
	ids := IDs()
	ids[0] = "diubah"
	if IDs()[0] == "diubah" {
		t.Fatal("IDs() mengembalikan slice internal")
	}
}

func TestOneLine(t *testing.T) {
	in := "a\n\tb   c"
	if got := OneLine(in); got != "a b c" {
		t.Fatalf("OneLine=%q", got)
	}
}

func TestOnLoginScriptRem(t *testing.T) {
	got, err := OnLoginScript(OnLoginOptions{
		Mode: ExpiryRemove, Price: "5000", Validity: "1d",
		SellingPrice: "7000", LockUser: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.ContainsAny(got, "\r\n") {
		t.Fatal("on-login mengandung newline")
	}
	if !strings.HasPrefix(got, `:put (",rem,5000,1d,7000,,Enable,Disable,")`) {
		t.Fatalf("metadata salah: %.80s", got)
	}
	if !strings.Contains(got, `:local mode "X"`) {
		t.Error("mode X tidak ada")
	}
	if !strings.Contains(got, `interval="1d"`) {
		t.Error("validity tidak masuk ke scheduler")
	}
	if !strings.Contains(got, "mac-address=$mac") {
		t.Error("lock MAC tidak ada")
	}
	if strings.Contains(got, "comment=mikhmon") {
		t.Error("mode rem tidak boleh menulis laporan")
	}
	if !strings.HasSuffix(got, "}}") {
		t.Errorf("script tidak ditutup: %s", got[len(got)-20:])
	}
}

func TestOnLoginScriptNoticeRecord(t *testing.T) {
	got, err := OnLoginScript(OnLoginOptions{
		Mode: ExpiryNoticeRecord, Price: "5000", Validity: "1d",
		SellingPrice: "7000", Profile: "1jam", LockServer: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, `:local mode "N"`) {
		t.Error("mode N tidak ada")
	}
	if !strings.Contains(got, "comment=mikhmon") {
		t.Error("laporan penjualan tidak ada")
	}
	if !strings.Contains(got, "-|-1jam-|-") {
		t.Error("nama profile tidak masuk laporan")
	}
	if !strings.Contains(got, "server=$srv") {
		t.Error("lock server tidak ada")
	}
	if strings.ContainsAny(got, "\r\n") {
		t.Fatal("on-login mengandung newline")
	}
}

func TestOnLoginScriptNone(t *testing.T) {
	got, err := OnLoginScript(OnLoginOptions{
		Mode: ExpiryNone, Price: "5000", Validity: "1d", SellingPrice: "7000",
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(got, ":local mode") || strings.Contains(got, "scheduler") {
		t.Fatalf("mode 0 tidak boleh punya logika expiry: %s", got)
	}
	if !strings.HasPrefix(got, `:put (",0,5000,1d,7000,,Disable,Disable,")`) {
		t.Fatalf("metadata salah: %s", got)
	}
}

func TestOnLoginScriptValidation(t *testing.T) {
	if _, err := OnLoginScript(OnLoginOptions{Mode: "bogus"}); err == nil {
		t.Error("mode tidak dikenal harus error")
	}
	if _, err := OnLoginScript(OnLoginOptions{Mode: ExpiryRemove}); err == nil {
		t.Error("validity kosong harus error")
	}
	if _, err := OnLoginScript(OnLoginOptions{
		Mode: ExpiryRemoveRecord, Validity: "1d",
	}); err == nil {
		t.Error("profile kosong pada remc harus error")
	}
	if _, err := OnLoginScript(OnLoginOptions{
		Mode: ExpiryRemove, Validity: "1d", Profile: "bad\"profile",
	}); err == nil {
		t.Error("profile dengan kutip harus error")
	}
}

func TestExpireMonitorScript(t *testing.T) {
	got := ExpireMonitorScript()
	if got == "" {
		t.Fatal("script kosong")
	}
	if strings.ContainsAny(got, "\r\n") {
		t.Fatal("expire monitor harus satu baris")
	}
	for _, want := range []string{
		"normdate", "dateint", "timeint",
		"limit-uptime=1s", "/ip hotspot user remove",
		"/ip hotspot active remove",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("script tidak memuat %q", want)
		}
	}
	if len(got) < 1500 {
		t.Errorf("script terlalu pendek: %d", len(got))
	}
}
