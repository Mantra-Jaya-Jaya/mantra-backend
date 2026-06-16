# Transaksi — Metode Pembayaran API

Auth (Customer & Admin). Kelola dan lihat metode pembayaran yang tersedia.

---

## GET /api/v1/customer/metode-pembayaran

Customer — daftar metode pembayaran aktif (urut berdasarkan `urutan`).

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "public_id": "uuid-...",
      "nama_metode": "QRIS",
      "kode_metode": "qris",
      "penyedia": "midtrans",
      "icon": "https://...",
      "urutan": 1
    }
  ]
}
```

---

## GET /api/v1/admin/metode-pembayaran

Admin — daftar semua metode (aktif dan non-aktif).

---

## POST /api/v1/admin/metode-pembayaran

Tambah metode pembayaran baru.

**Request:**
```json
{
  "nama_metode": "QRIS",
  "kode_metode": "qris",
  "penyedia": "midtrans",
  "icon": "https://...",
  "urutan": 1,
  "is_active": true
}
```

---

## PUT /api/v1/admin/metode-pembayaran/:public_id

Update metode pembayaran.

---

## DELETE /api/v1/admin/metode-pembayaran/:public_id

Hapus metode pembayaran.
