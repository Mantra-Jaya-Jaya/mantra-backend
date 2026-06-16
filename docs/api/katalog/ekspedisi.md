# Katalog — Ekspedisi API

Auth (Admin). Semua endpoint CRUD ekspedisi dan layanan pengiriman.

---

## GET /api/v1/admin/ekspedisi

Daftar semua ekspedisi.

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "public_id": "uuid-...",
      "nama_ekspedisi": "Ninja Xpress",
      "logo": "https://...",
      "deskripsi": "Layanan ekspedisi nasional",
      "is_active": true,
      "layanan": [
        {
          "public_id": "uuid-...",
          "nama_layanan": "REG",
          "estimasi_min": 2,
          "estimasi_max": 4,
          "is_active": true
        }
      ]
    }
  ]
}
```

---

## POST /api/v1/admin/ekspedisi

Tambah ekspedisi baru.

**Request:**
```json
{
  "nama_ekspedisi": "Ninja Xpress",
  "logo": "https://...",
  "deskripsi": "Layanan ekspedisi nasional",
  "kode_api": "",
  "is_active": true
}
```

---

## PUT /api/v1/admin/ekspedisi/:public_id

Update data ekspedisi.

---

## DELETE /api/v1/admin/ekspedisi/:public_id

Hapus ekspedisi (cascade menghapus layanan terkait).

---

## POST /api/v1/admin/ekspedisi/layanan

Tambah layanan baru untuk ekspedisi.

**Request:**
```json
{
  "id_ekspedisi": 1,
  "nama_layanan": "REG",
  "deskripsi": "Layanan Reguler",
  "estimasi_min": 2,
  "estimasi_max": 4,
  "is_active": true
}
```

---

## PUT /api/v1/admin/ekspedisi/layanan/:public_id

Update layanan.

---

## DELETE /api/v1/admin/ekspedisi/layanan/:public_id

Hapus layanan.
