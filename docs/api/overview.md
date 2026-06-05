# API Overview

Base URL dan format komunikasi untuk seluruh endpoint MANTRA.

## Base URL

| Environment | URL |
|-------------|-----|
| Local | `http://localhost:8080/api/v1` |
| Production | *(akan diisi)* |

## Authentication

Dua metode autentikasi tergantung client:

### Flutter (Customer & Kasir Apps)
```
Authorization: Bearer <access_token>
```

### Next.js (Admin Web)
Token dikirim via `httpOnly Cookie`:
```
Cookie: access_token=<token>; refresh_token=<token>
```

Deteksi client dilakukan otomatis via header `X-Client-Type`:
- `X-Client-Type: flutter` → response JSON dengan `access_token`
- `X-Client-Type: nextjs` → response JSON + set cookie

## Error Response Format

Semua error response memiliki struktur seragam:

```json
{
  "status": "error",
  "message": "Deskripsi error untuk user",
  "error": {
    "code": "ERR_XXX",
    "detail": "Penjelasan teknis untuk debugging"
  }
}
```

### Kode Error

| Kode | Arti |
|------|------|
| `VAL_001` | Validasi gagal — format input tidak sesuai |
| `VAL_002` | Validasi gagal — konfirmasi tidak cocok |
| `VAL_003` | Validasi gagal — format field tidak valid |
| `AUTH_001` | Token tidak valid / sudah expired |
| `AUTH_002` | Role tidak ditemukan di token |
| `AUTH_003` | Tidak punya izin akses resource |
| `AUTH_004` | Bukan pemilik resource |
| `CONF_001` | Username sudah terdaftar |
| `CONF_002` | Email sudah terdaftar |
| `REQ_003` | Password lama tidak sesuai |
| `REQ_004` | Data tidak ditemukan |
| `SERVER_001` | Internal server error |

## Success Response Format

```json
{
  "status": "success",
  "message": "Pesan sukses",
  "data": { }
}
```

## HTTP Status Codes

| Status | Deskripsi |
|--------|-----------|
| 200 | OK |
| 201 | Created |
| 400 | Bad Request |
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Not Found |
| 409 | Conflict |
| 422 | Unprocessable Entity |
| 500 | Internal Server Error |

## ID Format

- **Internal ID:** Integer auto-increment (private, jangan diekspos ke client)
- **External ID:** UUID `public_id` (untuk endpoint publik, mencegah ID enumeration)
