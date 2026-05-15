package models

import (
	"github.com/google/uuid"
)

// Memaksa nama tabel
func (Customer) TableName() string {
	return "customer"
}

type Customer struct {
	// Nama variabel WAJIB Kapital (IdRole), nama di DB diatur lewat Tag (id_role)
	IdCustomer uint      `gorm:"primaryKey;column:id_customer" json:"id_customer"`
	PublicId   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();column:public_id;uniqueIndex" json:"public_id"`
	NoTelp     string    `gorm:"column:no_telp" json:"no_telp"`

	//Relasi ke User
	UserId uint `gorm:"column:id_user;unique" json:"id_user"`

	//Relasi  ke tabel user
	User User `gorm:"foreignKey:UserId;references:IdUser" json:"user"`
}
