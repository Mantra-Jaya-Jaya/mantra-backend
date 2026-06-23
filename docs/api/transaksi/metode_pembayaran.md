# 💳 Transaksi — Metode Pembayaran API Contract

---
### 🧭 Navigasi Cepat
[🏠 Utama](../../README.md) | [🏛️ Arsitektur](../../architecture.md) | [🛠️ Deployment](../../deployment.md) | [💳 Midtrans](../../pembayaran.md) | [📦 Biteship](../../biteship.md) | [📡 API Contract](../overview.md) | [🗄️ Database](../../database/erd.md) | [🔒 Keamanan](../../security/README.md)
---

Pusat CRUD manajemen tipe/metode pembayaran yang didukung sistem (QRIS, VA BCA, Mandiri, Cash, dll).

---

## 🧭 Daftar Endpoint Metode Pembayaran

*   [**GET /api/v1/customer/metode-pembayaran**](#get-apiv1customermetode-pembayaran) - Mengambil metode pembayaran aktif untuk Customer
*   [**GET /api/v1/admin/metode-pembayaran**](#get-apiv1adminmetode-pembayaran) - Mengambil seluruh metode pembayaran (Admin)
*   [**POST /api/v1/admin/metode-pembayaran**](#post-apiv1adminmetode-pembayaran) - Menambahkan metode pembayaran baru (Admin)
*   [**PUT /api/v1/admin/metode-pembayaran/:public_id**](#put-apiv1adminmetode-pembayaranpublic_id) - Memperbarui status/urutan metode (Admin)
*   [**DELETE /api/v1/admin/metode-pembayaran/:public_id**](#delete-apiv1adminmetode-pembayaranpublic_id) - Menghapus metode pembayaran (Admin)

---

## GET /api/v1/customer/metode-pembayaran

Mengambil daftar opsi metode pembayaran aktif yang diurutkan berdasarkan prioritas tampilan (`urutan`).

*   **Autentikasi:** Wajib (Role: `Customer`)
*   **Header Wajib:** `Authorization: Bearer <access_token>`

### Response (200 OK)
```json
{
  "status": "success",
  "data": [
    {
      "public_id": "9e3c8162-...",
      "nama_metode": "QRIS (Gopay, OVO, ShopeePay)",
      "kode_metode": "qris",
      "penyedia": "midtrans",
      "icon": "https://minio-url/bucket/qris_logo.png",
      "urutan": 1
    }
  ]
}
```

---

## GET /api/v1/admin/metode-pembayaran

Mengambil daftar seluruh opsi metode pembayaran (termasuk yang dinonaktifkan) untuk panel administrasi.

*   **Autentikasi:** Wajib (Role: `Admin`)

---

## POST /api/v1/admin/metode-pembayaran

Menambahkan metode pembayaran baru yang didukung sistem.

*   **Autentikasi:** Wajib (Role: `Admin`)

### Request Payload
```json
{
  "nama_metode": "BCA Virtual Account",
  "kode_metode": "bca_va",
  "penyedia": "midtrans",
  "icon": "https://minio-url/bucket/bca_logo.png",
  "urutan": 2,
  "is_active": true
}
```

---

## PUT /api/v1/admin/metode-pembayaran/:public_id

Mengubah rincian data metode pembayaran (seperti merubah nama, icon, status aktif, atau urutan tampil) berdasarkan `public_id`.

*   **Autentikasi:** Wajib (Role: `Admin`)

---

## DELETE /api/v1/admin/metode-pembayaran/:public_id

Menghapus metode pembayaran dari database berdasarkan `public_id`.

*   **Autentikasi:** Wajib (Role: `Admin`)
