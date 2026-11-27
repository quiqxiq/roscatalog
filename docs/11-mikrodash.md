# 11 — Migrasi ke MikroDash

`roscatalog` menutup cakupan command MikroDash dengan `rosrun` sebagai transport data.

Strategi migrasi bertahap:

1. **Fase A — test harness di `roscatalog`:** verifikasi schema via `cmd/rosverify` dan `go test -tags live ./rosrun`.
2. **Fase B — ganti 1 menu:** mulai dari menu yang field CHR-nya terverifikasi (mis. `interface`) via `Builtin().Find(...)` dan `res.Build(...)`.
3. **Fase C — generalized fallback:** ganti kolektor per-menu secara bertahap dengan `Params.Only`/`Proplist`/`Raw`.
4. **Fase D — generic engine:** gunakan `Action.Sentence` untuk operasi tulis.
5. **Fase E — separasi responsibilitas:** TUI/REST tetap di aplikasi klien, katalog & transport di `roscatalog`/`rosrun`.

Pergi ke: [README](README.md).
