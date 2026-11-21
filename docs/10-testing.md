# 10 — Testing

## 10.1 Unit

```bash
go test ./...                # semua package
go test -race ./rosrun       # racing check utama
go vet ./...
go build ./...
```

Katalog assertions ada di `catalog_test.go` + `catalog_api_test.go`
(Validasi definition sudah ada: tambah test dengan `res.Validate()`, atau
duplikat resource → `ErrDuplicateResource`).

## 10.2 Live (ke CHR-GNS3 192.168.122.189:8728 — read-only constraints)

```bash
# 1. Connect + stream + exec + poll (sukses)
ROUTEROS_LIVE_ADDR=192.168.122.189:8728 ROUTEROS_LIVE_USER=admin ROUTEROS_LIVE_PASS=r00t \
  go test ./rosrun -run 'TestLiveSession' -v

# 2. Listen removal (telat del + penghapusan, realtime)
ROUTEROS_LIVE_ADDR=192.168.122.189:8728 ROUTEROS_LIVE_USER=admin ROUTEROS_LIVE_PASS=r00t \
  go test ./rosrun -run 'TestLiveListenRemoval' -v
```

`TestLiveSession*` menulis? NO — semua read-only. Status CHR (live):
`sess.Status().Connected == true`, `State == StatusWarning` (karena CHR
di 8728 tidak pakai TLS) — lihat [07-status](07-status.md).

## 10.3 roscheck / rosverify (verifikasi drift)

`cmd/rosverify` — setiap resource+view cek: path, verb, argumen native,
field `print` via `/console/inspect` (coba CHR pakai API user yang punya
permission semacam `api`). Status: `ok`/`warn`/`drift`/`unknown`/`error`.

```bash
# semua resource
go run ./cmd/rosverify -addr 192.168.122.189:8728 -user admin -pass r00t -json -strict

# per resource
go run ./cmd/rosverify -addr 192.168.122.189:8728 -user admin -pass r00t \
  -resource hotspot-active -json

# discovery (opt-in): path/verb completion untuk deploy
go run ./cmd/rosverify -addr 192.168.122.189:8728 -user admin -pass r00t \
  -discover ip,hotspot,active -completion
```

## 10.4 Golden schema di CI

Simpan hasil `rosverify -json` per CHR-versi sebagai golden file (lihat
`mikrotik_inspect/inspect.json` sebagai contoh), bandingkan secara manual
saat CHR update. Lihat juga `TestVerifiedHaveViews` (catalog test).

Pergi ke: [12-end](README.md).
