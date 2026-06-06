# Katalog — Kategori API

---

## GET /kategori (Customer & Kasir) / GET /admin/kategori (Admin)

Auth. Mendapatkan daftar semua kategori.

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id_kategori": 1,
      "nama_kategori": "Makanan",
      "icon_kategori": "https://storage..."
    }
  ]
}
```

---

## POST /admin/kategori

Auth (Admin). Menambah kategori baru.

**Request:**
```json
{
  "nama_kategori": "Minuman"
}
```

---

## PUT /admin/kategori/:id_kategori

Auth (Admin). Mengupdate nama kategori.

---

## DELETE /admin/kategori/:id_kategori

Auth (Admin). Menghapus kategori.

---

## POST /admin/kategori/upload

Auth (Admin). Upload icon kategori ke MinIO.

**Request:** `multipart/form-data` — field `icon`
