package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"
)

func SeedDetailSpesifikasi() {
	fmt.Println("⏳ Menyiapkan data Detail Spesifikasi...")

	// Mapping spesifikasi → values
	specValues := map[string][]string{
		"Warna":     {"Merah", "Kuning", "Hijau", "Biru", "Hitam", "Putih"},
		"Ukuran":    {"S", "M", "L", "XL", "XXL"},
		"RAM":       {"4GB", "8GB", "16GB", "32GB"},
		"Kapasitas": {"64GB", "128GB", "256GB", "512GB", "1TB"},
		"Rasa":      {"Original", "Manis", "Asam", "Pedas"},
		"Bahan":     {"Kain", "Katun", "Polyester", "Denim"},
	}

	totalCreated := 0
	for specName, values := range specValues {
		var spec models.Spesifikasi
		if err := config.DB.Where("nama_spesifikasi = ?", specName).First(&spec).Error; err != nil {
			fmt.Println("Spesifikasi", specName, "gak ketemu! Pastiin SeedSpesifikasi jalan duluan.")
			continue
		}

		for _, value := range values {
			detail := models.DetailSpesifikasi{
				NamaDetailSpesifikasi: value,
				SpesifikasiID:         spec.IdSpesifikasi,
			}
			if err := config.DB.Where("nama_detail_spesifikasi = ? AND id_spesifikasi = ?", value, spec.IdSpesifikasi).FirstOrCreate(&detail).Error; err == nil {
				totalCreated++
			}
		}
	}

	fmt.Printf("Yeyy, Berhasil seed %d detail spesifikasi!\n", totalCreated)
}
