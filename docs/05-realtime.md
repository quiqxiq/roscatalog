# 05 — Realtime vs Poll

`Plan` menentukan koneksi lazim: selama stream/poll hidup, subscriber bisa diputus
reconnect tanpa mengulang Dst.

```go
type Plan struct {
	Delivery Delivery      // once, stream, events, poll
	Every    time.Duration // interval stream (1s..5s) atau jeda poll (5s..30s)
	Diff     bool          // poll: hasilkan added/changed/removed
	Seed     bool          // listen: ambil snapshot awal lewat print
}
```

## 5.1 Delivery modes

| Mode | Fungsi |
|---|---|
| `Once()` (DeliveryOnce) | one-shot: Exec `req.Sentence` + first emit |
| `Stream()` (DeliveryStream) | stream native (interval) — **default baca; realtime-first** |
| `Events()` (DeliveryEvents) | listen — perubahan item saja (no initial snapshot) |
| `Poll(d)` (DeliveryPoll) | poller internal (Go polling per Once) — digunakan bila stream ditolak |
| `StreamEvery(d)` | alias `Plan{Delivery: DeliveryStream, Every: d}` |

Contoh:

```go
// live - interval 1s (default baca: stream)
req, _ := res.Build("live", ros.StreamEvery(time.Second), ros.Params{})

// hanya ada perubahan (no snapshot) — Events (listen)
req, _ := res.Build("events", ros.Events(), ros.Params{})

// poll dari Go — bila router menolak listen (Build sudah memakai diff)
req, _ := res.Build("events", ros.Poll(5*time.Second), ros.Params{})

// seed lalu poll: listen + seed
req, _ := res.Build("events", ros.Events().WithSeed(), ros.Params{})

// one-shot
req, _ := res.Build("list", ros.Once(), ros.Params{})
```

## 5.2 Config: realtime-first + Poll opt-in

```go
sess, _ := rosrun.Open(ctx, rosrun.Config{
	Address:   "192.168.88.1:8728",
	Username:  "admin",
	Password:  "r00t",
	PollEvery: 5 * time.Second, // polling dari Go hanya bila user minta
})
```

- Default: `Stream` native (realtime-first).
- `Config.Poll = true`: seluruh koneksi polling dari Go (misal monitor);
  **`Status` tetap native stream** (status monitor `/ping` selalu stream,
  lihat [Status](07-status.md)).
- `Config.Stream` (field deprecated) tidak lagi berarti poll; gunakan
  `Config.Poll`. `ros.Stream()` pada view listen tetap dinormalkan menjadi
  `Events` untuk kompatibilitas.

## 5.3 Fallback otomatis (`WatchAuto`)

Bila stream ditolak router sebelum ada data:

```go
ch, _ := sess.WatchAuto(ctx, res, "events", ros.Plan{Every: 5 * time.Second}, ros.Params{})
```

## 5.4 Difference lokal (`Params.Query`, `Diff`, `Seed`)

```go
req, _ := res.Build("list", ros.Poll(5*time.Second).WithDiff(), ros.Params{
	Query: []ros.Query{{Prop: "disabled", Value: "yes"}},
})
```

Pergi ke: [Connection Layer](06-connection-layer.md) | [Status](07-status.md).
