# Biteship Manual Check — Saat Ini

Status: **belum commit**, hanya catatan.

## Hasil Test Saat Ini

- Login Customer: berhasil
- Address customer: `4747116d-7f86-4fdb-b500-0126cb4b8cfd`
- Produk test: `id_barang=1`, `id_spesifikasi_barang=1`
- Ongkir external: berhasil (JNE, AnterAja, dst.)
- Checkout external (tunai): pesanan `382bd4e5-032d-403b-b4d1-9ca1ca6bc0a8`
- Status pesanan: `Dikemas`
- `/status-biteship`: return `pending` (belum ada `biteship_order_id`)

## Yang Perlu Dicek Lanjut

- [ ] Checkout non-tunai → Midtrans settlement → `processExternalShipment()` jalan
- [ ] Tracking `/lacak` untuk pesanan eksternal
- [ ] Webhook Biteship (`/webhook/biteship`) dengan signature header
- [ ] Cancel shipment (`/cancel-shipment`) untuk pesanan yang belum `Dikirim`
- [ ] Build `go build ./...` tanpa error di branch `feat/customer-biteship-tracking`

## Catatan

- `mantra-mobile` dan `mantra-admin-web`: **tanpa perubahan**
- `mantra-backend`: perubahan ada di branch `feat/customer-biteship-tracking`
- Migration baru: `add_biteship_order_id_to_pesanan.sql` belum dijalankan ke DB produksi
