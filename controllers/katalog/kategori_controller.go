package katalog

import (
	"net/http"

	"backend-mantra/config"
	"backend-mantra/models"

	"github.com/gin-gonic/gin"
)

// GetKategori mengambil semua kategori barang.
// Dipakai oleh: customer (GET /customer/katalog/kategori), kasir (GET /kasir/katalog/kategori), admin (GET /admin/katalog/kategori)
// Auth: Wajib login, semua role boleh akses (dikontrol di route)
func GetKategori(c *gin.Context) {
	kategori := []models.Kategori{}

	if err := config.DB.Find(&kategori).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil daftar kategori",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berhasil mengambil daftar kategori",
		"data":    kategori,
	})
}

// TambahKategori menambahkan kategori baru.
// Dipakai oleh: admin (POST /admin/katalog/kategori)
// Auth: Wajib login, role admin
func TambahKategori(c *gin.Context) {
	var inputKategori models.Kategori

	if err := c.ShouldBindJSON(&inputKategori); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Format inputan salah, pastikan pakai JSON yang benar",
		})
		return
	}

	if err := config.DB.Create(&inputKategori).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menyimpan data kategori ke database",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Kategori berhasil ditambahkan",
		"data":    inputKategori,
	})
}
