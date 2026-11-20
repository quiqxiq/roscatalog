package roscatalog

import "errors"

// Error sentinel yang bisa dibandingkan dengan errors.Is. Semua error dari
// Build/Sentence/Find membungkus salah satunya bila penyebabnya jelas.
var (
	// ErrInvalidResource menandai definisi resource yang tidak konsisten.
	ErrInvalidResource = errors.New("roscatalog: definisi resource tidak valid")
	// ErrDuplicateResource menandai ID yang sudah ada di catalog.
	ErrDuplicateResource = errors.New("roscatalog: resource duplikat")
	// ErrResourceNotFound dikembalikan Find/Lookup tingkat error.
	ErrResourceNotFound = errors.New("roscatalog: resource tidak ditemukan")
	// ErrViewNotFound menandai nama view yang tidak ada pada resource.
	ErrViewNotFound = errors.New("roscatalog: view tidak ditemukan")
	// ErrDeliveryUnsupported menandai Plan yang tidak valid untuk view.
	ErrDeliveryUnsupported = errors.New("roscatalog: delivery tidak didukung")
	// ErrUnknownFlag menandai flag print yang tidak dideklarasikan.
	ErrUnknownFlag = errors.New("roscatalog: flag tidak didukung")
	// ErrUnknownColumn menandai kolom di luar daftar field view.
	ErrUnknownColumn = errors.New("roscatalog: kolom tidak ada")
	// ErrMissingArgument menandai argumen wajib yang kosong.
	ErrMissingArgument = errors.New("roscatalog: argumen wajib belum diisi")
	// ErrSelectorRequired menandai aksi yang butuh selector tetapi tidak diisi.
	ErrSelectorRequired = errors.New("roscatalog: selector wajib diisi")
)
