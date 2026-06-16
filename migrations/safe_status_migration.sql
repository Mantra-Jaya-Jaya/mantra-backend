-- ==========================================
-- SAFE MIGRATION: Normalisasi StatusPesanan
-- Tidak destructive, preserve all data
-- ==========================================

-- STEP 1: Create tabel status_pesanan
CREATE TABLE IF NOT EXISTS status_pesanan (
    id SERIAL PRIMARY KEY,
    nama_status VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- STEP 2: Insert master data status
INSERT INTO status_pesanan (nama_status) VALUES
    ('Draft'),
    ('Menunggu Pembayaran'),
    ('Diproses'),
    ('Dikemas'),
    ('Dikirim'),
    ('Selesai'),
    ('Dibatalkan')
ON CONFLICT (nama_status) DO NOTHING;

-- STEP 3: Add new column (nullable first)
ALTER TABLE pesanan ADD COLUMN IF NOT EXISTS id_status_pesanan INTEGER;

-- STEP 4: Migrate existing data
UPDATE pesanan p
SET id_status_pesanan = (
    SELECT id FROM status_pesanan
    WHERE nama_status = p.status_pesanan
)
WHERE id_status_pesanan IS NULL;

-- STEP 5: Verify migration
SELECT
    COUNT(*) as total_records,
    COUNT(id_status_pesanan) as migrated_records,
    COUNT(*) - COUNT(id_status_pesanan) as unmigrated_records
FROM pesanan;

-- STEP 6: Check unmigrated records (should be 0)
SELECT DISTINCT status_pesanan, id_status_pesanan
FROM pesanan
WHERE id_status_pesanan IS NULL
LIMIT 10;

-- STEP 7: Make column NOT NULL (only if all migrated)
-- Uncomment after verification:
-- ALTER TABLE pesanan ALTER COLUMN id_status_pesanan SET NOT NULL;

-- STEP 8: Add foreign key constraint
ALTER TABLE pesanan
ADD CONSTRAINT fk_pesanan_status
FOREIGN KEY (id_status_pesanan)
REFERENCES status_pesanan(id)
ON DELETE RESTRICT;

-- STEP 9: Create index for performance
CREATE INDEX IF NOT EXISTS idx_pesanan_status_id
ON pesanan(id_status_pesanan);

-- STEP 10: Verify final result
SELECT
    p.id_status_pesanan,
    sp.nama_status,
    COUNT(*) as count
FROM pesanan p
JOIN status_pesanan sp ON p.id_status_pesanan = sp.id
GROUP BY p.id_status_pesanan, sp.nama_status
ORDER BY p.id_status_pesanan;
