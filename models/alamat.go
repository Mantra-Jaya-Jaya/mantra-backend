package models

import "time"

func (Alamat) TableName() string {
	return "alamat"
}

type Alamat struct {
	IdAlamat       uint    `gorm:"primaryKey;column:id_alamat" json:"id_alamat"`
	CustomerId     uint    `gorm:"column:id_customer" json:"id_customer"`
	NamaPenerima   string  `gorm:"column:nama_penerima" json:"nama_penerima"`
	LabelAlamat    string  `gorm:"column:label_alamat" json:"label_alamat"`
	NoTelpPenerima string  `gorm:"column:no_telp_penerima" json:"no_telp_penerima"`
	AlamatLengkap  string  `gorm:"column:alamat_lengkap" json:"alamat_lengkap"`
	Latitude       float64 `gorm:"column:latitude" json:"latitude"`
	Longitude      float64 `gorm:"column:longitude" json:"longitude"`
	CatatanLokasi  string  `gorm:"column:catatan_lokasi" json:"catatan_lokasi"`
	IsUtama        bool    `gorm:"column:is_utama" json:"is_utama"`

	// Relasi ke tabel customer
	Customer Customer `gorm:"foreignKey:CustomerId;references:IdCustomer" json:"customer"`
}
