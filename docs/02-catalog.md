# 02 — Catalog

Catalog adalah kumpulan resource yang bisa di-query dan dikembangkan. Dua hal:
mencari di catalog bawaan (`Builtin`) dan membangun catalog pemakai sendiri
(`NewCatalog` + `Add/Replace/Extend`).

## 2.1 API catalog

```go
// Builtin: catalog kurasi bawaan (125 resource saat ini).
cat := ros.Builtin()

res, err := cat.Find("hotspot-active")      // err = ErrResourceNotFound bila hilang
res, ok := cat.Lookup("hotspot-active")     // boolean lookup

all := cat.All()                            // urut kan tanah + ID
hs := cat.ByGroup(ros.GroupHotspot)         // urut ID (GroupHotspot, GroupPPP, GroupIP, GroupFirewall, GroupQueue, GroupInterface, GroupSystem, GroupTools)
n := cat.Len()
```

Error dari `Find` di-wrap `ErrResourceNotFound` (`errors.Is`).

## 2.2 Custom catalog

```go
cat := ros.NewCatalog(ros.Builtin().All()...)  // mulai dari katalog lama

if err := cat.Add(myResource); err != nil { ... }   // ErrDuplicateResource / ErrInvalidResource
if err := cat.Replace(replacement); err != nil { }  // Replace khusus ID yang sudah ada
if err := cat.Extend(other); err != nil { }         // atomik: tidak ada yang diambil sebelum ada error
if err := cat.Validate(); err != nil { }            // memeriksa semua (ID kosong/duplikat, path, selector...)
```

## 2.3 Anatomi Resource (untuk membaca atau menggunakan catalog)

```go
type Resource struct {
	ID, Group, Title, Path                 // identitas
	Singleton bool                         // menu settings (tidak punya .id per baris)
	Filter   []string                      // field aman sebagai filter ?prop=value
	Views    []View                        // data: list, live, events, tail, monitor, run...
	Actions  []Action                      // writes: set, add, remove, enable...
	Notes    string
}
```

Cara membaca: `View` bisa saya eksekusi lewat `Build`; `Action` lewat `Sentence(selector, args, tag)`.

Lihat [Custom Resources](03-custom-resources.md) untuk `Table/Single/Tool` (builder).
Perbedaan prioritas lihat [Views](04-views.md) (`Params.Raw > Proplist > Only`).

Pergi ke: [Custom Resources](03-custom-resources.md) | [Views](04-views.md).
