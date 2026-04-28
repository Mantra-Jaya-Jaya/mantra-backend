package models

func (StatusPengantaran) TableName() string {
  return "status_pengantaran"
}

type StatusPengantaran struct {
  IdStatusPengantaran      uint    `gorm:"primaryKey;column:id_status_pengantaran"`
  NamaStatus    		   string  `gorm:"column:nama_status"`
}