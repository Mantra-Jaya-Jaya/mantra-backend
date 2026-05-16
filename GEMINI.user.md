# GEMINI.user.md — MANTRA Backend: Controller User (Data Asli)

> Kerjakan di branch yang sama dengan GEMINI.models.md (setelah selesai).
> Sumber kebenaran adalah **kode yang ada** — models, query yang sudah ada,
> dan konvensi yang sudah dipakai di controller lain.
> Jangan asumsikan nama field, tabel, atau konvensi dari dokumentasi lama.

---

## LANGKAH 0 — AUDIT SEBELUM CODING (WAJIB)

Sebelum menulis satu baris kode, baca dan pahami:

### 0.1 Baca semua struct di `models/`
Catat nama field dan tag `db:` yang dipakai untuk tabel-tabel berikut:
- `user`
- `customer`
- `kasir`
- `kurir`
- `alamat` (atau nama lain yang merepresentasikan alamat pengiriman)
- `role`

**Ini adalah referensi utama untuk semua query yang akan ditulis.**

### 0.2 Baca controller yang sudah ada dan sudah bekerja
Cari controller yang sudah punya query ke database dan berjalan. Pelajari:
- Cara koneksi DB dipakai (package apa, cara inject, nama variabel)
- Pola query yang dipakai (`sqlx`, `database/sql`, ORM, dll)
- Pola response yang dipakai (struct response terpisah atau langsung `gin.H`)
- Cara ambil `user_id` dari JWT context (key dan tipe datanya)

### 0.3 Cek struktur folder `controllers/` yang sudah ada
Ikuti struktur dan konvensi yang sudah ada — jangan buat pola baru.

### 0.4 Tulis ringkasan audit sebelum lanjut
```
AUDIT RESULT:
- DB package: [nama package]
- Query style: [sqlx / database/sql / ORM]
- JWT context key untuk user_id: [nama key]
- Tipe data user_id dari context: [int / int64 / string]
- Response pattern: [gin.H / struct]
- Nama tabel user di DB: [user / users]
- Nama tabel alamat di DB: [hasil baca dari migration/model]
```

---

## LANGKAH 1 — BUAT STRUKTUR FILE

Buat folder dan file berikut jika belum ada, ikuti konvensi nama yang sudah ada:

```
controllers/user/
├── customer_controller.go
├── kasir_controller.go
├── kurir_controller.go
└── alamat_controller.go
```

Package name: `user` (sesuaikan jika konvensi project berbeda).

---

## LANGKAH 2 — IMPLEMENTASI (berdasarkan hasil audit)

Semua query wajib menggunakan nama field dan tabel **persis seperti yang ada
di struct models** — hasil audit Langkah 0.

### `customer_controller.go`

**`GetProfilCustomer`**
- Ambil `user_id` dari JWT context (gunakan key dan tipe yang ditemukan di audit)
- Query JOIN tabel `user` dan `customer` berdasarkan `user_id`
- Field yang dikembalikan: semua field public dari kedua tabel (kecuali `password`)
- Jika `foto_profil` tidak kosong, build URL absolut menggunakan `os.Getenv("BASE_URL")`
- Response 200

**`EditAkunCustomer`**
- Ambil `user_id` dari context
- Bind request body — field yang boleh diubah: `nama_lengkap`, `email`, `no_telp`, `foto_profil`
- Semua field opsional — hanya update field yang dikirim
- Jika `email` diubah, cek duplikasi dulu → 409 kode `CONF_002`
- Update tabel `user` dan/atau `customer` sesuai field yang berubah
- Response 200

---

### `kasir_controller.go`

**`GetProfilKasir`**
- Ambil `user_id` dari context
- Query JOIN tabel `user` dan `kasir`
- Response 200 tanpa field `password`

**`GetSemuaKasir`** — untuk admin
- Query semua kasir dengan JOIN `user` dan `kasir`
- Support query param: `search` (nama/username), `page`, `limit`
- Response 200 dengan pagination `meta`

**`TambahKasir`** — untuk admin
- Bind dan validasi field wajib: `username`, `email`, `password`, `nama_lengkap`
- Cek duplikasi `username` → 409 `CONF_001`
- Cek duplikasi `email` → 409 `CONF_002`
- Hash password bcrypt cost 12
- Query `id_role` untuk kasir dari tabel `role` — cek nama kolom dan nilai dari struct `role` yang ada
- Insert ke tabel `user`, ambil `id_user` yang dibuat
- Insert ke tabel `kasir`
- Response 201

**`DetailKasir`** — untuk admin
- Ambil `id_kasir` dari URL param
- Query JOIN `user` dan `kasir`
- 404 jika tidak ditemukan

**`EditKasir`** — untuk admin
- Ambil `id_kasir` dari URL param
- Semua field opsional
- Jika `email` diubah, cek duplikasi

**`HapusKasir`** — untuk admin
- Hapus dari tabel `kasir` dulu (FK), lalu tabel `user`
- Revoke semua refresh token user:
  ```sql
  UPDATE refresh_token SET revoked_at = NOW()
  WHERE id_user = $1 AND revoked_at IS NULL
  ```
- Cek nama tabel `refresh_token` dan kolom-kolomnya dari model yang ada

---

### `kurir_controller.go`

Pola identik dengan `kasir_controller.go`. Setelah baca struct `kurir` dari models,
ganti semua referensi tabel dan field sesuai yang ada di kode — jangan asumsikan
nama field sama dengan kasir.

Function yang harus ada:
- `GetProfilKurir` — role: kurir
- `GetSemuaKurir` — role: admin
- `TambahKurir` — role: admin
- `DetailKurir` — role: admin
- `EditKurir` — role: admin
- `HapusKurir` — role: admin (revoke refresh token juga)

---

### `alamat_controller.go`

**Ownership check — wajib ada di setiap function yang akses/edit alamat:**

```go
// Pola ownership check — sesuaikan nama tabel dan field dari hasil audit
// Selalu return 403 jika bukan miliknya, BUKAN 404
var count int
err := db.QueryRow(`
    SELECT COUNT(*) FROM <tabel_alamat> a
    JOIN <tabel_customer> c ON c.<pk_customer> = a.<fk_customer>
    WHERE a.<pk_alamat> = $1 AND c.<fk_user> = $2
`, idAlamat, userID).Scan(&count)

if count == 0 {
    c.JSON(403, gin.H{
        "status": "error",
        "message": "Anda tidak memiliki akses ke resource ini",
        "error": gin.H{"code": "AUTH_002", "detail": "Bukan milik user ini"},
    })
    return
}
```

**`GetAlamat`** — customer
- Ambil `user_id` dari context
- Query `id_customer` dari tabel customer
- Query semua alamat milik customer

**`TambahAlamat`** — customer
- Field wajib: `nama_penerima`, `no_telp_penerima`, `alamat_lengkap`
- Jika `is_utama = true` → reset semua alamat lain ke `false` dulu
- Insert ke tabel alamat

**`EditAlamat`** — customer
- Ownership check dulu
- Semua field opsional
- Jika `is_utama = true` → reset dulu

**`HapusAlamat`** — customer
- Ownership check dulu
- Hapus dari tabel alamat

---

## KONVENSI RESPONSE

Ikuti format yang sudah dipakai di controller lain yang sudah ada. Jika belum ada
referensi, gunakan:

```go
// Success
c.JSON(200, gin.H{
    "status":  "success",
    "message": "...",
    "data":    data,
})

// Success dengan pagination
c.JSON(200, gin.H{
    "status":  "success",
    "message": "...",
    "data":    data,
    "meta": gin.H{
        "page": page, "limit": limit,
        "total": total, "total_pages": totalPages,
    },
})

// Error
c.JSON(statusCode, gin.H{
    "status":  "error",
    "message": "Pesan untuk UI",
    "error":   gin.H{"code": "KODE", "detail": "..."},
})
```

---

## KODE ERROR YANG DIPAKAI

| Kode        | HTTP | Kondisi                          |
|-------------|------|----------------------------------|
| `AUTH_002`  | 403  | Ownership violation              |
| `CONF_001`  | 409  | Username duplikat                |
| `CONF_002`  | 409  | Email duplikat                   |
| `DATA_004`  | 404  | Resource tidak ditemukan         |
| `VAL_001`   | 422  | Validasi field gagal             |
| `SERVER_001`| 500  | Error tak terduga (log, jangan expose detail) |

---

## CHECKLIST

- [ ] Audit selesai dan ringkasan ditulis sebelum coding
- [ ] Semua query pakai nama field dari struct models yang ada (bukan asumsi)
- [ ] `GetProfilCustomer` dan `EditAkunCustomer` berjalan dengan data asli
- [ ] CRUD kasir semua berjalan dengan data asli
- [ ] CRUD kurir semua berjalan dengan data asli
- [ ] `HapusKasir` dan `HapusKurir` revoke refresh token
- [ ] Semua endpoint alamat ada ownership check
- [ ] Ownership violation return **403**, bukan 404
- [ ] `foto_profil` dikonversi ke URL absolut di response, disimpan path relatif di DB
- [ ] Tidak ada data dummy atau hardcode
- [ ] Password tidak pernah ada di response manapun

---

## YANG TIDAK BOLEH DILAKUKAN

- Jangan asumsikan nama field dari dokumentasi — baca dari struct models yang ada
- Jangan ambil `user_id` dari URL param/body untuk ownership — harus dari JWT context
- Jangan return 404 untuk ownership violation — selalu 403
- Jangan duplikasi logic yang sudah ada di controller lain
- Jangan ubah route, middleware, atau models
- Jangan expose detail error database ke client — log saja, return `SERVER_001`
