package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"

	fake "github.com/brianvoe/gofakeit/v7"
)

func SeedDetailPesanan() {
	fmt.Println("⏳ Menyiapkan rincian keranjang belanja (Detail Pesanan)...")

	var daftarPesanan []models.Pesanan
	if err := config.DB.Find(&daftarPesanan).Error; err != nil || len(daftarPesanan) == 0 {
		fmt.Println("Gagal: Data Pesanan masih kosong! Jalankan SeedPesanan dulu.")
		return
	}

	var daftarVarian []models.SpesifikasiBarang
	if err := config.DB.Find(&daftarVarian).Error; err != nil || len(daftarVarian) == 0 {
		fmt.Println("Gagal: Data Spesifikasi Barang (Varian) masih kosong!")
		return
	}

	totalDetailDibuat := 0
	var detailList []models.DetailPesanan

	// Pre-fetch which pesanan already have details
	var existingPesananIDs []uint
	config.DB.Model(&models.DetailPesanan{}).Distinct("id_pesanan").Pluck("id_pesanan", &existingPesananIDs)
	existingMap := make(map[uint]bool)
	for _, id := range existingPesananIDs {
		existingMap[id] = true
	}

	for _, pesanan := range daftarPesanan {
		if existingMap[pesanan.IdPesanan] {
			continue // Lewati jika pesanan sudah memiliki detail
		}

		var subtotalPesanan int
		numItems := fake.IntRange(2, 4)
		for i := 0; i < numItems; i++ {
			varian := daftarVarian[fake.IntRange(0, len(daftarVarian)-1)]
			qty := fake.IntRange(1, 5)
			hargaSatuan := varian.HargaBarang
			subtotal := qty * hargaSatuan
			subtotalPesanan += subtotal

			detail := models.DetailPesanan{
				Jumlah:              qty,
				HargaSatuan:         hargaSatuan,
				Subtotal:            subtotal,
				PesananID:           pesanan.IdPesanan,
				SpesifikasiBarangID: varian.IdSpesifikasiBarang,
			}
			detailList = append(detailList, detail)
		}

		// Update total_pembayaran pesanan dari akumulasi subtotal
		config.DB.Model(&models.Pesanan{}).Where("id_pesanan = ?", pesanan.IdPesanan).Update("total_pembayaran", subtotalPesanan)
	}

	if len(detailList) > 0 {
		if err := config.DB.CreateInBatches(detailList, 100).Error; err == nil {
			totalDetailDibuat = len(detailList)
		} else {
			fmt.Println("Error batch insert detail pesanan:", err)
		}
	}

	fmt.Printf("Yeyy, Berhasil seed %d detail pesanan!\n", totalDetailDibuat)
}
