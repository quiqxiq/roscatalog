# 09 — Error

Semua error dari `Build/Find/View/Action/Exec` di-wrap dengan **error sentinel**
dari package root (`errors.Is`). Library internal (routeros) punya
`DeviceError` (trap dari router, `!empty`, etc).

## Sentinel (package roscatalog)

```go
var (
	ErrInvalidResource      = errors.New("roscatalog: definisi resource tidak valid")
	ErrDuplicateResource    = errors.New("roscatalog: resource duplikat")
	ErrResourceNotFound     = errors.New("roscatalog: resource tidak ditemukan")
	ErrViewNotFound         = errors.New("roscatalog: view tidak ditemukan")
	ErrDeliveryUnsupported  = errors.New("roscatalog: delivery tidak didukung")
	ErrUnknownFlag          = errors.New("roscatalog: flag tidak didukung")
	ErrUnknownColumn        = errors.New("roscatalog: kolom tidak ada")
	ErrMissingArgument      = errors.New("roscatalog: argumen wajib belum diisi")
	ErrSelectorRequired     = errors.New("roscatalog: selector wajib diisi")
	ErrConnectionDown       = errors.New("rosrun: connection down")
	ErrStatusTimeout        = errors.New("rosrun: status ping timeout")
	ErrClosed               = errors.New("rosrun: closed")
)
```

(Catatan: `ErrActionNotFound` di-remove di v0.2.0 — `Resource.Action()` kembali `(Action, bool)`.)

## Penggunaan

```go
req, err := res.Build("list", ros.Params{Only: []string{"name", "bogus"}})
switch {
case errors.Is(err, ros.ErrViewNotFound):     // view tidak ada (buat dari catalog custom)
case errors.Is(err, ros.ErrUnknownColumn):    // Only ke luar field view
case errors.Is(err, ros.ErrMissingArgument):  // arg wajib terisi (kirim =name=, =value=)
case errors.Is(err, ros.ErrSelectorRequired): // set/remove butuh selector (.id/.id/.id)
case errors.Is(err, ros.ErrUnknownFlag):      // flag print tidak dideklarasikan
default:
	// error engine (router down, auth, trap, syntax error, timeout)
	// rosrun kirim DeviceError; lihat internal/routeros
}
```

## Penanganan umum

1. **Local validation** — Build/Find/Action/Sentence/Validate.
2. **Router rejection** — `rosrun` kirim request, coba `Raw` lalu cek proplist.
3. **`Fallback` di `Request`** — bila proplist ditolak router, ulang tanpa proplist.
4. **Cek drift** — `rosverify -addr CHR -resource X` (lihat [10-testing](10-testing.md)).

Pergi ke: [Fields](08-fields.md) | [Connection](06-connection-layer.md).
