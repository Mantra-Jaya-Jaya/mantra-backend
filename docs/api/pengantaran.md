# 🛵 Pengantaran (Courier & Shipping) API Contract

---
### 🧭 Navigasi Cepat
[🏠 Utama](../README.md) | [🏛️ Arsitektur](../architecture.md) | [🛠️ Deployment](../deployment.md) | [💳 Midtrans](../pembayaran.md) | [📦 Biteship](../biteship.md) | [📡 API Contract](overview.md) | [🗄️ Database](../database/erd.md) | [🔒 Keamanan](../security/README.md)
---

API untuk melacak pesanan oleh customer serta manajemen tugas pengantaran oleh kurir toko (internal).

---

## 🧭 Daftar Endpoint Pengantaran

*   [**GET /api/v1/customer/pesanan/:public_id/lacak**](#get-apiv1customerpesananpublic_idlacak) - Melacak status pengantaran paket (Customer)
*   [**PATCH /api/v1/kurir/pengantaran/:public_id/lokasi**](#patch-apiv1kurirpengantaranpublic_idlokasi) - Memperbarui koordinat GPS live kurir (Kurir)
*   [**POST /api/v1/kurir/pengantaran/:public_id/ambil**](#post-apiv1kurirpengantaranpublic_idambil) - Mengambil (claim) tugas pengiriman barang (Kurir)
*   [**POST /api/v1/kurir/pengantaran/:public_id/status**](#post-apiv1kurirpengantaranpublic_idstatus) - Memperbarui status pengantaran & upload bukti foto (Kurir)

---

## GET /api/v1/customer/pesanan/:public_id/lacak

Melacak riwayat status dan detail pengiriman suatu pesanan (Customer).

*   **Autentikasi:** Wajib (Role: `Customer`)
*   **Header Wajib:** `Authorization: Bearer <access_token>`
*   **Parameter URL:** `public_id` (UUID Pesanan)

### Response (Jika Ekspedisi Eksternal - Biteship)
```json
{
  "status": "success",
  "message": "Data lacak pesanan berhasil diambil",
  "data": {
    "id_pesanan": "9e3c8162-...",
    "nomor_resi": "JT1234567890",
    "ekspedisi": "J&T",
    "tipe_ekspedisi": "eksternal",
    "history": [
      {
        "time": "2026-06-12T10:00:00Z",
        "description": "Paket telah diterima oleh J&T Gateway",
        "status": "picked_up",
        "city": "Semarang"
      }
    ]
  }
}
```

### Response (Jika Ekspedisi Internal - Kurir Toko)
```json
{
  "status": "success",
  "message": "Data lacak pesanan berhasil diambil",
  "data": {
    "id_pesanan": "9e3c8162-...",
    "tipe_ekspedisi": "internal",
    "kurir": {
      "nama": "Ricardo Holahilo",
      "plat_nomor": "H 1234 AB",
      "foto": "http://minio:9000/profiles/avatar.png"
    },
    "lokasi_kurir": {
      "latitude": -7.02,
      "longitude": 110.43
    },
    "estimasi_tiba": "15 menit",
    "jarak_meter": 2300
  }
}
```

---

## PATCH /api/v1/kurir/pengantaran/:public_id/lokasi

Memperbarui koordinat GPS live kurir toko secara periodik selama di perjalanan.

*   **Autentikasi:** Wajib (Role: `Kurir`)
*   **Header Wajib:** `Authorization: Bearer <access_token>`
*   **Parameter URL:** `public_id` (UUID Pengantaran)

### Request Payload
```json
{
  "latitude": -7.02561,
  "longitude": 110.43128
}
```

### Response (200 OK)
```json
{
  "status": "success",
  "message": "Lokasi kurir berhasil diperbarui"
}
```

---

## POST /api/v1/kurir/pengantaran/:public_id/ambil

Mengambil tugas pengantaran pesanan oleh kurir toko untuk pesanan yang siap dikirim (`Dikemas`).

*   **Autentikasi:** Wajib (Role: `Kurir`)
*   **Header Wajib:** `Authorization: Bearer <access_token>`
*   **Parameter URL:** `public_id` (UUID Pesanan)

### Response (200 OK)
```json
{
  "status": "success",
  "message": "Pesanan berhasil diambil",
  "data": {
    "id_pengantaran": "8f3c8162-...",
    "status": "Dalam Perjalanan"
  }
}
```

---

## POST /api/v1/kurir/pengantaran/:public_id/status

Memperbarui status tugas pengantaran (misal dari "Dalam Perjalanan" menjadi "Selesai"). 

*   **Autentikasi:** Wajib (Role: `Kurir`)
*   **Header Wajib:** `Authorization: Bearer <access_token>`
*   **Parameter URL:** `public_id` (UUID Pengantaran)
*   **Tipe Request:** `multipart/form-data` (Jika status diubah ke **Selesai**, wajib menyertakan foto bukti pengantaran).

### Request Payload (Multipart Form Data)
- `status`: `Selesai` (Text)
- `foto`: `[Binary Image File]` (File upload gambar bukti penerimaan)

### Response (200 OK)
```json
{
  "status": "success",
  "message": "Status pengantaran berhasil diperbarui",
  "data": {
    "foto_bukti": "http://minio:9000/pengantaran/bukti_123.jpg",
    "status": "Selesai"
  }
}
```
