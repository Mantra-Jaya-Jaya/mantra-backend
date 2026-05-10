package models

import "time"

// Memaksa nama tabel
func (Customer) TableName() string {
	return "customer"
}

type Customer struct {
	// Nama variabel WAJIB Kapital (IdRole), nama di DB diatur lewat Tag (id_role)
	IdCustomer uint      `gorm:"primaryKey;column:id_customer"`
	NoTelp     string    `gorm:"column:no_telp"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`

	// Relasi ke User
	UserId uint `gorm:"column:id_user;unique"`
	User   User `gorm:"foreignKey:UserId;references:IdUser"`
}
