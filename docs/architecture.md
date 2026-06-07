# Architecture

## Tech Stack

| Layer | Teknologi |
|-------|-----------|
| Language | Go 1.26.2 |
| Framework | Gin v1.12.0 |
| ORM | GORM v1.31.1 |
| Database | PostgreSQL |
| Auth | JWT (golang-jwt v5) |
| Storage | MinIO (S3-compatible) |
| Hot Reload | Air |
| Migration | Atlas CLI + atlas-provider-gorm |

## Component Interaction

```text
┌─────────────┐     ┌──────────────┐     ┌────────────┐
│  Flutter App │────▶│              │────▶│            │
│ (Customer)   │     │              │     │ PostgreSQL │
│ (Kasir)      │     │   Gin HTTP   │     │            │
├─────────────┤     │   Server     │     └────────────┘
│  Next.js    │────▶│   :8080      │
│ (Admin Web) │     │              │     ┌────────────┐
└─────────────┘     │              │────▶│   MinIO    │
                    └──────────────┘     │  Storage   │
                                         └────────────┘
```

## Request Flow

```text
Request
  │
  ▼
┌─────────────┐
│  GIN Router  │  routes/routes.go → group per prefix (/api/v1/...)
└──────┬──────┘
       │
       ▼
┌──────────────────┐
│ AuthMiddleware    │  Validasi JWT, sliding expiration
│ (jika dipasang)  │  Set user_id, public_id, role ke context
└──────┬───────────┘
       │
       ▼
┌──────────────────┐
│ RoleMiddleware    │  Cek role (jika dipasang)
│ (jika dipasang)  │  Contoh: hanya admin yang boleh akses
└──────┬───────────┘
       │
       ▼
┌──────────────────┐
│ Controller       │  Handler bisnis logic
│ (per domain)     │  → models (query DB via GORM)
│                  │  → config.MinioClient (upload file)
└──────┬───────────┘
       │
       ▼
     Response (JSON)
```

## Database Model Relationships

```text
role ───< user ───< customer
              ├──< karyawan ───< kasir
              │               └──< kurir
              └──< refresh_token
              
customer ───< alamat
customer ───< keranjang ───> spesifikasi_barang
customer ───< pesanan ───< detail_pesanan ───> spesifikasi_barang
pesanan  ───< pembayaran
pesanan  ───< pengantaran ───> kurir
                            ───> ekspedisi
                            ───> status_pengantaran

kategori ───< barang ───< spesifikasi_barang ───< barcode
satuan   ───< barang              └──> detail_spesifikasi ───> spesifikasi
diskon   ───< barang

spesifikasi_barang ───< stok_opname

user ───< notifikasi
```

## Key Decisions

| Keputusan | Alasan |
|-----------|--------|
| Single-tenant | Satu instalasi = satu toko, tidak perlu isolasi multi-toko |
| int untuk uang | Menghindari floating point error pada perhitungan harga |
| UUID public_id | Mencegah ID enumeration di endpoint publik |
| Karyawan sebagai induk | Kasir & Kurir berbagi data karyawan (nama, alamat, NIK, dll) |
| Snapshot harga di detail_pesanan | Harga barang bisa berubah, invoice harus tetap akurat |
