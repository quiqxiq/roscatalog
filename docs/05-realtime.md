# 05 — Realtime vs Poll

`Plan` menentukan koneksi lazim: selama stream/poll hidup, subscriber bisa diputus
reconnect tanpa mengulang Dst.

```go
type Plan struct {
	Delivery Delivery // once, stream, events, poll
	Every    time.Duration // interval stream (MinPrintInterval..MaxPrintInterval)
	Interval time.Duration // override interval
	Diff     bool          // diff lokal &amp; bila gunakan diff bertubi (poll w/ diff)
	Seed     bool          // listen: ambil snapshot awal
	Fallback []string      // request tanpa .proplist= pada fallback
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

// poll dari Go — bila router menolak listen
req, _ := res.Build("events", ros.Poll(5*time.Second), ros.Params{})
req.Diff = true

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
- `Stream()` (deprecated): normalize jadi Events bila view listen.

## 5.3 Fallback otomatis (`WatchAuto`)

Bila stream ditolak router sebelum ada data:

```go
ch, _ := sess.WatchAuto(ctx, res, "events", ros.Plan{Every: 5 * time.Second}, ros.Params{})
```

## 5.4 Difference lokal (`Params.Query`, `Diff`, `Seed`)

```go
req, _ := res.Build("list", ros.Plan{Poll(5s).WithDiff()}, ros.Params{
	Query: []ros.Query{{Prop: "disabled", Value: "yes"}},
})
```

Pergi ke: [Connection Layer](06-connection-layer.md) | [Status](07-status.md).
