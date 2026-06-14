package utils

import (
	"backend-mantra/config"
	"backend-mantra/models"
)

func GetStatusKaryawanID(namaStatus string) uint {
	var status models.StatusKaryawan
	result := config.DB.Where("nama_status = ?", namaStatus).Take(&status)
	if result.Error != nil {
		panic("Status karyawan '" + namaStatus + "' tidak ditemukan di database")
	}
	return status.IdStatusKaryawan
}

func GetStatusKaryawanIDSafe(namaStatus string) uint {
	var status models.StatusKaryawan
	result := config.DB.Where("nama_status = ?", namaStatus).Take(&status)
	if result.Error != nil {
		return 0
	}
	return status.IdStatusKaryawan
}
