package models

func (StatusNotifikasi) TableName() string {
	return "status_notifikasi"
}

type StatusNotifikasi struct {
	IdStatusNotifikasi uint   `gorm:"primaryKey;column:id" json:"id_status_notifikasi"`
	NamaStatus         string `gorm:"column:nama_status;unique" json:"nama_status"`
}
