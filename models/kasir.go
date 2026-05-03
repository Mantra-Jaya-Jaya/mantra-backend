package models

import "time"

// Memaksa nama tabel
func (Kasir) TableName() string {
	return "kasir"
}

type Kasir struct {
	// Nama variabel WAJIB Kapital (IdRole), nama di DB diatur lewat Tag (id_role)
	IdKasir            uint      `gorm:"primaryKey;column:id_kasir"`
	NoTelp             string    `gorm:"column:no_telp"`
	TempatLahir        string    `gorm:"column:tempat_lahir"`
	TanggalLahir       time.Time `gorm:"type:date;column:tanggal_lahir"`
	JenisKelamin       string    `gorm:"column:jenis_kelamin"`
	Alamat             string    `gorm:"column:alamat"`
	PendidikanTerakhir string    `gorm:"column:pendidikan_terakhir"`
	Nik                string    `gorm:"type:varchar(16);column:nik;unique;not null"`

	//Relasi ke User
	UserId uint `gorm:"column:id_user;unique"`

	//Relasi  ke tabel user
	User User `gorm:"foreignKey:UserId;references:IdUser"`
}
