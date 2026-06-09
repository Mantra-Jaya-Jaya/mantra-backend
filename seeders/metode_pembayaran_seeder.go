package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"
)

func SeedMetodePembayaran() {
	fmt.Println("⏳ Menyiapkan data metode pembayaran...")

	daftarMetode := []models.MetodePembayaran{
		{NamaMetode: "Cash", KodeMetode: "cash", Penyedia: "internal", Urutan: 1, IsActive: true},
		{NamaMetode: "QRIS", KodeMetode: "qris", Penyedia: "midtrans", Urutan: 2, IsActive: true},
		{NamaMetode: "Virtual Account", KodeMetode: "va", Penyedia: "midtrans", Urutan: 3, IsActive: true},
		{NamaMetode: "E-Wallet", KodeMetode: "ewallet", Penyedia: "midtrans", Urutan: 4, IsActive: true},
		{NamaMetode: "COD (Bayar di Tempat)", KodeMetode: "cod", Penyedia: "internal", Urutan: 5, IsActive: true},
	}

	for _, m := range daftarMetode {
		if err := config.DB.Where("kode_metode = ?", m.KodeMetode).FirstOrCreate(&m).Error; err != nil {
			fmt.Println("Error insert metode", m.NamaMetode, ":", err)
			continue
		}
	}

	fmt.Println("Yeyy, berhasil seed metode pembayaran!")
}
