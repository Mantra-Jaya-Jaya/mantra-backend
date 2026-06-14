# 🚀 Dokumentasi Integrasi Pembayaran Midtrans (Core API)

**Metode Pembayaran:** QRIS & Virtual Account (BCA, BNI, BRI)  
**Sistem:** Golang (Backend) & Flutter (Frontend)

Dokumentasi ini berisi langkah-langkah untuk melakukan setup dan testing fitur pembayaran non-tunai menggunakan Midtrans Core API pada mode Sandbox.

---

## 🌐 1. Setup & Instalasi Ngrok (Jembatan Webhook)

Karena server Midtrans berada di internet, Midtrans tidak bisa mengirimkan notifikasi lunas ke Backend yang berjalan di `localhost`. Kita membutuhkan **Ngrok** untuk melakukan port-forwarding.

### A. Daftar dan Login

1. Buka website resmi Ngrok: [https://ngrok.com/](https://ngrok.com/)
2. Klik tombol **Sign Up** (atau Login jika sudah punya akun). Paling cepat gunakan opsi **Sign up with Google**
3. Setelah berhasil login, Anda akan masuk ke halaman Dashboard Ngrok

### B. Download Ngrok

1. Di halaman Dashboard, pada menu **Setup & Installation**, pilih OS yang Anda gunakan (Windows/Linux/Mac)
2. Download file `.zip` Ngrok
3. Ekstrak file `.zip` tersebut. Di dalamnya ada file aplikasi `ngrok` (atau `ngrok.exe` untuk Windows)
4. Pindahkan file `ngrok` ini ke folder yang mudah diakses, atau masukkan ke Environment Variables (PATH) OS Anda agar bisa dipanggil dari terminal mana saja

### C. Verifikasi Auth Token

1. Kembali ke Dashboard Ngrok di browser
2. Di menu sebelah kiri, klik **Getting Started** → **Your Authtoken**
3. Copy token panjang yang ada di halaman tersebut (klik icon Copy)
4. Buka terminal (Command Prompt/PowerShell/Bash), lalu ketikkan perintah berikut dan paste token Anda:

```bash
ngrok config add-authtoken TOKEN_ANDA_DI_SINI
```

5. Tekan Enter. Jika berhasil, akan muncul tulisan `Authtoken saved to configuration file`. Setup Ngrok selesai!
## 🛠️ 2. Setup Environment Variables (.env)

Pastikan Anda sudah memiliki akun Midtrans. Kita akan menggunakan Sandbox Environment untuk testing.

1. Login ke Midtrans Dashboard
2. Masuk ke menu **Settings** > **Access Keys**
3. Copy Server Key Sandbox Anda
4. Buka file `.env` di project Backend (Golang) Anda, dan tambahkan/ubah baris berikut:

```env
# Konfigurasi Midtrans
MIDTRANS_SERVER_KEY=Mid-server-XXXXXXXXXX
MIDTRANS_ENVIRONMENT=sandbox
MIDTRANS_CLIENT_KEY=Mid-XXXXXXXXXX
MERCHANID=XXXXXXXXX
```

> **Catatan:** Pastikan `MIDTRANS_SERVER_KEY` menggunakan awalan `SB-` yang menandakan akun Sandbox
## 🚀 3. Menjalankan Ngrok & Update Webhook Midtrans

Langkah ini bertujuan untuk memberi tahu Midtrans ke mana mereka harus mengirim pesan/notifikasi jika pelanggan sudah membayar (Settlement).

1. Buka terminal baru (biarkan terminal Backend Golang tetap berjalan di port 8080)

2. Jalankan perintah berikut:

```bash
ngrok http 8080
```

3. Tunggu hingga layar Ngrok muncul, lalu copy URL HTTPS yang ada di bagian **Forwarding** (Contoh: `https://a1b2-c3d4.ngrok-free.app`)

   > **⚠️ PENTING:** Jangan tutup terminal Ngrok ini selama proses testing berlangsung!

4. Buka Midtrans Sandbox Dashboard

5. Buka menu **Settings** > **Configuration** (Pengaturan > Pengaturan Umum)

6. Cari kolom **Payment Notification URL**

7. Masukkan URL Ngrok Anda ditambah dengan endpoint webhook Backend:

```
https://<ALAMAT-NGROK-ANDA>/api/v1/payment/notification
```

Contoh: `https://a1b2-c3d4.ngrok-free.app/api/v1/payment/notification`

8. Scroll ke bawah dan klik tombol **Update** / **Simpan**
## 🧪 4. Alur Testing (End-to-End)

Ikuti langkah-langkah ini untuk mensimulasikan pembayaran dari awal hingga selesai:

### Step 1: Jalankan Backend
- Pastikan server Golang sudah berjalan (`make run`)

### Step 2: Jalankan Ngrok
- Pastikan Ngrok sudah menyala dan URL sudah di-update di Midtrans

### Step 3: Jalankan Aplikasi Kasir (Flutter)
1. Masukkan barang ke keranjang
2. Pilih menu pembayaran Non-Tunai
3. Pilih metode pembayaran (misal: QRIS atau BCA Virtual Account)
4. Klik **Konfirmasi Pembayaran**. Aplikasi akan menampilkan URL Gambar QRIS atau Nomor Virtual Account

### Step 4: Simulasi Pembayaran
1. Buka browser dan masuk ke Midtrans Simulator
2. Jika QRIS: Pilih QRIS dan ikuti instruksi sukses
3. Jika VA: 
   - Pilih Bank yang sesuai (misal: BCA)
   - Masukkan nomor VA yang tampil di aplikasi kasir
   - Klik **Inquire**, lalu klik **Pay**

### Step 5: Cek Status di Aplikasi
1. Kembali ke aplikasi kasir (Flutter)
2. Klik tombol **Cek Status Pembayaran**
3. Jika berhasil, aplikasi akan otomatis memunculkan pop-up Sukses dan beralih ke halaman nota! 🎉
## 🐛 Troubleshooting

### Error 402 (Payment Channel Not Activated)

**Masalah:** Terjadi karena metode pembayaran belum aktif di akun Midtrans Sandbox Anda.

**Solusi:** 
- Buat akun email baru
- Daftar ulang Midtrans
- Gunakan Server Key yang baru

### Status Tidak Berubah Menjadi "Selesai"

**Checklist:**
- ✅ Pastikan terminal Ngrok masih menyala
- ✅ Pastikan **Payment Notification URL** di Midtrans sudah menggunakan URL Ngrok terbaru (URL Ngrok berubah setiap kali direstart)
- ✅ Cek log terminal Golang, pastikan ada notifikasi masuk `POST /api/v1/payment/notification` dengan status `200 OK`