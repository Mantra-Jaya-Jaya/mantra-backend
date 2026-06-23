# 📡 API Contract Overview

---
### 🧭 Navigasi Cepat
[🏠 Utama](../README.md) | [🏛️ Arsitektur](../architecture.md) | [🛠️ Deployment](../deployment.md) | [💳 Midtrans](../pembayaran.md) | [📦 Biteship](../biteship.md) | [📡 API Contract](overview.md) | [🗄️ Database](../database/erd.md) | [🔒 Keamanan](../security/README.md)
---

Halaman ini berisi base URL, format komunikasi JSON, daftar status HTTP, serta format penulisan response sukses dan error yang berlaku secara global untuk seluruh endpoint API **MANTRA**.

---

## 🌐 1. Base URL

| Environment | URL | Deskripsi |
| :--- | :--- | :--- |
| **Local Development** | `http://localhost:8080/api/v1` | Server lokal |
| **Production Server** | *(akan ditambahkan setelah deployment)* | Server live |

---

## 🔒 2. Authentication & Authorization

Autentikasi diatur secara otomatis tergantung pada tipe client yang menembak API (ditentukan via header `X-Client-Type`):

### A. Aplikasi Mobile (Flutter - Customer & Kasir)
*   **Metode:** Bearer Token.
*   **Header:** `Authorization: Bearer <access_token>`

### B. Dashboard Web Admin (Next.js)
*   **Metode:** HttpOnly Cookies.
*   **Header Cookie:** `access_token=<token>; refresh_token=<token>`

---

## 📥 3. Format Response Sukses (Success Response)

Seluruh response sukses memiliki struktur JSON seragam sebagai berikut:

```json
{
  "status": "success",
  "message": "Pesan sukses untuk konfirmasi tindakan",
  "data": {
    "key": "value"
  }
}
```

---

## ❌ 4. Format Response Error (Error Response)

Seluruh response error (HTTP 4xx & 5xx) memiliki struktur JSON seragam berikut:

```json
{
  "status": "error",
  "message": "Deskripsi error yang ramah untuk ditampilkan ke user",
  "error": {
    "code": "ERR_XXX",
    "detail": "Penjelasan detail teknis/sistem untuk keperluan debugging"
  }
}
```

### Daftar Standar Kode Error (Error Codes)

| Kode Error | Arti / Deskripsi |
| :--- | :--- |
| **`VAL_001`** | Validasi gagal — Format input payload tidak sesuai |
| **`VAL_002`** | Validasi gagal — Konfirmasi password atau input tidak cocok |
| **`VAL_003`** | Validasi gagal — Format field input tertentu tidak valid |
| **`AUTH_001`** | Token autentikasi tidak valid atau sudah kedaluwarsa (expired) |
| **`AUTH_002`** | Klaim Role tidak ditemukan di dalam token |
| **`AUTH_003`** | Tidak memiliki izin (permission) untuk mengakses resource |
| **`AUTH_004`** | Bukan pemilik (owner) dari resource yang diminta |
| **`CONF_001`** | Konflik data — Username sudah terdaftar di sistem |
| **`CONF_002`** | Konflik data — Email sudah terdaftar di sistem |
| **`REQ_003`** | Bad Request — Password lama tidak sesuai |
| **`REQ_004`** | Not Found — Data/Resource tidak ditemukan di database |
| **`SERVER_001`** | Internal Server Error — Kesalahan pada pemrosesan server |

---

## 📶 5. HTTP Status Codes

Berikut adalah daftar HTTP Status Codes yang digunakan pada API MANTRA:

| Status Code | Nama | Keterangan |
| :--- | :--- | :--- |
| **200** | OK | Request sukses diproses |
| **201** | Created | Sukses membuat data/resource baru |
| **400** | Bad Request | Request tidak valid (format payload salah) |
| **401** | Unauthorized | Belum melakukan login / token tidak valid |
| **403** | Forbidden | Token valid, tapi tidak memiliki hak akses role |
| **404** | Not Found | Resource/Data tidak ditemukan |
| **409** | Conflict | Terjadi bentrokan data (unique constraint violation) |
| **422** | Unprocessable Entity | Payload ter-parse tapi gagal validasi logika/bisnis |
| **500** | Internal Server Error | Kesalahan sistem internal pada server |

---

## 🆔 6. Format ID (Primary Key)

Untuk menjaga keamanan data dan mencegah eksploitasi peretasan (ID Enumeration / Insecure Direct Object Reference):
- **Internal ID (`id_xxx`):** Integer auto-increment bertipe data serial. Hanya boleh digunakan di internal database (sebagai Foreign Key) dan tidak boleh diekspos di URL API.
- **External ID (`public_id`):** Menggunakan tipe data **UUID v4**. Seluruh parameter URL API eksternal WAJIB menggunakan `public_id` untuk mencari data.
