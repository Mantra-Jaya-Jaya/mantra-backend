package seeders

import "fmt"

// Fungsi Master buat manggil semua seeder
func RunAllSeeders() {
	fmt.Println("======================================")
	fmt.Println("🚀 MEMULAI PROSES DATABASE SEEDING...")
	fmt.Println("======================================")

	// Panggil seeder Parent
	SeedRole()
	SeedKategori()
	SeedDiskon()

	// Panggil seeder Child
	SeedUser()
	SeedCustomer()
	SeedKasir()
	SeedKurir()
	SeedAlamat()


	fmt.Println("======================================")
	fmt.Println("✅ SEMUA DATA BERHASIL DI-SEED!")
	fmt.Println("======================================")
}