package models

func (TipePesanan) TableName() string {
	return "tipe_pesanan"
}

type TipePesanan struct {
	IdTipePesanan uint   `gorm:"primaryKey;column:id" json:"id_tipe_pesanan"`
	NamaTipe      string `gorm:"column:nama_tipe;unique" json:"nama_tipe"`
}
