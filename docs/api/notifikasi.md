# Notifikasi API

---

## GET /customer/notifikasi (Customer) / GET /kasir/notifikasi (Kasir)

Auth. Mendapatkan daftar notifikasi milik user yang login.

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id_notifikasi": 1,
      "judul": "Pesanan Dikemas",
      "pesan": "Pesanan #123 sedang dikemas",
      "status": "unread",
      "created_at": "2026-06-05T10:00:00Z"
    }
  ]
}
```

---

## GET /admin/notifikasi

Auth (Admin). Mendapatkan notifikasi untuk admin.

> Response structure sama dengan di atas.

---

## Notes

- Notifikasi difilter berdasarkan `id_user` dari JWT token — user hanya melihat notifikasinya sendiri.
- Status: `unread` | `read`
