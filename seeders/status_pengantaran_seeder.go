package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"
)

func SeedStatusPengantaran() {
	fmt.Println("⏳ Menyiapkan data status pengantaran...")

	daftarStatus := []string{
		"Menunggu Pickup",
		"Dalam Perjalanan",
		"Tiba di Tujuan",
		"Selesai",
		"Gagal Antar",
	}

	for _, nama := range daftarStatus {
		status := models.StatusPengantaran{NamaStatus: nama}
		if err := config.DB.Where("nama_status = ?", nama).FirstOrCreate(&status).Error; err != nil {
			fmt.Println("Error insert status pengantaran", nama, ":", err)
			continue
		}
	}

	fmt.Println("Yeyy, berhasil seed status pengantaran!")
}
