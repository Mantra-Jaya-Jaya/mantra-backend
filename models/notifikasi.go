package models

import "time"

func (Notifikasi) TableName() string {
	return "notifikasi"
}

type Notifikasi struct {
	IdNotifikasi uint      `gorm:"primaryKey;column:id_notifikasi"`
	Judul        string    `gorm:"column:judul"`
	Pesan        string    `gorm:"type:text;column:pesan"`
	Status       string    `gorm:"column:status"` // unread | read
	CreatedAt    time.Time `gorm:"column:created_at"`

	// Relasi ke User
	UserID uint `gorm:"column:id_user"`
	User   User `gorm:"foreignKey:UserID;references:IdUser"`
}
