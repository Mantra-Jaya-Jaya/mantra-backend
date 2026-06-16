package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"
)

func SeedStatusPesanan() {
	fmt.Println("Seeding status pesanan...")

	statusList := []models.StatusPesanan{
		{NamaStatus: "Menunggu Pembayaran"},
		{NamaStatus: "Dikemas"},
		{NamaStatus: "Dikirim"},
		{NamaStatus: "Selesai"},
		{NamaStatus: "Dibatalkan"},
	}

	for _, status := range statusList {
		result := config.DB.Where("nama_status = ?", status.NamaStatus).FirstOrCreate(&status)
		if result.Error != nil {
			fmt.Printf("Error seeding status pesanan '%s': %v\n", status.NamaStatus, result.Error)
		} else if result.RowsAffected > 0 {
			fmt.Printf("Created status pesanan: %s\n", status.NamaStatus)
		}
	}

	fmt.Println("Status pesanan seeding completed")
}
