package models

import "time"

func (Notifikasi) TableName() string {
	return "notifikasi"
}

type Notifikasi struct {
	IdNotifikasi uint `gorm:"primaryKey;column:id_notifikasi" json:"id_notifikasi"`

	UserID uint `gorm:"column:id_user" json:"id_user"`
	User   User `gorm:"foreignKey:UserID;references:IdUser" json:"user"`

	Judul               string `gorm:"column:judul" json:"judul"`
	Pesan               string `gorm:"column:pesan" json:"pesan"`
	Status              string `gorm:"column:status" json:"-"` // TODO: hapus setelah migration
	StatusNotifikasiID  uint   `gorm:"column:id_status_notifikasi" json:"id_status_notifikasi"`

	StatusNotifikasiRel *StatusNotifikasi `gorm:"foreignKey:StatusNotifikasiID;references:IdStatusNotifikasi" json:"status_notifikasi,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
