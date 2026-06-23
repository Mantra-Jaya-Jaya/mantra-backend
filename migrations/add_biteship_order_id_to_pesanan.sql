-- ==========================================
-- MIGRATION: Add Biteship Order ID to Pesanan
-- Untuk menyimpan order_id Biteship (bukan waybill)
-- ==========================================

-- Tambah kolom untuk menyimpan Biteship Order ID
ALTER TABLE pesanan
ADD COLUMN IF NOT EXISTS biteship_order_id VARCHAR(100);

-- Index untuk performa query
CREATE INDEX IF NOT EXISTS idx_pesanan_biteship_order_id
ON pesanan(biteship_order_id);