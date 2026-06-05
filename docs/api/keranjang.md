# Keranjang API

Auth (Customer). Semua endpoint keranjang untuk customer.

---

## POST /customer/keranjang

Tambah item ke keranjang.

**Request:**
```json
{
  "id_spesifikasi_barang": 1,
  "quantity": 2
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Item berhasil ditambahkan ke keranjang"
}
```

---

## PATCH /customer/keranjang/:id_keranjang

Update quantity item di keranjang.

**Request:**
```json
{
  "quantity": 5
}
```

---

## DELETE /customer/keranjang/:id_keranjang

Hapus item dari keranjang.

**Response:**
```json
{
  "status": "success",
  "message": "Item berhasil dihapus dari keranjang"
}
```
