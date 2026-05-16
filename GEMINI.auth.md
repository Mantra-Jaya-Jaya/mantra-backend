# GEMINI.auth.md — MANTRA Backend: Controller Auth (Data Asli)

> Kerjakan di branch yang sama setelah GEMINI.models.md selesai.
> Sumber kebenaran adalah **kode yang ada** — models, auth_helper.go, dan
> konvensi yang sudah dipakai di controller lain.
> Jangan asumsikan nama field, tabel, atau konvensi dari dokumentasi lama.

---

## LANGKAH 0 — AUDIT SEBELUM CODING (WAJIB)

Baca dan pahami semua file berikut sebelum menulis apapun:

### 0.1 Baca `controllers/auth/auth_helper.go`
Catat:
- Function apa saja yang sudah ada (generate token, validate token, dll)
- Signature masing-masing function (parameter dan return value)
- Nama dan tipe claims yang dipakai di JWT
- Cara ambil secret key (env variable apa yang dipakai)

### 0.2 Baca `controllers/auth/auth_controller.go`
Catat:
- Function apa saja yang sudah ada
- Mana yang sudah pakai data asli, mana yang masih dummy/placeholder
- Pola response yang sudah dipakai

### 0.3 Baca struct `User` dan `Role` di `models/`
Catat nama field dan tag `db:` persis seperti di kode — ini referensi untuk semua query.

### 0.4 Cari tabel refresh token di models
Temukan struct yang merepresentasikan refresh token. Catat:
- Nama tabel di DB (dari tag `db:` atau nama struct)
- Semua nama kolom yang ada (terutama: token, user id, expires, revoked)

### 0.5 Cek cara koneksi DB dipakai di controller lain
Ikuti pola yang sama — jangan buat cara baru.

### 0.6 Tulis ringkasan audit sebelum lanjut
```
AUDIT RESULT:
- JWT generate function: [nama function di auth_helper]
- JWT validate function: [nama function di auth_helper]
- JWT claims fields: [field apa saja]
- JWT secret env key: [nama env variable]
- Refresh secret env key: [nama env variable]
- Tabel refresh token: [nama tabel]
- Kolom token di refresh token: [nama kolom]
- Kolom revoked di refresh token: [nama kolom]
- Kolom expires di refresh token: [nama kolom]
- Function yang sudah ada di auth_controller: [list]
- Function yang masih dummy: [list]
- DB package dan cara inject: [cara yang dipakai]
```

---

## LANGKAH 1 — IDENTIFIKASI YANG PERLU DIUBAH

Berdasarkan audit, tentukan:
- Function mana yang perlu diimplementasi dari nol
- Function mana yang perlu diubah dari dummy ke data asli
- Function mana yang sudah benar dan tidak perlu disentuh

**Jangan ubah yang sudah benar.**

---

## LANGKAH 2 — IMPLEMENTASI

Kerjakan hanya function yang diidentifikasi perlu diubah di Langkah 1.
Semua query wajib pakai nama field dari hasil audit Langkah 0.

### `Login`
**Route:** `POST /api/v1/login` — public, semua role

**Alur:**
1. Bind dan validasi request body — field `username` dan `password` wajib ada
2. Query tabel `user` berdasarkan `username` — gunakan nama field dari audit
3. Jika tidak ditemukan → **401** kode `AUTH_001`
4. Bandingkan password dengan bcrypt (`bcrypt.CompareHashAndPassword`)
5. Jika salah → **401** kode `AUTH_001` (pesan sama — jangan bocorkan mana yang salah)
6. Generate access token dan refresh token menggunakan function dari `auth_helper.go`
7. Simpan refresh token ke tabel refresh token — gunakan nama kolom dari audit
8. Deteksi client:
   - Ada header `Authorization: Bearer` → client Flutter → token di response body
   - Tidak ada → client Next.js → token di `httpOnly Cookie`
9. Ambil `profile_id` dengan query ke tabel yang sesuai role (`customer`, `kasir`, `kurir`)
   berdasarkan `id_user` — cek nama FK di masing-masing tabel dari models

**Response Flutter:**
```json
{
  "status": "success",
  "message": "Login berhasil",
  "data": {
    "access_token": "...",
    "token_type": "Bearer",
    "expires_in": 900,
    "user": {
      "id_user": 1,
      "username": "...",
      "email": "...",
      "nama_lengkap": "...",
      "role": "customer",
      "profile_id": 10
    }
  }
}
```

**Response Next.js:**
```json
{
  "status": "success",
  "message": "Login berhasil",
  "data": {
    "user": {
      "id_user": 1,
      "username": "...",
      "nama_lengkap": "...",
      "role": "admin"
    }
  }
}
```

**Cookie Next.js:**
```
Set-Cookie: access_token=<token>; HttpOnly; Secure; SameSite=Strict; Path=/; Max-Age=900
Set-Cookie: refresh_token=<token>; HttpOnly; Secure; SameSite=Strict; Path=/api/v1/auth/refresh; Max-Age=604800
```

---

### `Register`
**Route:** `POST /api/v1/register` — public, hanya untuk customer

**Alur:**
1. Bind dan validasi — field wajib: `username`, `email`, `password`, `konfirmasi_password`, `nama_lengkap`, `no_telp`
2. Validasi `password` minimal 8 karakter → 422 `VAL_001`
3. Validasi `konfirmasi_password` cocok → 422 `VAL_002`
4. Validasi format email → 422 `VAL_003`
5. Cek duplikasi `username` → 409 `CONF_001`
6. Cek duplikasi `email` → 409 `CONF_002`
7. Hash password bcrypt cost 12
8. Query `id_role` untuk customer dari tabel `role` — gunakan nama kolom dari audit
9. Insert ke tabel `user`, ambil id yang dibuat
10. Insert ke tabel `customer`
11. Response **201** — jangan auto-login, arahkan ke halaman login

---

### `RefreshToken`
**Route:** `POST /api/v1/auth/refresh`

**Alur:**
1. Ambil refresh token:
   - Flutter: dari request body
   - Next.js: dari cookie `refresh_token`
2. Query tabel refresh token — gunakan nama kolom dari audit:
   ```sql
   SELECT * FROM <tabel_refresh_token>
   WHERE <kolom_token> = $1
     AND <kolom_expires> > NOW()
     AND <kolom_revoked> IS NULL
   ```
3. Jika tidak ditemukan atau expired → **401** `AUTH_003`
4. Generate access token baru menggunakan function dari `auth_helper.go`
5. Flutter: return token baru di body
6. Next.js: set cookie `access_token` baru, return 200

---

### `Logout`
**Route:** `POST /api/v1/logout` — butuh `AuthMiddleware`

**Alur:**
1. Ambil `user_id` dari JWT context — gunakan key yang ditemukan di audit
2. Ambil refresh token (Flutter: body, Next.js: cookie)
3. Revoke di DB — gunakan nama kolom dari audit:
   ```sql
   UPDATE <tabel_refresh_token>
   SET <kolom_revoked> = NOW()
   WHERE <kolom_token> = $1 AND <kolom_user_id> = $2
   ```
4. Next.js: clear cookie (`Max-Age=0`)
5. Response **200**

---

### `ChangePassword`
**Route:** `PUT /api/v1/change-password` — butuh `AuthMiddleware`

**Alur:**
1. Ambil `user_id` dari context
2. Bind body: `password_lama`, `password_baru`, `konfirmasi_password`
3. Validasi `password_baru` == `konfirmasi_password` → 422 `VAL_002`
4. Query password hash dari DB berdasarkan `user_id`
5. Bandingkan `password_lama` dengan hash → jika salah, **400** `REQ_003`
6. Hash `password_baru` bcrypt cost 12
7. Update kolom password di tabel `user`
8. Revoke **semua** refresh token user (logout semua device):
   ```sql
   UPDATE <tabel_refresh_token>
   SET <kolom_revoked> = NOW()
   WHERE <kolom_user_id> = $1 AND <kolom_revoked> IS NULL
   ```
9. Response **200**

---

## KODE ERROR YANG DIPAKAI

| Kode        | HTTP | Kondisi                                    |
|-------------|------|--------------------------------------------|
| `AUTH_001`  | 401  | Username/password salah, atau token invalid|
| `AUTH_003`  | 401  | Refresh token expired atau tidak valid     |
| `CONF_001`  | 409  | Username sudah terdaftar                   |
| `CONF_002`  | 409  | Email sudah terdaftar                      |
| `VAL_001`   | 422  | Validasi umum gagal                        |
| `VAL_002`   | 422  | Konfirmasi password tidak cocok            |
| `VAL_003`   | 422  | Format email tidak valid                   |
| `REQ_003`   | 400  | Password lama salah                        |
| `SERVER_001`| 500  | Error tak terduga — log, jangan expose     |

---

## CHECKLIST

- [ ] Audit selesai dan ringkasan ditulis sebelum coding dimulai
- [ ] Hanya function yang perlu diubah yang disentuh
- [ ] Semua query pakai nama field dari struct models (hasil audit)
- [ ] Semua token generate/validate pakai function dari `auth_helper.go`
- [ ] `Login` deteksi client Flutter vs Next.js dengan benar
- [ ] `Login` simpan refresh token ke DB
- [ ] `Register` hash bcrypt cost 12 dan insert ke dua tabel
- [ ] `RefreshToken` query dan validasi dari tabel refresh token
- [ ] `Logout` revoke refresh token di DB
- [ ] `ChangePassword` revoke semua refresh token user
- [ ] Password tidak pernah muncul di response manapun
- [ ] Error teknis tidak di-expose ke client — log saja

---

## YANG TIDAK BOLEH DILAKUKAN

- Jangan asumsikan nama field dari dokumentasi — baca dari kode yang ada
- Jangan hardcode JWT secret — selalu dari `os.Getenv(...)`
- Jangan ubah `auth_helper.go` kecuali ada bug yang perlu diperbaiki
- Jangan ubah route atau middleware
- Jangan simpan plain text password ke database
- Jangan return pesan berbeda untuk username salah vs password salah — bocorkan info
- Jangan expose stack trace atau detail query error ke client
