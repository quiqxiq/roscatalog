// 03 — Custom resource (katalog pemakai sendiri): CRUD + print + validate.
package main

import (
	"context"
	"fmt"
	"time"

	ros "example.com/roscatalog"
	"example.com/roscatalog/rosrun"
)

func main() {
	res := ros.Table("queue-custom", ros.GroupQueue, "Queue custom", "/queue/custom").
		Filter("name").
		List("name", "target", "rate", "comment").
		Count().
		Add(ros.P("name"), ros.P("target", "rate", "comment")).
		Set("target", "rate", "comment").
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

	ctx := context.Background()
	sess, err := rosrun.Open(ctx, rosrun.Config{
		Address:  "192.168.88.1:8728",
		Username: "admin",
		Password: "r00t",
		Timeout:  10 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	defer sess.Close()

	req, err := got.Build("list", ros.Once(), ros.Params{})
	if err != nil {
		panic(err)
	}
	rows, err := sess.Once(ctx, req)
	fmt.Println("baris:", len(rows), "err:", err)

	// Aksi tulis: add.
	add, ok := got.Action("add")
	if !ok {
		panic("action add hilang")
	}
	sent, err := add.Sentence("", map[string]string{"name": "q1", "target": "10.0.0.1", "rate": "10M"}, "")
	if err != nil {
		panic(err)
	}
	if _, err := sess.Exec(ctx, sent); err != nil {
		panic(err)
	}
	fmt.Println("add q1 ok")
}
