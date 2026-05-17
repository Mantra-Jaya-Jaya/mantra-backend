package katalog

import (
	"net/http"
	"strconv"

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

// UpdateKategori memperbarui data kategori berdasarkan ID.
// Dipakai oleh: admin (PUT /admin/katalog/kategori/:id_kategori)
// Auth: Wajib login, role admin
func UpdateKategori(c *gin.Context) {
	idStr := c.Param("id_kategori")
	idKategori, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "ID kategori tidak valid",
		})
		return
	}

	var kategori models.Kategori
	if err := config.DB.First(&kategori, "id_kategori = ?", idKategori).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Kategori tidak ditemukan",
		})
		return
	}

	var input struct {
		NamaKategori string `json:"nama_kategori"`
		IconKategori string `json:"icon_kategori"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Format inputan salah: " + err.Error(),
		})
		return
	}

	updates := map[string]interface{}{}
	if input.NamaKategori != "" {
		updates["nama_kategori"] = input.NamaKategori
	}
	if input.IconKategori != "" {
		updates["icon_kategori"] = input.IconKategori
	}

	if err := config.DB.Model(&kategori).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memperbarui data kategori",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Kategori berhasil diperbarui",
		"data":    kategori,
	})
}

// HapusKategori menghapus kategori berdasarkan ID.
// Dipakai oleh: admin (DELETE /admin/katalog/kategori/:id_kategori)
// Auth: Wajib login, role admin
func HapusKategori(c *gin.Context) {
	idStr := c.Param("id_kategori")
	idKategori, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "ID kategori tidak valid",
		})
		return
	}

	var kategori models.Kategori
	if err := config.DB.First(&kategori, "id_kategori = ?", idKategori).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Kategori tidak ditemukan",
		})
		return
	}

	if err := config.DB.Delete(&kategori).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menghapus kategori",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Kategori berhasil dihapus",
	})
}
