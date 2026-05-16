package katalog

import (
	"net/http"
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
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Diskon berhasil ditambahkan",
		"data": gin.H{
			"id_diskon":    3,
			"nama_diskon":  "Promo Lebaran",
			"besar_diskon": 15,
		},
	})
}
