# Error Codes & HTTP Status

## Error Response Format

Semua error response mengikuti format seragam:

```json
{
  "status": "error",
  "message": "Pesan untuk user",
  "error": {
    "code": "ERR_XXX",
    "detail": "Penjelasan teknis untuk debugging"
  }
}
```

## Error Code Table

| Kode | HTTP | Sumber | Arti |
|------|------|--------|------|
| `AUTH_001` | 401 | AuthMiddleware | Token tidak valid / expired |
| `AUTH_002` | 403 | RoleMiddleware | Role tidak ditemukan di token |
| `AUTH_003` | 403 | RoleMiddleware | Tidak punya izin akses resource |
| `AUTH_004` | 403 | OwnershipMiddleware | Bukan pemilik resource |
| `VAL_001` | 400 / 422 | Controller | Input tidak memenuhi aturan validasi |
| `VAL_002` | 422 | Controller | Konfirmasi password tidak cocok |
| `VAL_003` | 422 | Controller | Format email tidak valid |
| `CONF_001` | 409 | Controller | Username sudah terdaftar |
| `CONF_002` | 409 | Controller | Email sudah terdaftar |
| `REQ_003` | 400 | Controller | Password lama tidak sesuai |
| `REQ_004` | 404 | Controller | Data tidak ditemukan |
| `SERVER_001` | 500 | Any | Internal server error |

## HTTP Status Codes

| Status | Usage |
|--------|-------|
| 200 | OK — Request berhasil |
| 201 | Created — Resource berhasil dibuat |
| 400 | Bad Request — Input tidak valid |
| 401 | Unauthorized — Token invalid / expired |
| 403 | Forbidden — Tidak punya akses |
| 404 | Not Found — Resource tidak ditemukan |
| 409 | Conflict — Duplikasi data |
| 422 | Unprocessable Entity — Validasi gagal |
| 500 | Internal Server Error — Error sistem |
