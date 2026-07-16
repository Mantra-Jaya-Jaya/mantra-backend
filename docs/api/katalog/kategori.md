# 🗂️ Katalog — Kategori API Contract

---
### 🧭 Navigasi Cepat
[🏠 Utama](../../README.md) | [🏛️ Arsitektur](../../architecture.md) | [🛠️ Deployment](../../deployment.md) | [💳 Midtrans](../../pembayaran.md) | [📦 Biteship](../../biteship.md) | [📡 API Contract](../overview.md) | [🗄️ Database](../../database/erd.md) | [🔒 Keamanan](../../security/README.md)
---

Pusat pengelompokan produk/barang berdasarkan kategori (makanan, minuman, obat-obatan, dll.).

---

## 🧭 Daftar Endpoint Kategori

*   [**GET /api/v1/customer/kategori**](#get-apiv1customerkategori-customer--get-apiv1kasirkategori-kasir--get-apiv1adminkategori-admin) - Mendapatkan list semua Kategori (Semua Role)
*   [**POST /api/v1/admin/kategori**](#post-apiv1adminkategori) - Menambahkan kategori produk baru (Admin)
*   [**PUT /api/v1/admin/kategori/:public_id**](#put-apiv1adminkategoripublic_id) - Memperbarui nama/detail kategori (Admin)
*   [**DELETE /api/v1/admin/kategori/:public_id**](#delete-apiv1adminkategoripublic_id) - Menghapus kategori dari katalog (Admin)
*   [**POST /api/v1/admin/kategori/upload**](#post-apiv1adminkategoriupload) - Mengunggah icon kategori ke MinIO (Admin)

---

## GET /api/v1/customer/kategori (Customer) / GET /api/v1/kasir/kategori (Kasir) / GET /api/v1/admin/kategori (Admin)

Mendapatkan daftar kategori barang yang tersedia untuk penjelajahan katalog produk.

*   **Autentikasi:** Wajib (Semua Role)
*   **Header Wajib:** `Authorization: Bearer <access_token>`

### Response (200 OK)
```json
{
  "status": "success",
  "data": [
    {
      "id_kategori": 1,
      "nama_kategori": "Makanan Instan",
      "icon_kategori": "https://minio-url/bucket/kategori/makanan.png"
    }
  ]
}
```

---

## POST /api/v1/admin/kategori

Menambahkan kategori barang master baru.

*   **Autentikasi:** Wajib (Role: `Admin`)

### Request Payload
```json
{
  "nama_kategori": "Kebutuhan Rumah Tangga"
}
```

---

## PUT /api/v1/admin/kategori/:public_id

Mengubah nama atau deskripsi kategori barang berdasarkan `public_id`.

*   **Autentikasi:** Wajib (Role: `Admin`)

---

## DELETE /api/v1/admin/kategori/:public_id

Menghapus kategori barang berdasarkan `public_id`.

*   **Autentikasi:** Wajib (Role: `Admin`)

---

## POST /api/v1/admin/kategori/upload

Mengunggah gambar icon penanda kategori ke dalam cloud storage MinIO.

*   **Autentikasi:** Wajib (Role: `Admin`)
*   **Tipe Request:** `multipart/form-data`
*   **Field Form:** `icon` (Binary File)
