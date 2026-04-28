package models

func (DetailSpesifikasi) TableName() string {
	return "detail_spesifikasi"
}

type DetailSpesifikasi struct {
	IdDetailSpesifikasi   uint   `gorm:"primaryKey;column:id_detail_spesifikasi"`
	NamaDetailSpesifikasi string `gorm:"column:nama_detail_spesifikasi"`

	BarangID uint   `gorm:"column:id_barang"`
	Barang   Barang `gorm:"foreignKey:BarangID;references:IdBarang"`

	SpesifikasiID uint        `gorm:"column:id_spesifikasi"`
	Spesifikasi   Spesifikasi `gorm:"foreignKey:SpesifikasiID;references:IdSpesifikasi"`
}
