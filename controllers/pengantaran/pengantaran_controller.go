package pengantaran

import (
	"net/http"

	"backend-mantra/config"
	"backend-mantra/models"

	"github.com/gin-gonic/gin"
)

// GetDaftarPengantaran mengambil daftar pengantaran yang sedang aktif.
// Dipakai oleh: admin (GET /admin/pengantaran), kurir (GET /kurir/pengantaran)
// Auth: Wajib login, role admin atau kurir (dikontrol di route)
func GetDaftarPengantaran(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetInt64("user_id")

	var pengantarans []models.Pengantaran
	query := config.DB.Preload("Pesanan").Preload("StatusPengantaran").Preload("Ekspedisi")

	// Filter jika yang mengakses adalah kurir
	if role == "kurir" {
		var result struct{ IdKurir uint }
		if err := config.DB.Raw("SELECT id_kurir FROM kurir WHERE id_user = ?", userID).Scan(&result).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Gagal mengidentifikasi kurir",
			})
			return
		}
		query = query.Where("id_kurir = ?", result.IdKurir)
	}

	if err := query.Find(&pengantarans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil daftar pengantaran",
		})
		return
	}

	if pengantarans == nil {
		pengantarans = []models.Pengantaran{}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Daftar pengantaran berhasil diambil",
		"data":    pengantarans,
	})
}

// UpdateLokasiKurir memperbarui koordinat lokasi kurir yang sedang bertugas.
// Dipakai oleh: kurir (PATCH /kurir/pengantaran/:id_pengantaran/lokasi)
// Auth: Wajib login, role kurir
// Ownership: kurir hanya bisa update lokasi pengantaran yang ditugaskan kepadanya
func UpdateLokasiKurir(c *gin.Context) {
	idPengantaran := c.Param("id_pengantaran")
	userID := c.GetInt64("user_id")

	type UpdateLokasiInput struct {
		Latitude  float64 `json:"latitude" binding:"required"`
		Longitude float64 `json:"longitude" binding:"required"`
	}

	var input UpdateLokasiInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Format inputan salah, pastikan latitude dan longitude diisi dengan benar",
		})
		return
	}

	// 1. Cari id_kurir berdasarkan user_id dari JWT
	var result struct{ IdKurir uint }
	if err := config.DB.Raw("SELECT id_kurir FROM kurir WHERE id_user = ?", userID).Scan(&result).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengidentifikasi kurir",
		})
		return
	}
	kurirID := result.IdKurir

	// 2. Cari data pengantaran dan periksa kepemilikan (ownership check)
	var pengantaran models.Pengantaran
	if err := config.DB.Where("public_id = ?", idPengantaran).First(&pengantaran).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data pengantaran tidak ditemukan",
		})
		return
	}

	if pengantaran.KurirID != kurirID {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "Anda tidak memiliki akses ke resource ini",
			"error": gin.H{
				"code":   "AUTH_002",
				"detail": "Pengantaran ini tidak ditugaskan kepada Anda",
			},
		})
		return
	}

	// 3. Update koordinat
	pengantaran.LastLatitude = input.Latitude
	pengantaran.LastLongitude = input.Longitude

	if err := config.DB.Save(&pengantaran).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memperbarui lokasi kurir",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Lokasi kurir berhasil diperbarui",
	})
}
