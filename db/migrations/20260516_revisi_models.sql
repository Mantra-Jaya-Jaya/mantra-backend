-- Perubahan 1: Relasi barcode → spesifikasi_barang
ALTER TABLE barcode DROP CONSTRAINT IF EXISTS barcode_id_barang_fkey;
ALTER TABLE barcode DROP COLUMN IF EXISTS id_barang;
ALTER TABLE barcode ADD COLUMN id_spesifikasi_barang INT NOT NULL;
ALTER TABLE barcode
    ADD CONSTRAINT barcode_id_spesifikasi_barang_fkey
    FOREIGN KEY (id_spesifikasi_barang)
    REFERENCES spesifikasi_barang(id_spesifikasi_barang);

-- Perubahan 2: Tambah foto_profil ke tabel user
ALTER TABLE "user" ADD COLUMN IF NOT EXISTS foto_profil VARCHAR;
