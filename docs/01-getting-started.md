# 01 — Mulai

Install:

```bash
go get github.com/quiqxiq/roscatalog@latest
```

Pembuatan session — pakai `rosrun.Open` (satu koneksi asinkron + auto-reconnect):

```go
package main

import (
	"context"
	"fmt"
	"os"

	ros "github.com/quiqxiq/roscatalog"
	"github.com/quiqxiq/roscatalog/rosrun"
)

func main() {
	ctx := context.Background()
	sess, err := rosrun.Open(ctx, rosrun.Config{
		Address:  os.Getenv("ROS_ADDR"),
		Username: os.Getenv("ROS_USER"),
		Password: os.Getenv("ROS_PASS"),
		Status:   rosrun.StatusConfig{Disable: true}, // API health saja
	})
	if err != nil {
		panic(err)
	}
	defer sess.Close()

	res, ok := ros.Lookup("interface")
	if !ok {
		panic("interface tidak ada")
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
}
```

Persyaratan: Go 1.22+, RouterOS 7.9+, API user dengan policy `read` (atau `write` untuk aksi), port `8728` (API) atau `8729` (API-SSL / TLS).

Catatan penting:

- `res.Build(view, plan, params)` → `(ros.Request, error)`.
- Session (`*rosrun.MikroTik` / `*rosrun.Session`) = tepat satu koneksi asinkron per router.
- Default `StatusConfig` tanpa `Target` hanya memantau koneksi API; isi `Status.Target` (mis. `"8.8.8.8"`) bila ingin native `/ping` aktif.
- Lihat [Custom Resources](03-custom-resources.md) dan [Errors](09-errors.md).

Pergi ke: [Catalog](02-catalog.md) | [Realtime](05-realtime.md) | [Status](07-status.md).
