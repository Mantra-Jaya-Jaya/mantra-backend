# Database ERD

## Sumber Skema

Dua file berikut adalah representasi visual dari 28 tabel:

| File | Format | Cara Baca |
|------|--------|-----------|
| `mantra.dbml` | DBML (plain text) | Buka di VS Code + ekstensi DBML, atau upload ke dbdiagram.io |
| `mantra.dbdiagram` | JSON (layout) | Import ke dbdiagram.io untuk tampilan visual |

## Konvensi Penamaan

### Tabel

- Snake case, jamak: `role`, `user`, `refresh_token`, `spesifikasi_barang`, `stok_opname`, `ekspedisi_layanan`, `metode_pembayaran`, `detail_pembayaran`
- Nama tabel di struct GORM via method `TableName()`

### Kolom

- Snake case: `id_role`, `nama_barang`, `tanggal_pesanan`
- Primary key: `id_{tabel}` (contoh: `id_user`, `id_kategori`)
- Foreign key: `id_{tabel_referensi}` (contoh: `id_role` di tabel `user`)
- Timestamps: `created_at`, `updated_at`

### Tipe Data

- **Integer** untuk ID dan nilai uang (Rupiah). Tidak pernah pakai float untuk uang.
- **UUID** untuk `public_id` — digenerate otomatis via `gen_random_uuid()`.
- **Varchar** untuk string pendek, **text** untuk string panjang (deskripsi, alamat).
- **Date** untuk tanggal tanpa waktu (tanggal_lahir, tgl_mulai diskon).
- **Timestamp** untuk datetime (created_at, updated_at, tanggal_pesanan).
- **Boolean** untuk flag (status stok_opname, is_utama alamat).
- **Float** hanya untuk koordinat (latitude, longitude).

## Relasi

| Tipe | Simbol | Contoh |
|------|--------|--------|
| one-to-many | `>` | `user > refresh_token` |
| many-to-one | `<` | `pesanan < detail_pesanan` |
| one-to-one | `-` | `karyawan - kasir` |

## Tips Query

- Filter notifikasi selalu berdasarkan `id_user` dari JWT — jangan kirim semua notifikasi.
- Gunakan `Preload` untuk eager loading relasi (Role, User, dll).
- `public_id` adalah UUID — unik dan terindex. Gunakan untuk lookup di endpoint publik, bukan `id` integer.
- `harga_satuan` di `detail_pesanan` adalah snapshot — tidak otomatis berubah jika harga barang diupdate.
