package utils

import (
	"backend-mantra/config"
	"backend-mantra/models"
)

// GetStatusTransaksiID mencari ID status transaksi berdasarkan nama status.
// Panic jika tidak ditemukan.
func GetStatusTransaksiID(namaStatus string) uint {
	var status models.StatusTransaksi
	result := config.DB.Where("nama_status = ?", namaStatus).Take(&status)
	if result.Error != nil {
		panic("Status transaksi '" + namaStatus + "' tidak ditemukan di database")
	}
	return status.IdStatusTransaksi
}

// GetStatusTransaksiIDSafe mencari ID status transaksi, return 0 jika tidak ditemukan.
func GetStatusTransaksiIDSafe(namaStatus string) uint {
	var status models.StatusTransaksi
	result := config.DB.Where("nama_status = ?", namaStatus).Take(&status)
	if result.Error != nil {
		return 0
	}
	return status.IdStatusTransaksi
}
