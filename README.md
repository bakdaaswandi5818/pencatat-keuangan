# pencatat-keuangan

REST API pencatat keuangan menggunakan Go (Echo v4), GORM, dan SQLite.

## Authorization

Semua endpoint API (kecuali `/health`) membutuhkan header:

```
Authorization: Bearer <API_KEY>
```

Set `API_KEY` lewat environment variable saat menjalankan aplikasi.

Contoh:

```bash
API_KEY=my-secret-key CGO_ENABLED=1 go run cmd/main.go
```

`API_KEY` wajib diisi. Aplikasi akan gagal start jika `API_KEY` kosong.

## Build Linux (SQLite butuh CGO)

```bash
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/finance-api cmd/main.go
```
