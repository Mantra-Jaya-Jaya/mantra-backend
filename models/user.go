package models

import (
	"github.com/google/uuid"
)

func (User) TableName() string {
	return "user"
}

type User struct {
	IdUser      uint      `gorm:"primaryKey;column:id_user" json:"id_user"`
	PublicId    uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();column:public_id;uniqueIndex" json:"public_id"`
	Username    string    `gorm:"column:username" json:"username"`
	Email       string    `gorm:"unique;column:email" json:"email"`
	Password    string    `gorm:"column:password" json:"-"`
	NamaLengkap string    `gorm:"column:nama_lengkap" json:"nama_lengkap"`
	FotoProfil  string    `gorm:"column:foto_profil" json:"foto_profil"`

	RoleID uint `gorm:"column:id_role" json:"id_role"`
	Role   Role `gorm:"foreignKey:RoleID;references:IdRole" json:"role"`
}
