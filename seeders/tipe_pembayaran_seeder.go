package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"
)

func SeedTipePembayaran() {
	fmt.Println("Seeding tipe pembayaran...")

	tipeList := []string{"cash", "non-cash", "qris", "bank_transfer", "gopay"}
	for _, nama := range tipeList {
		tipe := models.TipePembayaran{NamaTipe: nama}
		if err := config.DB.Where("nama_tipe = ?", nama).FirstOrCreate(&tipe).Error; err != nil {
			fmt.Printf("Error seeding tipe pembayaran '%s': %v\n", nama, err)
		}
	}
	fmt.Println("Tipe pembayaran selesai di-seed!")
}
