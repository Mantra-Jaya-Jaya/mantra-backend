package models

func (TipePembayaran) TableName() string {
	return "tipe_pembayaran"
}

type TipePembayaran struct {
	IdTipePembayaran uint   `gorm:"primaryKey;column:id" json:"id_tipe_pembayaran"`
	NamaTipe         string `gorm:"column:nama_tipe;unique" json:"nama_tipe"`
}
