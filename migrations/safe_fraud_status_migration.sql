-- =============================================
-- MIGRATION: Normalisasi FraudStatus
-- =============================================

BEGIN;

-- 1. Buat lookup table
CREATE TABLE IF NOT EXISTS fraud_status (
    id SERIAL PRIMARY KEY,
    nama_status VARCHAR(50) UNIQUE NOT NULL
);

-- 2. Insert data fraud status
INSERT INTO fraud_status (nama_status) VALUES
    ('accept'),
    ('challenge'),
    ('deny')
ON CONFLICT (nama_status) DO NOTHING;

-- 3. Tambah kolom baru di tabel pembayaran
ALTER TABLE pembayaran
    ADD COLUMN IF NOT EXISTS id_fraud_status INTEGER;

-- 4. Update kolom baru dengan ID yang sesuai
UPDATE pembayaran p
SET id_fraud_status = fs.id
FROM fraud_status fs
WHERE p.fraud_status = fs.nama_status
  AND p.id_fraud_status IS NULL;

-- 5. Set default untuk data yang NULL (anggap 'accept')
UPDATE pembayaran
SET id_fraud_status = (SELECT id FROM fraud_status WHERE nama_status = 'accept')
WHERE id_fraud_status IS NULL;

-- 6. Set NOT NULL setelah data terisi
ALTER TABLE pembayaran
    ALTER COLUMN id_fraud_status SET NOT NULL;

-- 7. Tambah foreign key constraint
ALTER TABLE pembayaran
    ADD CONSTRAINT fk_pembayaran_fraud_status
    FOREIGN KEY (id_fraud_status)
    REFERENCES fraud_status(id)
    ON DELETE RESTRICT;

-- 8. Tambah index untuk performansi query
CREATE INDEX IF NOT EXISTS idx_pembayaran_fraud_status
    ON pembayaran(id_fraud_status);

COMMIT;