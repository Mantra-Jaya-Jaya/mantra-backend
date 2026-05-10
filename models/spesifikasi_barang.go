package models

func (SpesifikasiBarang) TableName() string {
	return "spesifikasi_barang"
}

// Junction table: barang + varian + stok + harga per varian.
// Barang tanpa varian tetap wajib punya 1 baris dengan detail "Default"
type SpesifikasiBarang struct {
	IdSpesifikasiBarang uint `gorm:"primaryKey;column:id_spesifikasi_barang"`
	Jumlah              int  `gorm:"column:jumlah"`       // Stok per varian
	HargaBarang         int  `gorm:"column:harga_barang"` // Harga per varian dalam Rupiah (integer)

	// Relasi ke Barang
	BarangID uint   `gorm:"column:id_barang"`
	Barang   Barang `gorm:"foreignKey:BarangID;references:IdBarang"`

	// Relasi ke DetailSpesifikasi (nilai varian: Hitam, XL, 16GB, dll)
	DetailSpesifikasiID uint              `gorm:"column:id_detail_spesifikasi"`
	DetailSpesifikasi   DetailSpesifikasi `gorm:"foreignKey:DetailSpesifikasiID;references:IdDetailSpesifikasi"`
}
