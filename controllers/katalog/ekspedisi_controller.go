package katalog

import (
	"net/http"

	"backend-mantra/config"
	"backend-mantra/models"
	"backend-mantra/services"

	"github.com/gin-gonic/gin"
)

func GetDaftarEkspedisi(c *gin.Context) {
	var ekspedisi []models.Ekspedisi
	if err := config.DB.Preload("Layanan").Find(&ekspedisi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil daftar ekspedisi",
		})
		return
	}
	if ekspedisi == nil {
		ekspedisi = []models.Ekspedisi{}
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Daftar ekspedisi berhasil diambil",
		"data":    ekspedisi,
	})
}

func TambahEkspedisi(c *gin.Context) {
	var input struct {
		NamaEkspedisi string `json:"nama_ekspedisi" binding:"required"`
		KodeApi       string `json:"kode_api" binding:"required"`
		Logo          string `json:"logo"`
		Deskripsi     string `json:"deskripsi"`
		IsActive      bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Input tidak valid, pastikan nama_ekspedisi dan kode_api diisi",
		})
		return
	}

	eks := models.Ekspedisi{
		NamaEkspedisi: input.NamaEkspedisi,
		KodeApi:       input.KodeApi,
		Logo:          input.Logo,
		Deskripsi:     input.Deskripsi,
		IsActive:      input.IsActive,
	}
	if err := config.DB.Create(&eks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menambah ekspedisi",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Ekspedisi berhasil ditambahkan",
		"data":    eks,
	})
}



func HapusEkspedisi(c *gin.Context) {
	publicID := c.Param("public_id")
	result := config.DB.Where("public_id = ?", publicID).Delete(&models.Ekspedisi{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menghapus ekspedisi",
		})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Ekspedisi tidak ditemukan",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Ekspedisi berhasil dihapus",
	})
}

func TambahLayanan(c *gin.Context) {
	var input struct {
		EkspedisiID uint   `json:"id_ekspedisi" binding:"required"`
		NamaLayanan string `json:"nama_layanan" binding:"required"`
		Deskripsi   string `json:"deskripsi"`
		EstimasiMin int    `json:"estimasi_min"`
		EstimasiMax int    `json:"estimasi_max"`
		IsActive    bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Input tidak valid, pastikan id_ekspedisi dan nama_layanan diisi",
		})
		return
	}

	layanan := models.EkspedisiLayanan{
		EkspedisiID: input.EkspedisiID,
		NamaLayanan: input.NamaLayanan,
		Deskripsi:   input.Deskripsi,
		EstimasiMin: input.EstimasiMin,
		EstimasiMax: input.EstimasiMax,
		IsActive:    input.IsActive,
	}
	if err := config.DB.Create(&layanan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menambah layanan ekspedisi",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Layanan ekspedisi berhasil ditambahkan",
		"data":    layanan,
	})
}



func HapusLayanan(c *gin.Context) {
	publicID := c.Param("public_id")
	result := config.DB.Where("public_id = ?", publicID).Delete(&models.EkspedisiLayanan{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menghapus layanan ekspedisi",
		})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Layanan ekspedisi tidak ditemukan",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Layanan ekspedisi berhasil dihapus",
	})
}

func SyncBiteshipCouriers(c *gin.Context) {
	if err := services.SyncCouriersFromBiteship(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal sinkronisasi dari Biteship: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Sinkronisasi ekspedisi dan layanan Biteship berhasil diselesaikan",
	})
}
