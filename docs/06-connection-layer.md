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
sess, _ := m.Add(ctx, rosrun.Config{ID: "r1", Address: "192.168.88.1:8728", Username: "admin", Password: "secret"})
id := sess.Stats().ID

_, _ = m.Acquire(ctx, id)                      // pin (refcount); membuka sesi dari Config terdaftar bila perlu
sub, err := m.Subscribe(ctx, id, req)         // stream/events per subscription
m.Release(id)
_ = m.Retain(id, "bg-job")                    // tahan sesi tanpa viewer
m.Drop(id, "bg-job")                          // lepas hold; router di-close bila refcount 0
_, _ = m.WatchAutoAll(ctx, res, view, plan, params) // Acquire+Release otomatis per router
st := m.Statuses()                            // []Status
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
	MaxInFlight int           // <= DefaultMaxInFlight (8) command/router
	Grace       time.Duration // jeda sebelum sesi idle ditutup, DefaultGrace = 30s
	StreamStale time.Duration // watchdog: idle stale, DefaultStreamStale = 10s
	StreamCheck time.Duration // interval check, DefaultStreamCheck = 5s
}
```

- **Stats()** — snapshot per-router: `ID, Address, Open, Epoch, Reconnects, Down,
  LastError, Uptime, InFlight, OpenStreams, Commands, BusiestMenus []MenuCount`.
- Watchdog stream: idle > 10s → re-open stream; stuck → check every 5s.
- Auth reconnect backoff: 30s hingga 5m.

Pergi ke: [Status](07-status.md) | [Realtime](05-realtime.md).
