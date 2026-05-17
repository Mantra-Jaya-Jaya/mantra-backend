# Mantra Backend Service

Backend service untuk aplikasi MANTRA. Dibangun dengan **Golang + Gin**, database **PostgreSQL**, dan autentikasi **JWT**.

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
- [6. Referensi Dokumentasi](#6-referensi-dokumentasi)

---

## 1. Overview & Arsitektur

- **Single-tenant:** Satu instalasi = satu toko. Tidak menggunakan multi-tenant.
- **Stateless JWT:** Access token 15 menit (Bearer/Cookie), refresh token 7 hari.
- **Multi-Client Support:** Flutter menggunakan Bearer Token, Next.js menggunakan httpOnly Cookie.
- **Currency Handling:** Semua nilai uang disimpan sebagai `int64` (Rupiah), formatting di sisi client.

---

## 2. Struktur Folder & Domain

### 2.1 Penjelasan Direktori

```text
backend/
├── config/             # Konfigurasi koneksi database
├── controllers/        # Handler HTTP, diorganisir per domain (bukan per role)
├── db/                 # Koneksi dan inisialisasi database
├── docs/               # Schema database ERD lengkap (.dbml)
├── middleware/         # Security checks (Auth, Role, Ownership, RateLimit, CORS)
├── models/             # Struct database (Representasi tabel GORM)
├── routes/             # Definisi route + pemasangan middleware per role
├── seeders/            # Data awal (seeder) untuk database
├── atlas.hcl           # Konfigurasi Atlas CLI
├── Makefile            # Daftar perintah Atlas CLI
└── main.go             # Entry point server Golang
```

> **Catatan Penting Tim:** Buat file struct (tabel) baru HANYA di dalam folder `models/`. Format nama struct wajib Kapital (contoh: `Keranjang`), dan nama field database wajib pakai tag column snake_case (contoh: `gorm:"column:id_kategori"`).

### 2.2 Domain Controller

Controller diorganisir per **domain bisnis**, bukan per role. Satu function controller bisa dipakai oleh beberapa role sekaligus. Yang mengontrol akses adalah middleware di layer `routes/`.

| Domain         | Folder              | Diakses oleh                  |
|----------------|---------------------|-------------------------------|
| Auth           | `auth/`             | Semua role (public sebagian)  |
| Katalog        | `katalog/`          | Customer, Kasir, Admin        |
| Transaksi      | `transaksi/`        | Customer, Kasir, Admin        |
| Pengantaran    | `pengantaran/`      | Customer (lacak), Admin       |
| Stok           | `stok/`             | Admin                         |
| User           | `user/`             | Masing-masing role (profil)   |
| Notifikasi     | `notifikasi/`       | Customer, Kasir, Admin        |
| Keranjang      | `keranjang/`        | Customer                      |

---

## 3. Panduan Setup Lokal

Ikuti langkah-langkah di bawah ini secara berurutan agar server berjalan lancar di komputer lokal Anda.

### 3.1 Prasyarat & Buat Database Kosong
1. Pastikan **PostgreSQL** sudah terinstall dan berjalan.
2. Buat database baru bernama `mantra_db`.
   ```sql
   CREATE DATABASE mantra_db;
   ```
   *(Cukup bikin databasenya aja, GORM akan otomatis membuatkan tabelnya).*

### 3.2 Konfigurasi Environment (.env)
1. Masuk ke folder backend: `cd backend`
2. Download semua library: `go mod tidy`
3. Sesuaikan konfigurasi koneksi database. Buka file `config/database.go`, cari baris DSN:
   ```go
   dsn := "host=localhost user=postgres password=123456 dbname=mantra_db port=5432 sslmode=disable"
   ```
   **WAJIB DIGANTI:** Ubah bagian `password=123456` menjadi password akun PostgreSQL di laptop Anda.
   *(Nantinya, kita akan migrasi full ke `.env` file untuk hal ini).*

Variabel environment standar yang digunakan:
```env
DB_HOST=
DB_PORT=
DB_USER=
DB_PASSWORD=
DB_NAME=
JWT_SECRET=
JWT_REFRESH_SECRET=
MIDTRANS_SERVER_KEY=
MIDTRANS_CLIENT_KEY=
ALLOWED_ORIGIN=
```

### 3.3 Menjalankan Server
Jalankan perintah ini di terminal:
```bash
go run main.go
```
Jika terminal menampilkan `Database Connected & Migrated Successfully!`, berarti API siap diakses di `http://localhost:8080`.

---

## 4. Migrasi Database (Atlas CLI)

Atlas CLI adalah engine untuk menyinkronkan struktur tabel di PostgreSQL secara otomatis berdasarkan GORM Structs (Single Source of Truth).

### 4.1 Instalasi Atlas

1. **Install Atlas CLI**:
   *   **Linux / macOS (atau Windows dengan Git Bash):**
       ```bash
       curl -sSf https://atlasgo.sh | sh
       ```
   *   **Windows (PowerShell):**
       ```powershell
       Invoke-WebRequest https://release.ariga.io/atlas/atlas-windows-amd64-latest.exe -OutFile atlas.exe
       ```
       *Catatan:* Pindahkan file `atlas.exe` ke sebuah folder khusus (misalnya `C:\atlas`), lalu tambahkan path folder tersebut ke dalam *Environment Variable* `Path` sistem Windows Anda.
2. **Install Atlas Provider GORM (Semua OS):**
   ```bash
   go install ariga.io/atlas-provider-gorm@latest
   ```
3. **Buat database "sandbox" bernama `mantra_dev` di PostgreSQL (Semua OS):**
   Database ini digunakan oleh Atlas untuk melakukan komparasi skema model GORM secara dinamis.
   ```sql
   CREATE DATABASE mantra_dev;
   ```

### 4.2 Workflow Migrasi (Makefile vs PowerShell)

> [!NOTE]
> Perintah `make` secara *default* langsung tersedia di Linux/macOS. Bagi pengguna **Windows**, disarankan menggunakan terminal **Git Bash** agar dapat langsung mengeksekusi `make` dengan mudah.

#### Opsi 1: Menggunakan `make` (Linux / macOS / Git Bash)
Jalankan perintah ini di dalam direktori `backend/`:
- `make db-diff` : Mendeteksi dan menampilkan perubahan skema dari model ke database sandbox.
- `make db-plan` : Simulasi *raw SQL* yang akan dieksekusi (*Dry Run*).
- `make db-apply`: Menerapkan perubahan struktur tabel secara permanen ke database utama (`mantra_db`).
- `make db-inspect`: Melihat representasi HCL/SQL dari tabel yang ada di DB utama.
- `make db-ui` : Membuka visualisasi relasi tabel secara interaktif di browser lokal.

#### Opsi 2: Menggunakan PowerShell Manual (Windows tanpa Make)
Jika menggunakan PowerShell biasa, lakukan ekspor variabel `.env` lalu panggil *binary* `atlas` secara langsung:
```powershell
# 1. Set environment variables (sesuaikan nilai password/port PostgreSQL Anda)
$env:DB_HOST="localhost"; $env:DB_USER="postgres"; $env:DB_PASSWORD="yourpassword"; $env:DB_NAME="mantra_db"; $env:DB_PORT="5432"

# 2. Jalankan perintah atlas langsung
# Setara dengan 'make db-diff'
atlas schema diff --env local --from "env://from" --to "env://to"

# Setara dengan 'make db-apply'
atlas schema apply --env local --to "env://to"
```

---

## 5. Keamanan & Middleware

Urutan eksekusi: `Request → RateLimit → CORS → AuthMiddleware → RoleMiddleware → OwnershipMiddleware → Controller`

- **AuthMiddleware:** Verifikasi JWT dari header/cookie.
- **RoleMiddleware:** Cek role (customer/kasir/admin/kurir).
- **OwnershipMiddleware:** Cek kepemilikan resource (contoh: pesanan hanya bisa diakses pembeli terkait).
- **ID Obfuscation:** Menggunakan UUID (`public_id`) untuk endpoint eksternal agar ID tidak bisa ditebak.

---

## 6. Referensi Dokumentasi

| Dokumen                             | Kegunaan                                         |
|-------------------------------------|--------------------------------------------------|
| `../docs/api-contract.md`           | Semua endpoint API, request/response, dan error  |
| `../docs/security-rbac.md`          | Kebijakan JWT, Rate limit, Middleware, RBAC      |
| `docs/mantra.dbml`                  | Skema ERD Database                               |