package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"
)

func SeedStatusTransaksi() {
	fmt.Println("Seeding status transaksi...")

	statusList := []string{"pending", "settlement", "deny", "cancel", "expire"}
	for _, nama := range statusList {
		status := models.StatusTransaksi{NamaStatus: nama}
		if err := config.DB.Where("nama_status = ?", nama).FirstOrCreate(&status).Error; err != nil {
			fmt.Printf("Error seeding status transaksi '%s': %v\n", nama, err)
		}
	}
	fmt.Println("Status transaksi selesai di-seed!")
}
