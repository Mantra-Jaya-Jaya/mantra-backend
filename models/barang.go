package models

func (Barang) TableName() string {
	return "barang"
}

type Barang struct {
	IdBarang     uint   `gorm:"primaryKey;column:id_barang"`
	NamaBarang   string `gorm:"column:nama_barang"`
	GambarBarang string `gorm:"column:gambar_barang"`

	// Relasi (DiskonId nullable — barang bisa tidak punya diskon)
	DiskonId   *uint `gorm:"column:id_diskon"`
	SatuanId   uint  `gorm:"column:id_satuan"`
	KategoriId uint  `gorm:"column:id_kategori"`

	// Relasi
	Diskon   *Diskon  `gorm:"foreignKey:DiskonId;references:IdDiskon"`
	Satuan   Satuan   `gorm:"foreignKey:SatuanId;references:IdSatuan"`
	Kategori Kategori `gorm:"foreignKey:KategoriId;references:IdKategori"`
}
