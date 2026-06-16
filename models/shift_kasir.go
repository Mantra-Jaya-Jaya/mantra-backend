package models

func (ShiftKasir) TableName() string {
	return "shift_kasir"
}

type ShiftKasir struct {
	IdShiftKasir uint   `gorm:"primaryKey;column:id" json:"id_shift_kasir"`
	NamaShift    string `gorm:"column:nama_shift;unique" json:"nama_shift"`
}
