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
		Mode: ExpiryRemove, Validity: "1d", LockUser: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.ContainsAny(got, "\r\n") {
		t.Fatal("on-login mengandung newline")
	}
	if strings.Contains(got, ":put") {
		t.Fatal("metadata harga tidak boleh ditulis ke on-login")
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
		t.Error("mode rem tanpa record_sales tidak boleh menulis laporan")
	}
	if !strings.HasSuffix(got, "}}") {
		t.Errorf("script tidak ditutup: %s", got[len(got)-20:])
	}
}

func TestOnLoginScriptHookAndFallback(t *testing.T) {
	got, err := OnLoginScript(OnLoginOptions{
		Mode: ExpiryNoticeRecord, Validity: "1d", Profile: "1jam",
		RecordSales: true, LockServer: true,
		PublicURL: "https://mikmon.example", HookToken: "tok_123", VerifyTLS: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(got, ":put") || strings.Contains(got, "5000") || strings.Contains(got, "7000") {
		t.Fatal("harga tidak boleh masuk on-login")
	}
	if !strings.Contains(got, `:local mode "N"`) {
		t.Error("mode N tidak ada")
	}
	if !strings.Contains(got, "/tool fetch") || !strings.Contains(got, "on-error=") {
		t.Error("fetch-first fallback tidak ada")
	}
	if !strings.Contains(got, "comment=mikhmon") {
		t.Error("fallback script tidak ada")
	}
	if !strings.Contains(got, "-|-1jam-|-") || !strings.Contains(got, "v2-|-") {
		t.Error("payload v2/profile tidak ada")
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
		Mode: ExpiryNone, Validity: "1d", Profile: "1jam", RecordSales: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(got, "interval=\"1d\"") {
		t.Fatalf("mode 0 tidak boleh punya logika expiry: %s", got)
	}
	if strings.Contains(got, ":put") {
		t.Fatal("metadata harga tidak boleh ada")
	}
	if !strings.Contains(got, "comment=mikhmon") {
		t.Error("record sales fallback tidak ada")
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
		Mode: ExpiryRemoveRecord, Validity: "1d", RecordSales: true,
	}); err == nil {
		t.Error("profile kosong pada record sales harus error")
	}
	if _, err := OnLoginScript(OnLoginOptions{
		Mode: ExpiryRemove, Validity: "1d", Profile: "bad\"profile",
	}); err == nil {
		t.Error("profile dengan kutip harus error")
	}
}

func TestParseOnLoginScript(t *testing.T) {
	old := `:put (",rem,5000,1d,7000,,Enable,Disable,")`
	meta := ParseOnLoginScript(old)
	if meta.Mode != "rem" || meta.Price != "5000" || meta.Validity != "1d" || !meta.LockUser {
		t.Fatalf("legacy metadata tidak terbaca: %+v", meta)
	}
	fresh := `:local mode "N"; interval="30d"; mac-address=$mac; /tool fetch`
	meta = ParseOnLoginScript(fresh)
	if meta.Mode != "N" || meta.Validity != "30d" || !meta.LockUser || !meta.RecordSales {
		t.Fatalf("metadata baru tidak terbaca: %+v", meta)
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
