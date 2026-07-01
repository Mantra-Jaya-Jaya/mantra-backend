package models

import "github.com/google/uuid"

func (DetailPembayaran) TableName() string {
	return "detail_pembayaran"
}

type DetailPembayaran struct {
	IdDetailPembayaran uint      `gorm:"primaryKey;column:id_detail_pembayaran" json:"id_detail_pembayaran"`
	PublicId           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();column:public_id;uniqueIndex" json:"public_id"`
	PembayaranID       uint      `gorm:"column:id_pembayaran;not null" json:"id_pembayaran"`
	KanalPembayaran    string    `gorm:"column:kanal_pembayaran" json:"kanal_pembayaran"`
	NomorVA            string    `gorm:"column:nomor_va" json:"nomor_va"`
	BillKey            string    `gorm:"column:bill_key" json:"bill_key"`
	BillCode           string    `gorm:"column:bill_code" json:"bill_code"`
	NamaBank           string    `gorm:"column:nama_bank" json:"nama_bank"`
	MerchantID         string    `gorm:"column:merchant_id" json:"merchant_id"`
	QrCode             string    `gorm:"column:qr_code_url" json:"qr_code_url"`

	Pembayaran Pembayaran `gorm:"foreignKey:PembayaranID;references:IdPembayaran" json:"pembayaran"`
}
