# MANTRA — Backend

Backend service untuk aplikasi MANTRA. Dibangun dengan **Golang + Gin**, database **PostgreSQL**, autentikasi **JWT**.

---

## Stack

| Komponen       | Teknologi                        |
|----------------|----------------------------------|
| Language       | Go                               |
| Framework      | Gin                              |
| Database       | PostgreSQL                       |
| Auth           | JWT (golang-jwt/jwt)             |
| Payment        | Midtrans Snap                    |
| Client Mobile  | Flutter (Customer & Kasir)       |
| Client Web     | Next.js (Admin)                  |

---

## Arsitektur Singkat

- **Single-tenant** — satu instalasi = satu toko
- **Stateless JWT** — access token 15 menit, refresh token 7 hari
- **Flutter** pakai Bearer Token, **Next.js** pakai httpOnly Cookie
- Semua nilai uang disimpan sebagai `int64` (Rupiah), formatting di sisi client

---

## Struktur Folder

```
backend/
├── controllers/        # Handler HTTP, diorganisir per domain (bukan per role)
│   ├── auth/           # Login, logout, register, refresh token, change password
│   ├── katalog/        # Barang, kategori, diskon, satuan, spesifikasi, barcode
│   ├── transaksi/      # Pesanan, pembayaran, checkout (online & POS)
│   ├── pengantaran/    # Pengantaran, ekspedisi, status pengantaran
│   ├── stok/           # Stok opname, riwayat stok
│   ├── user/           # Profil customer, kasir, kurir, alamat
│   ├── notifikasi/     # Notifikasi semua role
│   └── keranjang/      # Keranjang belanja customer
├── middleware/         # Auth, Role, Ownership, RateLimit, CORS
├── models/             # Struct database (sesuai ERD)
├── routes/             # Definisi route + pemasangan middleware per role
├── db/                 # Koneksi dan inisialisasi database
└── main.go
```

> **Penting:** Controller diorganisir per **domain bisnis**, bukan per role.
> Satu function controller bisa dipakai oleh beberapa role sekaligus.
> Yang mengontrol akses adalah middleware di layer `routes/`, bukan controller.

---

## Domain Controller & Aksesnya

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

## Middleware Stack

Urutan middleware untuk setiap request:

```
Request → RateLimit → CORS → AuthMiddleware → RoleMiddleware → OwnershipMiddleware → Controller
```

| Middleware           | Fungsi                                                    |
|----------------------|-----------------------------------------------------------|
| `RateLimit`          | Batasi request per IP (login 5/menit, global 60/menit)    |
| `AuthMiddleware`     | Verifikasi JWT dari Bearer header atau httpOnly Cookie     |
| `RoleMiddleware`     | Cek role dari JWT claims, tolak jika tidak sesuai         |
| `OwnershipMiddleware`| Cek kepemilikan resource personal (pesanan, alamat, dll)  |

---

## Autentikasi

| Client      | Mekanisme                        | Token disimpan di             |
|-------------|----------------------------------|-------------------------------|
| Flutter     | `Authorization: Bearer <token>`  | `flutter_secure_storage`      |
| Next.js     | `httpOnly Cookie`                | Cookie (tidak bisa diakses JS)|

JWT Claims: `user_id`, `role`, `iat`, `exp`

---

## Environment Variables

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
ALLOWED_ORIGIN=        # domain Next.js untuk CORS
```

---

## Dokumentasi Lengkap

| Dokumen             | Isi                                              |
|---------------------|--------------------------------------------------|
| `api-contract.md`   | Semua endpoint, request/response, error codes    |
| `security-rbac.md`  | JWT, middleware, RBAC, ownership rules           |
| `mantra.dbml`       | Schema database lengkap (buka di dbdiagram.io)   |
| `AGENTS.md`         | Panduan untuk AI agent mengerjakan codebase ini  |
