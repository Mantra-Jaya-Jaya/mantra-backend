package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"
	"time"

	fake "github.com/brianvoe/gofakeit/v7"
)

func SeedStokOpname() {
	fmt.Println("⏳ Menyiapkan data riwayat pergerakan stok (Stok Opname)...")

	var daftarVarian []models.SpesifikasiBarang
	if err := config.DB.Find(&daftarVarian).Error; err != nil || len(daftarVarian) == 0 {
		fmt.Println("Gak bisa bikin stok opname karena data varian barang kosong!")
		return
	}

	totalRiwayat := 0
	keteranganList := []string{"Stok masuk dari supplier", "Stok keluar penjualan", "Retur barang", "Stok opname bulanan"}

	for _, varian := range daftarVarian {
		var count int64
		config.DB.Model(&models.StokOpname{}).Where("id_spesifikasi_barang = ?", varian.IdSpesifikasiBarang).Count(&count)
		if count >= 3 {
			continue // Sudah memiliki minimal 3 log
		}

		for i := 0; i < 3; i++ {
			status := (i%2 == 0) // Mix true and false
			
			// Modal selalu lebih murah dari harga jual
			modal := varian.HargaBarang - fake.IntRange(1000, 15000)
			if modal < 0 {
				modal = varian.HargaBarang
			}

			stokOpname := models.StokOpname{
				HargaBeli:           modal,
				Status:              status,
				JumlahStok:          fake.IntRange(10, 100),
				Keterangan:          fake.RandomString(keteranganList),
				Tanggal:             time.Now().AddDate(0, 0, -fake.IntRange(1, 30)),
				SpesifikasiBarangID: varian.IdSpesifikasiBarang,
			}

			if err := config.DB.Create(&stokOpname).Error; err == nil {
				totalRiwayat++
			}
		}
	}

	fmt.Printf("Yeyy, Berhasil seed %d Stok Opname!\n", totalRiwayat)
}
