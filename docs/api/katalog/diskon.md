# 🏷️ Katalog — Diskon (Promo) API Contract

---
### 🧭 Navigasi Cepat
[🏠 Utama](../../README.md) | [🏛️ Arsitektur](../../architecture.md) | [🛠️ Deployment](../../deployment.md) | [💳 Midtrans](../../pembayaran.md) | [📦 Biteship](../../biteship.md) | [📡 API Contract](../overview.md) | [🗄️ Database](../../database/erd.md) | [🔒 Keamanan](../../security/README.md)
---

Pusat integrasi data diskon, promo banner, dan program promosi toko.

---

## 🧭 Daftar Endpoint Diskon

*   [**GET /api/v1/customer/promo**](#get-apiv1customerpromo-customer--get-apiv1admindiskon-admin) - Mendapatkan promo aktif untuk Customer
*   [**GET /api/v1/admin/diskon**](#get-apiv1customerpromo-customer--get-apiv1admindiskon-admin) - Mendapatkan diskon aktif untuk Admin
*   [**GET /api/v1/admin/diskon/semua**](#get-apiv1admindiskonsemua) - Mengambil seluruh data diskon (Admin)
*   [**POST /api/v1/admin/diskon**](#post-apiv1admindiskon) - Menambahkan program diskon baru (Admin)
*   [**DELETE /api/v1/admin/diskon/:public_id**](#delete-apiv1admindiskonpublic_id) - Menghapus program diskon (Admin)
*   [**POST /api/v1/admin/diskon/upload**](#post-apiv1admindiskonupload) - Mengunggah banner promo ke MinIO (Admin)

---

## GET /api/v1/customer/promo (Customer) / GET /api/v1/admin/diskon (Admin)

Mengambil daftar program diskon atau promosi aktif yang saat ini sedang berlangsung (berdasarkan tanggal mulai dan selesai).

*   **Autentikasi:** Wajib (Role: `Customer` atau `Admin`)
*   **Header Wajib:** `Authorization: Bearer <access_token>`

### Response (200 OK)
```json
{
  "status": "success",
  "data": [
    {
      "id_diskon": 1,
      "nama_diskon": "Diskon Akhir Tahun Super Hemat",
      "besar_diskon": 20,
      "banner_diskon": "https://minio-url/bucket/promo_banner.png",
      "tgl_mulai": "2026-12-01",
      "tgl_selesai": "2026-12-31"
    }
  ]
}
```

---

## GET /api/v1/admin/diskon/semua

Mengambil seluruh daftar program diskon tanpa filter masa aktif (termasuk promo masa lalu maupun rancangan promo masa depan).

*   **Autentikasi:** Wajib (Role: `Admin`)

---

## POST /api/v1/admin/diskon

Menambahkan jenis/program promo diskon baru ke dalam master data.

*   **Autentikasi:** Wajib (Role: `Admin`)

### Request Payload
```json
{
  "nama_diskon": "Promo Imlek Gembira",
  "besar_diskon": 15,
  "tgl_mulai": "2026-02-10",
  "tgl_selesai": "2026-02-25"
}
```

---

## DELETE /api/v1/admin/diskon/:public_id

Menghapus program promo diskon menggunakan UUID `public_id`.

*   **Autentikasi:** Wajib (Role: `Admin`)

---

## POST /api/v1/admin/diskon/upload

Mengunggah gambar/banner promo diskon ke dalam cloud storage MinIO.

*   **Autentikasi:** Wajib (Role: `Admin`)
*   **Tipe Request:** `multipart/form-data`
*   **Field Form:** `banner` (Binary File)
