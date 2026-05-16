package keranjang

import (
	"net/http"

	"backend-mantra/config"
	"backend-mantra/models"

	"github.com/gin-gonic/gin"
)

// TambahKeKeranjang menambahkan item ke keranjang belanja customer yang sedang login.
// Dipakai oleh: customer (POST /customer/keranjang)
// Auth: Wajib login, role customer
// Ownership: id_customer diambil dari JWT (user_id), bukan dari body request
func TambahKeKeranjang(c *gin.Context) {
	type TambahKeranjangInput struct {
		IdSpesifikasiBarang uint `json:"id_spesifikasi_barang" binding:"required"`
		Quantity            int  `json:"quantity" binding:"required"`
	}

	var input TambahKeranjangInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Input tidak valid, pastikan id_spesifikasi_barang dan quantity diisi",
		})
		return
	}

	userID := c.GetInt64("user_id")

	// Cari id_customer berdasarkan user_id dari JWT
	var result struct{ IdCustomer uint }
	if err := config.DB.Raw("SELECT id_customer FROM customer WHERE id_user = ?", userID).Scan(&result).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengidentifikasi customer",
		})
		return
	}
	customerID := result.IdCustomer

	var keranjang models.Keranjang
	err := config.DB.Where("id_customer = ? AND id_spesifikasi_barang = ?", customerID, input.IdSpesifikasiBarang).First(&keranjang).Error

	if err == nil {
		// Kalau sudah ada, tinggal tambahkan quantity-nya
		keranjang.Quantity += input.Quantity
		if err := config.DB.Save(&keranjang).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Gagal memperbarui item di keranjang",
			})
			return
		}
	} else {
		// Kalau belum ada, buat baris baru
		newItem := models.Keranjang{
			CustomerID:          customerID,
			SpesifikasiBarangID: input.IdSpesifikasiBarang,
			Quantity:            input.Quantity,
		}
		if err := config.DB.Create(&newItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Gagal menambahkan item ke keranjang",
			})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Barang berhasil ditambahkan ke keranjang",
	})
}

// UpdateKeranjang memperbarui quantity item di keranjang belanja.
// Dipakai oleh: customer (PATCH /customer/keranjang/:id_keranjang)
// Auth: Wajib login, role customer
// Ownership: item keranjang harus milik customer yang login
func UpdateKeranjang(c *gin.Context) {
	idKeranjang := c.Param("id_keranjang")

	type UpdateKeranjangInput struct {
		Quantity int `json:"quantity" binding:"required"`
	}

	var input UpdateKeranjangInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Input tidak valid, pastikan quantity diisi",
		})
		return
	}

	userID := c.GetInt64("user_id")

	var result struct{ IdCustomer uint }
	if err := config.DB.Raw("SELECT id_customer FROM customer WHERE id_user = ?", userID).Scan(&result).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengidentifikasi customer",
		})
		return
	}
	customerID := result.IdCustomer

	var keranjang models.Keranjang
	if err := config.DB.Where("id_keranjang = ? AND id_customer = ?", idKeranjang, customerID).First(&keranjang).Error; err != nil {
		// Ownership violation: selalu 403, bukan 404
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "Anda tidak memiliki akses ke resource ini",
			"error": gin.H{
				"code":   "AUTH_002",
				"detail": "Item keranjang ini bukan milik Anda",
			},
		})
		return
	}

	keranjang.Quantity = input.Quantity
	if err := config.DB.Save(&keranjang).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memperbarui item di keranjang",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Keranjang berhasil diperbarui",
	})
}

// HapusItemKeranjang menghapus item dari keranjang belanja customer yang login.
// Dipakai oleh: customer (DELETE /customer/keranjang/:id_keranjang)
// Auth: Wajib login, role customer
// Ownership: item keranjang harus milik customer yang login
func HapusItemKeranjang(c *gin.Context) {
	idKeranjang := c.Param("id_keranjang")

	userID := c.GetInt64("user_id")

	var result struct{ IdCustomer uint }
	if err := config.DB.Raw("SELECT id_customer FROM customer WHERE id_user = ?", userID).Scan(&result).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengidentifikasi customer",
		})
		return
	}
	customerID := result.IdCustomer

	deleteResult := config.DB.Where("id_keranjang = ? AND id_customer = ?", idKeranjang, customerID).Delete(&models.Keranjang{})

	if deleteResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menghapus item dari keranjang",
		})
		return
	}

	if deleteResult.RowsAffected == 0 {
		// Ownership violation: selalu 403, bukan 404
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "Anda tidak memiliki akses ke resource ini",
			"error": gin.H{
				"code":   "AUTH_002",
				"detail": "Item keranjang ini bukan milik Anda",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Item keranjang berhasil dihapus",
	})
}
