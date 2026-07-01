# 💳 Integrasi Pembayaran Midtrans (Core API)

---
### 🧭 Navigasi Cepat
[🏠 Utama](README.md) | [🏛️ Arsitektur](architecture.md) | [🛠️ Deployment](deployment.md) | [💳 Midtrans](pembayaran.md) | [📦 Biteship](biteship.md) | [📡 API Contract](api/overview.md) | [🗄️ Database](database/erd.md) | [🔒 Keamanan](security/README.md)
---

**Metode Pembayaran:** QRIS & Virtual Account (BCA, BNI, BRI, Permata)  
**Sistem:** Golang (Backend) & Flutter (Frontend)

Dokumentasi ini berisi panduan untuk melakukan setup, konfigurasi sandbox, penggunaan Ngrok sebagai jembatan webhook, serta alur pengujian fitur pembayaran non-tunai pada sistem MANTRA.

---

## 🌐 1. Setup & Instalasi Ngrok (Jembatan Webhook)

Karena server Midtrans berada di internet, server mereka tidak dapat mengirim notifikasi sukses pembayaran (webhook) ke backend yang berjalan di `localhost` komputer Anda. Kita memerlukan **Ngrok** sebagai terowongan (tunneling) untuk meneruskan request tersebut.

### A. Registrasi & Download
1. Kunjungi website resmi Ngrok: [https://ngrok.com/](https://ngrok.com/).
2. Buat akun baru (disarankan menggunakan opsi **Sign up with Google** agar lebih cepat).
3. Di dashboard Ngrok, unduh package zip Ngrok sesuai dengan Sistem Operasi Anda (Windows/Linux/MacOS).
4. Ekstrak file zip tersebut dan taruh binary `ngrok` (atau `ngrok.exe`) ke direktori yang mudah diakses.

### B. Konfigurasi Auth Token
1. Pada dashboard Ngrok web, buka menu **Getting Started** → **Your Authtoken**.
2. Salin token panjang yang disediakan.
3. Buka terminal (Bash/PowerShell/CMD) dan jalankan perintah:
   ```bash
   ngrok config add-authtoken TOKEN_ANDA_DI_SINI
   ```
4. Jika berhasil, pesan `Authtoken saved to configuration file` akan muncul.

---

## 🛠️ 2. Setup Environment Variables (.env)

Lakukan konfigurasi credentials Sandbox Midtrans Anda pada file `.env` backend:

1. Login ke [Midtrans Sandbox Dashboard](https://dashboard.sandbox.midtrans.com/).
2. Masuk ke menu **Settings** > **Access Keys**.
3. Salin Server Key, Client Key, dan Merchant ID.
4. Tambahkan konfigurasi berikut ke file `.env` backend Golang Anda:
   ```env
   MIDTRANS_SERVER_KEY=Mid-server-SB-XXXXXXXXXX
   MIDTRANS_CLIENT_KEY=Mid-client-SB-XXXXXXXXXX
   MIDTRANS_ENVIRONMENT=sandbox
   MERCHANID=GXXXXXXXXX
   ```

> [!IMPORTANT]
> Pastikan API Key yang Anda gunakan diawali dengan `Mid-server-SB-` untuk memastikan Anda berada dalam mode Sandbox (Testing).

---

## 🚀 3. Menjalankan Ngrok & Update Webhook Midtrans

Langkah ini bertujuan untuk mendaftarkan URL tujuan notifikasi ketika pelanggan telah menyelesaikan pembayarannya.

1. Buka terminal baru dan jalankan Ngrok pada port backend Anda (default: `8080`):
   ```bash
   ngrok http 8080
   ```
2. Salin URL publik HTTPS yang tampil di bagian **Forwarding** (contoh: `https://abcd-123.ngrok-free.app`).
   > [!WARNING]
   > Biarkan terminal Ngrok ini tetap berjalan selama Anda melakukan testing pembayaran!
3. Masuk ke Dashboard Midtrans Sandbox > menu **Settings** > **Configuration**.
4. Cari kolom **Payment Notification URL**.
5. Masukkan alamat webhook Anda dengan format berikut:
   ```text
   https://<DOMAIN-NGROK-ANDA>/api/v1/payment/notification
   ```
   *Contoh:* `https://abcd-123.ngrok-free.app/api/v1/payment/notification`
6. Gulir ke bawah halaman dan klik tombol **Update / Save**.

---

## 🧪 4. Alur Pengujian Pembayaran (End-to-End Testing)

Gunakan alur berikut untuk menyimulasikan transaksi non-tunai secara penuh:

### Langkah 1: Jalankan Backend & Ngrok
- Pastikan server backend Golang Anda sudah menyala (`go run main.go` atau `make run`).
- Pastikan Ngrok telah menyala dan alamat notifikasi di dashboard Midtrans sudah diperbarui dengan URL Ngrok terbaru.

### Langkah 2: Lakukan Checkout dari Aplikasi Mobile/Kasir
- Masukkan produk ke keranjang belanja.
- Pilih metode pembayaran non-tunai (misal: QRIS atau BCA Virtual Account).
- Selesaikan pesanan. Aplikasi akan menampilkan gambar QR Code (untuk QRIS) atau nomor rekening virtual (untuk VA).

### Langkah 3: Simulasi Pembayaran di Midtrans Simulator
- Buka browser dan kunjungi [Midtrans Sandbox Simulator](https://payment-simulator.sandbox.midtrans.com/).
- **Jika QRIS:** Pilih tab QRIS, salin URL gambar QRIS yang didapat dari API backend, tempel di simulator, lalu klik bayar.
- **Jika Virtual Account:**
  - Pilih bank yang sesuai (misal: BCA).
  - Masukkan nomor rekening virtual yang tampil pada aplikasi kasir/customer.
  - Klik **Inquire**, kemudian klik **Pay**.

### Langkah 4: Verifikasi Status Otomatis
- Kembali ke aplikasi klien Anda. Status pesanan akan otomatis terupdate.
- Pada log terminal backend Golang, Anda akan melihat request masuk ke `POST /api/v1/payment/notification` dengan response `200 OK`.
- Status transaksi di database backend akan berubah menjadi `settlement` (Lunas).

---

## 🐛 Troubleshooting

*   **Error 402 (Payment Channel Not Activated):**
    *   *Penyebab:* Tipe pembayaran (QRIS/VA) yang di-request belum diaktifkan di dashboard akun sandbox Anda.
    *   *Solusi:* Masuk ke dashboard Midtrans web, klik menu **Settings** > **Payment Methods**, lalu centang dan aktifkan metode pembayaran yang ingin digunakan.
*   **Status Transaksi Tidak Berubah Menjadi Lunas:**
    *   Pastikan terminal Ngrok masih aktif dan tidak terputus.
    *   Pastikan URL Notifikasi di dashboard Midtrans sama persis dengan URL Ngrok yang sedang berjalan (ingat bahwa URL Ngrok berubah setiap kali program Ngrok direstart).
    *   Periksa log terminal backend untuk melihat apakah ada error parsing notifikasi atau error koneksi DB.
