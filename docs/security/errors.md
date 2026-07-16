# Error Codes & HTTP Status

## Error Response Format

Semua error response mengikuti format seragam:

```json
{
  "status": "error",
  "message": "Pesan untuk user",
  "error": {
    "code": "XXX_000",
    "detail": "Penjelasan teknis untuk debugging"
  }
}
```

## Error Code Table

| Kode | HTTP | Sumber | Arti |
|------|------|--------|------|
| **Authentication & Authorization** | | | |
| `AUTH_001` | 401 | AuthMiddleware | Token tidak valid / expired |
| `AUTH_002` | 403 | RoleMiddleware | Role tidak ditemukan di token |
| `AUTH_003` | 403 | RoleMiddleware | Tidak punya izin akses resource |
| `AUTH_004` | 403 | OwnershipMiddleware | Bukan pemilik resource |
| `AUTH_005` | 401 | AuthMiddleware | Klien lupa mengirim token di Header |
| `AUTH_006` | 403 | Controller | Akun dinonaktifkan atau di-banned |
| `AUTH_007` | 401 | Controller | Email atau Password salah saat login |
| **Validation & Client Request** | | | |
| `VAL_001` | 400 / 422 | Controller | Input tidak memenuhi aturan validasi |
| `VAL_002` | 422 | Controller | Konfirmasi password tidak cocok |
| `VAL_003` | 422 | Controller | Format email tidak valid |
| `REQ_001` | 400 | Controller | Format JSON/Data hancur (Bad Request) |
| `REQ_003` | 400 | Controller | Password lama tidak sesuai |
| `REQ_004` | 404 | Controller | Data tidak ditemukan (Not Found) |
| `REQ_005` | 404 / 405 | Router | URL Endpoint salah ketik / Method salah |
| **Conflict & Database** | | | |
| `CONF_001` | 409 | Controller | Username sudah terdaftar |
| `CONF_002` | 409 | Controller | Email sudah terdaftar |
| `CONF_003` | 409 | Controller | Data tidak bisa dihapus (digunakan tabel lain) |
| **File Upload** | | | |
| `FILE_001` | 400 | Controller | Ukuran gambar/file melebihi batas (Too Large) |
| `FILE_002` | 400 | Controller | Format file salah (Misal unggah PDF padahal diminta PNG) |
| `FILE_003` | 400 | Controller | File corrupt / tidak bisa dibaca sistem |
| **E-Commerce & Transaksi** | | | |
| `TRX_001` | 400 | Controller | Checkout gagal karena stok barang habis |
| `TRX_002` | 400 | Controller | Mencoba checkout tapi keranjang kosong |
| `TRX_003` | 400 | Controller | Transaksi/Pembayaran gagal diproses gateway |
| `TRX_004` | 400 | Controller | Voucher / Diskon sudah kedaluwarsa |
| `TRX_005` | 400 | Controller | Kuota penggunaan voucher sudah habis |
| `TRX_006` | 400 | Controller | Total belanja belum mencapai minimum voucher |
| **Server & Infrastructure** | | | |
| `SERVER_001` | 500 | Any | Internal server error (Panic / Crash) |
| `SERVER_002` | 500 | Any | Gagal terhubung ke Database (Timeout) |
| `SERVER_003` | 500 | Any | Storage (MinIO) mati atau penuh |

## HTTP Status Codes

| Status | Usage |
|--------|-------|
| 200 | OK — Request berhasil |
| 201 | Created — Resource berhasil dibuat |
| 400 | Bad Request — Input tidak valid / Error Logika E-Commerce |
| 401 | Unauthorized — Token invalid / expired |
| 403 | Forbidden — Tidak punya akses |
| 404 | Not Found — Resource tidak ditemukan |
| 409 | Conflict — Duplikasi data / Constraint |
| 422 | Unprocessable Entity — Validasi form gagal |
| 500 | Internal Server Error — Error sistem |
