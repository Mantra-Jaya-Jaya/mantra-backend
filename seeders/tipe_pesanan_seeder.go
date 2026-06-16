package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"
)

func SeedTipePesanan() {
	fmt.Println("Seeding tipe pesanan...")

	tipeList := []string{"Online", "Offline"}
	for _, nama := range tipeList {
		tipe := models.TipePesanan{NamaTipe: nama}
		if err := config.DB.Where("nama_tipe = ?", nama).FirstOrCreate(&tipe).Error; err != nil {
			fmt.Printf("Error seeding tipe pesanan '%s': %v\n", nama, err)
		}
	}
	fmt.Println("Tipe pesanan selesai di-seed!")
}
