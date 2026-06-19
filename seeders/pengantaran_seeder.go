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

	var count int64
	config.DB.Model(&models.Pengantaran{}).Count(&count)
	if count > 0 {
		fmt.Println("Tabel pengantaran udah ada isinya, proses seeding dilewati.")
		return
	}

	dikirimID := utils.GetStatusPesananIDSafe("Dikirim")
	selesaiPesananID := utils.GetStatusPesananIDSafe("Selesai")
	if dikirimID == 0 || selesaiPesananID == 0 {
		fmt.Println("Gagal: Status pesanan 'Dikirim'/'Selesai' belum ada!")
		return
	}

	var daftarPesanan []models.Pesanan
	if err := config.DB.
		Where("id_tipe_pesanan = ?", utils.GetTipePesananID("Online")).
		Where("id_status_pesanan IN ?", []uint{dikirimID, selesaiPesananID}).
		Find(&daftarPesanan).Error; err != nil || len(daftarPesanan) == 0 {
		fmt.Println("Gagal: Data pesanan online dengan status Dikirim/Selesai masih kosong!")
		return
	}

	menungguID := utils.GetStatusPengantaranIDSafe("Menunggu Pickup")
	jalanID := utils.GetStatusPengantaranIDSafe("Dalam Perjalanan")
	selesaiPengID := utils.GetStatusPengantaranIDSafe("Selesai")
	if menungguID == 0 || jalanID == 0 || selesaiPengID == 0 {
		fmt.Println("Gagal: Setup status pengantaran belum lengkap!")
		return
	}

	var kurir models.Kurir
	adaKurir := config.DB.First(&kurir).Error == nil

	var ekspedisi models.Ekspedisi
	adaEkspedisi := config.DB.First(&ekspedisi).Error == nil

	now := time.Now()

	for i, pesanan := range daftarPesanan {
		roll := i % 10
		var (
			statusPengID  uint
			waktuPickup   *time.Time
			waktuSampai   *time.Time
			updateStatus  uint
		)

		switch {
		case roll < 2:
			// 20% — Menunggu Pickup (belum di-pickup kurir)
			statusPengID = menungguID
			waktuPickup = nil
			waktuSampai = nil
			updateStatus = dikirimID
		case roll < 5:
			// 30% — Dalam Perjalanan
			statusPengID = jalanID
			t := now.Add(-time.Duration(30+i) * time.Minute)
			waktuPickup = &t
			waktuSampai = nil
			updateStatus = dikirimID
		default:
			// 50% — Selesai
			statusPengID = selesaiPengID
			tPickup := now.Add(-time.Duration(120+i) * time.Minute)
			tSampai := now.Add(-time.Duration(60+i) * time.Minute)
			waktuPickup = &tPickup
			waktuSampai = &tSampai
			updateStatus = selesaiPesananID
		}

		var idKurir *uint
		if adaKurir {
			idKurir = &kurir.IdKurir
		}
		var idEkspedisi *uint
		if adaEkspedisi {
			idEkspedisi = &ekspedisi.IdEkspedisi
		}

		pengantaran := models.Pengantaran{
			WaktuPickup:         waktuPickup,
			WaktuSampai:         waktuSampai,
			LastLatitude:        -7.051410,
			LastLongitude:       110.438125,
			FotoBuktiPengiriman: "",
			PesananID:           pesanan.IdPesanan,
			KurirID:             idKurir,
			StatusPengantaranID: statusPengID,
			EkspedisiID:         idEkspedisi,
		}

		if err := config.DB.Create(&pengantaran).Error; err != nil {
			fmt.Println("Error insert pengantaran:", err)
			continue
		}

		if pesanan.StatusPesananID != updateStatus {
			config.DB.Model(&pesanan).Update("id_status_pesanan", updateStatus)
		}
	}

	fmt.Println("Yeyy, berhasil seed pengantaran!")
}
