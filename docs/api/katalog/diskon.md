# Katalog — Diskon API

---

## GET /api/v1/customer/promo (Customer) / GET /api/v1/admin/diskon (Admin)

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

## GET /api/v1/admin/diskon/semua

Auth (Admin). Mendapatkan semua data diskon (tidak difilter tanggal).

---

## POST /api/v1/admin/diskon

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

## DELETE /api/v1/admin/diskon/:public_id

Auth (Admin). Menghapus diskon.

---

## POST /api/v1/admin/diskon/upload

Auth (Admin). Upload banner diskon ke MinIO.

**Request:** `multipart/form-data` — field `banner`
