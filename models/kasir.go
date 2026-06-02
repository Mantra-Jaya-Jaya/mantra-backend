package models

import (
	"github.com/google/uuid"
)

func (Kasir) TableName() string {
	return "kasir"
}

type Kasir struct {
	IdKasir    uint      `gorm:"primaryKey;column:id_kasir" json:"id_kasir"`
	PublicId   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();column:public_id;uniqueIndex" json:"public_id"`
	Shift      string    `gorm:"column:shift" json:"shift"`

	KaryawanId uint     `gorm:"column:id_karyawan;unique" json:"id_karyawan"`
	Karyawan   Karyawan `gorm:"foreignKey:KaryawanId;references:IdKaryawan" json:"karyawan"`
}
