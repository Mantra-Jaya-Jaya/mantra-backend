package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
)

func SeedSpesifikasiBarang() {
	fmt.Println("⏳ Menyiapkan stok dan harga tiap varian barang...")

	gofakeit.Seed(0)

	// 1. Tarik semua data Barang
	var daftarBarang []models.Barang
	if err := config.DB.Find(&daftarBarang).Error; err != nil || len(daftarBarang) == 0 {
		fmt.Println("Barang masih kosong! Pastiin SeedBarang jalan duluan.")
		return
	}

	// 2. Tarik semua data Detail Spesifikasi
	var daftarDetailSpek []models.DetailSpesifikasi
	if err := config.DB.Find(&daftarDetailSpek).Error; err != nil || len(daftarDetailSpek) == 0 {
		fmt.Println("Detail Spesifikasi kosong! Pastiin SeedDetailSpesifikasi jalan.")
		return
	}

	totalVariasiDibuat := 0

	// 3. Looping ke setiap barang buat dikasih varian
	for _, barang := range daftarBarang {
		// Kita kasih masing-masing barang 3 varian secara random
		for i := 0; i < 3; i++ {
			// Ambil detail spesifikasi acak dari list
			randomIndex := gofakeit.Number(0, len(daftarDetailSpek)-1)
			detailSpekAcak := daftarDetailSpek[randomIndex]

			// Setup data varian, stok, dan harga
			spekBarang := models.SpesifikasiBarang{
				Jumlah:              gofakeit.Number(5, 150),
				HargaBarang:         gofakeit.Number(15000, 350000),
				BarangID:            barang.IdBarang,
				DetailSpesifikasiID: detailSpekAcak.IdDetailSpesifikasi,
			}

			// FirstOrCreate biar datanya gak dobel:
			if err := config.DB.Where("id_barang = ? AND id_detail_spesifikasi = ?", spekBarang.BarangID, spekBarang.DetailSpesifikasiID).FirstOrCreate(&spekBarang).Error; err != nil {
				fmt.Println("Error:", err)
				continue
			}

			totalVariasiDibuat++
		}
	}

	fmt.Printf("Yeyy, Berhasil seed spesifikasi barang!")
}
