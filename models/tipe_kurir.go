package models

func (TipeKurir) TableName() string {
	return "tipe_kurir"
}

type TipeKurir struct {
	IdTipeKurir uint   `gorm:"primaryKey;column:id" json:"id_tipe_kurir"`
	NamaTipe    string `gorm:"column:nama_tipe;unique" json:"nama_tipe"`
}
