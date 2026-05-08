package main

import (
	"fmt"
	"log"
	"net/http"
	"backend-mantra/config"
	"backend-mantra/routes"
	
)

func main() {
	// Mulai koneksi ke database
	config.ConnectDatabase()
	// Panggil daftar routes
	routes.SetupRoutes()

	fmt.Println("🚀 Server jalan di http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}