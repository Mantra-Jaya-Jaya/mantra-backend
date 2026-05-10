package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
)

func SeedKategori() {
	gofakeit.Seed(0)

	// Bikin 10 data
	for i := 0; i < 10; i++ {
		kategori := models.Kategori{
			// Pakai ProductCategory() dari docs biar namanya beneran kayak "Electronics", "Clothing", bukan kata random aneh
			NamaKategori: gofakeit.ProductCategory(),

			// (Karena di v7 fungsi Image() balikin byte, kita pakai trik URL kucing lucu aja)
			IconKategori: fmt.Sprintf("https://cataas.com/cat?width=100&height=100&random=%s", gofakeit.UUID()),
		}

		if err := config.DB.Create(&kategori).Error; err != nil {
			fmt.Println("Hadehh error:", err)
			return
		}
	}

	fmt.Println("yeeyyy, berhasil seed kategori!")
}
