package models

import (
	"time"
)

func (Pengantaran) TableName() string {
	return "pengantaran"
}

type Pengantaran struct {
	IdPengantaran       uint      `gorm:"primaryKey;column:id_pengantaran"`
	WaktuPickup         time.Time `gorm:"column:waktu_pickup"`
	WaktuSampai         time.Time `gorm:"column:waktu_sampai"`
	LastLatitude        float64   `gorm:"column:last_latitude"`
	LastLongitude       float64   `gorm:"column:last_longitude"`
	FotoBuktiPengiriman string    `gorm:"column:foto_bukti_pengiriman"`

	PesananID uint    `gorm:"column:id_pesanan"`
	Pesanan   Pesanan `gorm:"foreignKey:PesananID;references:IdPesanan"`

	KurirID uint  `gorm:"column:id_kurir"`
	Kurir   Kurir `gorm:"foreignKey:KurirID;references:IdKurir"`

	StatusPengantaranID uint              `gorm:"column:id_status_pengantaran"`
	StatusPengantaran   StatusPengantaran `gorm:"foreignKey:StatusPengantaranID;references:IdStatusPengantaran"`

	EkspedisiID uint      `gorm:"column:id_ekspedisi"`
	Ekspedisi   Ekspedisi `gorm:"foreignKey:EkspedisiID;references:IdEkspedisi"`
}
