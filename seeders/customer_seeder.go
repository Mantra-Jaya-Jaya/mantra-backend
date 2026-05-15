package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
)

func SeedCustomer() {
	// Inisialisasi gofakeit v7 buat generate nomor telp random
	gofakeit.Seed(0)

	// 1. Cari user yang emailnya customer@mantra.com (Akun Rajaba Hamim)
	var user models.User
	err := config.DB.Where("email = ?", "customer@mantra.com").First(&user).Error

	if err != nil {
		fmt.Println("Waduh Error Euyy, akun customer@mantra.com gak ketemu! Pastiin SeedUser jalan duluan.")
		return
	}

	// 2. Siapin data profil customer-nya
	customerProfil := models.Customer{
		NoTelp: "085" + gofakeit.DigitN(9),
		UserId: user.IdUser,
	}

	// 3. Simpan ke database (Pakai FirstOrCreate biar gak dobel pas di-run ulang)
	if err := config.DB.Where("id_user = ?", user.IdUser).FirstOrCreate(&customerProfil).Error; err != nil {
		fmt.Println("Walah Erorr :", err)
		return
	}

	fmt.Printf("Yeyy, Berhasil seed customer! Profil untuk %s siap dipakai.\n", user.NamaLengkap)
}
