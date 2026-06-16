package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"
)

func SeedShiftKasir() {
	fmt.Println("Seeding shift kasir...")

	shiftList := []string{"Pagi", "Siang", "Malam"}
	for _, nama := range shiftList {
		shift := models.ShiftKasir{NamaShift: nama}
		if err := config.DB.Where("nama_shift = ?", nama).FirstOrCreate(&shift).Error; err != nil {
			fmt.Printf("Error seeding shift kasir '%s': %v\n", nama, err)
		}
	}
	fmt.Println("Shift kasir selesai di-seed!")
}
