-- Migration: add indexes for ekspedisi sync performance
-- Fixes: slow INSERT/SELECT on ekspedisi_layanan (3s+) during Biteship sync
-- =====================================================================

-- Index untuk lookup courier by kode_api (WHERE kode_api = ?)
CREATE INDEX IF NOT EXISTS idx_ekspedisi_kode_api ON ekspedisi(kode_api);

-- Index untuk lookup layanan by id_ekspedisi (WHERE id_ekspedisi = ? AND nama_layanan = ?)
-- juga mempercepat INSERT karena FK check tidak lagi full-scan
CREATE INDEX IF NOT EXISTS idx_ekspedisi_layanan_id_ekspedisi ON ekspedisi_layanan(id_ekspedisi);
