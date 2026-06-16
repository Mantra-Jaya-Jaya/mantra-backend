-- =============================================
-- MIGRATION: Normalisasi Status Karyawan
-- =============================================

BEGIN;

CREATE TABLE IF NOT EXISTS status_karyawan (
    id SERIAL PRIMARY KEY,
    nama_status VARCHAR(50) UNIQUE NOT NULL
);

INSERT INTO status_karyawan (nama_status) VALUES
    ('Aktif'),
    ('Nonaktif'),
    ('Probation')
ON CONFLICT (nama_status) DO NOTHING;

ALTER TABLE karyawan
    ADD COLUMN IF NOT EXISTS id_status_karyawan INTEGER;

UPDATE karyawan k
SET id_status_karyawan = sk.id
FROM status_karyawan sk
WHERE k.status = sk.nama_status
  AND k.id_status_karyawan IS NULL;

UPDATE karyawan
SET id_status_karyawan = (SELECT id FROM status_karyawan WHERE nama_status = 'Aktif')
WHERE id_status_karyawan IS NULL;

ALTER TABLE karyawan
    ALTER COLUMN id_status_karyawan SET NOT NULL;

COMMIT;
