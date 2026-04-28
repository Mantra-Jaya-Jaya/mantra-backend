package models

func (Ekspedisi) TableName() string {
	return "ekspedisi"
}

type Ekspedisi struct {
	IdEkspedisi   uint   `gorm:"primaryKey;column:id_ekspedisi"`
	NamaEkspedisi string `gorm:"column:nama_ekspedisi"`
	KodeApi       string `gorm:"column:kode_api"`
}
