package main

import (
	"backend-mantra/config"
	"backend-mantra/seeders"
)

func main() {
	// Connect ke database
	config.ConnectDatabase()
	// Seed data ke tabel
	seeders.RunAllSeeders()
}