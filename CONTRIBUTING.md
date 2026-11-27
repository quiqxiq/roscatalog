# Contributing to `roscatalog`

## Prasyarat

- Go 1.22+
- Akses ke RouterOS 7.9+ (opsional, hanya untuk `//go:build live` dan `cmd/rosverify`)

## Alur Pengembangan

1. Jalankan formatter dan pemeriksaan statik sebelum membuat PR:
   ```bash
   gofmt -l .
   go vet ./...
   ```
2. Jalankan seluruh unit & integration test (termasuk race detector):
   ```bash
   go test ./... -race -count=1
   go build ./docs/examples/...
   ```
3. Untuk menguji ke perangkat RouterOS nyata (hati-hati: beberapa test menulis item sementara):
   ```bash
   ROUTEROS_LIVE_ADDR=192.168.88.1:8728 ROUTEROS_LIVE_USER=admin ROUTEROS_LIVE_PASS=secret \
     go test -tags live ./rosrun -v
   ```
4. Pastikan contoh dokumentasi dan `docs/examples/*` tetap dapat di-compile.
