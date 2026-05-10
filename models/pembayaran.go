package models

func (Pembayaran) TableName() string {
	return "pembayaran"
}

type Pembayaran struct {
	IdPembayaran     uint   	`gorm:"primaryKey;column:id_pembayaran" json:"id_pembayaran"`
	OrderIdMidtrans  string 	`gorm:"column:order_id_midtrans" json:"order_id_midtrans"` // ID dari Midtrans
	PaymentType      string 	`gorm:"column:payment_type" json:"payment_type"`      // gopay, bank_transfer, dll
	StatusTransaksi  string 	`gorm:"column:status_transaksi" json:"status_transaksi"`  // settlement, pending, deny
	FraudStatus      string 	`gorm:"column:fraud_status" json:"fraud_status"`      // accept, challenge

	PesananID        uint   	`gorm:"column:id_pesanan" json:"id_pesanan"` // Foreign key ke tabel Pesanan
	Pesanan       	 Pesanan	`gorm:"foreignKey:PesananID;references:IdPesanan" json:"pesanan"`
}