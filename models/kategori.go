package models

func (Kategori) TableName() string {
	return "kategori"
}

type Kategori struct {
	IdKategori   uint   `gorm:"primaryKey;column:id_kategori"`
	NamaKategori string `gorm:"column:nama_kategori"`
	IconKategori string `gorm:"column:icon_kategori"`
}
