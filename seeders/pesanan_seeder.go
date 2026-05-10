package seeders

import (
	"fmt"
	"time"
	"backend-mantra/config"
	"backend-mantra/models"
)

func SeedPesanan() {
	fmt.Println("⏳ Menyiapkan dummy data Pesanan...")

	// 1. Cek biar gak beranak-pinak kalau kena restart
	var count int64
	config.DB.Model(&models.Pesanan{}).Count(&count)
	if count >= 2 {
		fmt.Println("Tabel pesanan udah ada isinya, proses seeding dilewati.")
		return
	}

	// 2. Tarik aktor-aktornya (Customer, Kasir, Alamat)
	var customer models.Customer
	if err := config.DB.First(&customer).Error; err != nil {
		fmt.Println("Gagal: Data Customer belum ada!")
		return
	}

	var kasir models.Kasir
	if err := config.DB.First(&kasir).Error; err != nil {
		fmt.Println("Gagal: Data Kasir belum ada!")
		return
	}

	var alamat models.Alamat
	if err := config.DB.First(&alamat).Error; err != nil {
		fmt.Println("Gagal: Data Alamat belum ada!")
		return
	}

	// 3. Bikin Skenario 1: Pesanan Online (Pakai Alamat)
	pesananOnline := models.Pesanan{
		TotalPembayaran: 125000,
		TanggalPesanan:  time.Now(),
		TipePesanan:     "Online",
		StatusPesanan:   "Dikemas",
		CustomerId:      customer.IdCustomer,
		KasirId:         kasir.IdKasir,
		AlamatId:        &alamat.IdAlamat, 
	}

	// 4. Bikin Skenario 2: Pesanan Offline / Takeaway (Gak pakai alamat)
	pesananOffline := models.Pesanan{
		TotalPembayaran: 45000,
		TanggalPesanan:  time.Now().Add(-24 * time.Hour), // Ceritanya pesanan kemarin
		TipePesanan:     "Offline",
		StatusPesanan:   "Selesai",
		CustomerId:      customer.IdCustomer,
		KasirId:         kasir.IdKasir,
		AlamatId:        nil, // Aman banget di-set nil karena tipe datanya pointer
	}

	// 5. Eksekusi masukin ke Database
	if err := config.DB.Create(&pesananOnline).Error; err != nil {
		fmt.Println("Gagal insert pesanan online:", err)
	}
	if err := config.DB.Create(&pesananOffline).Error; err != nil {
		fmt.Println("Gagal insert pesanan offline:", err)
	}

	fmt.Println("Yeyy, Berhasil seed pesanan!")
}