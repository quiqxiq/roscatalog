# roscatalog

Katalog command MikroTik RouterOS API untuk tooling ISP/hotspot (gaya Mikhmon):
daftar menu, view baca (oneshot/interval/follow/listen/native), aksi tulis, dan
`Plan` stream/poll. Katalog hanya data + pembuat kata API; eksekusi ada di
adapter `rosrun`.

Repo ini adalah **satu workspace dengan dua module**:

```
go.work                     # use . dan ./goros
go.mod                      # module example.com/roscatalog
catalog_*.go                # 74 resource, 129 view baca, 245 aksi tulis
plan.go / types.go          # Plan, Build, Sentence, Query
rosrun/                     # adapter ke goros/v4 (satu-satunya import engine)
mikhmon/                    # overlay khusus Mikhmon: subset + workflow + script
cmd/example/                # jalankan satu view dari CLI
cmd/rosverify/              # verifikasi/discovery via /console/inspect
cmd/mikhmon-live/           # uji profile/voucher/sales/expire ke router
goros/                      # module github.com/quiqxiq/goros/v4 (fork engine)
```

Pembagian tanggung jawab:

| Lapisan | Pemilik |
|---|---|
| wire, async, tag, cancel, transport, validasi syntax/atribut, schema discovery | `goros/` |
| path/verb, view, plan, field dashboard, aksi aman, semantik domain | katalog `roscatalog` |
| eksekusi dan konversi `Request` → command | `rosrun/` |

## Mulai

```bash
go test -race ./...
go run ./cmd/example -addr 192.168.88.1:8728 -user admin -pass secret \
  -resource hotspot-active -view live -delivery stream -every 1s
```

Verifikasi katalog ke router sungguhan (RouterOS 7+):

```bash
go run ./cmd/rosverify -addr 192.168.88.1:8728 -user admin -pass secret
go run ./cmd/rosverify -addr 192.168.88.1:8728 -user admin -pass secret \
  -resource hotspot-user -all
go run ./cmd/rosverify -addr 192.168.88.1:8728 -user admin -pass secret \
  -discover ip,hotspot,active
```

Uji `listen` nyata (membuat lalu menghapus entri address-list sementara):

```bash
ROUTEROS_LIVE_ADDR=192.168.88.1:8728 ROUTEROS_LIVE_USER=admin \
ROUTEROS_LIVE_PASS=secret go test ./rosrun -run TestLiveListenRemoval -v
```

`rosverify` memakai `goros.Inspect` (`/console/inspect`) untuk memeriksa path,
verb, argumen native, dan field `print`. Status `unknown` bukan error: `listen`
adalah verb API-only, v6 tidak punya inspect, dan sebagian menu singleton tidak
mengembalikan field lewat completion.

## Paket `mikhmon`

`mikhmon/` adalah overlay, bukan katalog duplikat. Ia menyaring resource root
yang dipakai Mikhmon, menyediakan builder `on-login`/Expire Monitor, format
laporan penjualan/QuickPrint, dan `Client` workflow di atas `rosrun.Runner`.

```go
import (
    "context"

    "example.com/roscatalog/mikhmon"
)

ctx := context.Background()
client := mikhmon.New(run) // run = *rosrun.Runner

onLogin, _ := mikhmon.OnLoginScript(mikhmon.OnLoginOptions{
    Mode: mikhmon.ExpiryRemoveRecord, Price: "5000", Validity: "1d",
    SellingPrice: "7000", Profile: "1jam",
})

_, _ = client.AddProfile(ctx, mikhmon.Profile{
    Name: "1jam", RateLimit: "1M/1M", SharedUsers: "1", OnLogin: onLogin,
})
_, _ = client.AddUser(ctx, mikhmon.User{
    Name: "vc001", Password: "vc001", Profile: "1jam",
    Server: "all", Comment: mikhmon.MikhmonComment("vc001", "vc001", "batch"),
})
```

Uji langsung ke router (membuat data `zzmik...`, lalu menghapusnya):

```bash
go run ./cmd/mikhmon-live -addr 192.168.88.1:8728 -user admin -pass secret
go run ./cmd/mikhmon-live -addr 192.168.88.1:8728 -user admin -pass secret -json
```

Untuk login hotspot nyata lewat webterm GNS3:

```bash
go run ./cmd/mikhmon-live -addr 192.168.88.1:8728 -user admin -pass secret \
  -webterm <container> -hotspot http://10.10.8.1 -scheduler-wait
```

Tanpa `-webterm`, `on-login` dijalankan sebagai script sementara. Dengan
`-webterm`, command melakukan login HTTP CHAP ke hotspot, sehingga `on-login`
berjalan di konteks hotspot yang sebenarnya. `-scheduler-wait` menunggu tick
scheduler Expire Monitor, bukan menjalankan on-event manual.

Yang diuji: profile + `on-login`, voucher, record penjualan (`source` dan
`owner`), normalisasi `owner=*sys` ala RouterOS 7.24+, serta Expire Monitor
mode `X`/`N`.

## Dokumen

- `USAGE.md` — cara pakai, plan stream/poll, event, jebakan, verifikasi.
- `MIKHMON.md` — pemetaan command Mikhmon v4 dan contoh workflow.
- `REFERENCE.md` — referensi lengkap yang di-generate dari katalog.
- `goros/docs/` — desain, keputusan, riset, dan usage engine.

## Catatan

- Katalog belum sepenuhnya terverifikasi. Jalankan `cmd/rosverify` per versi
  RouterOS dan simpan laporan JSON-nya sebagai golden schema.
- `listen` tidak mengirim `.proplist`: di RouterOS 7.24.6 proplist menekan
  event `.dead=true`, sehingga `rosrun` menyaring kolom di sisi klien dan
  penghapusan tetap terdeteksi.
- `goros` adalah fork MIT dari `go-routeros/routeros`; atribusi dipertahankan.
