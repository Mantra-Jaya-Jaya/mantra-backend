package utils

import (
	"backend-mantra/config"
	"backend-mantra/models"
)

func GetFraudStatusID(namaStatus string) uint {
	var status models.FraudStatus
	result := config.DB.Where("nama_status = ?", namaStatus).Take(&status)
	if result.Error != nil {
		panic("Fraud status '" + namaStatus + "' tidak ditemukan di database")
	}
	return status.IdFraudStatus
}

func GetFraudStatusIDSafe(namaStatus string) uint {
	var status models.FraudStatus
	result := config.DB.Where("nama_status = ?", namaStatus).Take(&status)
	if result.Error != nil {
		return 0
	}
	return status.IdFraudStatus
}
