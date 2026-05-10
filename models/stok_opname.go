package models

import "time"

func (StokOpname) TableName() string {
	return "stok_opname"
}

// Riwayat pergerakan stok per varian barang
type StokOpname struct {
	IdStokOpname uint      `gorm:"primaryKey;column:id_stok_opname"`
	HargaBeli    int       `gorm:"column:harga_beli"` // Harga beli/modal dalam Rupiah
	HargaJual    int       `gorm:"column:harga_jual"` // Harga jual dalam Rupiah
	Status       bool      `gorm:"column:status"`     // true = stok masuk, false = stok keluar
	JumlahStok   int       `gorm:"column:jumlah_stok"`
	Keterangan   string    `gorm:"type:text;column:keterangan"`
	Tanggal      time.Time `gorm:"column:tanggal"`

	// Relasi ke SpesifikasiBarang
	SpesifikasiBarangID uint              `gorm:"column:id_spesifikasi_barang"`
	SpesifikasiBarang   SpesifikasiBarang `gorm:"foreignKey:SpesifikasiBarangID;references:IdSpesifikasiBarang"`
}
