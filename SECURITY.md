# Security Policy

## Reporting a Vulnerability

Jika Anda menemukan celah keamanan (misalnya kebocoran kredensial, XSS, atau issue bypass TLS) pada `roscatalog`, **jangan** membuka issue publik terlebih dahulu.

Silakan laporkan secara privat melalui GitHub Security Advisories pada repositori `github.com/quiqxiq/roscatalog` dengan menyertakan:

- Deskripsi masalah dan komponen terdampak (`rosrun`, `cmd/webtest`, dll.).
- Langkah reproduksi atau PoC minimal.
- Dampak potensial dan saran mitigasi bila ada.

## Security Defaults

- `cmd/roscheck` dan `cmd/multicheck` menggunakan `-insecure=false` secara default (verifikasi sertifikat TLS aktif).
- `cmd/webtest` adalah demo dashboard lokal yang secara default bind ke `127.0.0.1:8080`, meredaksi `password` pada DTO JSON, dan mewajibkan token bearer saat bind ke alamat non-loopback.
- `rosrun.StatusConfig` tidak melakukan ping eksternal ke `8.8.8.8` kecuali `Target` dikonfigurasi secara eksplisit.
