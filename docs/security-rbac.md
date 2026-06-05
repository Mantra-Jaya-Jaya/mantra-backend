# Security & RBAC

## Authentication Flow

### JWT Token

- **Access Token:** JWT HS256, masa berlaku 30 menit
- **Refresh Token:** Random 32-byte hex string, masa berlaku 7 hari, disimpan di DB
- **Password:** bcrypt cost factor 12

### Sliding Expiration

Jika sisa masa berlaku access token < 15 menit, token baru otomatis digenerate:

- **Flutter:** Token baru dikirim via header response `X-New-Access-Token`
- **Next.js:** Token baru di-set sebagai httpOnly cookie `access_token`

Client cukup membaca header/cookie baru dan menggunakannya untuk request selanjutnya.

### Multi-Client Detection

Deteksi otomatis via header `X-Client-Type`:

| Client | Header | Kirim Token Via |
|--------|--------|-----------------|
| Flutter | `X-Client-Type: flutter` | JSON response `access_token` |
| Next.js | `X-Client-Type: nextjs` (atau tidak ada) | httpOnly Cookie |

## Middleware Chain

```text
Request → RateLimit → CORS → AuthMiddleware → RoleMiddleware → Controller
```

| Middleware | File | Fungsi |
|-----------|------|--------|
| AuthMiddleware | `middleware/auth_middleware.go` | Validasi JWT, sliding expiration, set context |
| RoleMiddleware | `middleware/role_middleware.go` | Cek role user sesuai allowed roles |
| OwnershipMiddleware | `middleware/ownership_middleware.go` | Cek kepemilikan resource via public_id |

### AuthMiddleware Details

1. Extract token dari `Authorization: Bearer <token>` atau cookie `access_token`
2. Parse & validasi JWT dengan secret dari `JWT_SECRET` env (fallback: `rahasia_dapur_mantra`)
3. Jika valid dan sisa waktu < 15 menit → generate token baru
4. Set `user_id`, `public_id`, `role` ke Gin context

### OwnershipMiddleware Details

- Membandingkan `public_id` dari URL param dengan `public_id` dari JWT claims
- **Admin selalu bypass** pengecekan ini
- Error code: `AUTH_004`

## RBAC Matrix

| Endpoint | Customer | Kasir | Admin |
|----------|:--------:|:-----:|:-----:|
| `POST /login` | ✔ | ✔ | ✔ |
| `POST /register` | ✔ | ✘ | ✘ |
| `POST /auth/refresh` | ✔ | ✔ | ✔ |
| `POST /logout` | ✔ | ✔ | ✔ |
| `PUT /change-password` | ✔ | ✔ | ✔ |
| `GET /scan/:kode` | ✔ | ✔ | ✔ |
| `GET /customer/*` | ✔ | ✘ | ✘ |
| `GET /kasir/*` | ✘ | ✔ | ✘ |
| `GET /admin/*` | ✘ | ✘ | ✔ |

## Error Codes

| Kode | HTTP | Middleware | Arti |
|------|------|------------|------|
| `AUTH_001` | 401 | Auth | Token invalid / expired |
| `AUTH_002` | 403 | Role | Role tidak ditemukan di token |
| `AUTH_003` | 403 | Role | Tidak punya izin akses |
| `AUTH_004` | 403 | Ownership | Bukan pemilik resource |
| `VAL_001` | 400/422 | Controller | Validasi input gagal |
| `CONF_001` | 409 | Controller | Username sudah terdaftar |
| `CONF_002` | 409 | Controller | Email sudah terdaftar |
| `SERVER_001` | 500 | Any | Internal server error |

## Environment Variables

```env
JWT_SECRET=your_jwt_secret
JWT_REFRESH_SECRET=your_refresh_secret
```
