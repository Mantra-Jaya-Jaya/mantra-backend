package user

import (
	"net/http"
	"os"
	"strings"

	"backend-mantra/config"
	"backend-mantra/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetProfilKasir mengambil data profil kasir yang sedang login.
// Dipakai oleh: kasir (GET /kasir/profil)
// Auth: Wajib login, role kasir
func GetProfilKasir(c *gin.Context) {
	uid, exists := getUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "User belum login"})
		return
	}

	var kasir models.Kasir
	if err := config.DB.Joins("JOIN karyawan ON karyawan.id_karyawan = kasir.id_karyawan").
		Preload("Karyawan").Preload("Karyawan.User").
		Where("karyawan.id_user = ?", uid).First(&kasir).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Data kasir tidak ditemukan"})
		return
	}

	fotoProfil := kasir.Karyawan.User.FotoProfil
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
			"id_kasir":            kasir.IdKasir,
			"public_id":           kasir.Karyawan.PublicId,
			"no_telp":             kasir.Karyawan.NoTelp,
			"tempat_lahir":        kasir.Karyawan.TempatLahir,
			"tanggal_lahir":       kasir.Karyawan.TanggalLahir,
			"jenis_kelamin":       kasir.Karyawan.JenisKelamin,
			"alamat":              kasir.Karyawan.Alamat,
			"pendidikan_terakhir": kasir.Karyawan.PendidikanTerakhir,
			"nik":                 kasir.Karyawan.Nik,
			"id_user":             kasir.Karyawan.User.IdUser,
			"user_public_id":      kasir.Karyawan.User.PublicId,
			"username":            kasir.Karyawan.User.Username,
			"email":               kasir.Karyawan.User.Email,
			"nama_lengkap":        kasir.Karyawan.User.NamaLengkap,
			"foto_profil":         fotoProfil,
		},
	})
}

// UpdateProfilKasir memperbarui profil kasir yang sedang login.
// Dipakai oleh: kasir (PUT /kasir/profil)
// Auth: Wajib login, role kasir
func UpdateProfilKasir(c *gin.Context) {
	type UpdateProfilInput struct {
		NamaLengkap string `json:"nama_lengkap"`
		NoTelp      string `json:"no_telp"`
		Email       string `json:"email"`
		Alamat      string `json:"alamat"`
		FotoProfil  string `json:"foto_profil"`
	}

	var input UpdateProfilInput
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

	var kasir models.Kasir
	if err := config.DB.Joins("JOIN karyawan ON karyawan.id_karyawan = kasir.id_karyawan").
		Preload("Karyawan").Preload("Karyawan.User").
		Where("karyawan.id_user = ?", uid).First(&kasir).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data kasir tidak ditemukan",
			"error":   gin.H{"code": "DATA_004", "detail": "Kasir tidak ditemukan di database"},
		})
		return
	}

	// Cek duplikasi email jika email diubah
	if input.Email != "" && input.Email != kasir.Karyawan.User.Email {
		var existingUser models.User
		if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{
				"status":  "error",
				"message": "Email sudah terdaftar",
				"error":   gin.H{"code": "CONF_002", "detail": "Email telah digunakan"},
			})
			return
		}
		kasir.Karyawan.User.Email = input.Email
	}

	if input.NoTelp != "" {
		kasir.Karyawan.NoTelp = input.NoTelp
	}
	if input.NamaLengkap != "" {
		kasir.Karyawan.User.NamaLengkap = input.NamaLengkap
	}
	if input.Alamat != "" {
		kasir.Karyawan.Alamat = input.Alamat
	}
	if input.FotoProfil != "" {
		kasir.Karyawan.User.FotoProfil = input.FotoProfil
	}

	tx := config.DB.Begin()

	if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(&kasir.Karyawan).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memperbarui data kasir",
			"error":   gin.H{"code": "SERVER_001", "detail": err.Error()},
		})
		return
	}

	if err := tx.Save(&kasir.Karyawan.User).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memperbarui data user",
			"error":   gin.H{"code": "SERVER_001", "detail": err.Error()},
		})
		return
	}

	tx.Commit()

	fotoProfil := kasir.Karyawan.User.FotoProfil
	if fotoProfil != "" && !strings.HasPrefix(fotoProfil, "http") {
		baseURL := os.Getenv("BASE_URL")
		if baseURL != "" {
			fotoProfil = strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(fotoProfil, "/")
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Profil berhasil diperbarui",
		"data": gin.H{
			"nama_lengkap": kasir.Karyawan.User.NamaLengkap,
			"no_telp":      kasir.Karyawan.NoTelp,
			"email":        kasir.Karyawan.User.Email,
			"alamat":       kasir.Karyawan.Alamat,
			"foto_profil":  fotoProfil,
		},
	})
}
