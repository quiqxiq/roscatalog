# 04 — Views: mode, proplist, print presets

Setiap view adalah satu cara membaca dari resource. Mode + Command:

| Mode | Arti | Command tipikal |
|---|---|---|
| `list` | oneshot dump (seluruh baris) | `/menu/print` |
| `live` | terus-menerus per-interval, send-row on change | `/menu/print` interval |
| `events` | `/path/listen` — perubahan item saja (no seed) | `/menu/listen` |
| `tail` / `tail-only` | follow ke belakang (tidak punya .id baru) | `/menu/print` |
| `monitor` | native (monitor paket/koneksi) | `/interface/monitor-traffic`, `/tool/traceroute` |
| `run` | native (one-shot/self-terminating) | `/ping`, `/tool/fetch` |

## 4.1 Proplist

```go
res, _ := ros.Lookup("interface")
req, _ := res.Build("list", ros.Once(), ros.Params{})                    // Proplist dari view
req, _ := res.Build("list", ros.Once(), ros.Params{Only: []string{"name", "running"}}) // subset (validasi)
req, _ := res.Build("list", ros.Once(), ros.Params{Proplist: []string{"rx-byte"}})     // kolom luar katalog (tidak divalidasi)
req, _ := res.Build("list", ros.Once(), ros.Params{Raw: true})                         // semua kolom, tanpa .proplist
```

Prioritas lihat di `Params`: `Raw > Proplist > Only`. Bila `Only` isi field yang
tidak ada di view → `ErrUnknownColumn`.

## 4.2 Print presets (lihat `PrintFlags`)

MikroDash memakai preset: `ros.Print(name, props, flags...)` — bisa daftar
`"brief"`, `"detail"`, `"stats"`, `"stats-detail"` pada view print. Builder
katalog:

```go
res.Print("print-comment", P("name", "comment"), "detail")
res.Stats("name", "rx-byte")          // view = stat view, proplist interaktif
res.Brief("name", "type")
```

## 4.3 Listen — dilarang proplist

Di RouterOS 7.24.6, `=.proplist=` pada `/path/listen` menekan event `.dead=true`
(penghapusan item tidak terdeteksi). Oleh sebab itu **`ModeListen` tidak mengirim
`=.proplist=` ke router** (`NoServerProplist=true`); kolom tetap disaring di sisi klien.

Praktiknya: listen pakai `Params.Raw` atau plain Params.

## 4.4 Field ke user (lihat [Fields](08-fields.md))

```go
rows, err := sess.Once(ctx, req)
if err == nil {
	ifs, _ := ros.DecodeAll[fields.InterfaceRow](rows)
	_ = ifs
}
```

Pergi ke: [Fields](08-fields.md) | [Realtime](05-realtime.md).
