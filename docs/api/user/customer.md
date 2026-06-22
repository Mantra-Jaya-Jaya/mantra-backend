# User — Customer API

Auth (Customer). Profil dan manajemen akun customer.

---

## GET /api/v1/customer/profil

Mendapatkan profil customer yang login.

**Response:**

```json
{
  "status": "success",
  "data": {
    "id_user": 1,
    "public_id": "uuid-...",
    "username": "johndoe",
    "email": "john@example.com",
    "nama_lengkap": "John Doe",
    "foto_profil": null,
    "no_telp": "08123456789"
  }
}
```

---

## PUT /api/v1/customer/akun

Mengupdate data akun customer.

**Request:**

```json
{
  "nama_lengkap": "John Updated",
  "email": "john.baru@example.com",
  "no_telp": "08123456780"
}
```
