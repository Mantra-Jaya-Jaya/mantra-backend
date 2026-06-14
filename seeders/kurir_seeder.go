package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"backend-mantra/utils"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

func SeedKurir() {
	gofakeit.Seed(0)

	var role models.Role
	if err := config.DB.Where("nama_role = ?", "Kurir").First(&role).Error; err != nil {
		fmt.Println("Walah, Role Kurir gak ketemu!")
		return
	}

	var users []models.User
	if err := config.DB.Where("id_role = ?", role.IdRole).Find(&users).Error; err != nil || len(users) == 0 {
		fmt.Println("Walah, akun Kurir gak ketemu! Pastiin SeedUser jalan duluan.")
		return
	}

	totalKaryawan := 0
	totalKurir := 0

	for _, user := range users {
		tglLahir := gofakeit.DateRange(time.Date(1995, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2003, 12, 31, 0, 0, 0, 0, time.UTC))

		karyawanProfil := models.Karyawan{
			NoTelp:             "08" + gofakeit.DigitN(10),
			TempatLahir:        gofakeit.City(),
			TanggalLahir:       tglLahir,
			JenisKelamin:       gofakeit.RandomString([]string{"Laki-laki", "Perempuan"}),
			Alamat:             gofakeit.Address().Address,
			PendidikanTerakhir: gofakeit.RandomString([]string{"SMA/SMK", "D3", "S1"}),
			Nik:                "3374" + gofakeit.DigitN(12),
			StatusKaryawanID:   utils.GetStatusKaryawanID("Aktif"),
			UserID:             user.IdUser,
		}

		if err := config.DB.Where("id_user = ?", user.IdUser).FirstOrCreate(&karyawanProfil).Error; err != nil {
			fmt.Println("Error create karyawan untuk", user.NamaLengkap, ":", err)
			continue
		}
		totalKaryawan++

		kurirProfil := models.Kurir{
			KaryawanID: karyawanProfil.IdKaryawan,
		}
		if err := config.DB.Where("id_karyawan = ?", karyawanProfil.IdKaryawan).FirstOrCreate(&kurirProfil).Error; err != nil {
			fmt.Println("Error create kurir profil:", err)
			continue
		}
		totalKurir++
	}

	fmt.Printf("Yeyy, berhasil seed %d karyawan + %d kurir!\n", totalKaryawan, totalKurir)
}
