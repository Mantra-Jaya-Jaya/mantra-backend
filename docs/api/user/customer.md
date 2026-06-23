# 👤 User — Customer API Contract

---
### 🧭 Navigasi Cepat
[🏠 Utama](../../README.md) | [🏛️ Arsitektur](../../architecture.md) | [🛠️ Deployment](../../deployment.md) | [💳 Midtrans](../../pembayaran.md) | [📦 Biteship](../../biteship.md) | [📡 API Contract](../overview.md) | [🗄️ Database](../../database/erd.md) | [🔒 Keamanan](../../security/README.md)
---

Manajemen profil, foto avatar, dan pembaharuan informasi akun customer.

---

## 🧭 Daftar Endpoint Customer Account

*   [**GET /api/v1/customer/profil**](#get-apiv1customerprofil) - Mengambil data profil customer
*   [**PUT /api/v1/customer/akun**](#put-apiv1customerakun) - Memperbarui informasi akun customer

---

## GET /api/v1/customer/profil

Mengambil rincian informasi pengguna/customer yang saat ini sedang login.

*   **Autentikasi:** Wajib (Role: `Customer`)
*   **Header Wajib:** `Authorization: Bearer <access_token>`

### Response (200 OK)
```json
{
  "status": "success",
  "data": {
    "id_user": 1,
    "public_id": "9e3c8162-...",
    "username": "johndoe",
    "email": "john@example.com",
    "nama_lengkap": "John Doe",
    "foto_profil": "https://minio-url/bucket/profiles/johndoe.png",
    "no_telp": "08123456789"
  }
}
```

---

## PUT /api/v1/customer/akun

Mengubah detail informasi kontak akun pengguna yang sedang login.

*   **Autentikasi:** Wajib (Role: `Customer`)
*   **Header Wajib:** `Authorization: Bearer <access_token>`

### Request Payload
```json
{
  "nama_lengkap": "John Doe Updated",
  "email": "john.new@example.com",
  "no_telp": "08123456780"
}
```

### Response (200 OK)
```json
{
  "status": "success",
  "message": "Profil berhasil diperbarui"
}
```
