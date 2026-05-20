package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"
	"net/url"
	"time"
)

func SeedDiskon() {
	fmt.Println("⏳ Menyiapkan data diskon...")

	now := time.Now()

	diskons := []models.Diskon{
		{
			NamaDiskon:   "Flash Sale Elektronik",
			BesarDiskon:  15,
			BannerDiskon: fmt.Sprintf("https://placehold.co/800x300?text=%s", url.QueryEscape("Flash Sale Elektronik")),
			TglMulai:     now,
			TglSelesai:   now.AddDate(0, 0, 30),
		},
		{
			NamaDiskon:   "Promo Akhir Bulan",
			BesarDiskon:  20,
			BannerDiskon: fmt.Sprintf("https://placehold.co/800x300?text=%s", url.QueryEscape("Promo Akhir Bulan")),
			TglMulai:     now,
			TglSelesai:   now.AddDate(0, 0, 14),
		},
		{
			NamaDiskon:   "Diskon Member Baru",
			BesarDiskon:  10,
			BannerDiskon: fmt.Sprintf("https://placehold.co/800x300?text=%s", url.QueryEscape("Diskon Member Baru")),
			TglMulai:     now,
			TglSelesai:   now.AddDate(0, 0, 60),
		},
	}

	for _, d := range diskons {
		if err := config.DB.Where("nama_diskon = ?", d.NamaDiskon).FirstOrCreate(&d).Error; err != nil {
			fmt.Println("Error insert data diskon", d.NamaDiskon, ":", err)
			continue
		}
	}

	fmt.Println("Yeyy, berhasil seed diskon!")
}
