package models

import (
	"time"
)

func (Pesanan) TableName() string {
	return "pesanan"
}

type Pesanan struct {
	IdPesanan        uint      `gorm:"primaryKey;column:id_pesanan"`	
	TotalPembayaran  int       `gorm:"column:total_pembayaran"`
	TanggalPesanan   time.Time `gorm:"column:tanggal_pesanan"`
	TipePesanan      string    `gorm:"column:tipe_pesanan"`   // Online / Offline
	StatusPesanan    string    `gorm:"column:status_pesanan"` // Dikemas, Dikirim, Selesai, dll

	CustomerId       uint      `gorm:"column:id_customer"`
	KasirId          uint      `gorm:"column:id_kasir"`
	AlamatId         *uint     `gorm:"column:id_alamat"` // Pake pointer karena bisa NULL (jika takeaway/offline)\
	Customer         Customer  `gorm:"foreignKey:CustomerId;references:IdCustomer"`
	Kasir            Kasir		 `gorm:"foreignKey:KasirId;references:IdKasir"` 
	Alamat           *Alamat   `gorm:"foreignKey:AlamatId;references:IdAlamat"` // Pake pointer karena bisa NULL (jika takeaway/offline)
}