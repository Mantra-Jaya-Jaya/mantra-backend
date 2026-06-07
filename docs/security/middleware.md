# Middleware

## Chain Execution

```
Request → AuthMiddleware → RoleMiddleware → OwnershipMiddleware → Controller
```

## AuthMiddleware

**File:** `middleware/auth_middleware.go`

Berlaku untuk semua endpoint yang membutuhkan autentikasi.

### Proses
1. Extract token dari:
   - `Authorization: Bearer <token>` (prioritas utama — Flutter)
   - Cookie `access_token` (fallback — Next.js)
2. Parse JWT dengan secret dari `JWT_SECRET` env
3. Validasi signature & expiration
4. Jika valid:
   - Cek sliding expiration (sisa < 15 menit → generate token baru)
   - Set context: `user_id`, `public_id`, `role`
5. Jika tidak valid → `401 Unauthorized` dengan kode `AUTH_001`

### Error Response
```json
{
  "status": "error",
  "message": "Token tidak valid atau sudah expired",
  "error": {
    "code": "AUTH_001",
    "detail": "..."
  }
}
```

## RoleMiddleware

**File:** `middleware/role_middleware.go`

Cek apakah role user termasuk dalam daftar yang diizinkan.

### Penggunaan
```go
middleware.RoleMiddleware("admin")
middleware.RoleMiddleware("customer", "kasir")
```

### Error Codes
| Kode | Kondisi | HTTP |
|------|---------|------|
| `AUTH_002` | Role tidak ditemukan di token | 403 |
| `AUTH_003` | Role tidak punya izin akses | 403 |

## OwnershipMiddleware

**File:** `middleware/ownership_middleware.go`

Cek kepemilikan resource berdasarkan `public_id` di URL parameter.

### Alur
1. Ambil `public_id` dari JWT claims (context)
2. Ambil `public_id` dari URL parameter (contoh: `:public_id`)
3. Jika role **admin** → bypass (selalu diizinkan)
4. Jika cocok → lanjut
5. Jika tidak cocok → `403 Forbidden` dengan kode `AUTH_004`

### Penggunaan
```go
middleware.OwnershipMiddleware("public_id")
```
