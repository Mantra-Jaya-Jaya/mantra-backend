# Transaksi — Pembayaran API

Seluruh endpoint Auth (Kasir). Digunakan untuk alur POS di aplikasi Kasir.

---

## GET /api/v1/kasir/transaksi/checkout

Mendapatkan ringkasan checkout pesanan offline.

---

## POST /api/v1/kasir/transaksi/bayar/tunai

Memproses pembayaran tunai.

**Request:**
```json
{
  "total_bayar": 200000,
  "id_pesanan": 1
}
```

**Response:**
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

Memproses pembayaran non-tunai (QRIS, E-Wallet, VA, Kartu).

**Request:**
```json
{
  "payment_type": "qris",
  "id_pesanan": 1
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Pembayaran non-tunai berhasil diproses",
  "data": {
    "order_id_midtrans": "ORDER-123",
    "payment_type": "qris",
    "status_transaksi": "pending"
  }
}
```

---

## PATCH /api/v1/kasir/transaksi/item/update

Mengupdate quantity item di pesanan offline sebelum checkout.

**Request:**
```json
{
  "id_detail_pesanan": 1,
  "jumlah": 3
}
```

---

## GET /api/v1/kasir/dashboard

Mendapatkan data dashboard kasir (ringkasan penjualan hari ini).

---

## GET /api/v1/kasir/laporan

Mendapatkan laporan ringkasan.

---

## GET /api/v1/kasir/laporan/produk/:public_id

Mendapatkan detail laporan per produk.

---

## GET /api/v1/kasir/laporan/produk/:public_id/:pesanan_id

Mendapatkan detail pesanan dari laporan produk.
