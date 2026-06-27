package transaksi

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"backend-mantra/config"
	"backend-mantra/models"
	"backend-mantra/services"
)

func processExternalShipment(pesananID uint) {
	session, err := shipmentProcessingConn()
	if err != nil {
		fmt.Printf("❌ Shipment Error: Gagal menyiapkan lock untuk pesanan %d: %s\n", pesananID, err.Error())
		return
	}
	defer session.Close()

	locked, err := tryLockShipmentProcessing(session, pesananID)
	if err != nil {
		fmt.Printf("❌ Shipment Error: Gagal mengamankan proses pesanan %d: %s\n", pesananID, err.Error())
		return
	}
	if !locked {
		fmt.Printf("ℹ️ Shipment: Pesanan %d sedang diproses oleh worker lain, skip duplikat\n", pesananID)
		return
	}
	defer releaseShipmentProcessingLock(session, pesananID)

	var pesanan models.Pesanan
	if err := config.DB.
		Preload("TipeKurir").
		Preload("Ekspedisi").
		Preload("LayananEkspedisi").
		Preload("Alamat").
		Preload("Customer.User").
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
	if pesanan.BiteshipOrderID != nil && *pesanan.BiteshipOrderID != "" {
		fmt.Printf("ℹ️ Shipment: Pesanan %d sudah punya biteship order, skip\n", pesananID)
		return
	}

	var items []services.CreateShipmentItem
	for _, d := range pesanan.DetailPesanan {
		if d.SpesifikasiBarang.Barang.IdBarang == 0 {
			continue
		}
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

	if len(items) == 0 {
		fmt.Printf("❌ Shipment Error: Pesanan %d tidak memiliki detail barang valid untuk shipment\n", pesananID)
		return
	}

	// Parse postal code toko dari env
	originPostal, _ := strconv.Atoi(os.Getenv("BITESHIP_STORE_POSTAL_CODE"))
	originLat, _ := strconv.ParseFloat(os.Getenv("BITESHIP_STORE_COORDINATE_LAT"), 64)
	originLng, _ := strconv.ParseFloat(os.Getenv("BITESHIP_STORE_COORDINATE_LONG"), 64)

	// Parse postal code tujuan
	destPostal, _ := strconv.Atoi(pesanan.Alamat.KodePos)

	adapter := services.NewBiteshipAdapter()
	result, err := adapter.CreateShipment(services.CreateShipmentRequest{
		OriginAddress:           os.Getenv("BITESHIP_STORE_ADDRESS"),
		OriginLat:               originLat,
		OriginLng:               originLng,
		OriginPostalCode:        originPostal,
		DestinationAddress:      pesanan.Alamat.AlamatLengkap,
		DestinationLat:          pesanan.Alamat.Latitude,
		DestinationLng:          pesanan.Alamat.Longitude,
		DestinationPostalCode:   destPostal,
		DestinationContactName:  pesanan.Alamat.NamaPenerima,
		DestinationContactPhone: pesanan.Alamat.NoTelpPenerima,
		DestinationNote:         pesanan.Alamat.CatatanLokasi,
		CourierCode:             pesanan.Ekspedisi.KodeApi,
		CourierServiceCode:      pesanan.LayananEkspedisi.KodeLayanan,
		OrderNote:               pesanan.Catatan,
		Items:                   items,
	})
	if err != nil {
		fmt.Printf("❌ Shipment Error: Gagal create shipment pesanan %d: %s\n", pesananID, err.Error())
		return
	}
	if result == nil {
		fmt.Printf("❌ Shipment Error: Response shipment nil untuk pesanan %d\n", pesananID)
		return
	}

	// Debug log response Biteship
	fmt.Printf("📦 Biteship Response pesanan %d: success=%v id=%s waybill=%s status=%s\n",
		pesananID, result.Success, result.ID, result.WaybillID, result.Status)

	if result.ID == "" {
		fmt.Printf("❌ Shipment Error: Biteship tidak mengembalikan order ID untuk pesanan %d\n", pesananID)
		return
	}

	updates := map[string]any{
		"nomor_resi":        result.WaybillID,
		"biteship_order_id": result.ID,
	}
	if err := config.DB.Model(&models.Pesanan{}).Where("id_pesanan = ?", pesananID).Updates(updates).Error; err != nil {
		fmt.Printf("❌ Shipment Error: Gagal update status pesanan %d: %s\n", pesananID, err.Error())
		return
	}

	fmt.Printf("✅ Shipment Success: Pesanan %d → order_id=%s waybill=%s → status Dikirim\n",
		pesananID, result.ID, result.WaybillID)
}

func shipmentProcessingConn() (*sql.Conn, error) {
	db, err := config.DB.DB()
	if err != nil {
		return nil, err
	}

	return db.Conn(context.Background())
}

func tryLockShipmentProcessing(conn *sql.Conn, pesananID uint) (bool, error) {
	var locked bool
	if err := conn.QueryRowContext(context.Background(), "SELECT pg_try_advisory_lock($1)", int64(pesananID)).Scan(&locked); err != nil {
		return false, err
	}
	return locked, nil
}

func releaseShipmentProcessingLock(conn *sql.Conn, pesananID uint) {
	if _, err := conn.ExecContext(context.Background(), "SELECT pg_advisory_unlock($1)", int64(pesananID)); err != nil {
		fmt.Printf("⚠️ Shipment Warning: Gagal melepas lock pesanan %d: %s\n", pesananID, err.Error())
	}
}
