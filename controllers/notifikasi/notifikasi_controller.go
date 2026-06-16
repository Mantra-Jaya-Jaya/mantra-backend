package notifikasi

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"backend-mantra/config"
	"backend-mantra/models"
	"backend-mantra/utils"

	"github.com/gin-gonic/gin"
)

// GetNotifikasi mengambil daftar notifikasi untuk user yang sedang login.
// Ownership: hanya notifikasi milik user yang login yang dikembalikan.
// Dipakai oleh: customer (GET /customer/notifikasi), kasir (GET /kasir/notifikasi), admin (GET /admin/notifikasi)
// Auth: Wajib login, semua role boleh akses (dikontrol di route)
func GetNotifikasi(c *gin.Context) {
	// 1. Ambil dari context sebagai interface
	val, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "User ID tidak ditemukan"})
		return
	}

	// 2. Konversi ke uint secara eksplisit
	var userID uint
	switch v := val.(type) {
	case int:
		userID = uint(v)
	case int64:
		userID = uint(v)
	case uint:
		userID = v
	default:
		// Jika tipe data tidak dikenali, set ke 0
		userID = 0
	}

	// 3. Gunakan userID yang sudah pasti uint
	var notifikasis []models.Notifikasi
	if err := config.DB.Preload("StatusNotifikasiRel").Where("id_user = ?", userID).Find(&notifikasis).Error; err != nil {
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

	var notifikasis []models.Notifikasi
	if err := config.DB.Preload("StatusNotifikasiRel").Where("id_user = ?", userID).Find(&notifikasis).Error; err == nil {
		for _, n := range notifikasis {
			statusName := "unknown"
			if n.StatusNotifikasiRel != nil {
				statusName = n.StatusNotifikasiRel.NamaStatus
			}
			responseData = append(responseData, gin.H{
				"id_notifikasi": fmt.Sprintf("SYS-%d", n.IdNotifikasi),
				"id_barang":     nil,
				"public_id":     nil,
				"nama_barang":   nil,
				"varian":        nil,
				"stok_saat_ini": nil,
				"batas_minimum": nil,
				"pesan":         n.Pesan,
				"judul":         n.Judul,
				"status":        statusName,
				"created_at":    n.CreatedAt.Format(time.RFC3339),
			})
		}
	}

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
				"id_notifikasi": fmt.Sprintf("STK-%d", item.IdSpesifikasiBarang),
				"id_barang":     item.BarangID,
				"public_id":     item.Barang.PublicId,
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

// BacaNotifikasi menandai notifikasi sebagai "read".
// Dipakai oleh: admin (PATCH /admin/notifikasi/:id/baca)
func BacaNotifikasi(c *gin.Context) {
	userID := c.GetInt64("user_id")
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "ID notifikasi tidak valid"})
		return
	}

	result := config.DB.Model(&models.Notifikasi{}).
		Where("id_notifikasi = ? AND id_user = ?", id, userID).
		Update("id_status_notifikasi", utils.GetStatusNotifikasiID("read"))

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Notifikasi tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Notifikasi ditandai sudah dibaca"})
}

// HapusNotifikasi menghapus notifikasi dari database.
// Dipakai oleh: admin (DELETE /admin/notifikasi/:id)
func HapusNotifikasi(c *gin.Context) {
	userID := c.GetInt64("user_id")
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "ID notifikasi tidak valid"})
		return
	}

	result := config.DB.Where("id_notifikasi = ? AND id_user = ?", id, userID).Delete(&models.Notifikasi{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Notifikasi tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Notifikasi berhasil dihapus"})
}
