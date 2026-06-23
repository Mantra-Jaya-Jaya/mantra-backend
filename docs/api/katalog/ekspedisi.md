# 🚚 Katalog — Ekspedisi API Contract

---
### 🧭 Navigasi Cepat
[🏠 Utama](../../README.md) | [🏛️ Arsitektur](../../architecture.md) | [🛠️ Deployment](../../deployment.md) | [💳 Midtrans](../../pembayaran.md) | [📦 Biteship](../../biteship.md) | [📡 API Contract](../overview.md) | [🗄️ Database](../../database/erd.md) | [🔒 Keamanan](../../security/README.md)
---

Pusat CRUD manajemen kurir pihak ketiga (eksternal) dan sub-layanan pengirimannya.

---

## 🧭 Daftar Endpoint Ekspedisi

*   [**GET /api/v1/admin/ekspedisi**](#get-apiv1adminekspedisi) - Mengambil daftar seluruh ekspedisi & layanannya
*   [**POST /api/v1/admin/ekspedisi**](#post-apiv1adminekspedisi) - Menambahkan ekspedisi baru
*   [**PUT /api/v1/admin/ekspedisi/:public_id**](#put-apiv1adminekspedisipublic_id) - Memperbarui informasi ekspedisi
*   [**DELETE /api/v1/admin/ekspedisi/:public_id**](#delete-apiv1adminekspedisipublic_id) - Menghapus ekspedisi beserta layanannya
*   [**POST /api/v1/admin/ekspedisi/layanan**](#post-apiv1adminekspedisilayanan) - Menambahkan tipe layanan baru untuk ekspedisi
*   [**PUT /api/v1/admin/ekspedisi/layanan/:public_id**](#put-apiv1adminekspedisilayananpublic_id) - Memperbarui detail sub-layanan ekspedisi
*   [**DELETE /api/v1/admin/ekspedisi/layanan/:public_id**](#delete-apiv1adminekspedisilayananpublic_id) - Menghapus layanan ekspedisi

---

## GET /api/v1/admin/ekspedisi

Mengambil daftar kurir eksternal yang terdaftar di database beserta sub-layanan pengiriman mereka masing-masing.

*   **Autentikasi:** Wajib (Role: `Admin`)
*   **Header Wajib:** `Authorization: Bearer <access_token>`

### Response (200 OK)
```json
{
  "status": "success",
  "data": [
    {
      "public_id": "9e3c8162-...",
      "nama_ekspedisi": "J&T Express",
      "logo": "https://minio-url/bucket/logo_jnt.png",
      "deskripsi": "Layanan ekspedisi reguler nasional",
      "is_active": true,
      "layanan": [
        {
          "public_id": "8f3c8162-...",
          "nama_layanan": "EZ",
          "estimasi_min": 2,
          "estimasi_max": 4,
          "is_active": true
        }
      ]
    }
  ]
}
```

---

## POST /api/v1/admin/ekspedisi

Mendaftarkan perusahaan ekspedisi (kurir pihak ketiga) baru ke database.

*   **Autentikasi:** Wajib (Role: `Admin`)

### Request Payload
```json
{
  "nama_ekspedisi": "J&T Express",
  "logo": "https://minio-url/bucket/logo_jnt.png",
  "deskripsi": "Layanan ekspedisi reguler nasional",
  "kode_api": "jnt",
  "is_active": true
}
```

---

## PUT /api/v1/admin/ekspedisi/:public_id

Mengubah informasi utama dari perusahaan ekspedisi berdasarkan `public_id`.

*   **Autentikasi:** Wajib (Role: `Admin`)

---

## DELETE /api/v1/admin/ekspedisi/:public_id

Menghapus data ekspedisi. Penghapusan data ekspedisi secara otomatis akan menghapus seluruh data tipe layanan terkait (Cascade Delete).

*   **Autentikasi:** Wajib (Role: `Admin`)

---

## POST /api/v1/admin/ekspedisi/layanan

Menambahkan jenis layanan pengiriman baru untuk salah satu ekspedisi terdaftar.

*   **Autentikasi:** Wajib (Role: `Admin`)

### Request Payload
```json
{
  "id_ekspedisi": 1,
  "nama_layanan": "Super Fast",
  "deskripsi": "Layanan satu hari sampai tujuan",
  "estimasi_min": 1,
  "estimasi_max": 1,
  "is_active": true
}
```

---

## PUT /api/v1/admin/ekspedisi/layanan/:public_id

Mengubah/memperbarui informasi detail sub-layanan pengiriman berdasarkan `public_id` layanan.

*   **Autentikasi:** Wajib (Role: `Admin`)

---

## DELETE /api/v1/admin/ekspedisi/layanan/:public_id

Menghapus tipe layanan pengiriman tertentu berdasarkan `public_id`.

*   **Autentikasi:** Wajib (Role: `Admin`)
