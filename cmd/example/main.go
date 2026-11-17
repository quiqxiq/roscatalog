// Contoh CLI: jalankan satu view katalog ke router lewat goros/v4.
//
//	go run ./cmd/example -addr 192.168.88.1:8728 -user admin -pass '' \
//	  -resource hotspot-active -view live -delivery stream -every 1s
//
//	go run ./cmd/example -resource hotspot-active -view live -delivery poll -every 500ms -diff
//	go run ./cmd/example -resource hotspot-active -view events -delivery stream -seed
//	go run ./cmd/example -resource interface-monitor-traffic -view run -arg interface=ether1
//	go run ./cmd/example -resource hotspot-user -view count -delivery once -q profile=1jam
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/quiqxiq/goros/v4"

	ros "example.com/roscatalog"
	"example.com/roscatalog/rosrun"
)

func main() {
	var (
		addr     = flag.String("addr", "192.168.88.1:8728", "alamat API router")
		user     = flag.String("user", "admin", "user API")
		pass     = flag.String("pass", "", "password API")
		resource = flag.String("resource", "hotspot-active", "ID resource katalog")
		view     = flag.String("view", "live", "nama view")
		ver      = flag.String("ver", "", "v6 atau v7 (kosong = semua)")
		delivery = flag.String("delivery", "", "once | stream | poll (kosong = bawaan view)")
		every    = flag.Duration("every", time.Second, "interval stream atau jeda poll")
		diff     = flag.Bool("diff", false, "poll: keluarkan added/changed/removed")
		seed     = flag.Bool("seed", false, "listen: ambil snapshot awal")
	)
	params := ros.Params{Args: map[string]string{}}
	flag.Func("arg", "argumen command native, mis. interface=ether1 (boleh berulang)", func(s string) error {
		k, v, _ := strings.Cut(s, "=")
		params.Args[k] = v
		return nil
	})
	flag.Func("q", "query sisi server prop=nilai (boleh berulang)", func(s string) error {
		k, v, _ := strings.Cut(s, "=")
		params.Query = append(params.Query, ros.Query{Prop: k, Value: v})
		return nil
	})
	flag.Parse()

	res, ok := ros.Lookup(*resource)
	if !ok {
		log.Fatalf("resource %q tidak ada", *resource)
	}
	plan := ros.Plan{Delivery: ros.Delivery(*delivery), Every: *every, Diff: *diff, Seed: *seed}
	req, err := res.Build(*view, ros.Version(*ver), plan, params)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("kirim: %q (delivery=%s mode=%s)", req.Sentence, req.Delivery, req.Mode)

	c, err := routeros.Dial(*addr, *user, *pass)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	run := rosrun.New(c) // menyalakan async; satu Runner per koneksi
	run.Logf = log.Printf

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	events, err := run.Watch(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	for ev := range events {
		switch ev.Kind {
		case rosrun.KindSnapshot:
			fmt.Printf("%s snapshot %d baris\n", ev.At.Format("15:04:05.000"), len(ev.Rows))
		case rosrun.KindCount:
			fmt.Printf("jumlah: %d\n", ev.Count)
		case rosrun.KindError:
			log.Fatalf("berhenti: %v", ev.Err)
		default:
			fmt.Printf("%s %-7s %v\n", ev.At.Format("15:04:05.000"), ev.Kind, ev.Row)
		}
	}
}
