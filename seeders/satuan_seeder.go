package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"
)

func SeedSatuan() {
	// Siapkan Datanya
	daftarSatuan := []models.Satuan{
		{NamaSatuan: "Pcs"},
		{NamaSatuan: "Kg"},
		{NamaSatuan: "Gram"},
	}

	// Looping untuk masukin ke database
	for _, satuan := range daftarSatuan {
		// Mantra FirstOrCreate
		if err := config.DB.Where("nama_satuan = ?", satuan.NamaSatuan).FirstOrCreate(&satuan).Error; err != nil {
			fmt.Println("Error", satuan.NamaSatuan, ":", err)
			continue
		}
	}

	fmt.Println("Yeyy, Berhasil seed satuan!")
}
