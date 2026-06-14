package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"backend-mantra/utils"
	"fmt"
	"time"
)

func SeedPengantaran() {
	fmt.Println("⏳ Menyiapkan data pengantaran...")

	// 🚀 LOGIKA AMAN: Cek apakah tabel pengantaran udah ada isinya
	var count int64
	config.DB.Model(&models.Pengantaran{}).Count(&count)
	if count > 0 {
		fmt.Println("Tabel pengantaran udah ada isinya, proses seeding dilewati.")
		return
	}

	var daftarPesanan []models.Pesanan
	if err := config.DB.Where("id_tipe_pesanan = ?", utils.GetTipePesananID("Online")).Find(&daftarPesanan).Error; err != nil || len(daftarPesanan) == 0 {
		fmt.Println("Gagal: Data Pesanan Online masih kosong!")
		return
	}

	var kurir models.Kurir
	if err := config.DB.First(&kurir).Error; err != nil {
		fmt.Println("Gagal: Data Kurir masih kosong!")
		return
	}

	var ekspedisi models.Ekspedisi
	if err := config.DB.First(&ekspedisi).Error; err != nil {
		fmt.Println("Gagal: Data Ekspedisi masih kosong!")
		return
	}

	var statusSelesai models.StatusPengantaran
	if err := config.DB.Where("nama_status = ?", "Selesai").First(&statusSelesai).Error; err != nil {
		fmt.Println("Gagal: Status 'Selesai' belum ada!")
		return
	}

	var statusJalan models.StatusPengantaran
	if err := config.DB.Where("nama_status = ?", "Dalam Perjalanan").First(&statusJalan).Error; err != nil {
		fmt.Println("Gagal: Status 'Dalam Perjalanan' belum ada!")
		return
	}

	for i, pesanan := range daftarPesanan {
		// 🚀 1. BIKIN VARIABEL PENAMPUNG BUAT POINTER SEBELUM STRUCT
		waktuPickup := time.Now().Add(-2 * time.Hour)
		waktuSampai := time.Now().Add(-1 * time.Hour)
		idKurir := kurir.IdKurir
		idEkspedisi := ekspedisi.IdEkspedisi

		statusID := statusSelesai.IdStatusPengantaran
		var ptrWaktuSampai *time.Time = &waktuSampai // Default: Udah sampai

		if i%2 != 0 {
			statusID = statusJalan.IdStatusPengantaran
			// 🚀 LOGIKA SAKTI: Kalau masih di jalan, waktu sampainya kita set NULL!
			ptrWaktuSampai = nil
		}

		// 🚀 2. MASUKIN ALAMAT MEMORI (&) KE DALAM STRUCT PENGANTARAN
		pengantaran := models.Pengantaran{
			WaktuPickup:         &waktuPickup,   // Pakai & (Pointer)
			WaktuSampai:         ptrWaktuSampai, // Bisa terisi nilai pointer, atau nil (NULL)
			LastLatitude:        -7.051410,
			LastLongitude:       110.438125,
			FotoBuktiPengiriman: "https://picsum.photos/400/400",
			PesananID:           pesanan.IdPesanan,
			KurirID:             &idKurir, // Pakai & (Pointer)
			StatusPengantaranID: statusID,
			EkspedisiID:         &idEkspedisi, // Pakai & (Pointer)
		}

		if err := config.DB.Create(&pengantaran).Error; err != nil {
			fmt.Println("Error insert pengantaran:", err)
			continue
		}
	}

	fmt.Println("Yeyy, berhasil seed pengantaran!")
}
