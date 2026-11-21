# Changelog

Format mengikuti [Keep a Changelog](https://keepachangelog.com/) dan
[Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added — menutup gap dengan MikroDash

CHR-192 (CHR 7.24.6) dan CHR-103 (CHR 7.25-era) dianggap sumber field CHR;
katalog ambil proplist dari CHR lalu tambah field MikroDash (masih tersedia lama,
mis. `wifi`). CHR-192/103 hilang menu atau 0-field: katalog tetap ditambahkan
(CHR baru/v7.25 akan lengkap).

- Resource baru (File juga — keangkatan nama resource mengikuti path CHR):
  - `app` — /app (statis); `app-settings` — /app/settings
  - `user-group` — /user/group; `user-settings` — /user/settings
  - `device-mode` — /system/device-mode; `license` — /system/license
  - `routerboard-settings` — /system/routerboard/settings
  - `disk` — /disk (CHR: proplist ok pada CHR baru)
  - `ipsec-policy` — /ip/ipsec/policy; `ipsec-active-peers` — /ip/ipsec/active-peers;
    `ipsec-installed-sa` — /ip/ipsec/installed-sa
  - `discovery-settings` — /ip/neighbor/discovery-settings
  - `routing-rule` — /routing/rule
  - `ipv6-route` — /ipv6/route
  - `ip-cloud` — /ip/cloud
  - `interface-bridge-vlan` — /interface/bridge/vlan
  - `interface-detect-internet-state` — /interface/detect-internet/state
  - `wireguard-peers` — /interface/wireguard/peers
  - `wifi-ap` — /interface/wifi (SRS, CHR: 1 warn: ssid, cap);
    `wifi-channel` — /interface/wifi/channel; `wifi-configuration`;
    `wifi-security`; `wifi-datapath`; `wifi-provisioning`; `wifi-radio`;
    `wifi-registration-table` (CHR ok);
    `wifi-legacy` — /interface/wireless (CHR lama, gak ada disini CHR);
    `wifi-legacy-security-profiles`; `wifi-legacy-registration-table`
  - `wifi-capsman` — /caps-man (CHR lama ini: hilang);
    `wifi-capsman-channel` — /caps-man/channel; `wifi-capsman-configuration` — /caps-man/configuration;
    `wifi-capsman-datapath` — /caps-man/datapath; `wifi-capsman-interface` — /caps-man/interface;
    `wifi-capsman-provisioning` — /caps-man/provisioning;
    `wifi-capsman-registration-table` — /caps-man/registration-table;
    `wifi-capsman-remote-cap` — /caps-man/remote-cap; `wifi-capsman-security` — /caps-man/security
- Aksi (writes) baru:
  - `file`: `add` (name, type, contents, permission, owner, group), `set`
  - `system-package-update`: `install`
  - `system-package`: `apply-changes`

### Verified against CHR (192.168.122.189:8728 + 103.139.193.128:1012)

- `go run ./tools/probe_print -addr 192.168.122.189:8728 -user admin -pass r00t -out /tmp/opencode/pf_fresh.json`
  dan `-addr 103.139.193.128:1012 -user fandi1 -pass 001 -out /tmp/opencode/pf_103.json`
  (raw print CHR). CHR-103 data lengkap: /app/settings (6 field), /system/license
  (features, nlevel, software-id); /system/device-mode CHR = CHR-103 (24 field,
  sama-sama 0-field pada CHR-192 lama).
- CHR-103 tidak lengkap `/console/inspect` schema (kebanyakan `unknown`) — CHR data
  paling lengkap dari raw print (`probe_print`). CHR-192 = CHR 7.24.6; CHR-103 = CHR
  7.25-era dengan field yang lebih banyak (`/system/license` 3 field, `/app/settings`
  6 field, `/system/device-mode` sama-sama 24 field).
- CHR-103 tidak punya menu capaian penuh: `/caps-man`, `/interface/wireless`,
  `/interface/wifi` (CHR-192 sama — CHR tidak mengirim field, menu hadir dengan 0 field);
  CHR-192/103: `/app`, `/disk` 0 field. Proplist wifi/capsman sebagai konsistensi katalog.

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
