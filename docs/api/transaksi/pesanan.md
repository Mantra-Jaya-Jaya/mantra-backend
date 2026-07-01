# 📦 Transaksi — Pesanan API Contract

---
### 🧭 Navigasi Cepat
[🏠 Utama](../../README.md) | [🏛️ Arsitektur](../../architecture.md) | [🛠️ Deployment](../../deployment.md) | [💳 Midtrans](../../pembayaran.md) | [📦 Biteship](../../biteship.md) | [📡 API Contract](../overview.md) | [🗄️ Database](../../database/erd.md) | [🔒 Keamanan](../../security/README.md)
---

Pusat manajemen pembuatan pesanan, riwayat belanja, detail rincian belanjaan, pembatalan pesanan, dan pelacakan pengiriman untuk customer, kasir, dan kurir.

---

## 🧭 Daftar Endpoint Pesanan

*   [**GET /api/v1/customer/pesanan**](#get-apiv1customerpesanan-customer--get-apiv1kasirpesanan-kasir) - Mengambil riwayat transaksi belanja (Customer / Kasir)
*   [**POST /api/v1/customer/pesanan/checkout**](#post-apiv1customerpesanancheckout) - Checkout & buat pesanan online (Customer)
*   [**GET /api/v1/customer/pesanan/:public_id**](#get-apiv1customerpesananpublic_id-customer--get-apiv1kasirpesananpublic_id-kasir--get-apiv1kurirpesananpublic_id-kurir) - Mengambil rincian detail pesanan lengkap (Customer / Kasir / Kurir)
*   [**PATCH /api/v1/customer/pesanan/:public_id/batal**](#patch-apiv1customerpesananpublic_idbatal) - Membatalkan pesanan (Customer)

---

## GET /api/v1/customer/pesanan (Customer) / GET /api/v1/kasir/pesanan (Kasir)

Mengambil daftar riwayat pesanan (baik pesanan online maupun transaksi fisik offline kasir).

*   **Autentikasi:** Wajib (Role: `Customer` atau `Kasir`)
*   **Header Wajib:** `Authorization: Bearer <access_token>`

### Response (200 OK)
```json
{
  "status": "success",
  "data": [
    {
      "id_pesanan": 1,
      "public_id": "9e3c8162-...",
      "total_pembayaran": 150000,
      "status_pesanan": "Diproses",
      "tipe_pesanan": "Online",
      "tanggal_pesanan": "2026-06-05T10:00:00Z"
    }
  ]
}
```

---

## POST /api/v1/customer/pesanan/checkout

Membuat pesanan online baru dari isi keranjang belanja aktif customer dan men-generate link pembayaran Midtrans Snap.

*   **Autentikasi:** Wajib (Role: `Customer`)
*   **Header Wajib:** `Authorization: Bearer <access_token>`

### Request Payload
```json
{
  "id_alamat": "9e3c8162-...",
  "id_ekspedisi": 1,
  "id_layanan_ekspedisi": 1,
  "ongkos_kirim": 20000,
  "catatan": "Tolong dibungkus rapi",
  "id_metode_pembayaran": 2
}
```

### Response (201 Created)
```json
{
  "status": "success",
  "message": "Pesanan berhasil dibuat",
  "data": {
    "id_pesanan": "9e3c8162-...",
    "total_bayar": 172000,
    "ongkos_kirim": 20000,
    "pajak": 15000,
    "midtrans_token": "3f0d242e-...",
    "redirect_url": "https://app.sandbox.midtrans.com/snap/v4/redirection/..."
  }
}
```

---

## GET /api/v1/customer/pesanan/:public_id (Customer) / GET /api/v1/kasir/pesanan/:public_id (Kasir) / GET /api/v1/kurir/pesanan/:public_id (Kurir)

Mengambil detail rincian produk, alamat pengantaran, identitas kurir, serta rincian pembayaran untuk satu transaksi spesifik.

*   **Autentikasi:** Wajib (Semua Role terkait)
*   **Parameter URL:** `public_id` (UUID Pesanan)

### Response (200 OK)
```json
{
  "status": "success",
  "message": "Detail pesanan berhasil diambil",
  "data": {
    "no_pesanan": "9e3c8162-...",
    "status": "Diproses",
    "tanggal_pesan": "2026-06-05T10:00:00Z",
    "items": [
      {
        "id_barang": 1,
        "nama_barang": "Teh Sosro Kotak 300ml",
        "varian": "Ukuran: Reguler",
        "jumlah": 2,
        "harga_satuan": 5000,
        "gambar": "http://minio-url/bucket/products/teh.png"
      }
    ],
    "tujuan_pengantaran": {
      "nama_penerima": "Surya Pratama",
      "alamat_lengkap": "Jl. Prof Soedarto SH, Tembalang, Semarang"
    },
    "kurir": {
      "nama_kurir": "Ricardo Holahilo",
      "plat_nomor": "H 1234 AB",
      "ekspedisi": "Internal Toko",
      "foto_kurir": "http://minio-url/bucket/profiles/avatar.png"
    },
    "rincian_pembayaran": {
      "subtotal_items": 10000,
      "ongkir": 20000,
      "biaya_proteksi": 0,
      "total": 30000
    }
  }
}
```

---

## PATCH /api/v1/customer/pesanan/:public_id/batal

Membatalkan transaksi pesanan online yang belum terbayar/terkirim.

*   **Autentikasi:** Wajib (Role: `Customer` & Verifikasi Kepemilikan)

### Response (200 OK)
```json
{
  "status": "success",
  "message": "Pesanan berhasil dibatalkan"
}
```

---

## GET /api/v1/customer/pesanan/:public_id/lacak

*   **Info:** Digunakan untuk melacak status pengiriman kurir. Skema detail parameter response silakan merujuk pada file [pengantaran.md](../pengantaran.md).
