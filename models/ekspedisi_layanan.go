package models

import "github.com/google/uuid"

func (EkspedisiLayanan) TableName() string {
	return "ekspedisi_layanan"
}

type EkspedisiLayanan struct {
	IdEkspedisiLayanan uint      `gorm:"primaryKey;column:id_ekspedisi_layanan" json:"id_ekspedisi_layanan"`
	PublicId           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();column:public_id;uniqueIndex" json:"public_id"`
	EkspedisiID        uint      `gorm:"column:id_ekspedisi;not null;index" json:"id_ekspedisi"`
	NamaLayanan        string    `gorm:"column:nama_layanan" json:"nama_layanan"`
	KodeLayanan        string    `gorm:"column:kode_layanan" json:"kode_layanan"`
	Deskripsi          string    `gorm:"column:deskripsi" json:"deskripsi"`
	EstimasiMin        int       `gorm:"column:estimasi_min" json:"estimasi_min"`
	EstimasiMax        int       `gorm:"column:estimasi_max" json:"estimasi_max"`
	IsActive           bool      `gorm:"column:is_active;default:true" json:"is_active"`

	Ekspedisi Ekspedisi `gorm:"foreignKey:EkspedisiID;references:IdEkspedisi" json:"ekspedisi"`
}
