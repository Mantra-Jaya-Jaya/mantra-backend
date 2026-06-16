# Katalog — Kategori API

---

## GET /api/v1/customer/kategori (Customer) / GET /api/v1/kasir/kategori (Kasir) / GET /api/v1/admin/kategori (Admin)

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

## POST /api/v1/admin/kategori

Auth (Admin). Menambah kategori baru.

**Request:**
```json
{
  "nama_kategori": "Minuman"
}
```

---

## PUT /api/v1/admin/kategori/:public_id

Auth (Admin). Mengupdate nama kategori.

---

## DELETE /api/v1/admin/kategori/:public_id

Auth (Admin). Menghapus kategori.

---

## POST /api/v1/admin/kategori/upload

Auth (Admin). Upload icon kategori ke MinIO.

**Request:** `multipart/form-data` — field `icon`
