# Keranjang API

Auth (Customer). Semua endpoint keranjang untuk customer.

---

## POST /api/v1/customer/keranjang

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

## PATCH /api/v1/customer/keranjang/:public_id

Update quantity item di keranjang.

**Request:**

```json
{
  "quantity": 5
}
```

**Response:**

```json
{
  "status": "success",
  "message": "Jumlah item berhasil diperbarui"
}
```

---

## DELETE /api/v1/customer/keranjang/:public_id

Hapus item dari keranjang.

**Response:**

```json
{
  "status": "success",
  "message": "Item berhasil dihapus dari keranjang"
}
```
