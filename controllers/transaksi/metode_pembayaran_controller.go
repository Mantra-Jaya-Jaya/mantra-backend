package transaksi

import (
	"net/http"
	"strings"

	"backend-mantra/config"
	"backend-mantra/models"

	"github.com/gin-gonic/gin"
)

func sanitizeIcon(kodeMetode, icon string) string {
	ic := strings.ToLower(strings.TrimSpace(icon))
	switch kodeMetode {
	case "cash":
		if ic == "" || ic == "payments_outlined" {
			return "payments_outlined"
		}
	case "cod":
		if ic == "" || strings.Contains(ic, "xendit") || strings.Contains(ic, "iconcod") {
			return "package_outlined"
		}
	case "qris":
		if ic == "" || strings.Contains(ic, "xendit") || strings.Contains(ic, "iconqris") {
			return "qr_code_scanner"
		}
	}
	return icon
}

func GetMetodePembayaran(c *gin.Context) {
	var metode []models.MetodePembayaran
	if err := config.DB.Order("urutan asc").Find(&metode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengambil data metode pembayaran"})
		return
	}
	if metode == nil {
		metode = []models.MetodePembayaran{}
	}
	for i := range metode {
		metode[i].Icon = sanitizeIcon(metode[i].KodeMetode, metode[i].Icon)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Daftar metode pembayaran berhasil diambil",
		"data":    metode,
	})
}

func GetMetodePembayaranAktif(c *gin.Context) {
	var metode []models.MetodePembayaran
	if err := config.DB.Where("is_active = ?", true).Order("urutan asc").Find(&metode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengambil data metode pembayaran"})
		return
	}
	if metode == nil {
		metode = []models.MetodePembayaran{}
	}
	for i := range metode {
		metode[i].Icon = sanitizeIcon(metode[i].KodeMetode, metode[i].Icon)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Daftar metode pembayaran berhasil diambil",
		"data":    metode,
	})
}

func TambahMetodePembayaran(c *gin.Context) {
	var input struct {
		NamaMetode string `json:"nama_metode" binding:"required"`
		KodeMetode string `json:"kode_metode" binding:"required"`
		Penyedia   string `json:"penyedia" binding:"required"`
		Icon       string `json:"icon"`
		Urutan     int    `json:"urutan"`
		IsActive   bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Input tidak valid"})
		return
	}

	mt := models.MetodePembayaran{
		NamaMetode: input.NamaMetode,
		KodeMetode: input.KodeMetode,
		Penyedia:   input.Penyedia,
		Icon:       input.Icon,
		Urutan:     input.Urutan,
		IsActive:   input.IsActive,
	}
	if err := config.DB.Create(&mt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menambah metode pembayaran"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Metode pembayaran berhasil ditambahkan",
		"data":    mt,
	})
}

func UpdateMetodePembayaran(c *gin.Context) {
	publicID := c.Param("public_id")
	var mt models.MetodePembayaran
	if err := config.DB.Where("public_id = ?", publicID).First(&mt).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Metode pembayaran tidak ditemukan"})
		return
	}

	var input struct {
		NamaMetode string `json:"nama_metode"`
		KodeMetode string `json:"kode_metode"`
		Penyedia   string `json:"penyedia"`
		Icon       string `json:"icon"`
		Urutan     *int   `json:"urutan"`
		IsActive   *bool  `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Input tidak valid"})
		return
	}

	if input.NamaMetode != "" {
		mt.NamaMetode = input.NamaMetode
	}
	if input.KodeMetode != "" {
		mt.KodeMetode = input.KodeMetode
	}
	if input.Penyedia != "" {
		mt.Penyedia = input.Penyedia
	}
	if input.Icon != "" {
		mt.Icon = input.Icon
	}
	if input.Urutan != nil {
		mt.Urutan = *input.Urutan
	}
	if input.IsActive != nil {
		mt.IsActive = *input.IsActive
	}

	if err := config.DB.Save(&mt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengupdate metode pembayaran"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Metode pembayaran berhasil diupdate",
		"data":    mt,
	})
}

func HapusMetodePembayaran(c *gin.Context) {
	publicID := c.Param("public_id")
	result := config.DB.Where("public_id = ?", publicID).Delete(&models.MetodePembayaran{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menghapus metode pembayaran"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Metode pembayaran tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Metode pembayaran berhasil dihapus",
	})
}
