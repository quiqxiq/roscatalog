# Roadmap: proplist opsional berbasis struct

> Status: **usulan**, belum diimplementasikan. Dokumen ini memuat keputusan desain
> dan tahapan kerja. Contoh API di bawah bersifat indikatif dan boleh dirapikan
> saat implementasi.

## 1. Tujuan

Proyek ini adalah **library** yang harus:

1. bisa di-import dan dipakai oleh aplikasi lain;
2. mudah digunakan — pemakai tidak perlu hafal nama kolom RouterOS sebagai teks;
3. bisa disesuaikan — pemakai bisa memilih kolom yang mau diambil;
4. **difokuskan ke RouterOS v7.9+** (karena `/console/inspect` baru ada di
   v7.9). Untuk pemakai dengan flota v6, jalan keluar tanpa `Version` pertama-
   class: definisikan resource/view sendiri atau pakai `Raw`/`Params.Proplist`;
5. tetap berisi **command yang biasa dan penting saja**, bukan seluruh command MikroTik.

Perubahan yang direncanakan mendukung kelima poin itu tanpa mengubah tujuan
katalog: katalog tetap berisi data + pembuat kalimat API, eksekusi tetap di
`rosrun`.

## 2. Keadaan sekarang (masalah yang mau diselesaikan)

Alur proplist saat ini:

```
catalog_*.go   List("address", "network", "interface", ...)   ← hardcode []string
     ↓
View.Proplist + Request.Fields
     ↓
Sentence()     "=.proplist=address,network,interface,..."
     ↓
rosrun.trim()  row["address"], row["network"], ... → map[string]string
```

Konsekuensinya:

- daftar kolom hanya berupa string, tidak ada tipe Go yang bisa di-import;
- pemakai tidak bisa meminta subset kolom lewat API publik (`Params` hanya
  punya `Tag`, `Query`, `Interval`, `Args`);
- hasil baca tetap `map[string]string`, jadi nama kolom harus ditebak manual;
- tidak ada jalan eksplisit untuk mengembalikan **semua kolom apa adanya**.
  Yang ada hanya `Request.Fallback` (proplist dibuang kalau router menolak),
  itu mekanisme retry, bukan pilihan desain.

Catatan: `rosrun.trim()` sebenarnya **sudah** menyalin semua kolom ketika
`Request.Fields` kosong. Jadi mode "raw" tidak perlu mengubah `trim()` — cukup
ada sakelar di `Build`/`Sentence` supaya `=.proplist=` tidak dikirim dan
`Fields` dibiarkan kosong.

## 3. Konsep: dua mode hasil

| Mode | `=.proplist=` dikirim? | `Request.Fields` | Hasil dari `rosrun` | Dipakai untuk |
|---|---|---|---|---|
| **Raw** | tidak | kosong | semua kolom apa adanya (`map[string]string`) | eksplorasi, debug, field yang belum ada di struct, `rosverify` |
| **Proyeksi** | ya, berasal dari struct | diisi dari struct | kolom terpilih; bisa di-decode ke struct | pemakaian aplikasi sehari-hari |

Aturan:

- `Raw` adalah **opsi**, bukan pengganti. Mode proyeksi tetap menjadi default
  agar perilaku library tidak berubah mendadak.
- `Raw` menang atas opsi proyeksi lain (mis. `Params.Raw = true` membuat
  `Only`/struct diabaikan untuk kalimat dan untuk `Fields`).
- Untuk `listen`, proplist memang tidak pernah dikirim (alasan: di RouterOS
  7.24.6 `=.proplist=` menekan event `.dead=true`). Mode proyeksi tetap
  menyaring kolom di klien lewat `Fields` yang diambil dari view `list`.
  Mode `Raw` pada `listen` mengembalikan semua atribut event, termasuk
  `.dead=true`; konsumen bertugas menangannya sendiri.

## 4. Struct sebagai sumber proplist

### 4.1 Struct baca (untuk `print` / view baca)

Satu struct per resource menggantikan string kolom yang sekarang ada di
`List(...)`, `Live(...)`, `Tail(...)`, `Monitor(...)`.

Contoh:

```go
type HotspotActiveRow struct {
    ID      string `ros:".id"`
    User    string `ros:"user"`
    Address string `ros:"address"`
    Uptime  string `ros:"uptime"`
    Comment string `ros:"comment,omitempty"`
}
```

Aturan tag yang diusulkan:

- `ros:"nama"` → nama properti RouterOS. Kalau kosong, pakai nama field Go
  yang diubah ke kebab-case (`RateLimit` → `rate-limit`).
- `,omitempty` →kolom tidak dimasukkan kalau nilainya kosong (terutama untuk
  struct tulis).
- `,key` → menandai field kunci baris; default `.id`.
- `,v6` / `,v7` → **pengecualian**, bukan tag wajib. Default (tanpa tag)
  adalah v7.9+ — baseline yang dipakai `mikrotik_inspect/`. Field tanpa tag
  umumnya juga jalan di v6 karena command v7 relatif stabil; tag versi hanya
  dipakai bila memang diketahui beda (mis. `ip-route`: `routing-mark` di v6
  vs `routing-table` di v7). Proplist dan `Fields` tetap dipilih sesuai
  `Version` yang diberikan ke `Build`.

Semua field bertipe `string` dulu. RouterOS API mengirim teks, jadi ini
lossless dan menghindari kegagalan parsing untuk nilai seperti `1d2h`,
`1M/1M`, atau `yes/no`. Konversi ke `bool`/`int`/durasi bisa menjadi helper
terpisah di kemudian hari.

### 4.2 Struct tulis (untuk `add`, `set`, `remove`)

Baca dan tulis sengaja dibedakan karena bentuk kebutuhannya beda:

- `print` → daftar kolom yang ingin dilihat (termasuk field read-only/dinamis
  seperti `dynamic`, `invalid`, `active-address`);
- `add`/`set` → field yang boleh dikirim, dengan mana yang wajib dan mana yang
  opsional;
- `remove` → hanya butuh selector (`.id` / `numbers`), tanpa struct field.

Contoh:

```go
type IPAddressAdd struct {
    Address   string `ros:"address,required"`
    Interface string `ros:"interface,required"`
    Network   string `ros:"network,omitempty"`
    Comment   string `ros:"comment,omitempty"`
}

type IPAddressSet struct {
    Network string `ros:"network,omitempty"`
    Comment string `ros:"comment,omitempty"`
}
```

Field kosong dengan `omitempty` tidak dikirim — mengikuti pola yang sudah
dipakai `mikhmon.User.args()` / `mikhmon.Profile.args()`.

`Action.Unknown()` tetap dipertahankan untuk logging: menu firewall punya
ratusan properti, jadi struct hanya memuat yang penting dan pemakai boleh
tetap mengirim field di luar struct lewat `map[string]string` bila perlu.

### 4.3 Sumber referensi perintah: `mikrotik_inspect/`

Direktori `mikrotik_inspect/` sudah berisi hasil `/console/inspect` dari
RouterOS, terpisah per menu. `_index.json` adalah indeksnya (path, scope,
jumlah menu/command/arg per file); `_root.json` dan `*_cmds.json` berisi
command level root; file seperti `ip_address.json` berisi command per menu.
Total menurut indeks: **554 menu, 5134 command, 35657 arg**. Bentuk nodenya:

```json
"add": {
  "_type": "cmd",
  "address": {"_type": "arg", "desc": "A.B.C.D    (IP address)"},
  ...
},
"print": {
  "_type": "cmd",
  "count-only": {"_type": "arg"},
  "follow": {"_type": "arg"},
  "interval": {"_type": "arg", "desc": "00:00:00.200..    (time interval)"},
  "proplist": {"_type": "arg"},
  ...
}
```

Data ini dipakai sebagai **sumber kebenaran untuk perintah, verb, dan
argumen** saat membuat atau mengoreksi katalog dan struct:

- menentukan `add`/`set`/`remove`/`enable`/`disable`/`move`/... yang benar
  untuk sebuah menu beserta argumennya → dasar struct tulis dan
  `Action.Required` / `Action.Optional`;
- mengecek flag `print` yang benar-benar tersedia (`count-only`, `follow`,
  `follow-only`, `interval`, `proplist`, ...) → dasar `Deliveries()`;
- memberi `desc` untuk dokumentasi dan validasi nilai (IP address, time
  interval, script, dll);
- menjadi daftar calon resource baru yang belum masuk katalog (mis.
  `interface/wireguard`, `routing/ospf`, `user/group`) tanpa menebak path.

**Batasan penting:** file inspect berisi **command + argumen**, bukan
**kolom hasil `print`**. Nama kolom baris (properti seperti `address`,
`network`, `interface`, `dynamic`, `active-address`) tetap diambil dari
katalog yang ada dan/atau diverifikasi ke router sungguhan lewat `rosverify`
+ `ros.NeverSeen()`. Jadi:

- struct **tulis** → sebagian besar bisa dibantu langsung oleh `add`/`set`
  di inspect;
- struct **baca** → inspect hanya membantu memvalidasi flag `print`; daftar
  kolomnya tetap dari katalog + verifikasi live.

Catatan lain: `listen` adalah verb API-only dan tidak tercatat sebagai
command di inspect (pencarian hanya menemukan `listen` sebagai argumen di
`routing_bgp.json`), jadi ketiadaannya bukan berarti view `listen` tidak
valid. Karena `/console/inspect` baru ada di RouterOS **7.9+**, dump ini
berasal dari v7.9 ke atas. Di v7 command relatif stabil, jadi versi minor
tidak terlalu berpengaruh. Baseline "v7.9+" tetap perlu dicatat di metadata
sumber supaya tidak ada yang mengira data ini berlaku untuk v6.

**Hasil pengambilan data (sudah dilakukan):** `mikrotik_inspect/SOURCE.md`
berisi struct Go untuk semua menu, dibuat oleh `tools/gen_source_md.py` dari
tiga sumber: file inspect, hasil `print` + completion `print,proplist` dari
router CHR GNS3 `192.168.122.189` (`tools/probe_print`), dan dump katalog
(`tools/dump_catalog`). Data router mentah disimpan di
`mikrotik_inspect/print_fields.json`, dump katalog di
`mikrotik_inspect/catalog_dump.json`. Saat ini: 217 menu, 181 struct `Row`,
68 `Add`, 90 `Set`, semuanya sudah divalidasi dengan `gofmt`.

Struct di `SOURCE.md` sudah diwujudkan sebagai kode Go di folder `fields/`
(package `fields`, satu file per menu besar: `ip.go`, `ppp.go`, `queue.go`,
`interface.go`, `system.go`, dst — 25 file) oleh `tools/gen_fields_go.py`.
Pemakaiannya: `fields.IPAddressRow`. `go build ./...` dan `go test ./...`
lolos. Ini baru tahap "struct bisa di-import"; pemakaiannya sebagai proplist
dan mode raw masih menunggu perubahan builder/Plan.

Catatan: `_index.json` menyebut 243 file, tetapi hanya 99 file yang ikut
ter-attach. Untuk path tanpa file namun punya data `print`, struct `Row`
tetap dibuat dari data router; untuk resource katalog yang pathnya lebih
dalam dari file inspect, field diambil dari katalog.


## 5. Penyesuaian oleh pemakai

Tiga level, dari yang paling sederhana:

1. **Default dari katalog** — struct bawaan resource dipakai sebagai proplist.
2. **Subset kolom** — `Params.Only` (atau `Except`) membatasi field dari struct
   tanpa mengubah katalog.
3. **Struct sendiri** — `BuildOf[T]` memakai struct milik pemakai sebagai
   proplist dan sebagai cetakan hasil. Ini level lanjutan; pemakai bertanggung
   jawab atas nama properti (bantu verifikasi dengan `rosverify` /
   `ros.NeverSeen()`).

Kombinasi yang tidak boleh terjadi: `Raw: true` bersama `Only`/`BuildOf[T]`
— `Raw` menang dan sisanya diabaikan, atau `Build` mengembalikan error supaya
tidak ambigu. Perlu diputuskan saat implementasi.

## 6. Contoh pemakaian yang ditargetkan

```go
import (
    "context"

    ros "example.com/roscatalog"
    "example.com/roscatalog/rosrun"
)

res, _ := ros.Lookup("hotspot-active")
ctx := context.Background()

// 1) RAW — tanpa .proplist=, semua kolom dikembalikan apa adanya.
reqRaw, _ := res.Build("list", ros.AnyVersion, ros.Once(), ros.Params{Raw: true})
rawRows, _ := run.Once(ctx, reqRaw) // []map[string]string, semua kolom

// 2) PROYEKSI — proplist berasal dari struct, bukan string hardcode.
//    Only membatasi kolom yang diminta tanpa mengubah katalog.
req, _ := res.Build("list", ros.AnyVersion, ros.Once(),
    ros.Params{Only: []string{"user", "address", "uptime"}})
rows, _ := run.Once(ctx, req)
active := ros.DecodeAll[ros.HotspotActiveRow](rows)

// 3) TULIS — struct sebagai argumen; field kosong tidak dikirim.
add := ros.IPAddressAdd{Address: "10.0.0.2/32", Interface: "ether1"}
sentence, _ := res.Action("add").SentenceOf("", add, "")
ret, _ := run.Exec(ctx, sentence) // =ret=*ID
```

`DecodeAll[T]` mengisi field dari `map[string]string` memakai tag yang sama
dengan pembuat proplist. Field yang tidak hadir tetap kosong (RouterOS tidak
mengirim properti yang bernilai kosong, mis. `comment` yang belum diisi).
Sedyakan juga `DecodeStrict[T]` bila pemakai ingin error pada field yang hilang.

## 7. Dampak ke kode yang ada

| File | Perubahan |
|---|---|
| `fields/` (baru, package `fields`) | struct per menu: `ip.go`, `ppp.go`, `queue.go`, `interface.go`, ... (sudah dibuat) |
| `fields.go` (baru, package `roscatalog`) | helper generik: ambil nama field dari tag, `Decode[T]`, `DecodeAll[T]`, `ArgsOf[T]`, validasi `Only` |
| `fields_test.go` (baru) | unit test tag, omitempty, versi, decode |
| `tools/probe_print/` (baru) | probe kolom `print` + completion dari router → `print_fields.json` |
| `tools/dump_catalog/` (baru) | dump resource katalog → `catalog_dump.json` |
| `tools/gen_source_md.py` (baru) | gabungkan tiga sumber → `mikrotik_inspect/SOURCE.md` |
| `tools/gen_fields_go.py` (baru) | `SOURCE.md` → `fields/*.go` (sudah dibuat) |
| `builder.go` | tambah `ListOf[T]`, `LiveOf[T]`, `TailOf[T]`, `MonitorOf[T]` di samping `List`/`Live`/... yang lama |
| `types.go` | `Params` tambah `Raw`, `Only`/`Except`; `Action.SentenceOf` |
| `plan.go` | `Build` menaati `Raw`/`Only`; tambah `BuildOf[T]`; `Fingerprint` ikut mempertimbangkan opsi baru |
| `rosrun/runner.go` | **kemungkinan besar tidak berubah** — `trim()` sudah mempertahankan semua kolom kalau `Fields` kosong. Cukup tambah test raw mode |
| `catalog_*.go` | migrasi bertahap dari `List("a", "b")` ke `ListOf[T]` |
| `mikhmon/client.go` | setelah struct root tersedia, `User`/`Profile` dibungkus adapter ke struct root (tidak langsung dihapus) |
| `gen_reference.go` | REFERENCE.md menampilkan kolom dari struct, bukan hanya string |
| `README.md`, `USAGE.md` | dokumentasi mode raw/proyeksi dan contoh struct |

## 8. Jebakan yang harus dijaga

- **`listen` + `.dead=true`** — jangan pernah mengirim `=.proplist=` pada
  `/listen`; pertahankan alasan dan catatan yang sudah ada di `builder.go`
  dan `plan.go`.
- **`count-only`** — tidak menghasilkan baris; struct tidak relevan.
- **Native/monitor** — sebagian command menerima `=.proplist=`, sebagian tidak.
  Mekanisme `Request.Fallback` (retry tanpa proplist) tetap harus bekerja,
  termasuk saat proplist berasal dari struct.
- **Versi** — baseline adalah v7.9+. Jangan memberi tag `,v7` pada semua field
  hanya karena sumbernya inspect. Tag versi hanya untuk pengecualian yang
  memang diketahui beda di v6, dan tetap dalam satu resource yang sama.
- **Sumber inspect** — `mikrotik_inspect/` berisi command + argumen, bukan
  kolom hasil `print`. Jangan menjadikannya satu-satunya sumber proplist;
  `listen` juga tidak tercatat sebagai command di sana (API-only).
- **Lingkup** — struct hanya untuk command penting yang sudah dikurasi.
  Jangan membuka pintu "semua command MikroTik" hanya karena ada struct.
- **Kompatibilitas** — API lama berbasis string tetap didukung sampai migrasi
  selesai; aplikasi yang sudah memakai `List("...")` tidak boleh rusak.

## 9. Tahapan kerja

| Fase | Hasil | Kriteria selesai |
|---|---|---|
| 0. Dokumen | `ROADMAP.md` ini | keputusan desain disetujui |
| 1. Fondasi tipe | `fields/*.go` (sudah), `fields.go` helper + test | struct per menu bisa di-import sebagai `fields.IPAddressRow`; tag terbaca, decode struct berhasil, `omitempty` bekerja; output deterministik |
| 2. Builder & Plan | method `...Of[T]`, `Params.Raw`/`Only`, `Action.SentenceOf`, `BuildOf[T]` | mode raw benar-benar mengirim semua kolom; proyeksi memakai struct; test lama tetap lulus |
| 3. Migrasi katalog | struct untuk resource inti: `hotspot-active`, `hotspot-user`, `dhcp-lease`, `interface`, `system-resource`, `ip-address`, `ip-route`; calon resource baru dicek ke `mikrotik_inspect/` | katalog tidak lagi mengandung string kolom untuk resource tersebut; verb/argumen cocok dengan inspect |
| 4. Integrasi | contoh di `cmd/example`, adapter `mikhmon`, test `rosrun` untuk raw vs proyeksi | pemakai bisa import struct dan langsung memakainya |
| 5. Rapikan | `gen_reference.go`, `README.md`, `USAGE.md`, `rosverify` | dokumen dan katalog tetap satu sumber kebenaran |

Resource baru yang diambil dari `mikrotik_inspect/` tetap melewati kurasi
(yang penting saja) dan verifikasi `rosverify` sebelum masuk katalog — inspect
hanya menjawab "command/argumen apa yang ada", bukan "kolom apa yang perlu
ditampilkan".

**Progres Fase 1 (command/action dulu):** sudah ditambahkan `ArgKind`
(`ArgValue`/`ArgFlag`), `Action.Read`, `Action.Args`, `Action.SentenceWith`,
`View.AllowedFlags`, `Params.Flags`, builder `ReadDo`/`PrintFlags`, validasi
flag beserta konfliknya (`brief`/`detail`/`terse`, `follow`/`follow-only`,
`stats`/`count-only`), package data `commands/`, serta `zz_catalog_enrich.go`
yang mengisi flag `print` dan menambah aksi `print`, `comment`, `unset`,
`reset`, `export` untuk 34 resource. `go test ./...` lolos.

**Progres Fase 2 (mesin print & aksi):** `Params.Raw` (tanpa `=.proplist=`,
`Request.Fields` dikosongkan supaya semua kolom dikembalikan) dan `Params.Only`
(subset field, divalidasi) sudah bekerja di `Build`. Builder mendapat preset
`Print`, `Brief`, `Detail`, `Stats`, `StatsDetail`. `Action.SentenceWith`
dipakai untuk semua verb termasuk `print`. Test raw/only/preset ditambahkan;
`go test ./...` lolos. Struct masih ditunda sampai Fase 3.

## 10. Generator struct dengan Python

Keputusan: **boleh memakai Python**, dengan batasan tegas:

- Python hanya **tooling offline untuk maintainer**, bukan dependency runtime.
  Library tetap murni Go: `go build` / `go test ./...` tidak membutuhkan Python.
- File hasil generate (`fields_gen.go`) **di-commit**. Konsumen cukup
  `go get`, tidak perlu Python terpasang.
- Pakai **Python stdlib saja** (`re`, `pathlib`, `json`) supaya tidak ada
  `pip install` atau `requirements.txt`.
- Output harus **deterministik**: urutan resource mengikuti `GroupOrder`
  lalu ID, urutan field mengikuti urutan di katalog. Diff harus bersih.
- Header file hasil generate:
  `// Code generated by tools/gen_fields.py; DO NOT EDIT.`
- **Tidak** memakai directive `//go:generate` — itu membuat `go generate ./...`
  gagal pada kontributor tanpa Python. Perintah regenerasi ditulis di
  `tools/README.md`.

Dua mode generator:

1. **Mode migrasi (dipakai dulu).** Parser `catalog_*.go` membaca
   `Table(...)` / `Single(...)`, `List` / `Live` / `Tail` / `Monitor`, dan
   `CRUD(req, opt)`, lalu menghasilkan struct + tag. Ini langsung mengubah
   string hardcode menjadi struct tanpa membuat format spec baru, dan
   semantiknya persis mengikuti katalog yang sudah ada.
2. **Mode referensi inspect (dipakai untuk resource baru).** Generator bisa
   membaca `mikrotik_inspect/*.json` yang sudah ada (lihat §4.3) untuk
   menghasilkan struct tulis, `Action.Required`/`Optional`, dan daftar flag
   `print` yang valid. Tidak perlu membuat format spec baru. Kolom struct
   baca tetap berasal dari katalog + hasil `rosverify`.

Kenapa Python cocok: memetakan `rate-limit` → `RateLimit` dengan acronym map
(IP, DNS, DHCP, MAC, VLAN, RTT, CPU, HTTP, NAT, ARP, TTL, ...) dan menulis tag
`ros:"..."` secara eksplisit, plus validasi field duplikat, field read-only
(ada di `List` tapi tidak di `Add`), pengecualian v6/v7 (bukan default), serta
laporan field yang gagal dipetakan. Nama Go yang kurang rapi hanya berdampak kosmetik karena
nama kabel tetap berasal dari tag.

Lokasi: `tools/gen_fields.py` + `tools/README.md`; output `fields_gen.go` di
root (package `roscatalog`). Generator membaca dua sumber: `catalog_*.go`
(migrasi kolom baca) dan `mikrotik_inspect/*.json` (verb/argumen untuk struct
tulis dan resource baru). Laporan hasil generate dibandingkan dengan
`rosverify` terhadap router untuk menandai field yang belum terverifikasi.

## 11. Keputusan

**Sudah disepakati di diskusi:**

- proplist tidak lagi hardcode string; **struct** yang jadi sumber.
- **proplist opsional** — harus bisa mengembalikan data mentah tanpa
  `=.proplist=`.
- proplist **bisa disesuaikan** (subset + struct sendiri).
- struct **baca** dan **tulis** dibedakan.
- field bertipe `string` dulu.
- generator struct memakai **Python**, tetapi hanya sebagai tooling offline;
  `fields_gen.go` di-commit dan build Go tetap murni (lihat §10).
- `mikrotik_inspect/` menjadi referensi untuk command/verb/argumen yang akan
  dibuatkan katalog dan struct (lihat §4.3).
- kompatibilitas dengan API lama dijaga.
- fokus ke RouterOS **v7.9+** (karena `/console/inspect` baru ada di sana);
  v6 tetap didukung lewat `Version` tetapi bukan target utama; command v7
  cenderung stabil sehingga tidak perlu tag versi massal.
- tetap hanya command penting.

**Masih terbuka (diputuskan saat implementasi):**

- nama tag final: `ros:` atau `routeros:`.
- `Params.Only` menerima nama field Go atau nama properti RouterOS (usulan:
  nama properti, karena itu yang dikirim ke router).
- ke mana metadata baseline "v7.9+" ditaruh: field di `_index.json` atau
  file `mikrotik_inspect/SOURCE.md` (cukup satu, yang penting tercatat).
- `mikhmon.User` / `mikhmon.Profile` diganti atau dibungkus adapter dulu.
- konflik `Raw` + `Only`/`BuildOf[T]`: abaikan atau error.

---

## 12. Status implementasi (penanganan koneksi)

Sudah dikerjakan:

- **Satu koneksi async default**; `Config.DualConn=true` adalah escape hatch
  yang akan dihapus setelah satu rilis.
- **Default realtime**: `Config.Poll` (default false) → stream; `Stream`
  deprecated. `PlanFor` memilih stream lebih dulu.
- **`DeliveryEvents` terpisah dari `DeliveryStream`**: `ros.Events()` khusus
  `/listen`; tanpa watchdog diam (sepi = normal), `Stream()` pada view listen
  dinormalkan ke Events demi kompatibilitas.
- **Limiter per router** (`MaxInFlight`, default 8) + statistik command/menu.
- **Satu scheduler per sesi** untuk semua `DeliveryPoll`.
- **Union proplist per menu**: subscriber berbeda field berbagi satu upstream;
  tiap subscriber menerima hasil yang disaring ke field-nya sendiri.
- **Lifecycle**: `Manager.Acquire/Release/Retain/Drop` + grace.
- **Reconnect**: auth backoff, stream watchdog 10s/5s, level channel di-reset
  saat koneksi putus.
- **Helper struct opsional**: `FieldNames`, `Decode`, `DecodeAll`, `ArgsOf`,
  `Params.Proplist`.

Validasi:

- `go test -race ./rosrun` lolos.
- `go vet ./...` lolos.
- Uji read-only ke CHR GNS3 `192.168.122.189:8728`: connect, stream+exec,
  poll, union proplist, dan listen — semuanya lolos.
