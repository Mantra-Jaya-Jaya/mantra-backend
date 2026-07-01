package utils

import (
	"backend-mantra/config"
	"backend-mantra/models"
)

// GetStatusPengantaranID mencari ID status pengantaran berdasarkan nama.
// Panic jika nama tidak ditemukan — pastikan nama sesuai dengan seeder.
func GetStatusPengantaranID(namaStatus string) uint {
	var status models.StatusPengantaran
	result := config.DB.Where("nama_status = ?", namaStatus).Take(&status)
	if result.Error != nil {
		panic("Status pengantaran '" + namaStatus + "' tidak ditemukan di database")
	}
	return status.IdStatusPengantaran
}

// GetStatusPengantaranIDSafe sama seperti GetStatusPengantaranID tapi tidak panic.
// Return 0 jika tidak ditemukan.
func GetStatusPengantaranIDSafe(namaStatus string) uint {
	var status models.StatusPengantaran
	result := config.DB.Where("nama_status = ?", namaStatus).Take(&status)
	if result.Error != nil {
		return 0
	}
	return status.IdStatusPengantaran
}
