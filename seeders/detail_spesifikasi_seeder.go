package seeders

import (
	"fmt"
	"backend-mantra/config"
	"backend-mantra/models"
)

func SeedDetailSpesifikasi() {
	fmt.Println("⏳ Menyiapkan data Detail Spesifikasi (Warna & Ukuran)...")

	// 1. Cari ID Spesifikasi "Warna" dulu
	var spekWarna models.Spesifikasi
	if err := config.DB.Where("nama_spesifikasi = ?", "Warna").First(&spekWarna).Error; err != nil {
		fmt.Println("Spesifikasi 'Warna' gak ketemu! Pastiin SeedSpesifikasi jalan duluan.")
		return
	}

	// 2. Cari ID Spesifikasi "Ukuran"
	var spekUkuran models.Spesifikasi
	if err := config.DB.Where("nama_spesifikasi = ?", "Ukuran").First(&spekUkuran).Error; err != nil {
		fmt.Println("Spesifikasi 'Ukuran' gak ketemu! Pastiin SeedSpesifikasi jalan duluan.")
		return
	}

	// 3. Masukin Data Warna (Merah, Kuning, Hijau)
	daftarWarna := []string{"Merah", "Kuning", "Hijau"}
	for _, warna := range daftarWarna {
		detailWarna := models.DetailSpesifikasi{
			NamaDetailSpesifikasi: warna,
			SpesifikasiID:         spekWarna.IdSpesifikasi, 

		// FirstOrCreate biar gak dobel kalau Air restart
		if err := config.DB.Where("nama_detail_spesifikasi = ? AND id_spesifikasi = ?", warna, spekWarna.IdSpesifikasi).FirstOrCreate(&detailWarna).Error; err != nil {
			fmt.Println("Error", warna, ":", err)
		}
	}

	// 4. Masukin Data Ukuran (36 sampai 45) pakai looping
	for i := 36; i <= 45; i++ {
		ukuran := fmt.Sprintf("%d", i) // Ubah angka int jadi string
		detailUkuran := models.DetailSpesifikasi{
			NamaDetailSpesifikasi: ukuran,
			SpesifikasiID:         spekUkuran.IdSpesifikasi, // Tempelin ke ID Ukuran
		}

		if err := config.DB.Where("nama_detail_spesifikasi = ? AND id_spesifikasi = ?", ukuran, spekUkuran.IdSpesifikasi).FirstOrCreate(&detailUkuran).Error; err != nil {
			fmt.Println("Error", ukuran, ":", err)
		}
	}

	fmt.Println("Yeyy, Berhasil seed detail spesifikasi!")
}