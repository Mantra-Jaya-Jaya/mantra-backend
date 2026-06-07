# Katalog — Barang API

---

## GET /barang

Auth (Customer). Mendapatkan daftar barang untuk customer.

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id_barang": 1,
      "public_id": "uuid-...",
      "nama_barang": "Produk A",
      "gambar_barang": "https://storage...",
      "harga": 50000,
      "diskon": 10
    }
  ]
}
```

---

## GET /admin/barang

Auth (Admin). Mendapatkan daftar semua barang.

---

## POST /admin/barang

Auth (Admin). Menambah barang baru.

**Request:**
```json
{
  "nama_barang": "Produk Baru",
  "id_kategori": 1,
  "id_satuan": 1,
  "id_diskon": 1,
  "deskripsi": "Deskripsi produk"
}
```

---

## GET /admin/barang/detail/:public_id

Auth (Admin). Mendapatkan detail barang berdasarkan public_id UUID.

---

## PUT /admin/barang/:public_id

Auth (Admin). Mengupdate data barang.

---

## DELETE /admin/barang/:public_id

Auth (Admin). Menghapus barang.

---

## POST /admin/barang/upload

Auth (Admin). Upload gambar barang ke MinIO.

**Request:** `multipart/form-data` — field `gambar`

---

## GET /scan/:kode_barcode

Public. Mendapatkan detail barang berdasarkan kode barcode.

**Response:**
```json
{
  "status": "success",
  "data": {
    "nama_barang": "Produk A",
    "harga": 50000,
    "stok": 100
  }
}
```

---

## POST /kasir/transaksi/produk

Auth (Kasir). Mencari produk untuk transaksi POS.

**Request:**
```json
{
  "keyword": "nama atau barcode"
}
```
