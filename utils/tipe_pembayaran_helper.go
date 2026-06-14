package utils

import (
	"backend-mantra/config"
	"backend-mantra/models"
)

func GetTipePembayaranID(namaTipe string) uint {
	var tipe models.TipePembayaran
	result := config.DB.Where("nama_tipe = ?", namaTipe).Take(&tipe)
	if result.Error != nil {
		panic("Tipe pembayaran '" + namaTipe + "' tidak ditemukan di database")
	}
	return tipe.IdTipePembayaran
}

func GetTipePembayaranIDSafe(namaTipe string) uint {
	var tipe models.TipePembayaran
	result := config.DB.Where("nama_tipe = ?", namaTipe).Take(&tipe)
	if result.Error != nil {
		return 0
	}
	return tipe.IdTipePembayaran
}
