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
	// 1. Cari user lewat email kurir@mantra.com
	var user models.User
	err := config.DB.Where("email = ?", "kurir@mantra.com").First(&user).Error

	if err != nil {
		fmt.Println("Waduh, akun (kurir@mantra.com) gak ketemu! Pastiin SeedUser jalan duluan.")
		return
	}

	// 2. Siapin tanggal lahir (Misal: 10 Oktober 2003)
	tglLahir := time.Date(2003, time.October, 10, 0, 0, 0, 0, time.Local)

	// 3. Siapin profil Karyawan-nya
	karyawanProfil := models.Karyawan{
		NoTelp:             "08" + gofakeit.DigitN(10),
		TempatLahir:        "Wonogiri", // Sesuaikan dengan vibes daerah lu bro!
		TanggalLahir:       tglLahir,
		JenisKelamin:       "Perempuan",
		Alamat:             "Kecamatan Selogiri, Kabupaten Wonogiri",
		PendidikanTerakhir: "SMA Negeri 1 Wonogiri",
		Nik:                "3312" + gofakeit.DigitN(12), // 16 Digit NIK (Kode Wonogiri 3312)
		Status:             "Aktif",
		StatusKaryawanID:   utils.GetStatusKaryawanID("Aktif"),
		UserID:             user.IdUser,
	}

	if err := config.DB.Where("id_user = ?", user.IdUser).FirstOrCreate(&karyawanProfil).Error; err != nil {
		fmt.Println("Error:", err)
		return
	}

	kurirProfil := models.Kurir{
		KaryawanID: karyawanProfil.IdKaryawan,
	}
	if err := config.DB.Where("id_karyawan = ?", karyawanProfil.IdKaryawan).FirstOrCreate(&kurirProfil).Error; err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Yeyy, Berhasil seed kurir!")
}
