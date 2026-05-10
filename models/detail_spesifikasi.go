package models

func (DetailSpesifikasi) TableName() string {
	return "detail_spesifikasi"
}

type DetailSpesifikasi struct {
	IdDetailSpesifikasi   uint   `gorm:"primaryKey;column:id_detail_spesifikasi" json:"id_detail_spesifikasi"`
	NamaDetailSpesifikasi string `gorm:"column:nama_detail_spesifikasi" json:"nama_detail_spesifikasi"`

	BarangID uint   `gorm:"column:id_barang" json:"id_barang"`
	Barang   Barang `gorm:"foreignKey:BarangID;references:IdBarang" json:"barang"`

	SpesifikasiID uint        `gorm:"column:id_spesifikasi" json:"id_spesifikasi"`
	Spesifikasi   Spesifikasi `gorm:"foreignKey:SpesifikasiID;references:IdSpesifikasi" json:"spesifikasi"`
}
