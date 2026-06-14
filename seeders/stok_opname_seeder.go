package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

func SeedStokOpname() {
	fmt.Println("⏳ Menyiapkan data riwayat pergerakan stok (Stok Opname)...")

	gofakeit.Seed(0)

	var daftarVarian []models.SpesifikasiBarang
	if err := config.DB.Find(&daftarVarian).Error; err != nil || len(daftarVarian) == 0 {
		fmt.Println("Gak bisa bikin stok opname karena data varian barang kosong!")
		return
	}

	totalRiwayat := 0

	for _, varian := range daftarVarian {
		var count int64
		config.DB.Model(&models.StokOpname{}).Where("id_spesifikasi_barang = ?", varian.IdSpesifikasiBarang).Count(&count)
		if count >= 4 {
			continue
		}

		recordsToCreate := gofakeit.IntRange(2, 4)
		for i := 0; i < recordsToCreate; i++ {
			tipePergerakan := randomTipePergerakan()
			jumlahStok := gofakeit.IntRange(5, 80)
			if tipePergerakan == "keluar" {
				jumlahStok = gofakeit.IntRange(1, 40)
			}

			modal := varian.HargaBarang - gofakeit.IntRange(1000, 15000)
			if modal < 1000 {
				modal = varian.HargaBarang
			}

			daysAgo := gofakeit.IntRange(1, 90)
			stokOpname := models.StokOpname{
				HargaBeli:           modal,
				TipePergerakan:      tipePergerakan,
				JumlahStok:          jumlahStok,
				Keterangan:          keteranganTipePergerakan(tipePergerakan),
				Tanggal:             time.Now().AddDate(0, 0, -daysAgo),
				SpesifikasiBarangID: varian.IdSpesifikasiBarang,
			}

			if err := config.DB.Create(&stokOpname).Error; err == nil {
				totalRiwayat++
			}
		}
	}

	fmt.Printf("Yeyy, Berhasil seed %d Stok Opname!\n", totalRiwayat)
}

func randomTipePergerakan() string {
	rand := gofakeit.IntRange(1, 100)
	switch {
	case rand <= 60:
		return "masuk"
	case rand <= 90:
		return "keluar"
	case rand <= 95:
		return "retur"
	default:
		return "penyesuaian"
	}
}

func keteranganTipePergerakan(tipe string) string {
	switch tipe {
	case "masuk":
		return "Stok masuk dari supplier"
	case "keluar":
		return "Stok keluar karena penjualan"
	case "retur":
		return "Retur barang ke supplier"
	default:
		return "Penyesuaian stok hasil opname"
	}
}
