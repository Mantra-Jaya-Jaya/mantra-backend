# Transaksi — Ongkir API

Auth (Customer). Cek ongkos kirim via Biteship API.

---

## POST /api/v1/customer/ongkir/cek

Hitung ongkos kirim dari berbagai ekspedisi berdasarkan alamat tujuan dan berat barang.

**Request:**
```json
{
  "id_alamat": "uuid-...",
  "items": [
    {
      "id_spesifikasi_barang": 1,
      "quantity": 2
    }
  ]
}
```

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "nama_ekspedisi": "Ninja Xpress",
      "layanan": [
        {
          "id_ekspedisi_layanan": 1,
          "nama_layanan": "REG",
          "ongkir": 15000,
          "estimasi": "2-4 hari"
        }
      ]
    }
  ]
}
```
