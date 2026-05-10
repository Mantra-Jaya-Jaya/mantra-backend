package models

import (
	"time"
)

func (Pesanan) TableName() string {
	return "pesanan"
}

type Pesanan struct {
	IdPesanan       uint      `gorm:"primaryKey;column:id_pesanan"`
	TotalPembayaran int       `gorm:"column:total_pembayaran"`
	TanggalPesanan  time.Time `gorm:"column:tanggal_pesanan"`
	TipePesanan     string    `gorm:"column:tipe_pesanan"`   // online | offline
	StatusPesanan   string    `gorm:"column:status_pesanan"` // menunggu_pembayaran | diproses | dikemas | dikirim | selesai | dibatalkan
	UpdatedAt       time.Time `gorm:"column:updated_at"`

	// Foreign Keys (beberapa nullable)
	CustomerId *uint    `gorm:"column:id_customer"` // Null jika pesanan offline/walk-in
	KasirId    *uint    `gorm:"column:id_kasir"`    // Null jika pesanan online mandiri
	AlamatId   *uint    `gorm:"column:id_alamat"`   // Null jika pesanan offline

	// Relasi
	Customer *Customer `gorm:"foreignKey:CustomerId;references:IdCustomer"`
	Kasir    *Kasir    `gorm:"foreignKey:KasirId;references:IdKasir"`
	Alamat   *Alamat   `gorm:"foreignKey:AlamatId;references:IdAlamat"`
}