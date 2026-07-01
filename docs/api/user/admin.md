# User — Admin API

Auth (Admin). Dashboard, profil, dan chart admin.

---

## GET /api/v1/admin/dashboard

Mendapatkan ringkasan dashboard admin.

---

## GET /api/v1/admin/dashboard/chart

Mendapatkan data chart untuk dashboard admin.

---

## GET /api/v1/admin/profil

Mendapatkan profil admin yang login.

---

## PUT /api/v1/admin/profil

Mengupdate profil admin.

**Request:**

```json
{
  "nama_lengkap": "Admin Baru",
  "email": "admin@mantra.web.id"
}
```
