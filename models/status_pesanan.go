package models

func (StatusPesanan) TableName() string {
	return "status_pesanan"
}

type StatusPesanan struct {
	IdStatusPesanan uint   `gorm:"primaryKey;column:id" json:"id_status_pesanan"`
	NamaStatus      string `gorm:"column:nama_status" json:"nama_status"`
}
