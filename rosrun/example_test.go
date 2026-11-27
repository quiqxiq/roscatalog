package rosrun_test

import (
	"context"
	"log"
	"time"

	"github.com/quiqxiq/goros/v4"

	ros "github.com/quiqxiq/roscatalog"
	"github.com/quiqxiq/roscatalog/rosrun"
)

// Contoh ini hanya dikompilasi (tanpa // Output:), supaya potongan kode di
// USAGE.md selalu cocok dengan API yang sebenarnya.
func Example_usage() {
	c, err := routeros.Dial("192.168.88.1:8728", "api-user", "secret")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	run := rosrun.New(c)
	run.Logf = log.Printf
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	res, _ := ros.Lookup("hotspot-active")

	// stream interval 2s, poll 5s, poll + diff
	req, _ := res.Build("live", ros.StreamEvery(2*time.Second), ros.Params{})
	events, _ := run.Watch(ctx, req)
	_ = events
	req, _ = res.Build("live", ros.Poll(5*time.Second), ros.Params{})
	req, _ = res.Build("live", ros.Poll(5*time.Second).WithDiff(), ros.Params{})

	// listen + seed
	req, _ = res.Build("events", ros.Events().WithSeed(), ros.Params{})

	// native stream / poll
	mt, _ := ros.Lookup("interface-monitor-traffic")
	args := ros.Params{Args: map[string]string{"interface": "ether1,sfp1"}}
	req, _ = mt.Build("run", ros.Stream(), args)
	req, _ = mt.Build("run", ros.Poll(5*time.Second), args)

	// log
	lg, _ := ros.Lookup("log")
	req, _ = lg.Build("tail-only", ros.Stream(), ros.Params{})
	req, _ = lg.Build("tail-only", ros.Poll(5*time.Second), ros.Params{})

	// once, count, filter
	u, _ := ros.Lookup("hotspot-user")
	req, _ = u.Build("count", ros.Once(),
		ros.Params{Query: []ros.Query{{Prop: "profile", Value: "1jam"}}})
	n, _ := run.Count(ctx, req)
	_ = n
	req, _ = u.Build("list", ros.Once(), ros.Params{})
	rows, _ := run.Once(ctx, req)
	_ = rows

	// aksi tulis
	add, _ := u.Action("add")
	sent, _ := add.Sentence("", map[string]string{"name": "v001", "password": "x", "profile": "1jam"}, "")
	id, _ := run.Exec(ctx, sent)
	set, _ := u.Action("set")
	sent, _ = set.Sentence(id, map[string]string{"disabled": "yes"}, "")
	_, _ = run.Exec(ctx, sent)

	// fallback otomatis ke poll
	ev, _ := run.WatchAuto(ctx, res, "events", ros.Plan{Every: time.Second}, ros.Params{})
	_ = ev

	// validasi proplist
	req, _ = res.Build("list", ros.Once(), ros.Params{})
	rs, _ := run.Once(ctx, req)
	if missing := ros.NeverSeen(req.Fields, rs); len(missing) > 0 {
		log.Printf("%s: %v", res.ID, missing)
	}
}
