package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
)

func SeedCustomer() {
	gofakeit.Seed(0)

	var role models.Role
	if err := config.DB.Where("nama_role = ?", "Customer").First(&role).Error; err != nil {
		fmt.Println("Waduh, Role Customer gak ketemu!")
		return
	}

	var users []models.User
	if err := config.DB.Where("id_role = ?", role.IdRole).Find(&users).Error; err != nil || len(users) == 0 {
		fmt.Println("Waduh, akun Customer gak ketemu! Pastiin SeedUser jalan duluan.")
		return
	}

	totalCustomer := 0
	totalAlamat := 0

	for _, user := range users {
		customerProfil := models.Customer{
			NoTelp: "085" + gofakeit.DigitN(9),
			UserID: user.IdUser,
		}
		if err := config.DB.Where("id_user = ?", user.IdUser).FirstOrCreate(&customerProfil).Error; err != nil {
			fmt.Println("Error create customer untuk", user.NamaLengkap, ":", err)
			continue
		}
		totalCustomer++

		numAlamat := 1
		if gofakeit.Bool() {
			numAlamat = 2
		}
		labels := []string{"Rumah", "Kantor", "Kos"}
		for i := 0; i < numAlamat; i++ {
			alamat := models.Alamat{
				CustomerID:     customerProfil.IdCustomer,
				NamaPenerima:   user.NamaLengkap,
				LabelAlamat:    labels[i%len(labels)],
				NoTelpPenerima: "085" + gofakeit.DigitN(9),
				AlamatLengkap:  gofakeit.Address().Address,
				KodePos:        gofakeit.Zip(),
				Latitude:       gofakeit.Latitude(),
				Longitude:      gofakeit.Longitude(),
				CatatanLokasi:  gofakeit.Sentence(5),
				IsUtama:        i == 0,
			}
			if err := config.DB.Create(&alamat).Error; err != nil {
				fmt.Println("Error create alamat:", err)
				continue
			}
			totalAlamat++
		}
	}

	fmt.Printf("Yeyy, berhasil seed %d customer dengan %d alamat!\n", totalCustomer, totalAlamat)
}
