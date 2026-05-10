package models

func (Kategori) TableName() string {
	return "kategori"
}

type Kategori struct {
	IdKategori   uint   `gorm:"primaryKey;column:id_kategori" json:"id_kategori"`
	NamaKategori string `gorm:"column:nama_kategori" json:"nama_kategori"`
	IconKategori string `gorm:"column:icon_kategori" json:"icon_kategori"`
}
