# Auth API

Endpoint autentikasi — public dan protected.

---

## POST /api/v1/login

Public. Login dengan username atau email.

**Request:**
```json
{
  "username": "johndoe",
  "password": "password123"
}
```

**Response (Flutter):**
```json
{
  "status": "success",
  "message": "Login berhasil",
  "data": {
    "access_token": "eyJhbG...",
    "refresh_token": "a1b2c3d4...",
    "token_type": "Bearer",
    "expires_in": 1800,
    "user": {
      "id_user": 1,
      "username": "johndoe",
      "email": "john@example.com",
      "nama_lengkap": "John Doe",
      "role": "Customer",
      "profile_id": 1
    }
  }
}
```

**Response (Next.js):** Set cookie `access_token` + `refresh_token`, body tanpa token.

---

## POST /api/v1/register

Public. Registrasi customer baru.

**Request:**
```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "password123",
  "konfirmasi_password": "password123",
  "nama_lengkap": "John Doe",
  "no_telp": "08123456789"
}
```

**Response:** `201 Created`
```json
{
  "status": "success",
  "message": "Registrasi berhasil",
  "data": {
    "id_user": 2,
    "username": "johndoe",
    "email": "john@example.com",
    "nama_lengkap": "John Doe",
    "no_telp": "08123456789",
    "role": "customer"
  }
}
```

**Validasi:**
- Password minimal 8 karakter
- Konfirmasi password harus cocok
- Email format valid
- Username & email unique

---

## POST /api/v1/auth/refresh

Public. Mendapatkan access token baru menggunakan refresh token.

**Request (Flutter):**
```json
{
  "refresh_token": "a1b2c3d4..."
}
```

**Request (Next.js):** Kirim cookie `refresh_token` otomatis.

**Response (Flutter):**
```json
{
  "status": "success",
  "message": "Token berhasil diperbarui",
  "data": {
    "access_token": "eyJhbG...",
    "expires_in": 1800
  }
}
```

**Response (Next.js):** Set cookie baru, body tanpa token.

---

## POST /api/v1/logout

Auth required. Menonaktifkan refresh token.

**Request (Flutter):**
```json
{
  "refresh_token": "a1b2c3d4..."
}
```

**Request (Next.js):** Cukup panggil endpoint, cookie akan dihapus otomatis.

**Response:**
```json
{
  "status": "success",
  "message": "Logout berhasil"
}
```

---

## PUT /api/v1/change-password

Auth required. Mengubah password (semua session di-revoke).

**Request:**
```json
{
  "password_lama": "oldpass123",
  "password_baru": "newpass456",
  "konfirmasi_password": "newpass456"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Password berhasil diubah"
}
```

Semua refresh token milik user akan di-revoke (logout dari semua device).
