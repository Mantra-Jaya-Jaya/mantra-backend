package models

func (Satuan) TableName() string {
  return "satuan"
}

type Satuan struct {
  IdSatuan      uint   `gorm:"primaryKey;column:id_satuan" json:"id_satuan"`
  NamaSatuan    string `gorm:"column:nama_satuan" json:"nama_satuan"`
}