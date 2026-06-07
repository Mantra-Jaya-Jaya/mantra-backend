package models

import "github.com/google/uuid"

func (Kategori) TableName() string {
	return "kategori"
}

type Kategori struct {
	IdKategori   uint      `gorm:"primaryKey;column:id_kategori" json:"id_kategori"`
	PublicId     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();column:public_id;uniqueIndex" json:"public_id"`
	NamaKategori string    `gorm:"column:nama_kategori" json:"nama_kategori"`
	IconKategori string    `gorm:"column:icon_kategori" json:"icon_kategori"`
}
