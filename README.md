# Mantra Backend Service

Backend service untuk aplikasi **MANTRA**. Dibangun dengan **Golang + Gin**, database **PostgreSQL**, ORM **GORM**, autentikasi **JWT**, dan penyimpanan file **MinIO** (S3-compatible).

---

## Daftar Isi
- [1. Overview & Arsitektur](#1-overview--arsitektur)
- [2. Struktur Folder & Domain](#2-struktur-folder--domain)
  - [2.1 Penjelasan Direktori](#21-penjelasan-direktori)
  - [2.2 Domain Controller](#22-domain-controller)
- [3. Panduan Setup Lokal](#3-panduan-setup-lokal)
  - [3.1 Prasyarat & Buat Database Kosong](#31-prasyarat--buat-database-kosong)
  - [3.2 Konfigurasi Environment (.env)](#32-konfigurasi-environment-env)
  - [3.3 Menjalankan Server](#33-menjalankan-server)
- [4. Migrasi Database (Atlas CLI)](#4-migrasi-database-atlas-cli)
  - [4.1 Instalasi Atlas](#41-instalasi-atlas)
  - [4.2 Workflow Makefile](#42-workflow-makefile)
- [5. Keamanan & Middleware](#5-keamanan--middleware)
- [6. Storage (MinIO)](#6-storage-minio)
- [7. Referensi Dokumentasi](#7-referensi-dokumentasi)

---

## 1. Overview & Arsitektur

- **Single-tenant:** Satu instalasi = satu toko. Tidak menggunakan multi-tenant.
- **Stateless JWT:** Access token 30 menit (Bearer / httpOnly Cookie), refresh token 7 hari.
- **Multi-Client Support:**
  - **Flutter** (Customer & Kasir apps) menggunakan `Bearer Token` via `Authorization` header.
  - **Next.js** (Admin Web) menggunakan `httpOnly Cookie` (`access_token`, `refresh_token`).
  - Deteksi client via header `X-Client-Type: flutter | nextjs`.
- **Sliding Expiration:** Token diperbarui otomatis jika sisa masa berlaku < 15 menit.
- **Password:** bcrypt dengan cost factor 12.
- **ID Obfuscation:** UUID (`public_id`) digunakan di endpoint eksternal agar ID numerik tidak bisa ditebak.
- **Currency Handling:** Semua nilai uang disimpan sebagai `int` (Rupiah), formatting di sisi client.

---

## 2. Struktur Folder & Domain

### 2.1 Penjelasan Direktori

```text
mantra-backend/
├── .air.toml            # Konfigurasi Air (hot reload)
├── atlas.hcl            # Konfigurasi Atlas CLI (migrasi)
├── Makefile             # Perintah Atlas, Go, dan server
├── main.go              # Entry point server Golang
├── config/
│   ├── database.go      # Koneksi PostgreSQL + AutoMigrate (25 tabel)
│   └── minio.go         # Inisialisasi MinIO client
├── controllers/         # Handler HTTP, diorganisir per domain bisnis
│   ├── auth/            # Login, Register, Refresh Token, Logout, Change Password
│   ├── katalog/         # Barang, Kategori, Diskon, Promo, Scan Barcode
│   ├── keranjang/       # CRUD keranjang belanja (Customer)
│   ├── notifikasi/      # Notifikasi per user
│   ├── pengantaran/     # Lacak pengiriman
│   ├── stok/            # Stok opname (Admin)
│   ├── transaksi/       # Pesanan, Pembayaran, Dashboard Kasir, Laporan
│   └── user/            # Profil Customer, Kasir, Admin, Karyawan, Alamat
├── db/                  # Koneksi dan inisialisasi database
├── docs/
│   ├── mantra.dbml       # Schema database ERD (DBML)
│   └── mantra.dbdiagram  # Layout diagram dbdiagram.io
├── middleware/
│   ├── auth_middleware.go      # Validasi JWT + sliding expiration
│   ├── role_middleware.go      # Role-based access (Customer/Kasir/Admin/Kurir)
│   └── ownership_middleware.go # Cek kepemilikan resource (bypass untuk Admin)
├── models/              # 25 struct GORM (representasi tabel database)
├── routes/
│   └── routes.go        # Definisi route + pemasangan middleware per grup
├── seeders/             # 22 seeder + orchestrator untuk data awal
└── utils/
    └── minio_helper.go  # Helper upload file ke MinIO
```

> **Catatan Penting Tim:** Buat file struct (tabel) baru HANYA di dalam folder `models/`. Format nama struct wajib Kapital (contoh: `Keranjang`), dan nama field database wajib pakai tag column snake_case (contoh: `gorm:"column:id_kategori"`).

### 2.2 Domain Controller

Controller diorganisir per **domain bisnis**, bukan per role. Satu function controller bisa dipakai oleh beberapa role sekaligus. Yang mengontrol akses adalah middleware di layer `routes/`.

| Domain         | Package              | Diakses oleh                  |
|----------------|----------------------|-------------------------------|
| Auth           | `auth/`              | Semua role (public sebagian)  |
| Katalog        | `katalog/`           | Customer, Kasir, Admin        |
| Keranjang      | `keranjang/`         | Customer                      |
| Transaksi      | `transaksi/`         | Customer, Kasir, Admin        |
| Pengantaran    | `pengantaran/`       | Customer (lacak), Admin       |
| Stok           | `stok/`              | Admin                         |
| User & Profil  | `user/`              | Customer, Kasir, Admin        |
| Notifikasi     | `notifikasi/`        | Customer, Kasir, Admin        |

---

## 3. Panduan Setup Lokal

Ikuti langkah-langkah di bawah ini secara berurutan agar server berjalan lancar di komputer lokal Anda.

### 3.1 Prasyarat & Buat Database Kosong
1. Pastikan **PostgreSQL** sudah terinstall dan berjalan.
2. Buat database baru bernama `mantra_db`.
   ```sql
   CREATE DATABASE mantra_db;
   ```
   *(Cukup bikin databasenya aja, GORM akan otomatis membuatkan tabelnya via AutoMigrate).*

### 3.2 Konfigurasi Environment (.env)
1. Buka folder `mantra-backend/`.
2. Copy file `.env.example` menjadi `.env`:
   ```bash
   cp .env.example .env
   ```
3. Download semua library: `go mod tidy`
4. Sesuaikan isi `.env` dengan konfigurasi lokal Anda.

### 3.3 Menjalankan Server

**Opsi A — Langsung:**
```bash
go run main.go
```

**Opsi B — Hot reload (Air):**
```bash
# Install Air (satu kali)
go install github.com/air-verse/air@latest

# Jalankan dengan Air (auto-reload saat file berubah)
air
```

Jika terminal menampilkan `Database Connected & Migrated Successfully!`, berarti API siap diakses di `http://localhost:8080`.

---

## 4. Migrasi Database (Atlas CLI)

Atlas CLI adalah engine untuk menyinkronkan struktur tabel di PostgreSQL secara otomatis berdasarkan GORM Structs (Single Source of Truth). File konfigurasi: `atlas.hcl`.

### 4.1 Instalasi Atlas

1. **Install Atlas CLI**:
   - **Linux / macOS (atau Windows dengan Git Bash):**
     ```bash
     curl -sSf https://atlasgo.sh | sh
     ```
   - **Windows (PowerShell):**
     ```powershell
     Invoke-WebRequest https://release.ariga.io/atlas/atlas-windows-amd64-latest.exe -OutFile atlas.exe
     ```

2. **Install Atlas Provider GORM (Semua OS):**
   ```bash
   go install ariga.io/atlas-provider-gorm@latest
   ```

3. **Buat database sandbox `mantra_dev`** (digunakan Atlas untuk komparasi skema):
   ```sql
   CREATE DATABASE mantra_dev;
   ```

### 4.2 Workflow Makefile

Jalankan perintah di dalam direktori `mantra-backend/`:

| Perintah            | Kegunaan                                            |
|---------------------|-----------------------------------------------------|
| `make db-diff`      | Deteksi perubahan skema model → sandbox             |
| `make db-plan`      | Simulasi raw SQL (Dry Run) sebelum apply            |
| `make db-apply`     | Terapkan perubahan struktur ke database utama       |
| `make db-inspect`   | Lihat representasi HCL/SQL dari DB utama            |
| `make db-ui`        | Buka visualisasi ERD interaktif di browser          |
| `make db-clean`     | Hapus seluruh skema DB utama (⚠️ reset total)      |
| `make db-clean-dev` | Hapus skema sandbox `mantra_dev`                   |
| `make tidy`         | `go mod tidy` + `go fmt` + `go vet`                |
| `make run`          | Jalankan server (`go run main.go`)                  |
| `make build`        | Kompilasi binary ke `bin/mantra-backend`            |
| `make test`         | Jalankan seluruh unit test                          |
| `make clean`        | Hapus folder `bin/` + cache test                   |

---

## 5. Keamanan & Middleware

Urutan eksekusi middleware pada setiap request:

```
Request → RateLimit (future) → CORS → AuthMiddleware → RoleMiddleware → Controller
```

### AuthMiddleware (`middleware/auth_middleware.go`)
- Validasi JWT dari `Authorization: Bearer <token>` (Flutter) atau cookie `access_token` (Next.js).
- **Sliding expiration:** Jika masa token tersisa < 15 menit, token baru digenerate dan dikirim via header `X-New-Access-Token` + cookie.
- Set `user_id`, `public_id`, `role` ke context Gin.

### RoleMiddleware (`middleware/role_middleware.go`)
- Cek apakah role user (dari token) termasuk dalam daftar role yang diizinkan.
- Error code: `AUTH_002` (role tidak ditemukan), `AUTH_003` (tidak punya izin).

### OwnershipMiddleware (`middleware/ownership_middleware.go`)
- Cek kepemilikan resource berdasarkan `public_id` di URL.
- **Admin selalu bypass** pengecekan ini.
- Error code: `AUTH_004`.

### Authentication Flow
| Endpoint                          | Method | Middleware         |
|-----------------------------------|--------|--------------------|
| `POST /api/v1/login`              | Public | -                  |
| `POST /api/v1/register`           | Public | -                  |
| `POST /api/v1/auth/refresh`       | Public | -                  |
| `GET /api/v1/scan/:kode_barcode`  | Public | -                  |
| `POST /api/v1/logout`             | Auth   | AuthMiddleware     |
| `PUT /api/v1/change-password`     | Auth   | AuthMiddleware     |
| `/api/v1/customer/*`              | Auth   | AuthMiddleware     |
| `/api/v1/kasir/*`                 | Auth   | AuthMiddleware     |
| `/api/v1/admin/*`                 | Auth   | AuthMiddleware     |

> RBAC penuh via role checking di masing-masing controller. Untuk endpoint yang perlu dicek role-nya, gunakan `RoleMiddleware("admin", "kasir")`.

---

## 6. Storage (MinIO)

MANTRA menggunakan **MinIO** (self-hosted, S3-compatible) untuk menyimpan file:

- **Gambar barang** → folder `produk/`
- **Ikon kategori** → folder `kategori/`
- **Foto profil karyawan** → folder `profile/`
- **Banner diskon** → folder `diskon/`

Konfigurasi di `.env`:
```env
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=your_access_key
MINIO_SECRET_KEY=your_secret_key
MINIO_BUCKET=mantra-storage
```

URL publik file: `https://storage.mantra.web.id/mantra-storage/{folder}/{tahun}/{bulan}/{uuid}.{ext}`

---

## 7. Referensi Dokumentasi

| Dokumen                             | Kegunaan                                         |
|-------------------------------------|--------------------------------------------------|
| `docs/README.md`                    | Index seluruh dokumentasi                        |
| `docs/api/overview.md`              | Format API, error codes, autentikasi             |
| `docs/api/auth.md`                  | Endpoint login, register, refresh, logout        |
| `docs/api/katalog/`                 | Barang, kategori, diskon, promo, scan barcode    |
| `docs/api/keranjang.md`             | CRUD keranjang belanja customer                  |
| `docs/api/transaksi/`               | Pesanan, pembayaran, POS kasir                   |
| `docs/api/user/`                    | Profil customer, kasir, admin, karyawan, alamat  |
| `docs/api/notifikasi.md`            | Notifikasi per role                              |
| `docs/api/pengantaran.md`           | Lacak pengiriman pesanan                         |
| `docs/architecture.md`              | Component interaction, decision records          |
| `docs/security/`                    | JWT, middleware chain, RBAC matrix, error codes  |
| `docs/database/erd.md`              | Panduan baca ERD, konvensi tabel & kolom         |
| `docs/database/migrations.md`       | Workflow Atlas CLI, cara nambah tabel            |
| `docs/deployment.md`                | Build binary, env vars, production checklist     |
