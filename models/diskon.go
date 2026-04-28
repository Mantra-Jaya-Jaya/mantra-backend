package models

import "time"

func (Diskon) TableName() string {
	return "diskon"
}

type Diskon struct {
	IdDiskon    uint      `gorm:"primaryKey;column:id_diskon"`
	NamaDiskon  string    `gorm:"column:nama_diskon"`
	BesarDiskon int       `gorm:"column:besar_diskon"`
	TglMulai    time.Time `gorm:"column:tgl_mulai;type:date"`
	TglSelesai  time.Time `gorm:"column:tgl_selesai;type:date"`
}
