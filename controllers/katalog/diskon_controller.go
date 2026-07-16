package katalog

import (
	"net/http"
	"strconv"
	"time"

	"backend-mantra/config"
	"backend-mantra/models"
	"backend-mantra/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
			"public_id":   d.PublicId,
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
		TipeDiskon   string `json:"tipe_diskon" binding:"required"` // "persen" atau "nominal"
		BesarDiskon  int    `json:"besar_diskon" binding:"required"`
		BannerDiskon string `json:"banner_diskon"`
		TglMulai     string `json:"tgl_mulai" binding:"required"` // Format: YYYY-MM-DD
		TglSelesai   string `json:"tgl_selesai" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "error",
			"message": "Validasi gagal",
			"error":   err.Error(),
		})
		return
	}

	// Validasi tipe_diskon
	if input.TipeDiskon != "persen" && input.TipeDiskon != "nominal" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "tipe_diskon harus 'persen' atau 'nominal'",
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
		TipeDiskon:   input.TipeDiskon,
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
			"public_id":    diskon.PublicId,
			"nama_diskon":  diskon.NamaDiskon,
			"tipe_diskon":  diskon.TipeDiskon,
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
			"public_id":    d.PublicId,
			"nama_diskon":  d.NamaDiskon,
			"tipe_diskon":  d.TipeDiskon,
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
// Dipakai oleh: admin (DELETE /admin/diskon/:public_id)
// Auth: Wajib login, role admin
func HapusDiskon(c *gin.Context) {
	publicIdStr := c.Param("public_id")
	publicId, err := uuid.Parse(publicIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "ID diskon tidak valid (harus UUID)",
		})
		return
	}

	var diskon models.Diskon
	if err := config.DB.Where("public_id = ?", publicId).First(&diskon).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Diskon tidak ditemukan",
		})
		return
	}

	// PUTUS HUBUNGAN DULU (SET NULL)
	// Biar barang yang tadinya dapet diskon ini, balik ke harga normal, dan database gak error
	errLepasRelasi := config.DB.Model(&models.Barang{}).
		Where("id_diskon = ?", diskon.IdDiskon).
		Update("id_diskon", gorm.Expr("NULL")).Error

	if errLepasRelasi != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal melepaskan relasi promo dari barang",
		})
		return
	}

	// 🚀 BARU HAPUS DISKONNYA DENGAN AMAN
	if err := config.DB.Delete(&diskon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menghapus diskon",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Diskon berhasil dihapus dan barang terkait kembali ke harga normal",
	})
}

// UploadBannerDiskon mengunggah gambar banner promo/diskon ke MinIO.
// Dipakai oleh: admin (POST /admin/diskon/upload)
// Auth: Wajib login, role admin
func UploadBannerDiskon(c *gin.Context) {
	fileUrl, err := utils.UploadFileToMinio(c, "banner", "diskon")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Gagal mengunggah banner diskon: " + err.Error(),
		})
		return
	}

	// 2. Kembalikan URL publik MinIO ke Next.js
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Banner diskon berhasil diunggah ke server storage",
		"url":     fileUrl, // URL ini yang nanti disisipkan ke JSON TambahDiskon
	})
}
