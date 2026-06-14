package models

import "github.com/google/uuid"

func (MetodePembayaran) TableName() string {
	return "metode_pembayaran"
}

type MetodePembayaran struct {
	IdMetodePembayaran uint      `gorm:"primaryKey;column:id_metode_pembayaran" json:"id_metode_pembayaran"`
	PublicId           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();column:public_id;uniqueIndex" json:"public_id"`
	NamaMetode         string    `gorm:"column:nama_metode" json:"nama_metode"`
	KodeMetode         string    `gorm:"column:kode_metode" json:"kode_metode"`
	Penyedia           string    `gorm:"column:penyedia" json:"penyedia"`
	Icon               string    `gorm:"column:icon" json:"icon"`
	Urutan             int       `gorm:"column:urutan;default:0" json:"urutan"`
	IsActive           bool      `gorm:"column:is_active;default:true" json:"is_active"`
}
