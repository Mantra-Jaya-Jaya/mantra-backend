package seeders

import (
	"fmt"
	"backend-mantra/config"
	"backend-mantra/models"

	"github.com/brianvoe/gofakeit/v7"
)

func SeedBarang() {
	// 1. Cek dulu jumlah barang yang ada di database
	var count int64
	config.DB.Model(&models.Barang{}).Count(&count)

	// Kalau barang udah ada 10 atau lebih, stop! Gak usah nambah lagi.
	if count >= 10 {
		fmt.Printf("Tabel barang udah punya %d data, proses seeding dilewati.\n", count)
		return
	}

	// 2. Ambil data master buat dipasangin ke barang
	var diskon models.Diskon
	var satuan models.Satuan
	var kategori models.Kategori

	// Kita ambil satu aja buat contoh
	config.DB.First(&diskon)
	config.DB.First(&satuan)
	config.DB.First(&kategori)

	// 3. Kita bikin sisa barangnya biar genap jadi 10
	sisaPerluDibuat := 10 - int(count)

	for i := 0; i < sisaPerluDibuat; i++ {
		barangBaru := models.Barang{
			NamaBarang:   gofakeit.ProductName(),    
			GambarBarang: fmt.Sprintf("https://picsum.photos/seed/%d/400/400", gofakeit.Number(1, 1000)),
			DiskonId:     diskon.IdDiskon,
			SatuanId:     satuan.IdSatuan,
			KategoriId:   kategori.IdKategori,
		}

		if err := config.DB.Create(&barangBaru).Error; err != nil {
			fmt.Println("Error:", err)
			continue
		}
	}

	fmt.Printf("Yeyy, Berhasil seed barang!")
}