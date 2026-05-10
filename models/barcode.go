package models

// Memaksa nama tabel
func (Barcode) TableName() string {
	return "barcode"
}

type Barcode struct {
	IdBarcode   uint   `gorm:"primaryKey;column:id_barcode"`
	KodeBarcode string `gorm:"column:kode_barcode"` // String bukan int agar support leading zero
	Kuantitas   int    `gorm:"column:kuantitas"`    // Jumlah per scan, misal 1 dus = 12 pcs

	// Relasi ke Barang
	BarangId uint   `gorm:"column:id_barang"`
	Barang   Barang `gorm:"foreignKey:BarangId;references:IdBarang"`

	// Relasi ke Satuan
	SatuanId uint   `gorm:"column:id_satuan"`
	Satuan   Satuan `gorm:"foreignKey:SatuanId;references:IdSatuan"`
}

