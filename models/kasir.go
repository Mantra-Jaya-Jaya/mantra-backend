package models

import (
	"github.com/google/uuid"
)

func (Kasir) TableName() string {
	return "kasir"
}

type Kasir struct {
	IdKasir       uint        `gorm:"primaryKey;column:id_kasir" json:"id_kasir"`
	PublicId      uuid.UUID   `gorm:"type:uuid;default:gen_random_uuid();column:public_id;uniqueIndex" json:"public_id"`
	ShiftKasirID  uint        `gorm:"column:id_shift_kasir;not null" json:"id_shift_kasir"`
	ShiftKasirRel *ShiftKasir `gorm:"foreignKey:ShiftKasirID;references:IdShiftKasir" json:"shift,omitempty"`

	KaryawanID uint     `gorm:"column:id_karyawan;unique;not null" json:"id_karyawan"`
	Karyawan   Karyawan `gorm:"foreignKey:KaryawanID;references:IdKaryawan" json:"karyawan"`
}
