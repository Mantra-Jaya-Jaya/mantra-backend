package models

import (
	"github.com/google/uuid"
)

func (Barang) TableName() string {
	return "barang"
}

type Barang struct {
	IdBarang      uint      `gorm:"primaryKey;column:id_barang" json:"id_barang"`
	PublicId      uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();column:public_id;uniqueIndex" json:"public_id"`
	NamaBarang    string    `gorm:"column:nama_barang" json:"nama_barang"`
	GambarBarang  string    `gorm:"column:gambar_barang" json:"gambar_barang"`
	Deskripsi     string    `gorm:"column:deskripsi" json:"deskripsi"`
	PanjangBarang int       `gorm:"column:panjang_barang;default:0" json:"panjang_barang"`
	LebarBarang   int       `gorm:"column:lebar_barang;default:0" json:"lebar_barang"`
	TinggiBarang  int       `gorm:"column:tinggi_barang;default:0" json:"tinggi_barang"`

	//Relasi
	DiskonID   *uint `gorm:"column:id_diskon" json:"id_diskon"`
	SatuanID   uint  `gorm:"column:id_satuan;not null" json:"id_satuan"`
	KategoriID uint  `gorm:"column:id_kategori;not null" json:"id_kategori"`

	// Relasi
	Diskon   *Diskon  `gorm:"foreignKey:DiskonID;references:IdDiskon" json:"diskon"`
	Satuan   Satuan   `gorm:"foreignKey:SatuanID;references:IdSatuan" json:"satuan"`
	Kategori Kategori `gorm:"foreignKey:KategoriID;references:IdKategori" json:"kategori"`
}
