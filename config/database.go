package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"backend-mantra/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable search_path=public",
		host, user, password, dbname, port)

	// Logger kustom untuk menaikkan threshold SLOW SQL warning ke 2 detik
	customLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             2 * time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: customLogger,
	})

	if err != nil {
		log.Fatal("Waduh, gagal konek ke database nih!")
	}

	// Langsung sikat bikin tabelnya aja
	err = database.AutoMigrate(
		// Auth & User (urutan: Role -> User -> profil per role)
		&models.Role{},
		&models.User{},
		&models.RefreshToken{},
		&models.Customer{},
		&models.Karyawan{},
		&models.Kasir{},
		&models.Kurir{},
		// Alamat (butuh Customer)
		&models.Alamat{},
		// Notifikasi (butuh User)
		&models.Notifikasi{},
		// Katalog master (tidak ada dependensi)
		&models.Kategori{},
		&models.Satuan{},
		&models.Diskon{},
		&models.Spesifikasi{},
		// Katalog detail (butuh master di atas)
		&models.Barang{},
		&models.DetailSpesifikasi{},
		&models.SpesifikasiBarang{},
		&models.Barcode{},
		&models.StokOpname{},
		// Keranjang (butuh Customer & SpesifikasiBarang)
		&models.Keranjang{},
		// Transaksi (butuh Customer, Kasir, Alamat, SpesifikasiBarang)
		&models.Pesanan{},
		&models.DetailPesanan{},
		&models.Pembayaran{},
		// Metode Pembayaran (master data)
		&models.MetodePembayaran{},
		// Pengantaran (butuh Pesanan, Kurir, Ekspedisi)
		&models.StatusPengantaran{},
		&models.Ekspedisi{},
		&models.EkspedisiLayanan{},
		&models.Pengantaran{},
		// Detail Pembayaran (butuh Pembayaran)
		&models.DetailPembayaran{},
	)

	if err != nil {
		log.Fatal("Gagal Migrate Database!")
	}

	DB = database
	log.Println("Database Connected & Migrated Successfully!")
}
