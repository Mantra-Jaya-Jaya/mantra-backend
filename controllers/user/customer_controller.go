package user

import (
	"net/http"
	"os"
	"strings"

	"backend-mantra/config"
	"backend-mantra/models"
	"backend-mantra/utils"

	"github.com/gin-gonic/gin"
)

// Helper internal untuk mengekstrak user_id dari JWT context secara aman
func getUserIDFromContext(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	switch v := userID.(type) {
	case uint:
		return v, true
	case int64:
		return uint(v), true
	case float64:
		return uint(v), true
	}

	return 0, false
}

// GetProfilCustomer mengambil data profil customer yang sedang login.
// Dipakai oleh: customer (GET /customer/profil)
// Auth: Wajib login, role customer
func GetProfilCustomer(c *gin.Context) {
	uid, exists := getUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User belum login",
			"error":   gin.H{"code": "AUTH_001", "detail": "Token tidak valid"},
		})
		return
	}

	var customer models.Customer
	if err := config.DB.Preload("User").Where("id_user = ?", uid).First(&customer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data customer tidak ditemukan",
			"error":   gin.H{"code": "DATA_004", "detail": "Customer tidak ditemukan di database"},
		})
		return
	}

	// Build absolute URL for foto_profil
	fotoProfil := customer.User.FotoProfil
	if fotoProfil != "" && !strings.HasPrefix(fotoProfil, "http") {
		baseURL := os.Getenv("BASE_URL")
		if baseURL != "" {
			fotoProfil = strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(fotoProfil, "/")
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data profil berhasil diambil",
		"data": gin.H{
			"id_customer":    customer.IdCustomer,
			"public_id":      customer.PublicId,
			"no_telp":        customer.NoTelp,
			"id_user":        customer.User.IdUser,
			"user_public_id": customer.User.PublicId,
			"username":       customer.User.Username,
			"email":          customer.User.Email,
			"nama_lengkap":   customer.User.NamaLengkap,
			"foto_profil":    fotoProfil,
		},
	})
}

// EditAkunCustomer memperbarui informasi akun customer yang sedang login.
// Dipakai oleh: customer (PUT /customer/akun)
// Auth: Wajib login, role customer
func EditAkunCustomer(c *gin.Context) {
	type EditAkunInput struct {
		NamaLengkap string `json:"nama_lengkap"`
		NoTelp      string `json:"no_telp"`
		Email       string `json:"email"`
		FotoProfil  string `json:"foto_profil"`
	}

	var input EditAkunInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "error",
			"message": "Validasi gagal",
			"error":   gin.H{"code": "VAL_001", "detail": "Input tidak memenuhi aturan validasi"},
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

	var customer models.Customer
	if err := config.DB.Preload("User").Where("id_user = ?", uid).First(&customer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data customer tidak ditemukan",
			"error":   gin.H{"code": "DATA_004", "detail": "Customer tidak ditemukan di database"},
		})
		return
	}

	// Cek duplikasi email jika email diubah
	if input.Email != "" && input.Email != customer.User.Email {
		var existingUser models.User
		if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{
				"status":  "error",
				"message": "Email sudah terdaftar",
				"error":   gin.H{"code": "CONF_002", "detail": "Email telah digunakan"},
			})
			return
		}
		customer.User.Email = input.Email
	}

	if input.NoTelp != "" {
		customer.NoTelp = input.NoTelp
	}
	if input.NamaLengkap != "" {
		customer.User.NamaLengkap = input.NamaLengkap
	}
	if input.FotoProfil != "" {
		customer.User.FotoProfil = input.FotoProfil
	}

	tx := config.DB.Begin()

	if err := tx.Save(&customer).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memperbarui data customer",
			"error":   gin.H{"code": "SERVER_001", "detail": err.Error()},
		})
		return
	}

	if err := tx.Save(&customer.User).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memperbarui data user",
			"error":   gin.H{"code": "SERVER_001", "detail": err.Error()},
		})
		return
	}

	tx.Commit()

	fotoProfil := customer.User.FotoProfil
	if fotoProfil != "" && !strings.HasPrefix(fotoProfil, "http") {
		baseURL := os.Getenv("BASE_URL")
		if baseURL != "" {
			fotoProfil = strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(fotoProfil, "/")
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Informasi akun berhasil diperbarui",
		"data": gin.H{
			"nama_lengkap": customer.User.NamaLengkap,
			"no_telp":      customer.NoTelp,
			"email":        customer.User.Email,
			"foto_profil":  fotoProfil,
		},
	})
}

// UploadFotoProfilCustomer mengunggah foto profil ke MinIO dan memperbarui database.
// Dipakai oleh: customer (POST /customer/profil/foto)
// Auth: Wajib login, role customer
func UploadFotoProfilCustomer(c *gin.Context) {
	uid, exists := getUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User belum login",
			"error":   gin.H{"code": "AUTH_001", "detail": "Token tidak valid"},
		})
		return
	}

	// Cari customer berdasarkan id_user
	var customer models.Customer
	if err := config.DB.Preload("User").Where("id_user = ?", uid).First(&customer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data customer tidak ditemukan",
			"error":   gin.H{"code": "DATA_004", "detail": "Customer tidak ditemukan di database"},
		})
		return
	}

	// Gunakan utils untuk upload file (parameter nama field form-data: "foto", folder: "profil")
	// Pastikan kita sudah import "backend-mantra/utils" jika belum (tapi file ini package user, jadi mungkin butuh utils)
	// Wait, we need to import utils! I will ensure utils is imported if not already.
	// Ah, I need to check if 'backend-mantra/utils' is imported in customer_controller.go.
	// We'll see.
	fotoURL, err := utils.UploadFileToMinio(c, "foto", "profil")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Gagal mengunggah foto profil",
			"error":   gin.H{"code": "UPLOAD_001", "detail": err.Error()},
		})
		return
	}

	// Update foto profil di tabel user
	customer.User.FotoProfil = fotoURL
	if err := config.DB.Save(&customer.User).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menyimpan foto profil ke database",
			"error":   gin.H{"code": "SERVER_001", "detail": err.Error()},
		})
		return
	}

	// Build absolute URL for response
	fotoProfil := fotoURL
	if fotoProfil != "" && !strings.HasPrefix(fotoProfil, "http") {
		baseURL := os.Getenv("BASE_URL")
		if baseURL != "" {
			fotoProfil = strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(fotoProfil, "/")
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Foto profil berhasil diperbarui",
		"data": gin.H{
			"foto_profil": fotoProfil,
		},
	})
}
