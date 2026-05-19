package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"
	"time"

	fake "github.com/brianvoe/gofakeit/v7"
)

func SeedPesanan() {
	fmt.Println("⏳ Menyiapkan dummy data Pesanan (2 Tahun)...")

	var count int64
	config.DB.Model(&models.Pesanan{}).Count(&count)
	if count >= 400 {
		fmt.Println("Tabel pesanan udah punya minimal 400 data, proses seeding dilewati.")
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
		ordersPerDay := fake.IntRange(3, 8)
		for i := 0; i < ordersPerDay; i++ {
			datesToGenerate = append(datesToGenerate, now.AddDate(0, 0, -d))
		}
	}

	// Per bulan dalam 12 bulan terakhir (bulan 1 s/d 12 lalu): 15-30 pesanan per bulan
	for m := 1; m <= 12; m++ {
		ordersPerMonth := fake.IntRange(15, 30)
		for i := 0; i < ordersPerMonth; i++ {
			datesToGenerate = append(datesToGenerate, now.AddDate(0, -m, -fake.IntRange(0, 28)))
		}
	}

	// Per bulan dalam tahun sebelumnya (bulan 13 s/d 24 lalu): 10-20 pesanan per bulan
	for m := 13; m <= 24; m++ {
		ordersPerMonth := fake.IntRange(10, 20)
		for i := 0; i < ordersPerMonth; i++ {
			datesToGenerate = append(datesToGenerate, now.AddDate(0, -m, -fake.IntRange(0, 28)))
		}
	}

	totalCreated := 0
	kasirLen := len(kasirs)
	custLen := len(customers)
	alamatLen := len(alamats)

	for idx, tglPesanan := range datesToGenerate {
		randStatus := fake.IntRange(1, 100)
		var status string
		if randStatus <= 70 {
			status = "Selesai"
		} else if randStatus <= 80 {
			status = "Dikirim"
		} else if randStatus <= 90 {
			status = "Dikemas"
		} else if randStatus <= 95 {
			status = "Diproses"
		} else {
			status = "Dibatalkan"
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

		pesanan := models.Pesanan{
			TotalPembayaran: totalPembayaran,
			TanggalPesanan:  tglPesanan,
			TipePesanan:     tipePesanan,
			StatusPesanan:   status,
			CustomerId:      cId,
			KasirId:         kId,
			AlamatId:        alamatId,
		}

		if err := config.DB.Create(&pesanan).Error; err == nil {
			totalCreated++
		}
	}

	fmt.Printf("Yeyy, Berhasil seed %d pesanan (total historis)!\n", totalCreated)
}
