# 08 — Package `fields/`

Struct opsional per menu — pemakai **tidak wajib** pakai (`Params{Only: [...]}`, `Proplist`, `Raw` tetap jalan). Langsung pakai struct saat decoded row perlu bentuk tiped/consumer.

Semuanya pakai tag `ros:<proplist-name>` sebagai sumber kolom.

```go
package fields

type QueueSimpleRow struct {
	ID        string `ros:".id"`
	Name      string `ros:"name"`
	Target    string `ros:"target"`
	MaxLimit  string `ros:"max-limit"`
	Rate      string `ros:"rate"`
	PacketRate string `ros:"packet-rate,omitempty"`
	Comment   string `ros:"comment,omitempty"`
}
```

Helper (di root package `roscatalog`, file `fields.go`):

- `FieldNames[T]() []string` — kolom `ros:` yang ada di struct, urut replease.
- `Decode[T any](row map[string]string) (T, error)`
- `DecodeAll[T any](rows []map[string]string) ([]T, error)`
- `ArgsOf[T any](v T) (map[string]string, error)` — dari struct balik ke map untuk `Action.Sentence`.

```go
rows, _ := sess.Once(ctx, req)
q, _ := ros.DecodeAll[fields.QueueSimpleRow](rows)
for _, r := range q { fmt.Println(r.ID, r.Name, r.MaxLimit) }

// Bila hanya perlu kolom targeted, pakai Params.Only atau Params.Proplist:
req, _ := res.Build("list", ros.Params{Only: ros.FieldNames[fields.QueueSimpleRow]()[:2]})
```

Rule pakai structs:

1. Sebelum pakai, jalankan `ros.FieldNames[T]()` / katalog — sekali per struct.
2. `FieldNames` diambil dari CHR (proplist lengkap di CHR 7.24+ bisa pakai); `Params.Only` pakai subset yang benar-benar dipakai.
3. Perbedaan CHR antara versi — gunakan `Params.Proplist` saat kolom luar katalog atau CHR-custom.

---

## Struct Khusus Mikhmon (`fields/mikhmon.go`)

Paket `fields` juga menyediakan struct typed khusus untuk alur dan resource Mikhmon v4:

- `MikhmonSalesRow`, `MikhmonSalesAdd`, `MikhmonSalesSet` (alias `MikhmonSaleRow`, etc.) — laporan penjualan (`mikhmon-sales`, path `/system/script`).
- `MikhmonQuickPrintRow`, `MikhmonQuickPrintAdd`, `MikhmonQuickPrintSet` — preset cetak cepat voucher (`comment=QuickPrintMikhmon`).
- `MikhmonVoucherRow` — baris voucher hotspot lengkap dengan password (`/ip/hotspot/user` view `voucher`).
- `MikhmonUserAdd`, `MikhmonUserSet` — parameter tambah/ubah voucher hotspot standar Mikhmon v4.
- `MikhmonProfileRow`, `MikhmonProfileAdd`, `MikhmonProfileSet` — parameter dan hasil detail user profile paket hotspot Mikhmon.
- `MikhmonSchedulerRow`, `MikhmonSchedulerAdd`, `MikhmonSchedulerSet` — entri scheduler untuk `Monitor Profile <profile>` dan `Mikhmon-Expire-Monitor`.
- `MikhmonActiveRow` — sesi hotspot aktif untuk dashboard Mikhmon (`/ip/hotspot/active`).

```go
// Decode laporan penjualan Mikhmon
salesRows, err := ros.DecodeAll[fields.MikhmonSalesRow](rows)
for _, s := range salesRows {
    fmt.Println(s.Name, s.Owner, s.Source)
}

// Decode voucher hotspot termasuk password untuk dicetak
vouchers, err := ros.DecodeAll[fields.MikhmonVoucherRow](voucherRows)
for _, v := range vouchers {
    fmt.Println(v.Name, v.Password, v.Profile, v.LimitUptime)
}
```

Pergi ke: [Errors](09-errors.md) | [Catalog](02-catalog.md) | [Views](04-views.md) | [Mikhmon](12-mikhmon.md).
