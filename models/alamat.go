package models

func (Alamat) TableName() string {
	return "alamat"
}

type Alamat struct {
	IdAlamat       uint    `gorm:"primaryKey;column:id_alamat"`
	CustomerId     uint    `gorm:"column:id_customer"`
	NamaPenerima   string  `gorm:"column:nama_penerima"`
	LabelAlamat    string  `gorm:"column:label_alamat"`
	NoTelpPenerima int     `gorm:"column:no_telp_penerima"`
	AlamatLengkap  string  `gorm:"column:alamat_lengkap"`
	Latitude       float64 `gorm:"column:latitude"`
	Longitude      float64 `gorm:"column:longitude"`
	CatatanLokasi  string  `gorm:"column:catatan_lokasi"`
	IsUtama        bool    `gorm:"column:is_utama"`

	// Relasi ke tabel customer
	Customer Customer `gorm:"foreignKey:CustomerId;references:IdCustomer"`

}