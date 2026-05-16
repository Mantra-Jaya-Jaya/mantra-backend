package user

import (
	"net/http"

	"backend-mantra/config"
	"backend-mantra/models"

	"github.com/gin-gonic/gin"
)

// GetProfilCustomer mengambil data profil customer yang sedang login beserta daftar alamat.
// Dipakai oleh: customer (GET /customer/profil)
// Auth: Wajib login, role customer
// Ownership: id_customer diambil dari JWT (user_id), bukan dari request body/param
func GetProfilCustomer(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var customer models.Customer
	if err := config.DB.Preload("User").Where("id_user = ?", userID).First(&customer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data customer tidak ditemukan",
		})
		return
	}

	var alamat []models.Alamat
	config.DB.Where("id_customer = ?", customer.IdCustomer).Find(&alamat)

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

// UpdateAkunCustomer memperbarui informasi akun customer yang sedang login.
// Dipakai oleh: customer (PUT /customer/profil)
// Auth: Wajib login, role customer
// Ownership: id_customer diambil dari JWT (user_id), bukan dari request body/param
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

	userID := c.GetInt64("user_id")

	var customer models.Customer
	if err := config.DB.Preload("User").Where("id_user = ?", userID).First(&customer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data customer tidak ditemukan",
		})
		return
	}

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
