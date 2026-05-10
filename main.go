package main

import (
	// "fmt" karena tidak ada fungsi fmt.Something dalam func main() ini aku comment dahulu, semisal gak ada yan gpakai bisa hapus
	// "net/http" karena belum dipakai, bisa diaktifka nkembali jika di /routes terdapat file yang menggunakan ini
	"backend-mantra/config"
	"backend-mantra/routes"
	"backend-mantra/seeders"
	"backend-mantra/routes"
)

func main() {
	// Koneksi ke database (dan AutoMigrate)
	config.ConnectDatabase()

	// Seed data ke tabel
	seeders.RunAllSeeders()

	// Inisialisasi Gin router
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r)

	// Ambil port dari env atau default ke 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Jalankan server
	log.Printf("Server berjalan di port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}

