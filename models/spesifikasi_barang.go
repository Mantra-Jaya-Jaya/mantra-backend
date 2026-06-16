package models

func (SpesifikasiBarang) TableName() string {
	return "spesifikasi_barang"
}

// Junction table: barang + varian + stok + harga per varian.
// Barang tanpa varian tetap wajib punya 1 baris dengan detail "Default"
type SpesifikasiBarang struct {
	IdSpesifikasiBarang uint `gorm:"primaryKey;column:id_spesifikasi_barang" json:"id_spesifikasi_barang"`
	Jumlah              int  `gorm:"column:jumlah" json:"jumlah"`
	HargaBarang         int  `gorm:"column:harga_barang" json:"harga_barang"`
	BeratBarang         int  `gorm:"column:berat_barang;default:0" json:"berat_barang"`

	// Relasi ke Barang
	BarangID uint   `gorm:"column:id_barang;not null" json:"id_barang"`
	Barang   Barang `gorm:"foreignKey:BarangID;references:IdBarang" json:"barang"`

	// Relasi ke DetailSpesifikasi
	DetailSpesifikasiID uint              `gorm:"column:id_detail_spesifikasi;not null" json:"id_detail_spesifikasi"`
	DetailSpesifikasi   DetailSpesifikasi `gorm:"foreignKey:DetailSpesifikasiID;references:IdDetailSpesifikasi" json:"detail_spesifikasi"`
}
