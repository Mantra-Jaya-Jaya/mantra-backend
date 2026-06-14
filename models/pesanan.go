package models

import (
	"time"

	"github.com/google/uuid"
)

func (Pesanan) TableName() string {
	return "pesanan"
}

type Pesanan struct {
	IdPesanan       uint           `gorm:"primaryKey;column:id_pesanan" json:"id_pesanan"`
	PublicId        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();column:public_id;uniqueIndex" json:"public_id"`
	TotalPembayaran int            `gorm:"column:total_pembayaran" json:"total_pembayaran"`
	TanggalPesanan  time.Time      `gorm:"column:tanggal_pesanan" json:"tanggal_pesanan"`
	TipePesananID   uint           `gorm:"column:id_tipe_pesanan;not null" json:"id_tipe_pesanan"`
	TipePesananRel  *TipePesanan   `gorm:"foreignKey:TipePesananID;references:IdTipePesanan" json:"tipe_pesanan,omitempty"`
	StatusPesananID uint           `gorm:"column:id_status_pesanan;not null" json:"id_status_pesanan"`
	StatusPesanan   *StatusPesanan `gorm:"foreignKey:StatusPesananID;references:IdStatusPesanan" json:"status_pesanan,omitempty"`

	CustomerID uint     `gorm:"column:id_customer;not null" json:"id_customer"`
	KasirID    *uint    `gorm:"column:id_kasir" json:"id_kasir"`
	AlamatID   *uint    `gorm:"column:id_alamat" json:"id_alamat"`
	Customer   Customer `gorm:"foreignKey:CustomerID;references:IdCustomer" json:"customer"`
	Kasir      *Kasir   `gorm:"foreignKey:KasirID;references:IdKasir" json:"kasir,omitempty"`
	Alamat     *Alamat  `gorm:"foreignKey:AlamatID;references:IdAlamat" json:"alamat"`

	EkspedisiID        *uint             `gorm:"column:id_ekspedisi" json:"id_ekspedisi,omitempty"`
	Ekspedisi          *Ekspedisi        `gorm:"foreignKey:EkspedisiID;references:IdEkspedisi" json:"ekspedisi,omitempty"`
	LayananEkspedisiID *uint             `gorm:"column:id_layanan_ekspedisi" json:"id_layanan_ekspedisi,omitempty"`
	LayananEkspedisi   *EkspedisiLayanan `gorm:"foreignKey:LayananEkspedisiID;references:IdEkspedisiLayanan" json:"layanan_ekspedisi,omitempty"`
	OngkosKirim        int               `gorm:"column:ongkos_kirim;default:0" json:"ongkos_kirim"`
	Catatan            string            `gorm:"column:catatan" json:"catatan"`
	NomorResi          *string           `gorm:"column:nomor_resi" json:"nomor_resi,omitempty"`

	DetailPesanan []DetailPesanan `gorm:"foreignKey:PesananID;references:IdPesanan" json:"detail_pesanan,omitempty"`
	Pembayaran    *Pembayaran     `gorm:"foreignKey:PesananID" json:"pembayaran,omitempty"`
}
