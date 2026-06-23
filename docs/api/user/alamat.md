# 📍 User — Alamat API Contract

---
### 🧭 Navigasi Cepat
[🏠 Utama](../../README.md) | [🏛️ Arsitektur](../../architecture.md) | [🛠️ Deployment](../../deployment.md) | [💳 Midtrans](../../pembayaran.md) | [📦 Biteship](../../biteship.md) | [📡 API Contract](../overview.md) | [🗄️ Database](../../database/erd.md) | [🔒 Keamanan](../../security/README.md)
---

Pusat CRUD pengelolaan daftar alamat pengiriman customer (Rumah, Kantor, dll.) beserta data koordinat GPS untuk keperluan perhitungan ongkos kirim.

---

## 🧭 Daftar Endpoint Alamat

*   [**GET /api/v1/customer/alamat**](#get-apiv1customeralamat) - Mengambil semua alamat pengiriman aktif
*   [**POST /api/v1/customer/alamat**](#post-apiv1customeralamat) - Menambahkan alamat pengiriman baru
*   [**PUT /api/v1/customer/alamat/:public_id**](#put-apiv1customeralamatpublic_id) - Memperbarui detail alamat
*   [**DELETE /api/v1/customer/alamat/:public_id**](#delete-apiv1customeralamatpublic_id) - Menghapus alamat dari daftar

---

## GET /api/v1/customer/alamat

Mengambil daftar seluruh alamat pengiriman milik customer yang sedang login.

*   **Autentikasi:** Wajib (Role: `Customer`)
*   **Header Wajib:** `Authorization: Bearer <access_token>`

### Response (200 OK)
```json
{
  "status": "success",
  "data": [
    {
      "id_alamat": 1,
      "public_id": "9e3c8162-...",
      "nama_penerima": "John Doe",
      "label_alamat": "Rumah Utama",
      "no_telp_penerima": "08123456789",
      "alamat_lengkap": "Jl. Prof Soedarto SH No. 1, Tembalang, Semarang",
      "is_utama": true
    }
  ]
}
```

---

## POST /api/v1/customer/alamat

Menambahkan alamat baru untuk pengiriman barang belanjaan online.

*   **Autentikasi:** Wajib (Role: `Customer`)

### Request Payload
```json
{
  "nama_penerima": "John Doe",
  "label_alamat": "Kantor Cabang",
  "no_telp_penerima": "08123456789",
  "alamat_lengkap": "Jl. Gajah Mada No. 100, Semarang",
  "latitude": -6.9822,
  "longitude": 110.4223,
  "catatan_lokasi": "Gedung A Lantai Dasar, samping pos satpam",
  "is_utama": false
}
```

### Response (201 Created)
```json
{
  "status": "success",
  "message": "Alamat baru berhasil ditambahkan"
}
```

---

## PUT /api/v1/customer/alamat/:public_id

Mengubah/memperbarui data alamat pengantaran berdasarkan `public_id` alamat.

*   **Autentikasi:** Wajib (Role: `Customer` & Verifikasi Kepemilikan)

---

## DELETE /api/v1/customer/alamat/:public_id

Menghapus salah satu alamat pengantaran dari profil berdasarkan `public_id`.

*   **Autentikasi:** Wajib (Role: `Customer` & Verifikasi Kepemilikan)
