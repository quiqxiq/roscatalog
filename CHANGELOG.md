# Changelog

Format mengikuti [Keep a Changelog](https://keepachangelog.com/) dan
[Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added
- `ros.Events()` dan `DeliveryEvents`: `/path/listen` dipisahkan dari
  `DeliveryStream`. Listen hanya mengirim event perubahan, tanpa snapshot awal;
  diam lama tidak lagi dianggap channel rusak dan tidak memicu reopen.
- `DeliveryEvents` ditambahkan ke `Resource.Deliveries` untuk view listen
  (plus `poll` bila view `list` tersedia).
- Test: listen tidak dibuka ulang saat sepi; `PlanFor` memilih events untuk
  view listen; `Events()` ditolak untuk view non-listen.

### Changed
- `ros.Stream()` pada view listen dinormalkan menjadi `DeliveryEvents` untuk
  kompatibilitas. `Build` menolak `Events` pada view yang bukan listen.
- `MikroTik.PlanFor`, `WatchAuto`, dan cache menerima `DeliveryEvents`.
- Dokumentasi (`README.md`, `USAGE.md`, `ROADMAP.md`) menjelaskan perbedaan
  stream vs events.

### Fixed
- Watchdog stream tidak lagi membuka ulang channel `listen` hanya karena tidak
  ada event selama `StreamStale`.
- Refcount subscription dilepas saat upstream stream/events berakhir sendiri
  lewat `fireClose`, sehingga `Manager.Subscribe` tidak bocor.

## [0.1.0] - 2026-11-21

Rilis pertama yang diberi tag.

### Added
- Katalog RouterOS v7.9+ dengan view, aksi, dan flag `print`.
- Paket `fields/`: struct per menu, plus helper opsional `FieldNames`,
  `Decode`, `DecodeAll`, `ArgsOf`, dan `Params.Proplist`.
- Mode raw (tanpa `.proplist=`), `Only` untuk subset, dan `Proplist` untuk
  override struct buatan pemakai.
- Default realtime stream; `Config.Poll` untuk polling dari Go.
- Satu koneksi async per router (`DualConn` sebagai escape hatch).
- Limiter command in-flight per router dengan statistik.
- Satu scheduler poll per sesi; union proplist per menu dengan penyaringan
  field per subscriber.
- Lifecycle `Manager.Acquire/Release/Retain/Drop` dengan grace.
- Watchdog stream dan auth backoff.
- Dataset `mikrotik_inspect/`, `SOURCE.md`, dan generator struct/command.

### Verified
- `go test ./...`, `go test -race ./rosrun`, dan `go vet ./...`.
- Uji read-only ke CHR GNS3: connect, stream+command, poll, union proplist,
  listen.
