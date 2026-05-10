package models

import "time"

// Memaksa nama tabel
func (Kasir) TableName() string {
	return "kasir"
}

type Kasir struct {
	// Nama variabel WAJIB Kapital (IdRole), nama di DB diatur lewat Tag (id_role)
	IdKasir            uint      `gorm:"primaryKey;column:id_kasir" json:"id_kasir"`
	NoTelp             string    `gorm:"column:no_telp" json:"no_telp"`
	TempatLahir        string    `gorm:"column:tempat_lahir" json:"tempat_lahir"`
	TanggalLahir       time.Time `gorm:"type:date;column:tanggal_lahir" json:"tanggal_lahir"`
	JenisKelamin       string    `gorm:"column:jenis_kelamin" json:"jenis_kelamin"`
	Alamat             string    `gorm:"column:alamat" json:"alamat"`
	PendidikanTerakhir string    `gorm:"column:pendidikan_terakhir" json:"pendidikan_terakhir"`
	Nik                string    `gorm:"type:varchar(16);column:nik;unique;not null" json:"nik"`

	//Relasi ke User
	UserId uint `gorm:"column:id_user;unique" json:"id_user"`

	//Relasi  ke tabel user
	User User `gorm:"foreignKey:UserId;references:IdUser" json:"user"`
}
