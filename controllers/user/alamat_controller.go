package user

import (
	"net/http"

	"backend-mantra/config"
	"backend-mantra/models"

	"github.com/gin-gonic/gin"
)

// TambahAlamat menambahkan alamat pengiriman baru untuk customer yang sedang login.
// Dipakai oleh: customer (POST /customer/alamat)
// Auth: Wajib login, role customer
// Ownership: id_customer diambil dari JWT (user_id), bukan dari body request
func TambahAlamat(c *gin.Context) {
	type TambahAlamatInput struct {
		LabelAlamat    string  `json:"label_alamat" binding:"required"`
		NamaPenerima   string  `json:"nama_penerima" binding:"required"`
		NoTelpPenerima string  `json:"no_telp_penerima" binding:"required"`
		AlamatLengkap  string  `json:"alamat_lengkap" binding:"required"`
		Latitude       float64 `json:"latitude"`
		Longitude      float64 `json:"longitude"`
		CatatanLokasi  string  `json:"catatan_lokasi"`
		IsUtama        bool    `json:"is_utama"`
	}

	var input TambahAlamatInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Input tidak valid, pastikan field label, nama, telepon, dan alamat lengkap diisi",
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

	tx := config.DB.Begin()

	// Jika alamat baru diset sebagai UTAMA, matikan status is_utama di alamat lain
	if input.IsUtama {
		if err := tx.Model(&models.Alamat{}).Where("id_customer = ?", customerID).Update("is_utama", false).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Gagal memperbarui status alamat utama sebelumnya",
			})
			return
		}
	}

	newAlamat := models.Alamat{
		CustomerId:     customerID,
		LabelAlamat:    input.LabelAlamat,
		NamaPenerima:   input.NamaPenerima,
		NoTelpPenerima: input.NoTelpPenerima,
		AlamatLengkap:  input.AlamatLengkap,
		Latitude:       input.Latitude,
		Longitude:      input.Longitude,
		CatatanLokasi:  input.CatatanLokasi,
		IsUtama:        input.IsUtama,
	}

	if err := tx.Create(&newAlamat).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menambahkan alamat baru",
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Alamat baru berhasil ditambahkan",
		"data": gin.H{
			"id_alamat":    newAlamat.IdAlamat,
			"label_alamat": newAlamat.LabelAlamat,
			"is_utama":     newAlamat.IsUtama,
		},
	})
}

// UpdateAlamat memperbarui alamat pengiriman milik customer yang sedang login.
// Dipakai oleh: customer (PUT /customer/alamat/:id_alamat)
// Auth: Wajib login, role customer
// Ownership: alamat harus milik customer yang login (id_customer dari JWT)
func UpdateAlamat(c *gin.Context) {
	idAlamat := c.Param("id_alamat")

	type UpdateAlamatInput struct {
		LabelAlamat    string  `json:"label_alamat"`
		NamaPenerima   string  `json:"nama_penerima"`
		NoTelpPenerima string  `json:"no_telp_penerima"`
		AlamatLengkap  string  `json:"alamat_lengkap"`
		Latitude       float64 `json:"latitude"`
		Longitude      float64 `json:"longitude"`
		CatatanLokasi  string  `json:"catatan_lokasi"`
		IsUtama        bool    `json:"is_utama"`
	}

	var input UpdateAlamatInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Input tidak valid",
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

	var alamat models.Alamat
	if err := config.DB.Where("id_alamat = ? AND id_customer = ?", idAlamat, customerID).First(&alamat).Error; err != nil {
		// Ownership violation: selalu 403, bukan 404
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "Anda tidak memiliki akses ke resource ini",
			"error": gin.H{
				"code":   "AUTH_002",
				"detail": "Alamat ini bukan milik Anda",
			},
		})
		return
	}

	tx := config.DB.Begin()

	if input.IsUtama && !alamat.IsUtama {
		if err := tx.Model(&models.Alamat{}).Where("id_customer = ?", customerID).Update("is_utama", false).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Gagal memperbarui status alamat utama sebelumnya",
			})
			return
		}
	}

	if input.LabelAlamat != "" {
		alamat.LabelAlamat = input.LabelAlamat
	}
	if input.NamaPenerima != "" {
		alamat.NamaPenerima = input.NamaPenerima
	}
	if input.NoTelpPenerima != "" {
		alamat.NoTelpPenerima = input.NoTelpPenerima
	}
	if input.AlamatLengkap != "" {
		alamat.AlamatLengkap = input.AlamatLengkap
	}
	if input.Latitude != 0 {
		alamat.Latitude = input.Latitude
	}
	if input.Longitude != 0 {
		alamat.Longitude = input.Longitude
	}
	if input.CatatanLokasi != "" {
		alamat.CatatanLokasi = input.CatatanLokasi
	}
	alamat.IsUtama = input.IsUtama

	if err := tx.Save(&alamat).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memperbarui alamat",
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Alamat berhasil diperbarui",
		"data": gin.H{
			"id_alamat":    alamat.IdAlamat,
			"label_alamat": alamat.LabelAlamat,
			"is_utama":     alamat.IsUtama,
		},
	})
}

// HapusAlamat menghapus alamat pengiriman milik customer yang sedang login.
// Dipakai oleh: customer (DELETE /customer/alamat/:id_alamat)
// Auth: Wajib login, role customer
// Ownership: alamat harus milik customer yang login
func HapusAlamat(c *gin.Context) {
	idAlamat := c.Param("id_alamat")

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

	deleteResult := config.DB.Where("id_alamat = ? AND id_customer = ?", idAlamat, customerID).Delete(&models.Alamat{})

	if deleteResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menghapus alamat",
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
				"detail": "Alamat ini bukan milik Anda",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Alamat berhasil dihapus",
	})
}
