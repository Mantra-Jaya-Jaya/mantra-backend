package models

import "time"

// Memaksa nama tabel
func (Kurir) TableName() string {
	return "kurir"
}

type Kurir struct {
	IdKurir      uint      `gorm:"primaryKey;column:id_kurir"`
	NoTelp       string    `gorm:"column:no_telp"`
	Alamat       string    `gorm:"type:text;column:alamat"`
	TanggalLahir time.Time `gorm:"type:date;column:tanggal_lahir"`
	Foto         string    `gorm:"column:foto"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`

	// Relasi ke User
	UserId uint `gorm:"column:id_user;unique"`
	User   User `gorm:"foreignKey:UserId;references:IdUser"`
}
