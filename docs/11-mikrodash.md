# 11 — Migrasi ke MikroDash

roscatalog **telah menutup 117/117 command MikroDash** (verifikasi statik via
`tools/dump_catalog` + live CHR). Jalan tercepat: arahkan MikroDash ke catalog
untuk definisi/menu, gunakan `rosrun` sebagai transport data — lihat
[`MikroDash/ROADMAP.md`](../MikroDash/ROADMAP.md) dan
[`MikroDash/internal/server/resource.go`](../MikroDash/internal/server/resource.go)
sebagai calon home baru: `buildArgs → Action.Sentence`.

Strategi (per fase, pake dari small-step):

1. **Fase A — test harness di roscatalog.** Buat live test ke CHR dalam roscatalog
   (masukkan jadi test di `rosrun/live_*.go`): skenario 1 command+list+watch, 1 modify.
   Pack tolak error `errors.Is(err, ros.ErrViewNotFound)`.
2. **Fase B — ganti 1 menu.** Pilih menu dengan fields CHR terverifikasi (mis.
   `interface`) — define ulang di catalog via `Builtin().Find(...)`, lalu expose
   via API MikroDash `/api/<resource>` (hubungkan ke `res.Build`).
3. **Fase C — generalized fallback.** MikroDash `internal/collect/` per-menu, ganti
   satu-satu dengan `res.Build`; gunakan `Params.Only`/`Proplist`/`Raw`.
4. **Fase D — generic engine.** Single-form `internal/resource` generic terhasilkan dari
   `Action.Sentence` — lihat `MikroDash/internal/server/resource.go` pattern
   `GuardedBy/GuardedOn`.
5. **Fase E — separasi responsibilitas**: ownership halaman stay di MikroDash
   (TUI/REST), definisi command katalog, engine transport di roscatalog/rosrun.

Contoh integrasi ringan (satup menu, standar) — peta jalan fase A-E ada di
[`MikroDash/ROADMAP.md`](../MikroDash/ROADMAP.md).

Pergi ke: [README](README.md).
