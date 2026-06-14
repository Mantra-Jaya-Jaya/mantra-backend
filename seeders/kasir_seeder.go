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

	// 1. Cari user lewat email kasir@mantra.com
	var user models.User
	err := config.DB.Where("email = ?", "kasir@mantra.com").First(&user).Error

	if err != nil {
		fmt.Println("Walah, akun (kasir@mantra.com) gak ketemu! Pastiin SeedUser jalan duluan.")
		return
	}

	// 2. Siapin data tanggal lahir (Contoh: 15 Mei 2004)
	tglLahir := time.Date(2004, time.May, 15, 0, 0, 0, 0, time.Local)

	// 3. Siapin profil Karyawan-nya
	karyawanProfil := models.Karyawan{
		NoTelp:             "08" + gofakeit.DigitN(10),
		TempatLahir:        "Semarang",
		TanggalLahir:       tglLahir,
		JenisKelamin:       "Perempuan",
		Alamat:             "Jl. Prof. Sudarto, Tembalang, Kota Semarang",
		PendidikanTerakhir: "D3 Teknik Komputer",
		Nik:                "3374" + gofakeit.DigitN(12),
		Status:             "Aktif",
		StatusKaryawanID:   utils.GetStatusKaryawanID("Aktif"),
		UserID:             user.IdUser,
	}

	// 4. Simpan ke database
	if err := config.DB.Where("id_user = ?", user.IdUser).FirstOrCreate(&karyawanProfil).Error; err != nil {
		fmt.Println("Error:", err)
		return
	}

	kasirProfil := models.Kasir{
		ShiftKasirID: utils.GetShiftKasirID("Pagi"),
		KaryawanID: karyawanProfil.IdKaryawan,
	}
	if err := config.DB.Where("id_karyawan = ?", karyawanProfil.IdKaryawan).FirstOrCreate(&kasirProfil).Error; err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Yeyy, Berhasil seed kasir!")
}
