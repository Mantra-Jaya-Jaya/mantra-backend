# 🚚 Transaksi — Ongkir API Contract

---
### 🧭 Navigasi Cepat
[🏠 Utama](../../README.md) | [🏛️ Arsitektur](../../architecture.md) | [🛠️ Deployment](../../deployment.md) | [💳 Midtrans](../../pembayaran.md) | [📦 Biteship](../../biteship.md) | [📡 API Contract](../overview.md) | [🗄️ Database](../../database/erd.md) | [🔒 Keamanan](../../security/README.md)
---

Endpoint kalkulasi tarif pengiriman (cek ongkos kirim) secara real-time via Biteship.

---

## POST /api/v1/customer/ongkir/cek

Menghitung estimasi ongkos kirim dari berbagai perusahaan ekspedisi eksternal berdasarkan alamat pengiriman customer dan berat total barang belanjaan.

*   **Autentikasi:** Wajib (Role: `Customer`)
*   **Header Wajib:** `Authorization: Bearer <access_token>`

### Request Payload
```json
{
  "id_alamat": "9e3c8162-...",
  "items": [
    {
      "id_spesifikasi_barang": 1,
      "quantity": 2
    }
  ]
}
```

### Response (200 OK)
```json
{
  "status": "success",
  "data": [
    {
      "nama_ekspedisi": "J&T Express",
      "layanan": [
        {
          "id_ekspedisi_layanan": 1,
          "nama_layanan": "EZ",
          "ongkir": 12000,
          "estimasi": "2-3 hari"
        }
      ]
    }
  ]
}
```
