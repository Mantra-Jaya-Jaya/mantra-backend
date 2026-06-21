# Notifikasi API

---

## GET /api/v1/customer/notifikasi (Customer) / GET /api/v1/kasir/notifikasi (Kasir) / GET /api/v1/kurir/notifikasi (Kurir)

Auth. Mendapatkan daftar notifikasi milik user (customer, kasir, atau kurir) yang login.

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

## GET /api/v1/admin/notifikasi

Auth (Admin). Mendapatkan notifikasi untuk admin.

> Response structure sama dengan di atas.

---

## Notes

- Notifikasi difilter berdasarkan `id_user` dari JWT token — user hanya melihat notifikasinya sendiri.
- Status: `unread` | `read`
