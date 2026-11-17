# Memakai roscatalog dengan goros/v4

Alur selalu tiga langkah: **Lookup** resource, **Build** `Request` dengan `Plan` (stream atau poll), lalu **jalankan** lewat `rosrun`.

```
ros.Lookup("hotspot-active")  ->  res.Build("live", ver, plan, params)  ->  rosrun.Watch(ctx, req)
      katalog (data)                  Request: sentence, mode, kolom             goros/v4
```

Paket `roscatalog` tidak bergantung pada library API apa pun. `rosrun` adalah satu-satunya paket yang mengimpor `github.com/quiqxiq/goros/v4`. Ganti module path `example.com/roscatalog` dengan milik Anda, lalu `go mod tidy`.

Repo ini memakai dua module dalam satu workspace: module `roscatalog` (root) dan `goros/` (fork engine). Lihat §9.

## 1. Sambungan

```go
import routeros "github.com/quiqxiq/goros/v4"

c, err := routeros.Dial("192.168.88.1:8728", "api-user", "secret") // atau DialTLS
if err != nil { log.Fatal(err) }
defer c.Close()

run := rosrun.New(c) // menyalakan mode async; satu Runner per koneksi
run.Logf = log.Printf

ctx, cancel := context.WithCancel(context.Background())
defer cancel()
```

## 2. Stream atau poll

`Plan` menentukan siapa yang menggerakkan pengiriman:

| Plan | Arti | Yang dikirim ke router |
|---|---|---|
| `ros.Stream()` | Router mengirim terus, interval bawaan view | `print =interval=`, `follow`, `/listen`, atau command native |
| `ros.StreamEvery(500*time.Millisecond)` | Sama, dengan interval sendiri (>= 200 ms untuk print) | `=interval=500ms` |
| `ros.Poll(300*time.Millisecond)` | Go mengulang one-shot tiap 300 ms (>= 50 ms) | `print` biasa, atau command + `=once=` |
| `ros.Once()` | Satu kali | satu command |
| tambahan `.WithDiff()` | Poll menghasilkan added/changed/removed | (sisi Go) |
| tambahan `.WithSeed()` | Stream listen mulai dengan snapshot awal | `listen`, lalu `print` |

`res.Deliveries("live", ver)` memberi tahu pilihan yang valid untuk sebuah view.

### Interval (router-driven) dan padanan poll-nya

```go
res, _ := ros.Lookup("hotspot-active")

// Router mengirim snapshot penuh tiap 500 ms
req, _ := res.Build("live", ros.AnyVersion, ros.StreamEvery(500*time.Millisecond), ros.Params{})
events, _ := run.Watch(ctx, req)   // KindSnapshot per siklus

// Padanan poll: Go memanggil print biasa tiap 300 ms
req, _ = res.Build("live", ros.AnyVersion, ros.Poll(300*time.Millisecond), ros.Params{})

// Poll dengan diff: hanya perubahan yang keluar
req, _ = res.Build("live", ros.AnyVersion, ros.Poll(time.Second).WithDiff(), ros.Params{})
```

### Listen dengan snapshot awal

```go
req, _ := res.Build("events", ros.AnyVersion, ros.Stream().WithSeed(), ros.Params{})
events, _ := run.Watch(ctx, req)
// KindSnapshot (keadaan awal), lalu KindAdded / KindChanged / KindRemoved
```

`Seed` memulai `listen` lebih dulu, baru `print`, supaya perubahan di antaranya tidak terlewat. Perubahan parsial dari router digabung dengan state sebelumnya, jadi `Row` selalu berisi nilai terbaru.

### Native (monitor, torch, ping)

```go
mt, _ := ros.Lookup("interface-monitor-traffic")
args := ros.Params{Args: map[string]string{"interface": "ether1,sfp1"}}

req, _ := mt.Build("run", ros.AnyVersion, ros.Stream(), args)          // stream native
req, _ = mt.Build("run", ros.AnyVersion, ros.Poll(2*time.Second), args) // monitor + =once= tiap 2 detik
// KindRow: satu baris per interface per update
```

Command tanpa `once` (ping, torch, bandwidth-test) hanya bisa stream; batasi dengan `count` atau `duration`.

### Log

```go
lg, _ := ros.Lookup("log")
req, _ := lg.Build("tail-only", ros.AnyVersion, ros.Stream(), ros.Params{})       // follow-only
req, _ = lg.Build("tail-only", ros.AnyVersion, ros.Poll(time.Second), ros.Params{}) // print + diff; hasil pertama hanya baseline
```

### Sekali jalan, jumlah, dan filter

```go
u, _ := ros.Lookup("hotspot-user")
req, _ := u.Build("count", ros.AnyVersion, ros.Once(),
    ros.Params{Query: []ros.Query{{Prop: "profile", Value: "1jam"}}})
n, err := run.Count(ctx, req)

req, _ = u.Build("list", ros.AnyVersion, ros.Once(), ros.Params{})
rows, err := run.Once(ctx, req) // kolom sudah disaring ke Fields
```

### Aksi tulis

```go
add, _ := u.Action("add")
sent, err := add.Sentence("", map[string]string{"name": "v001", "password": "x", "profile": "1jam"}, "")
id, err := run.Exec(ctx, sent) // id berisi =ret= (mis. "*A1")

set, _ := u.Action("set")
sent, _ = set.Sentence(id, map[string]string{"disabled": "yes"}, "")
_, err = run.Exec(ctx, sent)
```

### Stream dengan fallback otomatis ke poll

```go
events, err := run.WatchAuto(ctx, res, "events", ros.AnyVersion,
    ros.Plan{Every: time.Second}, ros.Params{})
```

Bila router menolak stream sebelum ada data (mis. `listen` tidak jalan di suatu menu), `WatchAuto` beralih ke poll dengan jeda yang sama dan diff, sehingga bentuk event tetap serupa.

## 3. Event

| Kind | Isi | Muncul pada |
|---|---|---|
| `KindSnapshot` | `Rows`: seluruh baris satu siklus/poll | interval, poll tabel, seed, poll diff pertama |
| `KindRow` | `Row`: satu baris | native, follow, poll native |
| `KindAdded` / `KindChanged` / `KindRemoved` | `Row` | listen, poll dengan diff |
| `KindCount` | `Count` | view count |
| `KindError` | `Err`; channel ditutup setelahnya | stream berakhir karena error |

Membatalkan `ctx` menutup channel tanpa `KindError` dan mengirim `/cancel` ke router. Konsumen harus selalu membaca channel atau membatalkan `ctx`.

## 4. Rekomendasi pemilihan

| Kebutuhan | Pilihan |
|---|---|
| Grafik traffic interface | Stream native `monitor-traffic` 1 detik; bila banyak router, Poll 2 sampai 5 detik |
| Sesi aktif hotspot/PPPoE | `events` + Seed untuk login/logout, ditambah `live` untuk counter (uptime, bytes) |
| Counter queue dan firewall | Stream interval 1 detik, atau Poll bila router lemah |
| Log | `tail-only` stream |
| Daftar voucher/pelanggan besar | Poll 30 sampai 60 detik, plus `count` untuk angka dashboard |
| ping/torch/bandwidth-test | Stream native dengan `count` atau `duration` |

## 5. Jebakan yang sudah ditangani `rosrun`

Semuanya diperiksa pada source goros/v4 (turunan go-routeros v3.0.1), dan yang pertama dibuktikan lewat uji dengan server tiruan.

- **Konsumen lambat memblokir seluruh koneksi.** Channel stream library tidak berpenyangga (`Client.Queue` bawaan 0), jadi bila Anda tidak membaca, balasan semua command lain ikut berhenti. `rosrun` menguras channel ke antrean sendiri (batas `MaxBuffered`, bawaan 65536; bila terlampaui stream dibatalkan dan `KindError` dikirim).
- **Varian `*Context` library mematikan pembaca bersama.** Saat ctx-nya dibatalkan ia memanggil `c.r.Cancel()` untuk semua command. `rosrun` memakai `RunArgs`/`ListenArgs` dan membatalkan stream lewat `ListenReply.Cancel()`.
- **Jangan kirim `.tag`.** Library menetapkannya sendiri; biarkan `Params.Tag` kosong.
- **Batas siklus `interval` tidak ditandai router.** `rosrun` memakai `.section` bila ada, kunci yang berulang, atau jeda tanpa baris. Tabel kosong tidak mengirim apa pun; set `run.EmptySnapshots = true` untuk snapshot kosong setelah 3x interval tanpa baris (ini inferensi).
- **`.proplist` ditolak.** `Request.Fallback` berisi sentence yang sama tanpa `.proplist`; `rosrun` mengulang sekali bila router menolak di awal, dan tetap menyaring kolom di sisi klien.
- **API tidak punya `where` dan regex.** `Params.Query` hanya kesamaan persis; filter prefix/regex di klien.

## 6. Banyak router dan WebSocket

- Satu `routeros.Client` + satu `Runner` per router; `Watch` berkali-kali pada koneksi yang sama aman karena mode async.
- Pola umum: satu `Watch` per feed (resource + view + Plan + Params), disebar ke banyak pelanggan WebSocket lewat hub; batalkan `ctx` feed saat pelanggan terakhir pergi.
- Untuk menentukan Plan per tenant, simpan `Plan` bersama konfigurasi feed; `Build` menolak interval/poll di bawah batas sehingga router tidak terbebani.

## 7. Validasi terhadap router Anda

Katalog baru memverifikasi 13 resource ke dokumentasi (tanda ✅ di REFERENCE.md). Sisanya jalankan sekali per versi RouterOS:

```go
req, _ := res.Build("list", ver, ros.Once(), ros.Params{})
rows, _ := run.Once(ctx, req)
if missing := ros.NeverSeen(req.Fields, rows); len(missing) > 0 {
    log.Printf("%s: properti tidak pernah muncul: %v", res.ID, missing)
}
```

Properti kosong (mis. `comment` yang belum diisi) memang tidak dikirim router, jadi hasil ini hanya indikasi.

## 8. Contoh dan tes

- `go run ./cmd/example -resource hotspot-active -view live -delivery poll -every 500ms -diff` menjalankan view mana pun dari baris perintah.
- `go test -race ./...` di root menjalankan tes katalog dan tes integrasi `rosrun` terhadap server RouterOS tiruan yang berbicara protokol kawat asli lewat klien goros/v4 sungguhan.
- `ROUTEROS_LIVE_ADDR=192.168.88.1:8728 ROUTEROS_LIVE_USER=admin ROUTEROS_LIVE_PASS=secret go test ./rosrun -run TestLiveListenRemoval -v` menguji `listen` ke router sungguhan: membuat lalu menghapus satu entri address-list dan memastikan event `removed` (`.dead=true`) diterima tanpa `.proplist`.
- `cd goros && go test ./...` menjalankan tes engine (sebagian tes lab butuh env `ROUTEROS_TEST_*`).
- `go run ./cmd/rosverify ...` memverifikasi katalog ke router sungguhan lewat `/console/inspect` (§10).
- `go run ./cmd/mikhmon-live ...` menguji workflow Mikhmon: profile + on-login, voucher, sales record, dan Expire Monitor (§11). Tanpa `-webterm`, on-login dijalankan sebagai script sementara; dengan `-webterm` + `-hotspot`, command melakukan login hotspot nyata (HTTP CHAP) dari container webterm GNS3. Tambahkan `-scheduler-wait` untuk menunggu tick scheduler Expire Monitor.

Tes mock sudah mencakup race tag go-routeros yang lama; di goros/v4 tag didaftarkan sebelum request di-flush, jadi `go test -race -count=20 ./rosrun` lulus.

## 9. Struktur repo: dua module, satu workspace

```text
go.work                 # use . dan ./goros
go.mod                  # module example.com/roscatalog
catalog_*.go            # data katalog (tanpa dependensi API)
plan.go, types.go       # Build/Plan/Sentence
rosrun/                 # adapter: satu-satunya yang mengimpor goros
mikhmon/                # overlay + workflow khusus Mikhmon
cmd/example/            # CLI menjalankan satu view
cmd/rosverify/          # verifikasi/discovery via /console/inspect
cmd/mikhmon-live/       # uji profile/voucher/sales/expire ke router
goros/                  # module github.com/quiqxiq/goros/v4 (fork engine)
  schema/               # discovery /console/inspect + cache
  gate/                 # validasi syntax (:parse) dan atribut
  transport/            # native-api, SSH, kontrak terstruktur
```

Aturan dependensi: `roscatalog` tetap murni data; `rosrun` menjadi satu-satunya jembatan ke `goros`; `mikhmon` hanya bergantung pada katalog root dan interface `Executor` (dipenuhi `*rosrun.Runner`), bukan mengimpor `goros`; `cmd/rosverify` memakai `goros.Inspect` untuk memeriksa katalog, bukan untuk menggantinya. Untuk pengembangan lokal cukup `go test ./...` di root; `go.work` otomatis memakai `goros/` setempat.

## 10. Verifikasi katalog dengan `/console/inspect` (RouterOS 7+)

`cmd/rosverify` membandingkan deklarasi katalog dengan fakta router:

```bash
go run ./cmd/rosverify -addr 192.168.88.1:8728 -user admin -pass secret
go run ./cmd/rosverify -addr 192.168.88.1:8728 -user admin -pass secret -resource hotspot-user -all
go run ./cmd/rosverify -addr 192.168.88.1:8728 -user admin -pass secret -json -strict > schema.json
```

Status yang mungkin:

| Status | Arti |
|---|---|
| `ok` | path/verb dan atribut/field yang dideklarasikan ditemukan |
| `warn` | sebagian field/flag opsional tidak ditemukan; periksa manual |
| `drift` | atribut wajib yang dideklarasikan tidak ada di router |
| `unknown` | inspect tidak bisa resolve (mis. `listen` API-only, v6, menu CHR, command root) |
| `error` | panggilan inspect gagal (jaringan/auth/timeout) |

Default hanya memeriksa atribut wajib; `-optional` ikut memeriksa properti opsional (lebih berisik karena discovery atribut `add/set` di RouterOS tidak lengkap). `-strict` exit 1 hanya bila ada `drift`/`error`.

Discovery menu/argumen sebelum menambah resource:

```bash
go run ./cmd/rosverify -addr 192.168.88.1:8728 -user admin -pass secret \
  -discover ip,hotspot,active
go run ./cmd/rosverify -addr 192.168.88.1:8728 -user admin -pass secret \
  -discover ip,hotspot -verb add -completion
```

Batasan yang perlu diingat:

- `/console/inspect` adalah **skema console**, bukan skema API. Verb API-only seperti `listen` tidak muncul; katalog tetap sumber kebenaran untuk stream.
- Menu singleton (`system/resource`, `system/health`, `ip/dns`) sering tidak mengembalikan field lewat completion; verifier menandainya `unknown`, bukan `drift`.
- RouterOS v6 tidak punya `/console/inspect`; verifier akan `unknown`.
- Inspect tidak menangkap perilaku. Contoh nyata di CHR 7.24.6: `.proplist` pada `listen` menekan event `.dead=true`. Karena itu `Build` tidak mengirim `.proplist` untuk `ModeListen`; `rosrun` menyaring kolom di klien dan penghapusan tetap terdeteksi.
- Hasil inspect bisa berbeda antar versi; simpan laporan JSON per versi sebagai golden schema di CI.

## 11. Workflow Mikhmon (`mikhmon/`)

`mikhmon/` adalah overlay khusus Mikhmon v4, bukan duplikat katalog. Paket ini
menyaring resource root yang dipakai Mikhmon dan menambahkan:

- builder `on-login` dan Expire Monitor,
- format laporan penjualan (`Sale`) dan QuickPrint,
- `Client` workflow di atas `rosrun.Runner`.

`Client` hanya membutuhkan interface `Executor`; `*rosrun.Runner` memenuhinya,
jadi `mikhmon` tidak mengimpor `goros` secara langsung.

### 11.1 Menyiapkan client

```go
import (
    "context"

    "example.com/roscatalog/mikhmon"
    "example.com/roscatalog/rosrun"
)

run := rosrun.New(c) // c = *routeros.Client yang sudah login
client := mikhmon.New(run)
ctx := context.Background()
```

Subset resource bisa dilihat lewat `mikhmon.IDs()`/`mikhmon.Resources()`;
`mikhmon.Lookup` hanya mengembalikan resource yang termasuk subset Mikhmon.

### 11.2 Profile dan `on-login`

Gunakan `OnLoginScript`; jangan menulis `on-login` multi-baris sendiri.

```go
onLogin, err := mikhmon.OnLoginScript(mikhmon.OnLoginOptions{
    Mode:         mikhmon.ExpiryRemoveRecord, // remc: expiry + laporan
    Price:        "5000",
    Validity:     "1d",
    SellingPrice: "7000",
    Profile:      "1jam",
    LockUser:     false,
    LockServer:   false,
})
if err != nil {
    log.Fatal(err)
}

_, err = client.AddProfile(ctx, mikhmon.Profile{
    Name:        "1jam",
    AddressPool: "pool-hotspot",
    RateLimit:   "1M/1M",
    SharedUsers: "1",
    OnLogin:     onLogin,
})
```

Mode yang tersedia: `ExpiryNone` (`0`), `ExpiryRemove` (`rem`),
`ExpiryNotice` (`ntf`), `ExpiryRemoveRecord` (`remc`), dan
`ExpiryNoticeRecord` (`ntfc`).

Untuk mengganti script pada profile yang sudah ada:

```go
err = client.SetOnLogin(ctx, "1jam", onLogin)
```

`AddProfile`/`SetOnLogin` menolak script yang mengandung newline, karena driver
API memecah parameter pada newline.

### 11.3 Voucher

```go
name := "vc001"
comment := mikhmon.MikhmonComment(name, name, "batch-1") // vc-batch-1

_, err = client.AddUser(ctx, mikhmon.User{
    Name:        name,
    Password:    name,
    Profile:     "1jam",
    Server:      "all",
    Comment:     comment,
    LimitUptime: "1d",
})

rows, err := client.Vouchers(ctx, comment) // view voucher, termasuk password
```

`AddUsers` menambah beberapa user berurutan. `Users(ctx, profile)` membaca user
per profile, `KickUser` memutus sesi aktif, dan `RemoveExpiredUsers` menghapus
user ber-`limit-uptime=1s` (hati-hati bila dipakai di router produksi).

### 11.4 Expire Monitor

```go
if err := client.InstallExpireMonitor(ctx); err != nil {
    log.Fatal(err)
}

sched, err := client.ExpireMonitor(ctx) // map nama/interval/on-event/disabled
```

`InstallExpireMonitor` membuat scheduler `Mikhmon-Expire-Monitor` bila belum
ada, atau memperbarui interval/on-event lalu mengaktifkannya. On-event-nya
disalin dari `assets/js/func.js` dan tetap satu baris.

Mode kerja:

- `N`: saat kedaluwarsa, set `limit-uptime=1s` lalu putuskan sesi aktif.
- `X`: saat kedaluwarsa, hapus user lalu putuskan sesi aktif.

### 11.5 Laporan penjualan dan `owner=*sys`

`on-login` mode `remc`/`ntfc` menulis `/system/script` dengan:

- `name`: `tanggal-|-jam-|-user-|-harga-|-address-|-mac-|-validity-|-profile-|-comment`,
- `source`: tanggal legacy, mis. `nov/18/2026`,
- `owner`: bulan+tahun, mis. `nov2026`,
- `comment`: `mikhmon`.

Menulis laporan manual:

```go
_, err := client.AddSale(ctx, mikhmon.Sale{
    Date: "nov/18/2026", Time: "10:00:00", User: "vc001",
    Price: "5000", Address: "198.51.100.10", MAC: "AA:BB:CC:DD:EE:FF",
    Validity: "1d", Profile: "1jam", Comment: "vc-123-11.18.26",
})

rows, err := client.Sales(ctx, "nov2026")
```

Di RouterOS 7.24+, login hotspot nyata dapat menulis laporan dengan
`owner=*sys`. Akibatnya:

- query `?source=nov/18/2026` (dipakai `get_report.php`) tetap menemukan,
- query `?owner=nov2026` (dipakai dashboard) belum menemukan.

Normalisasi yang sama dengan `get_dashboard.php`:

```go
n, err := client.NormalizeSalesOwner(ctx, "nov2026")
// cari comment=mikhmon + owner=*sys, lalu set owner=nov2026
```

Jika `owner` kosong, helper menurunkan owner dari `source` tiap laporan.
Setelah normalisasi, `client.Sales(ctx, "nov2026")` menemukan laporan tersebut.

### 11.6 QuickPrint

```go
id, err := client.SaveQuickPrint(ctx, mikhmon.QuickPrintPreset{
    Package: "1jam", Server: "all", UserMode: "vc", UserLength: "4",
    Char: "mix", Profile: "1jam", TimeLimit: "1d", DataLimit: "0",
    Comment: "catatan", Validity: "1d",
    Price: "5000", SellingPrice: "7000", UserLock: "Disable",
})
```

`SaveQuickPrint` melakukan upsert berdasarkan `Quick_Print_<paket>` dan menulis
`comment=QuickPrintMikhmon`. `QuickPrints(ctx)` membaca semua preset beserta
`source`-nya.

### 11.7 Uji langsung ke router

Simulasi (tanpa login hotspot nyata):

```bash
go run ./cmd/mikhmon-live -addr 192.168.88.1:8728 -user admin -pass secret
go run ./cmd/mikhmon-live -addr 192.168.88.1:8728 -user admin -pass secret -json
```

Login hotspot nyata dari webterm GNS3 + tick scheduler sebenarnya:

```bash
go run ./cmd/mikhmon-live -addr 192.168.88.1:8728 -user admin -pass secret \
  -webterm <container> -hotspot http://10.10.8.1 -scheduler-wait -json
```

Container webterm harus berada di subnet hotspot yang sama. Command akan:

- membuat profile + voucher uji berprefix `zzmik`,
- login HTTP CHAP ke hotspot (mode `-webterm`),
- memeriksa sesi aktif dan comment expiry,
- membuat laporan penjualan dan mengecek query `source`/`owner`,
- menguji normalisasi `owner=*sys`,
- menunggu tick Expire Monitor (`-scheduler-wait`) atau menjalankan on-event
  manual,
- menghapus user/profile/script/sesi uji (pakai `-keep` untuk menahan data).

Hasil verifikasi 2026-11-18 dengan login hotspot nyata:

| Router | Versi | Profile/voucher | Real login | Sales owner | Expire X/N |
|---|---|---|---|---|---|
| 192.168.122.26 | 7.22.1 | ok | ok | `nov2026` | ok |
| 192.168.122.60 | 6.49.17 | ok | ok | `nov2026` | ok |
| 192.168.122.189 | 7.24.6 | ok | ok | `*sys` -> normalisasi | ok |

Di 7.24.6, `source` tetap ditemukan; dashboard `owner` perlu
`NormalizeSalesOwner`. Versi 7.24.2 tidak tersedia saat pengujian; perilaku
7.24.x diverifikasi di 7.24.6.

Detail pemetaan command dan contoh mentah ada di `MIKHMON.md`; API paket ada di
`mikhmon/README.md`.
