package models

import (
	"github.com/google/uuid"
)

func (Keranjang) TableName() string {
	return "keranjang"
}

// Keranjang belanja customer. FK ke spesifikasi_barang agar stok dan varian sinkron
type Keranjang struct {
	IdKeranjang uint      `gorm:"primaryKey;column:id_keranjang" json:"id_keranjang"`
	PublicId    uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();column:public_id;uniqueIndex" json:"public_id"`
	Quantity    int       `gorm:"column:quantity" json:"quantity"`

	// Relasi ke Customer
	CustomerID uint     `gorm:"column:id_customer"`
	Customer   Customer `gorm:"foreignKey:CustomerID;references:IdCustomer"`

	// Relasi ke SpesifikasiBarang (varian spesifik yang dipilih)
	SpesifikasiBarangID uint              `gorm:"column:id_spesifikasi_barang"`
	SpesifikasiBarang   SpesifikasiBarang `gorm:"foreignKey:SpesifikasiBarangID;references:IdSpesifikasiBarang"`
}
