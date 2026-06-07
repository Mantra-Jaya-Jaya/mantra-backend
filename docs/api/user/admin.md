# User — Admin API

Auth (Admin). Dashboard, profil, dan chart admin.

---

## GET /admin/dashboard

Mendapatkan ringkasan dashboard admin.

---

## GET /admin/dashboard/chart

Mendapatkan data chart untuk dashboard admin.

---

## GET /admin/profil

Mendapatkan profil admin yang login.

---

## PUT /admin/profil

Mengupdate profil admin.

**Request:**
```json
{
  "nama_lengkap": "Admin Baru",
  "email": "admin@mantra.web.id"
}
```
