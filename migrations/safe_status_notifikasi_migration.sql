-- =============================================
-- MIGRATION: Normalisasi Status Notifikasi
-- =============================================

BEGIN;

CREATE TABLE IF NOT EXISTS status_notifikasi (
    id SERIAL PRIMARY KEY,
    nama_status VARCHAR(50) UNIQUE NOT NULL
);

INSERT INTO status_notifikasi (nama_status) VALUES
    ('unread'),
    ('read'),
    ('aktif')
ON CONFLICT (nama_status) DO NOTHING;

ALTER TABLE notifikasi
    ADD COLUMN IF NOT EXISTS id_status_notifikasi INTEGER;

UPDATE notifikasi n
SET id_status_notifikasi = sn.id
FROM status_notifikasi sn
WHERE n.status = sn.nama_status
  AND n.id_status_notifikasi IS NULL;

UPDATE notifikasi
SET id_status_notifikasi = (SELECT id FROM status_notifikasi WHERE nama_status = 'unread')
WHERE id_status_notifikasi IS NULL;

ALTER TABLE notifikasi
    ALTER COLUMN id_status_notifikasi SET NOT NULL;

COMMIT;
