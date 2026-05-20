package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
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
	if err := config.DB.Find(&daftarPesanan).Error; err != nil || len(daftarPesanan) == 0 {
		fmt.Println("Gagal: Data Pesanan masih kosong!")
		return
	}

	onlinePaymentTypes := []string{"qris", "bank_transfer", "gopay"}
	totalCreated := 0

	for _, pesanan := range daftarPesanan {
		var countItem int64
		config.DB.Model(&models.Pembayaran{}).Where("id_pesanan = ?", pesanan.IdPesanan).Count(&countItem)
		if countItem > 0 {
			continue // Skip jika pembayaran untuk pesanan ini sudah ada
		}

		ptype := "cash"
		status := "settlement"
		orderIdMidtrans := ""

		if pesanan.TipePesanan == "Online" {
			ptype = fake.RandomString(onlinePaymentTypes)
			if pesanan.StatusPesanan == "Selesai" || pesanan.StatusPesanan == "Dikirim" {
				status = "settlement"
			} else if pesanan.StatusPesanan == "Dibatalkan" {
				status = "cancel"
			} else {
				status = "pending"
			}
			orderIdMidtrans = fmt.Sprintf("MANTRA-%d-%d", pesanan.IdPesanan, time.Now().UnixNano())
		} else {
			ptype = "cash"
			if pesanan.StatusPesanan == "Dibatalkan" {
				status = "cancel"
			} else {
				status = "settlement"
			}
		}

		pembayaran := models.Pembayaran{
			OrderIdMidtrans: orderIdMidtrans,
			PaymentType:     ptype,
			StatusTransaksi: status,
			FraudStatus:     "accept",
			PesananID:       pesanan.IdPesanan,
		}

		if err := config.DB.Create(&pembayaran).Error; err == nil {
			totalCreated++
		}
	}

	fmt.Printf("Yeyy, berhasil seed %d pembayaran!\n", totalCreated)
}
