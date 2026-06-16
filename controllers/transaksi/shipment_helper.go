package transaksi

import (
	"fmt"
	"os"

	"backend-mantra/config"
	"backend-mantra/models"
	"backend-mantra/services"
	"backend-mantra/utils"
)

func processExternalShipment(pesananID uint) {
	var pesanan models.Pesanan
	if err := config.DB.
		Preload("TipeKurir").
		Preload("Ekspedisi").
		Preload("LayananEkspedisi").
		Preload("Alamat").
		Preload("DetailPesanan.SpesifikasiBarang.Barang").
		First(&pesanan, pesananID).Error; err != nil {
		fmt.Printf("❌ Shipment Error: Gagal load pesanan %d: %s\n", pesananID, err.Error())
		return
	}

	if pesanan.TipeKurir == nil || pesanan.TipeKurir.NamaTipe != "external" {
		return
	}
	if pesanan.Ekspedisi == nil || pesanan.LayananEkspedisi == nil {
		fmt.Printf("❌ Shipment Error: Pesanan %d eksternal tapi data ekspedisi tidak lengkap\n", pesananID)
		return
	}
	if pesanan.Alamat == nil {
		fmt.Printf("❌ Shipment Error: Pesanan %d eksternal tapi alamat tidak ada\n", pesananID)
		return
	}
	if pesanan.NomorResi != nil && *pesanan.NomorResi != "" {
		fmt.Printf("ℹ️ Shipment: Pesanan %d sudah punya resi, skip\n", pesananID)
		return
	}

	var items []services.CreateShipmentItem
	for _, d := range pesanan.DetailPesanan {
		items = append(items, services.CreateShipmentItem{
			Name:     d.SpesifikasiBarang.Barang.NamaBarang,
			Weight:   d.SpesifikasiBarang.BeratBarang,
			Quantity: d.Jumlah,
			Value:    d.HargaSatuan,
			Length:   d.SpesifikasiBarang.PanjangBarang,
			Width:    d.SpesifikasiBarang.LebarBarang,
			Height:   d.SpesifikasiBarang.TinggiBarang,
		})
	}

	adapter := services.NewBiteshipAdapter()
	result, err := adapter.CreateShipment(services.CreateShipmentRequest{
		OriginAddress:        os.Getenv("BITESHIP_STORE_ADDRESS"),
		OriginCoordinate:     os.Getenv("BITESHIP_STORE_COORDINATE_LAT") + "," + os.Getenv("BITESHIP_STORE_COORDINATE_LONG"),
		DestinationAddress:   pesanan.Alamat.AlamatLengkap,
		DestinationCoordinate: fmt.Sprintf("%f,%f", pesanan.Alamat.Latitude, pesanan.Alamat.Longitude),
		CourierCode:          pesanan.Ekspedisi.KodeApi,
		CourierServiceCode:   pesanan.LayananEkspedisi.NamaLayanan,
		Items:                items,
	})
	if err != nil {
		fmt.Printf("❌ Shipment Error: Gagal create shipment pesanan %d: %s\n", pesananID, err.Error())
		return
	}
	if result == nil || result.WaybillID == "" {
		fmt.Printf("❌ Shipment Error: Response kosong untuk pesanan %d\n", pesananID)
		return
	}

	updates := map[string]interface{}{
		"nomor_resi":        result.WaybillID,
		"id_status_pesanan": utils.GetStatusPesananID("Dikirim"),
	}
	if err := config.DB.Model(&models.Pesanan{}).Where("id_pesanan = ?", pesananID).Updates(updates).Error; err != nil {
		fmt.Printf("❌ Shipment Error: Gagal update status pesanan %d: %s\n", pesananID, err.Error())
		return
	}

	fmt.Printf("✅ Shipment Success: Pesanan %d → waybill %s → status Dikirim\n", pesananID, result.WaybillID)
}
