package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"
)

func SeedPengaturanToko() {
	fmt.Println("Seeding pengaturan toko...")

	settings := map[string]string{
		"radius_kurir_internal": "5",
	}
	for key, value := range settings {
		setting := models.PengaturanToko{Key: key, Value: value}
		if err := config.DB.Where("key = ?", key).FirstOrCreate(&setting).Error; err != nil {
			fmt.Printf("Error seeding pengaturan '%s': %v\n", key, err)
		}
	}
	fmt.Println("Pengaturan toko selesai di-seed!")
}
