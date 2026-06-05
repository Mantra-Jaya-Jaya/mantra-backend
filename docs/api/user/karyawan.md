# User — Karyawan API

Auth (Admin). Manajemen data karyawan. Karyawan adalah induk dari kasir dan kurir.

---

## GET /admin/karyawan

Mendapatkan daftar semua karyawan.

---

## POST /admin/karyawan

Menambah karyawan baru (sekaligus membuat akun user + profil kasir/kurir).

**Request:**
```json
{
  "nama_lengkap": "Karyawan Baru",
  "username": "karyawan1",
  "email": "karyawan@mantra.web.id",
  "password": "password123",
  "no_telp": "08123456787",
  "tempat_lahir": "Jakarta",
  "tanggal_lahir": "1995-01-01",
  "jenis_kelamin": "L",
  "alamat": "Jl. Contoh No. 1",
  "pendidikan_terakhir": "SMA",
  "nik": "3201010101950001",
  "role": "Kasir",
  "shift": "Pagi"
}
```

---

## GET /admin/karyawan/:id

Mendapatkan detail karyawan.

---

## PUT /admin/karyawan/:id

Mengupdate data karyawan.

---

## DELETE /admin/karyawan/:id

Menghapus data karyawan.

---

## POST /admin/karyawan/upload

Upload foto profil karyawan ke MinIO.

**Request:** `multipart/form-data` — field `foto`
