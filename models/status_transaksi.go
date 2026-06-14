package models

func (StatusTransaksi) TableName() string {
	return "status_transaksi"
}

type StatusTransaksi struct {
	IdStatusTransaksi uint   `gorm:"primaryKey;column:id" json:"id_status_transaksi"`
	NamaStatus        string `gorm:"column:nama_status;unique" json:"nama_status"`
}
