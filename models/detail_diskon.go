package models

func (DetailDiskon) TableName() string {
	return "detail_diskon"
}

type DetailDiskon struct {
	IdDetailDiskon   uint     `gorm:"primaryKey;column:id_detail_diskon"`
	
	DiskonID         uint     `gorm:"column:id_diskon"`
	Diskon           Diskon   `gorm:"foreignKey:DiskonID;references:IdDiskon"`
	
	KategoriID       uint     `gorm:"column:id_kategori"`
	Kategori         Kategori `gorm:"foreignKey:KategoriID;references:IdKategori"`
	
	GambarDiskon     string   `gorm:"column:gambar_diskon"`
	KeteranganDiskon string   `gorm:"column:keterangan_diskon"`
}
