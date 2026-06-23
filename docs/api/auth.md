# 🔑 Auth API Contract

---
### 🧭 Navigasi Cepat
[🏠 Utama](../README.md) | [🏛️ Arsitektur](../architecture.md) | [🛠️ Deployment](../deployment.md) | [💳 Midtrans](../pembayaran.md) | [📦 Biteship](../biteship.md) | [📡 API Contract](overview.md) | [🗄️ Database](../database/erd.md) | [🔒 Keamanan](../security/README.md)
---

Pusat integrasi endpoint autentikasi dan otorisasi pengguna sistem MANTRA.

---

## 🧭 Daftar Endpoint Auth

*   [**POST /api/v1/login**](#post-apiv1login) - Masuk menggunakan username/email
*   [**POST /api/v1/register**](#post-apiv1register) - Mendaftar akun customer baru
*   [**POST /api/v1/auth/refresh**](#post-apiv1authrefresh) - Memperbarui token akses (access token)
*   [**POST /api/v1/logout**](#post-apiv1logout) - Keluar dan menonaktifkan refresh token
*   [**PUT /api/v1/change-password**](#put-apiv1change-password) - Mengubah password pengguna

---

## POST /api/v1/login

Melakukan login akun pengguna menggunakan username atau email.

*   **Autentikasi:** Tidak Ada (Public)
*   **Header Wajib:**
    *   `Content-Type: application/json`
    *   `X-Client-Type`: `flutter` (untuk apps) atau `nextjs` (untuk web admin)

### Request Payload
```json
{
  "username": "johndoe",
  "password": "password123"
}
```

### Response (Client Flutter App)
```json
{
  "status": "success",
  "message": "Login berhasil",
  "data": {
    "access_token": "eyJhbG...",
    "refresh_token": "a1b2c3d4...",
    "token_type": "Bearer",
    "expires_in": 1800,
    "user": {
      "id_user": 1,
      "username": "johndoe",
      "email": "john@example.com",
      "nama_lengkap": "John Doe",
      "role": "Customer",
      "profile_id": 1
    }
  }
}
```

### Response (Client Next.js Web Admin)
Response JSON sama dengan Flutter namun **tanpa** field `access_token` dan `refresh_token` di dalam body. Token secara otomatis di-set ke dalam browser client melalui header `Set-Cookie` (menggunakan flag `httpOnly`, `Secure`, dan `SameSite`).

---

## POST /api/v1/register

Membuat akun customer baru di database.

*   **Autentikasi:** Tidak Ada (Public)
*   **Header Wajib:** `Content-Type: application/json`

### Request Payload
```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "password123",
  "konfirmasi_password": "password123",
  "nama_lengkap": "John Doe",
  "no_telp": "08123456789"
}
```

### Response (201 Created)
```json
{
  "status": "success",
  "message": "Registrasi berhasil",
  "data": {
    "id_user": 2,
    "username": "johndoe",
    "email": "john@example.com",
    "nama_lengkap": "John Doe",
    "no_telp": "08123456789",
    "role": "customer"
  }
}
```

### Validasi Logika Bisnis
*   Password wajib minimal memiliki panjang **8 karakter**.
*   `konfirmasi_password` wajib bernilai sama persis dengan `password`.
*   Format input `email` harus merupakan email yang valid.
*   `username` dan `email` tidak boleh sama dengan yang sudah terdaftar (Unique Constraint).

---

## POST /api/v1/auth/refresh

Memperbarui `access_token` yang telah kedaluwarsa dengan menggunakan token penyegar (`refresh_token`).

*   **Autentikasi:** Tidak Ada (Public)
*   **Header Wajib:** 
    *   `Content-Type: application/json`
    *   `X-Client-Type`: `flutter` atau `nextjs`

### Request Payload (Khusus Flutter)
```json
{
  "refresh_token": "a1b2c3d4..."
}
```
*Catatan: Client Next.js mengirimkan refresh token otomatis via Cookie.*

### Response (Client Flutter App)
```json
{
  "status": "success",
  "message": "Token berhasil diperbarui",
  "data": {
    "access_token": "eyJhbG...",
    "expires_in": 1800
  }
}
```

### Response (Client Next.js Web Admin)
Mengembalikan status sukses dan secara otomatis memperbarui cookie access token yang baru.

---

## POST /api/v1/logout

Mengakhiri sesi pengguna dan menghapus/menonaktifkan status token aktif di database.

*   **Autentikasi:** Wajib (Semua Role)
*   **Header Wajib:**
    *   `Authorization: Bearer <access_token>` (Flutter)
    *   `X-Client-Type`: `flutter` atau `nextjs`

### Request Payload (Khusus Flutter)
```json
{
  "refresh_token": "a1b2c3d4..."
}
```
*Catatan: Client Next.js cukup menembak endpoint ini dan server otomatis akan menghapus cookie.*

### Response
```json
{
  "status": "success",
  "message": "Logout berhasil"
}
```

---

## PUT /api/v1/change-password

Mengubah password pengguna yang sedang login.

*   **Autentikasi:** Wajib (Semua Role)
*   **Header Wajib:** `Authorization: Bearer <access_token>`

### Request Payload
```json
{
  "password_lama": "oldpass123",
  "password_baru": "newpass456",
  "konfirmasi_password": "newpass456"
}
```

### Response
```json
{
  "status": "success",
  "message": "Password berhasil diubah"
}
```

> [!NOTE]
> Setelah password berhasil diubah, seluruh sesi/refresh token pengguna pada device lain otomatis akan di-revoke (logout dari semua perangkat demi keamanan).
