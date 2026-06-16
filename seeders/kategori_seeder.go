package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"
	"net/url"
)

func SeedKategori() {
	fmt.Println("⏳ Menyiapkan data kategori...")

	kategoris := []string{
		"Elektronik",
		"Fashion",
		"Makanan & Minuman",
		"Kesehatan",
		"Olahraga",
		"Peralatan Rumah",
		"Buku & Alat Tulis",
		"Kecantikan",
	}

	for _, nama := range kategoris {
		kategori := models.Kategori{
			NamaKategori: nama,
			IconKategori: fmt.Sprintf("https://placehold.co/64x64?text=%s", url.QueryEscape(nama)),
		}

		if err := config.DB.Where("nama_kategori = ?", nama).FirstOrCreate(&kategori).Error; err != nil {
			fmt.Println("Hadehh error insert kategori:", err)
			continue
		}
	}

	fmt.Println("Yeyy, berhasil seed kategori!")
}
