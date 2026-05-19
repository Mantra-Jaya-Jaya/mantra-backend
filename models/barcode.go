package models

// Memaksa nama tabel
func (Barcode) TableName() string {
  return "barcode"
}

type Barcode struct {
  IdBarcode           uint              `gorm:"primaryKey;column:id_barcode" json:"id_barcode"`
  KodeBarcode         string            `gorm:"column:kode_barcode;type:varchar" json:"kode_barcode"` 
  Kuantitas           uint              `gorm:"column:kuantitas" json:"kuantitas"`
  
  SpesifikasiBarangID uint              `gorm:"column:id_spesifikasi_barang"`
  SpesifikasiBarang   SpesifikasiBarang `gorm:"foreignKey:SpesifikasiBarangID;references:IdSpesifikasiBarang"`
}