package models

func (Notifikasi) TableName() string {
	return "notifikasi"
}

type Notifikasi struct {
	IdNotifikasi uint   `gorm:"primaryKey;column:id_notifikasi"`
	
	UserID       uint   `gorm:"column:id_user"`
	User         User   `gorm:"foreignKey:UserID;references:IdUser"`
	
	Judul        string `gorm:"column:judul"`
	Pesan        string `gorm:"column:pesan"`
	Status       string `gorm:"column:status"`
}
