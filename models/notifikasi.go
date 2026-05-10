package models

import "time"

func (Notifikasi) TableName() string {
	return "notifikasi"
}

type Notifikasi struct {
	IdNotifikasi uint   `gorm:"primaryKey;column:id_notifikasi" json:"id_notifikasi"`
	
	UserID       uint   `gorm:"column:id_user" json:"id_user"`
	User         User   `gorm:"foreignKey:UserID;references:IdUser" json:"user"`
	
	Judul        string `gorm:"column:judul" json:"judul"`
	Pesan        string `gorm:"column:pesan" json:"pesan"`
	Status       string `gorm:"column:status" json:"status"`
}
