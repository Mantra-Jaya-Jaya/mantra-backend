package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"backend-mantra/utils"
	"fmt"
	"time"

	fake "github.com/brianvoe/gofakeit/v7"
)

func SeedPembayaran() {
	fmt.Println("⏳ Menyiapkan data pembayaran historis...")

	var count int64
	config.DB.Model(&models.Pembayaran{}).Count(&count)
	if count >= 400 {
		fmt.Println("Tabel pembayaran udah punya minimal 400 data, proses seeding dilewati.")
		return
	}

	var daftarPesanan []models.Pesanan
	if err := config.DB.Preload("StatusPesanan").Find(&daftarPesanan).Error; err != nil || len(daftarPesanan) == 0 {
		fmt.Println("Gagal: Data Pesanan masih kosong!")
		return
	}

	// Pre-fetch which pesanan already have pembayaran records
	var existingPesananIDs []uint
	config.DB.Model(&models.Pembayaran{}).Distinct("id_pesanan").Pluck("id_pesanan", &existingPesananIDs)
	existingMap := make(map[uint]bool)
	for _, id := range existingPesananIDs {
		existingMap[id] = true
	}

	onlinePaymentTypes := []string{"qris", "bank_transfer", "gopay"}
	totalCreated := 0
	var pembayaranList []models.Pembayaran

	for _, pesanan := range daftarPesanan {
		if existingMap[pesanan.IdPesanan] {
			continue
		}

		ptype := "cash"
		status := "settlement"
		orderIdMidtrans := ""
		statusName := ""
		if pesanan.StatusPesanan != nil {
			statusName = pesanan.StatusPesanan.NamaStatus
		}

		if pesanan.TipePesananID == utils.GetTipePesananID("Online") {
			ptype = fake.RandomString(onlinePaymentTypes)
			switch statusName {
			case "Selesai", "Dikirim":
				status = "settlement"
			case "Dibatalkan":
				status = "cancel"
			default:
				status = "pending"
			}
			orderIdMidtrans = fmt.Sprintf("MANTRA-%d-%d", pesanan.IdPesanan, time.Now().UnixNano())
		} else {
			ptype = "cash"
			if statusName == "Dibatalkan" {
				status = "cancel"
			} else {
				status = "settlement"
			}
		}

		pembayaran := models.Pembayaran{
			OrderIdMidtrans:   orderIdMidtrans,
			TipePembayaranID:  utils.GetTipePembayaranID(ptype),
			StatusTransaksiID: utils.GetStatusTransaksiID(status),
			FraudStatusID:     utils.GetFraudStatusID("accept"),
			PesananID:         pesanan.IdPesanan,
		}

		pembayaranList = append(pembayaranList, pembayaran)
	}

	if len(pembayaranList) > 0 {
		if err := config.DB.CreateInBatches(pembayaranList, 100).Error; err == nil {
			totalCreated = len(pembayaranList)
		} else {
			fmt.Println("Error batch insert pembayaran:", err)
		}
	}

	fmt.Printf("Yeyy, berhasil seed %d pembayaran!\n", totalCreated)
}
