# 🔔 Notifikasi API Contract

---
### 🧭 Navigasi Cepat
[🏠 Utama](../README.md) | [🏛️ Arsitektur](../architecture.md) | [🛠️ Deployment](../deployment.md) | [💳 Midtrans](../pembayaran.md) | [📦 Biteship](../biteship.md) | [📡 API Contract](overview.md) | [🗄️ Database](../database/erd.md) | [🔒 Keamanan](../security/README.md)
---

Pusat kontrol notifikasi sistem MANTRA untuk mengabari pengguna mengenai status pesanan, pembayaran, atau notifikasi admin secara real-time.

---

## 🧭 Daftar Endpoint Notifikasi

*   [**GET /api/v1/customer/notifikasi**](#get-apiv1customernotifikasi-customer--get-apiv1kasirnotifikasi-kasir--get-apiv1kurirnotifikasi-kurir) - Mendapatkan daftar notifikasi untuk Customer
*   [**GET /api/v1/kasir/notifikasi**](#get-apiv1customernotifikasi-customer--get-apiv1kasirnotifikasi-kasir--get-apiv1kurirnotifikasi-kurir) - Mendapatkan daftar notifikasi untuk Kasir
*   [**GET /api/v1/kurir/notifikasi**](#get-apiv1customernotifikasi-customer--get-apiv1kasirnotifikasi-kasir--get-apiv1kurirnotifikasi-kurir) - Mendapatkan daftar notifikasi untuk Kurir
*   [**GET /api/v1/admin/notifikasi**](#get-apiv1adminnotifikasi) - Mendapatkan daftar notifikasi khusus Admin

---

## GET /api/v1/customer/notifikasi (Customer) / GET /api/v1/kasir/notifikasi (Kasir) / GET /api/v1/kurir/notifikasi (Kurir)

Mengambil daftar riwayat notifikasi yang dialamatkan kepada pengguna yang saat ini sedang login (sesuai peran/role masing-masing).

*   **Autentikasi:** Wajib (Sesuai Role masing-masing)
*   **Header Wajib:** `Authorization: Bearer <access_token>`

### Response (200 OK)
```json
{
  "status": "success",
  "data": [
    {
      "id_notifikasi": 1,
      "judul": "Pesanan Dikemas",
      "pesan": "Pesanan Anda #MNT-CUST-10 sedang dikemas dan dipersiapkan oleh toko.",
      "status": "unread",
      "created_at": "2026-06-05T10:00:00Z"
    }
  ]
}
```

---

## GET /api/v1/admin/notifikasi

Mengambil daftar notifikasi sistem khusus untuk role Administrator.

*   **Autentikasi:** Wajib (Role: `Admin`)
*   **Header Wajib:** `Authorization: Bearer <access_token>` (atau HttpOnly cookie)

### Response (200 OK)
Response sama dengan struktur di atas.

---

## 💡 Notes & Logika Bisnis

*   **Filter User:** Notifikasi secara otomatis difilter di server menggunakan `id_user` yang diekstrak dari JWT token pengguna. Pengguna tidak akan pernah bisa melihat notifikasi pengguna lain.
*   **Daftar Status:**
    - `unread`: Notifikasi baru yang belum dibaca.
    - `read`: Notifikasi yang sudah dibaca oleh user.
