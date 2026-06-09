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
		{NamaEkspedisi: "Ninja Xpress", KodeApi: "ninja"},
	}

	for _, eks := range daftarEkspedisi {
		if err := config.DB.Where("kode_api = ?", eks.KodeApi).FirstOrCreate(&eks).Error; err != nil {
			fmt.Println("Error insert ekspedisi", eks.NamaEkspedisi, ":", err)
			continue
		}

		// Seed layanan default untuk setiap ekspedisi
		seedLayananEkspedisi(eks)
	}

	fmt.Println("Yeyy, berhasil seed ekspedisi!")
}

func seedLayananEkspedisi(eks models.Ekspedisi) {
	daftarLayanan := map[string][]models.EkspedisiLayanan{
		"jne": {
			{NamaLayanan: "REG", Deskripsi: "Reguler", EstimasiMin: 2, EstimasiMax: 3},
			{NamaLayanan: "YES", Deskripsi: "Yakin Esok Sampai", EstimasiMin: 1, EstimasiMax: 1},
			{NamaLayanan: "OKE", Deskripsi: "Ongkos Kirim Ekonomis", EstimasiMin: 3, EstimasiMax: 5},
		},
		"jnt": {
			{NamaLayanan: "EZ", Deskripsi: "Economy", EstimasiMin: 2, EstimasiMax: 4},
			{NamaLayanan: "REG", Deskripsi: "Reguler", EstimasiMin: 2, EstimasiMax: 3},
		},
		"sicepat": {
			{NamaLayanan: "REG", Deskripsi: "Reguler", EstimasiMin: 1, EstimasiMax: 2},
			{NamaLayanan: "BEST", Deskripsi: "Besok Sampai Tujuan", EstimasiMin: 1, EstimasiMax: 1},
		},
		"anteraja": {
			{NamaLayanan: "REG", Deskripsi: "Reguler", EstimasiMin: 2, EstimasiMax: 4},
			{NamaLayanan: "ND", Deskripsi: "Next Day", EstimasiMin: 1, EstimasiMax: 1},
		},
		"spex": {
			{NamaLayanan: "REG", Deskripsi: "Reguler", EstimasiMin: 2, EstimasiMax: 4},
		},
		"ninja": {
			{NamaLayanan: "REG", Deskripsi: "Reguler", EstimasiMin: 2, EstimasiMax: 4},
		},
	}

	if layanan, exists := daftarLayanan[eks.KodeApi]; exists {
		for _, l := range layanan {
			l.EkspedisiID = eks.IdEkspedisi
			if err := config.DB.Where("id_ekspedisi = ? AND nama_layanan = ?", eks.IdEkspedisi, l.NamaLayanan).FirstOrCreate(&l).Error; err != nil {
				fmt.Println("Error insert layanan", l.NamaLayanan, "untuk", eks.NamaEkspedisi, ":", err)
			}
		}
	}
}
