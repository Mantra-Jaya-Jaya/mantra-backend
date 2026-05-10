package models

func (User) TableName() string {
  return "user"
}

type User struct {
  IdUser      uint   `gorm:"primaryKey;column:id_user" json:"id_user"`
  Username    string `gorm:"column:username" json:"username"`
  Email       string `gorm:"unique;column:email" json:"email"`
  Password    string `gorm:"column:password" json:"password"`
  NamaLengkap string `gorm:"column:nama_lengkap" json:"nama_lengkap"`

  RoleID      uint   `gorm:"column:id_role" json:"id_role"`
  Role        Role   `gorm:"foreignKey:RoleID;references:IdRole" json:"role"`
}