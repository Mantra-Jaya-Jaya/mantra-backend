# User — Alamat API

Auth (Customer). CRUD alamat pengiriman.

---

## GET /api/v1/customer/alamat

Mendapatkan daftar alamat customer yang login.

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id_alamat": 1,
      "public_id": "uuid-...",
      "nama_penerima": "John Doe",
      "label_alamat": "Rumah",
      "no_telp_penerima": "08123456789",
      "alamat_lengkap": "Jl. Contoh No. 1, Jakarta",
      "is_utama": true
    }
  ]
}
```

---

## POST /api/v1/customer/alamat

Tambah alamat baru.

**Request:**
```json
{
  "nama_penerima": "John Doe",
  "label_alamat": "Kantor",
  "no_telp_penerima": "08123456789",
  "alamat_lengkap": "Jl. Kantor No. 1, Jakarta",
  "latitude": -6.2,
  "longitude": 106.8,
  "catatan_lokasi": "Gedung Biru Lantai 3",
  "is_utama": false
}
```

---

## PUT /api/v1/customer/alamat/:public_id

Update data alamat.

---

## DELETE /api/v1/customer/alamat/:public_id

Hapus alamat.
