package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"backend-mantra/utils"
	"fmt"
	"time"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

func SeedPesanan() {
	fmt.Println("⏳ Menyiapkan dummy data Pesanan (2 Tahun)...")

	var count int64
	config.DB.Model(&models.Pesanan{}).Count(&count)
	if count >= 800 {
		fmt.Println("Tabel pesanan udah punya minimal 800 data, proses seeding dilewati.")
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

	// Per hari selama 1 tahun terakhir: density makin jarak ke belakang makin kecil
	for d := 0; d < 365; d++ {
		var ordersPerDay int
		switch {
		case d < 30: // 30 hari terakhir: 4-8 per hari
			ordersPerDay = fake.IntRange(4, 8)
		case d < 90: // 31-90 hari lalu: 3-5 per hari
			ordersPerDay = fake.IntRange(3, 5)
		case d < 180: // 91-180 hari lalu: 2-4 per hari
			ordersPerDay = fake.IntRange(2, 4)
		default: // 181-365 hari lalu: 1-3 per hari
			ordersPerDay = fake.IntRange(1, 3)
		}
		for i := 0; i < ordersPerDay; i++ {
			datesToGenerate = append(datesToGenerate, now.AddDate(0, 0, -d))
		}
	}

	totalCreated := 0
	kasirLen := len(kasirs)
	custLen := len(customers)
	alamatLen := len(alamats)

	var pesananList []models.Pesanan

	for idx, tglPesanan := range datesToGenerate {
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

		// Pilih status berdasarkan tipe pesanan
		randStatus := fake.IntRange(1, 100)
		var statusName string
		if tipePesanan == "Offline" {
			switch {
			case randStatus <= 80:
				statusName = "Selesai"
			case randStatus <= 93:
				statusName = "Dikemas"
			default:
				statusName = "Dibatalkan"
			}
		} else {
			switch {
			case randStatus <= 65:
				statusName = "Selesai"
			case randStatus <= 77:
				statusName = "Dikirim"
			case randStatus <= 90:
				statusName = "Dikemas"
			default:
				statusName = "Dibatalkan"
			}
		}

		// Tentukan tipe kurir
		tipeKurir := "internal"
		var ekspedisiID *uint = nil
		var layananEkspedisiID *uint = nil
		ongkir := 0
		var nomorResi *string = nil
		if tipePesanan == "Online" {
			if fake.IntRange(1, 100) <= 50 {
				tipeKurir = "external"
				// Ambil random ekspedisi
				var ekspedisi models.Ekspedisi
				if err := config.DB.Where("is_active = ?", true).Order("RANDOM()").First(&ekspedisi).Error; err == nil {
					ekspedisiID = &ekspedisi.IdEkspedisi
					var layanan models.EkspedisiLayanan
					if err := config.DB.Where("id_ekspedisi = ?", ekspedisi.IdEkspedisi).Order("RANDOM()").First(&layanan).Error; err == nil {
						layananEkspedisiID = &layanan.IdEkspedisiLayanan
					}
					ongkir = fake.IntRange(5000, 50000)
				}
				// External yang Dikirim/Selesai set nomor resi
				if statusName == "Dikirim" || statusName == "Selesai" {
					nr := fmt.Sprintf("%s-%d-%d", ekspedisi.KodeApi, fake.IntRange(100000, 999999), fake.IntRange(1000, 9999))
					nomorResi = &nr
				}
			} else {
				tipeKurir = "internal"
			}
		}
		// Tentukan kasir
		var kasirIdPtr *uint = nil
		setButuhKasir := map[string]bool{
			"Dikemas": true, "Dikirim": true, "Selesai": true,
		}
		if tipePesanan == "Offline" {
			kId := kasirs[fake.IntRange(0, kasirLen-1)].IdKasir
			kasirIdPtr = &kId
		} else if setButuhKasir[statusName] {
			kId := kasirs[fake.IntRange(0, kasirLen-1)].IdKasir
			kasirIdPtr = &kId
		}

		cId := customers[fake.IntRange(0, custLen-1)].IdCustomer
		totalPembayaran := fake.IntRange(50000, 5000000)

		pesanan := models.Pesanan{
			PublicId:           uuid.New(),
			TotalPembayaran:    totalPembayaran,
			TanggalPesanan:     tglPesanan,
			TipePesananID:      utils.GetTipePesananID(tipePesanan),
			StatusPesananID:    utils.GetStatusPesananID(statusName),
			TipeKurirID:        utils.GetTipeKurirID(tipeKurir),
			CustomerID:         &cId,
			KasirID:            kasirIdPtr,
			AlamatID:           alamatId,
			EkspedisiID:        ekspedisiID,
			LayananEkspedisiID: layananEkspedisiID,
			OngkosKirim:        ongkir,
			NomorResi:          nomorResi,
		}

		pesananList = append(pesananList, pesanan)
	}

	if len(pesananList) > 0 {
		if err := config.DB.CreateInBatches(pesananList, 100).Error; err == nil {
			totalCreated = len(pesananList)
		} else {
			fmt.Println("Error batch insert pesanan:", err)
		}
	}

	fmt.Printf("Yeyy, Berhasil seed %d pesanan (total historis)!\n", totalCreated)
}
