// 03 — Custom resource (katalog pemakai sendiri): CRUD + print + validate.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	ros "github.com/quiqxiq/roscatalog"
	"github.com/quiqxiq/roscatalog/rosrun"
)

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	addr := flag.String("addr", envOr("ROS_ADDR", "192.168.88.1:8728"), "alamat API RouterOS")
	user := flag.String("user", envOr("ROS_USER", "admin"), "username API")
	pass := flag.String("pass", envOr("ROS_PASS", ""), "password API")
	dryRun := flag.Bool("dry-run", true, "tampilkan rencana tanpa menulis ke router")
	flag.Parse()

	res := ros.Table("queue-custom", ros.GroupQueue, "Queue custom", "/queue/simple").
		Filter("name").
		List("name", "target", "max-limit", "comment").
		Count().
		Add(ros.P("name"), ros.P("target", "max-limit", "comment")).
		Set("target", "max-limit", "comment").
		Remove("hapus item")
	if err := res.Validate(); err != nil {
		panic(err)
	}

	// Buat catalog pribadi dari builtin + tambah custom resource.
	cat := ros.NewCatalog(ros.Builtin().All()...)
	if err := cat.Add(res); err != nil {
		panic(err)
	}
	got, err := cat.Find("queue-custom")
	if err != nil {
		panic(err)
	}

	req, err := got.Build("list", ros.Once(), ros.Params{})
	if err != nil {
		panic(err)
	}
	add, ok := got.Action("add")
	if !ok {
		panic("action add hilang")
	}
	addSent, err := add.Sentence("", map[string]string{"name": "roscatalog-ex-q1", "target": "10.0.0.1/32", "max-limit": "10M/10M"}, "")
	if err != nil {
		panic(err)
	}

	if *dryRun {
		fmt.Println("[dry-run] list sentence:", req.Sentence)
		fmt.Println("[dry-run] add sentence:", addSent)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	sess, err := rosrun.Open(ctx, rosrun.Config{
		Address:  *addr,
		Username: *user,
		Password: *pass,
		Timeout:  10 * time.Second,
		Status:   rosrun.StatusConfig{Disable: true},
	})
	if err != nil {
		panic(err)
	}
	defer sess.Close()

	rows, err := sess.Once(ctx, req)
	fmt.Println("baris:", len(rows), "err:", err)

	id, err := sess.Exec(ctx, addSent)
	if err != nil {
		panic(err)
	}
	fmt.Println("add q1 ok, id:", id)

	if rm, ok := got.Action("remove"); ok && id != "" {
		if rmSent, err := rm.Sentence(id, nil, ""); err == nil {
			_, _ = sess.Exec(ctx, rmSent)
			fmt.Println("cleanup remove ok")
		}
	}
}
