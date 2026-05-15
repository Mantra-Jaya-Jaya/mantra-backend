package seeders

import (
	"fmt"
	"backend-mantra/config"
	"backend-mantra/models"

	"github.com/brianvoe/gofakeit/v7"
)

func SeedDetailPesanan() {
	fmt.Println("⏳ Menyiapkan rincian keranjang belanja (Detail Pesanan)...")

	// 1. Cek biar gak kebanyakan data
	var count int64
	config.DB.Model(&models.DetailPesanan{}).Count(&count)
	if count > 0 {
		fmt.Println("Tabel detail_pesanan udah ada isinya, proses seeding dilewati.")
		return
	}

	// 2. Tarik semua Pesanan yang udah dibuat
	var daftarPesanan []models.Pesanan
	if err := config.DB.Find(&daftarPesanan).Error; err != nil || len(daftarPesanan) == 0 {
		fmt.Println("Gagal: Data Pesanan masih kosong! Jalankan SeedPesanan dulu.")
		return
	}

	// 3. Tarik beberapa Barang (SpesifikasiBarang) dari gudang
	var daftarVarian []models.SpesifikasiBarang
	if err := config.DB.Limit(5).Find(&daftarVarian).Error; err != nil || len(daftarVarian) == 0 {
		fmt.Println("Gagal: Data Spesifikasi Barang (Varian) masih kosong!")
		return
	}

	totalDetailDibuat := 0

	// 4. Looping: Masukin 2 barang ke setiap Nota Pesanan
	for _, pesanan := range daftarPesanan {
		
		// Beli Barang Pertama (Ambil varian index 0)
		jumlahBeli1 := gofakeit.Number(1, 3)
		hargaSatuan1 := daftarVarian[0].HargaBarang // Snapshot harga saat ini!
		
		detail1 := models.DetailPesanan{
			Jumlah:              jumlahBeli1,
			HargaSatuan:         hargaSatuan1,
			Subtotal:            jumlahBeli1 * hargaSatuan1, // Matematika dasar: Qty * Harga
			PesananId:           pesanan.IdPesanan,
			SpesifikasiBarangId: daftarVarian[0].IdSpesifikasiBarang,
		}

		// Beli Barang Kedua (Ambil varian index 1)
		jumlahBeli2 := gofakeit.Number(1, 2)
		hargaSatuan2 := daftarVarian[1].HargaBarang
		
		detail2 := models.DetailPesanan{
			Jumlah:              jumlahBeli2,
			HargaSatuan:         hargaSatuan2,
			Subtotal:            jumlahBeli2 * hargaSatuan2,
			PesananId:           pesanan.IdPesanan,
			SpesifikasiBarangId: daftarVarian[1].IdSpesifikasiBarang,
		}

		// Save ke database
		config.DB.Create(&detail1)
		config.DB.Create(&detail2)
		
		totalDetailDibuat += 2
	}

	fmt.Printf("Yeyy, Berhasil seed detail pesanan!")
}