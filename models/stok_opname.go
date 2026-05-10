package models

import "time"

func (StokOpname) TableName() string {
	return "stok_opname"
}

// Riwayat pergerakan stok per varian barang
type StokOpname struct {
	IdStokOpname uint      `gorm:"primaryKey;column:id_stok_opname" json:"id_stok_opname"`
	HargaBeli    int       `gorm:"column:harga_beli" json:"harga_beli"` // Harga beli/modal dalam Rupiah
	Status       bool      `gorm:"column:status" json:"status"`     // true = stok masuk, false = stok keluar
	JumlahStok   int       `gorm:"column:jumlah_stok" json:"jumlah_stok"` // yang sedang di masukkan atau keluar (bukan total stok)
	Keterangan   string    `gorm:"type:text;column:keterangan" json:"keterangan"`
	Tanggal      time.Time `gorm:"column:tanggal" json:"tanggal"`

	// Relasi ke SpesifikasiBarang
	SpesifikasiBarangID uint              `gorm:"column:id_spesifikasi_barang" json:"id_spesifikasi_barang"`
	SpesifikasiBarang   SpesifikasiBarang `gorm:"foreignKey:SpesifikasiBarangID;references:IdSpesifikasiBarang" json:"spesifikasi_barang"`
}
