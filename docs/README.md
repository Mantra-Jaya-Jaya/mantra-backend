# Dokumentasi MANTRA Backend

Pusat dokumentasi teknis backend service **MANTRA**.

## API Contract

Dokumentasi endpoint API diorganisir per domain — mirror struktur `controllers/`.

| Domain | File | Method | Role |
|--------|------|--------|------|
| Auth | `api/auth.md` | Public + Auth | Semua |
| Katalog | `api/katalog/barang.md` | Public + Auth | Customer, Kasir, Admin |
| Katalog | `api/katalog/kategori.md` | Public + Auth | Customer, Kasir, Admin |
| Katalog | `api/katalog/diskon.md` | Public + Auth | Customer, Kasir, Admin |
| Keranjang | `api/keranjang.md` | Auth | Customer |
| Notifikasi | `api/notifikasi.md` | Auth | Customer, Kasir, Admin |
| Pengantaran | `api/pengantaran.md` | Auth | Customer |
| Stok | `api/stok.md` | Auth | Admin |
| Transaksi | `api/transaksi/pesanan.md` | Auth | Customer, Kasir |
| Transaksi | `api/transaksi/pembayaran.md` | Auth | Kasir |
| User | `api/user/customer.md` | Auth | Customer |
| User | `api/user/kasir.md` | Auth | Kasir |
| User | `api/user/admin.md` | Auth | Admin |
| User | `api/user/karyawan.md` | Auth | Admin |
| User | `api/user/alamat.md` | Auth | Customer |

> Format umum request/response lihat `api/overview.md`.

## Arsitektur & Keamanan

| Dokumen | Isi |
|---------|-----|
| `architecture.md` | Overview arsitektur, decision records, component interaction |
| `security/README.md` | Index dokumentasi security |
| `security/jwt.md` | JWT flow, claims structure, sliding expiration |
| `security/middleware.md` | AuthMiddleware, RoleMiddleware, OwnershipMiddleware |
| `security/rbac.md` | RBAC matrix per endpoint & role |
| `security/errors.md` | Error codes, HTTP status, response format |

## Database

| Dokumen | Isi |
|---------|-----|
| `database/erd.md` | Panduan baca ERD, konvensi tabel & kolom |
| `database/migrations.md` | Workflow Atlas CLI, cara nambah tabel |
| `mantra.dbml` | Skema ERD (source of truth — 25 tabel) |
| `mantra.dbdiagram` | Layout diagram dbdiagram.io |

## Deployment

| Dokumen | Isi |
|---------|-----|
| `deployment.md` | Build binary, environment variables, production checklist |

## API Collections (Bruno)

Folder `api-collections/` berisi koleksi request API untuk **Bruno** (open-source API client, alternative Postman).

| Path | Isi |
|------|-----|
| `api-collections/collection.bru` | Global settings (baseUrl, auth bearer) |
| `api-collections/environments/` | Environment variables (local) |
| `api-collections/Auth/` | Login, Register, Refresh, Logout, Change Password |
| `api-collections/Katalog/` | Promo, Kategori, Barang, Barcode |
| `api-collections/Keranjang/` | Tambah, Update, Hapus item |
| `api-collections/Transaksi - Customer/` | Pesanan, Checkout, Lacak |
| `api-collections/Transaksi - Kasir/` | Dashboard, Laporan, Bayar |
| `api-collections/Admin - Katalog/` | CRUD Kategori, Barang, Satuan |
| `api-collections/Admin - Diskon/` | CRUD Diskon + Upload Banner |
| `api-collections/Admin - Karyawan/` | CRUD Karyawan + Upload Foto |
| `api-collections/Admin - Dashboard/` | Ringkasan, Chart |
| `api-collections/Profil/` | Customer, Kasir, Admin |
| `api-collections/Alamat/` | CRUD Alamat |
| `api-collections/Notifikasi/` | Customer, Admin |

**Cara pakai:** Buka Bruno → Import collection → Pilih folder `docs/api-collections/`. |
