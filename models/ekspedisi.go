package models

func (Ekspedisi) TableName() string {
	return "ekspedisi"
}

type Ekspedisi struct {
	IdEkspedisi   uint   `gorm:"primaryKey;column:id_ekspedisi" json:"id_ekspedisi"`
	NamaEkspedisi string `gorm:"column:nama_ekspedisi" json:"nama_ekspedisi"`
	KodeApi       string `gorm:"column:kode_api" json:"kode_api"`
}
