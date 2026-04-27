package models

func (User) TableName() string {
  return "user"
}

type User struct {
  IdUser      uint   `gorm:"primaryKey;column:id_user"`
  Username    string `gorm:"column:username"`
  Email       string `gorm:"unique;column:email"`
  Password    string `gorm:"column:password"`
  NamaLengkap string `gorm:"column:nama_lengkap"`
  
  RoleID      uint   `gorm:"column:id_role"`  
  Role        Role   `gorm:"foreignKey:RoleID;references:IdRole"` 
}