# Pengantaran API

---

## GET /customer/pesanan/:id_pesanan/lacak

Auth (Customer). Melacak status pengiriman pesanan.

**Response:**
```json
{
  "status": "success",
  "data": {
    "id_pesanan": 1,
    "status_pesanan": "dikirim",
    "pengantaran": {
      "status": "dalam_perjalanan",
      "kurir": "Kurir A",
      "last_latitude": -6.2,
      "last_longitude": 106.8,
      "waktu_pickup": "2026-06-05T08:00:00Z",
      "waktu_sampai": null
    }
  }
}
```
