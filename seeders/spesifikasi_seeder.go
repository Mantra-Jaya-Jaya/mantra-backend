package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"
)

func SeedSpesifikasi() {
	// Siapkan Datanya
	daftarSpesifikasi := []models.Spesifikasi{
		{NamaSpesifikasi: "Warna"},
		{NamaSpesifikasi: "Ukuran"},
	}

	// Looping untuk masukin ke database
	for _, spek := range daftarSpesifikasi {
		// FirstOrCreate
		if err := config.DB.Where("nama_spesifikasi = ?", spek.NamaSpesifikasi).FirstOrCreate(&spek).Error; err != nil {
			fmt.Println("Error", spek.NamaSpesifikasi, ":", err)
			continue
		}
	}

	fmt.Println("Yeyy, Berhasil seed spesifikasi!")
}
