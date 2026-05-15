package customer

import (
	"net/http"

	"backend-mantra/config"
	"backend-mantra/models"

	"github.com/gin-gonic/gin"
)

// GetNotifikasiCustomer handles fetching customer notifications
func GetNotifikasiCustomer(c *gin.Context) {
	// TODO: Ambil ID User dari context JWT jika middleware sudah aktif
	userID := uint(1)

	var notifikasis []models.Notifikasi
	// Ambil notifikasi milik user yang bersangkutan
	if err := config.DB.Where("id_user = ?", userID).Find(&notifikasis).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil notifikasi",
		})
		return
	}

	// Jika data kosong, pastikan return array kosong bukan null
	if notifikasis == nil {
		notifikasis = []models.Notifikasi{}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Notifikasi berhasil diambil",
		"data":    notifikasis,
	})
}
