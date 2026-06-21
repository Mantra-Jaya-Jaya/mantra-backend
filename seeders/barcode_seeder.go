package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
)

func SeedBarcode() {
	fmt.Println("⏳ Menyiapkan data barcode...")

	gofakeit.Seed(0)

	// 1. Tarik semua data Spesifikasi Barang
	var daftarSpesifikasi []models.SpesifikasiBarang
	if err := config.DB.Find(&daftarSpesifikasi).Error; err != nil || len(daftarSpesifikasi) == 0 {
		fmt.Println("Spesifikasi Barang masih kosong! Pastiin SeedSpesifikasiBarang jalan duluan.")
		return
	}

	var batch []models.Barcode
	for _, spek := range daftarSpesifikasi {
		for _, qty := range []uint{1, 12} {
			batch = append(batch, models.Barcode{
				KodeBarcode:         fmt.Sprintf("%012d", gofakeit.Number(100000000000, 999999999999)),
				Kuantitas:           qty,
				SpesifikasiBarangID: spek.IdSpesifikasiBarang,
			})
		}
	}

	if err := config.DB.Create(&batch).Error; err != nil {
		fmt.Printf("Gagal batch insert barcode: %v\n", err)
	} else {
		fmt.Printf("Yeyy, Berhasil seed %d barcode!\n", len(batch))
	}
}
