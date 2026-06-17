package utils

import (
	"backend-mantra/config"
	"backend-mantra/models"
)

func GetTipeKurirID(namaTipe string) uint {
	var tipe models.TipeKurir
	result := config.DB.Where("nama_tipe = ?", namaTipe).Take(&tipe)
	if result.Error != nil {
		panic("Tipe kurir '" + namaTipe + "' tidak ditemukan di database")
	}
	return tipe.IdTipeKurir
}

func GetTipeKurirIDSafe(namaTipe string) uint {
	var tipe models.TipeKurir
	result := config.DB.Where("nama_tipe = ?", namaTipe).Take(&tipe)
	if result.Error != nil {
		return 0
	}
	return tipe.IdTipeKurir
}

func GetNamaTipeKurir(id uint) string {
	var tipe models.TipeKurir
	result := config.DB.Where("id_tipe_kurir = ?", id).Take(&tipe)
	if result.Error != nil {
		return ""
	}
	return tipe.NamaTipe
}
