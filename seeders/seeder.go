package seeders

import (
	"backend-mantra/config"
	"fmt"
)

// Fungsi Master buat manggil semua seeder
func RunAllSeeders() {
	fmt.Println("======================================")
	fmt.Println("MEMULAI PROSES DATABASE SEEDING...")
	fmt.Println("======================================")

	// Panggil seeder Parent (Master)
	SeedRole()
	SeedKategori()
	SeedDiskon()
	SeedSatuan()
	SeedSpesifikasi()
	SeedEkspedisi()
	SeedStatusPesanan()
	SeedStatusTransaksi()
	SeedStatusNotifikasi()
	SeedStatusKaryawan()
	SeedShiftKasir()
	SeedFraudStatus()
	SeedTipePembayaran()
	SeedTipePesanan()
	SeedStatusPengantaran()
	SeedMetodePembayaran()
	SeedTipeKurir()
	SeedPengaturanToko()
	BackfillPublicId()

	// Panggil seeder Child
	SeedUser()
	SeedCustomer()
	SeedKasir()
	SeedKurir()
	// SeedAlamat()
	SeedBarang()
	SeedDetailSpesifikasi()
	SeedSpesifikasiBarang()
	SeedBarcode()
	SeedStokOpname()
	// SeedPesanan()
	// SeedDetailPesanan()
	// SeedPembayaran()
	// SeedPengantaran()
	SeedKeranjang()
	SeedNotifikasi()

	fmt.Println("======================================")
	fmt.Println("SEMUA DATA BERHASIL DI-SEED!")
	fmt.Println("======================================")
}

func BackfillPublicId() {
	fmt.Println("⏳ Mengecek data legacy tanpa PublicId...")

	type backfillTarget struct {
		tableName string
		label     string
	}

	targets := []backfillTarget{
		{tableName: "ekspedisi", label: "ekspedisi"},
		{tableName: "ekspedisi_layanan", label: "layanan ekspedisi"},
		{tableName: "metode_pembayaran", label: "metode pembayaran"},
		{tableName: "detail_pembayaran", label: "detail pembayaran"},
	}

	var totalUpdated int64
	for _, target := range targets {
		var missingCount int64
		if err := config.DB.Table(target.tableName).Where("public_id IS NULL").Count(&missingCount).Error; err != nil {
			fmt.Printf("⚠️ Gagal cek public_id untuk %s: %s\n", target.label, err.Error())
			continue
		}

		if missingCount == 0 {
			continue
		}

		result := config.DB.Exec(fmt.Sprintf("UPDATE %s SET public_id = gen_random_uuid() WHERE public_id IS NULL", target.tableName))
		if result.Error != nil {
			fmt.Printf("⚠️ Gagal backfill public_id untuk %s: %s\n", target.label, result.Error.Error())
			continue
		}

		totalUpdated += result.RowsAffected
	}

	if totalUpdated == 0 {
		fmt.Println("✅ PublicId sudah lengkap, skip backfill legacy.")
		return
	}

	fmt.Printf("✅ Backfill PublicId selesai: %d baris diperbarui\n", totalUpdated)
}
