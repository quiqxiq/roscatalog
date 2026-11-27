package main

import (
	"testing"
	"time"

	ros "github.com/quiqxiq/roscatalog"
)

func TestBuildAllCatalogRequests(t *testing.T) {
	// 1. /ip/hotspot/active
	hsActive, ok := ros.Lookup("hotspot-active")
	if !ok {
		t.Fatal("hotspot-active tidak ditemukan di katalog")
	}
	// One-shot
	r1, err := hsActive.Build("list", ros.Once(), ros.Params{})
	if err != nil {
		t.Fatalf("hsActive list: %v", err)
	}
	t.Logf("1. hotspot-active list: %v", r1.Sentence)

	// Follow (tambahkan view tail jika belum ada via builder)
	if _, ok := hsActive.View("tail"); !ok {
		hsActive.Tail("user", "address", "mac-address", "uptime", "bytes-in", "bytes-out")
	}
	r2, err := hsActive.Build("tail", ros.Stream(), ros.Params{})
	if err != nil {
		t.Fatalf("hsActive tail: %v", err)
	}
	t.Logf("2. hotspot-active tail: %v", r2.Sentence)

	// Interval 1s
	r3, err := hsActive.Build("live", ros.StreamEvery(time.Second), ros.Params{})
	if err != nil {
		t.Fatalf("hsActive live: %v", err)
	}
	t.Logf("3. hotspot-active live: %v", r3.Sentence)

	// 2. /ip/hotspot/user
	hsUser, ok := ros.Lookup("hotspot-user")
	if !ok {
		t.Fatal("hotspot-user tidak ditemukan")
	}
	r4, err := hsUser.Build("list", ros.Once(), ros.Params{})
	if err != nil {
		t.Fatalf("hsUser list: %v", err)
	}
	t.Logf("4. hotspot-user list: %v", r4.Sentence)

	if _, ok := hsUser.View("tail"); !ok {
		hsUser.Tail("name", "profile", "uptime", "bytes-in", "bytes-out")
	}
	r5, err := hsUser.Build("tail", ros.Stream(), ros.Params{})
	if err != nil {
		t.Fatalf("hsUser tail: %v", err)
	}
	t.Logf("5. hotspot-user tail: %v", r5.Sentence)

	if _, ok := hsUser.View("live-stats"); !ok {
		hsUser.Live(time.Second, true, "name", "profile", "uptime", "bytes-in", "bytes-out")
	}
	r6, err := hsUser.Build("live", ros.StreamEvery(time.Second), ros.Params{})
	if err != nil {
		t.Fatalf("hsUser live: %v", err)
	}
	t.Logf("6. hotspot-user live: %v", r6.Sentence)

	// 3. /ppp/secret
	pppSec, ok := ros.Lookup("ppp-secret")
	if !ok {
		t.Fatal("ppp-secret tidak ditemukan")
	}
	r7, err := pppSec.Build("list", ros.Once(), ros.Params{})
	if err != nil {
		t.Fatalf("pppSec list: %v", err)
	}
	t.Logf("7. ppp-secret list: %v", r7.Sentence)

	if _, ok := pppSec.View("tail"); !ok {
		pppSec.Tail("name", "service", "profile", "local-address", "remote-address")
	}
	r8, err := pppSec.Build("tail", ros.Stream(), ros.Params{})
	if err != nil {
		t.Fatalf("pppSec tail: %v", err)
	}
	t.Logf("8. ppp-secret tail: %v", r8.Sentence)

	if _, ok := pppSec.View("live"); !ok {
		pppSec.Live(time.Second, true, "name", "service", "profile", "limit-bytes-in", "limit-bytes-out")
	}
	r9, err := pppSec.Build("live", ros.StreamEvery(time.Second), ros.Params{})
	if err != nil {
		t.Fatalf("pppSec live: %v", err)
	}
	t.Logf("9. ppp-secret live: %v", r9.Sentence)

	// 4. /ppp/active
	pppAct, ok := ros.Lookup("ppp-active")
	if !ok {
		t.Fatal("ppp-active tidak ditemukan")
	}
	r10, err := pppAct.Build("list", ros.Once(), ros.Params{})
	if err != nil {
		t.Fatalf("pppAct list: %v", err)
	}
	t.Logf("10. ppp-active list: %v", r10.Sentence)

	if _, ok := pppAct.View("tail"); !ok {
		pppAct.Tail("name", "service", "caller-id", "address", "uptime")
	}
	r11, err := pppAct.Build("tail", ros.Stream(), ros.Params{})
	if err != nil {
		t.Fatalf("pppAct tail: %v", err)
	}
	t.Logf("11. ppp-active tail: %v", r11.Sentence)

	if _, ok := pppAct.View("live"); !ok {
		pppAct.Live(time.Second, true, "name", "service", "caller-id", "address", "uptime")
	}
	r12, err := pppAct.Build("live", ros.StreamEvery(time.Second), ros.Params{})
	if err != nil {
		t.Fatalf("pppAct live: %v", err)
	}
	t.Logf("12. ppp-active live: %v", r12.Sentence)

	// 5. /log
	logRes, ok := ros.Lookup("log")
	if !ok {
		t.Fatal("log tidak ditemukan")
	}
	r13, err := logRes.Build("tail", ros.Stream(), ros.Params{})
	if err != nil {
		t.Fatalf("log tail: %v", err)
	}
	t.Logf("13. log tail: %v", r13.Sentence)

	if _, ok := logRes.View("live"); !ok {
		logRes.Live(time.Second, false, "time", "topics", "message")
	}
	r14, err := logRes.Build("live", ros.StreamEvery(time.Second), ros.Params{})
	if err != nil {
		t.Fatalf("log live: %v", err)
	}
	t.Logf("14. log live: %v", r14.Sentence)

	// 6. /system/resource
	sysRes, ok := ros.Lookup("system-resource")
	if !ok {
		t.Fatal("system-resource tidak ditemukan")
	}
	r15, err := sysRes.Build("live", ros.StreamEvery(time.Second), ros.Params{})
	if err != nil {
		t.Fatalf("system-resource live: %v", err)
	}
	t.Logf("15. system-resource live: %v", r15.Sentence)

	// 7. /interface
	iface, ok := ros.Lookup("interface")
	if !ok {
		t.Fatal("interface tidak ditemukan")
	}
	r16, err := iface.Build("list", ros.Once(), ros.Params{Proplist: []string{".id", "name", "type"}})
	if err != nil {
		t.Fatalf("interface list proplist: %v", err)
	}
	t.Logf("16. interface list proplist: %v", r16.Sentence)

	r17, err := iface.Build("live", ros.StreamEvery(time.Second), ros.Params{
		Flags:    []string{"detail"},
		Proplist: []string{"tx-byte", "rx-byte"},
	})
	if err != nil {
		t.Fatalf("interface live detail: %v", err)
	}
	t.Logf("17. interface live stats detail: %v", r17.Sentence)
}
