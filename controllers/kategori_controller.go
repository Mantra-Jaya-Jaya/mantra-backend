package controllers

import (
	"encoding/json"
	"net/http"
	"backend-mantra/config"
	"backend-mantra/models"
)

// Fungsi buat Get Kategori
func GetKategori(w http.ResponseWriter, r *http.Request) {
	// Wajib set header biar dibaca sebagai JSON
	w.Header().Set("Content-Type", "application/json")

	var kategori []models.Kategori
	
	// Tarik data dari database pake GORM
	if err := config.DB.Find(&kategori).Error; err != nil {
		// Kalau error server (500)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "Gagal mengambil data kategori",
			"data":    nil,
		})
		return
	}

	// Kalau berhasil (200 OK)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Berhasil mengambil daftar kategori",
		"data":    kategori,
	})
}

// Fungsi buat nambah Kategori baru (Method: POST)
func CreateKategori(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1. Karena loket ini khusus nerima barang (POST), kita usir kalau ada yang pakai GET
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "Method tidak diizinkan, wajib pakai POST",
		})
		return
	}

	// 2. Siapin "laci" buat nangkep kiriman data JSON dari Postman
	var inputKategori models.Kategori

	// 3. Buka "surat" request (r.Body) dan masukin isinya ke laci (inputKategori)
	if err := json.NewDecoder(r.Body).Decode(&inputKategori); err != nil {
		w.WriteHeader(http.StatusBadRequest) // Error 400
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "Format inputan salah, pastikan pakai JSON yang bener",
			"data":    nil,
		})
		return
	}

	// 4. Simpan laci tadi ke database pakai mantra GORM: "Create"
	if err := config.DB.Create(&inputKategori).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Error 500
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "Gagal menyimpan data kategori ke database",
			"data":    nil,
		})
		return
	}

	// 5. Kalau sukses, kasih tau admin dan tampilin data yang berhasil disimpen
	w.WriteHeader(http.StatusCreated) // Status 201 (Created/Berhasil Dibuat)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Mantap! Kategori berhasil ditambahkan",
		"data":    inputKategori,
	})
}