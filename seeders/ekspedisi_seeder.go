package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"backend-mantra/services"
	"fmt"
)

// SeedEkspedisi menyinkronkan data ekspedisi dan layanan langsung dari Biteship.
func SeedEkspedisi() {
	fmt.Println("⏳ Sinkronisasi data ekspedisi dari Biteship...")

	adapter := services.NewBiteshipAdapter()
	biteshipCouriers, err := adapter.GetCouriers()
	if err != nil {
		fmt.Printf("⚠️ Gagal sinkronisasi ekspedisi dari Biteship: %s\n", err.Error())
		return
	}

	totalCourierCreated := 0
	totalCourierUpdated := 0
	totalServiceCreated := 0
	totalServiceUpdated := 0

	for _, bc := range biteshipCouriers {
		var eks models.Ekspedisi
		if err := config.DB.Where("kode_api = ?", bc.CourierCode).First(&eks).Error; err != nil {
			eks = models.Ekspedisi{
				NamaEkspedisi: bc.CourierName,
				KodeApi:       bc.CourierCode,
				Deskripsi:     bc.CourierName + " Shipping Service",
				IsActive:      true,
			}

			if err := config.DB.Create(&eks).Error; err != nil {
				fmt.Printf("Error insert ekspedisi %s: %v\n", bc.CourierName, err)
				continue
			}
			totalCourierCreated++
		} else {
			updates := map[string]any{}
			if bc.CourierName != "" && eks.NamaEkspedisi != bc.CourierName {
				updates["nama_ekspedisi"] = bc.CourierName
			}
			if eks.Deskripsi == "" && bc.CourierName != "" {
				updates["deskripsi"] = bc.CourierName + " Shipping Service"
			}
			if !eks.IsActive {
				updates["is_active"] = true
			}
			if len(updates) > 0 {
				if err := config.DB.Model(&eks).Updates(updates).Error; err != nil {
					fmt.Printf("Error update ekspedisi %s: %v\n", bc.CourierName, err)
					continue
				}
				totalCourierUpdated++
			}
		}

		var lay models.EkspedisiLayanan
		if err := config.DB.Where("id_ekspedisi = ? AND nama_layanan = ?", eks.IdEkspedisi, bc.CourierServiceName).First(&lay).Error; err != nil {
			lay = models.EkspedisiLayanan{
				EkspedisiID: eks.IdEkspedisi,
				NamaLayanan: bc.CourierServiceName,
				Deskripsi:   bc.Description,
				EstimasiMin: 1,
				EstimasiMax: 3,
				IsActive:    true,
			}
			if err := config.DB.Create(&lay).Error; err != nil {
				fmt.Printf("Error insert layanan %s untuk %s: %v\n", bc.CourierServiceName, bc.CourierName, err)
				continue
			}
			totalServiceCreated++
		} else {
			updates := map[string]any{}
			if bc.Description != "" && lay.Deskripsi != bc.Description {
				updates["deskripsi"] = bc.Description
			}
			if !lay.IsActive {
				updates["is_active"] = true
			}
			if len(updates) > 0 {
				if err := config.DB.Model(&lay).Updates(updates).Error; err != nil {
					fmt.Printf("Error update layanan %s untuk %s: %v\n", bc.CourierServiceName, bc.CourierName, err)
					continue
				}
				totalServiceUpdated++
			}
		}
	}

	fmt.Printf(
		"✅ Sinkronisasi ekspedisi selesai: %d courier baru, %d courier update, %d layanan baru, %d layanan update\n",
		totalCourierCreated,
		totalCourierUpdated,
		totalServiceCreated,
		totalServiceUpdated,
	)
}
