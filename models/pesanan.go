package models

import (
	"time"

	"github.com/google/uuid"
)

func (Pesanan) TableName() string {
	return "pesanan"
}

type Pesanan struct {
	IdPesanan       uint      `gorm:"primaryKey;column:id_pesanan" json:"id_pesanan"`
	PublicId        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();column:public_id;uniqueIndex" json:"public_id"`
	TotalPembayaran int       `gorm:"column:total_pembayaran" json:"total_pembayaran"`
	TanggalPesanan  time.Time `gorm:"column:tanggal_pesanan" json:"tanggal_pesanan"`
	TipePesanan     string    `gorm:"column:tipe_pesanan" json:"tipe_pesanan"`
	StatusPesanan   string    `gorm:"column:status_pesanan" json:"status_pesanan"`

	CustomerId uint     `gorm:"column:id_customer" json:"id_customer"`
	KasirId    *uint    `gorm:"column:id_kasir" json:"id_kasir"`
	AlamatId   *uint    `gorm:"column:id_alamat" json:"id_alamat"`
	Customer   Customer `gorm:"foreignKey:CustomerId;references:IdCustomer" json:"customer"`
	Kasir      Kasir    `gorm:"foreignKey:KasirId;references:IdKasir" json:"kasir"`
	Alamat     *Alamat  `gorm:"foreignKey:AlamatId;references:IdAlamat" json:"alamat"`

	EkspedisiID        *uint             `gorm:"column:id_ekspedisi" json:"id_ekspedisi,omitempty"`
	Ekspedisi          *Ekspedisi        `gorm:"foreignKey:EkspedisiID;references:IdEkspedisi" json:"ekspedisi,omitempty"`
	LayananEkspedisiID *uint             `gorm:"column:id_layanan_ekspedisi" json:"id_layanan_ekspedisi,omitempty"`
	LayananEkspedisi   *EkspedisiLayanan `gorm:"foreignKey:LayananEkspedisiID;references:IdEkspedisiLayanan" json:"layanan_ekspedisi,omitempty"`
	OngkosKirim        int               `gorm:"column:ongkos_kirim;default:0" json:"ongkos_kirim"`
	Catatan            string            `gorm:"column:catatan" json:"catatan"`

	DetailPesanan []DetailPesanan `gorm:"foreignKey:PesananId;references:IdPesanan" json:"detail_pesanan,omitempty"`
	Pembayaran    *Pembayaran     `gorm:"foreignKey:PesananID" json:"pembayaran,omitempty"`
}
