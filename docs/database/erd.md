# 🗄️ Database ERD & Konvensi Skema

---
### 🧭 Navigasi Cepat
[🏠 Utama](../README.md) | [🏛️ Arsitektur](../architecture.md) | [🛠️ Deployment](../deployment.md) | [💳 Midtrans](../pembayaran.md) | [📦 Biteship](../biteship.md) | [📡 API Contract](../api/overview.md) | [🗄️ Database](erd.md) | [🔒 Keamanan](../security/README.md)
---

Dokumen ini berisi informasi mengenai visualisasi ERD database MANTRA (terdiri dari 25+ tabel), konvensi penamaan tabel dan kolom, serta panduan praktis untuk query data.

---

## 🖼️ 1. Sumber Skema ERD (Visualisasi)

Skema database dimodelkan secara visual menggunakan tools **dbdiagram.io** menggunakan berkas-berkas berikut:

| Nama Berkas | Tipe / Format | Deskripsi & Cara Membuka |
| :--- | :--- | :--- |
| [**mantra.dbml**](../mantra.dbml) | DBML (Database Markup Language) | Plain-text markup skema database. Dapat dibuka dengan ekstensi DBML di VS Code, atau diunggah langsung ke [dbdiagram.io](https://dbdiagram.io/). |
| [**mantra.dbdiagram**](../mantra.dbdiagram) | JSON Layout | Berisi pengaturan koordinat visual diagram. Dapat di-import di menu Dashboard dbdiagram.io untuk mengembalikan tata letak grafik diagram. |

---

## 📝 2. Konvensi Penamaan (Naming Conventions)

### A. Tabel
*   Format penulisan menggunakan huruf kecil semua dengan pemisah underscore (**snake_case**) dan bersifat tunggal (singular).  
    *Contoh:* `role`, `user`, `refresh_token`, `spesifikasi_barang`, `stok_opname`, `ekspedisi_layanan`, `metode_pembayaran`.
*   Nama tabel pada struct model Golang disesuaikan dengan mengimplementasikan interface method `TableName()` bawaan GORM.

### B. Kolom & Kunci (Fields & Keys)
*   **Case Format:** Menggunakan **snake_case** (huruf kecil, dipisah underscore).  
    *Contoh:* `id_role`, `nama_barang`, `tanggal_pesanan`.
*   **Primary Key:** Ditulis dengan pola `id_{nama_tabel}`.  
    *Contoh:* `id_user`, `id_kategori`, `id_pesanan`.
*   **Foreign Key:** Ditulis dengan mencocokkan pola Primary Key tabel yang dituju.  
    *Contoh:* Kolom `id_role` di dalam tabel `user`.
*   **Timestamps:** Kolom tracking waktu wajib diberi nama `created_at` dan `updated_at`.

### C. Tipe Data
*   **Nilai Keuangan / Uang:** Wajib menggunakan tipe data **Integer** (`int`) untuk menyimpan nominal Rupiah. Dilarang keras menggunakan tipe data float untuk uang demi menghindari pembulatan nilai pecahan.
*   **External ID:** Wajib menggunakan tipe data **UUID** (`public_id`) yang ter-generate otomatis secara default via `gen_random_uuid()` pada sisi database.
*   **Koordinat Maps:** Menggunakan tipe data **Float / Double** untuk field latitude dan longitude.
*   **Flag / Status:** Menggunakan tipe data **Boolean** (`is_active`, `is_utama`).

---

## 🔄 3. Representasi Relasi DBML

Berikut adalah simbol relasi yang digunakan dalam memodelkan hubungan antar tabel:

| Tipe Relasi | Simbol | Contoh Kasus |
| :--- | :--- | :--- |
| **One-to-Many** | `>` | `user > refresh_token` (satu user bisa punya banyak refresh token) |
| **Many-to-One** | `<` | `pesanan < detail_pesanan` (banyak item pesanan merujuk ke satu pesanan) |
| **One-to-One** | `-` | `karyawan - kasir` (satu karyawan diposisikan sebagai satu kasir) |

---

## 💡 4. Tips Melakukan Query Database

*   **Penyaringan Notifikasi:** Notifikasi harus selalu difilter menggunakan query `id_user` yang diambil dari token JWT untuk menjamin keamanan privasi.
*   **Eager Loading:** Gunakan fitur `.Preload()` bawaan GORM untuk melakukan load data relasional secara efisien dan mencegah masalah query N+1.
*   **Lookups:** Gunakan kolom `public_id` (UUID) untuk query lookup data dari arah router publik (client app). Gunakan `id_pesanan` integer internal hanya untuk pemrosesan/relasi internal database.
*   **Snapshot Invoice:** Selalu simpan harga barang di tabel `detail_pesanan` sebagai snapshot harga saat check out. Hal ini mencegah nilai invoice lama berubah ketika admin mengubah harga barang di katalog produk.
