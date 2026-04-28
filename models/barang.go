package models

func (Barang) TableName() string {
	return "barang"
}

type Barang struct {
	IdBarang     uint   `gorm:"primaryKey;column:id_barang"`
	NamaBarang   string `gorm:"column:nama_barang"`
	HargaBarang  uint   `gorm:"column:harga_barang"`
	StokBarang   int    `gorm:"column:stok_barang"`
	GambarBarang string `gorm:"column:gambar_barang"`

	//Relasi
	DiskonId   uint `gorm:"column:id_diskon;unique"`
	SatuanId   uint `gorm:"column:id_satuan;unique"`
	KategoriId uint `gorm:"column:id_kategori;unique"`

	// Relasi
	Diskon   Diskon   `gorm:"foreignKey:DiskonId;references:IdDiskon"`
	Satuan   Satuan   `gorm:"foreignKey:SatuanId;references:IdSatuan"`
	Kategori Kategori `gorm:"foreignKey:KategoriId;references:IdKategori"`
}
