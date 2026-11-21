# Paket `mikhmon`

Overlay katalog dan workflow Go untuk **Mikhmon v4** di atas Mikrotik RouterOS.

Paket ini bukan duplikat katalog; sumber kebenaran path dan verb tetap katalog root `roscatalog`. Paket ini menyediakan:

1. **Subset Resource Mikhmon**: Menyaring 22 resource RouterOS yang benar-benar dipakai oleh Mikhmon (dashboard, hotspot user & profile, active sessions, hosts, bindings, cookies, DHCP leases, pool, simple queues, scheduler, scripts).
2. **On-Login Script Builder (`OnLoginScript`)**: Membangun properti `on-login` satu baris persis seperti `mikhmon_v3_onlogin_script` di Mikhmon v4, termasuk header metadata `:put (",<mode>,<price>,<validity>,<sprice>,,<getlock>,")`, normalisasi tanggal RouterOS v6/v7, scheduler sementara untuk kalkulasi tanggal expired, lock MAC, dan script pencatatan penjualan 9 kolom.
3. **Monitor Profile Scheduler (`ProfileMonitorScript`)**: Generator on-event scheduler `Monitor Profile <nama_profile>` persis seperti `mikhmon_v3_bgservice_script` di Mikhmon v4. Schedulernya dibuat per-profil sehingga indikator lampu di web UI Mikhmon v4 menyala hijau (aktif) dan pembersihan voucher kedaluwarsa berjalan otomatis di router.
4. **Global Expire Monitor (`ExpireMonitorScript`)**: On-event scheduler `Mikhmon-Expire-Monitor` global (~1,9 KB) untuk memindai seluruh user di router.
5. **Format Penjualan & QuickPrint**: Representasi struct `Sale` (9 field `-|-`) dan `QuickPrintPreset` (13 field `#`).
6. **Workflow Client (`Client`)**: Membungkus aksi add profile, add user/voucher, install monitor profile, add sale, simpan QuickPrint, normalisasi laporan `owner=*sys` (RouterOS 7.24+), kick user, dan hapus user kedaluwarsa (`limit-uptime=1s`).

`Client` hanya bergantung pada interface `Executor`:

```go
type Executor interface {
    Once(ctx context.Context, req ros.Request) ([]map[string]string, error)
    Count(ctx context.Context, req ros.Request) (int, error)
    Exec(ctx context.Context, sentence []string) (string, error)
}
```

Interface ini dipenuhi oleh `*rosrun.Runner`, sehingga paket ini tidak mengimpor `goros` secara langsung.

---

## Ringkasan API

### Resource Subset
- `Lookup(id string) (*ros.Resource, bool)`: Lookup resource terbatas pada subset Mikhmon.
- `Resources() []*ros.Resource`: Daftar seluruh resource dalam subset Mikhmon.
- `IDs() []string`: Daftar ID resource Mikhmon (22 resource).
- `IsResource(id string) bool`: Memeriksa apakah suatu resource masuk dalam subset Mikhmon.

### On-Login & Metadata Script
- `OnLoginScript(o OnLoginOptions) (string, error)`: Membangun script `on-login` satu baris sesuai format Mikhmon v4.
- `ParseOnLoginScript(script string) OnLoginMeta`: Membaca kembali metadata mode, harga modal, masa aktif, harga jual, dan lock dari script `on-login` yang ada.
- `OneLine(s string) string`: Merapatkan whitespace menjadi spasi tunggal agar aman dikirim via API RouterOS tanpa newline.

### Background Schedulers (Expiry)
- `ProfileMonitorScript(profile string, mode ExpiryMode) (string, error)`: Generator script background per-profile Mikhmon v4 (`mikhmon_v3_bgservice_script`). Mode `rem`/`remc` menggunakan `remove`, sedangkan `ntf`/`ntfc` menggunakan `set limit-uptime=1s`.
- `ExpireMonitorScript() string`: Generator script global `Mikhmon-Expire-Monitor`.

### Model Data
- `Sale`: Format laporan penjualan 9 kolom (`Date-|-Time-|-User-|-Price-|-Address-|-MAC-|-Validity-|-Profile-|-Comment`).
  - `Name() string`: Nama script di `/system/script`.
  - `Source() string`: Tanggal format legacy `mmm/dd/yyyy`.
  - `Owner() string`: Bulan dan tahun, misal `nov2026`.
  - `ParseSaleName(name string) (Sale, error)`: Mengurai kembali nama script penjualan.
- `QuickPrintPreset`: Format preset cetak cepat voucher 13 kolom `#`.
  - `Name() string`: `Quick_Print_<package>` di `/system/script`.
  - `Source() string`: String 13 kolom diawali `#`.
  - `ParseQuickPrintSource(source string) (QuickPrintPreset, error)`: Mengurai source preset.

### Client Workflow
- `New(run Executor) *Client`: Membuat instance Client baru.
- **User & Voucher**:
  - `AddUser(ctx, u User) (string, error)`: Tambah 1 user/voucher hotspot.
  - `AddUsers(ctx, users []User) ([]string, error)`: Tambah batch user/voucher berurutan.
  - `Vouchers(ctx, comment string) ([]map[string]string, error)`: Baca voucher berdasarkan comment batch.
  - `Users(ctx, profile string) ([]map[string]string, error)`: Daftar user (opsional filter profile).
  - `CountUsers(ctx, profile string) (int, error)`: Hitung total user.
  - `CountActive(ctx) (int, error)`: Hitung sesi hotspot aktif.
  - `KickUser(ctx, user string) (int, error)`: Putus semua sesi aktif milik user.
  - `RemoveExpiredUsers(ctx) (int, error)`: Hapus user ber-`limit-uptime=1s`.
- **Profile & Monitor**:
  - `AddProfile(ctx, p Profile) (string, error)`: Tambah hotspot user profile.
  - `SetOnLogin(ctx, profile, script string) error`: Ganti properti `on-login` profile yang sudah ada.
  - `InstallProfileMonitor(ctx, profile string, mode ExpiryMode) error`: Pasang/update scheduler `Monitor Profile <profile>` (interval `00:02:30`). Jika mode `0`, scheduler dihapus.
  - `RemoveProfileMonitor(ctx, profile string) error`: Hapus scheduler monitor profile.
  - `ProfileMonitor(ctx, profile string) (map[string]string, error)`: Baca status scheduler monitor profile.
  - `RemoveProfile(ctx, profile string) error`: Hapus profile hotspot sekaligus scheduler monitor profilnya.
  - `InstallExpireMonitor(ctx) error`: Pasang/update scheduler global `Mikhmon-Expire-Monitor`.
  - `ExpireMonitor(ctx) (map[string]string, error)`: Baca scheduler global `Mikhmon-Expire-Monitor`.
- **Laporan & QuickPrint**:
  - `AddSale(ctx, s Sale) (string, error)`: Tambah laporan penjualan ke `/system/script`.
  - `Sales(ctx, owner string) ([]map[string]string, error)`: Ambil laporan berdasarkan bulan (mis. `nov2026`).
  - `NormalizeSalesOwner(ctx, owner string) (int, error)`: Normalisasi laporan `owner=*sys` (RouterOS 7.24+) ke format bulan+tahun agar tampil di dashboard.
  - `SaveQuickPrint(ctx, p QuickPrintPreset) (string, error)`: Buat atau perbarui preset QuickPrint.
  - `QuickPrints(ctx) ([]map[string]string, error)`: Baca semua preset QuickPrint.

---

## Contoh Penggunaan Lengkap

```go
package main

import (
    "context"
    "log"

    routeros "github.com/quiqxiq/goros/v4"
    "example.com/roscatalog/mikhmon"
    "example.com/roscatalog/rosrun"
)

func main() {
    ctx := context.Background()

    // 1. Koneksi ke router
    c, err := routeros.Dial("192.168.88.1:8728", "admin", "secret")
    if err != nil {
        log.Fatal(err)
    }
    defer c.Close()

    run := rosrun.New(c)
    client := mikhmon.New(run)

    // 2. Buat script on-login dengan metadata harga Mikhmon v4
    onLogin, err := mikhmon.OnLoginScript(mikhmon.OnLoginOptions{
        Mode:         mikhmon.ExpiryRemoveRecord, // "remc": hapus user + catat penjualan
        Price:        "5000",                     // harga modal
        Validity:     "1d",                       // masa aktif 1 hari
        SellingPrice: "7000",                     // harga jual
        Profile:      "1jam",                     // nama profile
        LockUser:     true,                       // kunci MAC saat pertama kali login
    })
    if err != nil {
        log.Fatal(err)
    }

    // 3. Tambah Hotspot User Profile
    _, err = client.AddProfile(ctx, mikhmon.Profile{
        Name:              "1jam",
        AddressPool:       "hs-pool",
        RateLimit:         "1M/1M",
        SharedUsers:       "1",
        StatusAutorefresh: "1m",
        OnLogin:           onLogin,
    })
    if err != nil {
        log.Fatal(err)
    }

    // 4. Pasang Monitor Profile scheduler (agar indikator di web Mikhmon v4 hijau)
    if err := client.InstallProfileMonitor(ctx, "1jam", mikhmon.ExpiryRemoveRecord); err != nil {
        log.Fatal(err)
    }

    // 5. Buat voucher hotspot
    name := "vc001"
    comment := mikhmon.MikhmonComment(name, name, "batch-1") // format "vc-batch-1"
    _, err = client.AddUser(ctx, mikhmon.User{
        Name:        name,
        Password:    name,
        Profile:     "1jam",
        Server:      "all",
        Comment:     comment,
        LimitUptime: "1d",
    })
    if err != nil {
        log.Fatal(err)
    }

    // 6. Normalisasi laporan jika berjalan di RouterOS 7.24+
    // (di mana hotspot daemon internal otomatis memberi owner=*sys)
    if n, err := client.NormalizeSalesOwner(ctx, "nov2026"); err == nil && n > 0 {
        log.Printf("%d laporan dinormalisasi ke nov2026", n)
    }

    // 7. Ambil laporan penjualan bulan berjalan
    sales, err := client.Sales(ctx, "nov2026")
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Total penjualan tercatat: %d transaksi", len(sales))
}
```

---

## Kompatibilitas dengan Web UI Mikhmon v4

- **Source of Truth di RouterOS**: Web UI Mikhmon v4 membaca langsung konfigurasi profile dari string `:put` pada properti `on-login`. Format yang dihasilkan paket ini 100% kompatibel saat dibuka di menu Edit Profile, Add User, Quick User, Generate User, dan Cetak Voucher Mikhmon v4.
- **Indikator Monitor Profile**: Di Mikhmon v4, halaman daftar profil memeriksa `/system/scheduler/print ?name=<profile>`. Dengan memanggil `client.InstallProfileMonitor(ctx, profile, mode)`, status profil di web Mikhmon v4 akan terdeteksi aktif (lampu hijau).
- **Penjualan & SQLite Sync**: Script penjualan yang dicatat saat login menggunakan konvensi `-|-` 9 kolom standar. Fitur sinkronisasi database internal Mikhmon v4 (`mikhmon_sales_sync` / `sales_records`) dapat langsung mengimpor data tanpa kendala.

---

## Uji Langsung ke Router

Command `cmd/mikhmon-live` menyediakan verifikasi langsung terhadap alur kerja Mikhmon ke router nyata (RouterOS v6 maupun v7):

```bash
# Uji simulasi on-login dan verifikasi scheduler
go run ./cmd/mikhmon-live -addr 192.168.88.1:8728 -user admin -pass secret
go run ./cmd/mikhmon-live -addr 192.168.88.1:8728 -user admin -pass secret -json
```

Untuk login hotspot nyata melalui container webterm GNS3:

```bash
go run ./cmd/mikhmon-live -addr 192.168.88.1:8728 -user admin -pass secret \
  -webterm <container> -hotspot http://10.10.8.1 -scheduler-wait
```

Command ini menguji:
1. Pembuatan profile dengan tuple `:put` format v4.
2. Pemasangan scheduler `Monitor Profile <nama_profile>` dan verifikasi interval.
3. Pembuatan voucher dengan prefix comment `vc-`/`up-`.
4. Eksekusi `on-login` (simulasi atau CHAP nyata) dan verifikasi comment user berubah menjadi tanggal format legacy `mmm/dd/yyyy hh:mm:ss`.
5. Pengecekan script laporan penjualan 9 kolom dan normalisasi `owner=*sys`.
6. Pemasangan dan pembersihan scheduler Expire Monitor.
7. Seluruh data uji berprefix `zzmik` dibersihkan otomatis setelah pengujian selesai.
