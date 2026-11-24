# 06 — Connection Layer (rosrun)

`rosrun` adalah adapter ke engine transport (`goros/`). Dua pola utama:

- **Runner** (`New(c)`): satu koneksi asinkron per router, bawaan.
- **Manager** (`NewManager()`): banyak router; kontrol kuota/lifecycle.

## 6.1 Runner

```go
r := rosrun.New(c)                 // c = *routeros.Client (async mode)

rows, err := r.Once(ctx, req)      // one-shot: exec sentence, return all rows
sent, err := r.Exec(ctx, sentence) // exec, return !done.ret (string)
n, err := r.Count(ctx, req)        // exec count view
ch, err := r.Watch(ctx, req)       // stream/events chan
ch, err := r.WatchAuto(ctx, res, view, plan, params)  // auto-fallback to poll
```

Sessions bisa di-persist via `rosrun.Open(ctx, cfg)`.

## 6.2 Manager

```go
m := rosrun.NewManager()
sess, _ := m.Add(ctx, rosrun.Config{Address: "192.168.88.1:8728", Username: "admin", Password: "r00t"})

m.Acquire(ctx, sess.ID)     // pin (refcount); membuka sesi dari Config terdaftar bila perlu
sub, err := m.Subscribe(ctx, sess.ID, req)   // stream/events per subscription
m.Release(id)
m.Retain(id, reason)        // +1 ref, WatchAuto/WatchAll defers release
m.Drop(id, reason)          // -1 ref; router di-close bila refcount 0
m.WatchAutoAll(ctx, res, view, plan, params) // Acquire+Release otomatis per router
st := m.Statuses()          // []RouterStatus{ID, Status}
```

`Register(cfg)` hanya menyimpan konfigurasi; koneksi dibuka saat `Acquire`
pertama. Registrasi **tidak** hilang ketika sesi ditutup karena grace, jadi
`Acquire` berikutnya membuka sesi baru. `WatchAll`, `WatchAutoAll`,
`WatchStatuses`, dan `SubscribeAll` memakai `Acquire` selama stream hidup,
sehingga grace tidak bisa menutup sesi yang masih dipakai. `Configs()` berisi
semua registrasi; `IDs()`/`Routers()` hanya sesi yang sedang terbuka.

`Config.Grace` berlaku **per router**; router tanpa nilai memakai
`DefaultGrace` (30s).

Config:

```go
type Config struct {
	Address   string
	Username  string
	Password  string
	ID        string
	Timeout   time.Duration
	PollEvery time.Duration
	Poll      bool         // polling dari Go
	TLS, Insecure bool
	Status    StatusConfig // native /ping: target, timeout balasan, warn-packetloss
}
```

Setiap router memakai **tepat satu koneksi async**; command dan stream
dimultipleks lewat tag. Tidak ada opsi koneksi kedua. Feature switches
(`Config.Poll`, ...) default **off**.

## 6.3 Limiter & Watchdog

```go
type Config struct {
	MaxInFlight int          // ≤ DefaultMaxInFlight (8) command/router, Watch ↔ publish refcount
	Grace       time.Duration // antara drop koneksi, DefaultGrace = 30s
	StreamStale time.Duration // watchdog: idle stale, DefaultStreamStale = 10s
	StreamCheck time.Duration // interval check, DefaultStreamCheck = 5s
	Reconnect   time.Duration // bawaan
}
```

- **Stats()** — snapshot per-router: `ID, Address, Open, Epoch, Reconnects, Down,
  LastError, Uptime, InFlight, OpenStreams, Commands, BusiestMenus []MenuCount`.
- Watchdog stream: idle > 10s → re-open stream; stuck → check every 5s.
- Auth reconnect backoff (per router, memulai dari `Reconnect`).

Pergi ke: [Status](07-status.md) | [Realtime](05-realtime.md).
