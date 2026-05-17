package katalog

import (
	"net/http"
	"strconv"
	"time"

	"backend-mantra/config"
	"backend-mantra/models"

	"github.com/gin-gonic/gin"
)

// GetPromo mengambil semua diskon/promo yang sedang aktif.
// Dipakai oleh: customer (GET /customer/katalog/promo), admin (GET /admin/katalog/diskon)
// Auth: Wajib login, semua role boleh akses (dikontrol di route)
func GetPromo(c *gin.Context) {
	var diskons []models.Diskon
	now := time.Now()

	// Ambil diskon yang aktif (tgl_mulai <= now <= tgl_selesai)
	if err := config.DB.Where("tgl_mulai <= ? AND tgl_selesai >= ?", now, now).Find(&diskons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data promo",
		})
		return
	}

	var responseData []gin.H
	for _, d := range diskons {
		responseData = append(responseData, gin.H{
			"id_diskon":   d.IdDiskon,
			"nama_diskon": d.NamaDiskon,
			"banner_url":  d.BannerDiskon,
			"tgl_selesai": d.TglSelesai,
		})
	}

	// Jika data kosong, pastikan return array kosong bukan null
	if responseData == nil {
		responseData = []gin.H{}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berhasil mengambil data promo",
		"data":    responseData,
	})
}

// TambahDiskon menambahkan diskon baru untuk barang.
// Dipakai oleh: admin (POST /admin/katalog/diskon)
// Auth: Wajib login, role admin
func TambahDiskon(c *gin.Context) {
	var input struct {
		NamaDiskon   string `json:"nama_diskon" binding:"required"`
		BesarDiskon  int    `json:"besar_diskon" binding:"required"`
		BannerDiskon string `json:"banner_diskon"`
		TglMulai     string `json:"tgl_mulai" binding:"required"` // Format: YYYY-MM-DD
		TglSelesai   string `json:"tgl_selesai" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Format inputan salah: " + err.Error(),
		})
		return
	}

	const layoutDate = "2006-01-02"
	tglMulai, err := time.Parse(layoutDate, input.TglMulai)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Format tgl_mulai tidak valid, gunakan YYYY-MM-DD",
		})
		return
	}
	tglSelesai, err := time.Parse(layoutDate, input.TglSelesai)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Format tgl_selesai tidak valid, gunakan YYYY-MM-DD",
		})
		return
	}

	diskon := models.Diskon{
		NamaDiskon:   input.NamaDiskon,
		BesarDiskon:  input.BesarDiskon,
		BannerDiskon: input.BannerDiskon,
		TglMulai:     tglMulai,
		TglSelesai:   tglSelesai,
	}

	if err := config.DB.Create(&diskon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menyimpan data diskon",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Diskon berhasil ditambahkan",
		"data": gin.H{
			"id_diskon":    diskon.IdDiskon,
			"nama_diskon":  diskon.NamaDiskon,
			"besar_diskon": diskon.BesarDiskon,
			"tgl_mulai":    diskon.TglMulai.Format(layoutDate),
			"tgl_selesai":  diskon.TglSelesai.Format(layoutDate),
		},
	})
}

// GetAllDiskon mengambil semua diskon (aktif maupun tidak aktif) untuk keperluan manajemen admin.
// Dipakai oleh: admin (GET /admin/katalog/diskon/semua)
// Auth: Wajib login, role admin
func GetAllDiskon(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	offset := (page - 1) * limit

	var diskons []models.Diskon
	var total int64

	config.DB.Model(&models.Diskon{}).Count(&total)

	if err := config.DB.Limit(limit).Offset(offset).Find(&diskons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data diskon",
		})
		return
	}

	now := time.Now()
	var responseData []gin.H
	for _, d := range diskons {
		aktif := d.TglMulai.Before(now) && d.TglSelesai.After(now)
		responseData = append(responseData, gin.H{
			"id_diskon":    d.IdDiskon,
			"nama_diskon":  d.NamaDiskon,
			"besar_diskon": d.BesarDiskon,
			"banner_url":   d.BannerDiskon,
			"tgl_mulai":    d.TglMulai,
			"tgl_selesai":  d.TglSelesai,
			"aktif":        aktif,
		})
	}

	if responseData == nil {
		responseData = []gin.H{}
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berhasil mengambil semua diskon",
		"data":    responseData,
		"meta": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// HapusDiskon menghapus diskon berdasarkan ID.
// Dipakai oleh: admin (DELETE /admin/katalog/diskon/:id_diskon)
// Auth: Wajib login, role admin
func HapusDiskon(c *gin.Context) {
	idDiskonStr := c.Param("id_diskon")
	idDiskon, err := strconv.Atoi(idDiskonStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "ID diskon tidak valid",
		})
		return
	}

	var diskon models.Diskon
	if err := config.DB.First(&diskon, "id_diskon = ?", idDiskon).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Diskon tidak ditemukan",
		})
		return
	}

	if err := config.DB.Delete(&diskon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menghapus diskon",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Diskon berhasil dihapus",
	})
}
