package models

// Memaksa nama tabel
func (Barcode) TableName() string {
	return "barcode"
}

type Barcode struct {
	// Nama variabel WAJIB Kapital (IdRole), nama di DB diatur lewat Tag (id_role)
	IdBarcode uint `gorm:"primaryKey;column:id_barcode" json:"id_barcode"`
	Kuantitas uint `gorm:"column:kuantitas" json:"kuantitas"`

	//Relasi ke Spesifikasi Barang
	SpesifikasiBarangId uint `gorm:"column:id_spesifikasi_barang;unique" json:"id_spesifikasi_barang"`
	SatuanId            uint `gorm:"column:id_satuan;unique" json:"id_satuan"`

	//Relasi ke tabel spesifikasi barang
	SpesifikasiBarang SpesifikasiBarang `gorm:"foreignKey:SpesifikasiBarangId;references:IdSpesifikasiBarang" json:"spesifikasi_barang"`
}
