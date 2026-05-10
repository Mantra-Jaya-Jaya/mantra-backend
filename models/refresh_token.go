package models

import "time"

func (RefreshToken) TableName() string {
	return "refresh_token"
}

type RefreshToken struct {
	IdToken   uint       `gorm:"primaryKey;column:id_token"`
	Token     string     `gorm:"type:text;column:token"`
	ExpiresAt time.Time  `gorm:"column:expires_at"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	RevokedAt *time.Time `gorm:"column:revoked_at"` // Nullable — diisi saat logout, null berarti masih aktif

	// Relasi ke User
	UserID uint `gorm:"column:id_user"`
	User   User `gorm:"foreignKey:UserID;references:IdUser"`
}
