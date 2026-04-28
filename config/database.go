package config

import (
	"backend-mantra/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Sesuaikan user & password postgres di komputer kalian masing masing ya!
	dsn := "host=localhost user=postgres password=openpgpwd dbname=mantra_db port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Waduh, gagal konek ke database nih!")
	}

	// Langsung sikat bikin tabelnya aja
	err = database.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.Customer{},
		&models.Kurir{},
		&models.Kasir{},
		&models.Alamat{},
		&models.Barang{},
		&models.Barcode{},
	)

	if err != nil {
		log.Fatal("Gagal Migrate Database!")
	}

	DB = database
	log.Println("Database Connected & Migrated Successfully!")
}
