package notifikasi

import (
	"net/http"
	"time"

	"backend-mantra/config"
	"backend-mantra/models"

	"github.com/gin-gonic/gin"
)

// GetNotifikasi mengambil daftar notifikasi untuk user yang sedang login.
// Ownership: hanya notifikasi milik user yang login yang dikembalikan.
// Dipakai oleh: customer (GET /customer/notifikasi), kasir (GET /kasir/notifikasi), admin (GET /admin/notifikasi)
// Auth: Wajib login, semua role boleh akses (dikontrol di route)
func GetNotifikasi(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var notifikasis []models.Notifikasi
	if err := config.DB.Where("id_user = ?", userID).Find(&notifikasis).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil notifikasi",
		})
		return
	}

	if notifikasis == nil {
		notifikasis = []models.Notifikasi{}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Notifikasi berhasil diambil",
		"data":    notifikasis,
	})
}

// GetNotifikasiAdmin mengambil notifikasi khusus admin (stok menipis, dll).
// Dipakai oleh: admin (GET /admin/notifikasi)
// Auth: Wajib login, role admin
func GetNotifikasiAdmin(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var responseData []gin.H

	// 1. Ambil riwayat notifikasi admin dari database
	var notifikasis []models.Notifikasi
	if err := config.DB.Where("id_user = ?", userID).Find(&notifikasis).Error; err == nil {
		for _, n := range notifikasis {
			responseData = append(responseData, gin.H{
				"id_notifikasi": n.IdNotifikasi,
				"id_barang":     nil,
				"nama_barang":   nil,
				"varian":        nil,
				"stok_saat_ini": nil,
				"batas_minimum": nil,
				"pesan":         n.Pesan,
				"judul":         n.Judul,
				"status":        n.Status,
				"created_at":    time.Now().Add(-1 * time.Hour).Format(time.RFC3339), // Fallback time
			})
		}
	}

	// 2. Ambil barang-barang dengan stok menipis (jumlah <= 5) secara dinamis
	var lowStockItems []models.SpesifikasiBarang
	if err := config.DB.
		Preload("Barang").
		Preload("DetailSpesifikasi.Spesifikasi").
		Where("jumlah <= 5").
		Find(&lowStockItems).Error; err == nil {
		for _, item := range lowStockItems {
			varianName := ""
			if item.DetailSpesifikasi.Spesifikasi.NamaSpesifikasi != "" {
				varianName = item.DetailSpesifikasi.Spesifikasi.NamaSpesifikasi + " " + item.DetailSpesifikasi.NamaDetailSpesifikasi
			} else {
				varianName = "Default"
			}

			responseData = append(responseData, gin.H{
				"id_notifikasi": item.IdSpesifikasiBarang + 1000, // Offset agar tidak bentrok ID-nya
				"id_barang":     item.BarangID,
				"nama_barang":   item.Barang.NamaBarang,
				"varian":        varianName,
				"stok_saat_ini": item.Jumlah,
				"batas_minimum": 5,
				"pesan":         "Stok " + item.Barang.NamaBarang + " (" + varianName + ") hampir habis",
				"judul":         "Stok Menipis",
				"status":        "aktif",
				"created_at":    time.Now().Format(time.RFC3339),
			})
		}
	}

	if responseData == nil {
		responseData = []gin.H{}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Notifikasi admin berhasil diambil",
		"data":    responseData,
	})
}
