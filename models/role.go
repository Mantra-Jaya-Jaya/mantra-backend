package models

// Memaksa nama tabel
func (Role) TableName() string {
  return "role"
}

type Role struct {
  // Nama variabel WAJIB Kapital (IdRole), nama di DB diatur lewat Tag (id_role)
  IdRole   uint   `gorm:"primaryKey;column:id_role" json:"id_role"`
  NamaRole string `gorm:"unique;column:nama_role" json:"nama_role"`
}
