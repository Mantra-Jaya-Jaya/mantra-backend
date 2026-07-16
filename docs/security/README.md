# 🔒 Keamanan & RBAC (Role-Based Access Control)

---
### 🧭 Navigasi Cepat
[🏠 Utama](../README.md) | [🏛️ Arsitektur](../architecture.md) | [🛠️ Deployment](../deployment.md) | [💳 Midtrans](../pembayaran.md) | [📦 Biteship](../biteship.md) | [📡 API Contract](../api/overview.md) | [🗄️ Database](../database/erd.md) | [🔒 Keamanan](README.md)
---

Dokumentasi ini mencakup seluruh aspek keamanan yang diimplementasikan pada **Mantra Backend**, termasuk alur autentikasi token, validasi izin akses middleware, matriks peran, serta sistem standardisasi error.

---

## 🧭 Daftar Modul Keamanan

*   [**jwt.md**](jwt.md) - Alur autentikasi JWT token (Access Token & Refresh Token), struktur claims payload, mekanisme perpanjangan masa aktif (sliding expiration), serta pemisahan multi-client (Flutter & Next.js).
*   [**middleware.md**](middleware.md) - Dokumentasi sistem penyaringan request middleware:
    - `AuthMiddleware` (Validasi login pengguna).
    - `RoleMiddleware` (Validasi hak akses peran/role).
    - `OwnershipMiddleware` (Validasi kepemilikan resource data privat).
*   [**rbac.md**](rbac.md) - Matriks pemetaan Role-Based Access Control (RBAC) yang menjelaskan endpoint apa saja yang boleh diakses oleh masing-masing tipe pengguna (Customer, Kasir, Kurir, Admin).
*   [**errors.md**](errors.md) - Standardisasi format response kesalahan (error) dan daftar kode error global.
