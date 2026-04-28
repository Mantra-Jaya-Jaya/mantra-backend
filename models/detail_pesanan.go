package models

func (DetailPesanan) TableName() string {
  return "detail_pesanan"
}

type DetailPesanan struct {
  IdDetailPesanan uint 		`gorm:"primaryKey;column:id_detail_pesanan"`
  Jumlah          int  		`gorm:"column:jumlah"`
  HargaSatuan     int  		`gorm:"column:harga_satuan"`
  Subtotal        int  		`gorm:"column:subtotal"`

  // Foreign Key ke Pesanan
  PesananId       uint 		`gorm:"column:id_pesanan"`
  Pesanan         Pesanan `gorm:"foreignKey:PesananId;references:IdPesanan"`
  
  // Foreign Key ke Barang
  BarangId        uint 		`gorm:"column:id_barang"`
  Barang          Barang 	`gorm:"foreignKey:BarangId;references:IdBarang"`
}