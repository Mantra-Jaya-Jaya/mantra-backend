package models

// Memaksa nama tabel
func (Barcode) TableName() string {
	return "barcode"
}

type Barcode struct {
	// Nama variabel WAJIB Kapital (IdRole), nama di DB diatur lewat Tag (id_role)
	IdBarcode uint `gorm:"primaryKey;column:id_barcode"`
	Kuantitas uint `gorm:"column:kuantitas"`

	//Relasi ke Barang
	BarangId uint `gorm:"column:id_barang;unique"`
	SatuanId uint `gorm:"column:id_satuan;unique"`

	//Relasi  ke tabel barang
	Barang Barang `gorm:"foreignKey:BarangId;references:IdBarang"`
}
