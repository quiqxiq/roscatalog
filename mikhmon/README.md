# Paket `mikhmon`

Overlay katalog untuk **Mikhmon v4**. Paket ini tidak menyalin definisi
path/verb RouterOS; sumber kebenaran tetap katalog root `roscatalog`.

Isinya:

- `Lookup`/`Resources`/`IDs` — subset resource yang dipakai Mikhmon.
- `OnLoginScript` — builder properti `on-login` satu baris (mode `0`, `rem`,
  `ntf`, `remc`, `ntfc`, lock MAC/server, laporan penjualan).
- `ExpireMonitorScript` — on-event scheduler `Mikhmon-Expire-Monitor`.
- `Sale` — format nama/source/owner laporan penjualan.
- `QuickPrintPreset` — format source preset QuickPrint.
- `Client` — workflow add profile, add user/voucher, set on-login, install
  Expire Monitor, add sale, simpan QuickPrint, kick user, dan hapus user
  kedaluwarsa.

`Client` hanya bergantung pada interface `Executor`:

```go
type Executor interface {
    Once(ctx context.Context, req ros.Request) ([]map[string]string, error)
    Count(ctx context.Context, req ros.Request) (int, error)
    Exec(ctx context.Context, sentence []string) (string, error)
}
```

`*rosrun.Runner` memenuhi interface itu, sehingga paket ini tidak mengimpor
`goros` secara langsung.

## Contoh

```go
run := rosrun.New(c)
client := mikhmon.New(run)

onLogin, err := mikhmon.OnLoginScript(mikhmon.OnLoginOptions{
    Mode:         mikhmon.ExpiryRemoveRecord,
    Price:        "5000",
    Validity:     "1d",
    SellingPrice: "7000",
    Profile:      "1jam",
})
if err != nil {
    log.Fatal(err)
}

if _, err := client.AddProfile(ctx, mikhmon.Profile{
    Name: "1jam", RateLimit: "1M/1M", SharedUsers: "1",
    OnLogin: onLogin,
}); err != nil {
    log.Fatal(err)
}

name := "vc001"
if _, err := client.AddUser(ctx, mikhmon.User{
    Name: name, Password: name, Profile: "1jam", Server: "all",
    Comment: mikhmon.MikhmonComment(name, name, "batch-1"),
    LimitUptime: "1d",
}); err != nil {
    log.Fatal(err)
}

if err := client.InstallExpireMonitor(ctx); err != nil {
    log.Fatal(err)
}

if _, err := client.AddSale(ctx, mikhmon.Sale{
    Date: "nov/18/2026", Time: "10:00:00", User: name,
    Price: "5000", Address: "198.51.100.10", MAC: "AA:BB:CC:DD:EE:FF",
    Validity: "1d", Profile: "1jam", Comment: "vc-123-11.18.26",
}); err != nil {
    log.Fatal(err)
}
```

## Uji langsung

```bash
go run ./cmd/mikhmon-live -addr 192.168.88.1:8728 -user admin -pass secret
go run ./cmd/mikhmon-live -addr 192.168.88.1:8728 -user admin -pass secret -json
```

Command membuat data berprefix `zzmik`, menjalankan on-login secara manual
dengan variabel `$user`, `$address`, dan `$"mac-address"`, lalu memverifikasi
laporan penjualan lewat query `source`/`owner`. Data uji dihapus lagi; pakai
`-keep` untuk mempertahankannya.
