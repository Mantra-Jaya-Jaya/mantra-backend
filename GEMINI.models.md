# GEMINI.models.md — MANTRA Backend: Revisi Models

> Kerjakan di branch baru sebelum mulai.
> Sumber kebenaran adalah **kode yang ada di `models/`**, bukan dokumentasi DBML
> (dokumentasi tersebut sudah lama tidak diupdate dan tidak bisa dipercaya).

---

## LANGKAH 0 — BUAT BRANCH BARU

```bash
git checkout -b refactor/models-revisi
```

---

## LANGKAH 1 — AUDIT MODELS YANG ADA (WAJIB SEBELUM APAPUN)

Baca semua file di folder `models/`. Untuk setiap struct, catat:
- Nama struct
- Nama field dan tipenya persis seperti di kode
- Tag `db:` yang dipakai

Tulis hasil bacaan ini sebagai referensi sebelum melakukan perubahan apapun.
**Jangan asumsikan nama field — baca dari kode.**

---

## LANGKAH 2 — TERAPKAN PERUBAHAN

Setelah audit selesai, terapkan **hanya 2 perubahan berikut**:

### Perubahan 1 — Struct `Barcode`

Temukan struct yang merepresentasikan tabel `barcode`. Ubah field FK-nya:
- Hapus field yang mengarah ke `barang` (apapun nama field-nya di kode)
- Ganti dengan field FK ke `spesifikasi_barang`

Nama field baru mengikuti konvensi yang sudah dipakai di struct lain di project
(cek tag `db:` struct lain sebagai acuan penamaan).

**Tujuan perubahan:** Scan barcode harus bisa langsung tarik harga dan stok dari
varian spesifik (`spesifikasi_barang`), bukan dari induk barang.

### Perubahan 2 — Struct `User`

Temukan struct yang merepresentasikan tabel `user`. Tambahkan satu field baru:
- Nama field: `FotoProfil` (atau ikuti konvensi penamaan field yang sudah ada)
- Tipe data: `string`
- Tag db: `foto_profil`
- Posisi: sebelum `CreatedAt` / `UpdatedAt` (ikuti pola struct yang ada)
- Nullable: ya — user lama tidak wajib punya foto

---

## LANGKAH 3 — BUAT SQL MIGRATION

Buat file migrasi sesuai tool yang dipakai di project (cek folder migrasi yang sudah
ada untuk tahu format dan konvensi penamaannya).

Isi migrasi:

```sql
-- Perubahan 1: Relasi barcode → spesifikasi_barang
-- Sesuaikan nama constraint dengan konvensi yang ada di migration sebelumnya
ALTER TABLE barcode DROP CONSTRAINT IF EXISTS <nama_constraint_lama>;
ALTER TABLE barcode DROP COLUMN IF EXISTS <nama_kolom_lama>;
ALTER TABLE barcode ADD COLUMN id_spesifikasi_barang INT NOT NULL;
ALTER TABLE barcode
    ADD CONSTRAINT barcode_id_spesifikasi_barang_fkey
    FOREIGN KEY (id_spesifikasi_barang)
    REFERENCES spesifikasi_barang(id_spesifikasi_barang);

-- Perubahan 2: Tambah foto_profil ke tabel user
ALTER TABLE "user" ADD COLUMN IF NOT EXISTS foto_profil VARCHAR;
```

---

## LANGKAH 4 — CEK DAMPAK

Setelah mengubah struct, cari semua file di project yang **referensi field lama**
dari struct `Barcode` yang baru saja diubah:

```bash
grep -r "id_barang" --include="*.go" .
```

Jika ditemukan referensi di controller atau query — **catat saja, jangan ubah**.
Perubahan controller dikerjakan di sesi terpisah.

---

## CHECKLIST

- [ ] Branch baru sudah dibuat
- [ ] Semua struct di `models/` sudah dibaca sebelum mulai coding
- [ ] Struct `Barcode` sudah diubah FK-nya ke `spesifikasi_barang`
- [ ] Struct `User` sudah punya field `FotoProfil`
- [ ] File migrasi SQL sudah dibuat dengan format yang konsisten dengan migrasi lain
- [ ] Referensi field lama di controller sudah dicatat (belum diubah)
- [ ] Tidak ada perubahan di luar folder `models/` dan folder migrasi

---

## YANG TIDAK BOLEH DILAKUKAN

- Jangan percaya dokumentasi DBML — baca kode yang ada
- Jangan ubah struct lain selain `Barcode` dan `User`
- Jangan ubah controller, route, atau middleware
- Jangan asumsikan nama field — selalu baca dari kode yang ada
