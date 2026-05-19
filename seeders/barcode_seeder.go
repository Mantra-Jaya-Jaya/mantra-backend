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

	// 2. Tarik semua data Satuan untuk di-assign secara acak
	var daftarSatuan []models.Satuan
	if err := config.DB.Find(&daftarSatuan).Error; err != nil || len(daftarSatuan) == 0 {
		fmt.Println("Satuan masih kosong! Pastiin SeedSatuan jalan duluan.")
		return
	}

	totalBarcodeDibuat := 0

	// 3. Looping ke setiap spesifikasi barang buat dikasih barcode
	for _, spek := range daftarSpesifikasi {
		// Bikin 2 barcode untuk setiap spesifikasi dengan kuantitas 1 dan 12
		kuantitasList := []uint{1, 12}

		for _, qty := range kuantitasList {
			// Ambil satuan acak
			randomIndex := gofakeit.Number(0, len(daftarSatuan)-1)
			satuanAcak := daftarSatuan[randomIndex]

			// Buat ID Barcode berupa 9 digit angka (agar muat di kolom INT database)
			kodeBarcode := uint(gofakeit.Number(100000000, 999999999))

			barcode := models.Barcode{
				IdBarcode:           kodeBarcode,
				Kuantitas:           qty,
				SpesifikasiBarangId: spek.IdSpesifikasiBarang,
				SatuanId:            satuanAcak.IdSatuan,
			}

			// Simpan ke database
			if err := config.DB.Create(&barcode).Error; err != nil {
				fmt.Printf("Gagal buat barcode untuk Spek ID %d: %v\n", spek.IdSpesifikasiBarang, err)
				continue
			}

			totalBarcodeDibuat++
		}
	}

	fmt.Printf("Yeyy, Berhasil seed %d barcode!\n", totalBarcodeDibuat)
}
