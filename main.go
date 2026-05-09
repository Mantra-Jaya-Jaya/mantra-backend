package main

import (
	"fmt"
	"log"
	"net/http"
	"backend-mantra/config"
	"backend-mantra/seeders"
)

func main() {
	// Connect ke database
	config.ConnectDatabase()
	// Seed data ke tabel
	seeders.RunAllSeeders()
}