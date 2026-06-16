package utils

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"
)

// GetStatusPesananID mencari ID status pesanan berdasarkan nama status
// Panic jika tidak ditemukan (karena status harus selalu ada di DB)
func GetStatusPesananID(namaStatus string) uint {
	var status models.StatusPesanan
	result := config.DB.Where("nama_status = ?", namaStatus).Take(&status)
	if result.Error != nil {
		panic(fmt.Sprintf("Status pesanan '%s' tidak ditemukan di database", namaStatus))
	}
	return status.IdStatusPesanan
}

// GetStatusPesananIDSafe mencari ID status pesanan, return 0 jika tidak ditemukan
func GetStatusPesananIDSafe(namaStatus string) uint {
	var status models.StatusPesanan
	result := config.DB.Where("nama_status = ?", namaStatus).Take(&status)
	if result.Error != nil {
		return 0
	}
	return status.IdStatusPesanan
}
