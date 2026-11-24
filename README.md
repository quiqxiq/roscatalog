# roscatalog

Katalog command MikroTik RouterOS API untuk tooling ISP/hotspot (gaya Mikhmon):
daftar menu, view baca (oneshot/interval/follow/listen/native), aksi tulis, dan
`Plan` stream/poll. Katalog hanya data + pembuat kata API; eksekusi ada di
adapter `rosrun`.

Repo ini adalah **satu workspace dengan dua module**:

```
go.work                     # use . dan ./goros
go.mod                      # module example.com/roscatalog
catalog_*.go                # 125 resource, 175 view baca, 395 aksi tulis
plan.go / types.go          # Plan, Build, Sentence, Query
rosrun/                     # adapter ke goros/v4 (satu-satunya import engine)
mikhmon/                    # overlay khusus Mikhmon: subset + workflow + script
cmd/example/                # jalankan satu view dari CLI
cmd/multicheck/             # uji status/stream/poll di banyak router sekaligus
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

## Dokumentasi

Wiki lengkap untuk cara pakai: [`docs/`](docs/README.md) — 12 halaman:
mulai, catalog, custom resources, views, realtime/poll, connection layer,
status realtime, fields, errors, testing, laporan migrasi MikroDash, dan
workflow Mikhmon.
Contoh yang bisa dijalankan langsung: [`docs/examples/`](docs/examples).

## Mulai

```bash
go test -race ./...
go run ./cmd/example -addr 192.168.88.1:8728 -user admin -pass secret \
  -resource hotspot-active -view live -delivery stream -every 1s
```

Untuk pemakaian aplikasi (satu koneksi async + reconnect), pakai `rosrun.Open`:

```go
sess, err := rosrun.Open(ctx, rosrun.Config{
    Address: "192.168.88.1:8728", Username: "admin", Password: "secret",
    StreamEvery: time.Second, PollEvery: 5 * time.Second, // default: stream
})
if err != nil { log.Fatal(err) }
defer sess.Close()

events, err := sess.Watch(ctx, req) // lanjut otomatis setelah reconnect
st := sess.Status()                 // connected / warning / error, dari native /ping
_ = st
```

Banyak router:

```go
mgr, err := rosrun.OpenAll(ctx, []rosrun.Config{
    {ID: "r1", Address: "192.168.88.1:8728", Username: "admin", Password: "secret",
        StreamEvery: time.Second, PollEvery: 5 * time.Second},
    {ID: "r2", Address: "192.168.88.2:8728", Username: "admin", Password: "secret",
        Poll: true, PollEvery: 5 * time.Second},
})
if err != nil { log.Printf("partial: %v", err) }
defer mgr.Close()

res, _ := ros.Lookup("interface")
events, _ := mgr.WatchAutoAll(ctx, res, "live", ros.Plan{}, ros.Params{})
for ev := range events { fmt.Println(ev.RouterID, ev.Kind) }
for _, st := range mgr.Statuses() { fmt.Println(st.ID, st.State) }
```

Default delivery adalah **stream** (realtime): `print interval=`, `follow`,
`follow-only`, atau command native seperti `monitor`. View `/listen` memakai
**`ros.Events()`** yang terpisah: hanya event perubahan, tanpa snapshot awal,
dan tidak dibuka ulang saat sepi. Set `Poll: true` untuk polling dari Go. Tiap
router memakai tepat satu koneksi async, dan command in-flight dibatasi
`MaxInFlight` (default 8). `Manager.Acquire/Release/Retain/Drop` menyediakan
refcount, hold non-viewer, dan grace.

Batas proyek: `interval` **1s..5s** (print stats, ping, command native) dan
`PollEvery` **5s..30s**. View counter mengirim `print [stats] interval=N`
dalam satu sentence, mis. `/interface/print =interval=1s =stats=`.

Uji cepat ke router:

```bash
go run ./cmd/roscheck -addr 192.168.88.1:8728 -user admin -pass secret
go run ./cmd/multicheck -router 192.168.122.60:8728 \
  -router 192.168.122.26:8728 -router 192.168.122.189:8728 \
  -user admin -pass secret -duration 20s
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

ROUTEROS_LIVE_ADDRS=192.168.122.60:8728,192.168.122.26:8728,192.168.122.189:8728 \
ROUTEROS_LIVE_USER=admin ROUTEROS_LIVE_PASS=secret \
go test ./rosrun -run TestLiveMultiRouter -v
```

`rosverify` memakai `goros.Inspect` (`/console/inspect`) untuk memeriksa path,
verb, argumen native, dan field `print`. Status `unknown` bukan error: `listen`
adalah verb API-only, dan sebagian menu singleton tidak
mengembalikan field lewat completion.

## Paket `mikhmon`

`mikhmon/` adalah overlay, bukan katalog duplikat. Ia menyaring resource root
yang dipakai Mikhmon, menyediakan builder `on-login`/Monitor Profile/Expire Monitor, format
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
_ = client.InstallProfileMonitor(ctx, "1jam", mikhmon.ExpiryRemoveRecord)

_, _ = client.AddUser(ctx, mikhmon.User{
    Name: "vc001", Password: "vc001", Profile: "1jam",
    Server: "all", Comment: mikhmon.MikhmonComment("vc001", "vc001", "batch"),
})

// Di RouterOS 7.24+, laporan yang dibuat on-login bisa memiliki owner=*sys.
// Normalisasi ini sama dengan yang dilakukan get_dashboard.php.
_, _ = client.NormalizeSalesOwner(ctx, "nov2026")
```

Verifikasi login hotspot nyata lewat webterm GNS3 pada 7.22.1, 6.49.17, dan
7.24.6: profile + on-login, voucher, sales record, dan Expire Monitor lulus.
Di 7.24.6 laporan hasil login nyata ber-`owner=*sys`; query `source` tetap
menemukannya, sedangkan dashboard `owner` perlu `NormalizeSalesOwner`.

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

Yang diuji: profile + `on-login`, voucher, login hotspot nyata, record
penjualan (`source` dan `owner`), normalisasi `owner=*sys` ala RouterOS 7.24+,
serta tick scheduler Expire Monitor mode `X`/`N`. Detail workflow ada di
`USAGE.md` §11 dan `MIKHMON.md` §14.

## Dokumen

- `USAGE.md` — cara pakai, plan stream/poll, event, jebakan, verifikasi.
- `ROADMAP.md` — rencana proplist opsional berbasis struct (raw vs proyeksi).
- `CHANGELOG.md` — riwayat rilis dan perubahan.
- `fields/` — struct field RouterOS sebagai package Go, satu file per menu (ip, ppp, queue, ...).
- `commands/` — data perintah + argumen hasil `/console/inspect`, termasuk seluruh flag `print`.
- `MIKHMON.md` — pemetaan command Mikhmon v4 dan contoh workflow.
- `REFERENCE.md` — referensi lengkap; dibuat dengan `go run gen_reference.go > REFERENCE.md` (tidak disimpan di repo).
- Dump `/console/inspect` mentah (`mikrotik_inspect/`) dipakai saat regenerasi `commands/`, tetapi tidak ikut disimpan di repo.
- `goros/docs/` — desain, keputusan, riset, dan usage engine.

## Catatan

- Katalog belum sepenuhnya terverifikasi. Jalankan `cmd/rosverify` per versi
  RouterOS dan simpan laporan JSON-nya sebagai golden schema.
- `listen` tidak mengirim `.proplist`: di RouterOS 7.24.6 proplist menekan
  event `.dead=true`, sehingga `rosrun` menyaring kolom di sisi klien dan
  penghapusan tetap terdeteksi.
- `goros` adalah fork MIT dari `go-routeros/routeros`; atribusi dipertahankan.
