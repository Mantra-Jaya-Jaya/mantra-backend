package models

import (
	"time"

	"github.com/google/uuid"
)

func (Karyawan) TableName() string {
	return "karyawan"
}

type Karyawan struct {
	IdKaryawan         uint      `gorm:"primaryKey;column:id_karyawan" json:"id_karyawan"`
	PublicId           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();column:public_id;uniqueIndex" json:"public_id"`
	NoTelp             string    `gorm:"type:varchar(15);column:no_telp;unique;not null" json:"no_telp"`
	TempatLahir        string    `gorm:"column:tempat_lahir" json:"tempat_lahir"`
	TanggalLahir       time.Time `gorm:"type:date;column:tanggal_lahir" json:"tanggal_lahir"`
	JenisKelamin       string    `gorm:"column:jenis_kelamin" json:"jenis_kelamin"`
	Alamat             string    `gorm:"column:alamat" json:"alamat"`
	PendidikanTerakhir string    `gorm:"column:pendidikan_terakhir" json:"pendidikan_terakhir"`
	Nik                string    `gorm:"type:varchar(16);column:nik;unique;not null" json:"nik"`
	Status             string         `gorm:"column:status" json:"-"` // TODO: hapus setelah migration
	StatusKaryawanID   uint           `gorm:"column:id_status_karyawan" json:"id_status_karyawan"`
	StatusKaryawanRel  *StatusKaryawan `gorm:"foreignKey:StatusKaryawanID;references:IdStatusKaryawan" json:"status,omitempty"`

	UserID uint `gorm:"column:id_user;unique" json:"id_user"`
	User   User `gorm:"foreignKey:UserID;references:IdUser" json:"user"`
}
