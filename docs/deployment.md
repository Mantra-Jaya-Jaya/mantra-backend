# 🛠️ Alur & Panduan Deployment

---
### 🧭 Navigasi Cepat
[🏠 Utama](README.md) | [🏛️ Arsitektur](architecture.md) | [🛠️ Deployment](deployment.md) | [💳 Midtrans](pembayaran.md) | [📦 Biteship](biteship.md) | [📡 API Contract](api/overview.md) | [🗄️ Database](database/erd.md) | [🔒 Keamanan](security/README.md)
---

Dokumen ini menjelaskan proses build aplikasi backend MANTRA ke file binary executable, daftar konfigurasi environment variable yang wajib disiapkan, serta checklist sebelum menjalankan sistem di lingkungan produksi.

---

## 📦 1. Build Binary Aplikasi

Backend MANTRA ditulis menggunakan bahasa Go, yang menghasilkan satu file binary statis yang efisien dan siap dijalankan di server tanpa dependensi runtime eksternal.

Jalankan perintah berikut di direktori `mantra-backend` untuk melakukan proses kompilasi:

```bash
make build
```

**Hasil Build:**
- File binary akan dibuat di path: `bin/mantra-backend`.
- Binary ini dapat langsung dipindahkan (copy) ke server produksi (VPS/VM).

---

## 🛠️ 2. Environment Variables (.env)

Konfigurasi aplikasi diatur sepenuhnya menggunakan environment variables. Untuk development, Anda dapat menyalin file `.env.example` ke file `.env` di direktori root backend, lalu isi nilainya secara tepat:

### A. Database (PostgreSQL)
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_secure_password
DB_NAME=mantra_db
```

### B. Keamanan & Autentikasi (JWT)
```env
JWT_SECRET=generate_sebuah_random_secret_yang_panjang_di_sini
JWT_REFRESH_SECRET=generate_sebuah_secret_yang_berbeda_untuk_refresh_token
```

### C. Midtrans Payment Gateway
```env
MIDTRANS_SERVER_KEY=your_server_key_here
MIDTRANS_CLIENT_KEY=your_client_key_here
MIDTRANS_ENVIRONMENT=sandbox # Ubah ke 'production' di server live
MERCHANID=your_merchant_id_here
```

### D. Biteship Shipping Gateway
```env
BITESHIP_API_KEY=your_biteship_api_key_here
BITESHIP_MODE=sandbox # Ubah ke 'production' di server live
BITESHIP_STORE_NAME=Toko Mantra
BITESHIP_STORE_PHONE=081234567890
BITESHIP_STORE_ADDRESS=Jl. Prof. Soedarto, S.H., Tembalang
BITESHIP_STORE_CITY=Semarang
BITESHIP_STORE_POSTAL_CODE=50275
BITESHIP_STORE_REGION=Jawa Tengah
BITESHIP_STORE_COORDINATE_LAT=-7.046389
BITESHIP_STORE_COORDINATE_LONG=110.438333
BITESHIP_WEBHOOK_SECRET=your_webhook_secret_here
```

### E. MinIO (S3-Compatible Object Storage)
```env
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=your_minio_access_key
MINIO_SECRET_KEY=your_minio_secret_key
MINIO_BUCKET=mantra-storage
```

### F. CORS (Cross-Origin Resource Sharing)
```env
ALLOWED_ORIGIN=https://admin.mantra.web.id # Domain front-end Next.js admin web
```

---

## 🚀 3. Checklist Kesiapan Produksi (Production Checklist)

Sebelum menjalankan aplikasi di server live/produksi, pastikan poin-poin checklist berikut telah terpenuhi secara menyeluruh:

- [ ] **Kunci Keamanan JWT:** Ganti `JWT_SECRET` dan `JWT_REFRESH_SECRET` dengan string acak yang kuat (disarankan minimal 32 karakter alfanumerik).
- [ ] **Akses Publik MinIO:** Pastikan bucket `mantra-storage` telah dibuat pada instance MinIO server target, dan atur policy bucket ke mode *Public* atau *Read-Only* agar file/gambar upload dapat diakses oleh client app.
- [ ] **Konfigurasi Produksi Gateway:**
  - Ganti Server/Client key Midtrans ke mode **Production**, dan ubah `MIDTRANS_ENVIRONMENT=production`.
  - Ganti API Key Biteship ke mode **Production**, sesuaikan parameter alamat/koordinat fisik toko asli, dan ganti `BITESHIP_MODE=production`.
- [ ] **CORS Security:** Atur `ALLOWED_ORIGIN` secara spesifik ke domain Admin Web Anda (jangan gunakan `*`).
- [ ] **SSL / HTTPS Proxy:** Gunakan reverse proxy seperti **Nginx** atau **Caddy** di depan port backend (`:8080`) untuk menangani enkripsi SSL/TLS (HTTPS).
- [ ] **Process Manager:** Jalankan file binary menggunakan process manager seperti **systemd**, **Supervisor**, atau **PM2** agar backend dapat otomatis hidup kembali jika terjadi crash atau server reboot.
- [ ] **Backup Database:** Siapkan cron job berkala untuk melakukan backup database PostgreSQL (`pg_dump`) secara berkala.
