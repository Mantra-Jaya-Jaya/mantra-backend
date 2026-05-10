package controllers

import (
	"encoding/json"
	"net/http"
	"backend-mantra/config"
	"backend-mantra/models"
)

func GetPromoCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Harus GET
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "Method tidak diizinkan",
			"data":    nil,
		})
		return
	}

	// Siapin array kosong dari model Diskon 
	promos := []models.Diskon{}

	// Langsung tarik SEMUA data dari tabel diskon
	if err := config.DB.Find(&promos).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "Gagal mengambil data promo: koneksi database terputus",
			"data":    nil,
		})
		return
	}

	// Kalau datanya kosong di database (Error 404, tapi status success)
	if len(promos) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": "Saat ini tidak ada promo yang tersedia",
			"data":    promos, 
		})
		return
	}

	// Kalau berhasil narik data (200 OK)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Berhasil mengambil data promo",
		"data":    promos, 
	})
}