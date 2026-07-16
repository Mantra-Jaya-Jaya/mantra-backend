package services

import (
	"log"

	"backend-mantra/config"
	"backend-mantra/models"
)

// SyncCouriersFromBiteship mengambil data kurir dari Biteship API dan upsert ke DB lokal.
// Fungsi ini standalone (tidak terikat gin.Context) sehingga bisa dipanggil
// dari mana saja: controller admin, auto-sync saat cek ongkir, atau saat server start.
func SyncCouriersFromBiteship() error {
	adapter := NewBiteshipAdapter()
	biteshipCouriers, err := adapter.GetCouriers()
	if err != nil {
		return err
	}

	for _, bc := range biteshipCouriers {
		var eks models.Ekspedisi

		// 1. Upsert Ekspedisi (Courier)
		err := config.DB.Where("kode_api = ?", bc.CourierCode).First(&eks).Error
		if err != nil {
			// Not found, create new
			eks = models.Ekspedisi{
				NamaEkspedisi: bc.CourierName,
				KodeApi:       bc.CourierCode,
				Logo:          "",
				Deskripsi:     bc.CourierName + " Shipping Service",
				IsActive:      true,
			}
			if err := config.DB.Create(&eks).Error; err != nil {
				log.Printf("[SyncCouriers] Gagal create ekspedisi %s: %v", bc.CourierCode, err)
				continue
			}
		} else {
			updated := false
			if eks.NamaEkspedisi == "" {
				eks.NamaEkspedisi = bc.CourierName
				updated = true
			}
			if eks.Deskripsi == "" {
				eks.Deskripsi = bc.CourierName + " Shipping Service"
				updated = true
			}
			if updated {
				config.DB.Save(&eks)
			}
		}

		// 2. Upsert Layanan
		var lay models.EkspedisiLayanan
		err = config.DB.Where("id_ekspedisi = ? AND nama_layanan = ?", eks.IdEkspedisi, bc.CourierServiceName).First(&lay).Error
		if err != nil {
			// Not found, create new
			lay = models.EkspedisiLayanan{
				EkspedisiID: eks.IdEkspedisi,
				NamaLayanan: bc.CourierServiceName,
				KodeLayanan: bc.CourierServiceCode,
				Deskripsi:   bc.Description,
				EstimasiMin: 1,
				EstimasiMax: 3,
				IsActive:    true,
			}
			config.DB.Create(&lay)
		} else {
			// Found, update description/kode if changed or empty
			updated := false
			if lay.Deskripsi != bc.Description && bc.Description != "" {
				lay.Deskripsi = bc.Description
				updated = true
			}
			if lay.KodeLayanan == "" && bc.CourierServiceCode != "" {
				lay.KodeLayanan = bc.CourierServiceCode
				updated = true
			}
			if updated {
				config.DB.Save(&lay)
			}
		}
	}

	log.Println("[SyncCouriers] Sinkronisasi ekspedisi dari Biteship selesai")
	return nil
}
