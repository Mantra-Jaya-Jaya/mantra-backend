# Pengantaran (Courier & Shipping) API

API untuk melacak pesanan (Customer) dan mengelola pengantaran (Kurir).

---

## GET /api/v1/customer/pesanan/:public_id/lacak

Melacak status pengiriman pesanan (Customer).

**Auth:** Wajib Login (Customer)

**Response (Ekspedisi Eksternal - Biteship):**

```json
{
  "status": "success",
  "message": "Data lacak pesanan berhasil diambil",
  "data": {
    "id_pesanan": "9e3c8162-...",
    "nomor_resi": "JT1234567890",
    "ekspedisi": "J&T",
    "tipe_ekspedisi": "eksternal",
    "history": [
      {
        "time": "2026-06-12T10:00:00Z",
        "description": "Paket telah diterima oleh J&T Gateway",
        "status": "picked_up",
        "city": "Semarang"
      }
    ]
  }
}
```

**Response (Ekspedisi Internal - Kurir Toko):**

```json
{
  "status": "success",
  "message": "Data lacak pesanan berhasil diambil",
  "data": {
    "id_pesanan": "9e3c8162-...",
    "tipe_ekspedisi": "internal",
    "kurir": {
      "nama": "Ricardo Holahilo",
      "plat_nomor": "",
      "foto": "http://minio:9000/profiles/avatar.png"
    },
    "lokasi_kurir": {
      "latitude": -7.02,
      "longitude": 110.43
    },
    "estimasi_tiba": "",
    "jarak_meter": 0
  }
}
```

---

## PATCH /api/v1/kurir/pengantaran/:public_id/lokasi

Memperbarui koordinat GPS lokasi kurir secara berkala selama pengantaran.

**Auth:** Wajib Login (Kurir)

**Request:**

```json
{
  "latitude": -7.02561,
  "longitude": 110.43128
}
```

**Response:**

```json
{
  "status": "success",
  "message": "Lokasi kurir berhasil diperbarui"
}
```

---

## POST /api/v1/kurir/pengantaran/:public_id/ambil

Mengambil (claim) tugas pengantaran pesanan yang statusnya "Dikemas".

**Auth:** Wajib Login (Kurir)

**Response:**

```json
{
  "status": "success",
  "message": "Pesanan berhasil diambil",
  "data": {
    "id_pengantaran": "8f3c8162-...",
    "status": "Dalam Perjalanan"
  }
}
```

---

## POST /api/v1/kurir/pengantaran/:public_id/status

Memperbarui status pengantaran (misal: "Tiba di Tujuan", "Selesai").
Jika status diubah ke **Selesai**, kurir wajib menyertakan foto bukti pengiriman sebagai `multipart/form-data` dengan field name `"foto"`.

**Auth:** Wajib Login (Kurir)

**Request (Multipart/Form-Data untuk Selesai):**

- `status`: "Selesai" (Text)
- `foto`: [File Image] (Binary)

**Response:**

```json
{
  "status": "success",
  "message": "Status pengantaran berhasil diperbarui",
  "data": {
    "foto_bukti": "http://minio:9000/pengantaran/bukti_123.jpg",
    "status": "Selesai"
  }
}
```
