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

	// 3. Kita racik 5 alamat: dalam radius, luar kota, luar negeri
	daftarAlamat := []models.Alamat{
		{
			CustomerID:     customer.IdCustomer,
			NamaPenerima:   user.NamaLengkap,
			LabelAlamat:    "Kost",
			NoTelpPenerima: "08" + gofakeit.DigitN(10),
			AlamatLengkap:  "Jl. Banjarsari Selatan, Tembalang, Kota Semarang",
			KodePos:        "50275",
			Latitude:       -7.051410,
			Longitude:      110.438125,
			CatatanLokasi:  "Pagar hitam, samping warung burjo — dekat Polines",
			IsUtama:        true,
		},
		{
			CustomerID:     customer.IdCustomer,
			NamaPenerima:   user.NamaLengkap,
			LabelAlamat:    "Kontrakan",
			NoTelpPenerima: "08" + gofakeit.DigitN(10),
			AlamatLengkap:  "Perumahan Sambiroto Indah, Sambiroto, Tembalang, Semarang",
			KodePos:        "50276",
			Latitude:       -7.041000,
			Longitude:      110.442000,
			CatatanLokasi:  "Perumahan, blok C nomor 5",
			IsUtama:        false,
		},
		{
			CustomerID:     customer.IdCustomer,
			NamaPenerima:   user.NamaLengkap,
			LabelAlamat:    "Rumah",
			NoTelpPenerima: "08" + gofakeit.DigitN(10),
			AlamatLengkap:  "Kecamatan Selogiri, Kabupaten Wonogiri",
			KodePos:        "57612",
			Latitude:       -7.816667,
			Longitude:      110.916667,
			CatatanLokasi:  "Rumah cat hijau dekat pertigaan balai desa",
			IsUtama:        false,
		},
		{
			CustomerID:     customer.IdCustomer,
			NamaPenerima:   "Budi Santoso",
			LabelAlamat:    "Cabang",
			NoTelpPenerima: "08" + gofakeit.DigitN(10),
			AlamatLengkap:  "Jl. Soekarno-Hatta No. 123, Kaliwungu, Kabupaten Kendal",
			KodePos:        "51351",
			Latitude:       -6.983333,
			Longitude:      110.409722,
			CatatanLokasi:  "Sebelah utara pasar Kaliwungu",
			IsUtama:        false,
		},
		{
			CustomerID:     customer.IdCustomer,
			NamaPenerima:   user.NamaLengkap,
			LabelAlamat:    "Kantor",
			NoTelpPenerima: "08" + gofakeit.DigitN(10),
			AlamatLengkap:  "Kedutaan Besar RI, Jalan Tun Razak, Kuala Lumpur, Malaysia",
			KodePos:        "50450",
			Latitude:       3.139003,
			Longitude:      101.686855,
			CatatanLokasi:  "Gedung KBRI, lantai 2",
			IsUtama:        false,
		},
	}

	// 4. Hapus alamat lama customer ini, biar bisa di-create ulang
	config.DB.Where("id_customer = ?", customer.IdCustomer).Delete(&models.Alamat{})

	// 5. Looping buat masukin ke database
	for _, alamat := range daftarAlamat {
		if err := config.DB.Create(&alamat).Error; err != nil {
			fmt.Println("Error insert alamat", alamat.LabelAlamat, ":", err)
			continue
		}
	}

	fmt.Println("Yeyy, berhasil seed alamat!")
}
