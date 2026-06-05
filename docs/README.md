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
| `security-rbac.md` | JWT flow, sliding expiration, middleware chain, RBAC matrix |

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
