# Migrasi ID Integer ke Public ID (UUID)

> **Visi:** Semua endpoint REST API menggunakan `public_id` (UUID) sebagai identifier eksternal — tidak ada integer ID di URL. Konsisten, aman, dan siap untuk multi-tenant.

## Latar Belakang

Backend MANTRA menggunakan **ID obfuscation** — endpoint eksternal pakai `public_id` (UUID) bukan integer auto-increment. Tujuan:

- Tidak ada tebakan ID resource lain (enumeration attack)
- Konsisten antar tabel
- Siap untuk future multi-tenant

Saat ini **6 dari 8** path param sudah menggunakan `public_id`. Dua sisanya masih integer.

## Status Saat Ini

### Sudah Migrasi ✅ (controller sudah lookup by `public_id`)

| Path Param | Route | Model PublicId | Controller Lookup |
|---|---|---|---|
| `:public_id` (Barang) | `GET/PUT/DELETE /barang/:public_id` | ✅ | `public_id = ?` |
| `:id_keranjang` | `PATCH/DELETE /keranjang/:id_keranjang` | ✅ | `public_id = ?` |
| `:id_pesanan` | 5 route | ✅ | `public_id = ?` |
| `:id_alamat` | `PUT/DELETE /alamat/:id_alamat` | ✅ | `public_id = ?` |
| `:id_produk` | `GET /laporan/produk/:id_produk` | ✅ (Barang) | `public_id = ?` |
| `:id` (Karyawan) | 3 route | ✅ | `public_id = ?` |

### Belum Migrasi ❌ (masih integer)

| Path Param | Route | Model PublicId | Controller Lookup |
|---|---|---|---|
| `:id_kategori` | `PUT/DELETE /kategori/:id_kategori` | ❌ | `id_kategori = ?` (int) |
| `:id_diskon` | `DELETE /diskon/:id_diskon` | ❌ | `id_diskon = ?` (int) |

### Bug Ditemukan 🐛

| Bug | Lokasi | Masalah |
|---|---|---|
| `:id_order` vs `:id_pesanan` | `routes.go:67` — `pesanan_controller.go:126` | Nama param route `:id_order` tapi controller baca `c.Param("id_pesanan")` → **empty string** |
| GetRingkasanCheckout | `pembayaran_controller.go:24` | Query pakai `"id_pesanan = ?"` (integer) bukan `"public_id = ?"` |

## Dampak per Repo

### 🔴 `mantra-backend/` — 8-10 file

| # | File | Perubahan |
|---|---|---|
| 1 | `models/kategori.go` | Tambah field `PublicId uuid.UUID` |
| 2 | `models/diskon.go` | Tambah field `PublicId uuid.UUID` |
| 3 | `controllers/katalog/kategori_controller.go` | `UpdateKategori` & `HapusKategori` — hapus `strconv.Atoi`, lookup by `public_id` |
| 4 | `controllers/katalog/diskon_controller.go` | `HapusDiskon` — hapus `strconv.Atoi`, lookup by `public_id` |
| 5 | `controllers/transaksi/pembayaran_controller.go` | `GetRingkasanCheckout` — ganti `"id_pesanan = ?"` → `"public_id = ?"` |
| 6 | `routes/routes.go` | Ganti `:id_kategori` → `:public_id`, `:id_diskon` → `:public_id`, fix `:id_order` → `:id_pesanan` |
| 7 | `db/migrations/` | 2 migration baru: `ALTER TABLE kategori ADD COLUMN public_id UUID`, `ALTER TABLE diskon ADD COLUMN public_id UUID` |
| 8 | `docs/mantra.dbml` | Tambah `public_id uuid` di tabel kategori & diskon |
| 9 | `docs/api-collections/collection.bru` | Hapus `idKategori`, `idDiskon` dari vars (tidak perlu lagi) |
| 10 | `docs/task/` | File ini |

### 🟡 `mantra-admin-web/` — 2-3 halaman

| Halaman | File | Masalah |
|---|---|---|
| Kategori | `app/barang/kategori/page.tsx:129` | PUT `kategori/${editingId}` — `editingId` dari `item.id_kategori` → ganti ke `item.public_id` |
| Kategori | `app/barang/kategori/page.tsx:171` | DELETE `kategori/${kategoriToDelete}` — dari `item.id_kategori` → ganti |
| Diskon | `app/barang/diskon/page.tsx:149` | DELETE `diskon/${diskonToDelete}` — dari `item.id_diskon` → ganti |
| ⚠️ Karyawan | `app/karyawan/page.tsx:186` | DELETE `karyawan/${user.id}` — **perlu verifikasi: apakah `user.id` isi integer atau public_id?** |

### 🟡 `mantra-mobile/` — 1 file

| File | Baris | Potensi masalah |
|---|---|---|
| `lib/core/services/profile_service.dart` | 63, 77 | `$idAlamat` — parameter String, kemungkinan **sudah** public_id |

### ✅ Tidak terdampak

- `mantra-backend/controllers/keranjang/` — sudah pakai `public_id`
- `mantra-backend/controllers/transaksi/pesanan_controller.go` — sudah pakai `public_id`
- `mantra-backend/controllers/user/alamat_controller.go` — sudah pakai `public_id`
- `mantra-backend/controllers/user/karyawan_controller.go` — sudah pakai `public_id`
- `mantra-admin-web/app/barang/*` — semua pakai `public_id`
- `mantra-mobile/` — semua endpoint yang aktif sudah pakai String/UUID

## Rencana Implementasi

### Step 1: Backend — Model + Migration

Tambahkan field `PublicId` di 2 model dan buat migration.

```go
// models/kategori.go — tambah
PublicId uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();column:public_id;uniqueIndex" json:"public_id"`
```

```go
// models/diskon.go — tambah
PublicId uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();column:public_id;uniqueIndex" json:"public_id"`
```

Migration (Atlas CLI atau raw SQL):

```sql
ALTER TABLE kategori ADD COLUMN public_id UUID NOT NULL DEFAULT gen_random_uuid();
CREATE UNIQUE INDEX idx_kategori_public_id ON kategori(public_id);

ALTER TABLE diskon ADD COLUMN public_id UUID NOT NULL DEFAULT gen_random_uuid();
CREATE UNIQUE INDEX idx_diskon_public_id ON diskon(public_id);
```

### Step 2: Backend — Controller

#### `kategori_controller.go`

`UpdateKategori` (line 86-97):
- Hapus: `idStr := c.Param("id_kategori")` + `strconv.Atoi`
- Ganti: lookup `.Where("public_id = ?", c.Param("id_kategori"))`
- Hapus pemeriksaan `id_kategori == 0`

`HapusKategori` (line 153-164):
- Sama: ganti lookup + dependency check ke public_id

Cek barang terkait kategori (line 174):
- Ganti subquery dari `id_kategori = ?` ke lookup by public_id

#### `diskon_controller.go`

`HapusDiskon` (line 186-207):
- Hapus: `strconv.Atoi`
- Ganti: lookup `.Where("public_id = ?", c.Param("id_diskon"))`
- Ganti detach relasi barang: lookup `id_diskon` di Barang via public_id

#### `pembayaran_controller.go`

`GetRingkasanCheckout` (line 24):
- Ganti: `"id_pesanan = ?"` → `"public_id = ?"`

### Step 3: Backend — Routes

```go
// routes.go
adminGroup.PUT("/kategori/:public_id", ...)    // was :id_kategori
adminGroup.DELETE("/kategori/:public_id", ...) // was :id_kategori
adminGroup.DELETE("/diskon/:public_id", ...)   // was :id_diskon
kasirGroup.GET("/pesanan/:public_id", ...)     // was :id_order  (FIX BUG)
```

### Step 4: Backend — DBML + Bruno

- `docs/mantra.dbml`: tambah `public_id uuid` di tabel kategori & diskon
- `docs/api-collections/collection.bru`: hapus `idKategori`, `idDiskon` dari `vars:pre-request`
- `docs/api-collections/Admin/Kategori/Update.bru`: ganti `{{idKategori}}` → `{{publicId}}`
- `docs/api-collections/Admin/Kategori/Hapus.bru`: ganti `{{idKategori}}` → `{{publicId}}`
- `docs/api-collections/Admin/Diskon/Hapus.bru`: ganti `{{idDiskon}}` → `{{publicId}}`

### Step 5: Frontend — Admin Web

- `app/barang/kategori/page.tsx`: ganti `item.id_kategori` → `item.public_id`
- `app/barang/diskon/page.tsx`: ganti `item.id_diskon` → `item.public_id`
- Verifikasi `app/karyawan/page.tsx`: apakah `user.id` sudah public_id?

### Step 6: Frontend — Mobile

- Verifikasi `profile_service.dart`: `idAlamat` sebagai String sudah public_id

### Step 7: Testing

```bash
cd mantra-backend && go test ./...
```

Test manual via Bruno:
1. Login sebagai Admin
2. `Admin/Kategori/Daftar.bru` — dapatkan public_id kategori
3. `Admin/Kategori/Update.bru` — update dengan public_id
4. `Admin/Kategori/Hapus.bru` — hapus dengan public_id
5. `Admin/Diskon/Semua.bru` — dapatkan public_id diskon
6. `Admin/Diskon/Hapus.bru` — hapus dengan public_id
7. `Kasir/Transaksi/Cek Checkout.bru` — pastikan GetRingkasanCheckout OK

## Kriteria Selesai

- [ ] Tidak ada dependency `strconv.Atoi` untuk path param di controllers
- [ ] Semua path param di `routes.go` hanya `:public_id` atau `:kode_barcode`
- [ ] Kategori & Diskon model punya `PublicId` field
- [ ] Migration siap dijalankan (Atlas)
- [ ] `docs/mantra.dbml` sinkron
- [ ] `docs/api-collections/collection.bru` tidak ada `idKategori` / `idDiskon`
- [ ] Bruno collection Admin/Kategori & Admin/Diskon pake `{{publicId}}`
- [ ] `GET /api/v1/admin/kategori` return `public_id` di response JSON
- [ ] `GET /api/v1/admin/diskon/semua` return `public_id` di response JSON
- [ ] `go test ./...` lolos
- [ ] Bug `:id_order` fixed
- [ ] Bug `GetRingkasanCheckout` fixed
- [ ] Admin-web kategori & diskon CRUD work dengan public_id
- [ ] Mobile alamat CRUD work dengan public_id

## Branch Strategy

```bash
# Branch sama untuk semua repo yang berubah:
feat/public-id-migration
# Base: dev
# PR ke: dev
```

## Rollback Plan

Jika migrasi bermasalah:

1. Migration: `ALTER TABLE kategori DROP COLUMN public_id`, `ALTER TABLE diskon DROP COLUMN public_id`
2. Revert controller + routes ke versi sebelumnya
3. Frontend revert ke `item.id_kategori` / `item.id_diskon`
