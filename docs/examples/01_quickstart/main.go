// 01 — Mulai cepat: dial RouterOS, daftar /interface, decode dengan fields.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	ros "github.com/quiqxiq/roscatalog"
	"github.com/quiqxiq/roscatalog/fields"
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
	flag.Parse()

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

	res, ok := ros.Lookup("interface")
	if !ok {
		panic("interface tidak ada di catalog")
	}
	req, err := res.Build("list", ros.Once(), ros.Params{})
	if err != nil {
		panic(err)
	}

	rows, err := sess.Once(ctx, req)
	if err != nil {
		panic(err)
	}
	for _, r := range rows {
		fmt.Printf("%-16s %s\n", r["name"], r["type"])
	}

	// Decode dengan struct opsional.
	ifs, _ := ros.DecodeAll[fields.InterfaceRow](rows)
	for _, r := range ifs {
		fmt.Printf("- %s running=%s mtu=%s\n", r.ID, r.Running, r.MTU)
	}
}
