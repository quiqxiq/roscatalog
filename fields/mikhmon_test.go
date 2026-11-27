package fields_test

import (
	"reflect"
	"testing"

	ros "github.com/quiqxiq/roscatalog"
	"github.com/quiqxiq/roscatalog/fields"
)

func TestMikhmonSalesFields(t *testing.T) {
	names := ros.FieldNames[fields.MikhmonSalesRow]()
	wantPrefix := []string{".id", "name", "owner", "source", "comment"}
	for i, want := range wantPrefix {
		if i >= len(names) || names[i] != want {
			t.Fatalf("FieldNames[MikhmonSalesRow][%d] = %q, want %q (all: %v)", i, names[i], want, names)
		}
	}

	raw := map[string]string{
		".id":     "*10",
		"name":    "nov/21/2026-|-12:00:00-|-user1-|-5000-|-10.0.0.2-|-00:11:22:33:44:55-|-1d-|-paket1-|-vc-101",
		"owner":   "nov2026",
		"source":  "nov/21/2026",
		"comment": "mikhmon",
	}

	row, err := ros.Decode[fields.MikhmonSalesRow](raw)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if row.ID != "*10" || row.Owner != "nov2026" || row.Source != "nov/21/2026" || row.Comment != "mikhmon" {
		t.Fatalf("unexpected decode result: %+v", row)
	}

	add := fields.MikhmonSalesAdd{
		Name:    raw["name"],
		Source:  raw["source"],
		Owner:   raw["owner"],
		Comment: raw["comment"],
	}
	args, err := ros.ArgsOf(add)
	if err != nil {
		t.Fatalf("ArgsOf failed: %v", err)
	}
	if args["name"] != raw["name"] || args["owner"] != "nov2026" || args["comment"] != "mikhmon" {
		t.Fatalf("unexpected args: %v", args)
	}
}

func TestMikhmonVoucherFields(t *testing.T) {
	names := ros.FieldNames[fields.MikhmonVoucherRow]()
	wantPrefix := []string{".id", "name", "password", "profile", "limit-uptime", "limit-bytes-total", "comment"}
	for i, want := range wantPrefix {
		if i >= len(names) || names[i] != want {
			t.Fatalf("FieldNames[MikhmonVoucherRow][%d] = %q, want %q", i, names[i], want)
		}
	}

	raw := map[string]string{
		".id":               "*1",
		"name":              "vc101",
		"password":          "vc101",
		"profile":           "paket-1hari",
		"limit-uptime":      "1d",
		"limit-bytes-total": "1048576000",
		"comment":           "vc-batch-1",
	}
	v, err := ros.Decode[fields.MikhmonVoucherRow](raw)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if v.Name != "vc101" || v.Password != "vc101" || v.Profile != "paket-1hari" || v.LimitUptime != "1d" {
		t.Fatalf("unexpected decode result: %+v", v)
	}

	uAdd := fields.MikhmonUserAdd{
		Name:            "vc101",
		Password:        "vc101",
		Profile:         "paket-1hari",
		Comment:         "vc-batch-1",
		LimitUptime:     "1d",
		LimitBytesTotal: "1048576000",
	}
	args, err := ros.ArgsOf(uAdd)
	if err != nil {
		t.Fatalf("ArgsOf failed: %v", err)
	}
	if args["name"] != "vc101" || args["password"] != "vc101" || args["limit-uptime"] != "1d" {
		t.Fatalf("unexpected args: %v", args)
	}
}

func TestMikhmonProfileFields(t *testing.T) {
	raw := map[string]string{
		".id":                "*2",
		"name":               "paket-1hari",
		"rate-limit":         "2M/2M",
		"shared-users":       "1",
		"status-autorefresh": "1m",
		"on-login":           ":put (1,2,3)",
		"comment":            "",
	}
	p, err := ros.Decode[fields.MikhmonProfileRow](raw)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if p.Name != "paket-1hari" || p.RateLimit != "2M/2M" || p.OnLogin != ":put (1,2,3)" {
		t.Fatalf("unexpected decode result: %+v", p)
	}

	pAdd := fields.MikhmonProfileAdd{
		Name:              "paket-1hari",
		RateLimit:         "2M/2M",
		SharedUsers:       "1",
		StatusAutorefresh: "1m",
		OnLogin:           ":put (1,2,3)",
	}
	args, err := ros.ArgsOf(pAdd)
	if err != nil {
		t.Fatalf("ArgsOf failed: %v", err)
	}
	if args["name"] != "paket-1hari" || args["rate-limit"] != "2M/2M" || args["status-autorefresh"] != "1m" {
		t.Fatalf("unexpected args: %v", args)
	}
}

func TestMikhmonSchedulerFields(t *testing.T) {
	raw := map[string]string{
		".id":        "*3",
		"name":       "Mikhmon-Expire-Monitor",
		"start-time": "00:00:00",
		"interval":   "00:01:00",
		"on-event":   ":put test",
		"comment":    "Mikhmon Expire Monitor",
		"disabled":   "no",
	}
	s, err := ros.Decode[fields.MikhmonSchedulerRow](raw)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if s.Name != "Mikhmon-Expire-Monitor" || s.Interval != "00:01:00" || s.Comment != "Mikhmon Expire Monitor" {
		t.Fatalf("unexpected decode result: %+v", s)
	}

	schedAdd := fields.MikhmonSchedulerAdd{
		Name:      "Mikhmon-Expire-Monitor",
		StartTime: "00:00:00",
		Interval:  "00:01:00",
		OnEvent:   ":put test",
		Comment:   "Mikhmon Expire Monitor",
		Disabled:  "no",
	}
	args, err := ros.ArgsOf(schedAdd)
	if err != nil {
		t.Fatalf("ArgsOf failed: %v", err)
	}
	if args["name"] != "Mikhmon-Expire-Monitor" || args["interval"] != "00:01:00" || args["disabled"] != "no" {
		t.Fatalf("unexpected args: %v", args)
	}
}

func TestMikhmonQuickPrintFields(t *testing.T) {
	raw := map[string]string{
		".id":     "*4",
		"name":    "Quick_Print_Paket-Hemat",
		"source":  "#Paket Hemat#all#vc#4#vc#1234#paket1#1h#100M#comment#1d#5000_7000#lock",
		"comment": "QuickPrintMikhmon",
	}
	qp, err := ros.Decode[fields.MikhmonQuickPrintRow](raw)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if qp.Name != "Quick_Print_Paket-Hemat" || qp.Comment != "QuickPrintMikhmon" {
		t.Fatalf("unexpected decode result: %+v", qp)
	}

	qpAdd := fields.MikhmonQuickPrintAdd{
		Name:    raw["name"],
		Source:  raw["source"],
		Comment: raw["comment"],
	}
	args, err := ros.ArgsOf(qpAdd)
	if err != nil {
		t.Fatalf("ArgsOf failed: %v", err)
	}
	if args["name"] != raw["name"] || args["source"] != raw["source"] {
		t.Fatalf("unexpected args: %v", args)
	}
}

func TestMikhmonActiveFields(t *testing.T) {
	raw := map[string]string{
		".id":         "*5",
		"server":      "hotspot1",
		"user":        "vc101",
		"address":     "10.0.0.50",
		"mac-address": "00:11:22:33:44:55",
		"uptime":      "01:23:45",
		"bytes-in":    "1048576",
		"bytes-out":   "2097152",
		"login-by":    "http-chap",
	}
	act, err := ros.Decode[fields.MikhmonActiveRow](raw)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if act.User != "vc101" || act.Address != "10.0.0.50" || act.BytesIn != "1048576" || act.LoginBy != "http-chap" {
		t.Fatalf("unexpected decode result: %+v", act)
	}

	all, err := ros.DecodeAll[fields.MikhmonActiveRow]([]map[string]string{raw, raw})
	if err != nil || len(all) != 2 {
		t.Fatalf("DecodeAll failed: err=%v, len=%d", err, len(all))
	}
	if !reflect.DeepEqual(all[0], all[1]) {
		t.Fatalf("elements should be identical")
	}
}
