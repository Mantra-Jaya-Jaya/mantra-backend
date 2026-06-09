package models

import "time"

func (Pembayaran) TableName() string {
	return "pembayaran"
}

type Pembayaran struct {
	IdPembayaran    uint   `gorm:"primaryKey;column:id_pembayaran" json:"id_pembayaran"`
	OrderIdMidtrans string `gorm:"column:order_id_midtrans" json:"order_id_midtrans"`
	PaymentType     string `gorm:"column:payment_type" json:"payment_type"`
	StatusTransaksi string `gorm:"column:status_transaksi" json:"status_transaksi"`
	FraudStatus     string `gorm:"column:fraud_status" json:"fraud_status"`

	PesananID uint    `gorm:"column:id_pesanan" json:"id_pesanan"`
	Pesanan   Pesanan `gorm:"foreignKey:PesananID;references:IdPesanan" json:"pesanan"`

	MetodePembayaranID  *uint             `gorm:"column:id_metode_pembayaran" json:"id_metode_pembayaran,omitempty"`
	MetodePembayaran    *MetodePembayaran `gorm:"foreignKey:MetodePembayaranID;references:IdMetodePembayaran" json:"metode_pembayaran,omitempty"`
	TransaksiMidtransID string            `gorm:"column:transaksi_midtrans_id" json:"transaksi_midtrans_id,omitempty"`
	WaktuPembayaran     *time.Time        `gorm:"column:waktu_pembayaran" json:"waktu_pembayaran,omitempty"`
	TotalDibayar        int               `gorm:"column:total_dibayar;default:0" json:"total_dibayar"`

	DetailPembayaran []DetailPembayaran `gorm:"foreignKey:PembayaranID" json:"detail_pembayaran,omitempty"`
}
