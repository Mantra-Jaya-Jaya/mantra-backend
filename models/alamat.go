package models

import "time"

func (Alamat) TableName() string {
	return "alamat"
}

type Alamat struct {
	IdAlamat       uint      `gorm:"primaryKey;column:id_alamat"`
	NamaPenerima   string    `gorm:"column:nama_penerima"`
	LabelAlamat    string    `gorm:"column:label_alamat"` // Rumah | Kantor | Kos | Lainnya
	NoTelpPenerima string    `gorm:"column:no_telp_penerima"`
	AlamatLengkap  string    `gorm:"type:text;column:alamat_lengkap"`
	Latitude       float64   `gorm:"column:latitude"`
	Longitude      float64   `gorm:"column:longitude"`
	CatatanLokasi  string    `gorm:"type:text;column:catatan_lokasi"`
	IsUtama        bool      `gorm:"column:is_utama"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`

	// Relasi ke tabel customer
	CustomerId uint     `gorm:"column:id_customer"`
	Customer   Customer `gorm:"foreignKey:CustomerId;references:IdCustomer"`
}

