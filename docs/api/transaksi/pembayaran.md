# Transaksi — Pembayaran API

Seluruh endpoint Auth (Kasir). Digunakan untuk alur POS di aplikasi Kasir.

---

## GET /kasir/transaksi/checkout

Mendapatkan ringkasan checkout pesanan offline.

---

## POST /kasir/transaksi/bayar/tunai

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

## POST /kasir/transaksi/bayar/non-tunai

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

## PATCH /kasir/transaksi/item/update

Mengupdate quantity item di pesanan offline sebelum checkout.

**Request:**
```json
{
  "id_detail_pesanan": 1,
  "jumlah": 3
}
```

---

## GET /kasir/dashboard

Mendapatkan data dashboard kasir (ringkasan penjualan hari ini).

---

## GET /kasir/laporan

Mendapatkan laporan ringkasan.

---

## GET /kasir/laporan/produk/:id_produk

Mendapatkan detail laporan per produk.

---

## GET /kasir/laporan/produk/:id_produk/:id_pesanan

Mendapatkan detail pesanan dari laporan produk.
