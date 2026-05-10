package main

import (
	"fmt"
	"log"
	"net/http"
	"backend-mantra/config"
	"backend-mantra/seeders"
	"backend-mantra/routes"
)

func main() {
	// Connect ke database
	config.ConnectDatabase()
	// Setup routes
	routes.SetupRoutes()
	// Seed data ke tabel
	// seeders.RunAllSeeders()
	// Nyalain server
	fmt.Println("🚀 Server jalan di http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}