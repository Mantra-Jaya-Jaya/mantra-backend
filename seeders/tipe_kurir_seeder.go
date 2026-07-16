package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"
)

func SeedTipeKurir() {
	fmt.Println("Seeding tipe kurir...")

	tipeList := []string{"internal", "external"}
	for _, nama := range tipeList {
		tipe := models.TipeKurir{NamaTipe: nama}
		if err := config.DB.Where("nama_tipe = ?", nama).FirstOrCreate(&tipe).Error; err != nil {
			fmt.Printf("Error seeding tipe kurir '%s': %v\n", nama, err)
		}
	}
	fmt.Println("Tipe kurir selesai di-seed!")
}
