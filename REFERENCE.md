# Katalog command RouterOS untuk Go

Dibuat otomatis dari paket `roscatalog` (`go run gen_reference.go`). 74 resource, 129 view baca, 245 aksi tulis. Dokumen ini dan kode tidak bisa berbeda karena berasal dari data yang sama.

## Cara membaca

| Mode | Pemicu di API | Isi tiap pesan | Snapshot awal | Berhenti |
|---|---|---|---|---|
| `oneshot` | `/path/print` | Semua baris, lalu `!done` | Ya | Otomatis |
| `native` | command `monitor`, `monitor-traffic`, `ping`, `torch`, ... | Satu `!re` per update (nilai terkini) | Tidak relevan | `once`, `duration`, `count`, atau `/cancel` |
| `interval` | `/path/print` + `=interval=` (+ `=stats=`) | Snapshot penuh berulang | Ya, tiap siklus | `/cancel` |
| `follow` | `/path/print` + `=follow=` atau `=follow-only=` | Entri baru (`follow` kirim entri lama dulu) | `follow`: ya; `follow-only`: tidak | `/cancel` |
| `listen` | `/path/listen` | Hanya item yang berubah/muncul; hilang = `=.dead=true` | Tidak | `/cancel` |

Jumlah view per mode: oneshot 88, native 11, interval 14, follow 2, listen 14.

### Stream atau poll

Setiap view bisa dikirim lewat `Plan`: **stream** (router yang mengirim terus) atau **poll** (Go mengulang one-shot tiap N ms). `Resource.Build(view, versi, plan, params)` menghasilkan `Request` untuk `rosrun`.

| Mode view | Stream | Poll |
|---|---|---|
| `oneshot` | tidak ada | `print` tiap Every |
| `interval` | `print =interval=<Every>` (+`=stats=`) | `print` (+`=stats=`) tiap Every, tanpa interval |
| `native` | command tanpa `once` (opsional `=interval=`) | command + `=once=` tiap Every (hanya bila command punya `once`) |
| `follow` | `print =follow=` / `=follow-only=` | `print` tiap Every + diff (`follow-only`: hasil pertama hanya baseline) |
| `listen` | `/listen` (opsional Seed: print awal) | `print` view list tiap Every + diff |

Pilih **stream** untuk data yang sering berubah dengan latensi rendah. Pilih **poll** bila feed banyak, router lemah, jeda perlu diatur per klien/tenant, atau agar tidak menahan stream terbuka di router. Batas: interval stream >= 200ms, jeda poll >= 50ms.

Keterangan kolom dan tanda:

- **Pengiriman** = pilihan `Plan` yang valid untuk view itu; yang pertama adalah bawaan.
- **Proplist `server:`** dikirim sebagai `=.proplist=` (sudah termasuk `.id` untuk tabel). `rosrun` tetap menyaring ulang kolom di sisi klien dan, bila router menolak `.proplist` saat start, mengulang sekali tanpa `.proplist` (`Request.Fallback`).
- Tanda `*` pada argumen berarti wajib. Nilai kosong (mis. `once`) dikirim sebagai `=once=`.
- Aksi: kolom *Selector* menunjukkan argumen pemilih item (`=.id=`, `=numbers=`, `=number=`). Aksi **destruktif** mengubah atau menghapus data/layanan.
- ✅ = nama properti dicocokkan dengan CLI reference RouterOS 7.26. Resource tanpa tanda berasal dari pengetahuan umum RouterOS 6.4x/7.x: validasi di router target dengan `NeverSeen()`.
- ⚠ = tabel berpotensi besar: pakai filter, `count-only`, atau proplist sempit.
- `v6`/`v7` = view khusus mayor tersebut. Properti kosong (mis. `comment` yang belum diisi) memang tidak dikirim router.

Terverifikasi: 13 dari 74 resource (hotspot, PPP, RADIUS, DHCP lease, monitor-traffic).

### Aturan penanganan di Go (goros/v4)

- Library mengelola `.tag` sendiri; jangan kirim `Params.Tag`. Batalkan stream lewat `ListenReply.Cancel()` (`/cancel =tag=`); `!trap` kategori 2 sesudah cancel bukan error.
- Jangan pakai varian `*Context` library dengan ctx yang bisa dibatalkan: ia memanggil `c.r.Cancel()` yang mematikan pembaca bersama semua command. `rosrun` memakai `RunArgs`/`ListenArgs`.
- Channel stream library tidak berpenyangga (`Client.Queue` bawaan 0): konsumen lambat memblokir semua balasan di koneksi itu. `rosrun` menguras channel ke antrean sendiri (batas `MaxBuffered`).
- `listen` tidak punya snapshot awal: `Plan.WithSeed()` memulai listen dulu, baru `print`, lalu menggabungkan per `.id`. Perubahan parsial digabung dengan state sebelumnya.
- `interval` tidak menandai batas siklus. `rosrun` memakai `.section` bila ada, kunci yang berulang, atau jeda tanpa baris (Every/5, 20ms..300ms). Tabel kosong tidak mengirim apa pun; `Runner.EmptySnapshots` mengirim snapshot kosong setelah 3x Every tanpa baris (inferensi).
- API tidak punya `where` dan regex `~`: gunakan query word `?prop=value` (persis) dan filter prefix/regex di klien.
- Deteksi kemampuan menu lewat `InspectChild(...)` (`/console/inspect`).

---

## Hotspot dan Mikhmon

### hotspot-active: Hotspot active (sesi login) ✅

Path: `/ip/hotspot/active`

Filter server `?prop=nilai` (persis): `user`, `server`, `address`, `mac-address`, `login-by`

Listen cocok untuk login/logout; uptime dan bytes berubah terus dan perubahan counter tidak dijamin muncul di listen, jadi ambil lewat view live.

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/ip/hotspot/active/print` | - | server: `.id`, `user`, `address`, `mac-address`, `server`, `login-by`, `uptime`, `session-time-left`, `idle-time`, `bytes-in`, `bytes-out`, `radius` |
| live | interval | stream, poll | `/ip/hotspot/active/print` | `=interval=2s` | server: `.id`, `user`, `address`, `uptime`, `session-time-left`, `idle-time`, `bytes-in`, `bytes-out`<br>snapshot penuh tiap siklus; naikkan interval bila sesi ribuan |
| count | oneshot | once, poll | `/ip/hotspot/active/print` | `=count-only=` | -<br>jumlah ada di =ret=N pada !done |
| events | listen | stream, poll | `/ip/hotspot/active/listen` | - | diambil dari view list lalu disaring di klien<br>tanpa snapshot awal; .dead=true berarti item hilang; kolom diambil dari view list saat Build |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| remove (destruktif) | `/ip/hotspot/active/remove` | `=.id=` | - | -<br>kick: memutus sesi login user |

### hotspot-cookie: Hotspot cookie (mac-cookie) ✅

Path: `/ip/hotspot/cookie`

Filter server `?prop=nilai` (persis): `user`, `mac-address`

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/ip/hotspot/cookie/print` | - | server: `.id`, `user`, `mac-address`, `domain`, `expires-in` |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| remove (destruktif) | `/ip/hotspot/cookie/remove` | `=.id=` | - | -<br>hapus cookie agar user harus login ulang |

### hotspot-host: Hotspot host ✅

Path: `/ip/hotspot/host`

Filter server `?prop=nilai` (persis): `mac-address`, `address`, `server`

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list ⚠ | oneshot | once, poll | `/ip/hotspot/host/print` | - | server: `.id`, `mac-address`, `address`, `to-address`, `server`, `uptime`, `idle-time`, `found-by`, `authorized`, `bypassed`, `comment` |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| remove (destruktif) | `/ip/hotspot/host/remove` | `=.id=` | - | - |
| make-binding | `/ip/hotspot/host/make-binding` | `=numbers=` | - | `comment`, `type`<br>type: regular \| bypassed \| blocked |

### hotspot-ip-binding: Hotspot IP binding (bypass/blokir) ✅

Path: `/ip/hotspot/ip-binding`

Filter server `?prop=nilai` (persis): `mac-address`, `address`, `type`, `server`

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/ip/hotspot/ip-binding/print` | - | server: `.id`, `mac-address`, `address`, `to-address`, `server`, `type`, `disabled`, `comment` |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/ip/hotspot/ip-binding/add` | - | - | `mac-address`, `address`, `to-address`, `server`, `type`, `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/ip/hotspot/ip-binding/set` | `=.id=` | - | `mac-address`, `address`, `to-address`, `server`, `type`, `comment` |
| remove (destruktif) | `/ip/hotspot/ip-binding/remove` | `=.id=` | - | - |
| enable | `/ip/hotspot/ip-binding/enable` | `=.id=` | - | - |
| disable | `/ip/hotspot/ip-binding/disable` | `=.id=` | - | - |

### hotspot-server: Hotspot server ✅

Path: `/ip/hotspot`

Filter server `?prop=nilai` (persis): `name`, `interface`, `profile`, `disabled`

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/ip/hotspot/print` | - | server: `.id`, `name`, `interface`, `address-pool`, `profile`, `idle-timeout`, `keepalive-timeout`, `login-timeout`, `addresses-per-mac`, `disabled`, `comment` |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/ip/hotspot/add` | - | `name`, `interface` | `address-pool`, `profile`, `idle-timeout`, `keepalive-timeout`, `login-timeout`, `addresses-per-mac`, `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/ip/hotspot/set` | `=.id=` | - | `name`, `interface`, `address-pool`, `profile`, `idle-timeout`, `keepalive-timeout`, `login-timeout`, `addresses-per-mac`, ... (+1) |
| remove (destruktif) | `/ip/hotspot/remove` | `=.id=` | - | - |
| enable | `/ip/hotspot/enable` | `=.id=` | - | - |
| disable | `/ip/hotspot/disable` | `=.id=` | - | - |

### hotspot-server-profile: Hotspot server profile

Path: `/ip/hotspot/profile`

Filter server `?prop=nilai` (persis): `name`

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/ip/hotspot/profile/print` | - | server: `.id`, `name`, `hotspot-address`, `dns-name`, `html-directory`, `login-by`, `use-radius`, `rate-limit`, `http-cookie-lifetime`, `install-hotspot-queue`, `mac-auth-mode`, `comment` |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/ip/hotspot/profile/add` | - | `name` | `hotspot-address`, `dns-name`, `html-directory`, `login-by`, `use-radius`, `rate-limit`, `http-cookie-lifetime`, `install-hotspot-queue`, `mac-auth-mode`, `radius-mac-format`, `smtp-server`, `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/ip/hotspot/profile/set` | `=.id=` | - | `name`, `hotspot-address`, `dns-name`, `html-directory`, `login-by`, `use-radius`, `rate-limit`, `http-cookie-lifetime`, ... (+5) |
| remove (destruktif) | `/ip/hotspot/profile/remove` | `=.id=` | - | - |
| enable | `/ip/hotspot/profile/enable` | `=.id=` | - | - |
| disable | `/ip/hotspot/profile/disable` | `=.id=` | - | - |

### hotspot-user: Hotspot user (voucher) ✅

Path: `/ip/hotspot/user`

Filter server `?prop=nilai` (persis): `name`, `profile`, `server`, `mac-address`, `address`, `comment`, `disabled`

Query ?comment=... hanya cocok persis; pencarian prefix (mis. batch vc-) dilakukan di sisi klien.

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list ⚠ | oneshot | once, poll | `/ip/hotspot/user/print` | - | server: `.id`, `name`, `profile`, `server`, `mac-address`, `address`, `limit-uptime`, `limit-bytes-total`, `uptime`, `bytes-in`, `bytes-out`, `comment`, `disabled`<br>ribuan voucher: pakai Count atau filter profile/comment |
| voucher | oneshot | once, poll | `/ip/hotspot/user/print` | - | server: `.id`, `name`, `password`, `profile`, `limit-uptime`, `limit-bytes-total`, `comment`<br>berisi password; untuk cetak voucher saja |
| count | oneshot | once, poll | `/ip/hotspot/user/print` | `=count-only=` | -<br>jumlah ada di =ret=N pada !done |
| events | listen | stream, poll | `/ip/hotspot/user/listen` | - | diambil dari view list lalu disaring di klien<br>sinkronisasi add/remove/ubah user; perubahan counter (uptime, bytes) tidak dijamin muncul lewat listen |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/ip/hotspot/user/add` | - | `name` | `password`, `profile`, `server`, `address`, `mac-address`, `limit-uptime`, `limit-bytes-in`, `limit-bytes-out`, `limit-bytes-total`, `email`, `routes`, `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/ip/hotspot/user/set` | `=.id=` | - | `name`, `password`, `profile`, `server`, `address`, `mac-address`, `limit-uptime`, `limit-bytes-in`, ... (+5) |
| remove (destruktif) | `/ip/hotspot/user/remove` | `=.id=` | - | - |
| enable | `/ip/hotspot/user/enable` | `=.id=` | - | - |
| disable | `/ip/hotspot/user/disable` | `=.id=` | - | - |
| reset-counters | `/ip/hotspot/user/reset-counters` | `=numbers=` | - | -<br>nol-kan uptime dan bytes user |

### hotspot-user-profile: Hotspot user profile (paket voucher)

Path: `/ip/hotspot/user/profile`

Filter server `?prop=nilai` (persis): `name`

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/ip/hotspot/user/profile/print` | - | server: `.id`, `name`, `rate-limit`, `shared-users`, `session-timeout`, `idle-timeout`, `keepalive-timeout`, `status-autorefresh`, `address-pool`, `address-list`, `parent-queue`, `mac-cookie-timeout`, `add-mac-cookie`, `transparent-proxy`, `comment`<br>tanpa on-login/on-logout supaya ringan |
| detail | oneshot | once, poll | `/ip/hotspot/user/profile/print` | - | server: `.id`, `name`, `rate-limit`, `shared-users`, `session-timeout`, `idle-timeout`, `keepalive-timeout`, `status-autorefresh`, `address-pool`, `address-list`, `parent-queue`, `mac-cookie-timeout`, `add-mac-cookie`, `transparent-proxy`, `on-login`, `on-logout`, `insert-queue-before`, `comment`<br>on-login memuat script expire/harga Mikhmon; ambil hanya saat perlu |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/ip/hotspot/user/profile/add` | - | `name` | `rate-limit`, `shared-users`, `session-timeout`, `idle-timeout`, `keepalive-timeout`, `status-autorefresh`, `address-pool`, `address-list`, `parent-queue`, `mac-cookie-timeout`, `add-mac-cookie`, `transparent-proxy`, `on-login`, `on-logout`, `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/ip/hotspot/user/profile/set` | `=.id=` | - | `name`, `rate-limit`, `shared-users`, `session-timeout`, `idle-timeout`, `keepalive-timeout`, `status-autorefresh`, `address-pool`, ... (+8) |
| remove (destruktif) | `/ip/hotspot/user/profile/remove` | `=.id=` | - | - |
| enable | `/ip/hotspot/user/profile/enable` | `=.id=` | - | - |
| disable | `/ip/hotspot/user/profile/disable` | `=.id=` | - | - |

### hotspot-walled-garden: Walled garden (host/path)

Path: `/ip/hotspot/walled-garden`

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/ip/hotspot/walled-garden/print` | - | server: `.id`, `action`, `server`, `src-address`, `dst-address`, `dst-host`, `dst-port`, `path`, `method`, `hits`, `disabled`, `comment` |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/ip/hotspot/walled-garden/add` | - | - | `action`, `server`, `src-address`, `dst-address`, `dst-host`, `dst-port`, `path`, `method`, `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/ip/hotspot/walled-garden/set` | `=.id=` | - | `action`, `server`, `src-address`, `dst-address`, `dst-host`, `dst-port`, `path`, `method`, ... (+1) |
| remove (destruktif) | `/ip/hotspot/walled-garden/remove` | `=.id=` | - | - |
| enable | `/ip/hotspot/walled-garden/enable` | `=.id=` | - | - |
| disable | `/ip/hotspot/walled-garden/disable` | `=.id=` | - | - |
| reset-counters | `/ip/hotspot/walled-garden/reset-counters` | `=numbers=` | - | - |
| reset-counters-all | `/ip/hotspot/walled-garden/reset-counters-all` | - | - | - |

### hotspot-walled-garden-ip: Walled garden IP

Path: `/ip/hotspot/walled-garden/ip`

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/ip/hotspot/walled-garden/ip/print` | - | server: `.id`, `action`, `server`, `src-address`, `dst-address`, `dst-host`, `protocol`, `dst-port`, `disabled`, `comment` |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/ip/hotspot/walled-garden/ip/add` | - | - | `action`, `server`, `src-address`, `dst-address`, `dst-host`, `protocol`, `dst-port`, `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/ip/hotspot/walled-garden/ip/set` | `=.id=` | - | `action`, `server`, `src-address`, `dst-address`, `dst-host`, `protocol`, `dst-port`, `comment` |
| remove (destruktif) | `/ip/hotspot/walled-garden/ip/remove` | `=.id=` | - | - |
| enable | `/ip/hotspot/walled-garden/ip/enable` | `=.id=` | - | - |
| disable | `/ip/hotspot/walled-garden/ip/disable` | `=.id=` | - | - |

### mikhmon-sales: Laporan penjualan Mikhmon (/system/script)

Path: `/system/script`

Filter server `?prop=nilai` (persis): `owner`, `comment`, `name`

Konvensi Mikhmon, bukan perilaku RouterOS: name berformat tanggal-|-jam-|-user-|-harga-|-ip-|-mac-|-validity-|-profile-|-comment, owner berisi bulan+tahun (mis. nov2026), source berisi tanggal, comment=mikhmon. Verifikasi format terhadap versi Mikhmon yang dipakai.

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| report ⚠ | oneshot | once, poll | `/system/script/print` | - | server: `.id`, `name`, `owner`, `source`, `comment`<br>selalu filter ?owner=<bulantahun> |
| count | oneshot | once, poll | `/system/script/print` | `=count-only=` | -<br>jumlah ada di =ret=N pada !done |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/system/script/add` | - | `name`, `source` | `owner`, `comment`<br>balasan sukses berisi =ret=*ID |
| remove (destruktif) | `/system/script/remove` | `=.id=` | - | -<br>hapus laporan per .id; hapus per bulan = print ?owner= lalu remove satu per satu |

---

## PPP, PPPoE dan RADIUS

### ppp-aaa: PPP AAA (RADIUS untuk PPP) ✅

Path: `/ppp/aaa` (tanpa baris)

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/ppp/aaa/print` | - | server: `use-radius`, `accounting`, `interim-update`, `use-circuit-id-in-nas-port-id` |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| set | `/ppp/aaa/set` | - | - | `use-radius`, `accounting`, `interim-update`, `use-circuit-id-in-nas-port-id`, `enable-ipv6-accounting` |

### ppp-active: PPP active (sesi PPPoE aktif) ✅

Path: `/ppp/active`

Filter server `?prop=nilai` (persis): `name`, `service`, `caller-id`, `address`

Listen cocok untuk deteksi pelanggan online/offline; uptime dihitung klien dari waktu event.

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/ppp/active/print` | - | server: `.id`, `name`, `service`, `caller-id`, `address`, `uptime`, `encoding`, `session-id`, `limit-bytes-in`, `limit-bytes-out`, `radius` |
| count | oneshot | once, poll | `/ppp/active/print` | `=count-only=` | -<br>jumlah ada di =ret=N pada !done |
| events | listen | stream, poll | `/ppp/active/listen` | - | diambil dari view list lalu disaring di klien<br>tanpa snapshot awal; .dead=true berarti item hilang; kolom diambil dari view list saat Build |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| remove (destruktif) | `/ppp/active/remove` | `=.id=` | - | -<br>kick: memutus sesi PPPoE (client biasanya reconnect otomatis) |

### ppp-profile: PPP profile (paket) ✅

Path: `/ppp/profile`

Filter server `?prop=nilai` (persis): `name`

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/ppp/profile/print` | - | server: `.id`, `name`, `local-address`, `remote-address`, `rate-limit`, `parent-queue`, `address-list`, `dns-server`, `only-one`, `use-compression`, `use-encryption`, `session-timeout`, `idle-timeout`, `change-tcp-mss`, `comment` |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/ppp/profile/add` | - | `name` | `local-address`, `remote-address`, `rate-limit`, `parent-queue`, `address-list`, `dns-server`, `only-one`, `session-timeout`, `idle-timeout`, `change-tcp-mss`, `on-up`, `on-down`, `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/ppp/profile/set` | `=.id=` | - | `name`, `local-address`, `remote-address`, `rate-limit`, `parent-queue`, `address-list`, `dns-server`, `only-one`, ... (+6) |
| remove (destruktif) | `/ppp/profile/remove` | `=.id=` | - | - |
| enable | `/ppp/profile/enable` | `=.id=` | - | - |
| disable | `/ppp/profile/disable` | `=.id=` | - | - |

### ppp-secret: PPP secret (akun PPPoE) ✅

Path: `/ppp/secret`

Filter server `?prop=nilai` (persis): `name`, `service`, `profile`, `caller-id`, `disabled`

password sengaja tidak masuk proplist list; hanya dikirim saat add/set.

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list ⚠ | oneshot | once, poll | `/ppp/secret/print` | - | server: `.id`, `name`, `service`, `profile`, `caller-id`, `local-address`, `remote-address`, `limit-bytes-in`, `limit-bytes-out`, `last-logged-out`, `last-caller-id`, `last-disconnect-reason`, `disabled`, `comment`<br>ribuan pelanggan: filter name/profile atau Count |
| count | oneshot | once, poll | `/ppp/secret/print` | `=count-only=` | -<br>jumlah ada di =ret=N pada !done |
| events | listen | stream, poll | `/ppp/secret/listen` | - | diambil dari view list lalu disaring di klien<br>tanpa snapshot awal; .dead=true berarti item hilang; kolom diambil dari view list saat Build |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/ppp/secret/add` | - | `name`, `password` | `service`, `profile`, `caller-id`, `local-address`, `remote-address`, `routes`, `limit-bytes-in`, `limit-bytes-out`, `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/ppp/secret/set` | `=.id=` | - | `name`, `password`, `service`, `profile`, `caller-id`, `local-address`, `remote-address`, `routes`, ... (+3) |
| remove (destruktif) | `/ppp/secret/remove` | `=.id=` | - | - |
| enable | `/ppp/secret/enable` | `=.id=` | - | - |
| disable | `/ppp/secret/disable` | `=.id=` | - | - |

### pppoe-server: PPPoE server (service)

Path: `/interface/pppoe-server/server`

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/interface/pppoe-server/server/print` | - | server: `.id`, `service-name`, `interface`, `default-profile`, `authentication`, `keepalive-timeout`, `max-sessions`, `one-session-per-host`, `disabled`, `comment` |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/interface/pppoe-server/server/add` | - | `service-name`, `interface` | `default-profile`, `authentication`, `keepalive-timeout`, `max-sessions`, `one-session-per-host`, `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/interface/pppoe-server/server/set` | `=.id=` | - | `service-name`, `interface`, `default-profile`, `authentication`, `keepalive-timeout`, `max-sessions`, `one-session-per-host`, `comment` |
| remove (destruktif) | `/interface/pppoe-server/server/remove` | `=.id=` | - | - |
| enable | `/interface/pppoe-server/server/enable` | `=.id=` | - | - |
| disable | `/interface/pppoe-server/server/disable` | `=.id=` | - | - |

### radius: RADIUS client ✅

Path: `/radius`

secret tidak masuk proplist list. Monitor menunjukkan counter request/accept/reject per server.

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/radius/print` | - | server: `.id`, `service`, `address`, `authentication-port`, `accounting-port`, `timeout`, `realm`, `src-address`, `accounting-backup`, `status`, `disabled`, `comment` |
| monitor | native | stream, once, poll | `/radius/monitor` | `numbers`* `once` `duration` `interval` | server: `pending`, `requests`, `accepts`, `rejects`, `resends`, `timeouts`, `bad-replies`, `last-request-rtt` |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/radius/add` | - | `service`, `address`, `secret` | `authentication-port`, `accounting-port`, `timeout`, `realm`, `src-address`, `accounting-backup`, `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/radius/set` | `=.id=` | - | `service`, `address`, `secret`, `authentication-port`, `accounting-port`, `timeout`, `realm`, `src-address`, ... (+2) |
| remove (destruktif) | `/radius/remove` | `=.id=` | - | - |
| enable | `/radius/enable` | `=.id=` | - | - |
| disable | `/radius/disable` | `=.id=` | - | - |
| reset-counters | `/radius/reset-counters` | `=numbers=` | - | -<br>nol-kan counter monitor |

### radius-incoming: RADIUS incoming (CoA / Disconnect)

Path: `/radius/incoming` (tanpa baris)

Aktifkan accept=yes agar FreeRADIUS bisa mengirim Disconnect-Request/CoA (port 3799).

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/radius/incoming/print` | - | server: `accept`, `port` |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| set | `/radius/incoming/set` | - | - | `accept`, `port` |

---

## IP (address, route, pool, ARP, DNS, DHCP)

### dhcp-client: DHCP client

Path: `/ip/dhcp-client`

release/renew memakai =numbers=.

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/ip/dhcp-client/print` | - | server: `.id`, `interface`, `status`, `address`, `gateway`, `dhcp-server`, `expires-after`, `add-default-route`, `use-peer-dns`, `disabled`, `comment` |
| events | listen | stream, poll | `/ip/dhcp-client/listen` | - | diambil dari view list lalu disaring di klien<br>tanpa snapshot awal; .dead=true berarti item hilang; kolom diambil dari view list saat Build |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/ip/dhcp-client/add` | - | `interface` | `add-default-route`, `use-peer-dns`, `use-peer-ntp`, `default-route-distance`, `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/ip/dhcp-client/set` | `=.id=` | - | `interface`, `add-default-route`, `use-peer-dns`, `use-peer-ntp`, `default-route-distance`, `comment` |
| remove (destruktif) | `/ip/dhcp-client/remove` | `=.id=` | - | - |
| enable | `/ip/dhcp-client/enable` | `=.id=` | - | - |
| disable | `/ip/dhcp-client/disable` | `=.id=` | - | - |
| release | `/ip/dhcp-client/release` | `=numbers=` | - | -<br>lepas binding lalu restart client |
| renew | `/ip/dhcp-client/renew` | `=numbers=` | - | -<br>perbarui lease |

### dhcp-lease: DHCP lease ✅

Path: `/ip/dhcp-server/lease`

Filter server `?prop=nilai` (persis): `address`, `mac-address`, `server`, `status`, `host-name`, `comment`, `disabled`

active-address dan active-mac-address diisi untuk lease yang sedang terpakai; address/mac-address adalah konfigurasi lease statis.

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list ⚠ | oneshot | once, poll | `/ip/dhcp-server/lease/print` | - | server: `.id`, `address`, `mac-address`, `client-id`, `host-name`, `status`, `server`, `active-address`, `active-mac-address`, `expires-after`, `last-seen`, `dynamic`, `disabled`, `block-access`, `comment` |
| count | oneshot | once, poll | `/ip/dhcp-server/lease/print` | `=count-only=` | -<br>jumlah ada di =ret=N pada !done |
| events | listen | stream, poll | `/ip/dhcp-server/lease/listen` | - | diambil dari view list lalu disaring di klien<br>lease baru/berubah/hilang; perubahan expires-after dan last-seen tidak dijamin muncul lewat listen |
| live | interval | stream, poll | `/ip/dhcp-server/lease/print` | `=interval=5s` | server: `.id`, `address`, `mac-address`, `host-name`, `status`, `expires-after`, `last-seen` |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/ip/dhcp-server/lease/add` | - | `address`, `mac-address` | `server`, `client-id`, `comment`, `block-access`, `rate-limit`<br>balasan sukses berisi =ret=*ID |
| set | `/ip/dhcp-server/lease/set` | `=.id=` | - | `address`, `mac-address`, `server`, `client-id`, `comment`, `block-access`, `rate-limit` |
| remove (destruktif) | `/ip/dhcp-server/lease/remove` | `=.id=` | - | - |
| enable | `/ip/dhcp-server/lease/enable` | `=.id=` | - | - |
| disable | `/ip/dhcp-server/lease/disable` | `=.id=` | - | - |
| make-static | `/ip/dhcp-server/lease/make-static` | `=numbers=` | - | -<br>ubah lease dinamis jadi statis |
| check-status | `/ip/dhcp-server/lease/check-status` | `=numbers=` | - | -<br>bebaskan lease berstatus conflict |

### dhcp-server: DHCP server

Path: `/ip/dhcp-server`

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/ip/dhcp-server/print` | - | server: `.id`, `name`, `interface`, `lease-time`, `address-pool`, `disabled`, `comment` |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/ip/dhcp-server/add` | - | `name`, `interface` | `lease-time`, `address-pool`, `authoritative`, `use-radius`, `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/ip/dhcp-server/set` | `=.id=` | - | `name`, `interface`, `lease-time`, `address-pool`, `authoritative`, `use-radius`, `comment` |
| remove (destruktif) | `/ip/dhcp-server/remove` | `=.id=` | - | - |
| enable | `/ip/dhcp-server/enable` | `=.id=` | - | - |
| disable | `/ip/dhcp-server/disable` | `=.id=` | - | - |

### dhcp-server-network: DHCP server network

Path: `/ip/dhcp-server/network`

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/ip/dhcp-server/network/print` | - | server: `.id`, `address`, `gateway`, `dns-server`, `domain`, `comment` |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/ip/dhcp-server/network/add` | - | `address` | `gateway`, `dns-server`, `domain`, `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/ip/dhcp-server/network/set` | `=.id=` | - | `address`, `gateway`, `dns-server`, `domain`, `comment` |
| remove (destruktif) | `/ip/dhcp-server/network/remove` | `=.id=` | - | - |
| enable | `/ip/dhcp-server/network/enable` | `=.id=` | - | - |
| disable | `/ip/dhcp-server/network/disable` | `=.id=` | - | - |

### ip-address: IP address

Path: `/ip/address`

Filter server `?prop=nilai` (persis): `interface`, `address`, `disabled`, `dynamic`

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/ip/address/print` | - | server: `.id`, `address`, `network`, `interface`, `disabled`, `dynamic`, `invalid`, `comment` |
| events | listen | stream, poll | `/ip/address/listen` | - | diambil dari view list lalu disaring di klien<br>tanpa snapshot awal; .dead=true berarti item hilang; kolom diambil dari view list saat Build |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/ip/address/add` | - | `address`, `interface` | `network`, `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/ip/address/set` | `=.id=` | - | `address`, `interface`, `network`, `comment` |
| remove (destruktif) | `/ip/address/remove` | `=.id=` | - | - |
| enable | `/ip/address/enable` | `=.id=` | - | - |
| disable | `/ip/address/disable` | `=.id=` | - | - |

### ip-arp: ARP

Path: `/ip/arp`

Filter server `?prop=nilai` (persis): `address`, `mac-address`, `interface`, `dynamic`

RouterOS 7 menambah properti status (reachable/stale/...) pada print.

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list ⚠ | oneshot | once, poll | `/ip/arp/print` | - | server: `.id`, `address`, `mac-address`, `interface`, `dynamic`, `disabled`, `published`, `comment` |
| count | oneshot | once, poll | `/ip/arp/print` | `=count-only=` | -<br>jumlah ada di =ret=N pada !done |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/ip/arp/add` | - | `address`, `mac-address`, `interface` | `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/ip/arp/set` | `=.id=` | - | `address`, `mac-address`, `interface`, `comment` |
| remove (destruktif) | `/ip/arp/remove` | `=.id=` | - | - |
| enable | `/ip/arp/enable` | `=.id=` | - | - |
| disable | `/ip/arp/disable` | `=.id=` | - | - |

### ip-dns: DNS (pengaturan)

Path: `/ip/dns` (tanpa baris)

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/ip/dns/print` | - | server: `servers`, `dynamic-servers`, `allow-remote-requests`, `cache-size`, `cache-used` |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| set | `/ip/dns/set` | - | - | `servers`, `allow-remote-requests`, `cache-size`, `cache-max-ttl` |

### ip-dns-cache: DNS cache

Path: `/ip/dns/cache`

Filter server `?prop=nilai` (persis): `name`, `type`

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list ⚠ | oneshot | once, poll | `/ip/dns/cache/print` | - | server: `.id`, `name`, `type`, `data`, `ttl` |
| count | oneshot | once, poll | `/ip/dns/cache/print` | `=count-only=` | -<br>jumlah ada di =ret=N pada !done |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| flush | `/ip/dns/cache/flush` | - | - | -<br>kosongkan seluruh cache DNS |

### ip-dns-static: DNS static

Path: `/ip/dns/static`

Filter server `?prop=nilai` (persis): `name`, `address`, `disabled`

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/ip/dns/static/print` | - | server: `.id`, `name`, `address`, `ttl`, `disabled`, `comment` |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/ip/dns/static/add` | - | `name`, `address` | `ttl`, `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/ip/dns/static/set` | `=.id=` | - | `name`, `address`, `ttl`, `comment` |
| remove (destruktif) | `/ip/dns/static/remove` | `=.id=` | - | - |
| enable | `/ip/dns/static/enable` | `=.id=` | - | - |
| disable | `/ip/dns/static/disable` | `=.id=` | - | - |

### ip-neighbor: Neighbor discovery (MNDP/CDP/LLDP)

Path: `/ip/neighbor`

Filter server `?prop=nilai` (persis): `interface`, `mac-address`, `identity`

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/ip/neighbor/print` | - | server: `.id`, `interface`, `address`, `mac-address`, `identity`, `platform`, `version`, `board`, `age` |
| events | listen | stream, poll | `/ip/neighbor/listen` | - | diambil dari view list lalu disaring di klien<br>tanpa snapshot awal; .dead=true berarti item hilang; kolom diambil dari view list saat Build |

### ip-pool: IP pool

Path: `/ip/pool`

Filter server `?prop=nilai` (persis): `name`

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/ip/pool/print` | - | server: `.id`, `name`, `ranges`, `next-pool`, `comment` |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/ip/pool/add` | - | `name`, `ranges` | `next-pool`, `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/ip/pool/set` | `=.id=` | - | `name`, `ranges`, `next-pool`, `comment` |
| remove (destruktif) | `/ip/pool/remove` | `=.id=` | - | - |
| enable | `/ip/pool/enable` | `=.id=` | - | - |
| disable | `/ip/pool/disable` | `=.id=` | - | - |

### ip-pool-used: IP pool terpakai

Path: `/ip/pool/used`

Filter server `?prop=nilai` (persis): `pool`, `address`, `owner`

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list ⚠ | oneshot | once, poll | `/ip/pool/used/print` | - | server: `.id`, `pool`, `address`, `owner`, `info`<br>selalu filter ?pool=<nama> |
| count | oneshot | once, poll | `/ip/pool/used/print` | `=count-only=` | -<br>jumlah ada di =ret=N pada !done |

### ip-route: Routing table (static/dinamis)

Path: `/ip/route`

Filter server `?prop=nilai` (persis): `dst-address`, `gateway`, `distance`, `active`, `dynamic`, `disabled`

RouterOS 7 memakai routing-table, RouterOS 6 memakai routing-mark. Rute default aktif: ?dst-address=0.0.0.0/0 lalu cek active=true.

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list (v7) | oneshot | once, poll | `/ip/route/print` | - | server: `.id`, `dst-address`, `gateway`, `distance`, `scope`, `target-scope`, `routing-table`, `active`, `dynamic`, `disabled`, `comment` |
| list (v6) | oneshot | once, poll | `/ip/route/print` | - | server: `.id`, `dst-address`, `gateway`, `distance`, `scope`, `target-scope`, `routing-mark`, `active`, `dynamic`, `disabled`, `comment` |
| count | oneshot | once, poll | `/ip/route/print` | `=count-only=` | -<br>jumlah ada di =ret=N pada !done |
| events | listen | stream, poll | `/ip/route/listen` | - | diambil dari view list lalu disaring di klien<br>tanpa snapshot awal; .dead=true berarti item hilang; kolom diambil dari view list saat Build |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/ip/route/add` | - | `dst-address`, `gateway` | `distance`, `scope`, `target-scope`, `routing-table`, `routing-mark`, `check-gateway`, `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/ip/route/set` | `=.id=` | - | `dst-address`, `gateway`, `distance`, `scope`, `target-scope`, `routing-table`, `routing-mark`, `check-gateway`, ... (+1) |
| remove (destruktif) | `/ip/route/remove` | `=.id=` | - | - |
| enable | `/ip/route/enable` | `=.id=` | - | - |
| disable | `/ip/route/disable` | `=.id=` | - | - |

---

## Firewall

### fw-address-list: Firewall address-list (isolir, blacklist)

Path: `/ip/firewall/address-list`

Filter server `?prop=nilai` (persis): `list`, `address`, `dynamic`, `disabled`

timeout pada print adalah sisa waktu; entri berakhir otomatis dan hilang sendiri.

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list ⚠ | oneshot | once, poll | `/ip/firewall/address-list/print` | - | server: `.id`, `list`, `address`, `timeout`, `creation-time`, `dynamic`, `disabled`, `comment`<br>selalu filter ?list=<nama> |
| count | oneshot | once, poll | `/ip/firewall/address-list/print` | `=count-only=` | -<br>jumlah ada di =ret=N pada !done |
| events | listen | stream, poll | `/ip/firewall/address-list/listen` | - | diambil dari view list lalu disaring di klien<br>tanpa snapshot awal; .dead=true berarti item hilang; kolom diambil dari view list saat Build |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/ip/firewall/address-list/add` | - | `list`, `address` | `timeout`, `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/ip/firewall/address-list/set` | `=.id=` | - | `list`, `address`, `timeout`, `comment` |
| remove (destruktif) | `/ip/firewall/address-list/remove` | `=.id=` | - | - |
| enable | `/ip/firewall/address-list/enable` | `=.id=` | - | - |
| disable | `/ip/firewall/address-list/disable` | `=.id=` | - | - |

### fw-connection: Connection tracking

Path: `/ip/firewall/connection`

Filter server `?prop=nilai` (persis): `protocol`, `tcp-state`

API tidak punya regex/prefix, jadi satu IP tidak bisa difilter di sisi router (src-address berisi ip:port). Pakai Count untuk dashboard; untuk detail ambil proplist sempit lalu filter di klien.

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list ⚠ | oneshot | once, poll | `/ip/firewall/connection/print` | - | server: `.id`, `protocol`, `src-address`, `dst-address`, `reply-src-address`, `reply-dst-address`, `tcp-state`, `timeout`, `orig-bytes`, `repl-bytes`<br>tabel bisa puluhan ribu baris; jangan dipakai tanpa Count/proplist |
| count | oneshot | once, poll | `/ip/firewall/connection/print` | `=count-only=` | -<br>jumlah ada di =ret=N pada !done |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| remove (destruktif) | `/ip/firewall/connection/remove` | `=.id=` | - | -<br>putus koneksi per .id |

### fw-filter: Firewall filter

Path: `/ip/firewall/filter`

Filter server `?prop=nilai` (persis): `chain`, `action`, `disabled`, `dynamic`

Counter bytes/packets hanya muncul pada view live (print stats). place-before pada add berisi *ID rule yang akan ditaruh di depannya.

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/ip/firewall/filter/print` | - | server: `.id`, `chain`, `action`, `protocol`, `src-address`, `dst-address`, `src-port`, `dst-port`, `in-interface`, `out-interface`, `in-interface-list`, `out-interface-list`, `src-address-list`, `dst-address-list`, `connection-state`, `log`, `disabled`, `comment` |
| live | interval | stream, poll | `/ip/firewall/filter/print` | `=stats=` `=interval=1s` | server: `.id`, `chain`, `action`, `bytes`, `packets`, `disabled`, `comment`<br>pantau hit rule; cocokkan baris dengan .id |
| count | oneshot | once, poll | `/ip/firewall/filter/print` | `=count-only=` | -<br>jumlah ada di =ret=N pada !done |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/ip/firewall/filter/add` | - | `chain`, `action` | `protocol`, `src-address`, `dst-address`, `src-port`, `dst-port`, `in-interface`, `out-interface`, `in-interface-list`, `out-interface-list`, `src-address-list`, `dst-address-list`, `connection-state`, `log`, `log-prefix`, `place-before`, `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/ip/firewall/filter/set` | `=.id=` | - | `chain`, `action`, `protocol`, `src-address`, `dst-address`, `src-port`, `dst-port`, `in-interface`, ... (+10) |
| remove (destruktif) | `/ip/firewall/filter/remove` | `=.id=` | - | - |
| enable | `/ip/firewall/filter/enable` | `=.id=` | - | - |
| disable | `/ip/firewall/filter/disable` | `=.id=` | - | - |
| move | `/ip/firewall/filter/move` | `=numbers=` | `destination` | -<br>destination berisi *ID item tujuan; item dipindah ke depannya |
| reset-counters | `/ip/firewall/filter/reset-counters` | `=numbers=` | - | - |
| reset-counters-all | `/ip/firewall/filter/reset-counters-all` | - | - | -<br>semua item di menu |

### fw-mangle: Firewall mangle

Path: `/ip/firewall/mangle`

Filter server `?prop=nilai` (persis): `chain`, `action`, `disabled`, `dynamic`

Counter bytes/packets hanya muncul pada view live (print stats). place-before pada add berisi *ID rule yang akan ditaruh di depannya.

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/ip/firewall/mangle/print` | - | server: `.id`, `chain`, `action`, `protocol`, `src-address`, `dst-address`, `src-port`, `dst-port`, `in-interface`, `out-interface`, `in-interface-list`, `out-interface-list`, `src-address-list`, `dst-address-list`, `connection-state`, `connection-mark`, `packet-mark`, `new-connection-mark`, `new-packet-mark`, `new-routing-mark`, `passthrough`, `log`, `disabled`, `comment` |
| live | interval | stream, poll | `/ip/firewall/mangle/print` | `=stats=` `=interval=1s` | server: `.id`, `chain`, `action`, `bytes`, `packets`, `disabled`, `comment`<br>pantau hit rule; cocokkan baris dengan .id |
| count | oneshot | once, poll | `/ip/firewall/mangle/print` | `=count-only=` | -<br>jumlah ada di =ret=N pada !done |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/ip/firewall/mangle/add` | - | `chain`, `action` | `protocol`, `src-address`, `dst-address`, `src-port`, `dst-port`, `in-interface`, `out-interface`, `in-interface-list`, `out-interface-list`, `src-address-list`, `dst-address-list`, `connection-state`, `connection-mark`, `packet-mark`, `new-connection-mark`, `new-packet-mark`, `new-routing-mark`, `passthrough`, `log`, `log-prefix`, `place-before`, `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/ip/firewall/mangle/set` | `=.id=` | - | `chain`, `action`, `protocol`, `src-address`, `dst-address`, `src-port`, `dst-port`, `in-interface`, ... (+16) |
| remove (destruktif) | `/ip/firewall/mangle/remove` | `=.id=` | - | - |
| enable | `/ip/firewall/mangle/enable` | `=.id=` | - | - |
| disable | `/ip/firewall/mangle/disable` | `=.id=` | - | - |
| move | `/ip/firewall/mangle/move` | `=numbers=` | `destination` | -<br>destination berisi *ID item tujuan; item dipindah ke depannya |
| reset-counters | `/ip/firewall/mangle/reset-counters` | `=numbers=` | - | - |
| reset-counters-all | `/ip/firewall/mangle/reset-counters-all` | - | - | -<br>semua item di menu |

### fw-nat: Firewall NAT

Path: `/ip/firewall/nat`

Filter server `?prop=nilai` (persis): `chain`, `action`, `disabled`, `dynamic`

Counter bytes/packets hanya muncul pada view live (print stats). place-before pada add berisi *ID rule yang akan ditaruh di depannya.

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/ip/firewall/nat/print` | - | server: `.id`, `chain`, `action`, `protocol`, `src-address`, `dst-address`, `src-port`, `dst-port`, `in-interface`, `out-interface`, `in-interface-list`, `out-interface-list`, `src-address-list`, `dst-address-list`, `connection-state`, `to-addresses`, `to-ports`, `log`, `disabled`, `comment` |
| live | interval | stream, poll | `/ip/firewall/nat/print` | `=stats=` `=interval=1s` | server: `.id`, `chain`, `action`, `bytes`, `packets`, `disabled`, `comment`<br>pantau hit rule; cocokkan baris dengan .id |
| count | oneshot | once, poll | `/ip/firewall/nat/print` | `=count-only=` | -<br>jumlah ada di =ret=N pada !done |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/ip/firewall/nat/add` | - | `chain`, `action` | `protocol`, `src-address`, `dst-address`, `src-port`, `dst-port`, `in-interface`, `out-interface`, `in-interface-list`, `out-interface-list`, `src-address-list`, `dst-address-list`, `connection-state`, `to-addresses`, `to-ports`, `log`, `log-prefix`, `place-before`, `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/ip/firewall/nat/set` | `=.id=` | - | `chain`, `action`, `protocol`, `src-address`, `dst-address`, `src-port`, `dst-port`, `in-interface`, ... (+12) |
| remove (destruktif) | `/ip/firewall/nat/remove` | `=.id=` | - | - |
| enable | `/ip/firewall/nat/enable` | `=.id=` | - | - |
| disable | `/ip/firewall/nat/disable` | `=.id=` | - | - |
| move | `/ip/firewall/nat/move` | `=numbers=` | `destination` | -<br>destination berisi *ID item tujuan; item dipindah ke depannya |
| reset-counters | `/ip/firewall/nat/reset-counters` | `=numbers=` | - | - |
| reset-counters-all | `/ip/firewall/nat/reset-counters-all` | - | - | -<br>semua item di menu |

### fw-raw: Firewall raw

Path: `/ip/firewall/raw`

Filter server `?prop=nilai` (persis): `chain`, `action`, `disabled`, `dynamic`

Counter bytes/packets hanya muncul pada view live (print stats). place-before pada add berisi *ID rule yang akan ditaruh di depannya.

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/ip/firewall/raw/print` | - | server: `.id`, `chain`, `action`, `protocol`, `src-address`, `dst-address`, `src-port`, `dst-port`, `in-interface`, `out-interface`, `in-interface-list`, `out-interface-list`, `src-address-list`, `dst-address-list`, `connection-state`, `log`, `disabled`, `comment` |
| live | interval | stream, poll | `/ip/firewall/raw/print` | `=stats=` `=interval=1s` | server: `.id`, `chain`, `action`, `bytes`, `packets`, `disabled`, `comment`<br>pantau hit rule; cocokkan baris dengan .id |
| count | oneshot | once, poll | `/ip/firewall/raw/print` | `=count-only=` | -<br>jumlah ada di =ret=N pada !done |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/ip/firewall/raw/add` | - | `chain`, `action` | `protocol`, `src-address`, `dst-address`, `src-port`, `dst-port`, `in-interface`, `out-interface`, `in-interface-list`, `out-interface-list`, `src-address-list`, `dst-address-list`, `connection-state`, `log`, `log-prefix`, `place-before`, `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/ip/firewall/raw/set` | `=.id=` | - | `chain`, `action`, `protocol`, `src-address`, `dst-address`, `src-port`, `dst-port`, `in-interface`, ... (+10) |
| remove (destruktif) | `/ip/firewall/raw/remove` | `=.id=` | - | - |
| enable | `/ip/firewall/raw/enable` | `=.id=` | - | - |
| disable | `/ip/firewall/raw/disable` | `=.id=` | - | - |
| move | `/ip/firewall/raw/move` | `=numbers=` | `destination` | -<br>destination berisi *ID item tujuan; item dipindah ke depannya |
| reset-counters | `/ip/firewall/raw/reset-counters` | `=numbers=` | - | - |
| reset-counters-all | `/ip/firewall/raw/reset-counters-all` | - | - | -<br>semua item di menu |

---

## Queue

### queue-simple: Simple queue

Path: `/queue/simple`

Filter server `?prop=nilai` (persis): `name`, `target`, `parent`, `disabled`, `dynamic`

rate dan packet-rate berformat upload/download (tx/rx) dalam bit/s.

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/queue/simple/print` | - | server: `.id`, `name`, `target`, `max-limit`, `limit-at`, `parent`, `priority`, `queue`, `disabled`, `comment` |
| live | interval | stream, poll | `/queue/simple/print` | `=stats=` `=interval=1s` | server: `.id`, `name`, `target`, `rate`, `packet-rate`, `bytes`, `packets`, `dropped`, `queued-bytes`, `queued-packets`<br>pantau throughput per queue |
| count | oneshot | once, poll | `/queue/simple/print` | `=count-only=` | -<br>jumlah ada di =ret=N pada !done |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/queue/simple/add` | - | `name`, `target` | `max-limit`, `limit-at`, `burst-limit`, `burst-threshold`, `burst-time`, `parent`, `priority`, `queue`, `place-before`, `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/queue/simple/set` | `=.id=` | - | `name`, `target`, `max-limit`, `limit-at`, `burst-limit`, `burst-threshold`, `burst-time`, `parent`, ... (+4) |
| remove (destruktif) | `/queue/simple/remove` | `=.id=` | - | - |
| enable | `/queue/simple/enable` | `=.id=` | - | - |
| disable | `/queue/simple/disable` | `=.id=` | - | - |
| move | `/queue/simple/move` | `=numbers=` | `destination` | -<br>destination berisi *ID item tujuan; item dipindah ke depannya |
| reset-counters | `/queue/simple/reset-counters` | `=numbers=` | - | - |
| reset-counters-all | `/queue/simple/reset-counters-all` | - | - | -<br>semua item di menu |

### queue-tree: Queue tree

Path: `/queue/tree`

Filter server `?prop=nilai` (persis): `name`, `parent`, `packet-mark`, `disabled`

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/queue/tree/print` | - | server: `.id`, `name`, `parent`, `packet-mark`, `limit-at`, `max-limit`, `priority`, `queue`, `disabled`, `comment` |
| live | interval | stream, poll | `/queue/tree/print` | `=stats=` `=interval=1s` | server: `.id`, `name`, `parent`, `rate`, `packet-rate`, `bytes`, `packets`, `dropped`, `queued-packets` |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/queue/tree/add` | - | `name`, `parent` | `packet-mark`, `limit-at`, `max-limit`, `priority`, `queue`, `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/queue/tree/set` | `=.id=` | - | `name`, `parent`, `packet-mark`, `limit-at`, `max-limit`, `priority`, `queue`, `comment` |
| remove (destruktif) | `/queue/tree/remove` | `=.id=` | - | - |
| enable | `/queue/tree/enable` | `=.id=` | - | - |
| disable | `/queue/tree/disable` | `=.id=` | - | - |
| reset-counters | `/queue/tree/reset-counters` | `=numbers=` | - | - |
| reset-counters-all | `/queue/tree/reset-counters-all` | - | - | -<br>semua item di menu |

---

## Interface

### interface: Interface (semua jenis)

Path: `/interface`

Filter server `?prop=nilai` (persis): `name`, `type`, `running`, `disabled`

Counter hanya ada pada view live (print stats). Interface dinamis (PPPoE) kadang tidak mengirim .dead pada listen; lengkapi dengan list berkala.

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/interface/print` | - | server: `.id`, `name`, `type`, `mtu`, `actual-mtu`, `running`, `disabled`, `comment`, `last-link-up-time`, `last-link-down-time`, `link-downs` |
| live | interval | stream, poll | `/interface/print` | `=stats=` `=interval=1s` | server: `.id`, `name`, `running`, `rx-byte`, `tx-byte`, `rx-packet`, `tx-packet`, `rx-drop`, `tx-drop`, `tx-queue-drop`, `rx-error`, `tx-error`<br>counter kumulatif; hitung bps di klien dari selisih, atau pakai monitor-traffic |
| events | listen | stream, poll | `/interface/listen` | - | diambil dari view list lalu disaring di klien<br>tanpa snapshot awal; .dead=true berarti item hilang; kolom diambil dari view list saat Build |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| set | `/interface/set` | `=.id=` | - | `name`, `mtu`, `comment` |
| enable | `/interface/enable` | `=.id=` | - | - |
| disable | `/interface/disable` | `=.id=` | - | - |

### interface-bridge: Bridge

Path: `/interface/bridge`

Filter server `?prop=nilai` (persis): `name`

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/interface/bridge/print` | - | server: `.id`, `name`, `protocol-mode`, `vlan-filtering`, `mtu`, `running`, `disabled`, `comment` |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/interface/bridge/add` | - | `name` | `protocol-mode`, `vlan-filtering`, `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/interface/bridge/set` | `=.id=` | - | `name`, `protocol-mode`, `vlan-filtering`, `comment` |
| remove (destruktif) | `/interface/bridge/remove` | `=.id=` | - | - |
| enable | `/interface/bridge/enable` | `=.id=` | - | - |
| disable | `/interface/bridge/disable` | `=.id=` | - | - |

### interface-bridge-host: Bridge host (tabel MAC)

Path: `/interface/bridge/host`

Filter server `?prop=nilai` (persis): `mac-address`, `on-interface`, `bridge`, `local`

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list ⚠ | oneshot | once, poll | `/interface/bridge/host/print` | - | server: `.id`, `mac-address`, `on-interface`, `bridge`, `vid`, `age`, `local`<br>cari MAC tertentu: ?mac-address=AA:BB:... |

### interface-bridge-port: Bridge port

Path: `/interface/bridge/port`

Filter server `?prop=nilai` (persis): `bridge`, `interface`

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/interface/bridge/port/print` | - | server: `.id`, `interface`, `bridge`, `pvid`, `disabled`, `comment` |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/interface/bridge/port/add` | - | `interface`, `bridge` | `pvid`, `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/interface/bridge/port/set` | `=.id=` | - | `interface`, `bridge`, `pvid`, `comment` |
| remove (destruktif) | `/interface/bridge/port/remove` | `=.id=` | - | - |
| enable | `/interface/bridge/port/enable` | `=.id=` | - | - |
| disable | `/interface/bridge/port/disable` | `=.id=` | - | - |

### interface-ethernet: Ethernet

Path: `/interface/ethernet`

Filter server `?prop=nilai` (persis): `name`, `running`

Monitor memberi status link aktual dan data modul SFP (rx/tx power, suhu, tegangan).

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/interface/ethernet/print` | - | server: `.id`, `name`, `default-name`, `mac-address`, `mtu`, `l2mtu`, `speed`, `auto-negotiation`, `running`, `disabled`, `comment` |
| monitor | native | stream, once, poll | `/interface/ethernet/monitor` | `numbers`* `once` `duration` `interval` | server: `name`, `status`, `rate`, `full-duplex`, `auto-negotiation`, `tx-flow-control`, `rx-flow-control`, `sfp-module-present`, `sfp-type`, `sfp-vendor-name`, `sfp-wavelength`, `sfp-temperature`, `sfp-supply-voltage`, `sfp-tx-bias-current`, `sfp-tx-power`, `sfp-rx-power` |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| set | `/interface/ethernet/set` | `=.id=` | - | `mtu`, `speed`, `auto-negotiation`, `arp`, `advertise`, `tx-flow-control`, `rx-flow-control`, `comment` |
| enable | `/interface/ethernet/enable` | `=.id=` | - | - |
| disable | `/interface/ethernet/disable` | `=.id=` | - | - |
| reset-counters | `/interface/ethernet/reset-counters` | `=numbers=` | - | -<br>nol-kan counter interface |

### interface-list: Interface list

Path: `/interface/list`

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/interface/list/print` | - | server: `.id`, `name`, `comment` |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/interface/list/add` | - | `name` | `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/interface/list/set` | `=.id=` | - | `name`, `comment` |
| remove (destruktif) | `/interface/list/remove` | `=.id=` | - | - |
| enable | `/interface/list/enable` | `=.id=` | - | - |
| disable | `/interface/list/disable` | `=.id=` | - | - |

### interface-list-member: Interface list member

Path: `/interface/list/member`

Filter server `?prop=nilai` (persis): `list`, `interface`

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/interface/list/member/print` | - | server: `.id`, `list`, `interface`, `disabled`, `comment` |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/interface/list/member/add` | - | `list`, `interface` | `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/interface/list/member/set` | `=.id=` | - | `list`, `interface`, `comment` |
| remove (destruktif) | `/interface/list/member/remove` | `=.id=` | - | - |
| enable | `/interface/list/member/enable` | `=.id=` | - | - |
| disable | `/interface/list/member/disable` | `=.id=` | - | - |

### interface-monitor-traffic: Monitor traffic (bps/pps per interface) ✅

Path: `/interface/monitor-traffic` (tanpa baris)

Native streaming: satu !re per interface per detik; =.proplist= didukung.

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| run | native | stream, once, poll | `/interface/monitor-traffic` | `interface`* `once` `duration` | server: `name`, `rx-bits-per-second`, `tx-bits-per-second`, `rx-packets-per-second`, `tx-packets-per-second`, `rx-drops-per-second`, `tx-drops-per-second`, `rx-errors-per-second`, `tx-errors-per-second`, `tx-queue-drops-per-second` |

### interface-vlan: VLAN

Path: `/interface/vlan`

Filter server `?prop=nilai` (persis): `name`, `interface`, `vlan-id`

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/interface/vlan/print` | - | server: `.id`, `name`, `vlan-id`, `interface`, `mtu`, `running`, `disabled`, `comment` |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/interface/vlan/add` | - | `name`, `vlan-id`, `interface` | `mtu`, `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/interface/vlan/set` | `=.id=` | - | `name`, `vlan-id`, `interface`, `mtu`, `comment` |
| remove (destruktif) | `/interface/vlan/remove` | `=.id=` | - | - |
| enable | `/interface/vlan/enable` | `=.id=` | - | - |
| disable | `/interface/vlan/disable` | `=.id=` | - | - |

### pppoe-client: PPPoE client (uplink)

Path: `/interface/pppoe-client`

Field monitor mengikuti pppoe-server/monitor yang terverifikasi, ditambah service-name dan ac-name; validasi di router target.

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/interface/pppoe-client/print` | - | server: `.id`, `name`, `interface`, `user`, `service-name`, `running`, `disabled`, `comment` |
| monitor | native | stream, once, poll | `/interface/pppoe-client/monitor` | `numbers`* `once` `duration` `interval` | server: `status`, `uptime`, `encoding`, `service-name`, `ac-name`, `local-address`, `remote-address` |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/interface/pppoe-client/add` | - | `interface`, `user`, `password` | `name`, `service-name`, `add-default-route`, `use-peer-dns`, `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/interface/pppoe-client/set` | `=.id=` | - | `interface`, `user`, `password`, `name`, `service-name`, `add-default-route`, `use-peer-dns`, `comment` |
| remove (destruktif) | `/interface/pppoe-client/remove` | `=.id=` | - | - |
| enable | `/interface/pppoe-client/enable` | `=.id=` | - | - |
| disable | `/interface/pppoe-client/disable` | `=.id=` | - | - |

---

## Sistem dan operasional

### file: File

Path: `/file`

Filter server `?prop=nilai` (persis): `name`, `type`

Jangan minta properti contents lewat API untuk file besar.

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/file/print` | - | server: `.id`, `name`, `type`, `size`, `creation-time` |
| count | oneshot | once, poll | `/file/print` | `=count-only=` | -<br>jumlah ada di =ret=N pada !done |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| remove (destruktif) | `/file/remove` | `=.id=` | - | -<br>hapus file di router |

### ip-service: Service manajemen (ssh, api, winbox)

Path: `/ip/service`

Filter server `?prop=nilai` (persis): `name`, `disabled`

Matikan service yang tidak dipakai; batasi address ke jaringan manajemen.

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/ip/service/print` | - | server: `.id`, `name`, `port`, `address`, `disabled` |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| set | `/ip/service/set` | `=.id=` | - | `port`, `address` |
| enable | `/ip/service/enable` | `=.id=` | - | - |
| disable | `/ip/service/disable` | `=.id=` | - | - |

### log: Log router

Path: `/log`

API tidak punya regex, jadi filter topics/message dilakukan di klien. follow dan listen hanya mengirim entri baru; list/tail mengirim isi buffer log dulu.

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list ⚠ | oneshot | once, poll | `/log/print` | - | server: `.id`, `time`, `topics`, `message`<br>buffer log memory; batasi dengan membaca sekali saat start |
| tail | follow | stream, poll | `/log/print` | `=follow=` | server: `.id`, `time`, `topics`, `message`<br>kirim entri yang ada dulu, lalu entri baru |
| tail-only | follow | stream, poll | `/log/print` | `=follow-only=` | server: `.id`, `time`, `topics`, `message`<br>hanya entri baru setelah command dijalankan |
| events | listen | stream, poll | `/log/listen` | - | diambil dari view list lalu disaring di klien<br>tanpa snapshot awal; .dead=true berarti item hilang; kolom diambil dari view list saat Build |

### system-backup: Backup

Path: `/system/backup` (tanpa baris)

File backup ada di /file; unduh dengan /tool/fetch (upload) atau FTP/SFTP. load menimpa konfigurasi dan me-reboot router.

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| save | `/system/backup/save` | - | - | `name`, `password`, `dont-encrypt` |
| load (destruktif) | `/system/backup/load` | - | `name` | `password`<br>menimpa konfigurasi |

### system-clock: Clock

Path: `/system/clock` (tanpa baris)

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/system/clock/print` | - | server: `time`, `date`, `time-zone-name`, `time-zone-autodetect` |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| set | `/system/clock/set` | - | - | `time`, `date`, `time-zone-name`, `time-zone-autodetect` |

### system-export: Export konfigurasi

Path: `/export` (tanpa baris)

Gunakan =file= lalu ambil file-nya. RouterOS 7 menyembunyikan secret secara default (show-sensitive untuk menampilkan); RouterOS 6 memakai hide-sensitive.

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| export | `/export/export` | - | - | `file`, `compact`, `terse`, `verbose`, `hide-sensitive`, `show-sensitive` |

### system-health: Health (suhu, tegangan, fan)

Path: `/system/health` (tanpa baris)

RouterOS 7 mengirim baris name/value/type; RouterOS 6 mengirim properti langsung. Sensor yang tersedia bergantung model.

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list (v7) | oneshot | once, poll | `/system/health/print` | - | server: `name`, `value`, `type` |
| live (v7) | interval | stream, poll | `/system/health/print` | `=interval=5s` | server: `name`, `value`, `type` |
| list (v6) | oneshot | once, poll | `/system/health/print` | - | server: `voltage`, `temperature`, `cpu-temperature`, `power-consumption`, `fan1-speed`, `fan2-speed` |
| live (v6) | interval | stream, poll | `/system/health/print` | `=interval=5s` | server: `voltage`, `temperature`, `cpu-temperature`, `power-consumption`, `fan1-speed`, `fan2-speed` |

### system-identity: Identity

Path: `/system/identity` (tanpa baris)

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/system/identity/print` | - | server: `name` |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| set | `/system/identity/set` | - | - | `name` |

### system-package: Paket terpasang

Path: `/system/package`

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/system/package/print` | - | server: `.id`, `name`, `version`, `build-time`, `disabled`, `scheduled` |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| enable | `/system/package/enable` | `=.id=` | - | - |
| disable | `/system/package/disable` | `=.id=` | - | - |

### system-package-update: Update paket

Path: `/system/package/update` (tanpa baris)

check-for-updates mengirim beberapa !re status lalu !done (self-terminating).

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/system/package/update/print` | - | server: `channel`, `installed-version`, `latest-version`, `status` |
| check | native | stream | `/system/package/update/check-for-updates` | - | klien: `channel`, `installed-version`, `latest-version`, `status` |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| set | `/system/package/update/set` | - | - | `channel` |

### system-power: Reboot / shutdown

Path: `/system` (tanpa baris)

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| reboot (destruktif) | `/system/reboot` | - | - | -<br>router putus sampai selesai boot |
| shutdown (destruktif) | `/system/shutdown` | - | - | -<br>router mati; perlu akses fisik untuk menyalakan |

### system-resource: Resource (CPU, RAM, uptime)

Path: `/system/resource` (tanpa baris)

Satu baris tanpa .id. Pada view live, router mengirim baris yang sama berulang.

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/system/resource/print` | - | server: `uptime`, `version`, `build-time`, `architecture-name`, `board-name`, `platform`, `cpu`, `cpu-count`, `cpu-frequency`, `cpu-load`, `free-memory`, `total-memory`, `free-hdd-space`, `total-hdd-space` |
| live | interval | stream, poll | `/system/resource/print` | `=interval=1s` | server: `uptime`, `cpu-load`, `free-memory`, `total-memory`, `free-hdd-space`, `total-hdd-space` |

### system-resource-cpu: CPU per core

Path: `/system/resource/cpu`

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/system/resource/cpu/print` | - | server: `.id`, `cpu`, `load`, `irq`, `disk` |
| live | interval | stream, poll | `/system/resource/cpu/print` | `=interval=1s` | server: `.id`, `cpu`, `load`, `irq`, `disk` |

### system-routerboard: RouterBOARD

Path: `/system/routerboard` (tanpa baris)

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/system/routerboard/print` | - | server: `routerboard`, `model`, `serial-number`, `firmware-type`, `factory-firmware`, `current-firmware`, `upgrade-firmware` |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| upgrade | `/system/routerboard/upgrade` | - | - | -<br>tulis firmware baru; berlaku setelah reboot |

### system-scheduler: Scheduler

Path: `/system/scheduler`

Filter server `?prop=nilai` (persis): `name`, `disabled`

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/system/scheduler/print` | - | server: `.id`, `name`, `start-date`, `start-time`, `interval`, `on-event`, `next-run`, `run-count`, `policy`, `disabled`, `comment` |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/system/scheduler/add` | - | `name`, `on-event` | `start-date`, `start-time`, `interval`, `policy`, `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/system/scheduler/set` | `=.id=` | - | `name`, `on-event`, `start-date`, `start-time`, `interval`, `policy`, `comment` |
| remove (destruktif) | `/system/scheduler/remove` | `=.id=` | - | - |
| enable | `/system/scheduler/enable` | `=.id=` | - | - |
| disable | `/system/scheduler/disable` | `=.id=` | - | - |

### system-script: Script

Path: `/system/script`

Filter server `?prop=nilai` (persis): `name`, `owner`

run memakai =number= (nama atau .id) dan tidak mengembalikan output script.

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/system/script/print` | - | server: `.id`, `name`, `owner`, `policy`, `last-started`, `run-count`, `comment` |
| detail | oneshot | once, poll | `/system/script/print` | - | server: `.id`, `name`, `owner`, `policy`, `dont-require-permissions`, `last-started`, `run-count`, `source`, `comment`<br>source bisa panjang; ambil hanya saat perlu |
| count | oneshot | once, poll | `/system/script/print` | `=count-only=` | -<br>jumlah ada di =ret=N pada !done |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/system/script/add` | - | `name`, `source` | `owner`, `policy`, `dont-require-permissions`, `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/system/script/set` | `=.id=` | - | `name`, `source`, `owner`, `policy`, `dont-require-permissions`, `comment` |
| remove (destruktif) | `/system/script/remove` | `=.id=` | - | - |
| enable | `/system/script/enable` | `=.id=` | - | - |
| disable | `/system/script/disable` | `=.id=` | - | - |
| run | `/system/script/run` | `=number=` | - | -<br>jalankan script |

### user: User router

Path: `/user`

Filter server `?prop=nilai` (persis): `name`, `group`, `disabled`

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/user/print` | - | server: `.id`, `name`, `group`, `address`, `last-logged-in`, `disabled`, `comment` |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/user/add` | - | `name`, `password`, `group` | `address`, `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/user/set` | `=.id=` | - | `name`, `password`, `group`, `address`, `comment` |
| remove (destruktif) | `/user/remove` | `=.id=` | - | - |
| enable | `/user/enable` | `=.id=` | - | - |
| disable | `/user/disable` | `=.id=` | - | - |

### user-active: Sesi admin aktif

Path: `/user/active`

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/user/active/print` | - | server: `.id`, `name`, `when`, `address`, `via`, `group`, `radius` |
| events | listen | stream, poll | `/user/active/listen` | - | diambil dari view list lalu disaring di klien<br>tanpa snapshot awal; .dead=true berarti item hilang; kolom diambil dari view list saat Build |

---

## Troubleshooting (tools)

### tool-bandwidth-test: Bandwidth test

Path: `/tool/bandwidth-test` (tanpa baris)

Menghabiskan bandwidth dan CPU di kedua router; jangan dijalankan tanpa izin pada link produksi.

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| run | native | stream | `/tool/bandwidth-test` | `address`* `user` `password` `protocol` `direction` `local-udp-tx-size` `remote-udp-tx-size` `local-tx-speed` `remote-tx-speed` `tcp-connection-count` `random-data` `duration` | server: `status`, `duration`, `tx-current`, `tx-10-second-average`, `tx-total-average`, `rx-current`, `rx-10-second-average`, `rx-total-average`, `lost-packets`, `direction`, `connection-count`, `local-cpu-load`, `remote-cpu-load` |

### tool-fetch: Fetch (unduh/unggah file, HTTP)

Path: `/tool/fetch` (tanpa baris)

Mengirim beberapa !re progres lalu !done. Pakai upload=yes untuk mengirim backup/export ke server.

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| run | native | stream | `/tool/fetch` | `url`* `dst-path` `mode` `check-certificate` `keep-result` `upload` `src-path` `user` `password` `http-method` `http-data` | server: `status`, `downloaded`, `total`, `duration` |

### tool-netwatch: Netwatch (monitor host up/down)

Path: `/tool/netwatch`

Filter server `?prop=nilai` (persis): `host`, `status`, `type`, `disabled`

Perubahan status muncul sebagai perubahan properti status/since; listen layak dicoba, fallback ke view live bila event tidak datang.

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| list | oneshot | once, poll | `/tool/netwatch/print` | - | server: `.id`, `host`, `type`, `interval`, `timeout`, `status`, `since`, `disabled`, `comment` |
| live | interval | stream, poll | `/tool/netwatch/print` | `=interval=5s` | server: `.id`, `host`, `status`, `since` |
| events | listen | stream, poll | `/tool/netwatch/listen` | - | diambil dari view list lalu disaring di klien<br>tanpa snapshot awal; .dead=true berarti item hilang; kolom diambil dari view list saat Build |

| Aksi | Command | Selector | Wajib | Opsional |
|---|---|---|---|---|
| add | `/tool/netwatch/add` | - | `host` | `type`, `interval`, `timeout`, `up-script`, `down-script`, `comment`<br>balasan sukses berisi =ret=*ID |
| set | `/tool/netwatch/set` | `=.id=` | - | `host`, `type`, `interval`, `timeout`, `up-script`, `down-script`, `comment` |
| remove (destruktif) | `/tool/netwatch/remove` | `=.id=` | - | - |
| enable | `/tool/netwatch/enable` | `=.id=` | - | - |
| disable | `/tool/netwatch/disable` | `=.id=` | - | - |

### tool-ping: Ping

Path: `/ping` (tanpa baris)

Satu !re per probe; field sent/received/packet-loss/min/avg/max bersifat kumulatif. status berisi timeout bila probe gagal.

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| run | native | stream | `/ping` | `address`* `count` `interval` `size` `ttl` `interface` `src-address` `do-not-fragment` `routing-table` | server: `seq`, `host`, `size`, `ttl`, `time`, `status`, `sent`, `received`, `packet-loss`, `min-rtt`, `avg-rtt`, `max-rtt` |

### tool-profile: Profile CPU per proses

Path: `/tool/profile` (tanpa baris)

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| run | native | stream | `/tool/profile` | `cpu` `duration` | server: `name`, `cpu`, `usage` |

### tool-torch: Torch (flow per sumber/tujuan)

Path: `/tool/torch` (tanpa baris)

Tiap update berisi seluruh flow aktif; urutkan di klien menurut tx/rx. Berat di CPU router pada interface yang padat.

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| run | native | stream | `/tool/torch` | `interface`* `src-address` `dst-address` `mac-protocol` `protocol` `port` `ip-protocol` `vlan-id` `duration` | server: `src-address`, `dst-address`, `mac-protocol`, `ip-protocol`, `src-port`, `dst-port`, `vlan-id`, `tx`, `rx`, `tx-packets`, `rx-packets` |

### tool-traceroute: Traceroute

Path: `/tool/traceroute` (tanpa baris)

Router mengirim ulang seluruh daftar hop tiap update.

| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |
|---|---|---|---|---|---|
| run | native | stream | `/tool/traceroute` | `address`* `count` `max-hops` `timeout` `size` `protocol` `port` `src-address` `interface` `use-dns` | server: `address`, `loss`, `sent`, `last`, `avg`, `best`, `worst`, `std-dev`, `status` |

