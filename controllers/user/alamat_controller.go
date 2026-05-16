package user

import (
	"net/http"

	"backend-mantra/config"
	"backend-mantra/models"

	"github.com/gin-gonic/gin"
)

// GetAlamat mengambil semua alamat milik customer yang sedang login.
// Dipakai oleh: customer (GET /customer/alamat)
// Auth: Wajib login, role customer
func GetAlamat(c *gin.Context) {
	uid, exists := getUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User belum login",
			"error":   gin.H{"code": "AUTH_001", "detail": "Token tidak valid"},
		})
		return
	}

	var result struct{ IdCustomer uint }
	if err := config.DB.Raw("SELECT id_customer FROM customer WHERE id_user = ?", uid).Scan(&result).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengidentifikasi customer",
			"error":   gin.H{"code": "SERVER_001", "detail": err.Error()},
		})
		return
	}

	var alamats []models.Alamat
	if err := config.DB.Where("id_customer = ?", result.IdCustomer).Find(&alamats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil daftar alamat",
			"error":   gin.H{"code": "SERVER_001", "detail": err.Error()},
		})
		return
	}

	var data []gin.H
	for _, a := range alamats {
		data = append(data, gin.H{
			"id_alamat":        a.IdAlamat,
			"label_alamat":     a.LabelAlamat,
			"nama_penerima":    a.NamaPenerima,
			"no_telp_penerima": a.NoTelpPenerima,
			"alamat_lengkap":   a.AlamatLengkap,
			"latitude":         a.Latitude,
			"longitude":        a.Longitude,
			"catatan_lokasi":   a.CatatanLokasi,
			"is_utama":         a.IsUtama,
		})
	}

	if data == nil {
		data = []gin.H{}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Daftar alamat berhasil diambil",
		"data":    data,
	})
}

// TambahAlamat menambahkan alamat pengiriman baru untuk customer yang sedang login.
// Dipakai oleh: customer (POST /customer/alamat)
// Auth: Wajib login, role customer
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
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "error",
			"message": "Input tidak valid, pastikan field label, nama, telepon, dan alamat lengkap diisi",
			"error":   gin.H{"code": "VAL_001", "detail": "Validasi gagal"},
		})
		return
	}

	uid, exists := getUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User belum login",
			"error":   gin.H{"code": "AUTH_001", "detail": "Token tidak valid"},
		})
		return
	}

	var result struct{ IdCustomer uint }
	if err := config.DB.Raw("SELECT id_customer FROM customer WHERE id_user = ?", uid).Scan(&result).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengidentifikasi customer",
			"error":   gin.H{"code": "SERVER_001", "detail": err.Error()},
		})
		return
	}
	customerID := result.IdCustomer

	tx := config.DB.Begin()

	if input.IsUtama {
		if err := tx.Model(&models.Alamat{}).Where("id_customer = ?", customerID).Update("is_utama", false).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Gagal memperbarui status alamat utama sebelumnya",
				"error":   gin.H{"code": "SERVER_001", "detail": err.Error()},
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
			"error":   gin.H{"code": "SERVER_001", "detail": err.Error()},
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
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "error",
			"message": "Input tidak valid",
			"error":   gin.H{"code": "VAL_001", "detail": "Validasi gagal"},
		})
		return
	}

	uid, exists := getUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User belum login",
			"error":   gin.H{"code": "AUTH_001", "detail": "Token tidak valid"},
		})
		return
	}

	// Ownership check
	var count int
	err := config.DB.Raw(`
		SELECT COUNT(*) FROM alamat a
		JOIN customer c ON c.id_customer = a.id_customer
		WHERE a.id_alamat = ? AND c.id_user = ?
	`, idAlamat, uid).Scan(&count).Error

	if err != nil || count == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "Anda tidak memiliki akses ke resource ini",
			"error":   gin.H{"code": "AUTH_002", "detail": "Bukan milik user ini"},
		})
		return
	}

	var alamat models.Alamat
	if err := config.DB.First(&alamat, idAlamat).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Alamat tidak ditemukan",
			"error":   gin.H{"code": "DATA_004", "detail": "Data tidak ada di database"},
		})
		return
	}

	tx := config.DB.Begin()

	if input.IsUtama && !alamat.IsUtama {
		if err := tx.Model(&models.Alamat{}).Where("id_customer = ?", alamat.CustomerId).Update("is_utama", false).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Gagal memperbarui status alamat utama sebelumnya",
				"error":   gin.H{"code": "SERVER_001", "detail": err.Error()},
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
			"error":   gin.H{"code": "SERVER_001", "detail": err.Error()},
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
func HapusAlamat(c *gin.Context) {
	idAlamat := c.Param("id_alamat")

	uid, exists := getUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User belum login",
			"error":   gin.H{"code": "AUTH_001", "detail": "Token tidak valid"},
		})
		return
	}

	// Ownership check
	var count int
	err := config.DB.Raw(`
		SELECT COUNT(*) FROM alamat a
		JOIN customer c ON c.id_customer = a.id_customer
		WHERE a.id_alamat = ? AND c.id_user = ?
	`, idAlamat, uid).Scan(&count).Error

	if err != nil || count == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "Anda tidak memiliki akses ke resource ini",
			"error":   gin.H{"code": "AUTH_002", "detail": "Bukan milik user ini"},
		})
		return
	}

	if err := config.DB.Delete(&models.Alamat{}, idAlamat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menghapus alamat",
			"error":   gin.H{"code": "SERVER_001", "detail": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Alamat berhasil dihapus",
	})
}
