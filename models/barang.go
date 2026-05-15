package models

func (Barang) TableName() string {
	return "barang"
}

type Barang struct {
	IdBarang     uint   `gorm:"primaryKey;column:id_barang" json:"id_barang"`
	NamaBarang   string `gorm:"column:nama_barang" json:"nama_barang"`
	GambarBarang string `gorm:"column:gambar_barang" json:"gambar_barang"`

	//Relasi
	DiskonId   uint `gorm:"column:id_diskon;unique" json:"id_diskon"`
	SatuanId   uint `gorm:"column:id_satuan;unique" json:"id_satuan"`
	KategoriId uint `gorm:"column:id_kategori;unique" json:"id_kategori"`

	// Relasi
	Diskon   Diskon   `gorm:"foreignKey:DiskonId;references:IdDiskon" json:"diskon"`
	Satuan   Satuan   `gorm:"foreignKey:SatuanId;references:IdSatuan" json:"satuan"`
	Kategori Kategori `gorm:"foreignKey:KategoriId;references:IdKategori" json:"kategori"`
}
