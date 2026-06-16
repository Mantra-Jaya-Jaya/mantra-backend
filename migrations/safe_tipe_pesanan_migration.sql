-- =============================================
-- MIGRATION: Normalisasi Tipe Pesanan
-- =============================================

BEGIN;

CREATE TABLE IF NOT EXISTS tipe_pesanan (
    id SERIAL PRIMARY KEY,
    nama_tipe VARCHAR(50) UNIQUE NOT NULL
);

INSERT INTO tipe_pesanan (nama_tipe) VALUES
    ('Online'),
    ('Offline')
ON CONFLICT (nama_tipe) DO NOTHING;

ALTER TABLE pesanan
    ADD COLUMN IF NOT EXISTS id_tipe_pesanan INTEGER;

UPDATE pesanan p
SET id_tipe_pesanan = tp.id
FROM tipe_pesanan tp
WHERE p.tipe_pesanan = tp.nama_tipe
  AND p.id_tipe_pesanan IS NULL;

UPDATE pesanan
SET id_tipe_pesanan = (SELECT id FROM tipe_pesanan WHERE nama_tipe = 'Offline')
WHERE id_tipe_pesanan IS NULL;

ALTER TABLE pesanan
    ALTER COLUMN id_tipe_pesanan SET NOT NULL;

COMMIT;
