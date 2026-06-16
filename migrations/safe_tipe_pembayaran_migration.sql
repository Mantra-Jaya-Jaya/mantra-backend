-- =============================================
-- MIGRATION: Normalisasi Tipe Pembayaran (RawPaymentType)
-- =============================================

BEGIN;

CREATE TABLE IF NOT EXISTS tipe_pembayaran (
    id SERIAL PRIMARY KEY,
    nama_tipe VARCHAR(50) UNIQUE NOT NULL
);

INSERT INTO tipe_pembayaran (nama_tipe) VALUES
    ('cash'),
    ('non-cash'),
    ('qris'),
    ('bank_transfer'),
    ('gopay')
ON CONFLICT (nama_tipe) DO NOTHING;

ALTER TABLE pembayaran
    ADD COLUMN IF NOT EXISTS id_tipe_pembayaran INTEGER;

UPDATE pembayaran p
SET id_tipe_pembayaran = tp.id
FROM tipe_pembayaran tp
WHERE p.payment_type = tp.nama_tipe
  AND p.id_tipe_pembayaran IS NULL;

UPDATE pembayaran
SET id_tipe_pembayaran = (SELECT id FROM tipe_pembayaran WHERE nama_tipe = 'non-cash')
WHERE id_tipe_pembayaran IS NULL;

ALTER TABLE pembayaran
    ALTER COLUMN id_tipe_pembayaran SET NOT NULL;

COMMIT;
