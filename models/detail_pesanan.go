package models

import "gorm.io/gorm"

func (DetailPesanan) TableName() string {
	return "detail_pesanan"
}

type DetailPesanan struct {
	IdDetailPesanan uint `gorm:"primaryKey;column:id_detail_pesanan" json:"id_detail_pesanan"`
	Jumlah          int  `gorm:"column:jumlah" json:"jumlah"`
	HargaSatuan     int  `gorm:"column:harga_satuan" json:"harga_satuan"` // Snapshot harga saat transaksi
	Subtotal        int  `gorm:"column:subtotal" json:"subtotal"`

	// Foreign Key ke Pesanan
	PesananID uint    `gorm:"column:id_pesanan;not null" json:"id_pesanan"`
	Pesanan   Pesanan `gorm:"foreignKey:PesananID;references:IdPesanan" json:"pesanan"`

	// Foreign Key ke SpesifikasiBarang (bukan ke Barang langsung — menyimpan varian yang dipilih)
	SpesifikasiBarangID uint              `gorm:"column:id_spesifikasi_barang;not null" json:"id_spesifikasi_barang"`
	SpesifikasiBarang   SpesifikasiBarang `gorm:"foreignKey:SpesifikasiBarangID;references:IdSpesifikasiBarang" json:"spesifikasi_barang"`
}

// BeforeCreate hook untuk auto-compute Subtotal saat INSERT
func (d *DetailPesanan) BeforeCreate(tx *gorm.DB) error {
	d.Subtotal = d.Jumlah * d.HargaSatuan
	return nil
}

// BeforeUpdate hook untuk auto-compute Subtotal saat UPDATE
func (d *DetailPesanan) BeforeUpdate(tx *gorm.DB) error {
	d.Subtotal = d.Jumlah * d.HargaSatuan
	return nil
}
