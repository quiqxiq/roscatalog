//go:build live

package rosrun_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	routeros "github.com/quiqxiq/goros/v4"

	ros "github.com/quiqxiq/roscatalog"
	"github.com/quiqxiq/roscatalog/rosrun"
)

// TestLiveListenRemoval membuktikan listen menerima event .dead=true ketika
// item dihapus, tanpa mengirim .proplist. Di RouterOS 7.24.6 proplist pada
// listen menekan event penghapusan.
//
// Uji ini hanya berjalan bila ROUTEROS_LIVE_ADDR diisi:
//
//	ROUTEROS_LIVE_ADDR=192.168.88.1:8728 ROUTEROS_LIVE_USER=admin \
//	ROUTEROS_LIVE_PASS=secret go test ./rosrun -run TestLiveListenRemoval -v
//
// Uji memakai /ip/firewall/address-list dengan nama list acak dan menghapus
// kembali entrinya.
func TestLiveListenRemoval(t *testing.T) {
	addr := os.Getenv("ROUTEROS_LIVE_ADDR")
	if addr == "" {
		t.Skip("set ROUTEROS_LIVE_ADDR untuk uji ke router sungguhan")
	}
	user := os.Getenv("ROUTEROS_LIVE_USER")
	if user == "" {
		user = "admin"
	}
	pass := os.Getenv("ROUTEROS_LIVE_PASS")

	c, err := routeros.Dial(addr, user, pass)
	if err != nil {
		t.Fatalf("dial %s: %v", addr, err)
	}
	defer c.Close()

	res, ok := ros.Lookup("fw-address-list")
	if !ok {
		t.Fatal("resource fw-address-list tidak ada")
	}
	req, err := res.Build("events", ros.Events(), ros.Params{})
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	for _, w := range req.Sentence {
		if strings.HasPrefix(w, "=.proplist=") {
			t.Fatalf("listen mengirim .proplist: %q", req.Sentence)
		}
	}

	run := rosrun.New(c)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	events, err := run.Watch(ctx, req)
	if err != nil {
		t.Fatalf("watch: %v", err)
	}

	list := fmt.Sprintf("roscatalog-live-%d", time.Now().UnixNano())
	const address = "198.51.100.99"
	defer func() {
		// Bersihkan semua sisa entri, termasuk bila uji gagal di tengah.
		r, err := c.RunArgs([]string{"/ip/firewall/address-list/print", "?list=" + list})
		if err != nil {
			return
		}
		for _, s := range r.Re {
			_, _ = c.RunArgs([]string{"/ip/firewall/address-list/remove", "=.id=" + s.Map[".id"]})
		}
	}()

	// Beri waktu goroutine stream mengirim /listen sebelum add dikirim.
	time.Sleep(500 * time.Millisecond)

	if _, err := c.RunArgs([]string{
		"/ip/firewall/address-list/add",
		"=list=" + list,
		"=address=" + address,
		"=comment=roscatalog live listen test",
	}); err != nil {
		t.Fatalf("add: %v", err)
	}
	added := waitLiveEvent(t, events, rosrun.KindAdded, list)
	if added.Row["address"] != address {
		t.Errorf("added address = %q, mau %q; row=%v", added.Row["address"], address, added.Row)
	}
	id := added.Row[".id"]
	if id == "" {
		t.Fatalf("event added tanpa .id: %v", added.Row)
	}

	if _, err := c.RunArgs([]string{"/ip/firewall/address-list/remove", "=.id=" + id}); err != nil {
		t.Fatalf("remove: %v", err)
	}
	removed := waitLiveEvent(t, events, rosrun.KindRemoved, list)
	if removed.Row["address"] != address {
		t.Errorf("removed address = %q, mau %q; row=%v", removed.Row["address"], address, removed.Row)
	}
}

func waitLiveEvent(t *testing.T, ch <-chan rosrun.Event, kind rosrun.Kind, list string) rosrun.Event {
	t.Helper()
	timer := time.NewTimer(10 * time.Second)
	defer timer.Stop()
	for {
		select {
		case ev, ok := <-ch:
			if !ok {
				t.Fatalf("channel event tertutup saat menunggu %s", kind)
			}
			if ev.Kind == rosrun.KindError {
				t.Fatalf("stream error: %v", ev.Err)
			}
			if ev.Kind == kind && ev.Row["list"] == list {
				return ev
			}
		case <-timer.C:
			t.Fatalf("timeout menunggu event %s untuk list %s", kind, list)
		}
	}
}
