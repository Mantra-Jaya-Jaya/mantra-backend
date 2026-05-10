package models

import "time"

// Memaksa nama tabel
func (Kurir) TableName() string {
	return "kurir"
}

type Kurir struct {
	IdKurir            uint      `gorm:"primaryKey;column:id_kurir"`
	NoTelp             string    `gorm:"column:no_telp"`
	Alamat             string    `gorm:"type:text;column:alamat"`
	TanggalLahir       time.Time `gorm:"type:date;column:tanggal_lahir"`
	JenisKelamin       string    `gorm:"column:jenis_kelamin"`
	PendidikanTerakhir string    `gorm:"column:pendidikan_terakhir"`
	TempatLahir        string    `gorm:"column:tempat_lahir"`
	Nik                string    `gorm:"type:varchar(16);column:nik;unique;not null"`
	Foto               string    `gorm:"column:foto"`
	CreatedAt          time.Time `gorm:"column:created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at"`

	// Relasi ke User
	UserId uint `gorm:"column:id_user;unique"`
	User   User `gorm:"foreignKey:UserId;references:IdUser"`
}
