package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"
	"time"
)

func SeedPengantaran() {
	fmt.Println("⏳ Menyiapkan data pengantaran...")

	var count int64
	config.DB.Model(&models.Pengantaran{}).Count(&count)
	if count > 0 {
		fmt.Println("Tabel pengantaran udah ada isinya, proses seeding dilewati.")
		return
	}

	var daftarPesanan []models.Pesanan
	if err := config.DB.Where("tipe_pesanan = ?", "Online").Find(&daftarPesanan).Error; err != nil || len(daftarPesanan) == 0 {
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
		statusID := statusSelesai.IdStatusPengantaran
		if i%2 != 0 {
			statusID = statusJalan.IdStatusPengantaran
		}

		pengantaran := models.Pengantaran{
			WaktuPickup:         time.Now().Add(-2 * time.Hour),
			WaktuSampai:         time.Now().Add(-1 * time.Hour),
			LastLatitude:        -7.051410,
			LastLongitude:       110.438125,
			FotoBuktiPengiriman: "https://picsum.photos/400/400",
			PesananID:           pesanan.IdPesanan,
			KurirID:             kurir.IdKurir,
			StatusPengantaranID: statusID,
			EkspedisiID:         ekspedisi.IdEkspedisi,
		}

		if err := config.DB.Create(&pengantaran).Error; err != nil {
			fmt.Println("Error insert pengantaran:", err)
			continue
		}
	}

	fmt.Println("Yeyy, berhasil seed pengantaran!")
}
