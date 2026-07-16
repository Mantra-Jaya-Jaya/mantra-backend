# Include file .env dan export semua variabel di dalamnya agar terbaca oleh Atlas
include .env
export

# Cross-platform detection (Linux / Git Bash / cmd / PowerShell)
BINARY := bin/mantra-backend
RM     := rm -rf bin
PSQL   := PGPASSWORD="$(DB_PASSWORD)" psql -h "$(DB_HOST)" -p "$(DB_PORT)" -U "$(DB_USER)" -d "$(DB_NAME)"

ifeq ($(OS),Windows_NT)
    BINARY := bin\mantra-backend.exe
    ifeq ($(findstring bash,$(SHELL)),bash)
        # Git Bash — syntax sama kayak Linux
    else ifeq ($(findstring powershell,$(SHELL)),powershell)
        RM   := powershell -Command "Remove-Item -Recurse -Force 'bin'"
        PSQL := powershell -Command "$$env:PGPASSWORD='$(DB_PASSWORD)'; psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME)"
    else
        # cmd.exe
        RM   := cmd /c rmdir /s /q bin
        PSQL := cmd /c set PGPASSWORD=$(DB_PASSWORD) && psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME)
    endif
endif

# ==========================================
# DATABASE & MIGRATION (ATLAS)
# ==========================================

# Melihat perbandingan (diff) antara GORM Models dan Database saat ini
db-diff:
	atlas schema diff --env local --from "env://from" --to "env://to"

# Perintah untuk melihat raw SQL (Dry Run) sebelum apply
db-plan:
	atlas schema apply --env local --to "env://to" --dry-run

# Perintah untuk mengeksekusi skema ke database utama
db-apply:
	atlas schema apply --env local --to "env://to" --auto-approve

# Melihat struktur db saat ini via CLI
db-inspect:
	atlas schema inspect --url "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

# Membuka visualisasi skema/ERD secara interaktif di Web UI (Lokal)
db-ui:
	atlas schema inspect --env local --web

# Menghapus seluruh skema (Drop All) di database utama
# PERINGATAN: Hanya gunakan ini di environment lokal saat butuh reset total!
db-clean:
	atlas schema clean --env local --auto-approve
	$(PSQL) -c "CREATE SCHEMA IF NOT EXISTS public; CREATE EXTENSION IF NOT EXISTS \"pgcrypto\";"

db-clean-dev:
	atlas schema clean --url "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" --auto-approve
	$(PSQL) -c "CREATE SCHEMA IF NOT EXISTS public; CREATE EXTENSION IF NOT EXISTS \"pgcrypto\";"
# Membersihkan skema dari database utama dan sandbox sekaligus
db-clean-all: db-clean db-clean-dev

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
	go build -o $(BINARY) main.go

# Menjalankan seluruh unit test di dalam project
test:
	go test ./... -v

# Membersihkan file binary hasil build dan cache testing
clean:
	$(RM)
	go clean -testcache