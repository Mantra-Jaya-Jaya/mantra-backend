package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"
)

func SeedKeranjang() {
	fmt.Println("⏳ Menyiapkan data keranjang...")

	var count int64
	config.DB.Model(&models.Keranjang{}).Count(&count)
	if count > 0 {
		fmt.Println("Tabel keranjang udah ada isinya, proses seeding dilewati.")
		return
	}

	var customer models.Customer
	if err := config.DB.First(&customer).Error; err != nil {
		fmt.Println("Gagal: Data Customer belum ada!")
		return
	}

	var daftarSpek []models.SpesifikasiBarang
	if err := config.DB.Limit(3).Find(&daftarSpek).Error; err != nil || len(daftarSpek) == 0 {
		fmt.Println("Gagal: Data Spesifikasi Barang masih kosong!")
		return
	}

	for _, spek := range daftarSpek {
		keranjang := models.Keranjang{
			Quantity:            2,
			CustomerID:          customer.IdCustomer,
			SpesifikasiBarangID: spek.IdSpesifikasiBarang,
		}

		if err := config.DB.Create(&keranjang).Error; err != nil {
			fmt.Println("Error insert keranjang:", err)
			continue
		}
	}

	fmt.Println("Yeyy, berhasil seed keranjang!")
}
