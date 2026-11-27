# 12 — Workflow Mikhmon (`mikhmon/`)

Paket `mikhmon/` adalah **overlay katalog dan workflow Go untuk Mikhmon v4** di atas Mikrotik RouterOS. Paket ini bukan duplikat katalog root; sumber kebenaran path dan verb tetap katalog `roscatalog`.

Paket ini menjembatani interoperabilitas penuh dua arah antara aplikasi Go dan instalasi web PHP Mikhmon v4:

1. **RouterOS Sebagai Source of Truth**: Web UI Mikhmon v4 membaca konfigurasi profil dari header `:put` properti `on-login`, bukan dari database eksternal. Paket `mikhmon/` memastikan format ini 100% kompatibel.
2. **Monitor Profile Scheduler**: Mendukung mekanisme scheduler per-profil bawaan Mikhmon v4 (`Monitor Profile <nama_profile>`) agar indikator lampu profil di web UI Mikhmon v4 menyala hijau dan pembersihan user kedaluwarsa berjalan otomatis di router.
3. **Format Penjualan 9 Kolom**: Kompatibel penuh dengan parser `mikhmon_script_parse` dan sinkronisasi SQLite internal Mikhmon v4 (`mikhmon_sales_sync`).

---

## 1. Subset Resource Mikhmon

Mikhmon hanya memakai 22 dari sekian banyak resource RouterOS:

```go
import "github.com/quiqxiq/roscatalog/mikhmon"

// Periksa atau ambil resource dalam subset Mikhmon
res, ok := mikhmon.Lookup("hotspot-user")
ids := mikhmon.IDs()              // 22 ID resource
all := mikhmon.Resources()        // []*ros.Resource subset Mikhmon
isMikh := mikhmon.IsResource("interface") // true
```

Daftar 22 resource:
- **Dashboard & Monitoring**: `system-resource`, `system-routerboard`, `system-health`, `system-clock`, `system-identity`, `interface`, `interface-monitor-traffic`, `log`.
- **Hotspot**: `hotspot-server`, `hotspot-user`, `hotspot-user-profile`, `hotspot-active`, `hotspot-host`, `hotspot-ip-binding`, `hotspot-cookie`.
- **Infrastruktur Jaringan**: `dhcp-lease`, `ip-pool`, `fw-nat`, `queue-simple`.
- **Scripting & Laporan**: `system-scheduler`, `system-script`, `mikhmon-sales`.

---

## 2. Profile dan Script `on-login`

Mikhmon menyimpan metadata voucher (mode kedaluwarsa, harga modal, masa aktif, harga jual, status lock user) di awal string `on-login` menggunakan tuple `:put (...)`:

```routeros
:put (",remc,5000,1d,7000,,Enable,"); { ... }
```

Gunakan `mikhmon.OnLoginScript` untuk membangunnya secara aman dalam satu baris (tanpa newline):

```go
onLogin, err := mikhmon.OnLoginScript(mikhmon.OnLoginOptions{
    Mode:         mikhmon.ExpiryRemoveRecord, // "remc": hapus user saat expired + catat penjualan
    Price:        "5000",                     // harga modal
    Validity:     "1d",                       // masa aktif
    SellingPrice: "7000",                     // harga jual
    Profile:      "1jam",                     // nama profile
    LockUser:     true,                       // kunci MAC user saat login pertama kali
})
if err != nil {
    log.Fatal(err)
}

// Tambah profile ke router
_, err = client.AddProfile(ctx, mikhmon.Profile{
    Name:              "1jam",
    AddressPool:       "hs-pool",
    RateLimit:         "1M/1M",
    SharedUsers:       "1",
    StatusAutorefresh: "1m",
    OnLogin:           onLogin,
})
```

### Mode Kedaluwarsa (`ExpiryMode`)
- `ExpiryNone` (`"0"`): Tidak ada logika kedaluwarsa (mis. paket unlimited). Format header: `:put (",,<price>,,,noexp,<getlock>,")`.
- `ExpiryRemove` (`"rem"`): Menghapus user saat masa aktif habis.
- `ExpiryNotice` (`"ntf"`): Menyetel `limit-uptime=1s` saat masa aktif habis (user masih ada, sesi diputus).
- `ExpiryRemoveRecord` (`"remc"`): Sama seperti `rem` + membuat script rekaman penjualan saat login pertama.
- `ExpiryNoticeRecord` (`"ntfc"`): Sama seperti `ntf` + membuat script rekaman penjualan saat login pertama.

### Parsing Metadata yang Sudah Ada
```go
meta := mikhmon.ParseOnLoginScript(profileOnLogin)
fmt.Println(meta.Mode, meta.Price, meta.Validity, meta.Selling, meta.LockUser)
```

---

## 3. Monitor Profile Scheduler (Standar Mikhmon v4)

Di Mikhmon v4, setiap user-profile dengan masa aktif memiliki scheduler background independen bernama `<nama_profile>` dengan comment `Monitor Profile <nama_profile>` dan interval `00:02:30`.

```go
// Pasang scheduler monitor profile (indikator di web Mikhmon v4 menyala hijau)
err := client.InstallProfileMonitor(ctx, "1jam", mikhmon.ExpiryRemoveRecord)

// Cek status scheduler monitor profile
sched, err := client.ProfileMonitor(ctx, "1jam")

// Hapus scheduler monitor profile
err = client.RemoveProfileMonitor(ctx, "1jam")

// Hapus user-profile hotspot sekaligus scheduler monitor profilnya
err = client.RemoveProfile(ctx, "1jam")
```

Script background di-generate oleh `mikhmon.ProfileMonitorScript(profile, mode)` yang mengimplementasikan `mikhmon_v3_bgservice_script`. Script ini menormalisasi tanggal RouterOS v6/v7 dan memindai user berprofil tersebut untuk mengeksekusi aksi `remove` atau `set limit-uptime=1s`.

*(Catatan: Fungsi `InstallExpireMonitor()` dan `ExpireMonitorScript()` tetap tersedia untuk kebutuhan pemantau global `Mikhmon-Expire-Monitor`).*

---

## 4. Voucher & Manajemen Pengguna

```go
// Tambah voucher baru
name := "vc001"
comment := mikhmon.MikhmonComment(name, name, "batch-1") // menghasilkan "vc-batch-1"

_, err = client.AddUser(ctx, mikhmon.User{
    Name:        name,
    Password:    name,
    Profile:     "1jam",
    Server:      "all",
    Comment:     comment,
    LimitUptime: "1d",
})

// Query voucher berdasarkan comment batch
vouchers, err := client.Vouchers(ctx, comment)

// Hitung user dan sesi aktif
totalUsers, _ := client.CountUsers(ctx, "1jam")
activeSessions, _ := client.CountActive(ctx)

// Putus sesi aktif user tertentu
kicked, err := client.KickUser(ctx, name)

// Hapus user yang sudah kedaluwarsa (limit-uptime=1s)
removed, err := client.RemoveExpiredUsers(ctx)
```

---

## 5. Laporan Penjualan & Normalisasi `owner=*sys`

Script penjualan disimpan di `/system/script` dengan `comment=mikhmon` dan nama 9 kolom:
`Date-|-Time-|-User-|-Price-|-Address-|-MAC-|-Validity-|-Profile-|-Comment`

```go
// Catat penjualan manual
_, err := client.AddSale(ctx, mikhmon.Sale{
    Date:     "nov/21/2026",
    Time:     "14:00:00",
    User:     "vc001",
    Price:    "5000",
    Address:  "192.168.88.50",
    MAC:      "AA:BB:CC:DD:EE:FF",
    Validity: "1d",
    Profile:  "1jam",
    Comment:  "vc-batch-1",
})

// Query laporan penjualan bulan berjalan
sales, err := client.Sales(ctx, "nov2026")
```

### Normalisasi RouterOS 7.24+ (`owner=*sys`)
Pada RouterOS 7.24+, script yang ditambahkan oleh daemon internal hotspot otomatis diberi properti `owner=*sys`. Hal ini menyebabkan query frontend `?owner=nov2026` tidak menemukannya sampai dinormalisasi:

```go
// Normalisasi owner=*sys menjadi bulan berjalan (mis. nov2026)
n, err := client.NormalizeSalesOwner(ctx, "nov2026")
```

---

## 6. QuickPrint Presets

Preset cetak cepat voucher disimpan di `/system/script` dengan nama `Quick_Print_<package>` dan `comment=QuickPrintMikhmon`:

```go
// Simpan preset QuickPrint
id, err := client.SaveQuickPrint(ctx, mikhmon.QuickPrintPreset{
    Package:      "Paket-1Jam",
    Server:       "all",
    UserMode:     "vc",
    UserLength:   "4",
    Prefix:       "",
    Char:         "mix",
    Profile:      "1jam",
    TimeLimit:    "1d",
    DataLimit:    "0",
    Comment:      "promo",
    Validity:     "1d",
    Price:        "5000",
    SellingPrice: "7000",
    UserLock:     "Disable",
})

// Baca semua preset
presets, err := client.QuickPrints(ctx)
```

---

## 7. Typed Fields & Struct Adapters (`fields/`)

Paket `github.com/quiqxiq/roscatalog/fields` menyediakan struct khusus untuk seluruh resource Mikhmon (`fields/mikhmon.go`). Klien `mikhmon` menyediakan method typed langsung di samping metode map:

```go
import "github.com/quiqxiq/roscatalog/fields"

// Baca laporan penjualan bertipe fields.MikhmonSalesRow
sales, err := client.SalesRows(ctx, "nov2026")
for _, s := range sales {
    saleObj, _ := mikhmon.SaleFromRow(s)
    fmt.Println(saleObj.User, saleObj.Price, saleObj.Validity)
}

// Baca voucher bertipe fields.MikhmonVoucherRow (termasuk password untuk cetak)
vouchers, err := client.VouchersTyped(ctx, "vc-batch-1")
for _, v := range vouchers {
    fmt.Printf("User: %s, Pass: %s, Profile: %s\n", v.Name, v.Password, v.Profile)
}

// Baca preset QuickPrint bertipe fields.MikhmonQuickPrintRow
presets, err := client.QuickPrintRows(ctx)

// Baca status monitor scheduler bertipe fields.MikhmonSchedulerRow
pMon, err := client.ProfileMonitorRow(ctx, "1jam")
expMon, err := client.ExpireMonitorRow(ctx)
```

Adapter domain ↔ fields:
- `mikhmon.User.ToFields() fields.MikhmonUserAdd`
- `mikhmon.UserFromVoucher(fields.MikhmonVoucherRow) mikhmon.User`
- `mikhmon.Profile.ToFields() fields.MikhmonProfileAdd`
- `mikhmon.Sale.ToFields() fields.MikhmonSalesAdd`
- `mikhmon.SaleFromRow(fields.MikhmonSalesRow) (mikhmon.Sale, error)`
- `mikhmon.QuickPrintPreset.ToFields() fields.MikhmonQuickPrintAdd`
- `mikhmon.PresetFromRow(fields.MikhmonQuickPrintRow) (mikhmon.QuickPrintPreset, error)`

---

## 8. Verifikasi Langsung (`cmd/mikhmon-live`)

Untuk memverifikasi alur kerja ke router nyata:

```bash
# Uji simulasi on-login dan verifikasi scheduler
go run ./cmd/mikhmon-live -addr 192.168.88.1:8728 -user admin -pass secret

# Uji login HTTP CHAP nyata via container webterm GNS3
go run ./cmd/mikhmon-live -addr 192.168.88.1:8728 -user admin -pass secret \
  -webterm webterm-1 -hotspot http://10.10.8.1 -scheduler-wait
```

Pergi ke: [README](README.md) | [Mulai](01-getting-started.md) | [Connection Layer](06-connection-layer.md).
