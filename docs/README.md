# 📚 Dokumentasi MANTRA Backend

Selamat datang di pusat dokumentasi teknis untuk **MANTRA Backend Service**. Dokumentasi ini dirancang untuk memudahkan developer dalam memahami arsitektur, database, alur deployment, serta integrasi pihak ketiga pada sistem MANTRA.

---
### 🧭 Navigasi Cepat
[🏠 Utama](README.md) | [🏛️ Arsitektur](architecture.md) | [🛠️ Deployment](deployment.md) | [💳 Midtrans](pembayaran.md) | [📦 Biteship](biteship.md) | [📡 API Contract](api/overview.md) | [🗄️ Database](database/erd.md) | [🔒 Keamanan](security/README.md)
---

## 🏛️ Fondasi Sistem & Arsitektur

*   [**architecture.md**](architecture.md) - Penjelasan detail mengenai arsitektur sistem, Tech Stack, request flow, hubungan model database, dan keputusan arsitektural utama.
*   [**deployment.md**](deployment.md) - Panduan cara membangun binary, daftar environment variables (`.env`), serta checklist kesiapan produksi.

---

## 💳 & 📦 Integrasi Pihak Ketiga (Gateway)

*   [**pembayaran.md**](pembayaran.md) - Panduan integrasi pembayaran non-tunai (QRIS & Virtual Account) via **Midtrans Core API**, konfigurasi sandbox, setup Ngrok, serta alur testing pembayaran.
*   [**biteship.md**](biteship.md) - Panduan integrasi pengiriman eksternal via **Biteship API** (Rates, Order Shipment, Tracking), setup Ngrok webhook, saldo virtual, dan fallback mock.

---

## 📡 API Contract & Collections

*   [**api/overview.md**](api/overview.md) - Base URL, format komunikasi request/response (sukses & error), HTTP Status Codes, dan format ID (public vs private).
*   **API per Domain:**
    *   [Auth API](api/auth.md) - Autentikasi user (Login, Logout, Refresh Token, Change Password)
    *   [Katalog API](api/katalog/barang.md) - CRUD & Pencarian Barang, Kategori, Diskon, Ekspedisi
    *   [Keranjang API](api/keranjang.md) - Manajemen keranjang belanja customer
    *   [Transaksi API](api/transaksi/pesanan.md) - Checkout, Cek Ongkir, Pembayaran, Notification
    *   [Pengantaran API](api/pengantaran.md) - Manajemen status pengiriman oleh kurir internal
    *   [Notifikasi API](api/notifikasi.md) - Pengiriman notifikasi real-time & history
    *   [User API](api/user/customer.md) - Profil & Manajemen Customer, Karyawan, Kasir, Kurir, Alamat
*   [**api-collections/**](api-collections/) - Folder koleksi API untuk **Bruno** (alternatif Postman/Insomnia) untuk mempermudah testing manual.

---

## 🗄️ Database & Migrasi

*   [**database/erd.md**](database/erd.md) - Panduan membaca skema database, relasi tabel utama, dan konvensi kolom.
*   [**database/migrations.md**](database/migrations.md) - Workflow penggunaan Atlas CLI & GORM AutoMigrate untuk modifikasi skema tabel.
*   [**mantra.dbml**](mantra.dbml) - Source of truth skema database MANTRA (25 tabel).
*   [**mantra.dbdiagram**](mantra.dbdiagram) - Konfigurasi visual layout ERD untuk dbdiagram.io.

---

## 🔒 Keamanan & Penanganan Error

*   [**security/README.md**](security/README.md) - Index dokumentasi sistem keamanan MANTRA.
*   [**security/jwt.md**](security/jwt.md) - Mekanisme JWT, claims structure, dan sliding expiration token.
*   [**security/middleware.md**](security/middleware.md) - Penjelasan Middleware: Auth, Role, dan Ownership.
*   [**security/rbac.md**](security/rbac.md) - Matrix Role-Based Access Control (RBAC) per role & endpoint.
*   [**security/errors.md**](security/errors.md) - Daftar lengkap kode error standar, deskripsi user-facing, dan status HTTP.
