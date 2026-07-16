# Transaksi — Payment Notification API

Public (Webhook). Endpoint untuk Midtrans mengirim notifikasi status pembayaran.

---

## POST /api/v1/payment/notification

Menerima notifikasi dari Midtrans tentang perubahan status transaksi.

Verifikasi dilakukan dengan SHA-512 hash dari `order_id + status_code + gross_amount + server_key`.

**Request:**

```json
{
  "transaction_time": "2026-06-09 12:00:00",
  "transaction_id": "trx-midtrans-001",
  "transaction_status": "settlement",
  "status_code": "200",
  "signature_key": "sha512hash...",
  "payment_type": "qris",
  "order_id": "MID-20260609-227",
  "merchant_id": "M913788146",
  "gross_amount": "24270556.00",
  "fraud_status": "accept",
  "currency": "IDR"
}
```

**Status yang ditangani:**

| transaction_status | Aksi |
|--------------------|------|
| `settlement` | Update status pesanan → "Dikemas" (jika pesanan Online) atau "Selesai" (jika pesanan POS/offline), catat waktu bayar |
| `capture` | Sama dengan settlement (kartu kredit) |
| `pending` | Tidak ada aksi (menunggu pembayaran) |
| `deny` | Update status transaksi → deny |
| `cancel` | Update status transaksi → cancel |
| `expire` | Update status transaksi → expire |

**Response:**

```json
{
  "status": "success",
  "message": "Notifikasi berhasil diproses"
}
```
