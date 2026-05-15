package notifikasi

import (
	"net/http"

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
// Dipakai oleh: admin (GET /admin/notifikasi/stok)
// Auth: Wajib login, role admin
func GetNotifikasiAdmin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Notifikasi admin berhasil diambil",
		"data": []gin.H{
			{
				"id_notifikasi": 1,
				"id_barang":     1,
				"nama_barang":   "Laptop Gaming X",
				"varian":        "16GB RAM",
				"stok_saat_ini": 3,
				"batas_minimum": 5,
				"pesan":         "Stok Laptop Gaming X (16GB RAM) hampir habis",
				"created_at":    "2026-05-09T08:00:00Z",
			},
		},
	})
}
