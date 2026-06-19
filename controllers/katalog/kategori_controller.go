package katalog

import (
	"net/http"
	"strconv"

	"backend-mantra/config"
	"backend-mantra/models"
	"backend-mantra/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetKategori mengambil semua kategori barang.
// Dipakai oleh: customer (GET /customer/katalog/kategori), kasir (GET /kasir/katalog/kategori), admin (GET /admin/katalog/kategori)
// Auth: Wajib login, semua role boleh akses (dikontrol di route)
func GetKategori(c *gin.Context) {
	// 🚀 1. PENJINAK TOKEN (Ambil user_id dengan aman)
	val, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "User belum login"})
		return
	}

	var userID int64
	switch v := val.(type) {
	case float64:
		userID = int64(v)
	case int64:
		userID = v
	case int:
		userID = int64(v)
	case uint:
		userID = int64(v)
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Format token tidak valid"})
		return
	}

	// 🚀 2. INTELIJEN DATABASE: Cari Role User Langsung dari Tabel Master!
	var namaRole string
	// Asumsi tabel user lu namanya 'user' atau 'users', sesuaikan kalau beda ya bosku!
	errRole := config.DB.Raw("SELECT role.nama_role FROM user JOIN role ON role.id_role = user.id_role WHERE user.id_user = ?", userID).Scan(&namaRole).Error
	
	if errRole != nil || namaRole == "" {
		// Fallback (Jaga-jaga kalau query gagal, kita tetep coba ambil dari token JWT)
		namaRole = c.GetString("role")
	}

	kategori := []models.Kategori{}

	query := config.DB.
		Order("id_kategori ASC")

	// Admin lihat semua kategori; customer/kasir lihat hanya yang punya barang
	role, _ := c.Get("role")
	if roleStr, ok := role.(string); !ok || roleStr != "Admin" {
		query = query.Where("id_kategori IN (?)",
			config.DB.Table("barang").Select("DISTINCT id_kategori"),
		)
	}

	limitStr := c.Query("limit")
	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err == nil && limit > 0 {
			query = query.Limit(limit)
		}
	}

	// 7. EKSEKUSI QUERY
	if err := query.Find(&kategori).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil daftar kategori",
		})
		return
	}

	// 8. BUNGKUS DAN KIRIM!
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
	var input struct {
		NamaKategori string `json:"nama_kategori" binding:"required"`
		IconKategori string `json:"icon_kategori"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Format inputan salah atau data tidak lengkap: " + err.Error(),
		})
		return
	}

	kategori := models.Kategori{
		NamaKategori: input.NamaKategori,
		IconKategori: input.IconKategori,
	}

	if err := config.DB.Create(&kategori).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menyimpan data kategori ke database",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Kategori berhasil ditambahkan",
		"data":    kategori,
	})
}

// UpdateKategori memperbarui data kategori berdasarkan ID.
// Dipakai oleh: admin (PUT /admin/katalog/kategori/:public_id)
// Auth: Wajib login, role admin
func UpdateKategori(c *gin.Context) {
	publicIdStr := c.Param("public_id")
	publicId, err := uuid.Parse(publicIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "ID kategori tidak valid (harus UUID)",
		})
		return
	}

	var kategori models.Kategori
	if err := config.DB.Where("public_id = ?", publicId).First(&kategori).Error; err != nil {
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

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Tidak ada data yang diubah",
		})
		return
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
// Dipakai oleh: admin (DELETE /admin/katalog/kategori/:public_id)
// Auth: Wajib login, role admin
func HapusKategori(c *gin.Context) {
	publicIdStr := c.Param("public_id")
	publicId, err := uuid.Parse(publicIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "ID kategori tidak valid (harus UUID)",
		})
		return
	}

	var kategori models.Kategori
	if err := config.DB.Where("public_id = ?", publicId).First(&kategori).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Kategori tidak ditemukan",
		})
		return
	}

	// Cek apakah kategori masih digunakan oleh barang
	var count int64
	if err := config.DB.Model(&models.Barang{}).Where("id_kategori = ?", kategori.IdKategori).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memeriksa keterkaitan barang dengan kategori",
		})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Kategori tidak dapat dihapus karena masih digunakan oleh produk lain",
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

// UploadIconKategori mengunggah gambar/icon kategori ke MinIO.
// Dipakai oleh: admin (POST /admin/kategori/upload)
// Auth: Wajib login, role admin
func UploadIconKategori(c *gin.Context) {
	fileUrl, err := utils.UploadFileToMinio(c, "icon", "kategori")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Gagal mengunggah icon kategori: " + err.Error(),
		})
		return
	}

	// 2. Kembalikan URL publik MinIO ke Next.js
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Icon berhasil diunggah ke server storage",
		"url":     fileUrl, // URL ini yang nanti dikirim Next.js
	})
}
