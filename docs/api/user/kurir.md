# User — Kurir API

Auth (Kurir). Profil kurir.

---

## GET /api/v1/kurir/profile

Mendapatkan profil kurir yang login.

**Response:**

```json
{
  "status": "success",
  "message": "Data profil berhasil diambil",
  "data": {
    "id_kurir": 1,
    "public_id": "uuid-...",
    "no_telp": "08123456787",
    "tempat_lahir": "Jakarta",
    "tanggal_lahir": "1995-01-01",
    "jenis_kelamin": "L",
    "alamat": "Jl. Contoh No. 1",
    "pendidikan_terakhir": "SMA",
    "nik": "3201010101950001",
    "id_user": 3,
    "user_public_id": "uuid-...",
    "username": "ricardo1",
    "email": "ricardo@mantra.web.id",
    "nama_lengkap": "Ricardo Holahilo",
    "foto_profil": "http://localhost:8080/uploads/foto.jpg"
  }
}
```
