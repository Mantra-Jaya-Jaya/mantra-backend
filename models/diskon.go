package models

import "time"

func (Diskon) TableName() string {
	return "diskon"
}

type Diskon struct {
	IdDiskon     uint      `gorm:"primaryKey;column:id_diskon" json:"id_diskon"`
	NamaDiskon   string    `gorm:"column:nama_diskon" json:"nama_diskon"`
	BesarDiskon  int       `gorm:"column:besar_diskon" json:"besar_diskon"`
	BannerDiskon string    `gorm:"column:banner_diskon" json:"banner_diskon"`
	TglMulai     time.Time `gorm:"column:tgl_mulai;type:date" json:"tgl_mulai"`
	TglSelesai   time.Time `gorm:"column:tgl_selesai;type:date" json:"tgl_selesai"`
}
