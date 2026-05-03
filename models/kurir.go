package models

import "time"

// Memaksa nama tabel
func (Kurir) TableName() string {
	return "kurir"
}

type Kurir struct {
	// Nama variabel WAJIB Kapital (IdRole), nama di DB diatur lewat Tag (id_role)
	IdKurir            uint      `gorm:"primaryKey;column:id_kurir"`
	NoTelp             string    `gorm:"type:varchar(15);column:no_telp;unique;not null"`
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
