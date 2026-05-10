package models

func (DetailPesanan) TableName() string {
	return "detail_pesanan"
}

type DetailPesanan struct {
	IdDetailPesanan uint `gorm:"primaryKey;column:id_detail_pesanan"`
	Jumlah          int  `gorm:"column:jumlah"`
	HargaSatuan     int  `gorm:"column:harga_satuan"` // Snapshot harga saat transaksi
	Subtotal        int  `gorm:"column:subtotal"`

	// Foreign Key ke Pesanan
	PesananId uint    `gorm:"column:id_pesanan"`
	Pesanan   Pesanan `gorm:"foreignKey:PesananId;references:IdPesanan"`

	// Foreign Key ke SpesifikasiBarang (bukan ke Barang langsung — menyimpan varian yang dipilih)
	SpesifikasiBarangId uint              `gorm:"column:id_spesifikasi_barang"`
	SpesifikasiBarang   SpesifikasiBarang `gorm:"foreignKey:SpesifikasiBarangId;references:IdSpesifikasiBarang"`
}