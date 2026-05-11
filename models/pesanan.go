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
	TipePesanan     string    `gorm:"column:tipe_pesanan" json:"tipe_pesanan"`     // Online / Offline
	StatusPesanan   string    `gorm:"column:status_pesanan" json:"status_pesanan"` // Dikemas, Dikirim, Selesai, dll

	CustomerId uint     `gorm:"column:id_customer" json:"id_customer"`
	KasirId    uint     `gorm:"column:id_kasir" json:"id_kasir"`
	AlamatId   *uint    `gorm:"column:id_alamat" json:"id_alamat"` // Pake pointer karena bisa NULL (jika takeaway/offline)\
	Customer   Customer `gorm:"foreignKey:CustomerId;references:IdCustomer" json:"customer"`
	Kasir      Kasir    `gorm:"foreignKey:KasirId;references:IdKasir" json:"kasir"`
	Alamat     *Alamat  `gorm:"foreignKey:AlamatId;references:IdAlamat" json:"alamat"` // Pake pointer karena bisa NULL (jika takeaway/offline)
}
