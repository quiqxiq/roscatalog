# 07 — Realtime Status (connected / warning / error)

`Runner`/`MikroTik` memiliki *monitor status realtime-first*: satu stream
native `/ping` ke target (default 8.8.8.8) — selalu **native stream** (bukan
poll; `Config.Poll` tidak mempengaruhi — lihat
`TestStatusUsesNativePingRegardlessOfStream`). Tanpa `PingInterval`, tidak ada
`=interval=` yang dikirim dan RouterOS memakai cadence bawaannya (1 detik).
`ReplyTimeout` (default 15 detik) menentukan kapan tidak adanya balasan ping
dianggap error; nilainya tidak diturunkan dari interval ping.

```go
st := sess.Status()

type Status struct {
	ID, Address, Target string
	State   StatusState     // StatusConnected / StatusWarning / StatusError
	Connected bool          // koneksi hidup + ping baru saja
	Warnings []string       // alasan warning (TLS, packet loss, status ping)
	Err      error          // di-populate bila State=error
	Ping    PingStatus      // cumulative data ping
	Down bool               // koneksi sesi sedang putus
	Epoch, Reconnects uint64 / int64
	Uptime  time.Duration
	Updated time.Time
}
```

`PingStatus` (dari status.go):

```go
type PingStatus struct {
	Seq, Sent, Received    int
	PacketLoss             float64       // persen
	MinRTT, AvgRTT, MaxRTT, LastRTT time.Duration
	LastStatus             string        // "success" atau status router
	Updated                time.Time
}
```

## 7.1 WatchStatus (realtime)

```go
ctx, cancel := context.WithCancel(ctx)
ch := sess.WatchStatus(ctx)
for st := range ch {
	if st.Connected {
		fmt.Println(st.State, st.Ping.PacketLoss, st.Ping.LastRTT)
	}
}
// kanal tutup saat sesi di-close atau ctx batal.
```

## 7.2 Live result (CHR 7.24.6)

Dari CHR 192.168.122.189 (live probe tadi):

```
WATCH update 1: state=warning connected=true warnings=[TLS disabled: API traffic is plaintext] ping={Seq:14 Sent:15 Received:15 PacketLoss:0 MinRTT:29.8ms AvgRTT:37.8ms LastStatus:success ...}
WATCH update 2: state=warning connected=true ping={Seq:15 Sent:16 Received:16 ...}
WATCH update 3: state=warning connected=true ping={Seq:16 Sent:17 Received:17 ...}
```

- `connected=true`, ping berjalan (RTT ~30ms, no loss), update tiap detik.
- `warning` = karena CHR di port 8728 tidak pakai TLS (warning biasa — router tetap terhubung).
- **Err = `ErrStatusTimeout`** bila tidak ada reply ke 8.8.8.8 > `Status.ReplyTimeout` (default 15s);
  bila **koneksi down** → `ErrConnectionDown`; bila sesi tertutup → `ErrClosed`.
- `StatusTimeoutBecomesError` unit test PASS.

## 7.3 Transition (dari error)

Kalau CHR mengganti IP:

1. stream dicolokkan kembali secara otomatis (backoff 1s→10s).
2. ping: `PingStatus` frozen, `Warnings += "status ping: ..."`; kalau tidak ada balasan 15s → `ErrStatusTimeout`.
3. `State: error` muncul realtime lewat `WatchStatus`.

Pergi ke: [Errors](09-errors.md) | [Connection](06-connection-layer.md).
