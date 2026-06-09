# Katalog — Ekspedisi API

Auth (Admin). Semua endpoint CRUD ekspedisi dan layanan pengiriman.

---

## GET /admin/ekspedisi

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
          "id_ekspedisi_layanan": 1,
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

## POST /admin/ekspedisi

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

## PUT /admin/ekspedisi/:public_id

Update data ekspedisi.

---

## DELETE /admin/ekspedisi/:public_id

Hapus ekspedisi (cascade menghapus layanan terkait).

---

## POST /admin/ekspedisi/layanan

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

## PUT /admin/ekspedisi/layanan/:id

Update layanan (id = id_ekspedisi_layanan numeric).

---

## DELETE /admin/ekspedisi/layanan/:id

Hapus layanan.
