# 🛒 Keranjang API Contract

---
### 🧭 Navigasi Cepat
[🏠 Utama](../README.md) | [🏛️ Arsitektur](../architecture.md) | [🛠️ Deployment](../deployment.md) | [💳 Midtrans](../pembayaran.md) | [📦 Biteship](../biteship.md) | [📡 API Contract](overview.md) | [🗄️ Database](../database/erd.md) | [🔒 Keamanan](../security/README.md)
---

Pusat kontrol dan manajemen keranjang belanjaan (Shopping Cart) pelanggan.

---

## 🧭 Daftar Endpoint Keranjang

*   [**POST /api/v1/customer/keranjang**](#post-apiv1customerkeranjang) - Menambahkan item baru ke keranjang
*   [**PATCH /api/v1/customer/keranjang/:public_id**](#patch-apiv1customerkeranjangpublic_id) - Memperbarui jumlah (quantity) barang
*   [**DELETE /api/v1/customer/keranjang/:public_id**](#delete-apiv1customerkeranjangpublic_id) - Menghapus barang dari keranjang

---

## POST /api/v1/customer/keranjang

Menambahkan produk dengan spesifikasi tertentu ke keranjang belanja customer.

*   **Autentikasi:** Wajib (Role: `Customer`)
*   **Header Wajib:** `Authorization: Bearer <access_token>`

### Request Payload
```json
{
  "id_spesifikasi_barang": 1,
  "quantity": 2
}
```

### Response (200 OK)
```json
{
  "status": "success",
  "message": "Item berhasil ditambahkan ke keranjang"
}
```

---

## PATCH /api/v1/customer/keranjang/:public_id

Memperbarui jumlah kuantitas (`quantity`) item yang sudah ada di dalam keranjang belanja.

*   **Autentikasi:** Wajib (Role: `Customer` & Verifikasi Kepemilikan Resource)
*   **Header Wajib:** `Authorization: Bearer <access_token>`
*   **Parameter URL:** `public_id` (UUID item keranjang)

### Request Payload
```json
{
  "quantity": 5
}
```

### Response (200 OK)
```json
{
  "status": "success",
  "message": "Jumlah item berhasil diperbarui"
}
```

---

## DELETE /api/v1/customer/keranjang/:public_id

Menghapus salah satu item produk dari keranjang belanja.

*   **Autentikasi:** Wajib (Role: `Customer` & Verifikasi Kepemilikan Resource)
*   **Header Wajib:** `Authorization: Bearer <access_token>`
*   **Parameter URL:** `public_id` (UUID item keranjang)

### Response (200 OK)
```json
{
  "status": "success",
  "message": "Item berhasil dihapus dari keranjang"
}
```
