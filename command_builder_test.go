package roscatalog

import (
	"reflect"
	"testing"
	"time"

	"github.com/quiqxiq/roscatalog/fields"
)

func TestPathNormalize(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"/interface", "/interface"},
		{"interface", "/interface"},
		{"/ip/address/", "/ip/address"},
		{"ip address", "/ip/address"},
		{"ip firewall filter", "/ip/firewall/filter"},
		{"/ping", "/ping"},
		{"", "/"},
	}
	for _, tc := range cases {
		got := Path(tc.input).path
		if got != tc.want {
			t.Errorf("Path(%q).path = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestPathPrintBasic(t *testing.T) {
	cmd := Path("/interface").Print()
	sent := cmd.Sentence()
	want := []string{"/interface/print"}
	if !reflect.DeepEqual(sent, want) {
		t.Fatalf("sent = %v, want %v", sent, want)
	}

	req, err := cmd.Build()
	if err != nil {
		t.Fatal(err)
	}
	if req.Delivery != DeliveryOnce || req.Mode != ModeOneShot {
		t.Fatalf("Delivery=%v, Mode=%v", req.Delivery, req.Mode)
	}
}

func TestPathPrintSwitchesAndProplist(t *testing.T) {
	cmd := Path("/interface").Print("stats").
		Switch("detail", "without-paging").
		Proplist("name", "rx-byte,tx-byte").
		From("0,1").
		Where("type", "=", "ether").
		Where("running", "=", "yes")

	sent := cmd.Sentence()
	// Word 0: /interface/print
	if sent[0] != "/interface/print" {
		t.Errorf("word 0 = %s", sent[0])
	}

	// Switches harus masuk
	for _, sw := range []string{"=detail=", "=stats=", "=without-paging="} {
		if !containsWord(sent, sw) {
			t.Errorf("sent tidak memuat switch %s: %v", sw, sent)
		}
	}

	// Proplist diawali titik
	if !containsWord(sent, "=.proplist=name,rx-byte,tx-byte") {
		t.Errorf("proplist salah: %v", sent)
	}

	// From
	if !containsWord(sent, "=from=0,1") {
		t.Errorf("from salah: %v", sent)
	}

	// Queries
	if !containsWord(sent, "?type=ether") || !containsWord(sent, "?running=yes") {
		t.Errorf("queries salah: %v", sent)
	}

	req := cmd.Request()
	if len(req.Fields) != 3 {
		t.Fatalf("Fields di request salah: %v", req.Fields)
	}
}

func TestPathPrintIntervalAndFollow(t *testing.T) {
	cmdLive := Path("/interface").Print("stats", "detail").
		Interval(time.Second)

	sent := cmdLive.Sentence()
	if !containsWord(sent, "=interval=1s") {
		t.Errorf("interval tidak ada di sentence: %v", sent)
	}

	req, err := cmdLive.Build()
	if err != nil {
		t.Fatal(err)
	}
	if req.Delivery != DeliveryStream || req.Mode != ModeInterval || req.Every != time.Second {
		t.Fatalf("req live salah: %+v", req)
	}

	cmdFollow := Path("/log").Print().Follow()
	reqFollow := cmdFollow.Request()
	if reqFollow.Delivery != DeliveryStream || reqFollow.Mode != ModeFollow || reqFollow.Shape != ShapeRows {
		t.Fatalf("req follow salah: %+v", reqFollow)
	}
	if !containsWord(reqFollow.Sentence, "=follow=") {
		t.Errorf("follow switch hilang: %v", reqFollow.Sentence)
	}
}

func TestPathSetWithMap(t *testing.T) {
	cmd := Path("/interface").Set("ether1", map[string]any{
		"mtu":      1500,
		"comment":  "WAN Uplink",
		"disabled": false,
	})

	sent := cmd.Sentence()
	if sent[0] != "/interface/set" {
		t.Fatalf("word 0 = %s", sent[0])
	}
	if !containsWord(sent, "=numbers=ether1") {
		t.Errorf("numbers ether1 hilang: %v", sent)
	}
	if !containsWord(sent, "=mtu=1500") {
		t.Errorf("mtu hilang: %v", sent)
	}
	if !containsWord(sent, "=comment=WAN Uplink") {
		t.Errorf("comment hilang: %v", sent)
	}
	if !containsWord(sent, "=disabled=no") {
		t.Errorf("disabled no salah: %v", sent)
	}
}

func TestPathSetWithStruct(t *testing.T) {
	row := fields.IPHotspotUserSet{
		Profile:     "1jam",
		LimitUptime: "1d",
	}
	cmd := Path("/ip/hotspot/user").Set("u1", row)
	sent := cmd.Sentence()

	if sent[0] != "/ip/hotspot/user/set" {
		t.Fatalf("word 0 = %s", sent[0])
	}
	if !containsWord(sent, "=profile=1jam") || !containsWord(sent, "=limit-uptime=1d") {
		t.Errorf("field struct tidak masuk: %v", sent)
	}
}

func TestPathSetWithVariadic(t *testing.T) {
	cmd := Path("/interface").Set("ether1", "mtu", 1460, "comment", "test")
	sent := cmd.Sentence()
	if !containsWord(sent, "=mtu=1460") || !containsWord(sent, "=comment=test") {
		t.Errorf("variadic set gagal: %v", sent)
	}
}

func TestPathAddWithMapAndStruct(t *testing.T) {
	cmd := Path("/ip/address").Add(map[string]any{
		"address":   "192.168.88.1/24",
		"interface": "ether2",
		"comment":   "LAN",
	})
	sent := cmd.Sentence()
	if sent[0] != "/ip/address/add" {
		t.Fatalf("word 0 = %s", sent[0])
	}
	if !containsWord(sent, "=address=192.168.88.1/24") || !containsWord(sent, "=interface=ether2") {
		t.Errorf("add fields salah: %v", sent)
	}

	addStruct := fields.IPHotspotUserAdd{
		Name:     "voucher01",
		Password: "secret",
		Profile:  "1jam",
	}
	cmdStruct := Path("/ip/hotspot/user").Add(addStruct)
	sentStruct := cmdStruct.Sentence()
	if !containsWord(sentStruct, "=name=voucher01") || !containsWord(sentStruct, "=password=secret") {
		t.Errorf("add struct gagal: %v", sentStruct)
	}
}

func TestPathRemoveEnableDisable(t *testing.T) {
	rem := Path("/ip/address").Remove("*1").Sentence()
	if rem[0] != "/ip/address/remove" || !containsWord(rem, "=numbers=*1") {
		t.Errorf("remove salah: %v", rem)
	}

	en := Path("/interface").Enable("ether1").Sentence()
	if en[0] != "/interface/enable" || !containsWord(en, "=numbers=ether1") {
		t.Errorf("enable salah: %v", en)
	}

	dis := Path("/interface").Disable("ether1").Sentence()
	if dis[0] != "/interface/disable" || !containsWord(dis, "=numbers=ether1") {
		t.Errorf("disable salah: %v", dis)
	}

	comm := Path("/interface").Comment("ether1", "LAN Office").Sentence()
	if comm[0] != "/interface/comment" || !containsWord(comm, "=numbers=ether1") || !containsWord(comm, "=comment=LAN Office") {
		t.Errorf("comment salah: %v", comm)
	}

	rst := Path("/interface").Reset("ether1").Sentence()
	if rst[0] != "/interface/reset" || !containsWord(rst, "=numbers=ether1") {
		t.Errorf("reset salah: %v", rst)
	}
}

func TestPathDoChildCommands(t *testing.T) {
	// 1. /interface/monitor-traffic ether1,ether2 once
	mt := Path("/interface").Do("monitor-traffic", "ether1,ether2").
		Switch("once").
		Sentence()

	if mt[0] != "/interface/monitor-traffic" {
		t.Fatalf("word 0 = %s", mt[0])
	}
	if !containsWord(mt, "=numbers=ether1,ether2") {
		t.Errorf("numbers ether1,ether2 hilang: %v", mt)
	}
	if !containsWord(mt, "=once=") {
		t.Errorf("switch once hilang: %v", mt)
	}

	// 2. /system/package/update/check-for-updates once
	chk := Path("/system/package/update").Do("check-for-updates").
		Switch("once").
		Sentence()

	if chk[0] != "/system/package/update/check-for-updates" || !containsWord(chk, "=once=") {
		t.Errorf("check-for-updates salah: %v", chk)
	}

	// 3. /ping 8.8.8.8 count=4
	png := Path("/ping").Do("", "8.8.8.8", map[string]any{"count": 4}).
		Sentence()

	if png[0] != "/ping" || !containsWord(png, "=numbers=8.8.8.8") || !containsWord(png, "=count=4") {
		t.Errorf("ping salah: %v", png)
	}
}

func TestRequestableInterface(t *testing.T) {
	var r1 Requestable = Path("/interface").Print()
	var r2 Requestable = Path("/interface").Set("ether1", "mtu", 1500)

	req1 := r1.Request()
	req2 := r2.Request()

	if req1.Sentence[0] != "/interface/print" {
		t.Errorf("req1 salah: %v", req1.Sentence)
	}
	if req2.Sentence[0] != "/interface/set" {
		t.Errorf("req2 salah: %v", req2.Sentence)
	}
}
