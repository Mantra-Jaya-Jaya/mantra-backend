# User — Kasir API

Auth (Kasir). Profil kasir.

---

## GET /api/v1/kasir/profil

Mendapatkan profil kasir yang login.

**Response:**
```json
{
  "status": "success",
  "data": {
    "id_user": 2,
    "username": "kasir1",
    "nama_lengkap": "Kasir Satu",
    "shift": "Pagi",
    "no_telp": "08123456788"
  }
}
```
