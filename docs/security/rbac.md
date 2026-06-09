# RBAC Matrix

## Roles

| Role | Aplikasi | Registrasi |
|------|----------|------------|
| Customer | Flutter App | Mandiri via endpoint `/register` |
| Admin | Next.js Web | Dari seed data |
| Kasir | Flutter App | Dibuat oleh Admin (via `TambahKaryawan`) |
| Kurir | - | Dibuat oleh Admin (belum ada endpoint) |

## Endpoint Access

| Endpoint | Method | Public | Customer | Kasir | Admin |
|----------|--------|:------:|:--------:|:-----:|:-----:|
| **Auth** ||||||
| `/login` | POST | ✔ | - | - | - |
| `/register` | POST | ✔ | - | - | - |
| `/auth/refresh` | POST | ✔ | - | - | - |
| `/logout` | POST | - | ✔ | ✔ | ✔ |
| `/change-password` | PUT | - | ✔ | ✔ | ✔ |
| `/scan/:kode_barcode` | GET | ✔ | - | - | - |
| **Customer** ||||||
| `/customer/promo` | GET | - | ✔ | - | - |
| `/customer/kategori` | GET | - | ✔ | - | - |
| `/customer/barang` | GET | - | ✔ | - | - |
| `/customer/keranjang` | * | - | ✔ | - | - |
| `/customer/notifikasi` | GET | - | ✔ | - | - |
| `/customer/pesanan` | * | - | ✔ | - | - |
| `/customer/profil` | GET | - | ✔ | - | - |
| `/customer/akun` | PUT | - | ✔ | - | - |
| `/customer/alamat` | * | - | ✔ | - | - |
| **Kasir** ||||||
| `/kasir/dashboard` | GET | - | - | ✔ | - |
| `/kasir/laporan` | * | - | - | ✔ | - |
| `/kasir/pesanan` | * | - | - | ✔ | - |
| `/kasir/kategori` | GET | - | - | ✔ | - |
| `/kasir/transaksi` | * | - | - | ✔ | - |
| `/kasir/profil` | GET | - | - | ✔ | - |
| `/kasir/notifikasi` | GET | - | - | ✔ | - |
| **Admin** ||||||
| `/admin/dashboard` | * | - | - | - | ✔ |
| `/admin/kategori` | * | - | - | - | ✔ |
| `/admin/barang` | * | - | - | - | ✔ |
| `/admin/satuan` | GET | - | - | - | ✔ |
| `/admin/diskon` | * | - | - | - | ✔ |
| `/admin/karyawan` | * | - | - | - | ✔ |
| `/admin/notifikasi` | GET | - | - | - | ✔ |
| `/admin/profil` | * | - | - | - | ✔ |
| `/admin/ekspedisi` | * | - | - | - | ✔ |
| `/admin/ekspedisi/layanan/:id` | PUT, DELETE | - | - | - | ✔ |
| `/admin/metode-pembayaran` | * | - | - | - | ✔ |
| **Customer (baru)** ||||||
| `/customer/ongkir/cek` | POST | - | ✔ | - | - |
| `/customer/metode-pembayaran` | GET | - | ✔ | - | - |
| **Public (baru)** ||||||
| `/payment/notification` | POST | ✔ | - | - | - |

> `*` = Multiple HTTP methods (GET, POST, PUT, PATCH, DELETE)

## Ownership Rules

| Resource | Owner | Admin Access |
|----------|-------|:------------:|
| Profil customer | Customer sendiri | ✔ |
| Alamat | Customer pemilik alamat | ✔ |
| Keranjang | Customer pemilik keranjang | ✔ |
| Pesanan | Customer yang membuat | ✔ |
| Notifikasi | User pemilik notifikasi | ✔ |
| Semua resource admin | - | ✔ |
