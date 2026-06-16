package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
)

func SeedKeranjang() {
	fmt.Println("⏳ Menyiapkan data keranjang...")

	gofakeit.Seed(0)

	var customers []models.Customer
	if err := config.DB.Find(&customers).Error; err != nil || len(customers) == 0 {
		fmt.Println("Gagal: Data Customer belum ada!")
		return
	}

	var daftarSpek []models.SpesifikasiBarang
	if err := config.DB.Find(&daftarSpek).Error; err != nil || len(daftarSpek) == 0 {
		fmt.Println("Gagal: Data Spesifikasi Barang masih kosong!")
		return
	}

	totalKeranjang := 0
	for _, customer := range customers {
		numItems := gofakeit.IntRange(2, 4)
		usedSpek := make(map[uint]bool)

		for i := 0; i < numItems; i++ {
			spekIdx := gofakeit.IntRange(0, len(daftarSpek)-1)
			spek := daftarSpek[spekIdx]

			// Prevent duplicates in same cart
			if usedSpek[spek.IdSpesifikasiBarang] {
				continue
			}
			usedSpek[spek.IdSpesifikasiBarang] = true

			keranjang := models.Keranjang{
				Quantity:            gofakeit.IntRange(1, 5),
				CustomerID:          customer.IdCustomer,
				SpesifikasiBarangID: spek.IdSpesifikasiBarang,
			}

			if err := config.DB.Create(&keranjang).Error; err != nil {
				fmt.Println("Error insert keranjang:", err)
				continue
			}
			totalKeranjang++
		}
	}

	fmt.Printf("Yeyy, berhasil seed %d keranjang!\n", totalKeranjang)
}
