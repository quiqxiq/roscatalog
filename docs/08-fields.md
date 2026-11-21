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

Pergi ke: [Errors](09-errors.md) | [Catalog](02-catalog.md) | [Views](04-views.md).
