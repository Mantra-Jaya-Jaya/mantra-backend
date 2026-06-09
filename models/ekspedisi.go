package models

import "github.com/google/uuid"

func (Ekspedisi) TableName() string {
	return "ekspedisi"
}

type Ekspedisi struct {
	IdEkspedisi   uint      `gorm:"primaryKey;column:id_ekspedisi" json:"id_ekspedisi"`
	PublicId      uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();column:public_id;uniqueIndex" json:"public_id"`
	NamaEkspedisi string    `gorm:"column:nama_ekspedisi" json:"nama_ekspedisi"`
	KodeApi       string `gorm:"column:kode_api" json:"kode_api"`
	Logo          string `gorm:"column:logo" json:"logo"`
	Deskripsi     string `gorm:"column:deskripsi" json:"deskripsi"`
	IsActive      bool   `gorm:"column:is_active;default:true" json:"is_active"`

	Layanan []EkspedisiLayanan `gorm:"foreignKey:EkspedisiID;constraint:OnDelete:CASCADE" json:"layanan,omitempty"`
}
