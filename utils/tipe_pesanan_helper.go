package utils

import (
	"backend-mantra/config"
	"backend-mantra/models"
)

func GetTipePesananID(namaTipe string) uint {
	var tipe models.TipePesanan
	result := config.DB.Where("nama_tipe = ?", namaTipe).Take(&tipe)
	if result.Error != nil {
		panic("Tipe pesanan '" + namaTipe + "' tidak ditemukan di database")
	}
	return tipe.IdTipePesanan
}

func GetTipePesananIDSafe(namaTipe string) uint {
	var tipe models.TipePesanan
	result := config.DB.Where("nama_tipe = ?", namaTipe).Take(&tipe)
	if result.Error != nil {
		return 0
	}
	return tipe.IdTipePesanan
}
