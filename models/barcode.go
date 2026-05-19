package models

// Memaksa nama tabel
func (Barcode) TableName() string {
  return "barcode"
}

type Barcode struct {
  // Nama variabel WAJIB Kapital (IdRole), nama di DB diatur lewat Tag (id_role)
  IdBarcode           uint `gorm:"primaryKey;column:id_barcode" json:"id_barcode"`
  Kuantitas           uint `gorm:"column:kuantitas" json:"kuantitas"`
  SpesifikasiBarangId uint `gorm:"column:id_spesifikasi_barang" json:"id_spesifikasi_barang"`
  
  // Relasi ke Satuan (Diambil dari branch feat/flutter-auth-integration)
  SatuanId            uint `gorm:"column:id_satuan" json:"id_satuan"`

  // Relasi ke Barang/SpesifikasiBarang (Mengikuti penamaan terbaru di branch dev)
  SpesifikasiBarang SpesifikasiBarang `gorm:"foreignKey:SpesifikasiBarangId;references:IdSpesifikasiBarang" json:"spesifikasi_barang"`
}
