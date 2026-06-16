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
| Katalog | `api/katalog/ekspedisi.md` | Auth | Admin |
| Ongkir | `api/transaksi/ongkir.md` | Auth | Customer |
| Metode Pembayaran | `api/transaksi/metode_pembayaran.md` | Auth | Customer, Admin |
| Payment Notification | `api/transaksi/payment_notification.md` | Public (Webhook) | Midtrans |
| Keranjang | `api/keranjang.md` | Auth | Customer |
| Notifikasi | `api/notifikasi.md` | Auth | Customer, Kasir, Admin, Kurir |
| Pengantaran | `api/pengantaran.md` | Auth | Customer, Kurir |
| Stok | `api/stok.md` | Auth | Admin |
| Transaksi | `api/transaksi/pesanan.md` | Auth | Customer, Kasir |
| Transaksi | `api/transaksi/pembayaran.md` | Auth | Kasir |
| User | `api/user/customer.md` | Auth | Customer |
| User | `api/user/kasir.md` | Auth | Kasir |
| User | `api/user/kurir.md` | Auth | Kurir |
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
| `api-collections/environments/` | Environment variables (local, dev, production) |
| `api-collections/Public/` | Scan Barcode (no auth) |
| `api-collections/Customer/` | Login, Register, Logout, Change Password, Refresh Token, Profil, Notifikasi, Promo, Kategori, Barang, Keranjang, Pesanan, Alamat |
| `api-collections/Kasir/` | Login, Logout, Change Password, Refresh Token, Profil, Dashboard, Laporan, Transaksi |
| `api-collections/Kurir/` | Login, Profil, Notifikasi, Pesanan, Pengantaran |
| `api-collections/Admin/` | Login, Logout, Change Password, Refresh Token, Profil, Notifikasi, Dashboard, Kategori, Barang, Diskon, Karyawan, Satuan |
| `api-collections/Admin/Ekspedisi/` | CRUD Ekspedisi & Layanan (Tambah, Update, Hapus, Daftar) |
| `api-collections/Admin/Metode Pembayaran/` | CRUD Metode Pembayaran (Tambah, Update, Hapus, Daftar) |
| `api-collections/Customer/Ongkir/` | Cek Ongkos Kirim (via Biteship) |
| `api-collections/Customer/Metode Pembayaran/` | Daftar Metode Pembayaran Aktif |
| `api-collections/Public/Payment Notification/` | Midtrans Webhook (200, 400, 401) |

**Cara pakai:** Buka Bruno → Import collection → Pilih folder `docs/api-collections/`. |
