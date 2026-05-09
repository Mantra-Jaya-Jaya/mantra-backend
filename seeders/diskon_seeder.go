package seeders

import (
	"fmt"
	"time"
	"backend-mantra/config"
	"backend-mantra/models"
)

func SeedDiskon() {
	// buat 3 data diskon 
	daftarDiskon := []models.Diskon{
		{
			NamaDiskon:   "Promo Back to Campus",
			BesarDiskon:  15, 
			BannerDiskon: "https://cataas.com/cat?text=Promo+Kampus&width=400&height=200", 

			TglMulai:     time.Date(2026, time.May, 1, 0, 0, 0, 0, time.Local),
			TglSelesai:   time.Date(2026, time.May, 31, 0, 0, 0, 0, time.Local),
		},
		{
			NamaDiskon:   "Flash Sale Mahasiswa Polines",
			BesarDiskon:  50, // Diskon gila-gilaan 50%
			BannerDiskon: "https://cataas.com/cat?text=Flash+Sale&width=400&height=200",
			// Promo super singkat (misal: 10 - 15 Mei 2026)
			TglMulai:     time.Date(2026, time.May, 10, 0, 0, 0, 0, time.Local),
			TglSelesai:   time.Date(2026, time.May, 15, 0, 0, 0, 0, time.Local),
		},
		{
			NamaDiskon:   "Diskon Akhir Tahun",
			BesarDiskon:  25,
			BannerDiskon: "https://cataas.com/cat?text=Akhir+Tahun&width=400&height=200",
			// Promo buat akhir tahun nanti
			TglMulai:     time.Date(2026, time.December, 1, 0, 0, 0, 0, time.Local),
			TglSelesai:   time.Date(2026, time.December, 31, 0, 0, 0, 0, time.Local),
		},
	}

	for _, diskon := range daftarDiskon {
		if err := config.DB.Where("nama_diskon = ?", diskon.NamaDiskon).FirstOrCreate(&diskon).Error; err != nil {
			fmt.Println("Error insert data diskon", diskon.NamaDiskon, ":", err)
			continue
		}
	}

	fmt.Println("Yeyy, berhasil seed diskon!")
}