package seeders

import (
  "backend-mantra/config"
  "backend-mantra/models"
  "fmt"

  "github.com/brianvoe/gofakeit/v7"
)

func SeedBarcode() {
  fmt.Println("⏳ Menyiapkan data barcode...")

  gofakeit.Seed(0)

  // 1. Tarik semua data Spesifikasi Barang
  var daftarSpesifikasi []models.SpesifikasiBarang
  if err := config.DB.Find(&daftarSpesifikasi).Error; err != nil || len(daftarSpesifikasi) == 0 {
    fmt.Println("Spesifikasi Barang masih kosong! Pastiin SeedSpesifikasiBarang jalan duluan.")
    return
  }

  totalBarcodeDibuat := 0

  // 2. Looping ke setiap spesifikasi barang buat dikasih barcode
  for _, spek := range daftarSpesifikasi {
    // Bikin 2 barcode untuk setiap spesifikasi (misal: qty 1 buat eceran, qty 12 buat lusinan)
    kuantitasList := []uint{1, 12}

    for _, qty := range kuantitasList {
      
      // Bikin Kode Barcode beneran (Format String/Varchar)
      // Pakai %012d biar kalau angkanya kurang, otomatis ditambahin 0 di depannya
      kodeBarcodeString := fmt.Sprintf("%012d", gofakeit.Number(100000000000, 999999999999))

      barcode := models.Barcode{
        // 🚀 IdBarcode KITA HAPUS! Biarin Postgres yang ngasih nomor urut otomatis
        KodeBarcode:         kodeBarcodeString, 
        Kuantitas:           qty,
        SpesifikasiBarangID: spek.IdSpesifikasiBarang, // Pastikan ini pakai 'ID' besar di ujung sesuai model lu
        // 🚀 SatuanId KITA HAPUS! Karena emang gak ada di ERD tabel barcode
      }

      // Simpan ke database
      if err := config.DB.Create(&barcode).Error; err != nil {
        fmt.Printf("Gagal buat barcode untuk Spek ID %d: %v\n", spek.IdSpesifikasiBarang, err)
        continue
      }

      totalBarcodeDibuat++
    }
  }

  fmt.Printf("Yeyy, Berhasil seed %d barcode!\n", totalBarcodeDibuat)
}