package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"

	fake "github.com/brianvoe/gofakeit/v7"
)

func SeedSpesifikasiBarang() {
	fmt.Println("⏳ Menyiapkan stok dan harga tiap varian barang...")

	var count int64
	config.DB.Model(&models.SpesifikasiBarang{}).Count(&count)
	if count >= 30 {
		fmt.Println("Tabel spesifikasi_barang udah punya cukup data, proses seeding dilewati.")
		return
	}

	var daftarBarang []models.Barang
	if err := config.DB.Preload("Kategori").Find(&daftarBarang).Error; err != nil || len(daftarBarang) == 0 {
		fmt.Println("Barang masih kosong! Pastiin SeedBarang jalan duluan.")
		return
	}

	var daftarDetailSpek []models.DetailSpesifikasi
	if err := config.DB.Find(&daftarDetailSpek).Error; err != nil || len(daftarDetailSpek) == 0 {
		fmt.Println("Detail Spesifikasi kosong! Pastiin SeedDetailSpesifikasi jalan.")
		return
	}

	totalVariasiDibuat := 0
	for _, barang := range daftarBarang {
		baseMin, baseMax := 10000, 100000
		kat := barang.Kategori.NamaKategori
		
		if kat == "Elektronik" {
			baseMin, baseMax = 1500000, 15000000
		} else if kat == "Fashion" {
			baseMin, baseMax = 150000, 800000
		} else if kat == "Makanan & Minuman" {
			baseMin, baseMax = 5000, 200000
		} else if kat == "Kesehatan" {
			baseMin, baseMax = 15000, 150000
		} else if kat == "Olahraga" {
			baseMin, baseMax = 50000, 500000
		}

		for i := 0; i < 2; i++ {
			// Pilih ID spesifikasi yang unik per barang
			detailIdx := (i + int(barang.IdBarang)) % len(daftarDetailSpek)
			detail := daftarDetailSpek[detailIdx]

			// Buat beberapa barang memiliki stok menipis (<10) atau kritis (<=5)
			jumlahStok := fake.IntRange(15, 50)
			if totalVariasiDibuat%5 == 0 {
				jumlahStok = fake.IntRange(1, 5) // Kritis
			} else if totalVariasiDibuat%7 == 0 {
				jumlahStok = fake.IntRange(6, 10) // Warning
			}

			spekBarang := models.SpesifikasiBarang{
				Jumlah:              jumlahStok,
				HargaBarang:         fake.IntRange(baseMin, baseMax),
				BarangID:            barang.IdBarang,
				DetailSpesifikasiID: detail.IdDetailSpesifikasi,
			}

			if err := config.DB.Where("id_barang = ? AND id_detail_spesifikasi = ?", spekBarang.BarangID, spekBarang.DetailSpesifikasiID).FirstOrCreate(&spekBarang).Error; err == nil {
				totalVariasiDibuat++
			}
		}
	}

	fmt.Printf("Yeyy, Berhasil seed spesifikasi barang (%d varian)!\n", totalVariasiDibuat)
}
