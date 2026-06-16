package models

func (PengaturanToko) TableName() string {
	return "pengaturan_toko"
}

type PengaturanToko struct {
	Key   string `gorm:"column:key;primaryKey" json:"key"`
	Value string `gorm:"column:value" json:"value"`
}
