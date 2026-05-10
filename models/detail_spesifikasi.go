package models

func (DetailSpesifikasi) TableName() string {
	return "detail_spesifikasi"
}

type DetailSpesifikasi struct {
	IdDetailSpesifikasi   uint   `gorm:"primaryKey;column:id_detail_spesifikasi" json:"id_detail_spesifikasi"`
	NamaDetailSpesifikasi string `gorm:"type:text;column:nama_detail_spesifikasi" json:"nama_detail_spesifikasi"` // Hitam, Merah, XL, 16GB, dll

	// Relasi ke master tipe spesifikasi (Warna, Ukuran, RAM, dll)
	SpesifikasiID uint        `gorm:"column:id_spesifikasi" json:"spesifikasi_id"`
	Spesifikasi   Spesifikasi `gorm:"foreignKey:SpesifikasiID;references:IdSpesifikasi" json:"spesifikasi"`
}

