# 03 — Custom Resources

Membangun menu kustom pemakai (tidak perlu diedit di repo roscatalog).

## 3.1 Membangun resource

```go
res := ros.Table("queue-custom", ros.GroupQueue, "Queue custom", "/queue/custom").
	Describe("Menu pemakai; lihat [Views](04-views.md) saat pakai print").
	Filter("name", "target").
	List("name", "target", "rate", "comment").
	MarkHeavy().
	KeyBy("name").
	UseServerProplist().
	Count().
	Live(time.Second, true, "name", "rate").
	Tail("time", "topics", "message").
	Events().
	Print("print-comment", ros.P("name", "comment"), "detail").
	Add(ros.P("name"), ros.P("target", "rate", "comment")).
	Set("target", "rate").
	Remove("hapus item").
	Toggle()
```

## 3.2 Validate

Sebelum menyentuh router, jalankan `res.Validate()` atau `cat.Validate()`;
error yang dijatuhkan satu dari sentinel dari [Errors](09-errors.md)
(`ErrInvalidResource`, `ErrDuplicateResource`, ...). `Catalog.Add` validasi atomik.

## 3.3 Panggil Build

```go
req, err := res.Build("list", ros.Once(), ros.Params{
	Only: []string{"name", "target"}, // kolom, harus ada di daftar view
})
```

## 3.4 Minta di catalog: 


```go
cat := ros.NewCatalog(ros.Builtin().All()...)
_ = cat.Add(res)                       // Assumption: nil error
```

Pergi ke: [Views](04-views.md) | [Realtime](05-realtime.md) | [Catalog](02-catalog.md).
