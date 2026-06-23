# 🛍️ Katalog — Barang API Contract

---
### 🧭 Navigasi Cepat
[🏠 Utama](../../README.md) | [🏛️ Arsitektur](../../architecture.md) | [🛠️ Deployment](../../deployment.md) | [💳 Midtrans](../../pembayaran.md) | [📦 Biteship](../../biteship.md) | [📡 API Contract](../overview.md) | [🗄️ Database](../../database/erd.md) | [🔒 Keamanan](../../security/README.md)
---

Pusat integrasi manajemen produk/barang katalog untuk customer, kasir, dan admin.

---

## 🧭 Daftar Endpoint Barang

*   [**GET /api/v1/customer/barang**](#get-apiv1customerbarang) - Mengambil daftar katalog produk untuk Customer
*   [**GET /api/v1/admin/barang**](#get-apiv1adminbarang) - Mengambil daftar seluruh produk (Admin)
*   [**POST /api/v1/admin/barang**](#post-apiv1adminbarang) - Menambah barang katalog baru (Admin)
*   [**GET /api/v1/admin/barang/detail/:public_id**](#get-apiv1adminbarangdetailpublic_id) - Mengambil detail produk (Admin)
*   [**PUT /api/v1/admin/barang/:public_id**](#put-apiv1adminbarangpublic_id) - Memperbarui data produk (Admin)
*   [**DELETE /api/v1/admin/barang/:public_id**](#delete-apiv1adminbarangpublic_id) - Menghapus produk dari katalog (Admin)
*   [**POST /api/v1/admin/barang/upload**](#post-apiv1adminbarangupload) - Mengunggah gambar produk ke MinIO (Admin)
*   [**GET /api/v1/scan/:kode_barcode**](#get-apiv1scankode_barcode) - Mencari produk via scan barcode (Public)
*   [**POST /api/v1/kasir/transaksi/produk**](#post-apiv1kasirtransaksiproduk) - Mencari produk untuk POS (Kasir)

---

## GET /api/v1/customer/barang

Mendapatkan daftar katalog produk aktif yang dapat dibeli oleh customer.

*   **Autentikasi:** Wajib (Role: `Customer`)
*   **Header Wajib:** `Authorization: Bearer <access_token>`

### Response (200 OK)
```json
{
  "status": "success",
  "data": [
    {
      "id_barang": 1,
      "public_id": "9e3c8162-...",
      "nama_barang": "Teh Kotak Sosro 300ml",
      "gambar_barang": "https://minio-url/bucket/teh_kotak.jpg",
      "harga": 5000,
      "diskon": 10
    }
  ]
}
```

---

## GET /api/v1/admin/barang

Mendapatkan list seluruh produk tanpa filter status untuk keperluan dashboard.

*   **Autentikasi:** Wajib (Role: `Admin`)

---

## POST /api/v1/admin/barang

Menambahkan produk master baru ke dalam database katalog.

*   **Autentikasi:** Wajib (Role: `Admin`)

### Request Payload
```json
{
  "nama_barang": "Kopi Tubruk 100g",
  "id_kategori": 2,
  "id_satuan": 1,
  "id_diskon": null,
  "deskripsi": "Kopi lokal kualitas premium tanpa ampas"
}
```

---

## GET /api/v1/admin/barang/detail/:public_id

Mengambil detail lengkap informasi barang menggunakan UUID `public_id`.

*   **Autentikasi:** Wajib (Role: `Admin`)

---

## PUT /api/v1/admin/barang/:public_id

Mengubah/memperbarui informasi produk master berdasarkan `public_id`.

*   **Autentikasi:** Wajib (Role: `Admin`)

---

## DELETE /api/v1/admin/barang/:public_id

Menghapus data produk master dari database berdasarkan `public_id`.

*   **Autentikasi:** Wajib (Role: `Admin`)

---

## POST /api/v1/admin/barang/upload

Mengunggah file foto produk ke dalam cloud storage MinIO.

*   **Autentikasi:** Wajib (Role: `Admin`)
*   **Tipe Request:** `multipart/form-data`
*   **Field Form:** `gambar` (Binary File)

---

## GET /api/v1/scan/:kode_barcode

Mendapatkan informasi dasar produk secara instan berdasarkan kode barcode (biasanya dipakai untuk kasir scanning barang fisik).

*   **Autentikasi:** Tidak Ada (Public / POS Scan)

### Response (200 OK)
```json
{
  "status": "success",
  "data": {
    "nama_barang": "Teh Kotak Sosro 300ml",
    "harga": 5000,
    "stok": 120
  }
}
```

---

## POST /api/v1/kasir/transaksi/produk

Pencarian produk cepat untuk kasir kas register POS dengan mencocokkan nama barang atau kode barcode.

*   **Autentikasi:** Wajib (Role: `Kasir`)

### Request Payload
```json
{
  "keyword": "Teh Sosro"
}
```
