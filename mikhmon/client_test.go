package mikhmon

import (
	"context"
	"strings"
	"testing"

	ros "example.com/roscatalog"
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
