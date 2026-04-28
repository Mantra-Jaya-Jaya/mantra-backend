package config

import (
	"log"
	"fmt"
	"os"
	
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"backend-mantra/models"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Load file .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Gagal load file .env!")
	}

	// Ambil data dari variabel environment
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// Susun DSN secara dinamis
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			host, user, password, dbname, port)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Waduh, gagal konek ke database nih!")
	}

	// Langsung sikat bikin tabelnya aja
	err = database.AutoMigrate(
		&models.Role{}, 
		&models.User{},
	)

	if err != nil {
		log.Fatal("Gagal Migrate Database!")
	}

	DB = database
	log.Println("Database Connected & Migrated Successfully!")
}