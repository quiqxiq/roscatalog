# 03 — Custom Resources

Membangun menu kustom pemakai (tidak perlu diedit di repo roscatalog).

## 3.1 Membangun resource

```go
res := ros.Table("queue-custom", ros.GroupQueue, "Queue custom", "/queue/custom").
	Describe("Menu pemakai; lihat [Docs](04-views.md) saat pakai print").
	Filter("name", "target").                     // filter diterima sebagai ?prop=value
	MarkHeavy().                                  // saran: kalau data banyak pakai count/trim
	KeyBy("name").                                // field utama (list), pakai list watch.
	UseServerProplist()

	// view read — lihat [Views](04-views.md)
	.List("name", "target", "rate", "comment").
	Count().
	Live(time.Second, true, "name", "rate").
	Tail("time", "topics", "message").            // tail (follow ke back)
	Events().                                     // listen — punya beberapa router
	Print("print-comment", P("name", "comment"), "detail").

	// aksi write
	.Add(P("name"), P("target", "rate", "comment")).
	Set(P("target", "rate")).
	Remove("hapus item").
	Toggle()
```

## 3.2 Validate

Sebelum menyentuh router, jalankan `res.Validate()` atau `cat.Validate()`;
error yang dijatuhkan satu dari sentinel dari [Errors](09-errors.md)
(`ErrInvalidResource`, `ErrDuplicateResource`, ...). `Catalog.Add` validasi atomik.

## 3.3 Panggil Build

```go
req, err := res.Build("list", ros.Params{
	Only:     []string{"name", "target"},  // kolom, harus ada di daftar view
	Proplist: []string{"name"},            // kolom luar katalog — pakai Raw/Proplist
	Raw:      false,
})
```

## 3.4 Minta di catalog: 


```go
cat := ros.NewCatalog(ros.Builtin().All()...)
_ = cat.Add(res)                       // Assumption: nil error
```

Pergi ke: [Views](04-views.md) | [Realtime](05-realtime.md) | [Catalog](02-catalog.md).
