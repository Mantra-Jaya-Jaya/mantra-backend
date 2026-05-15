package seeders

import (
	"fmt"
	"time"
	"backend-mantra/config"
	"backend-mantra/models"

	"github.com/brianvoe/gofakeit/v7"
)

func SeedStokOpname() {
	fmt.Println("⏳ Menyiapkan data riwayat pergerakan stok (Stok Opname)...")

	// 1. Ambil semua data SpesifikasiBarang yang udah dibuat
	var daftarVarian []models.SpesifikasiBarang
	if err := config.DB.Find(&daftarVarian).Error; err != nil || len(daftarVarian) == 0 {
		fmt.Println("Gak bisa bikin stok opname karena data varian barang kosong!")
		return
	}

	totalRiwayat := 0

	// 2. Looping setiap varian barang
	for _, varian := range daftarVarian {
		// Kita bikin 3 catatan riwayat per varian
		for i := 0; i < 3; i++ {
			status := gofakeit.Bool() // Acak: true (masuk) atau false (keluar)
			
			// Logika Harga Beli: biasanya lebih murah dari harga jual (sekitar 60-80%)
			modal := int(float64(varian.HargaBarang) * gofakeit.Float64Range(0.6, 0.8))

			keterangan := "Restock barang dari supplier"
			if !status {
				keterangan = "Penyesuaian stok / barang rusak"
			}

			stokOpname := models.StokOpname{
				HargaBeli:           modal,
				Status:              status,
				JumlahStok:          gofakeit.Number(5, 50),
				Keterangan:          keterangan,
				Tanggal:             gofakeit.DateRange(time.Now().AddDate(0, 0, -30), time.Now()),
				SpesifikasiBarangID: varian.IdSpesifikasiBarang,
			}

			// Tambahkan Datanya ke database
			if err := config.DB.Create(&stokOpname).Error; err != nil {
				fmt.Println("Error:", err)
				continue
			}
			totalRiwayat++
		}
	}

	fmt.Printf("Yeyy, Berhasil Seed Stok Opname!")
}