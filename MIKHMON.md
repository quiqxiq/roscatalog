# Mikhmon v4 → roscatalog

Dokumen ini memetakan command RouterOS yang dipakai **Mikhmon v4** (PHP) ke
katalog **roscatalog** (Go). Fokusnya bukan UI Mikhmon, melainkan lapisan API:
command apa yang dikirim, kenapa, dan bagaimana bentuk padanannya di
`roscatalog` + `rosrun`.

Referensi sumber (checkout lokal, tidak disalin ke repo ini):

| Area | File Mikhmon |
|---|---|
| Koneksi & dashboard | `get/get_dashboard.php`, `config/connection.php` |
| User/voucher | `post/post_add_user.php`, `post/post_generate_voucher.php`, `get/get_users.php` |
| Profile & on-login | `post/post_add_userprofile.php`, `post/post_update_userprofile.php`, `core/mikhmon_scripts.php` |
| Expire monitor | `assets/js/func.js` (`setExpMon`), `post/post_expire_monitor.php` |
| Laporan penjualan | script `$record` di `core/mikhmon_scripts.php` + `get/get_report.php`, `get/get_dashboard.php` |
| QuickPrint | `post/post_quick_print.php` |
| Active/host/binding/cookie/DHCP | `post/post_hotspot_remove.php`, `post/post_ip_binding.php`, `post/post_dhcp_lease.php` |
| Status voucher | `status/status.php` |

---

## 1. Peta cepat

| Kebutuhan Mikhmon | Command RouterOS | roscatalog |
|---|---|---|
| Resource/CPU/RAM/uptime | `/system/resource/print` | `system-resource` → view `list`, `live` |
| Model board | `/system/routerboard/print` | `system-routerboard` → view `list` |
| Health/suhu | `/system/health/print` | `system-health` → view `list`, `live` |
| Jam & timezone | `/system/clock/print` | `system-clock` → view `list` |
| Identitas | `/system/identity/print` | `system-identity` → view `list` |
| Jumlah user/aktif | `/ip/hotspot/user/print count-only`, `/ip/hotspot/active/print count-only` | `hotspot-user`/`hotspot-active` → view `count` |
| Traffic interface | `/interface/monitor-traffic` | `interface-monitor-traffic` → view `run` |
| Log hotspot | `/system/logging/*`, `/log/print` | `log` → `list`/`tail`/`tail-only`; `/system/logging` belum ada resource |
| Daftar user/voucher | `/ip/hotspot/user/print [?profile]` | `hotspot-user` → view `list`, `voucher` |
| Tambah user/voucher | `/ip/hotspot/user/add` | `hotspot-user` → action `add` |
| Ubah user | `/ip/hotspot/user/set` | action `set` |
| Reset counter | `/ip/hotspot/user/reset-counters` | action `reset-counters` |
| Hapus user | `/ip/hotspot/user/remove` | action `remove` |
| Disable/enable user | `/ip/hotspot/user/set disabled=...` | action `enable`/`disable` |
| Profile user (harga, validity, on-login) | `/ip/hotspot/user/profile/*` | `hotspot-user-profile` → `list`/`detail`, action `add`/`set`/`remove` |
| Sesi aktif | `/ip/hotspot/active/print` | `hotspot-active` → `list`, `live`, `events` |
| Kick sesi | `/ip/hotspot/active/remove` | action `remove` |
| Host hotspot | `/ip/hotspot/host/print` | `hotspot-host` → view `list`, action `remove`/`make-binding` |
| IP binding | `/ip/hotspot/ip-binding/*` | `hotspot-ip-binding` |
| Cookie | `/ip/hotspot/cookie/*` | `hotspot-cookie` |
| DHCP lease | `/ip/dhcp-server/lease/*` | `dhcp-lease` |
| Scheduler (Expire Monitor) | `/system/scheduler/*` | `system-scheduler` |
| QuickPrint | `/system/script/*` | `system-script` |
| Laporan penjualan | `/system/script/*` | `mikhmon-sales` (path sama: `/system/script`) |
| Baca NAT/queue/pool | `/ip/firewall/nat/print`, `/queue/simple/print`, `/ip/pool/print` | `fw-nat`, `queue-simple`, `ip-pool` |

Poin penting: `roscatalog` **tidak** menyalin data Mikhmon. Ia menyediakan
path/verb, field, dan aksi yang sama; Mikhmon tetap aplikasi PHP terpisah.

---

## 2. Pola dasar di Go

```go
import (
    "context"
    "fmt"
    "log"

    routeros "github.com/quiqxiq/goros/v4"
    ros "example.com/roscatalog"
    "example.com/roscatalog/rosrun"
)

c, err := routeros.Dial("192.168.88.1:8728", "admin", "secret")
if err != nil {
    log.Fatal(err)
}
defer c.Close()

run := rosrun.New(c)
ctx := context.Background()
```

Membaca:

```go
res, _ := ros.Lookup("hotspot-user")
req, _ := res.Build("list", ros.AnyVersion, ros.Once(),
    ros.Params{Query: []ros.Query{{Prop: "profile", Value: "1jam"}}})
rows, err := run.Once(ctx, req)
```

Menulis:

```go
add, _ := res.Action("add")
sent, err := add.Sentence("", map[string]string{
    "name": "vc001", "password": "vc001", "profile": "1jam",
}, "")
id, err := run.Exec(ctx, sent) // =ret= dari /ip/hotspot/user/add
```

`Action.Sentence(selector, args, tag)` mengisi selector (`.id`/`numbers`) dan
memvalidasi argumen wajib. `run.Exec` mengembalikan `=ret=` (mis. `.id` baru).

---

## 3. Dashboard dan monitoring

Padanan `get_dashboard.php`:

```go
// Resource + versi + uptime
res, _ := ros.Lookup("system-resource")
req, _ := res.Build("list", ros.AnyVersion, ros.Once(), ros.Params{})
rows, _ := run.Once(ctx, req)
fmt.Println(rows[0]["version"], rows[0]["uptime"], rows[0]["cpu-load"])

// Health: view-nya ber-versi, jadi pilih V6 atau V7.
health, _ := ros.Lookup("system-health")
req, _ = health.Build("list", ros.V7, ros.Once(), ros.Params{})
rows, _ = run.Once(ctx, req)

// Identitas
id, _ := ros.Lookup("system-identity")
req, _ = id.Build("list", ros.AnyVersion, ros.Once(), ros.Params{})
rows, _ = run.Once(ctx, req)
fmt.Println(rows[0]["name"])
```

Jumlah user dan sesi aktif:

```go
u, _ := ros.Lookup("hotspot-user")
req, _ := u.Build("count", ros.AnyVersion, ros.Once(), ros.Params{})
n, _ := run.Count(ctx, req)

a, _ := ros.Lookup("hotspot-active")
req, _ = a.Build("count", ros.AnyVersion, ros.Once(), ros.Params{})
active, _ := run.Count(ctx, req)
```

Mikhmon mengurangi jumlah user dengan `-1` di dashboard; itu quirk UI, jangan
disalin ke katalog.

Traffic interface (native, satu baris per interface):

```go
tr, _ := ros.Lookup("interface-monitor-traffic")
req, _ := tr.Build("run", ros.AnyVersion, ros.Once(), ros.Params{
    Args: map[string]string{"interface": "ether1"},
})
rows, _ := run.Once(ctx, req) // rx-bits-per-second, tx-bits-per-second, ...
```

Streaming terus-menerus:

```go
req, _ = tr.Build("run", ros.AnyVersion, ros.Stream(), ros.Params{
    Args: map[string]string{"interface": "ether1"},
})
events, _ := run.Watch(ctx, req)
for ev := range events {
    fmt.Println(ev.Row["rx-bits-per-second"], ev.Row["tx-bits-per-second"])
}
```

Log: Mikhmon menambahkan rule `/system/logging` dengan prefix `->`, lalu
membaca `/log/print ?topics=hotspot,info,debug`. Katalog belum punya resource
`system-logging`, jadi bagian itu pakai sentence mentah:

```go
_, err := run.Exec(ctx, []string{
    "/system/logging/add",
    "=action=memory", "=prefix=->", "=topics=hotspot,info,debug",
})

lg, _ := ros.Lookup("log")
req, _ := lg.Build("list", ros.AnyVersion, ros.Once(), ros.Params{
    Query: []ros.Query{{Prop: "topics", Value: "hotspot,info,debug"}},
})
rows, _ := run.Once(ctx, req)
```

---

## 4. User dan voucher

Mikhmon (`post/post_add_user.php`) mengirim:

```php
$API->comm("/ip/hotspot/user/add", [
  "server" => $server, "name" => $name, "password" => $password,
  "profile" => $profile, "disabled" => "no", "comment" => $comment,
  // opsional, hanya bila tidak kosong:
  "mac-address" => ..., "limit-uptime" => ..., "limit-bytes-total" => ...,
]);
```

Comment diberi prefix:

- `vc-...` bila `name == password` (voucher)
- `up-...` bila berbeda (user biasa)

Padanan Go:

```go
u, _ := ros.Lookup("hotspot-user")
add, _ := u.Action("add")
sent, err := add.Sentence("", map[string]string{
    "server":       "all",
    "name":         "vc001",
    "password":     "vc001",
    "profile":      "1jam",
    "comment":      "vc-123-11.18.26",
    "limit-uptime": "1h",          // opsional
    // "limit-bytes-total": "1073741824", // 1 GiB, opsional
}, "")
if err != nil {
    log.Fatal(err)
}
id, err := run.Exec(ctx, sent)
```

Mikhmon memvalidasi agar properti opsional yang kosong tidak dikirim (RouterOS
menolak `limit-uptime=` kosong). Di Go, cukup jangan masukkan key-nya.

### Batch voucher

Mikhmon (`post_generate_voucher.php`) membuat nama/password acak, comment
`<mode>-<gencode>-<mm.dd.yy>-<catatan>` (mode `vc`/`up`), lalu loop `add`.
Versi sederhana:

```go
func addVouchers(ctx context.Context, run *rosrun.Runner, qty int,
    prefix, profile, server, comment string) error {

    u, _ := ros.Lookup("hotspot-user")
    add, _ := u.Action("add")
    for i := 1; i <= qty; i++ {
        name := fmt.Sprintf("%s%04d", prefix, i)
        sent, err := add.Sentence("", map[string]string{
            "name": name, "password": name, "profile": profile,
            "server": server, "comment": comment,
        }, "")
        if err != nil {
            return err
        }
        if _, err := run.Exec(ctx, sent); err != nil {
            return err
        }
    }
    return nil
}
```

Mikhmon memakai `randLC/randUC/randULC/randNLC` dan password terpisah; logika
generator itu milik aplikasi, bukan katalog. Yang penting bentuk command-nya
sama.

### Cetak voucher

Mikhmon mengambil ulang berdasarkan comment batch, lalu memakai view
`voucher` (berisi password):

```go
req, _ := u.Build("voucher", ros.AnyVersion, ros.Once(), ros.Params{
    Query: []ros.Query{{Prop: "comment", Value: "vc-123-11.18.26"}},
})
rows, _ := run.Once(ctx, req)
for _, row := range rows {
    fmt.Println(row["name"], row["password"], row["limit-uptime"])
}
```

Catatan: query API hanya cocok persis. Mikhmon kadang mencari prefix comment
(`vc-`, `up-`) di sisi klien; di roscatalog lakukan filter yang sama setelah
`run.Once`.

### Ubah, reset, hapus

```go
set, _ := u.Action("set")
sent, _ := set.Sentence(id, map[string]string{"profile": "2jam"}, "")
_, err = run.Exec(ctx, sent)

reset, _ := u.Action("reset-counters")
sent, _ = reset.Sentence(id, nil, "")
_, err = run.Exec(ctx, sent)

remove, _ := u.Action("remove")
sent, _ = remove.Sentence(id, nil, "")
_, err = run.Exec(ctx, sent)
```

Aksi `remove` memakai selector `.id`; `reset-counters` memakai `=numbers=`, dan
nilai `.id` (`*A1`) juga diterima di sana. Jangan tertukar dengan `=number=`
yang dipakai `/system/script/run` (nama script).

Mikhmon juga memakai `disabled=yes/no` untuk toggle. Katalog menyediakan aksi
`enable`/`disable`. Di menu `hotspot-user` dan `system-scheduler`, kedua verb
itu terverifikasi `ok`; di menu seperti `hotspot-user-profile` dan
`system-script`, `/console/inspect` tidak bisa me-resolve-nya sehingga
verifier menandainya `unknown`, tetapi command-nya tetap valid di API.

---

## 5. Profile dan `on-login`

### Apa yang dikirim Mikhmon

`post/post_add_userprofile.php`:

```php
$API->comm("/ip/hotspot/user/profile/add", [
  "name" => $name,
  "address-pool" => $addrpool,
  "rate-limit" => $ratelimit,
  "shared-users" => $sharedusers,
  "status-autorefresh" => "1m",
  "on-login" => $onlogin,
  "parent-queue" => $parent,
]);
```

Field-field itu semuanya sudah ada di katalog:

```go
prof, _ := ros.Lookup("hotspot-user-profile")
add, _ := prof.Action("add")
sent, err := add.Sentence("", map[string]string{
    "name":              "1jam",
    "address-pool":      "pool-hotspot",
    "rate-limit":        "1M/1M",
    "shared-users":      "1",
    "status-autorefresh": "1m",
    "parent-queue":      "none",
    "on-login":          onLogin,
}, "")
id, err := run.Exec(ctx, sent)
```

### Anatomi script `on-login` Mikhmon

Sumber: `core/mikhmon_scripts.php`. Script ini **bukan** system script; ia
disimpan di properti `on-login` profile user dan dijalankan RouterOS setiap
kali user login. Mikhmon membangunnya sebagai **satu baris** karena driver API
memecah parameter pada newline.

Bagian-bagiannya:

1. **Penanda harga/validity**

   ```routeros
   :put (",rem,5000,1d,7000,,Disable,Disable,")
   ```

   `print_voucher.php` membelah string `on-login` dengan koma: index 2 = harga
   modal, index 3 = validity, index 4 = harga jual. Jadi properti ini juga
   dipakai sebagai metadata voucher. Dua field setelah harga jual adalah
   `lockuser` dan `lockserver`.

2. **Normalisasi tanggal v6/v7**

   RouterOS 6 memberi `nov/18/2026`, RouterOS 7 memberi `2026-11-18`. Script
   menormalkan keduanya ke bentuk legacy `mmm/dd/yyyy` agar laporan Mikhmon
   konsisten.

3. **Hitung expiry sekali saja**

   ```routeros
   :local mode "X";
   :local comment [/ip hotspot user get [find where name=$user] comment];
   :local ucode [:pick $comment 0 2];
   :if ($ucode = "vc" or $ucode = "up" or $comment = "") do={
     /system scheduler add name=$user start-date=$date interval=1d;
     :delay 2s;
     :local exp [/system scheduler get [find where name=$user] next-run];
     /ip hotspot user set comment="$exp $mode" [find where name=$user];
     /system scheduler remove [find where name=$user];
   }
   ```

   Triknya: scheduler sementara dipakai untuk **menghitung** `next-run`
   (tanggal kedaluwarsa), lalu langsung dihapus. Comment user menjadi sumber
   kebenaran expiry. Login berikutnya tidak menghitung ulang karena comment
   sudah berisi tanggal — ini mencegah masa aktif bertambah setiap login.

4. **Mode expiry**

   - `N` (`ntf`/`ntfc`): saat kedaluwarsa, set `limit-uptime=1s` lalu putuskan
     sesi aktif. User tetap ada.
   - `X` (`rem`/`remc`): hapus user dan putuskan sesi aktif.

5. **Lock MAC/server** (opsional)

   ```routeros
   :local mac $"mac-address";
   /ip hotspot user set mac-address=$mac [find where name=$user];
   :local srv [/ip hotspot host get [find where mac-address="$mac"] server];
   /ip hotspot user set server=$srv [find where name=$user]
   ```

6. **Rekaman penjualan** (opsional, mode `remc`/`ntfc`)

   ```routeros
   :local mac $"mac-address";
   :local time [/system clock get time];
   :local sname "$legacy-|-$time-|-$user-|-5000-|-$address-|-$mac-|-1d-|-1jam-|-$comment";
   /system script add name=$sname owner="$monthname$year" source=$legacy comment=mikhmon;
   /system script set [find name=$sname] owner="$monthname$year"
   ```

   Baris di atas ditampilkan terlipat; di router semuanya menjadi satu baris
   yang dipisah `;`. Lihat §7.

Hotspot menyediakan `$user`, `$address`, dan `$mac-address` saat `on-login`.
Script Mikhmon menambahkan variabel lokal `$comment` (dibaca dari comment
user), `$legacy`, `$monthname`, `$year`, `$mode`, `$validity`, dan seterusnya.
Jadi `$comment` di script itu **bukan** variabel bawaan hotspot.

### Contoh `on-login` yang ringkas

Untuk memahami alurnya, ini versi sederhana (bukan pengganti script Mikhmon
yang menangani v6/v7 dan laporan):

```go
const onLogin = `:put (",rem,5000,1d,7000,,,,"); ` +
    `:local mode "X"; ` +
    `:local date [/system clock get date]; ` +
    `:local time [/system clock get time]; ` +
    `:local comment [/ip hotspot user get [find where name=$user] comment]; ` +
    `:if ([:pick $comment 0 2] = "vc" or [:pick $comment 0 2] = "up" or $comment = "") do={ ` +
    `/system scheduler add name=$user start-date=$date interval=1d; ` +
    `:delay 2s; ` +
    `:local exp [/system scheduler get [find where name=$user] next-run]; ` +
    `/ip hotspot user set comment="$exp $mode" [find where name=$user]; ` +
    `/system scheduler remove [find where name=$user] ` +
    `}`
```

Pasang ke profile:

```go
sent, err := add.Sentence("", map[string]string{
    "name":              "1jam",
    "rate-limit":        "1M/1M",
    "shared-users":      "1",
    "status-autorefresh": "1m",
    "address-pool":      "pool-hotspot",
    "on-login":          onLogin,
}, "")
```

Penting:

- Jangan pernah mengirim `on-login` multi-baris lewat API. Ganti newline
  dengan `;` atau spasi, seperti Mikhmon (`preg_replace('/\s+/', ' ', ...)`).
- Untuk script lengkap, gunakan hasil `mikhmon_onlogin_script()` dari
  `core/mikhmon_scripts.php`; jangan tulis ulang dari nol.
- Setelah dipasang, verifikasi dengan view `detail`:

  ```go
  req, _ := prof.Build("detail", ros.AnyVersion, ros.Once(), ros.Params{
      Query: []ros.Query{{Prop: "name", Value: "1jam"}},
  })
  rows, _ := run.Once(ctx, req)
  fmt.Println(rows[0]["on-login"])
  ```

---

## 6. Expire Monitor

Mikhmon memasang scheduler bernama `Mikhmon-Expire-Monitor` yang berjalan
setiap menit (`post/post_expire_monitor.php`, script dibangun di
`assets/js/func.js` `setExpMon`):

```php
$API->comm("/system/scheduler/add", [
  "name" => "Mikhmon-Expire-Monitor",
  "start-time" => "00:00:00",
  "interval" => "00:01:00",
  "on-event" => $expire_monitor_src,
  "disabled" => "no",
  "comment" => "Mikhmon Expire Monitor",
]);
```

Alur script `on-event`:

1. Normalisasi tanggal v6/v7.
2. Iterasi user yang comment-nya mengandung tahun berjalan/sebelumnya.
3. Parse tanggal + jam kedaluwarsa dari comment.
4. Jika sudah lewat dan `limit-uptime != 00:00:01`:
   - mode `N`: `limit-uptime=1s` + hapus sesi aktif.
   - mode `X`: hapus user + hapus sesi aktif.

Pseudocode yang disederhanakan:

```routeros
:foreach i in=[/ip hotspot user find where comment~"/2026" || comment~"/2025"] do={
  :local c [/ip hotspot user get $i comment];
  :local name [/ip hotspot user get $i name];
  :if (<tanggal-kedaluwarsa> < <sekarang>) do={
    :if ([:pick $c 21] = "N") do={
      /ip hotspot user set limit-uptime=1s $i;
      /ip hotspot active remove [find where user=$name];
    } else={
      /ip hotspot user remove $i;
      /ip hotspot active remove [find where user=$name];
    }
  }
}
```

Script aslinya ~1,9 KB dan memakai fungsi `normdate`, `dateint`, dan `timeint`;
salin dari `assets/js/func.js` bila perlu dipasang ulang.

Memasang lewat katalog:

```go
sched, _ := ros.Lookup("system-scheduler")
add, _ := sched.Action("add")
sent, err := add.Sentence("", map[string]string{
    "name":       "Mikhmon-Expire-Monitor",
    "start-time": "00:00:00",
    "interval":   "00:01:00",
    "on-event":   expireMonitorScript,
    "comment":    "Mikhmon Expire Monitor",
}, "")
_, err = run.Exec(ctx, sent)
```

Mikhmon juga mengirim `disabled=no`; nilainya default, jadi contoh di atas
tidak menyertakannya. `Action.Sentence` tetap meneruskan argumen tambahan bila
Anda ingin menyalin command Mikhmon persis.

Membaca/mengubah:

```go
req, _ := sched.Build("list", ros.AnyVersion, ros.Once(), ros.Params{
    Query: []ros.Query{{Prop: "name", Value: "Mikhmon-Expire-Monitor"}},
})
rows, _ := run.Once(ctx, req)

// Matikan lewat aksi disable (selector .id).
off, _ := sched.Action("disable")
sent, _ = off.Sentence(rows[0][".id"], nil, "")
_, err = run.Exec(ctx, sent)
```

Catatan: Mikhmon memakai `/system/scheduler/set disabled=...`; katalog
menyediakan `enable`/`disable`. Keduanya mengubah properti yang sama.

---

## 7. Laporan penjualan (`$record`)

Saat user login dan profile memakai mode `remc`/`ntfc`, `on-login` menulis
system script:

```routeros
:local mac  $"mac-address";
:local time [/system clock get time];
:local sname "$legacy-|-$time-|-$user-|-5000-|-$address-|-$mac-|-1d-|-1jam-|-$comment";
/system script add name=$sname owner="$monthname$year" source=$legacy comment=mikhmon;
/system script set [find name=$sname] owner="$monthname$year";
```

Konvensi Mikhmon:

| Properti | Isi |
|---|---|
| `name` | `tanggal-|-jam-|-user-|-harga-|-address-|-mac-|-validity-|-profile-|-comment` |
| `source` | tanggal legacy `mmm/dd/yyyy` |
| `owner` | bulan+tahun, mis. `nov2026` |
| `comment` | `mikhmon` |

`get_report.php` membacanya dengan `?source=<hari>`. Laporan live di
`get_dashboard.php` memakai `?owner=<bulan><tahun>`; RouterOS 7.24+ kadang
menulis `owner=*sys`, dan Mikhmon memperbaikinya dengan `set owner=<bulan>`.

Katalog punya resource khusus untuk konvensi ini: `mikhmon-sales`, dengan path
yang sama (`/system/script`), view `report`, dan aksi `add`/`remove`.

Padanan Go:

```go
sales, _ := ros.Lookup("mikhmon-sales")
add, _ := sales.Action("add")
sent, err := add.Sentence("", map[string]string{
    "name":    "nov/18/2026-10:00:00-vc001-|-5000-|-198.51.100.10-|-AA:BB:CC:DD:EE:FF-|-1d-|-1jam-|-vc-batch01",
    "source":  "nov/18/2026",
    "owner":   "nov2026",
    "comment": "mikhmon",
}, "")
_, err = run.Exec(ctx, sent)
```

Membaca laporan bulanan (view `report` karena butuh `source`):

```go
req, _ := sales.Build("report", ros.AnyVersion, ros.Once(), ros.Params{
    Query: []ros.Query{{Prop: "owner", Value: "nov2026"}},
})
rows, _ := run.Once(ctx, req)
for _, row := range rows {
    fmt.Println(row["name"], row["source"])
}
```

Filter `mikhmon-sales`: `owner`, `comment`, `name`; query `source` juga bisa
dikirim lewat `Params.Query` karena API menerima properti apa pun. Ingat, query
API hanya kesamaan persis — bukan regex seperti `?comment~"..."` di dalam
script RouterOS. Jika ingin mengelola semua `/system/script` (termasuk
QuickPrint), gunakan resource umum `system-script`; `mikhmon-sales` hanya
membungkus konvensi laporannya.

---

## 8. QuickPrint

Mikhmon menyimpan preset cetak di `/system/script`:

- `name`: `Quick_Print_<nama paket>`
- `comment`: `QuickPrintMikhmon`
- `source`: string internal `#paket#server#mode#panjang#prefix#char#profile#...`

Padanan Go:

```go
scr, _ := ros.Lookup("system-script")
add, _ := scr.Action("add")

source := "#1jam#all#vc#4##mix#1jam#1d#0#catatan#1d#5000_7000#Disable"
sent, err := add.Sentence("", map[string]string{
    "name":    "Quick_Print_1jam",
    "source":  source,
    "comment": "QuickPrintMikhmon",
}, "")
_, err = run.Exec(ctx, sent)
```

Hapus preset:

```go
req, _ := scr.Build("list", ros.AnyVersion, ros.Once(), ros.Params{
    Query: []ros.Query{{Prop: "comment", Value: "QuickPrintMikhmon"}},
})
rows, _ := run.Once(ctx, req)
del, _ := scr.Action("remove")
for _, row := range rows {
    sent, _ := del.Sentence(row[".id"], nil, "")
    _, _ = run.Exec(ctx, sent)
}
```

`source` QuickPrint adalah format internal Mikhmon; RouterOS hanya
menyimpannya sebagai string. Jangan menganggapnya command API.

---

## 9. Active, host, binding, cookie, DHCP

### Putuskan sesi (kick)

```go
a, _ := ros.Lookup("hotspot-active")
req, _ := a.Build("list", ros.AnyVersion, ros.Once(), ros.Params{
    Query: []ros.Query{{Prop: "user", Value: "vc001"}},
})
rows, _ := run.Once(ctx, req)
remove, _ := a.Action("remove")
for _, row := range rows {
    sent, _ := remove.Sentence(row[".id"], nil, "")
    _, _ = run.Exec(ctx, sent)
}
```

### Host dan IP binding

```go
h, _ := ros.Lookup("hotspot-host")
req, _ := h.Build("list", ros.AnyVersion, ros.Once(), ros.Params{
    Query: []ros.Query{{Prop: "mac-address", Value: "AA:BB:CC:DD:EE:FF"}},
})
rows, _ := run.Once(ctx, req)

bind, _ := h.Action("make-binding")
sent, _ := bind.Sentence(rows[0][".id"], map[string]string{"type": "bypassed"}, "")
_, _ = run.Exec(ctx, sent)
```

IP binding manual:

```go
b, _ := ros.Lookup("hotspot-ip-binding")
add, _ := b.Action("add")
sent, _ := add.Sentence("", map[string]string{
    "mac-address": "AA:BB:CC:DD:EE:FF",
    "type":        "bypassed",
    "comment":     "printer",
}, "")
_, _ = run.Exec(ctx, sent)
```

### Cookie

Hapus semua cookie (Mikhmon `cookie_all`):

```go
ck, _ := ros.Lookup("hotspot-cookie")
req, _ := ck.Build("list", ros.AnyVersion, ros.Once(), ros.Params{})
rows, _ := run.Once(ctx, req)
del, _ := ck.Action("remove")
for _, row := range rows {
    sent, _ := del.Sentence(row[".id"], nil, "")
    _, _ = run.Exec(ctx, sent)
}
```

### DHCP lease

```go
d, _ := ros.Lookup("dhcp-lease")
mk, _ := d.Action("make-static")
sent, _ := mk.Sentence(leaseID, nil, "")
_, _ = run.Exec(ctx, sent)

del, _ := d.Action("remove")
sent, _ = del.Sentence(leaseID, nil, "")
_, _ = run.Exec(ctx, sent)
```

---

## 10. Hapus user kedaluwarsa ala Mikhmon

Mikhmon menandai user kedaluwarsa dengan `limit-uptime=1s`, lalu menghapus
semua yang cocok:

```go
u, _ := ros.Lookup("hotspot-user")
req, _ := u.Build("list", ros.AnyVersion, ros.Once(), ros.Params{
    Query: []ros.Query{{Prop: "limit-uptime", Value: "1s"}},
})
rows, _ := run.Once(ctx, req)
del, _ := u.Action("remove")
for _, row := range rows {
    sent, _ := del.Sentence(row[".id"], nil, "")
    _, _ = run.Exec(ctx, sent)
}
```

Untuk menghapus per batch comment:

```go
req, _ = u.Build("list", ros.AnyVersion, ros.Once(), ros.Params{
    Query: []ros.Query{{Prop: "comment", Value: "vc-123-11.18.26"}},
})
```

Mikhmon juga bisa menghapus user per comment batch dengan query persis
(`?comment=...`). Pencarian prefix `vc-`/`up-` dilakukan di UI, bukan oleh API.

---

## 11. Status voucher (`status/status.php`)

Mikhmon menampilkan expiry dari:

1. comment user, bila sudah berbentuk tanggal; atau
2. `next-run` scheduler dengan `name=<user>`.

Padanan Go:

```go
req, _ := u.Build("list", ros.AnyVersion, ros.Once(), ros.Params{
    Query: []ros.Query{{Prop: "name", Value: "vc001"}},
})
rows, _ := run.Once(ctx, req)
if len(rows) > 0 {
    fmt.Println("comment:", rows[0]["comment"])
}
```

Kalau comment belum bertanggal, cari scheduler:

```go
req, _ = sched.Build("list", ros.AnyVersion, ros.Once(), ros.Params{
    Query: []ros.Query{{Prop: "name", Value: "vc001"}},
})
```

---

## 12. Yang belum ada di katalog

- `/system/logging` (print/add/set) belum menjadi resource. Mikhmon memakainya
  untuk menambahkan rule log `prefix=->`; sementara pakai `run.Exec` mentah
  seperti contoh §3.
- Query API tidak mendukung regex `~`. Script Expire Monitor memakai `~` di
  dalam RouterOS script (itu sah), tetapi dari Go gunakan filter klien atau
  query kesamaan persis.
- `enable`/`disable` adalah verb bawaan menu. Di `hotspot-user` dan
  `system-scheduler`, verifier mendapat `ok`; di `hotspot-user-profile` dan
  `system-script`, `/console/inspect` tidak bisa me-resolve-nya sehingga
  statusnya `unknown`. Aksinya tetap tersedia di katalog dan valid di API.
- Mikhmon mengandalkan scheduler sementara untuk menghitung expiry. Jangan
  menggantinya dengan aritmetika Go saja: tanggal/jam router dan timezone
  adalah sumber kebenaran.

---

## 13. Ringkasan perbedaan pendekatan

| Aspek | Mikhmon v4 | roscatalog |
|---|---|---|
| Bahasa | PHP | Go |
| Sumber kebenaran | UI + session + polling | katalog data + `Plan` |
| Update user/aktif | polling `print` | `list`, `live`, atau `events` (listen) |
| Expiry | scheduler 1 menit + comment | sama; script tetap milik aplikasi |
| Laporan | `/system/script` dengan `owner`/`source` | resource `mikhmon-sales` |
| Voucher | loop `add` + comment batch | action `add` + query comment |
| On-login | string satu baris di profile | action `set` field `on-login` |
| Verifikasi | tidak ada | `cmd/rosverify` via `/console/inspect` |

Intinya: command-command Mikhmon sudah tercakup di katalog, kecuali
`/system/logging`. `on-login`, Expire Monitor, laporan penjualan, dan QuickPrint
bukan command baru — semuanya adalah kombinasi `/ip/hotspot/user/profile`,
`/system/scheduler`, dan `/system/script`.

---

## 14. Paket `mikhmon/`

Katalog root sengaja tetap umum. Untuk pemakaian sehari-hari ala Mikhmon,
tersedia paket `example.com/roscatalog/mikhmon` di folder `mikhmon/`.

Paket itu adalah **overlay**, bukan duplikat katalog:

- `mikhmon.Lookup`/`Resources`/`IDs` hanya mengekspos 22 resource yang dipakai
  Mikhmon (dashboard, hotspot, DHCP, pool/NAT/queue, scheduler/script).
- `mikhmon.OnLoginScript` membangun `on-login` satu baris, termasuk metadata
  harga, normalisasi tanggal v6/v7, scheduler sementara, mode N/X, lock, dan
  laporan `remc`/`ntfc`.
- `mikhmon.ExpireMonitorScript` mengembalikan on-event
  `Mikhmon-Expire-Monitor` yang sama dengan `assets/js/func.js`.
- `mikhmon.Sale` dan `mikhmon.QuickPrintPreset` menyusun/membaca format name,
  source, owner, dan source QuickPrint.
- `mikhmon.Client` membungkus workflow add profile, add user/voucher,
  set on-login, install Expire Monitor, add sale, simpan QuickPrint, normalize
  laporan `owner=*sys`, kick user, dan hapus user `limit-uptime=1s`.

`Client` hanya membutuhkan interface `Executor`; `*rosrun.Runner` sudah
memenuhinya. Contoh lengkap ada di `mikhmon/README.md`.

### Verifikasi langsung

Command `cmd/mikhmon-live` punya dua mode:

1. **Simulasi** (tanpa `-webterm`): `on-login` dijalankan sebagai script
   sementara dengan variabel hotspot (`$user`, `$address`,
   `$"mac-address"`).
2. **Login hotspot nyata** (`-webterm` + `-hotspot`): login HTTP CHAP dilakukan
   dari container webterm GNS3, sehingga `on-login` berjalan di konteks hotspot
   yang sebenarnya.

Yang diperiksa:

- `on-login` tersimpan satu baris dan tidak terpotong;
- sesi hotspot aktif muncul setelah login;
- comment user berubah menjadi tanggal kedaluwarsa + mode `X`;
- scheduler sementara `name=<user>` sudah dihapus;
- script laporan memiliki `comment=mikhmon`, `source=<tanggal>`, dan
  `owner=<bulan><tahun>`;
- query `?source=` (dipakai `get_report.php`) dan `?owner=`
  (dipakai dashboard) menemukan laporan;
- skenario `owner=*sys` ala RouterOS 7.24+ bisa ditemukan lewat
  `?comment=mikhmon&?owner=*sys` lalu dinormalisasi dengan
  `/system/script/set owner=<bulan><tahun>`;
- Expire Monitor mode `X` menghapus user dan mode `N` menyetel
  `limit-uptime=1s`; dengan `-scheduler-wait`, command menunggu tick scheduler
  yang sebenarnya, bukan menjalankan on-event manual.

Hasil uji 2026-11-18 (login hotspot nyata lewat webterm + tick scheduler):

| Router | Versi | Profile + on-login | Voucher | Real login | Sales source | Sales owner | `owner=*sys` | Expire tick X/N |
|---|---|---|---|---|---|---|---|---|
| 192.168.122.26 | 7.22.1 | ok, 1.899 byte | ok | ok | ok | `nov2026` | ok | ok |
| 192.168.122.60 | 6.49.17 | ok, 1.899 byte | ok | ok | ok | `nov2026` | ok | ok |
| 192.168.122.189 | 7.24.6 | ok, 1.899 byte | ok | ok | ok | `*sys` | ok | ok |

Temuan penting dari login hotspot nyata:

- Di 7.22.1 dan 6.49.17, `on-login` menulis laporan dengan
  `owner=nov2026`; dashboard langsung menemukannya.
- Di 7.24.6, login nyata menulis laporan dengan **`owner=*sys`**. Query
  `?owner=nov2026` tidak menemukannya, tetapi query `?source=nov/18/2026`
  (dipakai `get_report.php`) tetap menemukannya. Normalisasi
  `?comment=mikhmon&?owner=*sys` lalu `set owner=nov2026` membuat laporan
  muncul di dashboard.
- Ini mereproduksi gejala “penjualan tidak tampil di frontend” pada keluarga
  7.24.x: record-nya ada, tetapi owner-nya `*sys` sampai dinormalisasi.

Catatan: 7.24.2 tidak tersedia saat pengujian; yang diuji adalah 7.24.6
(long-term). Semua data uji berprefix `zzmik` dihapus kembali, termasuk sesi
hotspot aktif, sehingga tidak ada sisa user/profile/script/scheduler.

Untuk mengulang:

```bash
# simulasi on-login
go run ./cmd/mikhmon-live -addr 192.168.88.1:8728 -user admin -pass secret

# login hotspot nyata dari webterm GNS3
go run ./cmd/mikhmon-live -addr 192.168.88.1:8728 -user admin -pass secret \
  -webterm <container> -hotspot http://10.10.8.1 -scheduler-wait

# laporan JSON
go run ./cmd/mikhmon-live -addr 192.168.88.1:8728 -user admin -pass secret \
  -webterm <container> -hotspot http://10.10.8.1 -scheduler-wait -json
```
