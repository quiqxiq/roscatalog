# Changelog

Format mengikuti [Keep a Changelog](https://keepachangelog.com/) dan
[Semantic Versioning](https://semver.org/).

## [Unreleased]

## [0.2.0] - 2026-11-21

Baseline **v7.9+**. Penghapusan `Version` adalah breaking change: Pemakai
v6 yang butuh field lama pakai resource/view sendiri atau `Raw`/`Params.Proplist`.

### Removed

- `Version` dan `AnyVersion`; konstanta `V6`/`V7`; field `View.Version`.
- Versi-parameter di permukaan library:
  - `Resource.View(name, ver)` → `Resource.View(name)`,
  - `Resource.Build(view, ver, plan, params)` → `Resource.Build(view, plan, params)`,
  - `Resource.Deliveries(view, ver)` → `Resource.Deliveries(view)`,
  - `MikroTik.PlanFor(res, view, ver)` → `MikroTik.PlanFor(res, view)`,
  - `MikroTik.WatchAuto` / `Manager.WatchAutoAll` tanpa parameter versi,
  - flag `-ver` di `cmd/rosverify` dan `cmd/roscheck`.
- View versi di katalog: `ip-route` (field v7: `routing-table`, bukan `routing-mark`
  v6) dan `system-health` (baris `name/value/type`) hanya v7 sekarang;
  katalog `zz_catalog_enrich.go` dibuang.
- `ErrActionNotFound` (was dead; `Resource.Action` kembali `(Action, bool)`).

### Added

- `Catalog`: Pemakai bisa membuat catalog sendiri dan menambahkan menu/command
  tanpa mengubah global. API: `NewCatalog(base ...*Resource)`, `Builtin()`,
  `Add`, `Replace`, `Extend`, `Find` (`ErrResourceNotFound`), `Lookup`,
  `All`, `ByGroup`, `Len`, `Validate`.
- Error sentinel: `ErrInvalidResource`, `ErrDuplicateResource`,
  `ErrViewNotFound`, `ErrUnknownFlag`, `ErrUnknownColumn`, `ErrMissingArgument`,
  `ErrSelectorRequired`, `ErrDeliveryUnsupported`. Error dari `Build`/`Find`
  di-wrap dengan `%w` sehingga pemakai bisa `errors.Is`.
- `Catalog.Validate()` / `Resource.Validate()`: kesalahan definisi sebelum
  menyentuh router (ID kosong/duplikat, path tidak diawali `/`, selector
  hilang, dll).
- Uji: Catalog (find/add/duplikat/invalid/replace/extend), sentinel errors,
  resource kustom CRUD/print (`catalog_api_test.go`).

### Added (dari 0.1.0)

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
  stream vs events; serta contoh `Build`/`WatchAuto` tanpa `AnyVersion`.

### Fixed

- Watchdog stream tidak lagi membuka ulang channel `listen` hanya karena tidak
  ada event selama `StreamStale`.
- Refcount subscription dilepas saat upstream stream/events berakhir sendiri
  lewat `fireClose`, sehingga `Manager.Subscribe` tidak bocor.

### Verified

- `go test ./...`, `go test -race ./rosrun`, dan `go vet ./...`.
- Uji read-only ke CHR GNS3: connect, stream+command, poll, union proplist,
  listen, listen-removal.
