package utils

import (
	"backend-mantra/config"
	"backend-mantra/models"
)

func GetShiftKasirID(namaShift string) uint {
	var shift models.ShiftKasir
	result := config.DB.Where("nama_shift = ?", namaShift).Take(&shift)
	if result.Error != nil {
		panic("Shift kasir '" + namaShift + "' tidak ditemukan di database")
	}
	return shift.IdShiftKasir
}

func GetShiftKasirIDSafe(namaShift string) uint {
	var shift models.ShiftKasir
	result := config.DB.Where("nama_shift = ?", namaShift).Take(&shift)
	if result.Error != nil {
		return 0
	}
	return shift.IdShiftKasir
}
