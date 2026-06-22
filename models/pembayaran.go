package models

import "time"

func (Pembayaran) TableName() string {
	return "pembayaran"
}

type Pembayaran struct {
	IdPembayaran       uint             `gorm:"primaryKey;column:id_pembayaran" json:"id_pembayaran"`
	OrderIdMidtrans    string           `gorm:"column:order_id_midtrans" json:"order_id_midtrans"`
	TipePembayaranID   uint             `gorm:"column:id_tipe_pembayaran;not null" json:"id_tipe_pembayaran"`
	TipePembayaranRel  *TipePembayaran  `gorm:"foreignKey:TipePembayaranID;references:IdTipePembayaran" json:"tipe_pembayaran,omitempty"`
	StatusTransaksiID  uint             `gorm:"column:id_status_transaksi;not null" json:"id_status_transaksi"`
	StatusTransaksiRel *StatusTransaksi `gorm:"foreignKey:StatusTransaksiID;references:IdStatusTransaksi" json:"status_transaksi,omitempty"`
	FraudStatusID      uint             `gorm:"column:id_fraud_status;not null" json:"id_fraud_status"`
	FraudStatusRel     *FraudStatus     `gorm:"foreignKey:FraudStatusID;references:IdFraudStatus" json:"fraud_status,omitempty"`

	PesananID uint    `gorm:"column:id_pesanan;not null" json:"id_pesanan"`
	Pesanan   Pesanan `gorm:"foreignKey:PesananID;references:IdPesanan" json:"pesanan"`

	MetodePembayaranID  *uint             `gorm:"column:id_metode_pembayaran" json:"id_metode_pembayaran,omitempty"`
	MetodePembayaran    *MetodePembayaran `gorm:"foreignKey:MetodePembayaranID;references:IdMetodePembayaran" json:"metode_pembayaran,omitempty"`
	TransaksiMidtransID string            `gorm:"column:transaksi_midtrans_id" json:"transaksi_midtrans_id,omitempty"`
	WaktuPembayaran     *time.Time        `gorm:"column:waktu_pembayaran" json:"waktu_pembayaran,omitempty"`
	BatasWaktuPembayaran *time.Time       `gorm:"column:batas_waktu_pembayaran" json:"batas_waktu_pembayaran,omitempty"`
	TotalDibayar        int               `gorm:"column:total_dibayar;default:0" json:"total_dibayar"`

	DetailPembayaran []DetailPembayaran `gorm:"foreignKey:PembayaranID" json:"detail_pembayaran,omitempty"`
}
