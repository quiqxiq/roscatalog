# Wiki: roscatalog

roscatalog adalah library **katalog + pembangun sentence** untuk command RouterOS
API (milik-style Mikhmon): mengandung daftar menu/resource, view baca
(oneshot/interval/follow/listen/native), aksi tulis, dan konversi ke `Plan`
stream/poll/execute.

Ada dua module dalam project ini:

| Module | Fungsi |
|---|---|
| `example.com/roscatalog` | katalog + `Request` builder + `fields/` struct + sentinel error |
| `example.com/roscatalog/rosrun` | adapter ke engine transport (`goros/`), transport asinkron, watchdog, Manager |

Daftar halaman (baca dari awal):

1. [Mulai](01-getting-started.md) — install, dial, daftar pertama
2. [Catalog](02-catalog.md) — API catalog, Builtin, custom catalogs
3. [Custom Resources](03-custom-resources.md) — Table/Single/Tool + aksi tulis
4. [Views](04-views.md) — mode view, proplist, Only/Raw/Presets, Listen-no-proplist
5. [Realtime vs Poll](05-realtime.md) — `Plan`, `Delivery`, `Poll`, `Stream`, `Events`
6. [Connection Layer](06-connection-layer.md) — `Runner`, `Watch`, `Exec`, `Manager`
7. [Status](07-status.md) — `Status`, `WatchStatus`, state machine, live result di CHR
8. [Fields](08-fields.md) — `fields/` structs + helper
9. [Errors](09-errors.md) — sentinel errors + `errors.Is`
10. [Testing](10-testing.md) — unit + live tests, roscheck/rosverify
11. [Migrasi MikroDash](11-mikrodash.md) — jalan menggunakan roscatalog di MikroDash

Contoh code yang bisa dijalankan: [docs/examples/](examples).

Baseline: RouterOS **v7.9+** (tidak ada `Version`; `Raw`/`Proplist`/`Only` untuk
field v6-style bisa pakai sebagai bridge).
