package models

import (
	"time"

	"github.com/google/uuid"
)

func (Pengantaran) TableName() string {
	return "pengantaran"
}

type Pengantaran struct {
	IdPengantaran uint      `gorm:"primaryKey;column:id_pengantaran" json:"id_pengantaran"`
	PublicId      uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();column:public_id;uniqueIndex" json:"public_id"`

	WaktuPickup *time.Time `gorm:"column:waktu_pickup" json:"waktu_pickup,omitempty"`
	WaktuSampai *time.Time `gorm:"column:waktu_sampai" json:"waktu_sampai,omitempty"`

	LastLatitude        float64 `gorm:"column:last_latitude" json:"last_latitude,omitempty"`
	LastLongitude       float64 `gorm:"column:last_longitude" json:"last_longitude,omitempty"`
	FotoBuktiPengiriman string  `gorm:"column:foto_bukti_pengiriman" json:"foto_bukti_pengiriman,omitempty"`

	PesananID uint     `gorm:"column:id_pesanan;not null" json:"id_pesanan"`
	Pesanan   *Pesanan `gorm:"foreignKey:PesananID;references:IdPesanan" json:"pesanan,omitempty"`

	KurirID *uint  `gorm:"column:id_kurir" json:"id_kurir,omitempty"`
	Kurir   *Kurir `gorm:"foreignKey:KurirID;references:IdKurir" json:"kurir,omitempty"`

	StatusPengantaranID uint               `gorm:"column:id_status_pengantaran;not null" json:"id_status_pengantaran"`
	StatusPengantaran   *StatusPengantaran `gorm:"foreignKey:StatusPengantaranID;references:IdStatusPengantaran" json:"status_pengantaran,omitempty"`

	EkspedisiID *uint      `gorm:"column:id_ekspedisi" json:"id_ekspedisi,omitempty"`
	Ekspedisi   *Ekspedisi `gorm:"foreignKey:EkspedisiID;references:IdEkspedisi" json:"ekspedisi,omitempty"`
}
