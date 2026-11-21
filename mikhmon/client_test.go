package mikhmon

import (
	"context"
	"strings"
	"testing"

	ros "example.com/roscatalog"
	"example.com/roscatalog/fields"
	"example.com/roscatalog/rosrun"
)

var _ Executor = (*rosrun.Runner)(nil)

type fakeExec struct {
	onceRows  []map[string]string
	onceErr   error
	count     int
	countErr  error
	execRet   string
	execErr   error
	onceReqs  []ros.Request
	countReqs []ros.Request
	execs     [][]string
}

func (f *fakeExec) Once(_ context.Context, req ros.Request) ([]map[string]string, error) {
	f.onceReqs = append(f.onceReqs, req)
	return f.onceRows, f.onceErr
}

func (f *fakeExec) Count(_ context.Context, req ros.Request) (int, error) {
	f.countReqs = append(f.countReqs, req)
	return f.count, f.countErr
}

func (f *fakeExec) Exec(_ context.Context, sentence []string) (string, error) {
	f.execs = append(f.execs, append([]string(nil), sentence...))
	return f.execRet, f.execErr
}

func joined(s []string) string { return strings.Join(s, " ") }

func TestAddUser(t *testing.T) {
	f := &fakeExec{execRet: "*A1"}
	c := New(f)
	id, err := c.AddUser(context.Background(), User{
		Name: "vc001", Password: "vc001", Profile: "1jam",
		Server: "all", Comment: MikhmonComment("vc001", "vc001", "batch"),
		LimitUptime: "1d",
	})
	if err != nil {
		t.Fatal(err)
	}
	if id != "*A1" {
		t.Fatalf("id=%q", id)
	}
	got := joined(f.execs[0])
	for _, want := range []string{
		"/ip/hotspot/user/add", "=name=vc001", "=password=vc001",
		"=profile=1jam", "=comment=vc-batch", "=limit-uptime=1d",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("sentence tidak memuat %q: %s", want, got)
		}
	}
}

func TestAddUserRequiresName(t *testing.T) {
	c := New(&fakeExec{})
	if _, err := c.AddUser(context.Background(), User{}); err == nil {
		t.Fatal("nama kosong harus error")
	}
}

func TestAddProfile(t *testing.T) {
	f := &fakeExec{execRet: "*P1"}
	c := New(f)
	_, err := c.AddProfile(context.Background(), Profile{
		Name: "1jam", RateLimit: "1M/1M", SharedUsers: "1",
		OnLogin: ":put 1",
	})
	if err != nil {
		t.Fatal(err)
	}
	got := joined(f.execs[0])
	for _, want := range []string{
		"/ip/hotspot/user/profile/add", "=name=1jam",
		"=status-autorefresh=1m", "=on-login=:put 1",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("sentence tidak memuat %q: %s", want, got)
		}
	}
}

func TestSetOnLogin(t *testing.T) {
	f := &fakeExec{onceRows: []map[string]string{{".id": "*P1"}}}
	c := New(f)
	if err := c.SetOnLogin(context.Background(), "1jam", ":put 1"); err != nil {
		t.Fatal(err)
	}
	got := joined(f.execs[0])
	if !strings.Contains(got, "/ip/hotspot/user/profile/set") ||
		!strings.Contains(got, "=.id=*P1") ||
		!strings.Contains(got, "=on-login=:put 1") {
		t.Fatalf("sentence salah: %s", got)
	}
}

func TestSetOnLoginRejectsMultiline(t *testing.T) {
	c := New(&fakeExec{})
	if err := c.SetOnLogin(context.Background(), "1jam", "a\nb"); err == nil {
		t.Fatal("on-login multi-baris harus error")
	}
}

func TestInstallExpireMonitorAdd(t *testing.T) {
	f := &fakeExec{execRet: "*S1"}
	c := New(f)
	if err := c.InstallExpireMonitor(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(f.execs) != 1 {
		t.Fatalf("exec=%d, want 1", len(f.execs))
	}
	got := joined(f.execs[0])
	for _, want := range []string{
		"/system/scheduler/add", "=name=Mikhmon-Expire-Monitor",
		"=interval=00:01:00", "=on-event=", "normdate",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("sentence tidak memuat %q", want)
		}
	}
}

func TestInstallExpireMonitorUpdate(t *testing.T) {
	f := &fakeExec{onceRows: []map[string]string{{".id": "*S1"}}}
	c := New(f)
	if err := c.InstallExpireMonitor(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(f.execs) != 2 {
		t.Fatalf("exec=%d, want 2 (set + enable)", len(f.execs))
	}
	if got := joined(f.execs[0]); !strings.Contains(got, "/system/scheduler/set") {
		t.Errorf("exec pertama bukan set: %s", got)
	}
	if got := joined(f.execs[1]); !strings.Contains(got, "/system/scheduler/enable") {
		t.Errorf("exec kedua bukan enable: %s", got)
	}
}

func TestAddSale(t *testing.T) {
	f := &fakeExec{execRet: "*R1"}
	c := New(f)
	_, err := c.AddSale(context.Background(), Sale{
		Date: "nov/18/2026", Time: "10:00:00", User: "vc001",
		Price: "5000", Address: "198.51.100.10", MAC: "AA:BB:CC:DD:EE:FF",
		Validity: "1d", Profile: "1jam", Comment: "vc-123-11.18.26",
	})
	if err != nil {
		t.Fatal(err)
	}
	got := joined(f.execs[0])
	for _, want := range []string{
		"/system/script/add", "=source=nov/18/2026",
		"=owner=nov2026", "=comment=mikhmon",
		"=name=nov/18/2026-|-10:00:00-|-vc001-|-5000-|-",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("sentence tidak memuat %q: %s", want, got)
		}
	}
}

func TestNormalizeSalesOwner(t *testing.T) {
	f := &fakeExec{onceRows: []map[string]string{
		{".id": "*R1", "source": "nov/18/2026"},
		{".id": "*R2", "source": "2026-11-18"},
	}}
	c := New(f)
	n, err := c.NormalizeSalesOwner(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 || len(f.execs) != 2 {
		t.Fatalf("n=%d exec=%d, want 2", n, len(f.execs))
	}
	for _, sent := range f.execs {
		got := joined(sent)
		if !strings.Contains(got, "/system/script/set") ||
			!strings.Contains(got, "=owner=nov2026") {
			t.Errorf("sentence normalisasi salah: %s", got)
		}
	}
	if got := joined(f.onceReqs[0].Sentence); !strings.Contains(got, "?owner=*sys") ||
		!strings.Contains(got, "?comment=mikhmon") {
		t.Errorf("query normalisasi salah: %s", got)
	}
}

func TestNormalizeSalesOwnerExplicit(t *testing.T) {
	f := &fakeExec{onceRows: []map[string]string{{".id": "*R1", "source": "bad"}}}
	c := New(f)
	n, err := c.NormalizeSalesOwner(context.Background(), "nov2026")
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("n=%d", n)
	}
}

func TestNormalizeSalesOwnerEmpty(t *testing.T) {
	c := New(&fakeExec{})
	n, err := c.NormalizeSalesOwner(context.Background(), "nov2026")
	if err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("n=%d", n)
	}
}

func TestRemoveExpiredUsers(t *testing.T) {
	f := &fakeExec{onceRows: []map[string]string{
		{".id": "*U1"}, {".id": "*U2"},
	}}
	c := New(f)
	n, err := c.RemoveExpiredUsers(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 || len(f.execs) != 2 {
		t.Fatalf("n=%d exec=%d", n, len(f.execs))
	}
	for _, sent := range f.execs {
		if !strings.Contains(joined(sent), "/ip/hotspot/user/remove") {
			t.Errorf("bukan remove: %v", sent)
		}
	}
	if got := joined(f.onceReqs[0].Sentence); !strings.Contains(got, "?limit-uptime=1s") {
		t.Errorf("query limit-uptime salah: %s", got)
	}
}

func TestSaveQuickPrintAdd(t *testing.T) {
	f := &fakeExec{execRet: "*Q1"}
	c := New(f)
	_, err := c.SaveQuickPrint(context.Background(), QuickPrintPreset{
		Package: "1jam", Server: "all", UserMode: "vc",
	})
	if err != nil {
		t.Fatal(err)
	}
	got := joined(f.execs[0])
	for _, want := range []string{
		"/system/script/add", "=name=Quick_Print_1jam",
		"=comment=QuickPrintMikhmon", "=source=#1jam#all#",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("sentence tidak memuat %q: %s", want, got)
		}
	}
}

func TestSaveQuickPrintUpdate(t *testing.T) {
	f := &fakeExec{onceRows: []map[string]string{{".id": "*Q1"}}}
	c := New(f)
	_, err := c.SaveQuickPrint(context.Background(), QuickPrintPreset{Package: "1jam"})
	if err != nil {
		t.Fatal(err)
	}
	got := joined(f.execs[0])
	if !strings.Contains(got, "/system/script/set") ||
		!strings.Contains(got, "=.id=*Q1") {
		t.Fatalf("bukan update: %s", got)
	}
}

func TestKickUser(t *testing.T) {
	f := &fakeExec{onceRows: []map[string]string{{".id": "*A1"}}}
	c := New(f)
	n, err := c.KickUser(context.Background(), "vc001")
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 || len(f.execs) != 1 {
		t.Fatalf("n=%d exec=%d", n, len(f.execs))
	}
	if got := joined(f.execs[0]); !strings.Contains(got, "/ip/hotspot/active/remove") {
		t.Fatalf("bukan remove active: %s", got)
	}
}

func TestInstallProfileMonitorAdd(t *testing.T) {
	f := &fakeExec{execRet: "*S1"}
	c := New(f)
	if err := c.InstallProfileMonitor(context.Background(), "1jam", ExpiryRemoveRecord); err != nil {
		t.Fatal(err)
	}
	if len(f.execs) != 1 {
		t.Fatalf("exec=%d, want 1", len(f.execs))
	}
	got := joined(f.execs[0])
	for _, want := range []string{
		"/system/scheduler/add", "=name=1jam",
		"=comment=Monitor Profile 1jam", "=interval=00:02:30",
		"=on-event=", `profile="1jam"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("sentence tidak memuat %q: %s", want, got)
		}
	}
}

func TestInstallProfileMonitorUpdate(t *testing.T) {
	f := &fakeExec{onceRows: []map[string]string{{".id": "*S2"}}}
	c := New(f)
	if err := c.InstallProfileMonitor(context.Background(), "1jam", ExpiryNoticeRecord); err != nil {
		t.Fatal(err)
	}
	if len(f.execs) != 2 {
		t.Fatalf("exec=%d, want 2 (set + enable)", len(f.execs))
	}
	if got := joined(f.execs[0]); !strings.Contains(got, "/system/scheduler/set") || !strings.Contains(got, "=.id=*S2") {
		t.Errorf("exec pertama bukan set: %s", got)
	}
	if got := joined(f.execs[1]); !strings.Contains(got, "/system/scheduler/enable") || !strings.Contains(got, "=.id=*S2") {
		t.Errorf("exec kedua bukan enable: %s", got)
	}
}

func TestInstallProfileMonitorNoneRemoves(t *testing.T) {
	f := &fakeExec{onceRows: []map[string]string{{".id": "*S3"}}}
	c := New(f)
	if err := c.InstallProfileMonitor(context.Background(), "1jam", ExpiryNone); err != nil {
		t.Fatal(err)
	}
	if len(f.execs) != 1 {
		t.Fatalf("exec=%d, want 1 (remove)", len(f.execs))
	}
	if got := joined(f.execs[0]); !strings.Contains(got, "/system/scheduler/remove") || !strings.Contains(got, "=.id=*S3") {
		t.Errorf("sentence bukan remove: %s", got)
	}
}

func TestRemoveProfileMonitor(t *testing.T) {
	f := &fakeExec{onceRows: []map[string]string{{".id": "*S4"}}}
	c := New(f)
	if err := c.RemoveProfileMonitor(context.Background(), "1jam"); err != nil {
		t.Fatal(err)
	}
	if len(f.execs) != 1 {
		t.Fatalf("exec=%d, want 1", len(f.execs))
	}
	if got := joined(f.execs[0]); !strings.Contains(got, "/system/scheduler/remove") || !strings.Contains(got, "=.id=*S4") {
		t.Errorf("sentence salah: %s", got)
	}
}

func TestProfileMonitor(t *testing.T) {
	f := &fakeExec{onceRows: []map[string]string{{".id": "*S5", "name": "1jam"}}}
	c := New(f)
	row, err := c.ProfileMonitor(context.Background(), "1jam")
	if err != nil {
		t.Fatal(err)
	}
	if row[".id"] != "*S5" || row["name"] != "1jam" {
		t.Fatalf("row salah: %v", row)
	}
}

func TestRemoveProfile(t *testing.T) {
	f := &fakeExec{
		onceRows: []map[string]string{{".id": "*P1"}}, // for profile lookup
	}
	c := New(f)
	if err := c.RemoveProfile(context.Background(), "1jam"); err != nil {
		t.Fatal(err)
	}
	// first action is remove profile
	if got := joined(f.execs[0]); !strings.Contains(got, "/ip/hotspot/user/profile/remove") || !strings.Contains(got, "=.id=*P1") {
		t.Errorf("exec remove profile salah: %s", got)
	}
}

func TestUserToFieldsAndFromVoucher(t *testing.T) {
	u := User{
		Name:            "vc101",
		Password:        "pass101",
		Profile:         "paket1",
		Server:          "all",
		Comment:         "vc-batch-1",
		Address:         "10.0.0.10",
		MACAddress:      "00:11:22:33:44:55",
		LimitUptime:     "1d",
		LimitBytesTotal: "1048576000",
		Email:           "test@example.com",
	}
	f := u.ToFields()
	if f.Name != u.Name || f.Password != u.Password || f.LimitUptime != u.LimitUptime {
		t.Fatalf("ToFields salah: %+v", f)
	}

	vRow := fields.MikhmonVoucherRow{
		ID:              "*101",
		Name:            u.Name,
		Password:        u.Password,
		Profile:         u.Profile,
		LimitUptime:     u.LimitUptime,
		LimitBytesTotal: u.LimitBytesTotal,
		Comment:         u.Comment,
		Server:          u.Server,
		Address:         u.Address,
		MACAddress:      u.MACAddress,
	}
	uFromV := UserFromVoucher(vRow)
	if uFromV.Name != u.Name || uFromV.Password != u.Password || uFromV.Profile != u.Profile {
		t.Fatalf("UserFromVoucher salah: %+v", uFromV)
	}
}

func TestProfileToFields(t *testing.T) {
	p := Profile{
		Name:        "paket1",
		RateLimit:   "2M/2M",
		SharedUsers: "1",
		OnLogin:     ":put test",
	}
	f := p.ToFields()
	if f.Name != "paket1" || f.StatusAutorefresh != "1m" || f.RateLimit != "2M/2M" {
		t.Fatalf("ToFields profile salah: %+v", f)
	}
}

func TestClientTypedMethods(t *testing.T) {
	f := &fakeExec{
		onceRows: []map[string]string{
			{
				".id":     "*1",
				"name":    "nov/18/2026-|-10:00:00-|-u1-|-5000-|-10.0.0.1-|-AA:BB-|-1d-|-p1-|-vc-1",
				"owner":   "nov2026",
				"source":  "nov/18/2026",
				"comment": "mikhmon",
			},
		},
	}
	c := New(f)
	sales, err := c.SalesRows(context.Background(), "nov2026")
	if err != nil || len(sales) != 1 {
		t.Fatalf("SalesRows failed: err=%v, len=%d", err, len(sales))
	}
	if sales[0].Owner != "nov2026" || sales[0].Comment != "mikhmon" {
		t.Fatalf("unexpected sales row: %+v", sales[0])
	}

	f.onceRows = []map[string]string{
		{
			".id":     "*2",
			"name":    "Quick_Print_Paket1",
			"source":  "#Paket1#all#vc#4#vc#mix#paket1#1d#0#c#1d#5000_7000#Disable",
			"comment": "QuickPrintMikhmon",
		},
	}
	qps, err := c.QuickPrintRows(context.Background())
	if err != nil || len(qps) != 1 {
		t.Fatalf("QuickPrintRows failed: err=%v, len=%d", err, len(qps))
	}
	if qps[0].Name != "Quick_Print_Paket1" {
		t.Fatalf("unexpected quickprint row: %+v", qps[0])
	}

	f.onceRows = []map[string]string{
		{
			".id":          "*3",
			"name":         "vc101",
			"password":     "1234",
			"profile":      "paket1",
			"limit-uptime": "1d",
			"comment":      "vc-batch-1",
		},
	}
	vouchers, err := c.VouchersTyped(context.Background(), "vc-batch-1")
	if err != nil || len(vouchers) != 1 {
		t.Fatalf("VouchersTyped failed: err=%v, len=%d", err, len(vouchers))
	}
	if vouchers[0].Password != "1234" {
		t.Fatalf("unexpected voucher row: %+v", vouchers[0])
	}

	f.onceRows = []map[string]string{
		{
			".id":        "*4",
			"name":       "paket1",
			"start-time": "01:00:00",
			"interval":   "00:02:30",
			"on-event":   ":put test",
			"comment":    "Monitor Profile paket1",
			"disabled":   "no",
		},
	}
	pMon, err := c.ProfileMonitorRow(context.Background(), "paket1")
	if err != nil || pMon.Comment != "Monitor Profile paket1" {
		t.Fatalf("ProfileMonitorRow failed: err=%v, row=%+v", err, pMon)
	}

	f.onceRows = []map[string]string{
		{
			".id":        "*5",
			"name":       ExpireMonitorName,
			"start-time": "00:00:00",
			"interval":   "00:01:00",
			"on-event":   ":put test",
			"comment":    ExpireMonitorComment,
			"disabled":   "no",
		},
	}
	expMon, err := c.ExpireMonitorRow(context.Background())
	if err != nil || expMon.Comment != ExpireMonitorComment {
		t.Fatalf("ExpireMonitorRow failed: err=%v, row=%+v", err, expMon)
	}
}
