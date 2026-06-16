package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"backend-mantra/utils"
	"fmt"
	"time"

	fake "github.com/brianvoe/gofakeit/v7"
)

func SeedPesanan() {
	fmt.Println("⏳ Menyiapkan dummy data Pesanan (2 Tahun)...")

	var count int64
	config.DB.Model(&models.Pesanan{}).Count(&count)
	if count >= 100 {
		fmt.Println("Tabel pesanan udah punya minimal 100 data, proses seeding dilewati.")
		return
	}

	var customers []models.Customer
	if err := config.DB.Find(&customers).Error; err != nil || len(customers) == 0 {
		fmt.Println("Gagal: Data Customer belum ada!")
		return
	}

	var kasirs []models.Kasir
	if err := config.DB.Find(&kasirs).Error; err != nil || len(kasirs) == 0 {
		fmt.Println("Gagal: Data Kasir belum ada!")
		return
	}

	var alamats []models.Alamat
	config.DB.Find(&alamats)

	now := time.Now()
	var datesToGenerate []time.Time

	// Per hari dalam 30 hari terakhir: 3-8 pesanan per hari
	for d := 0; d < 30; d++ {
		ordersPerDay := fake.IntRange(5, 10)
		for i := 0; i < ordersPerDay; i++ {
			datesToGenerate = append(datesToGenerate, now.AddDate(0, 0, -d))
		}
	}

	totalCreated := 0
	kasirLen := len(kasirs)
	custLen := len(customers)
	alamatLen := len(alamats)

	for idx, tglPesanan := range datesToGenerate {
		randStatus := fake.IntRange(1, 100)
		var statusName string
		if randStatus <= 70 {
			statusName = "Selesai"
		} else if randStatus <= 80 {
			statusName = "Dikirim"
		} else if randStatus <= 90 {
			statusName = "Dikemas"
		} else if randStatus <= 95 {
			statusName = "Diproses"
		} else {
			statusName = "Dibatalkan"
		}

		randType := fake.IntRange(1, 100)
		tipePesanan := "Offline"
		var alamatId *uint = nil

		if randType <= 60 {
			tipePesanan = "Online"
			if alamatLen > 0 {
				alID := alamats[idx%alamatLen].IdAlamat
				alamatId = &alID
			}
		}

		kId := kasirs[idx%kasirLen].IdKasir
		cId := customers[idx%custLen].IdCustomer
		totalPembayaran := fake.IntRange(50000, 5000000)

		var kasirIdPtr *uint = &kId
		if tipePesanan == "Online" {
			kasirIdPtr = nil // Online order doesn't have a cashier initially
		}

		pesanan := models.Pesanan{
			TotalPembayaran: totalPembayaran,
			TanggalPesanan:  tglPesanan,
			TipePesananID:   utils.GetTipePesananID(tipePesanan),
			StatusPesananID: utils.GetStatusPesananID(statusName),
			CustomerID:      cId,
			KasirID:         &kId,
			AlamatID:        alamatId,
		}

		if err := config.DB.Create(&pesanan).Error; err == nil {
			totalCreated++
		}
	}

	fmt.Printf("Yeyy, Berhasil seed %d pesanan (total historis)!\n", totalCreated)
}
