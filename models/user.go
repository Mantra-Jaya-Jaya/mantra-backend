package models

import "time"

func (User) TableName() string {
	return "user"
}

type User struct {
	IdUser      uint      `gorm:"primaryKey;column:id_user"`
	Username    string    `gorm:"column:username"`
	Email       string    `gorm:"unique;column:email"`
	Password    string    `gorm:"column:password"`
	NamaLengkap string    `gorm:"column:nama_lengkap"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`

	// Relasi ke Role
	RoleID uint `gorm:"column:id_role"`
	Role   Role `gorm:"foreignKey:RoleID;references:IdRole"`
}