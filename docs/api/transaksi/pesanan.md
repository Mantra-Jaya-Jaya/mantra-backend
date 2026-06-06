# Transaksi — Pesanan API

---

## GET /customer/pesanan (Customer) / GET /kasir/pesanan (Kasir)

Auth. Mendapatkan daftar pesanan.

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id_pesanan": 1,
      "public_id": "uuid-...",
      "total_pembayaran": 150000,
      "status_pesanan": "diproses",
      "tipe_pesanan": "online",
      "tanggal_pesanan": "2026-06-05T10:00:00Z"
    }
  ]
}
```

---

## POST /customer/pesanan/checkout

Auth (Customer). Membuat pesanan baru dari keranjang.

**Request:**
```json
{
  "id_alamat": 1,
  "catatan": "Tolong dibungkus rapih"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Pesanan berhasil dibuat",
  "data": {
    "id_pesanan": 1,
    "public_id": "uuid-...",
    "total_pembayaran": 150000
  }
}
```

---

## GET /customer/pesanan/:id_pesanan (Customer) / GET /kasir/pesanan/:id_order (Kasir)

Auth. Mendapatkan detail pesanan.

**Response:**
```json
{
  "status": "success",
  "data": {
    "id_pesanan": 1,
    "public_id": "uuid-...",
    "status_pesanan": "diproses",
    "detail_pesanan": [
      {
        "id_barang": 1,
        "nama_barang": "Produk A",
        "jumlah": 2,
        "harga_satuan": 50000,
        "subtotal": 100000
      }
    ],
    "total_pembayaran": 150000
  }
}
```

---

## PATCH /customer/pesanan/:id_pesanan/batal

Auth (Customer). Membatalkan pesanan.

**Response:**
```json
{
  "status": "success",
  "message": "Pesanan berhasil dibatalkan"
}
```

---

## GET /customer/pesanan/:id_pesanan/lacak

Auth (Customer). Lihat status pengiriman. Detail di `pengantaran.md`.
