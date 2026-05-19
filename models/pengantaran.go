package models

import (
	"time"

	"github.com/google/uuid"
)

func (Pengantaran) TableName() string {
	return "pengantaran"
}

type Pengantaran struct {
	IdPengantaran       uint      `gorm:"primaryKey;column:id_pengantaran" json:"id_pengantaran"`
	PublicId            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();column:public_id;uniqueIndex" json:"public_id"`
	WaktuPickup         time.Time `gorm:"column:waktu_pickup" json:"waktu_pickup"`
	WaktuSampai         time.Time `gorm:"column:waktu_sampai" json:"waktu_sampai"`
	LastLatitude        float64   `gorm:"column:last_latitude" json:"last_latitude"`
	LastLongitude       float64   `gorm:"column:last_longitude" json:"last_longitude"`
	FotoBuktiPengiriman string    `gorm:"column:foto_bukti_pengiriman" json:"foto_bukti_pengiriman"`

	PesananID uint    `gorm:"column:id_pesanan" json:"id_pesanan"`
	Pesanan   Pesanan `gorm:"foreignKey:PesananID;references:IdPesanan" json:"pesanan"`

	KurirID uint  `gorm:"column:id_kurir" json:"id_kurir"`
	Kurir   Kurir `gorm:"foreignKey:KurirID;references:IdKurir" json:"kurir"`

	StatusPengantaranID uint              `gorm:"column:id_status_pengantaran" json:"id_status_pengantaran"`
	StatusPengantaran   StatusPengantaran `gorm:"foreignKey:StatusPengantaranID;references:IdStatusPengantaran" json:"status_pengantaran"`

	EkspedisiID uint      `gorm:"column:id_ekspedisi" json:"id_ekspedisi"`
	Ekspedisi   Ekspedisi `gorm:"foreignKey:EkspedisiID;references:IdEkspedisi" json:"ekspedisi"`
}
