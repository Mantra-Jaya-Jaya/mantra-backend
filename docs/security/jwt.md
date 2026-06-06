# JWT Authentication

## Token Types

| Token | Masa Berlaku | Penyimpanan |
|-------|-------------|-------------|
| Access Token | 30 menit | Memory (Flutter) / httpOnly Cookie (Next.js) |
| Refresh Token | 7 hari | Database (`refresh_token` table) + httpOnly Cookie |

## Access Token (JWT)

Algorithm: **HS256**

### Claims Structure

```json
{
  "user_id": 1,
  "public_id": "uuid-...",
  "role": "Customer",
  "exp": 1712345678,
  "iat": 1712343878
}
```

| Claim | Type | Deskripsi |
|-------|------|-----------|
| `user_id` | uint | ID internal user |
| `public_id` | string | UUID untuk endpoint eksternal |
| `role` | string | Nama role (Customer, Admin, Kasir, Kurir) |
| `exp` | int64 | Expiration time (UNIX timestamp) |
| `iat` | int64 | Issued at time (UNIX timestamp) |

### Secret Key

- Diambil dari env `JWT_SECRET`
- Fallback: `rahasia_dapur_mantra` (hanya untuk development)

## Refresh Token

- Random 32-byte hex string (64 karakter hex)
- Disimpan di tabel `refresh_token`
- Memiliki kolom `revoked_at` — diisi saat logout
- Token expired otomatis tidak bisa dipakai (validasi `expires_at`)
- Saat ganti password: semua refresh token user di-revoke (logout semua device)

## Sliding Expiration

Di `middleware/auth_middleware.go`:

1. Parse token, dapatkan `exp` claim
2. Hitung sisa waktu: `time.Until(claims.ExpiresAt.Time)`
3. Jika sisa waktu > 0 **dan** < 15 menit:
   - Generate access token baru
   - **Flutter:** Kirim via header `X-New-Access-Token`
   - **Next.js:** Set cookie baru `access_token`
4. Client wajib membaca header/cookie baru untuk request selanjutnya

## Multi-Client Authentication

| Client | Kirim Token | Deteksi |
|--------|-------------|---------|
| Flutter | `Authorization: Bearer <token>` | Header `X-Client-Type: flutter` |
| Next.js | httpOnly Cookie (`access_token`) | Header tidak ada atau `X-Client-Type: nextjs` |

### Flow Login

1. User kirim username + password ke `POST /login`
2. Server validasi credential (bcrypt compare)
3. Generate access token (JWT, 30 menit)
4. Generate refresh token (random hex, 7 hari) → simpan di DB
5. Response sesuai client type:
   - Flutter: JSON dengan `access_token` + `refresh_token`
   - Next.js: JSON (tanpa token di body) + set httpOnly cookies
