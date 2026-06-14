-- =============================================
-- MIGRATION: Normalisasi StatusTransaksi
-- 
-- Tujuan: Convert kolom status_transaksi (string)
-- menjadi id_status_transaksi (uint) dengan lookup
-- table status_transaksi.
-- =============================================

BEGIN;

-- 1. Buat lookup table
CREATE TABLE IF NOT EXISTS status_transaksi (
    id SERIAL PRIMARY KEY,
    nama_status VARCHAR(50) UNIQUE NOT NULL
);

-- 2. Insert data status transaksi
INSERT INTO status_transaksi (nama_status) VALUES
    ('pending'),
    ('settlement'),
    ('deny'),
    ('cancel'),
    ('expire')
ON CONFLICT (nama_status) DO NOTHING;

-- 3. Tambah kolom baru di tabel pembayaran
ALTER TABLE pembayaran
    ADD COLUMN IF NOT EXISTS id_status_transaksi INTEGER;

-- 4. Update kolom baru dengan ID yang sesuai
UPDATE pembayaran p
SET id_status_transaksi = st.id
FROM status_transaksi st
WHERE p.status_transaksi = st.nama_status
  AND p.id_status_transaksi IS NULL;

-- 5. Set default untuk data yang NULL (anggap 'pending')
UPDATE pembayaran
SET id_status_transaksi = (SELECT id FROM status_transaksi WHERE nama_status = 'pending')
WHERE id_status_transaksi IS NULL;

-- 6. Set NOT NULL setelah data terisi
ALTER TABLE pembayaran
    ALTER COLUMN id_status_transaksi SET NOT NULL;

COMMIT;
