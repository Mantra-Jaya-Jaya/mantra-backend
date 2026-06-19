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

	semarangCoords := [][2]float64{
		{-7.051410, 110.438125},  // Tembalang (dalam radius ~500m)
		{-7.041000, 110.442000},  // Sambiroto (dalam radius ~1.5km)
		{-7.050000, 110.430000},  // Polines area (dalam)
		{3.139003, 101.686855},   // KBRI Kuala Lumpur (LUAR NEGERI)
		{-6.983333, 110.409722},  // Kendal (luar radius ~15km)
		{-7.816667, 110.916667},  // Wonogiri (luar radius ~100km)
	}
	semarangAlamat := []string{
		"Jl. Banjarsari Selatan, Tembalang, Kota Semarang",
		"Perumahan Sambiroto Indah, Tembalang, Semarang",
		"Jl. Prof. Soedarto, S.H., Tembalang, Semarang",
		"Kedutaan Besar RI, Jalan Tun Razak, Kuala Lumpur, Malaysia",
		"Jl. Soekarno-Hatta No. 123, Kaliwungu, Kabupaten Kendal",
		"Kecamatan Selogiri, Kabupaten Wonogiri",
	}
	semarangKodePos := []string{
		"50275",
		"50276",
		"50275",
		"50450",
		"51351",
		"57612",
	}
	catatan := []string{
		"Pagar hitam, samping warung burjo",
		"Perumahan, blok C nomor 5",
		"Depan kampus Polines",
		"Gedung KBRI, lantai 2",
		"Sebelah utara pasar",
		"Rumah cat hijau dekat pertigaan balai desa",
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

		// Lepas FK dulu, baru hapus alamat lama biar bisa di-create ulang
		config.DB.Exec("UPDATE pesanan SET id_alamat = NULL WHERE id_customer = ?", customerProfil.IdCustomer)
		config.DB.Where("id_customer = ?", customerProfil.IdCustomer).Delete(&models.Alamat{})

		numAlamat := 1
		if gofakeit.Bool() {
			numAlamat = 2
		}
		labels := []string{"Rumah", "Kantor", "Kos"}
		for i := 0; i < numAlamat; i++ {
			idx := (totalAlamat + i) % len(semarangCoords)
			alamat := models.Alamat{
				CustomerID:     customerProfil.IdCustomer,
				NamaPenerima:   user.NamaLengkap,
				LabelAlamat:    labels[i%len(labels)],
				NoTelpPenerima: "085" + gofakeit.DigitN(9),
				AlamatLengkap:  semarangAlamat[idx],
				KodePos:        semarangKodePos[idx],
				Latitude:       semarangCoords[idx][0],
				Longitude:      semarangCoords[idx][1],
				CatatanLokasi:  catatan[idx],
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
