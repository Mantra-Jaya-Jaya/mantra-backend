# Include file .env dan export semua variabel di dalamnya agar terbaca oleh Atlas
include .env
export

# ==========================================
# DATABASE & MIGRATION (ATLAS)
# ==========================================

# Melihat perbandingan (diff) antara GORM Models dan Database saat ini
db-diff:
	atlas schema diff --env local

# Perintah untuk melihat raw SQL (Dry Run) sebelum apply
db-plan:
	atlas schema apply --env local --to "env://to" --dry-run

# Perintah untuk mengeksekusi skema ke database utama
db-apply:
	atlas schema apply --env local --to "env://to"

# Melihat struktur db saat ini via CLI
db-inspect:
	atlas schema inspect --url "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

# Membuka visualisasi skema/ERD secara interaktif di Web UI (Lokal)
db-ui:
	atlas schema inspect --env local --web

# Menghapus seluruh skema (Drop All) di database utama
# PERINGATAN: Hanya gunakan ini di environment lokal saat butuh reset total!
db-clean:
	atlas schema clean --env local

# ==========================================
# GOLANG UTILITIES
# ==========================================

# Merapikan dependencies dan auto-format kode untuk membuang unused imports
tidy:
	go mod tidy
	go fmt ./...
	go vet ./...

# Menjalankan server backend
run:
	go run main.go

# Melakukan kompilasi binary backend (disimpan di folder bin/)
build:
	go build -o bin/mantra-backend main.go

# Menjalankan seluruh unit test di dalam project
test:
	go test ./... -v

# Membersihkan file binary hasil build dan cache testing
clean:
	rm -rf bin/
	go clean -testcache