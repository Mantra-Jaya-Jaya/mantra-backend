# 🏛️ Arsitektur Sistem & Hubungan Model

---
### 🧭 Navigasi Cepat
[🏠 Utama](README.md) | [🏛️ Arsitektur](architecture.md) | [🛠️ Deployment](deployment.md) | [💳 Midtrans](pembayaran.md) | [📦 Biteship](biteship.md) | [📡 API Contract](api/overview.md) | [🗄️ Database](database/erd.md) | [🔒 Keamanan](security/README.md)
---

Halaman ini menjelaskan teknologi, alur request data, struktur relasi database, dan keputusan arsitektural penting yang diambil dalam pembangunan **Mantra Backend**.

---

## 🛠️ Tech Stack

Mantra Backend dibangun menggunakan komponen-komponen software berikut:

| Layer | Teknologi | Deskripsi / Versi |
| :--- | :--- | :--- |
| **Language** | Go | v1.26.2 - Bahasa pemrograman utama |
| **Framework** | Gin | v1.12.0 - HTTP Web Framework untuk REST API |
| **ORM** | GORM | v1.31.1 - Library pemetaan relasi objek ke database |
| **Database** | PostgreSQL | Sistem database relasional utama |
| **Auth** | JWT | golang-jwt v5 - Autentikasi berbasis token |
| **Storage** | MinIO | S3-compatible Object Storage untuk upload media/file |
| **Hot Reload** | Air | Utilitas Live-reload untuk development server |
| **Migration** | Atlas CLI | `atlas-provider-gorm` - Manajemen skema database deklaratif |

---

## 📡 Hubungan Komponen (Component Interaction)

Diagram berikut menjelaskan hubungan komunikasi antara aplikasi client (Mobile & Admin Web), server backend, database, dan storage:

```text
┌─────────────────┐     ┌──────────────────┐     ┌────────────────┐
│  Flutter App    │────▶│                  │────▶│   PostgreSQL   │
│  - Customer     │     │                  │     │   Database     │
│  - Kasir        │     │     Gin HTTP     │     └────────────────┘
├─────────────────┤     │     Server       │
│  Next.js App    │────▶│     :8080        │     ┌────────────────┐
│  - Admin Web    │     │                  │────▶│     MinIO      │
└─────────────────┘     │                  │     │    Storage     │
                        └──────────────────┘     └────────────────┘
```

---

## 🔄 Alur Request (Request Flow)

Setiap request HTTP yang masuk ke server backend akan diproses secara berurutan sesuai alur berikut:

```text
Request HTTP
   │
   ▼
┌──────────────────┐
│    GIN Router    │  routes/routes.go → Grouping berdasarkan prefix (/api/v1/...)
└────────┬─────────┘
         │
         ▼
┌──────────────────┐
│  AuthMiddleware  │  Validasi token JWT & menyegarkan masa aktif token (sliding expiration).
│  (Opsional)      │  Menyimpan user_id, public_id, dan role ke context request.
└────────┬─────────┘
         │
         ▼
┌──────────────────┐
│  RoleMiddleware  │  Memverifikasi hak akses berdasarkan role (misal: Admin-only).
│  (Opsional)      │
└────────┬─────────┘
         │
         ▼
┌──────────────────┐
│    Controller    │  Tempat logika bisnis diproses:
│   (Per Domain)   │  - Melakukan query DB via GORM model.
│                  │  - Mengunggah file ke MinIO via config.MinioClient.
└────────┬─────────┘
         │
         ▼
Response JSON (Format Sukses / Error)
```

---

## 🗄️ Relasi Model Database

Berikut adalah pemetaan relasi konseptual antar model (tabel) pada database MANTRA:

```text
role ───< user ───< customer
               ├──< karyawan ───< kasir
               │               └──< kurir
               └──< refresh_token
              
customer ───< alamat
customer ───< keranjang ───> spesifikasi_barang
customer ───< pesanan ───< detail_pesanan ───> spesifikasi_barang
pesanan  ───< pembayaran ───< detail_pembayaran
pesanan  ───< pengantaran ───> kurir
                             ───> ekspedisi
                             ───> status_pengantaran
pesanan  ───> ekspedisi
pesanan  ───> ekspedisi_layanan
pembayaran ──> metode_pembayaran

ekspedisi ───< ekspedisi_layanan

kategori ───< barang ───< spesifikasi_barang ───< barcode
satuan   ───< barang              └──> detail_spesifikasi ───> spesifikasi
diskon   ───< barang

spesifikasi_barang ───< stok_opname

user ───< notifikasi
```

---

## 💡 Keputusan Arsitektural Utama (Key Decisions)

Beberapa keputusan desain arsitektur penting yang digunakan di project ini:

*   **Single-tenant:** Satu instance dideploy untuk satu toko. Tidak memerlukan arsitektur multi-tenant/isolasi data kompleks.
*   **Integer untuk Keuangan:** Semua nominal uang menggunakan tipe data `int` (Rupiah) untuk mencegah error pembulatan floating-point.
*   **UUID public_id:** ID privat bertipe integer auto-increment tidak pernah dipublikasikan ke klien. Menggunakan UUID `public_id` pada URL API untuk keamanan.
*   **Karyawan sebagai Tabel Induk:** Model `Kasir` dan `Kurir` ditautkan ke tabel `Karyawan` untuk membagikan data dasar (nama, NIK, alamat, dll).
*   **Snapshot Harga di Detail Pesanan:** Menyimpan salinan harga produk saat transaksi dibuat, sehingga invoice masa lalu tidak berubah saat harga produk di-update.
*   **Ongkir via Biteship API:** Kalkulasi tarif pengiriman pihak ketiga secara real-time langsung memanggil API Biteship.
*   **Pembayaran via Midtrans Snap:** Keamanan dan validasi metode pembayaran QRIS, Virtual Account, dan E-Wallet diserahkan sepenuhnya ke Midtrans.
*   **WEBHOOK Midtrans & Biteship:** Pembaruan status pembayaran dan pengiriman dilakukan secara asinkron via Webhook dengan validasi hash SHA-512 & HMAC-SHA256 untuk keamanan.
