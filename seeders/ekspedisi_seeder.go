package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"
)

func SeedEkspedisi() {
	fmt.Println("⏳ Menyiapkan data ekspedisi...")

	daftarEkspedisi := []models.Ekspedisi{
		{NamaEkspedisi: "SPEX Express", KodeApi: "spex"},
		{NamaEkspedisi: "JNE", KodeApi: "jne"},
		{NamaEkspedisi: "J&T Express", KodeApi: "jnt"},
		{NamaEkspedisi: "SiCepat", KodeApi: "sicepat"},
		{NamaEkspedisi: "Anteraja", KodeApi: "anteraja"},
	}

	for _, eks := range daftarEkspedisi {
		if err := config.DB.Where("kode_api = ?", eks.KodeApi).FirstOrCreate(&eks).Error; err != nil {
			fmt.Println("Error insert ekspedisi", eks.NamaEkspedisi, ":", err)
			continue
		}
	}

	fmt.Println("Yeyy, berhasil seed ekspedisi!")
}
