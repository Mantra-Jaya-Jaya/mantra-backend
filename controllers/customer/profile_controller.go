package customer

import (
	"net/http"

	"backend-mantra/config"
	"backend-mantra/models"

	"github.com/gin-gonic/gin"
)

// GetProfilCustomer handles fetching customer profile
func GetProfilCustomer(c *gin.Context) {
	// TODO: Ambil ID Customer dari context JWT jika middleware sudah aktif
	customerID := uint(1)

	var customer models.Customer
	// Ambil data customer beserta user-nya
	if err := config.DB.Preload("User").Where("id_customer = ?", customerID).First(&customer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data customer tidak ditemukan",
		})
		return
	}

	var alamat []models.Alamat
	// Ambil semua alamat milik customer
	config.DB.Where("id_customer = ?", customerID).Find(&alamat)

	var responseAlamat []gin.H
	for _, a := range alamat {
		responseAlamat = append(responseAlamat, gin.H{
			"id_alamat":        a.IdAlamat,
			"label_alamat":     a.LabelAlamat,
			"nama_penerima":    a.NamaPenerima,
			"no_telp_penerima": a.NoTelpPenerima,
			"alamat_lengkap":   a.AlamatLengkap,
			"is_utama":         a.IsUtama,
		})
	}

	if responseAlamat == nil {
		responseAlamat = []gin.H{}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data profil berhasil diambil",
		"data": gin.H{
			"user": gin.H{
				"nama_lengkap": customer.User.NamaLengkap,
				"no_telp":      customer.NoTelp,
				"email":        customer.User.Email,
				"username":     customer.User.Username,
			},
			"daftar_alamat": responseAlamat,
		},
	})
}

// UpdateAkunCustomer handles updating customer account info
func UpdateAkunCustomer(c *gin.Context) {
	type UpdateAkunInput struct {
		NamaLengkap string `json:"nama_lengkap"`
		NoTelp      string `json:"no_telp"`
		Email       string `json:"email"`
		Username    string `json:"username"`
	}

	var input UpdateAkunInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Input tidak valid",
		})
		return
	}

	// TODO: Ambil ID Customer dari context JWT jika middleware sudah aktif
	customerID := uint(1)

	var customer models.Customer
	// Ambil data customer beserta user-nya
	if err := config.DB.Preload("User").Where("id_customer = ?", customerID).First(&customer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data customer tidak ditemukan",
		})
		return
	}

	// Update field jika dikirim di JSON
	if input.NoTelp != "" {
		customer.NoTelp = input.NoTelp
	}
	if input.NamaLengkap != "" {
		customer.User.NamaLengkap = input.NamaLengkap
	}
	if input.Email != "" {
		customer.User.Email = input.Email
	}
	if input.Username != "" {
		customer.User.Username = input.Username
	}

	// Mulai transaksi untuk mengupdate dua tabel
	tx := config.DB.Begin()

	if err := tx.Save(&customer).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memperbarui data customer",
		})
		return
	}

	if err := tx.Save(&customer.User).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memperbarui data user",
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Informasi akun berhasil diperbarui",
		"data": gin.H{
			"nama_lengkap": customer.User.NamaLengkap,
			"no_telp":      customer.NoTelp,
			"email":        customer.User.Email,
			"username":     customer.User.Username,
		},
	})
}

// TambahAlamat handles adding new address
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

	// TODO: Ambil ID Customer dari context JWT jika middleware sudah aktif
	customerID := uint(1)

	// Mulai transaksi
	tx := config.DB.Begin()

	// Jika alamat baru diset sebagai UTAMA, matikan status is_utama di alamat lain milik customer ini
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

// UpdateAlamat handles updating address
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

	// TODO: Ambil ID Customer dari context JWT jika middleware sudah aktif
	customerID := uint(1)

	var alamat models.Alamat
	// Cari alamat dan pastikan milik customer yang bersangkutan
	if err := config.DB.Where("id_alamat = ? AND id_customer = ?", idAlamat, customerID).First(&alamat).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Alamat tidak ditemukan",
		})
		return
	}

	// Mulai transaksi
	tx := config.DB.Begin()

	// Jika alamat ini diubah menjadi UTAMA, matikan status is_utama di alamat lain
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

	// Update field jika dikirim
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

// HapusAlamat handles deleting address
func HapusAlamat(c *gin.Context) {
	idAlamat := c.Param("id_alamat")

	// TODO: Ambil ID Customer dari context JWT jika middleware sudah aktif
	customerID := uint(1)

	// Hapus alamat jika cocok ID dan Customer ID
	result := config.DB.Where("id_alamat = ? AND id_customer = ?", idAlamat, customerID).Delete(&models.Alamat{})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menghapus alamat",
		})
		return
	}

	// Jika tidak ada baris yang terhapus, berarti alamat tidak ditemukan
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Alamat tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Alamat berhasil dihapus",
	})
}
