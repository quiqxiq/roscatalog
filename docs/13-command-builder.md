# Universal Command Builder: `ros.Path` / `ros.Menu`

Di RouterOS, terdapat dua cara berinteraksi dengan API:
1. **Catalog (`ros.Lookup`)**: Pendekatan deklaratif terstruktur untuk resource umum yang terdaftar (seperti `hotspot-active`, `hotspot-user`, `interface`, dll.).
2. **Universal Command Builder (`ros.Path` / `ros.Menu`)**: Pendekatan terstruktur bebas untuk membentuk command RouterOS apa pun — baik baca (`print`), tulis (`set`, `add`, `remove`, `enable`, `disable`), maupun aksi/tools (`monitor-traffic`, `ping`, `check-for-updates`) — sesuai dengan spesifikasi RouterOS CLI dan Binary API sentence protocol.

---

## 1. Hubungan Struktur CLI dan Binary API

Berdasarkan [MIKROTIK-ROUTEROS-COMMAND-STRUCTURE.md](../MIKROTIK-ROUTEROS-COMMAND-STRUCTURE.md), kalimat command RouterOS CLI tersusun atas:

```text
[menu/path] command [uparam] [switch ...] [nama=nilai ...] [where ...]
```

Pada Binary API protocol (`goros/v4`), struktur tersebut diterjemahkan secara presisi:
- **Command word**: `/{menu/path}/{command}` (contoh: `/interface/print`, `/ip/hotspot/user/set`, `/ping`).
- **Switches / Flags tanpa nilai**: `=switch=` (contoh: `=stats=`, `=detail=`, `=count-only=`).
- **Positional Selector (Numbers)**: `=numbers={selector}` (contoh: `=numbers=ether1`, `=numbers=*1,*2`).
- **Proplist**: ` =.proplist=kolom1,kolom2` (diawali titik `.`).
- **Named Arguments**: `={nama}={nilai}`.
- **Query Filter**: `?{nama}={nilai}`, `?>{nama}={nilai}`, `?{nama}`, `?-{nama}`.

---

## 2. Inisialisasi: `ros.Path` & `ros.Menu`

`ros.Path(path)` (atau aliasnya `ros.Menu(path)`) menerima path dalam format slash (`/ip/hotspot/user`), spasi CLI (`ip hotspot user`), maupun root (`/`):

```go
p1 := ros.Path("/ip/hotspot/user")
p2 := ros.Menu("ip firewall filter")
p3 := ros.Path("/")
```

---

## 3. Print Command Builder: `.Print(...)`

Memanggil `.Print(switches ...string)` pada path mengembalikan `*PrintCmdBuilder`.

### Switches & Flags
Switches diterjemahkan menjadi kata `=switch=`:
```go
// Menghasilkan sentence: /interface/print =stats= =detail=
cmd := ros.Path("/interface").Print("stats", "detail")
```

Switch bawaan yang didukung:
- `"stats"`, `"detail"`, `"terse"`, `"brief"`, `"count-only"`, `"value-list"`, `"without-paging"`, dll.

### Kolom Proplist: `.Proplist(...)`
Memilih kolom yang diambil langsung di router (` =.proplist=...`):
```go
cmd := ros.Path("/interface").Print("stats", "detail").
    Proplist(".id", "name", "tx-byte", "rx-byte")
```

### Realtime Streaming: `.Interval(...)` & `.Follow()`
Mengatur stream realtime dengan interval atau perubahan baris:
```go
// Realtime streaming setiap 1 detik
cmdLive := ros.Path("/interface").Print("stats", "detail").
    Interval(time.Second).
    Proplist("tx-byte", "rx-byte")

// Event streaming follow
cmdFollow := ros.Path("/ip/hotspot/active").Print().
    Follow()
```

### Server-side Queries & Where
Mendukung filter di sisi router:
```go
// WhereEq: ?profile=1jam
// WhereHas: ?comment
// WhereNotHas: ?-disabled
cmd := ros.Path("/ip/hotspot/user").Print().
    WhereEq("profile", "1jam").
    WhereHas("comment")
```

---

## 4. Aksi Tulis: `Set`, `Add`, `Remove`, `Enable`, `Disable`

Setiap method mutasi menerima argumen fleksibel dalam **satu panggilan**, menghindari pengulangan method chaining (`.Set().Set().Set()`):

Argumen dapat berupa:
1. `map[string]any` atau `map[string]string`
2. Struct bertag `ros:"..."` (dari paket `fields` atau struct kustom)
3. Pasangan variadic key-value (`"key", val, "key2", val2`)

### Menambah Item (`.Add(...)`)
```go
// Menggunakan map[string]any
cmdAdd := ros.Path("/ip/hotspot/user").Add(map[string]any{
    "name":     "voucher01",
    "password": "123",
    "profile":  "1jam",
    "comment":  "vc-nov-24",
})

// Atau menggunakan struct dari package fields
row := fields.IPHotspotUserAdd{
    Name:     "voucher02",
    Password: "123",
    Profile:  "1jam",
}
cmdAddStruct := ros.Path("/ip/hotspot/user").Add(row)
```

### Mengubah Item (`.Set(selector, inputs...)`)
```go
// Menggunakan map
cmdSet := ros.Path("/ip/hotspot/user").Set("voucher01", map[string]any{
    "profile":      "2jam",
    "limit-uptime": "2h",
})

// Menggunakan variadic key-value
cmdSet2 := ros.Path("/interface/ethernet").Set("ether1",
    "mtu", 1500,
    "comment", "WAN Primary",
)
```

### Kontrol Item (`Remove`, `Enable`, `Disable`, `Comment`, `Reset`)
```go
p := ros.Path("/ip/hotspot/user")

cmdRemove  := p.Remove("voucher01")              // /ip/hotspot/user/remove =numbers=voucher01
cmdEnable  := p.Enable("voucher01")              // /ip/hotspot/user/enable =numbers=voucher01
cmdDisable := p.Disable("voucher01")             // /ip/hotspot/user/disable =numbers=voucher01
cmdComment := p.Comment("voucher01", "expired")  // /ip/hotspot/user/comment =numbers=voucher01 =comment=expired
cmdReset   := ros.Path("/system/history").Reset("*1") // /system/history/reset =numbers=*1
```

---

## 5. Child Commands & Tools: `.Do(command, inputs...)`

Untuk command non-CRUD atau tool di submenu mana pun:

```go
// Menjalankan ping
cmdPing := ros.Path("/").Do("ping", map[string]any{
    "address": "8.8.8.8",
    "count":   4,
})

// Monitor traffic interface
cmdTraffic := ros.Path("/interface").Do("monitor-traffic", map[string]any{
    "interface": "ether1",
    "once":      "",
})

// Reboot router
cmdReboot := ros.Path("/system").Do("reboot")
```

---

## 6. Integrasi dengan `rosrun`

Semua builder mengimplementasikan antarmuka `ros.Requestable`:
```go
type Requestable interface {
    Request() Request
    Build() (Request, error)
}
```

Paket `rosrun` menyediakan method langsung tanpa perlu konversi manual:

### Eksekusi di MikroTik (`*rosrun.MikroTik`)
```go
// 1. One-shot query (baca)
rows, err := sess.OnceCmd(ctx, ros.Path("/system/resource").Print())

// 2. Realtime watch / interval stream
ch, err := sess.WatchCmd(ctx, ros.Path("/interface").Print("stats", "detail").Interval(time.Second))
for ev := range ch {
    // Tangani ev.Kind (snapshot/row/error)
}

// 3. Eksekusi aksi tulis (mengembalikan =ret= bila ada)
id, err := sess.ExecCmd(ctx, ros.Path("/ip/hotspot/user").Add(map[string]any{
    "name": "u1", "password": "p1",
}))
```

### Eksekusi Multi-Router (`*rosrun.Manager`)
```go
// Membuka stream ke semua router sekaligus
ch, err := mgr.WatchAllCmd(ctx, ros.Path("/interface").Print("stats").Interval(time.Second))

// Mendaftarkan shared-subscription ke semua router
subs, err := mgr.SubscribeAllCmd(ctx, ros.Path("/system/resource").Print().Interval(2*time.Second))
```
