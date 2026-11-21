// 01 — Mulai cepat: dial CHR, daftar /interface, decode dengan fields.
package main

import (
	"context"
	"fmt"
	"time"

	ros "example.com/roscatalog"
	"example.com/roscatalog/fields"
	"example.com/roscatalog/rosrun"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

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

	res, ok := ros.Lookup("interface")
	if !ok {
		panic("interface tidak ada di catalog")
	}
	req, err := res.Build("list", ros.Plan{Delivery: ros.DeliveryOnce}, ros.Params{})
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
