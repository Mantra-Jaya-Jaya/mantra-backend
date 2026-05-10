package models

import "time"

func (Pembayaran) TableName() string {
	return "pembayaran"
}

type Pembayaran struct {
	IdPembayaran    uint      `gorm:"primaryKey;column:id_pembayaran"`
	OrderIdMidtrans string    `gorm:"column:order_id_midtrans"` // ID dari Midtrans, null jika bayar tunai
	PaymentType     string    `gorm:"column:payment_type"`      // cash | qris | ewallet | va | kartu
	StatusTransaksi string    `gorm:"column:status_transaksi"`  // pending | settlement | cancel | expire | deny
	FraudStatus     string    `gorm:"column:fraud_status"`      // accept | challenge | deny
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`

	// Foreign Key ke Pesanan
	PesananID uint    `gorm:"column:id_pesanan"`
	Pesanan   Pesanan `gorm:"foreignKey:PesananID;references:IdPesanan"`
}