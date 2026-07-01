# 💳 Transaksi — Pembayaran (POS Kasir) API Contract

---
### 🧭 Navigasi Cepat
[🏠 Utama](../../README.md) | [🏛️ Arsitektur](../../architecture.md) | [🛠️ Deployment](../../deployment.md) | [💳 Midtrans](../../pembayaran.md) | [📦 Biteship](../../biteship.md) | [📡 API Contract](../overview.md) | [🗄️ Database](../../database/erd.md) | [🔒 Keamanan](../../security/README.md)
---

Seluruh endpoint transaksi kasir offline/Point of Sales (POS) untuk mengelola checkout, pembayaran tunai/non-tunai, dan laporan ringkasan kasir.

---

## 🧭 Daftar Endpoint Pembayaran Kasir

*   [**GET /api/v1/kasir/transaksi/checkout**](#get-apiv1kasirtransaksicheckout) - Mengambil rincian checkout POS
*   [**POST /api/v1/kasir/transaksi/bayar/tunai**](#post-apiv1kasirtransaksibayartunai) - Memproses transaksi tunai (Cash)
*   [**POST /api/v1/kasir/transaksi/bayar/non-tunai**](#post-apiv1kasirtransaksibayarnon-tunai) - Memproses transaksi non-tunai (Midtrans)
*   [**PATCH /api/v1/kasir/transaksi/item/update**](#patch-apiv1kasirtransaksiitemupdate) - Memperbarui kuantitas item keranjang kasir
*   [**GET /api/v1/kasir/dashboard**](#get-apiv1kasirdashboard) - Mengambil ringkasan penjualan kasir hari ini
*   [**GET /api/v1/kasir/laporan**](#get-apiv1kasirlaporan) - Mengambil ringkasan laporan performa produk
*   [**GET /api/v1/kasir/laporan/produk/:public_id**](#get-apiv1kasirlaporanprodukpublic_id) - Detail laporan penjualan per produk
*   [**GET /api/v1/kasir/laporan/produk/:public_id/:pesanan_id**](#get-apiv1kasirlaporanprodukpublic_idpesanan_id) - Detail pesanan terkait dari laporan produk

---

## GET /api/v1/kasir/transaksi/checkout

Mengambil ringkasan pesanan offline Point of Sale (POS) yang sedang aktif sebelum dibayarkan.

*   **Autentikasi:** Wajib (Role: `Kasir`)
*   **Header Wajib:** `Authorization: Bearer <access_token>`

---

## POST /api/v1/kasir/transaksi/bayar/tunai

Memproses pembayaran fisik menggunakan uang tunai (cash) di meja kasir.

*   **Autentikasi:** Wajib (Role: `Kasir`)

### Request Payload
```json
{
  "total_bayar": 200000,
  "id_pesanan": 1
}
```

### Response (200 OK)
```json
{
  "status": "success",
  "message": "Pembayaran tunai berhasil",
  "data": {
    "kembalian": 50000
  }
}
```

---

## POST /api/v1/kasir/transaksi/bayar/non-tunai

Memproses pembayaran non-tunai (seperti QRIS, E-Wallet, Kartu Kredit, atau Virtual Account) melalui integrasi Midtrans Core API.

*   **Autentikasi:** Wajib (Role: `Kasir`)

### Request Payload
```json
{
  "payment_type": "qris",
  "id_pesanan": 1
}
```

### Response (200 OK)
```json
{
  "status": "success",
  "message": "Pembayaran non-tunai berhasil diproses",
  "data": {
    "order_id_midtrans": "MNT-POS-10-178000000",
    "payment_type": "qris",
    "status_transaksi": "pending"
  }
}
```

---

## PATCH /api/v1/kasir/transaksi/item/update

Mengubah kuantitas item produk yang dimasukkan ke keranjang kasir POS sebelum checkout.

*   **Autentikasi:** Wajib (Role: `Kasir`)

### Request Payload
```json
{
  "id_detail_pesanan": 1,
  "jumlah": 3
}
```

---

## GET /api/v1/kasir/dashboard

Mengambil total omset harian, jumlah transaksi sukses, dan data ringkasan kinerja kasir hari ini.

*   **Autentikasi:** Wajib (Role: `Kasir`)

---

## GET /api/v1/kasir/laporan

Mengambil ringkasan laporan produk terlaris (top selling products) dan data statistik performa produk.

*   **Autentikasi:** Wajib (Role: `Kasir`)

---

## GET /api/v1/kasir/laporan/produk/:public_id

Mengambil laporan statistik penjualan mendalam untuk satu produk spesifik berdasarkan `public_id`.

*   **Autentikasi:** Wajib (Role: `Kasir`)

---

## GET /api/v1/kasir/laporan/produk/:public_id/:pesanan_id

Mengambil detail faktur pesanan offline berdasarkan ID pesanan yang terdaftar pada riwayat laporan produk.

*   **Autentikasi:** Wajib (Role: `Kasir`)
