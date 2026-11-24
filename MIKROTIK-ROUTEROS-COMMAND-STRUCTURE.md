# Struktur Command MikroTik RouterOS

**Panduan berbasis dokumentasi resmi: sintaks, jenis argumen, action, parameter `print`, API, dan REST**

> Disusun 24 November 2026.
> Acuan utama: dokumentasi resmi di `manual.mikrotik.com` (CLI Reference menampilkan versi **7.26**). Situs lama `help.mikrotik.com` sudah dibekukan oleh MikroTik, tetapi isinya masih dipakai sebagai pembanding.
> Forum dan wiki pihak ketiga dipakai hanya sebagai pelengkap dan selalu diberi label.

---

## 0. Cara membaca laporan ini

Setiap contoh command diberi salah satu penanda:

| Penanda | Arti |
|---|---|
| **(dok)** | Contoh diambil langsung dari dokumentasi resmi. |
| **(forum)** | Contoh diambil dari thread forum.mikrotik.com atau sumber komunitas. |
| **(rakitan)** | Contoh disusun dari parameter yang terdokumentasi, tetapi **belum diuji** di perangkat. Uji dulu di router uji sebelum dipakai di produksi. |

Kode sumber (`[D1]`, `[F1]`, dan seterusnya) merujuk ke daftar di Bagian 17.

---

## 1. Ringkasan satu halaman

```
[prefix] [path] command [uparam] [param=[value]] .. [param=[value]]      <- sintaks resmi [D2]
```

Dalam praktik, sebuah baris command terdiri dari:

```
/<menu>/<submenu>  <command>  [uparam]  [switch ...]  [nama=nilai ...]  [where <ekspresi>]
```

| Bagian | Fungsi | Contoh |
|---|---|---|
| `prefix` | `/` menandai path, `:` menandai command global | `/ip/address`, `:put` |
| `path` | Menu bertingkat. Pemisahnya spasi atau `/` | `/ip firewall filter`, `/ip/firewall/filter` |
| `command` | Action di menu itu | `print`, `add`, `set`, `monitor-traffic` |
| `uparam` | Argumen tanpa nama (positional) | `set 0 ...`, `set ether1 ...`, `ping 8.8.8.8` |
| switch | Argumen tanpa nilai (tipe `switch`) | `detail`, `stats`, `terse`, `without-paging` |
| `nama=nilai` | Argumen bernama | `interval=1s`, `proplist=name,rx-byte` |
| `where` | Filter item dengan ekspresi | `where type=ether && running` |

Lima hal yang paling sering terlewat:

1. Tidak ada spasi di antara `nama`, `=`, dan `nilai` (`gateway=3.3.3.3`, bukan `gateway = 3.3.3.3`). [D2]
2. Nomor item (`0`, `1`, ...) dibuat oleh `print`, berlaku per sesi, dan bisa berubah. Untuk script, pakai nama atau `find`. [D3]
3. Parameter `print` tidak seragam di semua menu. Cek dengan **F1** (RouterOS v7), atau lewat `/console/inspect`. [D3][D7]
4. API biner dan REST tidak mendukung `where` dan regex. Filter dilakukan dengan query word (`?type=ether`) atau `.query`. [D5][D6]
5. REST memutus request yang berjalan 60 detik. Command yang terus berjalan harus dibatasi (`count`, `duration`, `once`). [D6]

---

## 2. Sumber dan tingkat keandalan

| Kode | Sumber | Tingkat |
|---|---|---|
| D1 | CLI Reference, halaman indeks: tipe argumen dan modifier `print` standar | Resmi |
| D2 | Scripting manual: sintaks baris command, operator, command global, `where` | Resmi |
| D3 | Console (help.mikrotik.com, beku): hierarki, item, action umum, bantuan F1 | Resmi (arsip) |
| D4 | Command Line Interface (help.mikrotik.com, beku): login, `print` common parameters | Resmi (arsip) |
| D5 | API manual (versi 7.25) | Resmi |
| D6 | REST API manual (versi 7.26) | Resmi |
| D7 | CLI Reference: `console/inspect` | Resmi |
| D8 | CLI Reference: `interface` dan `interface/monitor-traffic` | Resmi |
| D9 | Interface stats and monitor-traffic (`manual.mikrotik.com/docs/diagnostics-monitoring-and-troubleshooting/interface-stats-and-monitor-traffic/`) | Resmi |
| D10 | Referensi command yang dihasilkan otomatis di help.mikrotik.com (contoh `/docs/x/eADWAg`) | Resmi, tetapi hanya dibaca lewat potongan hasil pencarian karena halaman penuhnya menolak akses otomatis |
| F1 | Forum: "Missing SNMP Entries On hAP ac^3" (tampilan F1 untuk `print`) | Komunitas |
| F2 | Forum: "New command in RouterOs 7" (contoh `/console/inspect`) | Komunitas |
| F3 | Forum: "RouterOS LSP ..." (`/console/inspect` untuk syntax check dan completion) | Komunitas |
| F4 | Forum: "Pear2 broken signature" (`print` dengan `follow` lewat API) | Komunitas |
| F5 | Forum: "Insert PPPoe client bandwidth into mysql" (contoh API dengan `.proplist` dan `stats`) | Komunitas |
| C1 | tik4net wiki: "CLI tips and tricks" (perilaku `print as-value`) | Pihak ketiga |

Catatan: dokumentasi baru menyediakan **Markdown mentah per halaman** dengan menambahkan `.md` di belakang alamat halaman. Ini berguna bila dokumentasi dipakai sebagai bahan alat atau asisten. [D-intro, `manual.mikrotik.com/docs/introduction`]

---

## 3. Sintaks dasar

### 3.1 Tata bahasa resmi

Menurut Scripting manual [D2]:

- **`prefix`**: karakter `:` atau `/` yang menandai apakah kata berikutnya command global atau path. Prefix tidak wajib bila command dipanggil dari menu root atau dari menu yang relatif terhadap path-nya.
- **`path`**: path relatif ke level menu yang dituju. Bisa dihilangkan.
- **`command`**: salah satu command yang tersedia di level menu tersebut.
- **`uparam`**: *unnamed parameter*, wajib diisi bila command membutuhkannya.
- **`params`**: rangkaian parameter bernama beserta nilainya (bila parameter itu butuh nilai).
- Akhir baris command ditandai `;` atau baris baru. Command di dalam `()`, `[]`, atau `{}` tidak perlu penanda akhir.

CLI Reference menyebut argumen yang tidak memakai nilai sebagai tipe **`switch`**, yaitu toggle on/off tanpa nilai eksplisit [D1]. Dokumentasi tidak memakai istilah "flag" untuk ini. Kata "flag" di dokumentasi dipakai untuk **penanda status item** di output `print` (`X`, `D`, `R`, dan seterusnya, lihat 4.4). [D8]

### 3.2 Path dan navigasi

Hierarki menu bekerja seperti direktori [D3][D4]:

```
[admin@MikroTik] > /ip route print              <- path penuh + command           (dok)
[admin@MikroTik] > /ip route                    <- masuk ke menu
[admin@MikroTik] /ip/route> print               <- jalankan command di menu itu   (dok)
[admin@MikroTik] > /ip/route/                   <- pemisah "/" juga boleh         (dok)
[admin@MikroTik] /ip/route> ..                  <- naik satu level                (dok)
[admin@MikroTik] /ip/route> /                   <- kembali ke root                (dok)
[admin@MikroTik] /ip/route> /ping 10.0.0.1      <- jalankan dari level lain tanpa pindah  (dok)
[admin@MikroTik] /ip/firewall/nat> .. service-port print   <- relatif dengan ..   (dok)
```

Setiap kata di path boleh dipisah **spasi** atau **`/`** [D4]. Prompt berubah mengikuti posisi menu.

### 3.3 Aturan penulisan nilai

| Aturan | Keterangan | Sumber |
|---|---|---|
| Tanpa spasi di sekitar `=` | `/ip route add gateway=3.3.3.3` benar, `gateway = 3.3.3.3` salah. Berlaku juga untuk `from=`, `to=`, `step=`, `in=`, `do=`, `else=` | D2 |
| Hapus (unset) parameter dengan `!` | Pada `set`, awalan `!` di depan nama parameter mengosongkan parameter itu | D2 |
| String bertanda kutip ganda | Di dalam kutip, `$nama` diganti nilai variabel. Escape: `\"`, `\\`, `\$`, `\n`, `\r`, `\t`, `\_` (spasi) | D2 |
| Komentar | Diawali `#` sampai akhir baris. Tidak ada komentar multibaris | D2 |
| Penyambung baris | Backslash `\` di akhir baris menyambung ke baris berikutnya. Baris yang diakhiri `\` tidak boleh memuat komentar | D2 |
| Substitusi command | `[ ... ]` menjalankan command dan memakai hasilnya | D2 |
| Daftar nomor item | Beberapa item dipisah koma: `set 0,1,2 mtu=1460` | D3 |

Contoh (dok) [D2][D3]:

```
/interface set 0,1,2 mtu=1460
/ip/route/print where gateway~"^[0-9 \\.]*202\$"      # regex, tanda $ di-escape
```

### 3.4 Singkatan dan Tab

Console mengenali dua cara mengetik lebih cepat [D3][D4]:

- **Tab**: melengkapi kata bila awalannya unik. Tab dua kali menampilkan semua kandidat. Tab bekerja untuk nama command, nama argumen, dan nilai yang bersifat daftar tertutup (nama item, nama protokol). Angka dan alamat IP tidak bisa dilengkapi.
- **Singkatan**: awalan kata yang tidak ambigu diterima sebagai nama penuh. Selain awalan, console juga mencari huruf-huruf yang berurutan di dalam nama. Contoh dokumentasi: `interface mt` + Tab menjadi `interface monitor-traffic`.

Praktik yang aman: **jangan memakai singkatan di dalam script**, karena arti singkatan bisa berubah ketika menu mendapat command baru. (Ini saran praktik, bukan pernyataan dokumentasi.)

### 3.5 Tombol dan mode console

| Tombol | Fungsi | Sumber |
|---|---|---|
| **F1** | Bantuan konteks: menampilkan apa yang boleh diketik di posisi kursor (di v7, pengganti `?`) | D3, D4 |
| **Tab** (2x) | Kandidat completion | D3 |
| **Ctrl-C** | Menghentikan command yang sedang berjalan (misalnya `print interval=...`) | D3 |
| **Ctrl-D** | Keluar dari sesi (saat baris kosong) | D3, D4 |
| **Ctrl-X** atau **F4** | Safe mode: perubahan baru disimpan permanen setelah safe mode ditutup | D3 |
| **F7** | Hot-lock: melengkapi command otomatis | D3 |
| **Ctrl-R** atau **F3** | Mencari di riwayat command | D3 |
| **F8** | Menampilkan seluruh input multibaris | D3 |

---

## 4. Jenis argumen

### 4.1 Empat bentuk argumen

| Jenis | Bentuk | Contoh | Catatan |
|---|---|---|---|
| **Unnamed** (`uparam`, positional) | Nilai saja | `set ether1 ...`, `set 0 ...`, `ping 8.8.8.8`, `move 5 destination=0` | Di REST dan API, argumen ini bernama `numbers` atau `.id` |
| **Switch** | Nama saja | `print detail`, `print stats`, `print terse` | Di REST ditulis sebagai string kosong, misalnya `"once":""` |
| **Named** | `nama=nilai` | `interval=1s`, `proplist=name,rx-byte`, `file=backup` | Nilai bertipe sesuai Bagian 4.2 |
| **Query** | Kata `where` lalu ekspresi | `where type=ether && running` | Tipe `query_arg` di referensi. Tidak ada di API biner dan REST (lihat Bagian 12 dan 13) |

Dokumentasi tidak menyatakan aturan urutan antar jenis argumen untuk CLI (kecuali untuk API: urutan attribute word tidak penting). Semua contoh resmi menulis `where` paling akhir. Perlakukan itu sebagai kebiasaan, bukan aturan yang terjamin. [D2][D5]

### 4.2 Tipe nilai di CLI Reference [D1]

| Tipe | Arti |
|---|---|
| `string` | Teks |
| `num` | Bilangan bulat tak bertanda; bisa dibatasi rentang, misalnya `[3 .. 20]`; 32-bit bila rentang tidak disebut |
| `bool` | `true`/`false` atau `yes`/`no` |
| `enum` | Pilihan tetap, nilai valid dicantumkan dalam tanda kurung |
| `time` | Durasi, umumnya dalam detik; boleh memakai satuan `s`, `m`, `h`, `d` |
| `date` | Tanggal dan waktu |
| `switch` | Toggle tanpa nilai |
| `ipAddr`, `ip6Addr` | Alamat IPv4 dan IPv6 |
| `ipPrefix`, `ip6Prefix` | Alamat dengan panjang prefix |
| `ipRange` | Rentang alamat IPv4 |
| `macAddr` | Alamat MAC (boleh dengan wildcard) |
| `address` | Alamat universal, dirender `address (flags=...)`; flag `4`, `6`, `D` (DNS), `i` (nama interface), `/` (CIDR), `v` (VRF), `%` (zona), `+` (address list), dan lainnya |
| `iface_enum` | Nama interface atau nilai khusus seperti `all`, `none`, `dynamic`, `slave` |
| `multi` | Beberapa nilai dari satu himpunan |
| `array`, `object` | Daftar nilai atau struktur nilai |
| `script` | Nama script di `/system script` atau kode script langsung |
| `alt`, `super`, `composite` | Tipe gabungan; subtipe ditulis dalam kurung kurawal dan bisa bersarang |

Tipe data di bahasa script (untuk variabel): `num` (bilangan bulat 64-bit bertanda), `bool`, `str`, `ip`, `ip-prefix`, `ip6`, `ip6-prefix`, `id` (ID internal berawalan `*`), `time`, `array`, `nil`. [D2]

### 4.3 Properti tambahan yang berlaku di semua menu berisi daftar [D1]

| Properti | Fungsi |
|---|---|
| `comment` | Komentar bebas untuk item (maksimum 16 KB) |
| `disabled` | `yes` menonaktifkan item tanpa menghapusnya; item ditandai `X` |
| `copy-from` | Argumen `add`: membuat item baru sebagai salinan item yang ada |
| `place-before` | Argumen `add` dan `move`: menaruh item sebelum item lain |

### 4.4 Flag status di output `print`

Baris `Flags:` di atas output `print` menjelaskan huruf yang muncul di kolom kiri. Untuk `/interface` [D8]:

| Flag | Nama | Arti |
|---|---|---|
| `D` | dynamic | Dibuat otomatis |
| `X` | disabled | Nonaktif |
| `I` | inactive | Tidak aktif |
| `R` | running | Berjalan |
| `S` | slave | Anggota bridge atau bond |
| `P` | passthrough | Passthrough |
| `w` | power-limited | Dibatasi anggaran daya PoE |

Flag berbeda per menu. Misalnya tabel rute memakai `A` (active), `C` (connect), `S` (static), `b` (bgp), dan seterusnya. [D3] Flag ini bisa dipakai di `find` dan `where` sebagai properti bernilai `yes`/`no` atau dipakai langsung, misalnya `where running` [D3][C1].

---

## 5. Menunjuk item

Daftar item (interface, rute, user, dan seterusnya) bisa dirujuk dengan tiga cara. [D3][D4]

| Cara | Contoh | Sifat |
|---|---|---|
| **Nama** | `/interface set ether1 name=WAN` | Stabil karena nama adalah properti item. Hanya ada di menu yang item-nya bernama |
| **Nomor** | `/interface set 0,1,2 mtu=1460` | Diberikan oleh `print`. Berlaku **per sesi dan per daftar**; bisa berbeda di `print` berikutnya |
| **ID internal** | `*1`, `*1A` | Dipakai API dan REST (`=.id=*1`), bertipe `id` |
| **`find`** | `/ip address remove [find interface=ether2]` | Mengembalikan nomor internal item yang cocok; paling aman untuk script |

Hal-hal yang dinyatakan resmi:

- Nomor item tetap bisa dipakai setelah `add`, `remove`, dan `move` selama belum ada `print` baru di sesi itu. Di console v7, nomor juga bisa dipakai tanpa menjalankan `print` terlebih dulu. [D3]
- **Jangan memakai nomor item di script.** Dokumentasi menyarankan `find`. [D3]
- Nomor dibuat terpisah per daftar: `/ip address print` tidak mengubah penomoran daftar interface. [D3]

Contoh gabungan (dok) [D3]:

```
/ip firewall/filter/add chain=forward place-before=[find where comment=CommentX]
/ip/firewall/filter/add chain=forward place-before="CommentX"      # bentuk ringkas yang sama
```


---

## 6. Action umum (berlaku di hampir semua menu)

Dokumentasi menyebut action berikut sebagai command yang berperilaku mirip di banyak level menu: `add`, `edit`, `find`, `move`, `print`, `remove`, `set`, `reset`, `export`, `get`, `enable`, `disable`, `comment`. [D3]

| Action | Fungsi | Argumen penting | Contoh |
|---|---|---|---|
| `print` | Menampilkan item atau pengaturan di menu | Lihat Bagian 7 | `/ip address print` (dok) |
| `add` | Menambah item. Argumennya umumnya sama dengan `set` tanpa nomor item; mengembalikan nomor internal item baru | `copy-from`, `place-before`, `disabled`, `comment` | `/ip address add address=192.168.88.1/24 interface=ether2 comment="LAN"` (rakitan) |
| `set` | Mengubah nilai parameter item (atau pengaturan umum di menu tanpa daftar item). Tidak mengembalikan apa pun. Bisa mengubah beberapa parameter sekaligus | Argumen unnamed berupa nomor item atau daftar nomor; `!param` untuk mengosongkan parameter | `/interface set 0,1,2 mtu=1460` (dok) |
| `remove` | Menghapus item terpilih | Nomor, nama, atau `[find ...]` | `/ip address remove [find interface=ether2]` (rakitan) |
| `enable` / `disable` | Mengaktifkan atau menonaktifkan item | Nomor atau nama | `/ip firewall filter disable 3` (rakitan) |
| `find` | Mencari item. Argumennya sama dengan `set`, ditambah flag bernilai `yes`/`no` (misalnya `disabled`, `active`). Mengembalikan nomor internal item yang cocok | Pasangan `nama=nilai`, atau `where` | `:put [/ip route get [find gateway=1.1.1.1]]` (dok) |
| `get` | Mengambil nilai satu properti. Tidak menampilkannya ke terminal kecuali digabung dengan `:put` | Item, lalu nama properti | `:put [/system/resource/get version]` (dok) |
| `edit` | Mengedit properti bernilai teks panjang (misalnya script) di editor. Untuk script console, batasnya 30 ribu karakter | Item, lalu nama properti | `/system script edit 0 source` (rakitan) |
| `move` | Mengubah urutan item. Argumen pertama: item yang dipindah. Argumen kedua: item yang di depannya item-item tersebut ditempatkan. Bila argumen kedua dihilangkan, item pindah ke akhir | `destination` | `/ip firewall filter move 5 destination=0` (rakitan); `move` di REST: `{".id":"*9","destination":"*C"}` (dok) |
| `reset` | Mengembalikan parameter ke nilai bawaan; boleh sekaligus memberi nilai baru di command yang sama | Bergantung menu | Tidak ada contoh di dokumentasi yang dibaca; cek F1 di menu yang dituju |
| `export` | Mengekspor konfigurasi menu dan submenunya. Dengan `file=` ditulis ke berkas `.rsc`; tanpa `file=` dicetak ke console | `file`, `compact`, `terse`, `verbose`, `hide-sensitive`, `show-sensitive` | `/ip firewall filter export file=filter` (rakitan) |
| `comment` | Memberi komentar pada item | Item, lalu teks | `/ip address comment 0 comment="LAN"` (rakitan) |

Sumber parameter `export`: D10 (potongan), D6 untuk `show-sensitive`.

### 6.1 Contoh `add`, `set`, `find` (dok)

```
/ip/firewall/filter/add chain=forward place-before=[find where comment=CommentX]
/ip firewall filter add chain=blah action=accept protocol=tcp port=123 nth=4,2
/ip firewall filter set 0 !port chain=blah2 !nth protocol=udp       # "!" mengosongkan port dan nth
```

(Dua contoh terakhir mengikuti contoh di Scripting manual [D2].)

---

## 7. `print` secara lengkap

`print` adalah action paling kaya argumen. Daftar di bawah menggabungkan tiga sumber: halaman Console dan Command Line Interface [D3][D4], referensi otomatis di help.mikrotik.com [D10], dan keluaran F1 yang diposting di forum [F1].

### 7.1 Parameter `print`

| Parameter | Jenis | Fungsi | Sumber |
|---|---|---|---|
| `brief` | switch | Output tabel dengan kolom ringkas | D3, D10 |
| `detail` | switch | Output `properti=nilai` untuk tiap item. Tergolong modifier standar yang tidak diulang di tiap halaman menu | D1, D3 |
| `terse` | switch | Format ringkas satu baris per item, ramah mesin | D10, F1 |
| `value-list` | switch | Satu properti per baris | D10, F1 |
| `stats` | switch | **Khusus menu tertentu.** Di `/interface` menampilkan kolom counter (lihat 7.3). Di forum F1, untuk `registration-table` diberi keterangan "subset properti dalam bentuk detail" | D9, D10, F1 |
| `stats-detail` | switch | Khusus `/interface`: menambah waktu link up/down dan jumlah `link-downs` | D9 |
| `without-paging` | switch | Mencetak seluruh output tanpa berhenti tiap layar | D3, D10 |
| `count-only` | switch | Hanya menampilkan jumlah item | D3, D10 |
| `as-value` | switch | Mengembalikan output sebagai nilai (array) untuk scripting. Tergolong modifier standar | D1, D10 |
| `from` | item | Hanya menampilkan item yang disebut, dalam urutan yang diberikan | D3, D10, F1 |
| `where` | query | Menyaring item dengan ekspresi. Sintaksnya mirip `find` | D3, D10 |
| `proplist` | daftar | Daftar nama properti yang disertakan, dipisah koma dan **berurutan** | D4, D10 |
| `interval` | time | Menyegarkan output setiap selang waktu tertentu sampai dihentikan | D3, D10 |
| `follow` | switch | Menampilkan item lalu terus menampilkan perubahan secara real-time | D4, D10 |
| `follow-only` | switch | Hanya menampilkan perubahan yang terjadi setelah command dijalankan | D4, D10 |
| `follow-strict` | switch | Terdaftar di referensi tanpa penjelasan | D4, D10 |
| `file` | string | Menulis hasil `print` ke berkas di router | D3, D10 |
| `append` | switch | Dengan `file`: menambah ke berkas yang ada, bukan menimpa | D4, D10 |
| `oid` | switch | Menampilkan OID SNMP untuk properti yang bisa diakses lewat SNMP | D3, F1 |
| `show-ids` | switch | Menampilkan ID statis item di submenu itu | D4, D10 |

Catatan:

- Daftar parameter **berbeda per menu**. Keluaran F1 di forum memperlihatkan contohnya: `registration-table` di paket `wifi` tidak menawarkan `terse`, `value-list`, dan `where` seperti di paket `wireless`. Selalu cek dengan F1 di menu yang dituju. [F1]
- `proplist` menerima urutan; urutan kolom output mengikuti daftar yang diberikan. [D4]
- Halaman Console [D3] hanya mencantumkan subset (`from`, `where`, `brief`, `detail`, `count-only`, `file`, `interval`, `oid`, `without-paging`). Halaman Command Line Interface [D4] menambah `append`, `follow`, `follow-only`, `follow-strict`, `proplist`, `show-ids`. Referensi otomatis [D10] menambah `as-value`, `terse`, `value-list`. Gabungan ketiganya yang dipakai di tabel.

### 7.2 Contoh `print`

```
# Bentuk dasar dan tampilan
/ip address print                                            (dok)
/ip/route/print where gateway~"^[0-9 \\.]*202\$"              (dok)
/interface/print where (name~"ether")=false                  (dok)

# Counter interface
/interface/print stats                                       (dok)
/interface/print stats-detail                                (dok)

# Kombinasi parameter (disusun dari parameter terdokumentasi)
/interface print stats interval=1s proplist=name,rx-byte,tx-byte   (rakitan)
/interface print detail without-paging where type=ether            (rakitan)
/ip address print terse where interface=ether2                     (rakitan)
/ip address print count-only                                       (rakitan; REST memakai "count-only":"" lalu membalas {"ret":"1"})
/log print follow                                                  (rakitan)

# Dari script
:put [/ip address print as-value]                                  (forum/komunitas, lihat 7.4)
:put [/interface print detail as-value]                            (komunitas [C1])
/system/resource/print oid                                         (rakitan; REST: {"oid":""} membalas OID tiap properti)
```

Bila `interval=` dipakai, output diperbarui terus-menerus sampai dihentikan dengan **Ctrl-C**. [D3]

### 7.3 `stats` dan `stats-detail` pada `/interface`

Dokumentasi resmi [D9] menyebut dua variasi untuk melihat counter interface:

- `/interface print stats`: kolom `NAME`, `RX-BYTE`, `TX-BYTE`, `RX-PACKET`, `TX-PACKET`, `TX-QUEUE-DROP`. Properti berawalan `fp-` adalah counter Fast Path; `tx-queue-drop` menghitung paket yang dibuang antrean interface (sinyal kemacetan).
- `/interface print stats-detail`: menambah `last-link-down-time`, `last-link-up-time`, dan `link-downs` (berguna untuk mendeteksi kabel yang flapping).

Properti read-only di `/interface` yang tampil di output (D8): `rx-byte`, `tx-byte`, `rx-packet`, `tx-packet`, `rx-drop`, `tx-drop`, `tx-queue-drop`, `rx-error`, `tx-error`, `fp-rx-byte`, `fp-tx-byte`, `fp-rx-packet`, `fp-tx-packet`, `fp-rps-drop`, serta `last-link-down-time`, `last-link-up-time`, `link-downs`. Argumen yang bisa diubah: `name`, `mtu`, `l2mtu`.

### 7.4 Catatan `as-value`

- Dokumentasi menjelaskan `as-value` sebagai output berupa nilai untuk scripting. [D10]
- Wiki tik4net (pihak ketiga) mencatat bahwa `print as-value` yang berdiri sendiri di console tidak menampilkan apa pun; hasilnya dibungkus `:put [ ... ]` agar terlihat. Perilaku ini sebaiknya dikonfirmasi di router Anda. [C1]
- Contoh dokumentasi yang memakai hasil `print as-value` sebagai input fungsi lain: `:put [:serialize to=dsv options=dsv.remap delimiter="#" [/ip/address/print as-value]]`. [D2]

---

## 8. Ekspresi `where`

Sintaks `where` pada `print` mirip dengan `find`. [D3] Operator yang terdokumentasi di Scripting manual [D2]:

| Kelompok | Operator | Contoh |
|---|---|---|
| Relasional | `=`, `!=`, `<`, `>`, `<=`, `>=` | `where rx-byte>1000000` (rakitan) |
| Regex | `~` (POSIX extended regular expression) | `where gateway~"^[0-9 \\.]*202\$"` (dok) |
| Logika | `!`, `&&` atau `and`, `\|\|` atau `or`, `in` | `where type=ether && running` (rakitan; pola `and` ada di [C1]) |
| Negasi ekspresi | `<ekspresi>=false` | `where (name~"ether")=false` (dok) |
| Kebalikannya | `<ekspresi>=true` | `where (name~"ether")=true` (dok) |
| Keanggotaan | `in` | `:put (1.1.1.1/32 in 1.0.0.0/8)` (dok) |

### 8.1 Jebakan nama variabel di dalam `where` [D2]

Semua nama properti bawaan RouterOS adalah **nama terpakai**. Variabel lokal yang namanya sama dengan properti akan salah dibaca sebagai properti item yang sedang diperiksa. Contoh resmi: `:local type "ether1"; /interface print where name=$type` tidak mencetak apa pun karena `$type` dibaca sebagai properti `type` tiap interface. Solusinya memakai nama variabel lain, misalnya `$customname`.

### 8.2 `where` tidak ada di API dan REST

Di API biner dan REST, filter dibuat dengan query word `?` atau kunci `.query`. Regex tidak didukung di sana. Lihat Bagian 12 dan 13. [D5][D6]

---

## 9. Action khusus (per menu)

Selain action umum, banyak menu punya command sendiri. CLI Reference mengelompokkan setiap node sebagai **Directory** (menu) atau **Command**. Tabel **Argument** berisi input, tabel **Read-only Argument** berisi field output. [D1][D8]

### 9.1 Contoh yang terdokumentasi

| Command | Fungsi | Argumen yang disebut dokumentasi | Contoh |
|---|---|---|---|
| `/interface monitor-traffic` | Laju trafik per detik, per interface. Beberapa interface tampil berdampingan | `interface` (satu atau lebih, atau `aggregate`), `duration`, `once` (lihat catatan) | `/interface/monitor-traffic ether2 duration=3` (dok [D9]); `/interface monitor-traffic [find]` (dok [D9]) |
| `/ping` | Uji konektivitas | `address`, `count` | `/ping 10.0.0.1` (dok [D3]); REST: `{"address":"192.168.88.10","count":"3"}` (dok [D6]) |
| `/tool traffic-monitor add` | Memicu script saat laju melewati ambang | `name`, `interface`, `traffic`, `trigger`, `threshold`, `on-event` | `/tool/traffic-monitor/add name=high-load interface=ether2 traffic=received trigger=above threshold=100000000 on-event=rate-alert` (dok [D9]) |
| `/tool bandwidth-test` | Uji bandwidth | `address`, `user`, `password`, `duration` | REST (dok [D6]) |
| `/system script run` | Menjalankan script tersimpan | `.id` (nama atau ID) | REST: `{".id":"daily-report"}` (dok [D6]) |
| `/execute` | Menjalankan script sebagai job; `as-string` menunggu dan mengembalikan output | `script`, `as-string`, `file` | `{"script":":put [/system/clock/get date]","as-string":""}` (dok [D6]) |
| `/system package update check-for-updates` | Cek versi RouterOS terbaru | `once` (opsional) | REST (dok [D6]) |
| `/ip firewall nat move` | Memindah aturan | `.id`, `destination` | REST: `{".id":"*9","destination":"*C"}` (dok [D6]) |
| `/file read` | Membaca isi berkas | `file`, `chunk-size` (maks. 32768 byte per request), `offset` | REST (dok [D6]) |
| `/password` | Mengganti password sendiri; command akan meminta input tambahan | `old-password`, `new-password`, `confirm-new-password` | REST (dok [D6]); di CLI, console menampilkan prompt (dok [D4]) |
| `/console inspect` | Menjelaskan struktur command | `request`, `path`, `input` | Bagian 11 |

Field output `/interface monitor-traffic` (D8): `name`, `rx-packets-per-second`, `rx-bits-per-second`, `rx-drops-per-second`, `rx-errors-per-second`, `tx-packets-per-second`, `tx-bits-per-second`, `tx-drops-per-second`, `tx-queue-drops-per-second`, `tx-errors-per-second`, serta versi `fp-` untuk Fast Path.

**Catatan `once`.** Halaman referensi `interface/monitor-traffic` [D8] hanya mencantumkan `interface` (dan menyebut `duration` di deskripsi). Namun halaman REST API menyatakan bahwa command `monitor` menerima `once` untuk satu sampel, dan memberi contoh `/interface/monitor-traffic` dengan `{"interface":"ether1,ether2","once":""}`. [D6] Jadi `once` terdokumentasi, tetapi di REST. Bentuk CLI-nya, `/interface monitor-traffic ether1 once`, adalah rakitan.

**Catatan `interval` pada `monitor-traffic`.** Tidak disebut di halaman resmi yang dibaca, tetapi contoh kode di forum mengirim argumen `interval` ke `/interface/monitor-traffic` lewat API. [F4] Anggap sebagai informasi komunitas; verifikasi di versi RouterOS Anda.

### 9.2 Cara menemukan command khusus suatu menu

Gunakan F1 di menu itu, atau `/console inspect` (Bagian 11), atau halaman CLI Reference menu tersebut (`manual.mikrotik.com/docs/cli-reference/<menu>/`).

---

## 10. Command global (`:`) yang sering dipakai bersama `print`

Command global diawali `:` (bukan path menu). Beberapa yang terdokumentasi di Scripting manual [D2]:

| Command | Fungsi | Contoh |
|---|---|---|
| `:put` | Mencetak ke terminal | `:put "Hello world"` |
| `:log` | Menulis ke log sistem (topik `debug`, `error`, `info`, `warning`) | `:log info "Hello from script"` |
| `:delay` | Menunggu | `:delay 3s` |
| `:len` | Panjang string atau jumlah elemen array | `:put [:len [/ip route find]]` |
| `:typeof` | Jenis data | `:put [:typeof 4]` |
| `:serialize` | Mengubah nilai atau array menjadi JSON atau dsv | `:put [:serialize to=json value=a,b,c]` |
| `:deserialize` | Kebalikannya | `:put [:deserialize from=json value="[\"a\",\"b\"]"]` |
| `:execute` | Menjalankan script di latar belakang (maks. 64 kB) | lihat `/execute` di REST |
| `:parse` | Mengubah string menjadi fungsi | `:global myFunc [:parse ":put hello!"]; $myFunc;` |
| `:onerror` | Menangkap error | `:onerror e in={ :error "failure" } do={ :put "Critical $e" }` |
| `:for`, `:if` | Perulangan dan percabangan | `:for i from=1 to=2 do={ :put $i }`, `:if ( true ) do={ :put "lala" }` |
| `:environment print` | Menampilkan variabel yang diinisialisasi | `:global myVar true; :environment print` |
| `:grep` | Menjalankan script dan memfilter baris keluarannya | `:grep script="/interface print" pattern="ether" after=1 before=1 filename=results.txt` |

Pola dasar yang dipakai bersama `print`/`find`:

```
:put [/ip route get [find gateway=1.1.1.1]]          # tiga baris command bersarang (dok)
:put [/system/resource/get version]                  # get + :put (dok)
```

Variabel: `:global nama [nilai]` untuk global, `:local nama [nilai]` untuk lokal. Tiap baris yang diketik di console interaktif adalah scope lokal tersendiri. Nama variabel yang mengandung selain huruf dan angka harus dikutip, misalnya `:global "my-var"`. [D2]

---

## 11. Menemukan struktur command secara interaktif dan terprogram

### 11.1 Interaktif

- **F1** setelah menulis command: menampilkan argumen yang tersedia di posisi itu beserta penjelasannya (v7). Tab dua kali memberi bentuk ringkas. [D3]
- Contoh keluaran F1 untuk `print` di forum (potongan; menunjukkan daftar `append`, `as-value`, `brief`, `count-only`, `detail`, `file`, `follow`, `follow-only`, `follow-strict`, `from`, `interval`, `oid`, `proplist`, `show-ids`, `stats`, `terse`, `value-list`, `where`). [F1]

### 11.2 Terprogram: `/console/inspect`

Dokumentasi 7.26 menjelaskan `console/inspect` sebagai command yang mengembalikan pohon command, bantuan bawaan, kandidat tab-completion, dan data syntax highlighting yang dipakai console. Script dan alat memakainya untuk menemukan menu, command, argumen, dan tipe nilai, atau untuk memeriksa baris command sebelum dijalankan. [D7]

| Argumen | Tipe | Fungsi |
|---|---|---|
| `request` | ubit: `self`, `child`, `completion`, `highlight`, `syntax`, `error` | Jenis permintaan; boleh digabung dengan koma, misalnya `child,completion` |
| `path` | daftar nama dipisah koma | Menu atau command yang dituju, misalnya `ip,address,add`. **Tanpa slash dan tanpa slash di depan.** Kosong berarti root |
| `input` | string | Baris command untuk permintaan `completion`, `highlight`, `syntax` |

Field output penting: `type`, `name`, `node-type` (`path`, `dir`, `cmd`, `arg`), `completion`, `style`, `offset`, `preference`, `show`, `highlight`, `symbol`, `symbol-type` (`collection`, `explanation`, `definition`), `nested`, `text`. [D7]

Contoh:

```
# Daftar anak dari /ip/address/add: command, argumen (rakitan dari contoh REST di D6)
/console inspect request=child path=ip,address,add

# Penjelasan sintaks command tersebut (forum [F2])
/console inspect request=syntax path=ip,address,add

# Memeriksa kebenaran satu baris command (forum [F3])
:put [:serialize to=json [/console/inspect request=highlight input="/ip/address/add interface=ether1 address=1.1.1.1/24" as-value ]]

# Daftar kandidat completion (forum [F3])
/console/inspect request=completion input="/ip/address/add interface=ether1 address=1.1.1.1/24"
```

Via REST (dok [D6]):

```bash
curl --cacert rest-ca.crt -u monitoring -X POST \
    https://192.168.88.1/rest/console/inspect \
    -H "Content-Type: application/json" \
    --data '{"request":"child","path":"ip,address,add",
        ".proplist":"name,node-type"}'
```

Balasan contoh di dokumentasi: `[{"name":"add","node-type":"cmd"},{"name":"address","node-type":"arg"}, ...]`.

Catatan sejarah: thread forum F3 (editor LSP) menyebut `/console/inspect` sebagai fitur yang saat itu tidak terdokumentasi. Di versi 7.26, command ini sudah punya halaman resmi di CLI Reference [D7]. Rincian permintaan dan contohnya ada di halaman "Console inspect" di bagian Scripting.


---

## 12. API biner (port 8728 / 8729)

Halaman yang dibaca adalah versi **7.25** [D5]; versi 7.26 bisa memiliki perbedaan kecil. API mengikuti sintaks CLI, dan layanan API harus diaktifkan lebih dulu. Port bawaan: **8728** dan **8729** (aman). [D5]

### 12.1 Struktur kalimat dan kata

- Komunikasi berupa **sentence**: rangkaian **word** yang diakhiri word berpanjang nol. Router baru mulai mengevaluasi setelah word kosong diterima. Tanpa itu, semua input berikutnya dianggap bagian kalimat yang sama.
- Setiap word dikodekan sebagai panjang lalu isi. Panjang 1 sampai 4 byte atau 5 byte (awalan `0xF0`) tergantung nilainya; byte pertama `>= 0xF8` adalah byte kontrol cadangan.
- Urutan: word pertama harus **command word**. Atribut dan query boleh dalam urutan bebas, kecuali query yang **urutannya penting**.

| Jenis word | Bentuk | Contoh | Catatan |
|---|---|---|---|
| Command | `/` + path CLI dengan spasi diganti `/` | `/interface/print`, `/user/active/listen` | Perintah khusus API: `login`, `cancel` |
| Attribute | `=nama=nilai` | `=address=10.0.0.1`, `=comment=` (nilai kosong), `=name=iu=c3Eeg` (nilai boleh memuat `=`) | Urutan tidak penting dan jangan diandalkan |
| API attribute | `.nama=nilai` | `.tag=2` | Saat ini satu-satunya: `tag` |
| Query | `?...` | `?type=ether`, `?#\|` | Hanya `print` dan `getall` yang memprosesnya |
| Reply | `!re`, `!done`, `!trap`, `!empty`, `!fatal` | | `!re` data, `!done` akhir, `!trap` error, `!empty` tanpa data, `!fatal` koneksi akan ditutup |

### 12.2 Login, tag, cancel, listen, getall

```
/login
=name=admin
=password=
```

- Tag: bila sebuah kalimat membawa `.tag=`, semua balasannya membawa tag yang sama. Berguna untuk menjalankan beberapa command bersamaan. [D5]
- `/cancel`: menghentikan command yang berjalan. Argumen opsional `=tag=<tag>`; tanpa argumen, semua command yang berjalan dihentikan. `/cancel` tidak membatalkan dirinya sendiri dan bisa punya `.tag` sendiri. Command yang dibatalkan biasanya membalas `!trap` (category 2, "interrupted") dan `!done`. [D5]
- `listen`: tersedia di mana `print` tersedia (tidak selalu berfungsi). Mengirim `!re` setiap ada perubahan di daftar item; item yang hilang ditandai `=.dead=yes`. Tidak berhenti sendiri; hentikan dengan `/cancel`. [D5]
- `getall`: alias `print`; balasannya menyertakan `=.id=`. [D5]

### 12.3 Perbedaan `print` di API dibanding CLI [D5]

| CLI | API |
|---|---|
| `where` | **Tidak didukung.** Gunakan query word |
| `proplist=a,b` | `=.proplist=a,b` (daftar dipisah koma) |
| Nomor item | `=.id=*1` (ID internal) |
| Regex (`~`) | **Tidak didukung** |
| `print` tanpa `proplist` | Semua properti dikembalikan, termasuk yang lambat (isi berkas, counter kinerja). Dokumentasi menganjurkan selalu memakai `.proplist`; tanpa itu, penalti kinerja bisa besar terutama bila `=detail=` dipakai |
| `set` dengan `[find ...]` | `set` **tidak menerima query**. Jalankan `print` dengan `=.proplist=.id` dan query, lalu `set` per ID |

Aturan `.proplist`: properti yang tidak dimiliki item tidak muncul (dan `?name` bernilai false untuk item itu); item bisa menyertakan properti tambahan; urutan properti tidak boleh diandalkan; entri ganda berperilaku tak terdefinisi. [D5]

### 12.4 Query word

Query dievaluasi per item memakai **stack** boolean (awalnya berisi `true` tak terhingga). Item diproses bila di akhir tidak ada nilai `false` di stack. [D5]

| Query | Arti |
|---|---|
| `?name` | `true` bila item punya properti `name` |
| `?-name` | `true` bila item **tidak** punya properti `name` |
| `?name=x` | `true` bila nilai `name` sama dengan `x` |
| `?<name=x` | `true` bila nilai `name` lebih kecil dari `x` |
| `?>name=x` | `true` bila nilai `name` lebih besar dari `x` |
| `?#ops` | Operasi stack dari kiri ke kanan: `!` (NOT), `&` (AND), `\|` (OR); angka desimal = indeks stack (puncak = 0); `.` setelah karakter menyalin nilai puncak |

Contoh resmi [D5]:

```
# CLI: /interface/print where type=ether || type=vlan
/interface/print
?type=ether
?type=vlan
?#|

# CLI: /ip/route/print where comment      (rute dengan komentar tidak kosong)
/ip/route/print
?>comment=

# CLI: /ip/route/print where !(distance>1 && gateway=172.16.1.1)
/ip/route/print
?>distance=1
?gateway=172.16.1.1
?#&!
```

### 12.5 OID dan `!trap`

- OID: di CLI, `print oid`. Di API, properti OID bernama `<properti>.oid` dan diminta lewat `.proplist`, misalnya `=.proplist=uptime,cpu-load,uptime.oid,cpu-load.oid`. [D5]
- `!trap` membawa `=message=` dan `=category=`: `0` item atau command tidak ada, `1` nilai argumen gagal, `2` eksekusi diinterupsi, `3` kegagalan script, `4` kegagalan umum, `5` terkait API, `6` terkait TTY, `7` nilai dari `:return`. [D5]

### 12.6 Mengubah item lewat API [D5]

```
# CLI: /interface set [find where type=ether] mtu=1500
/interface/print
=.proplist=.id
?type=ether
# -> !re =.id=*1 ... !re =.id=*2 ...

/interface/set
=.id=*1
=mtu=1500

/interface/set
=.id=*2
=mtu=1500
```

### 12.7 Catatan komunitas (forum)

- Thread F4: penulis library PEAR2 menyatakan bahwa `print` dengan argumen `follow` adalah command yang valid di API, hampir setara `listen`, tetapi **menerima query** dan **juga menyertakan item yang sudah ada**. Ini pernyataan forum, bukan halaman API resmi. Verifikasi di versi Anda sebelum bergantung padanya. [F4]
- Thread F5: kode klien memakai `/interface/print` dengan `.proplist=name,rx-byte,tx-byte` ditambah argumen `stats`, serta `/interface/listen`. Menunjukkan bahwa `stats` dikirim sebagai argumen API oleh klien komunitas (era RouterOS 5/6). [F5]

---

## 13. REST API (HTTP/JSON)

REST API memberi akses HTTP ke console. Setiap menu menjadi URL di bawah `/rest` (`/ip/address` menjadi `https://<router>/rest/ip/address`). REST adalah pembungkus JSON dari API biner: web server login ke API dengan akun terotentikasi lalu menjalankan tiap request sebagai command API. Karena itu menu, parameter, dan izinnya sama dengan console. [D6]

### 13.1 Metode HTTP ke command console [D6]

| Metode | Command | Balasan sukses |
|---|---|---|
| `GET` | `print` | `200 OK`; satu record atau daftar record |
| `PUT` | `add` | `201 Created`; record baru, URL-nya di header `Location` |
| `PATCH` | `set` | `200 OK`; record yang berubah |
| `DELETE` | `remove` | `204 No Content` |
| `POST` | Command apa pun | `200 OK`; balasan command |

Body request berupa objek JSON dengan header `Content-Type: application/json`; tanpa header itu, router membalas `415`.

### 13.2 Membaca data

```bash
# print proplist=.id,address,interface            (dok [D6])
curl --cacert rest-ca.crt -u monitoring -G \
    https://192.168.88.1/rest/ip/address \
    -d .proplist=.id,address,interface

# satu record, lewat ID atau nama
curl ... "https://192.168.88.1/rest/ip/address/*1" -d .proplist=address,network,interface
curl ... https://192.168.88.1/rest/interface/ether1 -d .proplist=name,type,running,mtu

# filter: semua parameter harus cocok persis
curl ... -G https://192.168.88.1/rest/interface -d type=ether -d running=true -d .proplist=name
```

- ID berawalan `*`; kutip URL agar shell tidak memperluasnya, dan **jangan** di-encode sebagai `%2A`. [D6]
- Karakter terpesan lain (misalnya spasi) di-encode (`%20`). [D6]
- `GET` pada daftar selalu menyertakan `.id`. Bila tidak ada record yang cocok, balasannya `[]`. [D6]
- Menu tanpa record (`/system/identity`, `/ip/dns`) mengembalikan satu objek. [D6]
- Kondisi selain "sama persis" (OR, NOT, lebih kecil, dan seterusnya) memakai `.query` di body `POST ... /print`. `.query` **bukan** query string di URL. [D6]

### 13.3 Menjalankan command dengan `POST`

URL adalah path menu plus nama command; body berisi parameternya. [D6]

```bash
# print dengan properti dan filter
curl --cacert rest-ca.crt -u monitoring -X POST \
    https://192.168.88.1/rest/interface/print \
    -H "Content-Type: application/json" \
    --data '{".proplist":"name,type",
        ".query":["type=ether","type=vlan","#|!"]}'

# count-only: balasan {"ret":"1"}
curl ... -X POST https://192.168.88.1/rest/ip/address/print \
    -H "Content-Type: application/json" --data '{"count-only":""}'

# satu sampel trafik (monitor + once)
curl ... -X POST https://192.168.88.1/rest/interface/monitor-traffic \
    -H "Content-Type: application/json" \
    --data '{"interface":"ether1,ether2","once":"",
        ".proplist":"name,rx-bits-per-second,tx-bits-per-second"}'

# oid
curl ... -X POST https://192.168.88.1/rest/system/resource/print \
    -H "Content-Type: application/json" --data '{"oid":""}'
```

Aturan nilai [D6]:

- Balasan **selalu berupa string**, termasuk angka dan boolean (`"true"`, `"false"`). Konversi sendiri.
- Request boleh mengirim string atau angka/boolean JSON (`"disabled":"yes"` sama dengan `"disabled":true`).
- Switch tanpa nilai (`once`, `compact`, `as-string`, `count-only`, `oid`) dikirim sebagai **string kosong**: `"once":""`.
- Argumen unnamed console menjadi `numbers`, misalnya `{"numbers":"ether2","once":""}` untuk `/interface/ethernet/monitor ether2`.
- `.proplist` (string dipisah koma atau daftar string) berlaku untuk command apa pun yang mengembalikan record; `.query` hanya untuk `print`. `.query` adalah daftar query word API tanpa awalan `?`.
- Command yang mengembalikan data di `!done` menghasilkan objek, biasanya nilai di `ret`: ID record baru untuk `add`, jumlah record untuk `print count-only`, keluaran `execute` dengan `as-string`.

### 13.4 Batas dan perilaku penting [D6]

| Topik | Keterangan |
|---|---|
| Batas 60 detik | Router menghentikan request yang berjalan 60 detik dengan `400 Bad Request` ("Session closed"). `duration` tidak memperpanjang batas ini |
| Command yang tak berhenti | `ping` tanpa `count`, command `monitor`, `bandwidth-test` tanpa `duration` selalu kena batas. Batasi dengan `count`, `once`, atau `duration` |
| Tanpa streaming | Balasan baru tiba setelah command selesai |
| Serial per user | Request satu user berjalan satu per satu; user berbeda berjalan paralel. Beri dashboard dan script otomasi user yang terpisah |
| Sesi | Login ke API disimpan 10 menit. Perubahan policy grup berlaku saat login diperbarui (hingga 10 menit) |
| Tanpa CORS | Router tidak mengirim header CORS dan menjawab preflight `OPTIONS` dengan `503`; browser tidak bisa memanggil langsung, perlu backend atau proxy |
| Otentikasi | HTTP Basic; gunakan HTTPS (`www-ssl`) dengan sertifikat |
| Izin user | Grup perlu policy `rest-api` **dan** `api`, ditambah policy untuk aksinya (`read` untuk membaca, `write` untuk mengubah, `test` untuk `ping`, `ftp` untuk menulis berkas, `sensitive` untuk melihat password dan kunci). Contoh grup baca saja: `/user/group/add name=rest-read policy=read,api,rest-api` |

### 13.5 Kode error yang sering muncul [D6]

| Balasan | Sebab |
|---|---|
| `401 Unauthorized` | Kredensial salah, grup tanpa `rest-api`, atau `address` user tidak mengizinkan klien |
| `500` "not allowed" | Grup punya `rest-api` tapi tidak punya `api` |
| `400` "no such command or directory (x)" | Menu, command, atau nama record `x` tidak ada |
| `400` "missing or invalid resource identifier" | `PATCH`/`DELETE` tanpa ID record, nilai ditolak, atau user tanpa `write` |
| `400` "Session closed" | Command melewati 60 detik atau menunggu di belakang request lain user yang sama |
| `415` | Body tanpa header `Content-Type: application/json` |

---

## 14. Peta padanan CLI, API biner, dan REST

| Tujuan | CLI | API biner | REST |
|---|---|---|---|
| Pilih properti | `print proplist=a,b` | `/path/print` + `=.proplist=a,b` | `GET /rest/path?.proplist=a,b` atau `POST .../print {".proplist":"a,b"}` |
| Filter sama dengan | `print where type=ether` | `?type=ether` | `GET /rest/interface?type=ether` |
| OR | `where type=ether \|\| type=vlan` | `?type=ether` `?type=vlan` `?#\|` | `POST .../print {".query":["type=ether","type=vlan","#\|"]}` |
| Regex | `where name~"ether"` | **Tidak ada** | **Tidak ada** |
| Hitung item | `print count-only` | `=count-only=` (rakitan, REST membungkus API) | `POST .../print {"count-only":""}` → `{"ret":"1"}` |
| Detail | `print detail` | `=detail=` (disebut di [D5]) | (tak diperlukan; balasan sudah berisi properti) |
| `stats` | `/interface print stats` | `stats` dikirim klien komunitas [F5]; tidak dijelaskan di [D5] | Tidak terdokumentasi; kemungkinan `{"stats":""}` (rakitan) |
| OID | `print oid` | `=.proplist=uptime,uptime.oid` | `{"oid":""}` |
| Satu sampel monitor | `monitor-traffic ether1 once` (rakitan) | `/interface/monitor-traffic` + `=interface=ether1` + `=once=` (rakitan) | `{"interface":"ether1","once":""}` |
| Pembaruan berkala | `print interval=1s` atau `follow` | `listen`, lalu `/cancel` [D5]; `print follow` menurut forum [F4] | Tidak ada streaming; polling dengan `once`/`count` |
| Ubah item | `set 0 mtu=1500` | `/set` + `=.id=*1` + `=mtu=1500` | `PATCH /rest/path/*1` atau `POST .../set` |
| Tambah | `add ...` | `/add` + `=param=nilai` | `PUT /rest/path` |
| Hapus | `remove 0` | `/remove` + `=.id=*1` | `DELETE /rest/path/*1` |
| Hentikan command berjalan | Ctrl-C | `/cancel` (`=tag=`) | Batas 60 detik; tidak ada cancel |
| Cari item untuk `set`/`remove` | `[find where ...]` | `print` dengan `.proplist=.id` dan query, lalu loop | `GET` dengan filter, lalu `PATCH`/`DELETE` per ID |

Baris dengan label "rakitan" disusun dari aturan terdokumentasi dan belum diuji.

---

## 15. Jebakan dan praktik yang aman

1. **Tanpa spasi di sekitar `=`.** `gateway = 3.3.3.3` adalah error. [D2]
2. **Nomor item tidak stabil.** Gunakan nama, `[find ...]` di CLI, atau `=.id=` di API/REST. [D3][D5]
3. **Parameter `print` berbeda per menu.** Cek F1 atau `/console/inspect` sebelum menulis script yang bergantung pada `stats`, `terse`, `value-list`, dan sebagainya. [F1][D7]
4. **Nama variabel lokal jangan sama dengan nama properti** di dalam `where`. [D2]
5. **Regex hanya ada di CLI.** API biner dan REST tidak mendukung `~`. [D5]
6. **Selalu pakai `.proplist` di API/REST.** Tanpanya semua properti, termasuk yang mahal, ikut dikirim. [D5]
7. **`set` di API memerlukan ID.** Ambil ID lewat `print` dan query dulu. [D5]
8. **REST tidak cocok untuk streaming.** Pakai `once`, `count`, `duration`, atau poll berkala; ingat batas 60 detik dan antrean per user. [D6]
9. **Konversi tipe di klien.** Balasan REST berupa string. Laju trafik dalam bit/detik, ukuran dalam byte, `uptime` berformat `1d7h44m20s`. [D6]
10. **Jangan buka `*` sebagai `%2A`** di URL REST. [D6]
11. **Gunakan user baca-saja** (`read,api,rest-api`) untuk dashboard dan pemantauan; jangan berikan kredensial `admin` ke aplikasi atau asisten. [D6]
12. **Jangan memakai singkatan command di script.** Tulis nama lengkap. (Saran praktik.)
13. **`print as-value` di console** mungkin tidak menampilkan apa pun bila tidak dibungkus `:put [ ... ]`. [C1]
14. **Gunakan Safe Mode** (Ctrl-X/F4) saat mengubah konfigurasi lewat sesi yang bisa terputus. [D3]

---

## 16. Hal yang belum terkonfirmasi di dokumentasi

Bagian ini sengaja dibuat agar tidak ada klaim yang terlihat pasti padahal tidak:

- **Urutan argumen di CLI.** Tidak ada pernyataan resmi bahwa `where` harus paling akhir atau bahwa switch dan named boleh dicampur bebas. Yang jelas tertulis hanya: urutan attribute word di API tidak penting, urutan query word penting. [D5]
- **`stats` di luar `/interface`.** Dokumentasi resmi menjelaskan `stats` dan `stats-detail` untuk `/interface`. Keluaran F1 di forum menunjukkan `stats` juga ada di `registration-table`. Cakupan lengkapnya tergantung menu; cek F1. [D9][F1]
- **`follow-strict`.** Terdaftar di referensi tanpa penjelasan. [D4][D10]
- **Satuan dan batas `interval`.** Referensi otomatis menampilkan rentang numerik, tetapi tafsirannya tidak dijelaskan. Dokumentasi CLI Reference hanya menyebut tipe `time`. [D1][D10]
- **`once` untuk `monitor-traffic` di CLI.** Terdokumentasi di halaman REST, tidak di halaman referensi `monitor-traffic`. [D6][D8]
- **`stats`, `interval`, `terse` lewat API biner dan REST.** Halaman API hanya menjelaskan `.proplist`, query, `=detail=`, dan `listen`. Bukti penggunaan `stats` dan `interval` datang dari kode komunitas (F4, F5), bukan dari halaman resmi.
- **`print follow` di API.** Hanya dari forum (F4).
- **Halaman API versi 7.26.** Yang dibaca adalah 7.25. Pembaruan kecil mungkin ada di 7.26.
- **Referensi otomatis help.mikrotik.com (D10).** Halaman penuhnya menolak akses otomatis; isi tabel parameternya hanya dibaca dari potongan hasil pencarian.
- **Semua contoh bertanda (rakitan)** belum dijalankan di perangkat.

---

## 17. Daftar referensi

**Dokumentasi resmi**

- [D1] CLI Reference (tipe argumen, modifier `print`): https://manual.mikrotik.com/docs/cli-reference/
- [D2] Scripting: https://manual.mikrotik.com/docs/developer-guides/scripting/
- [D3] Console (arsip beku): https://help.mikrotik.com/docs/spaces/ROS/pages/8978498/Console
- [D4] Command Line Interface (arsip beku): https://help.mikrotik.com/docs/spaces/ROS/pages/328134/Command+Line+Interface
- [D5] API (versi 7.25): https://manual.mikrotik.com/docs/7.25/developer-guides/api/
- [D6] REST API: https://manual.mikrotik.com/docs/developer-guides/rest-api
- [D7] `console/inspect`: https://manual.mikrotik.com/docs/cli-reference/console/inspect
- [D8] `interface` dan `interface/monitor-traffic`: https://manual.mikrotik.com/docs/cli-reference/interface/ , https://manual.mikrotik.com/docs/cli-reference/interface/monitor-traffic
- [D9] Interface stats and monitor-traffic: https://manual.mikrotik.com/docs/diagnostics-monitoring-and-troubleshooting/interface-stats-and-monitor-traffic/
- [D10] Referensi command otomatis (arsip): https://help.mikrotik.com/docs/x/eADWAg , https://help.mikrotik.com/docs/x/gwDWAg
- Introduction (tip file `.md` per halaman): https://manual.mikrotik.com/docs/introduction

**Forum MikroTik (komunitas)**

- [F1] Missing SNMP Entries On hAP ac^3: https://forum.mikrotik.com/t/missing-snmp-entries-on-hap-ac-3/178748
- [F2] New command in RouterOs 7: https://forum.mikrotik.com/t/new-command-in-routeros-7/169237
- [F3] RouterOS LSP (syntax checking dan completion): https://forum.mikrotik.com/t/routeros-lsp-for-better-syntax-checking-command-completion-in-editors-like-vscode-neovim/184067
- [F4] Pear2 broken signature: https://forum.mikrotik.com/t/pear2-broken-signature/94032
- [F5] Insert PPPoe client bandwidth into mysql: https://forum.mikrotik.com/t/insert-pppoe-client-bandwidth-into-mysql/59856

**Pihak ketiga**

- [C1] tik4net wiki, CLI tips and tricks: https://github.com/danikf/tik4net/wiki/CLI-tips-and-tricks