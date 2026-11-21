# 01 — Mulai

Install (module-local):

```bash
go get example.com/roscatalog@latest
go get example.com/roscatalog/rosrun@latest
```

Pembuatan session — pakai `rosrun.Open` (satu koneksi asinkron + auto-reconnect):

```go
package main

import (
	"context"
	"fmt"
	"time"

	"example.com/roscatalog"
	"example.com/roscatalog/rosrun"
)

func main() {
	ctx := context.Background()

	// Koneksi satu per router, asynchrony — runner dial, async, autoreconnect.
	sess, err := rosrun.Open(ctx, rosrun.Config{
		Address:  "192.168.88.1:8728", // CHR di 192.168.122.189:8728 (jalurs?)
		Username: "admin",
		Password: "r00t",
		Timeout:  10 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	defer sess.Close()

	// Daftar satu view: /interface/print (list, proplist CHR-terverifikasi).
	res, ok := ros.Lookup("interface")
	if !ok {
		panic("interface tidak ada")
	}
	rows, err := sess.Once(ctx, must(res.Build("list", ros.Plan{Delivery: ros.DeliveryOnce})))
	for _, r := range rows {
		fmt.Println(r["name"], r["type"])
	}
}

func must(req ros.Request, err error) ros.Request {
	if err != nil {
		panic(err)
	}
	return req
}
```

Catatan penting:

- `res.Build(view, plan, params)` → `(ros.Request, error)`.
- Session = satu koneksi asinkron per router (pengganti dualconn).
- Default delivery baca: **realtime** (`Stream`); minta `Once()` bila satu dump saja.
- Command yang tidak ada di catalog → ditangkal: lihat [Custom Resources](03-custom-resources.md)
  (`Table`, `Tool`, `Single`) dan [Errors](09-errors.md) (`errors.Is(err, ros.ErrViewNotFound)`).
- Error dari router (command tidak ada, arg salah, trap) kembali sebagai
  `*routeros.DeviceError` via `rosrun` — dihandle seperti `errors.Is(err, routeros.Err?...)`.

Pergi ke: [Catalog](02-catalog.md) | [Realtime](05-realtime.md) | [Status](07-status.md).
