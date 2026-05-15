package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
)

func SeedAlamat() {

	gofakeit.Seed(0)

	// 1. Cari dulu user-nya
	var user models.User
	if err := config.DB.Where("email = ?", "customer@mantra.com").First(&user).Error; err != nil {
		fmt.Println("Waduh, user Rajaba Hamim gak ketemu! Pastiin SeedUser jalan.")
		return
	}

	// 2. Cari ID Customer yang nempel sama user tersebut
	var customer models.Customer
	if err := config.DB.Where("id_user = ?", user.IdUser).First(&customer).Error; err != nil {
		fmt.Println("Waduh, profil Customer belum ada! Pastiin SeedCustomer jalan duluan.")
		return
	}

	// 3. Kita racik 2 alamat sakti (Kost & Rumah)
	daftarAlamat := []models.Alamat{
		{
			CustomerId:     customer.IdCustomer,
			NamaPenerima:   user.NamaLengkap,
			LabelAlamat:    "Kost",
			NoTelpPenerima: "08" + gofakeit.DigitN(10),
			AlamatLengkap:  "Jl. Banjarsari Selatan, Tembalang, Kota Semarang",
			Latitude:       -7.051410,
			Longitude:      110.438125,
			CatatanLokasi:  "Pagar hitam, samping warung burjo",
			IsUtama:        true,
		},
		{
			CustomerId:     customer.IdCustomer,
			NamaPenerima:   user.NamaLengkap,
			LabelAlamat:    "Rumah",
			NoTelpPenerima: "08" + gofakeit.DigitN(10),
			AlamatLengkap:  "Kecamatan Selogiri, Kabupaten Wonogiri",
			Latitude:       -7.816667,
			Longitude:      110.916667,
			CatatanLokasi:  "Rumah cat hijau dekat pertigaan balai desa",
			IsUtama:        false,
		},
	}

	// 4. Looping buat masukin ke database
	for _, alamat := range daftarAlamat {
		if err := config.DB.Where("id_customer = ? AND label_alamat = ?", alamat.CustomerId, alamat.LabelAlamat).FirstOrCreate(&alamat).Error; err != nil {
			fmt.Println("Error insert alamat", alamat.LabelAlamat, ":", err)
			continue
		}
	}

	fmt.Println("Yeyy, berhasil seed alamat!")
}
