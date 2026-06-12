# Transaksi — Pesanan API

---

## GET /api/v1/customer/pesanan (Customer) / GET /api/v1/kasir/pesanan (Kasir)

Auth. Mendapatkan daftar pesanan online/offline.

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id_pesanan": 1,
      "public_id": "9e3c8162-...",
      "total_pembayaran": 150000,
      "status_pesanan": "Diproses",
      "tipe_pesanan": "Online",
      "tanggal_pesanan": "2026-06-05T10:00:00Z"
    }
  ]
}
```

---

## POST /api/v1/customer/pesanan/checkout

Auth (Customer). Membuat pesanan online baru dari isi keranjang.

**Request:**
```json
{
  "id_alamat": "9e3c8162-...",
  "id_ekspedisi": 1,
  "id_layanan_ekspedisi": 1,
  "ongkos_kirim": 20000,
  "catatan": "Tolong dibungkus rapi",
  "id_metode_pembayaran": 2
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Pesanan berhasil dibuat",
  "data": {
    "id_pesanan": "9e3c8162-...",
    "total_bayar": 172000,
    "ongkos_kirim": 20000,
    "pajak": 15000,
    "midtrans_token": "3f0d242e-...",
    "redirect_url": "https://app.sandbox.midtrans.com/snap/v4/redirection/..."
  }
}
```

---

## GET /api/v1/customer/pesanan/:public_id (Customer) / GET /api/v1/kasir/pesanan/:public_id (Kasir) / GET /api/v1/kurir/pesanan/:public_id (Kurir)

Auth. Mendapatkan detail pesanan lengkap beserta info pengiriman dan rincian item.

**Response:**
```json
{
  "status": "success",
  "message": "Detail pesanan berhasil diambil",
  "data": {
    "no_pesanan": "9e3c8162-...",
    "status": "Diproses",
    "tanggal_pesan": "2026-06-05T10:00:00Z",
    "items": [
      {
        "id_barang": 1,
        "nama_barang": "Produk A",
        "varian": "Ukuran: L",
        "jumlah": 2,
        "harga_satuan": 50000,
        "gambar": "http://minio:9000/products/item.png"
      }
    ],
    "tujuan_pengantaran": {
      "nama_penerima": "Surya",
      "alamat_lengkap": "Semarang"
    },
    "kurir": {
      "nama_kurir": "Ricardo Holahilo",
      "plat_nomor": "",
      "ekspedisi": "Internal Toko",
      "foto_kurir": ""
    },
    "rincian_pembayaran": {
      "subtotal_items": 100000,
      "ongkir": 20000,
      "biaya_proteksi": 0,
      "total": 120000
    }
  }
}
```

---

## PATCH /api/v1/customer/pesanan/:public_id/batal

Auth (Customer). Membatalkan pesanan.

**Response:**
```json
{
  "status": "success",
  "message": "Pesanan berhasil dibatalkan"
}
```

---

## GET /api/v1/customer/pesanan/:public_id/lacak

Auth (Customer). Lihat status pengiriman. Detail skema kembalian dapat dilihat di `pengantaran.md`.
