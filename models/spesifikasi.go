package models

func (Spesifikasi) TableName() string {
	return "spesifikasi"
}

type Spesifikasi struct {
	IdSpesifikasi   uint   `gorm:"primaryKey;column:id_spesifikasi" json:"id_spesifikasi"`
	NamaSpesifikasi string `gorm:"column:nama_spesifikasi" json:"nama_spesifikasi"`
}
