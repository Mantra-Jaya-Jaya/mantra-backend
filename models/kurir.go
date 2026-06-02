package models

import (
	"github.com/google/uuid"
)

func (Kurir) TableName() string {
	return "kurir"
}

type Kurir struct {
	IdKurir    uint      `gorm:"primaryKey;column:id_kurir" json:"id_kurir"`
	PublicId   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();column:public_id;uniqueIndex" json:"public_id"`

	KaryawanId uint     `gorm:"column:id_karyawan;unique" json:"id_karyawan"`
	Karyawan   Karyawan `gorm:"foreignKey:KaryawanId;references:IdKaryawan" json:"karyawan"`
}
