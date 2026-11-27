package mikhmon

import (
	"strings"
	"testing"

	ros "github.com/quiqxiq/roscatalog"
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
	if !strings.HasPrefix(got, `:put (",rem,0,1d,0,,Enable,")`) {
		t.Errorf("header put salah: %s", got[:40])
	}
	if !strings.Contains(got, `interval="1d"`) {
		t.Error("validity tidak masuk ke scheduler")
	}
	if !strings.Contains(got, "mac-address=$mac") {
		t.Error("lock MAC tidak ada")
	}
	if !strings.Contains(got, `comment="$expLegacy"`) {
		t.Errorf("comment user harus $expLegacy murni: %s", got)
	}
	if strings.Contains(got, "$expLegacy $mode") {
		t.Error("comment user tidak boleh memiliki suffix $mode")
	}
	if strings.Contains(got, "comment=mikhmon") {
		t.Error("mode rem tanpa record_sales tidak boleh menulis laporan")
	}
	if !strings.HasSuffix(got, "}}") {
		t.Errorf("script tidak ditutup: %s", got[len(got)-20:])
	}
}

func TestOnLoginScriptRemcWithPrice(t *testing.T) {
	got, err := OnLoginScript(OnLoginOptions{
		Mode: ExpiryRemoveRecord, Price: "5000", Validity: "1d",
		SellingPrice: "7000", Profile: "1jam", LockUser: false, LockServer: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.ContainsAny(got, "\r\n") {
		t.Fatal("on-login mengandung newline")
	}
	if !strings.HasPrefix(got, `:put (",remc,5000,1d,7000,,Disable,")`) {
		t.Errorf("header put salah: %s", got[:40])
	}
	if !strings.Contains(got, `sname "$legacy-|-$time-|-$user-|-5000-|-$address-|-$mac-|-1d-|-1jam-|-$comment"`) {
		t.Errorf("script laporan salah: %s", got)
	}
	if !strings.Contains(got, "comment=mikhmon") {
		t.Error("script laporan harus memiliki comment=mikhmon")
	}
	if !strings.Contains(got, "server=$srv") {
		t.Error("lock server tidak ada")
	}
	if strings.Contains(got, "/tool fetch") {
		t.Error("tidak boleh ada webhook /tool fetch")
	}
	if !strings.HasSuffix(got, "}}") {
		t.Errorf("script tidak ditutup: %s", got[len(got)-20:])
	}
}

func TestOnLoginScriptNtfcWithPrice(t *testing.T) {
	got, err := OnLoginScript(OnLoginOptions{
		Mode: ExpiryNoticeRecord, Price: "3000", Validity: "2h",
		SellingPrice: "4000", Profile: "2jam", LockUser: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(got, `:put (",ntfc,3000,2h,4000,,Enable,")`) {
		t.Errorf("header put salah: %s", got[:40])
	}
	if !strings.Contains(got, `sname "$legacy-|-$time-|-$user-|-3000-|-$address-|-$mac-|-2h-|-2jam-|-$comment"`) {
		t.Errorf("script laporan salah: %s", got)
	}
}

func TestOnLoginScriptNone(t *testing.T) {
	got, err := OnLoginScript(OnLoginOptions{
		Mode: ExpiryNone, Price: "5000", LockUser: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(got, "interval=") {
		t.Fatalf("mode 0 tidak boleh punya logika expiry: %s", got)
	}
	if !strings.HasPrefix(got, `:put (",,5000,,,noexp,Enable,")`) {
		t.Errorf("header put mode 0 salah: %s", got)
	}
	if !strings.Contains(got, "mac-address=$mac") {
		t.Error("lock script harus disertakan")
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
	if meta.Mode != "rem" || meta.Price != "5000" || meta.Validity != "1d" || meta.Selling != "7000" || !meta.LockUser || meta.LockServer {
		t.Fatalf("legacy metadata tidak terbaca: %+v", meta)
	}

	v4Remc := `:put (",remc,5000,1d,7000,,Enable,"); { :local date ... comment=mikhmon }`
	meta = ParseOnLoginScript(v4Remc)
	if meta.Mode != "remc" || meta.Price != "5000" || meta.Validity != "1d" || meta.Selling != "7000" || !meta.LockUser || !meta.RecordSales {
		t.Fatalf("v4 remc metadata tidak terbaca: %+v", meta)
	}

	v4None := `:put (",,3000,,,noexp,Disable,")`
	meta = ParseOnLoginScript(v4None)
	if meta.Price != "3000" || meta.LockUser {
		t.Fatalf("v4 mode 0 metadata tidak terbaca: %+v", meta)
	}
}

func TestProfileMonitorScript(t *testing.T) {
	got, err := ProfileMonitorScript("1jam", ExpiryRemoveRecord)
	if err != nil {
		t.Fatal(err)
	}
	if strings.ContainsAny(got, "\r\n") {
		t.Fatal("profile monitor harus satu baris")
	}
	for _, want := range []string{
		"normdate", "dateint", "timeint",
		`profile="1jam"`, "[ /ip hotspot user remove $i ]",
		"[ /ip hotspot active remove [find where user=$name] ]",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("script tidak memuat %q", want)
		}
	}

	gotNtf, err := ProfileMonitorScript("2jam", ExpiryNotice)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(gotNtf, "[ /ip hotspot user set limit-uptime=1s $i ]") {
		t.Errorf("mode notice harus set limit-uptime=1s: %s", gotNtf)
	}

	gotNone, err := ProfileMonitorScript("free", ExpiryNone)
	if err != nil {
		t.Fatal(err)
	}
	if gotNone != "" {
		t.Errorf("mode none harus string kosong, got %q", gotNone)
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
