package models

import "time"

// Memaksa nama tabel
func (Kasir) TableName() string {
	return "kasir"
}

type Kasir struct {
	IdKasir            uint      `gorm:"primaryKey;column:id_kasir"`
	NoTelp             string    `gorm:"column:no_telp"`
	Alamat             string    `gorm:"type:text;column:alamat"`
	TanggalLahir       time.Time `gorm:"type:date;column:tanggal_lahir"`
	Foto               string    `gorm:"column:foto"`
	JenisKelamin       string    `gorm:"column:jenis_kelamin"`
	PendidikanTerakhir string    `gorm:"column:pendidikan_terakhir"`
	Nik                string    `gorm:"column:nik"`
	Shift              string    `gorm:"column:shift"` // Pagi | Siang | Malam
	CreatedAt          time.Time `gorm:"column:created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at"`

	// Relasi ke User
	UserId uint `gorm:"column:id_user;unique"`
	User   User `gorm:"foreignKey:UserId;references:IdUser"`
}
