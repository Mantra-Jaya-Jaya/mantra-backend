# Katalog — Diskon API

---

## GET /promo (Customer) / GET /admin/diskon (Admin)

Auth. Mendapatkan daftar promo/diskon yang aktif.

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id_diskon": 1,
      "nama_diskon": "Diskon Akhir Tahun",
      "besar_diskon": 20,
      "banner_diskon": "https://storage...",
      "tgl_mulai": "2026-01-01",
      "tgl_selesai": "2026-01-31"
    }
  ]
}
```

---

## GET /admin/diskon/semua

Auth (Admin). Mendapatkan semua data diskon (tidak difilter tanggal).

---

## POST /admin/diskon

Auth (Admin). Menambah diskon baru.

**Request:**
```json
{
  "nama_diskon": "Diskon Akhir Tahun",
  "besar_diskon": 20,
  "tgl_mulai": "2026-01-01",
  "tgl_selesai": "2026-01-31"
}
```

---

## DELETE /admin/diskon/:id_diskon

Auth (Admin). Menghapus diskon.

---

## POST /admin/diskon/upload

Auth (Admin). Upload banner diskon ke MinIO.

**Request:** `multipart/form-data` — field `banner`
