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
	for _, pesanan := range daftarPesanan {
		var count int64
		config.DB.Model(&models.DetailPesanan{}).Where("id_pesanan = ?", pesanan.IdPesanan).Count(&count)
		if count > 0 {
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
				PesananId:           pesanan.IdPesanan,
				SpesifikasiBarangId: varian.IdSpesifikasiBarang,
			}

			if err := config.DB.Where("id_pesanan = ? AND id_spesifikasi_barang = ?", detail.PesananId, detail.SpesifikasiBarangId).FirstOrCreate(&detail).Error; err == nil {
				totalDetailDibuat++
			}
		}

		// Update total_pembayaran pesanan dari akumulasi subtotal
		config.DB.Model(&pesanan).Update("total_pembayaran", subtotalPesanan)
	}

	fmt.Printf("Yeyy, Berhasil seed %d detail pesanan!\n", totalDetailDibuat)
}
