# 09 — Error

Semua error dari `Build/Find/View/Action/Exec` di-wrap dengan **error sentinel** (`errors.Is`). Library transport (`goros/v4`) memiliki `*routeros.DeviceError` (trap dari router, `!empty`, dll).

## Sentinel

Package root (`roscatalog`):

```go
var (
	ErrInvalidResource     = errors.New("roscatalog: invalid resource definition")
	ErrDuplicateResource   = errors.New("roscatalog: duplicate resource")
	ErrResourceNotFound    = errors.New("roscatalog: resource not found")
	ErrViewNotFound        = errors.New("roscatalog: view not found")
	ErrDeliveryUnsupported = errors.New("roscatalog: delivery unsupported")
	ErrUnknownFlag         = errors.New("roscatalog: unknown flag")
	ErrInvalidFlag         = errors.New("roscatalog: invalid flag")
	ErrUnknownColumn       = errors.New("roscatalog: unknown column")
	ErrMissingArgument     = errors.New("roscatalog: missing required argument")
	ErrSelectorRequired    = errors.New("roscatalog: selector required")
	ErrInvalidInterval     = errors.New("roscatalog: invalid interval")
	ErrInvalidQuery        = errors.New("roscatalog: invalid query")
)
```

Package `rosrun`:

```go
var (
	ErrClosed         = errors.New("rosrun: session closed")
	ErrReconnecting   = errors.New("rosrun: reconnecting")
	ErrUnknownRouter  = errors.New("rosrun: unknown router")
	ErrStatusTimeout  = errors.New("rosrun: status ping timeout")
	ErrConnectionDown = errors.New("rosrun: connection down")
	ErrQueueOverflow  = errors.New("rosrun: stream queue overflow")
	ErrCacheFull      = errors.New("rosrun: cache full")
)
```

## Penggunaan

```go
req, err := res.Build("list", ros.Once(), ros.Params{Only: []string{"name", "bogus"}})
switch {
case errors.Is(err, ros.ErrViewNotFound):     // view tidak ada
case errors.Is(err, ros.ErrUnknownColumn):    // Only di luar field view
case errors.Is(err, ros.ErrMissingArgument):  // argumen wajib kosong
case errors.Is(err, ros.ErrSelectorRequired): // set/remove butuh selector
case errors.Is(err, ros.ErrUnknownFlag):      // flag print tidak dideklarasikan
case errors.Is(err, ros.ErrInvalidFlag):      // kombinasi flag saling eksklusif
case errors.Is(err, ros.ErrInvalidInterval):  // interval di luar batas
default:
	// error transport / DeviceError dari router
}
```

## Penanganan umum

1. **Local validation** — `Build`/`Find`/`Action.Sentence`/`Validate`.
2. **Router rejection** — `rosrun` otomatis mencoba `Fallback` tanpa `.proplist` bila router menolak.
3. **Cek drift** — `rosverify -addr <host:port> -resource <id>` (lihat [10-testing](10-testing.md)).

Pergi ke: [Fields](08-fields.md) | [Connection](06-connection-layer.md).
