-- =============================================
-- MIGRATION: Normalisasi FraudStatus
-- =============================================

BEGIN;

CREATE TABLE IF NOT EXISTS fraud_status (
    id SERIAL PRIMARY KEY,
    nama_status VARCHAR(50) UNIQUE NOT NULL
);

INSERT INTO fraud_status (nama_status) VALUES
    ('accept'),
    ('challenge'),
    ('deny')
ON CONFLICT (nama_status) DO NOTHING;

ALTER TABLE pembayaran
    ADD COLUMN IF NOT EXISTS id_fraud_status INTEGER;

UPDATE pembayaran p
SET id_fraud_status = fs.id
FROM fraud_status fs
WHERE p.fraud_status = fs.nama_status
  AND p.id_fraud_status IS NULL;

UPDATE pembayaran
SET id_fraud_status = (SELECT id FROM fraud_status WHERE nama_status = 'accept')
WHERE id_fraud_status IS NULL;

ALTER TABLE pembayaran
    ALTER COLUMN id_fraud_status SET NOT NULL;

COMMIT;
