# 📦 Integrasi & Pengujian Biteship (Sandbox)

---
### 🧭 Navigasi Cepat
[🏠 Utama](README.md) | [🏛️ Arsitektur](architecture.md) | [🛠️ Deployment](deployment.md) | [💳 Midtrans](pembayaran.md) | [📦 Biteship](biteship.md) | [📡 API Contract](api/overview.md) | [🗄️ Database](database/erd.md) | [🔒 Keamanan](security/README.md)
---

**Sistem:** Golang (Backend), Flutter (Frontend), & Biteship API Gateway  
**Status Integrasi:** **SELESAI & MERGED** (Sudah digabungkan ke branch `dev`).

Dokumen ini berisi rangkuman status akhir integrasi Biteship, hasil pengujian sandbox, serta panduan setup mandiri untuk keperluan testing lokal maupun staging.

---

## 🧪 1. Hasil Pengujian Sandbox

Integrasi Biteship telah diuji secara end-to-end dengan hasil sebagai berikut:

*   **Cek Ongkir Eksternal (`POST /customer/ongkir/cek`):** Berhasil memanggil API Rates Biteship untuk menghitung tarif ongkos kirim berbagai ekspedisi (JNE, J&T, SiCepat, AnterAja, Ninja Xpress).
    *   *Fallback:* Jika saldo sandbox Anda habis/bernilai 0, adapter otomatis melakukan fallback ke sistem Mock Ongkir lokal agar testing aplikasi tidak terganggu.
*   **Buat Pengiriman Otomatis (`processExternalShipment`):** Berhasil membuat shipment baru di Biteship secara otomatis setelah:
    *   Transaksi non-tunai lunas via Webhook Midtrans (`settlement`).
    *   Kasir menyelesaikan transaksi tunai (cash).
*   **Lacak Resi (`GET /customer/pesanan/:id/lacak`):** Berhasil mengembalikan riwayat pengiriman (tracking history) berbasis teks dari Biteship.
*   **Update Resi & Status Otomatis (`POST /webhook/biteship`):** Berhasil memproses callback dari Biteship untuk mencatat `waybill_id` (resi) dan memperbarui status pesanan menjadi `Dikirim` / `Selesai` secara asinkron.

---

## 🌐 2. Panduan Setup Mandiri (Sandbox Testing)

Karena Biteship Sandbox bergantung pada akun developer Anda, berikut adalah langkah-langkah manual yang harus dilakukan oleh developer:

### A. Konfigurasi Environment Variables (`.env`)
Pastikan file `.env` di `mantra-backend` dikonfigurasi dengan credentials sandbox Anda:
```env
# Biteship Credentials
BITESHIP_API_KEY=biteship_test.xxx... (Gunakan API Key Sandbox dari Dashboard Biteship Anda)
BITESHIP_MODE=sandbox

# Alamat Pengirim (Toko) - Wajib valid agar penghitungan ongkir akurat
BITESHIP_STORE_NAME=Toko Mantra
BITESHIP_STORE_PHONE=081234567890
BITESHIP_STORE_ADDRESS=Jl. Prof. Soedarto, S.H.
BITESHIP_STORE_CITY=Semarang
BITESHIP_STORE_POSTAL_CODE=50275
BITESHIP_STORE_REGION=Jawa Tengah
BITESHIP_STORE_COORDINATE_LAT=-7.046389
BITESHIP_STORE_COORDINATE_LONG=110.438333

# Secret untuk Webhook (Opsional, untuk verifikasi signature request)
BITESHIP_WEBHOOK_SECRET=your_webhook_secret_here
```

### B. Isi Saldo Virtual di Dashboard Biteship
API Cek Ongkir (Rates) dan Lacak Resi (Tracking) pada Biteship Sandbox tetap memotong saldo virtual akun Anda (**Rp5/req** untuk Rates, **Rp10/req** untuk Tracking).
1. Masuk ke [Dashboard Biteship](https://dashboard.biteship.com).
2. Switch ke mode **Sandbox / Testing**.
3. Lakukan **Top Up Saldo Virtual** agar API tidak menghasilkan error *No sufficient balance*.

### C. Konfigurasi Webhook via Ngrok (Update Status & Resi Otomatis)
Biteship mengirimkan nomor resi (`waybill_id`) dan status kurir secara asinkron via webhook. Untuk testing di komputer lokal:
1. Jalankan **Ngrok** pada port backend Anda (misal port 8080):
   ```bash
   ngrok http 8080
   ```
2. Salin URL publik HTTPS yang diberikan Ngrok (contoh: `https://abcd-123.ngrok-free.app`).
3. Masuk ke Dashboard Biteship Sandbox > menu **Integrations** > **Webhook**.
4. Masukkan URL Webhook:
   ```text
   https://<DOMAIN-NGROK-ANDA>/api/v1/webhook/biteship
   ```
5. Aktifkan event: `order.status` dan `order.waybill_id` lalu simpan.

> [!WARNING]
> Setiap kali Ngrok dijalankan ulang, URL publik akan berubah. Anda perlu memperbarui URL Webhook di dashboard Biteship (dan dashboard Midtrans) dengan domain baru tersebut.

---

## 🗄️ 3. Database Migration

Pastikan tabel `pesanan` Anda sudah memiliki kolom `biteship_order_id`.
*   Secara default, GORM AutoMigrate pada file `config/database.go` akan mendeteksi struct `models.Pesanan` dan menambahkan kolom tersebut secara otomatis pada lingkungan development.
*   Jika database Anda tidak menggunakan AutoMigrate (misalnya di staging/produksi), Anda wajib menjalankan migrasi manual menggunakan script berikut:
    `mantra-backend/migrations/add_biteship_order_id_to_pesanan.sql`
    ```sql
    ALTER TABLE pesanan ADD COLUMN IF NOT EXISTS biteship_order_id VARCHAR(100);
    CREATE INDEX IF NOT EXISTS idx_pesanan_biteship_order_id ON pesanan(biteship_order_id);
    ```
