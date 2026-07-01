package utils

import (
	"backend-mantra/config"
	"backend-mantra/models"
)

func GetStatusNotifikasiID(namaStatus string) uint {
	var status models.StatusNotifikasi
	result := config.DB.Where("nama_status = ?", namaStatus).Take(&status)
	if result.Error != nil {
		panic("Status notifikasi '" + namaStatus + "' tidak ditemukan di database")
	}
	return status.IdStatusNotifikasi
}

func GetStatusNotifikasiIDSafe(namaStatus string) uint {
	var status models.StatusNotifikasi
	result := config.DB.Where("nama_status = ?", namaStatus).Take(&status)
	if result.Error != nil {
		return 0
	}
	return status.IdStatusNotifikasi
}
