-- =============================================
-- MIGRATION: Normalisasi Shift Kasir
-- =============================================

BEGIN;

CREATE TABLE IF NOT EXISTS shift_kasir (
    id SERIAL PRIMARY KEY,
    nama_shift VARCHAR(50) UNIQUE NOT NULL
);

INSERT INTO shift_kasir (nama_shift) VALUES
    ('Pagi'),
    ('Siang'),
    ('Malam')
ON CONFLICT (nama_shift) DO NOTHING;

ALTER TABLE kasir
    ADD COLUMN IF NOT EXISTS id_shift_kasir INTEGER;

UPDATE kasir k
SET id_shift_kasir = sk.id
FROM shift_kasir sk
WHERE k.shift = sk.nama_shift
  AND k.id_shift_kasir IS NULL;

UPDATE kasir
SET id_shift_kasir = (SELECT id FROM shift_kasir WHERE nama_shift = 'Pagi')
WHERE id_shift_kasir IS NULL;

ALTER TABLE kasir
    ALTER COLUMN id_shift_kasir SET NOT NULL;

COMMIT;
