package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"backend-mantra/utils"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

func SeedKasir() {
	gofakeit.Seed(0)

	var role models.Role
	if err := config.DB.Where("nama_role = ?", "Kasir").First(&role).Error; err != nil {
		fmt.Println("Walah, Role Kasir gak ketemu!")
		return
	}

	var users []models.User
	if err := config.DB.Where("id_role = ?", role.IdRole).Find(&users).Error; err != nil || len(users) == 0 {
		fmt.Println("Walah, akun Kasir gak ketemu! Pastiin SeedUser jalan duluan.")
		return
	}

	shiftList := []string{"Pagi", "Siang", "Malam"}

	totalKaryawan := 0
	totalKasir := 0

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

		shift := gofakeit.RandomString(shiftList)
		kasirProfil := models.Kasir{
			KaryawanID:   karyawanProfil.IdKaryawan,
			ShiftKasirID: utils.GetShiftKasirID(shift),
		}
		if err := config.DB.Where("id_karyawan = ?", karyawanProfil.IdKaryawan).FirstOrCreate(&kasirProfil).Error; err != nil {
			fmt.Println("Error create kasir profil:", err)
			continue
		}
		totalKasir++
	}

	fmt.Printf("Yeyy, berhasil seed %d karyawan + %d kasir!\n", totalKaryawan, totalKasir)
}
