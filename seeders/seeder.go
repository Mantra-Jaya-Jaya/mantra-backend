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
	SeedStatusPengantaran()
	SeedMetodePembayaran()
	BackfillPublicId()

	// Panggil seeder Child
	SeedUser()
	SeedCustomer()
	SeedKasir()
	SeedKurir()
	SeedAlamat()
	SeedBarang()
	SeedDetailSpesifikasi()
	SeedSpesifikasiBarang()
	SeedBarcode()
	SeedStokOpname()
	//SeedPesanan()
	SeedDetailPesanan()
	SeedPembayaran()
	SeedPengantaran()
	SeedKeranjang()
	SeedNotifikasi()

	fmt.Println("======================================")
	fmt.Println("SEMUA DATA BERHASIL DI-SEED!")
	fmt.Println("======================================")
}

func BackfillPublicId() {
	fmt.Println("⏳ Backfill PublicId untuk data lama...")
	config.DB.Exec("UPDATE ekspedisi SET public_id = gen_random_uuid() WHERE public_id IS NULL")
	config.DB.Exec("UPDATE ekspedisi_layanan SET public_id = gen_random_uuid() WHERE public_id IS NULL")
	config.DB.Exec("UPDATE metode_pembayaran SET public_id = gen_random_uuid() WHERE public_id IS NULL")
	config.DB.Exec("UPDATE detail_pembayaran SET public_id = gen_random_uuid() WHERE public_id IS NULL")
	fmt.Println("Yeyy, berhasil backfill PublicId!")
}
