package models

func (DetailPesanan) TableName() string {
  return "detail_pesanan"
}

type DetailPesanan struct {
  IdDetailPesanan uint 		`gorm:"primaryKey;column:id_detail_pesanan" json:"id_detail_pesanan"`
  Jumlah          int  		`gorm:"column:jumlah" json:"jumlah"`
  HargaSatuan     int  		`gorm:"column:harga_satuan" json:"harga_satuan"`
  Subtotal        int  		`gorm:"column:subtotal" json:"subtotal"`

  // Foreign Key ke Pesanan
  PesananId       uint 		`gorm:"column:id_pesanan" json:"id_pesanan"`
  Pesanan         Pesanan `gorm:"foreignKey:PesananId;references:IdPesanan" json:"pesanan"`
  
  // Foreign Key ke Barang
  BarangId        uint 		`gorm:"column:id_barang" json:"id_barang"`
  Barang          Barang 	`gorm:"foreignKey:BarangId;references:IdBarang" json:"barang"`
}