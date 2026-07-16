package models

func (StatusKaryawan) TableName() string {
	return "status_karyawan"
}

type StatusKaryawan struct {
	IdStatusKaryawan uint   `gorm:"primaryKey;column:id" json:"id_status_karyawan"`
	NamaStatus       string `gorm:"column:nama_status;unique" json:"nama_status"`
}
