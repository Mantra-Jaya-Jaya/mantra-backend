package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"
	// "net/url"

	fake "github.com/brianvoe/gofakeit/v7"
)

func SeedBarang() {
	fmt.Println("⏳ Menyiapkan data barang...")

	var count int64
	config.DB.Model(&models.Barang{}).Count(&count)
	if count >= 15 {
		fmt.Printf("Tabel barang udah punya %d data, proses seeding dilewati.\n", count)
		return
	}

	var satuan models.Satuan
	config.DB.First(&satuan)
	if satuan.IdSatuan == 0 {
		satuan.IdSatuan = 1
	}

	var diskon models.Diskon
	config.DB.First(&diskon)
	if diskon.IdDiskon == 0 {
		diskon.IdDiskon = 1
	}

	items := []struct{ Kat, Nama string }{
		{"Elektronik", "Laptop Gaming ASUS ROG"},
		{"Elektronik", "Smartphone Samsung A55"},
		{"Elektronik", "TWS Earphone Bluetooth"},
		{"Fashion", "Sepatu Lari Nike Air"},
		{"Fashion", "Kaos Polos Premium"},
		{"Fashion", "Jaket Hoodie Fleece"},
		{"Makanan & Minuman", "Mie Instan Box 40pcs"},
		{"Makanan & Minuman", "Kopi Sachet Box"},
		{"Makanan & Minuman", "Susu UHT Full Cream 1L"},
		{"Kesehatan", "Vitamin C 1000mg 30 tablet"},
		{"Kesehatan", "Masker KN95 Box isi 20"},
		{"Kesehatan", "Hand Sanitizer 500ml"},
		{"Olahraga", "Dumbbell Set 5kg"},
		{"Olahraga", "Matras Yoga Anti Slip"},
		{"Olahraga", "Raket Badminton Carbon"},
	}

	totalAdded := 0
	for _, item := range items {
		var kategori models.Kategori
		if err := config.DB.Where("nama_kategori = ?", item.Kat).First(&kategori).Error; err != nil {
			continue // Lewati jika kategori belum ada
		}

		barang := models.Barang{
			NamaBarang:   item.Nama,
			GambarBarang: fmt.Sprintf("https://picsum.photos/seed/%d/400/400", fake.Number(1, 1000)),
			Deskripsi:    fake.Sentence(10),
			DiskonId:     &diskon.IdDiskon,
			SatuanId:     satuan.IdSatuan,
			KategoriId:   kategori.IdKategori,
		}

		if err := config.DB.Where("nama_barang = ?", item.Nama).FirstOrCreate(&barang).Error; err == nil {
			totalAdded++
		}
	}

	fmt.Printf("Yeyy, Berhasil seed barang! Tambah: %d record.\n", totalAdded)
}
