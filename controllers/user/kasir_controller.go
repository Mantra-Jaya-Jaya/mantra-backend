package user

import (
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"backend-mantra/config"
	"backend-mantra/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// GetProfilKasir mengambil data profil kasir yang sedang login.
// Dipakai oleh: kasir (GET /kasir/profil)
// Auth: Wajib login, role kasir
func GetProfilKasir(c *gin.Context) {
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
	if err := config.DB.Preload("User").Where("id_user = ?", uid).First(&kasir).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data kasir tidak ditemukan",
			"error":   gin.H{"code": "DATA_004", "detail": "Kasir tidak ditemukan di database"},
		})
		return
	}

	fotoProfil := kasir.User.FotoProfil
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
			"public_id":           kasir.PublicId,
			"no_telp":             kasir.NoTelp,
			"tempat_lahir":        kasir.TempatLahir,
			"tanggal_lahir":       kasir.TanggalLahir,
			"jenis_kelamin":       kasir.JenisKelamin,
			"alamat":              kasir.Alamat,
			"pendidikan_terakhir": kasir.PendidikanTerakhir,
			"nik":                 kasir.Nik,
			"id_user":             kasir.User.IdUser,
			"user_public_id":      kasir.User.PublicId,
			"username":            kasir.User.Username,
			"email":               kasir.User.Email,
			"nama_lengkap":        kasir.User.NamaLengkap,
			"foto_profil":         fotoProfil,
		},
	})
}

// GetDaftarKasir mengambil daftar semua kasir.
// Dipakai oleh: admin (GET /admin/user/kasir)
// Auth: Wajib login, role admin
func GetDaftarKasir(c *gin.Context) {
	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	var kasirs []models.Kasir
	var total int64

	query := config.DB.Model(&models.Kasir{}).Joins("User")

	if search != "" {
		query = query.Where("User.nama_lengkap ILIKE ? OR User.username ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&total)

	if err := query.Preload("User").Offset(offset).Limit(limit).Find(&kasirs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data",
			"error":   gin.H{"code": "SERVER_001", "detail": err.Error()},
		})
		return
	}

	var response []gin.H
	baseURL := os.Getenv("BASE_URL")
	
	for _, k := range kasirs {
		fotoProfil := k.User.FotoProfil
		if fotoProfil != "" && !strings.HasPrefix(fotoProfil, "http") && baseURL != "" {
			fotoProfil = strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(fotoProfil, "/")
		}

		response = append(response, gin.H{
			"id_kasir":     k.IdKasir,
			"nama_lengkap": k.User.NamaLengkap,
			"username":     k.User.Username,
			"foto_profil":  fotoProfil,
		})
	}
	if response == nil {
		response = []gin.H{}
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Daftar kasir berhasil diambil",
		"data":    response,
		"meta": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// TambahKasir membuat akun kasir baru.
// Dipakai oleh: admin (POST /admin/user/kasir)
// Auth: Wajib login, role admin
func TambahKasir(c *gin.Context) {
	type TambahKasirInput struct {
		Username           string `json:"username" binding:"required"`
		Email              string `json:"email" binding:"required"`
		Password           string `json:"password" binding:"required"`
		NamaLengkap        string `json:"nama_lengkap" binding:"required"`
		NoTelp             string `json:"no_telp"`
		TempatLahir        string `json:"tempat_lahir"`
		TanggalLahir       string `json:"tanggal_lahir"`
		JenisKelamin       string `json:"jenis_kelamin"`
		Alamat             string `json:"alamat"`
		PendidikanTerakhir string `json:"pendidikan_terakhir"`
		Nik                string `json:"nik"`
	}

	var input TambahKasirInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "error",
			"message": "Validasi gagal",
			"error":   gin.H{"code": "VAL_001", "detail": "Input tidak memenuhi aturan validasi"},
		})
		return
	}

	var existingUser models.User
	if err := config.DB.Where("username = ?", input.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"status":  "error",
			"message": "Username sudah terdaftar",
			"error":   gin.H{"code": "CONF_001", "detail": "Username duplikat"},
		})
		return
	}

	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"status":  "error",
			"message": "Email sudah terdaftar",
			"error":   gin.H{"code": "CONF_002", "detail": "Email duplikat"},
		})
		return
	}

	var role models.Role
	if err := config.DB.Where("nama_role = ?", "Kasir").First(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Role tidak ditemukan",
			"error":   gin.H{"code": "SERVER_001", "detail": "Role Kasir tidak ada di database"},
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memproses password",
			"error":   gin.H{"code": "SERVER_001", "detail": err.Error()},
		})
		return
	}

	tx := config.DB.Begin()
	newUser := models.User{
		Username:    input.Username,
		Email:       input.Email,
		Password:    string(hashedPassword),
		NamaLengkap: input.NamaLengkap,
		RoleID:      role.IdRole,
	}

	if err := tx.Create(&newUser).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal membuat data user",
			"error":   gin.H{"code": "SERVER_001", "detail": err.Error()},
		})
		return
	}

	var tglLahir time.Time
	if input.TanggalLahir != "" {
		tglLahir, _ = time.Parse("2006-01-02", input.TanggalLahir)
	}

	newKasir := models.Kasir{
		NoTelp:             input.NoTelp,
		TempatLahir:        input.TempatLahir,
		TanggalLahir:       tglLahir,
		JenisKelamin:       input.JenisKelamin,
		Alamat:             input.Alamat,
		PendidikanTerakhir: input.PendidikanTerakhir,
		Nik:                input.Nik,
		UserId:             newUser.IdUser,
	}

	if err := tx.Create(&newKasir).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal membuat data kasir",
			"error":   gin.H{"code": "SERVER_001", "detail": err.Error()},
		})
		return
	}
	
	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Kasir berhasil ditambahkan",
		"data": gin.H{
			"id_kasir":     newKasir.IdKasir,
			"nama_lengkap": newUser.NamaLengkap,
		},
	})
}

// GetDetailKasir mengambil detail profil satu kasir.
// Dipakai oleh: admin (GET /admin/user/kasir/:id_kasir)
// Auth: Wajib login, role admin
func GetDetailKasir(c *gin.Context) {
	idKasir := c.Param("id_kasir")

	var kasir models.Kasir
	if err := config.DB.Preload("User").Where("id_kasir = ?", idKasir).First(&kasir).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Detail kasir tidak ditemukan",
			"error":   gin.H{"code": "DATA_004", "detail": "Kasir tidak ada di database"},
		})
		return
	}

	fotoProfil := kasir.User.FotoProfil
	if fotoProfil != "" && !strings.HasPrefix(fotoProfil, "http") {
		baseURL := os.Getenv("BASE_URL")
		if baseURL != "" {
			fotoProfil = strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(fotoProfil, "/")
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Detail kasir berhasil diambil",
		"data": gin.H{
			"id_kasir":            kasir.IdKasir,
			"public_id":           kasir.PublicId,
			"no_telp":             kasir.NoTelp,
			"tempat_lahir":        kasir.TempatLahir,
			"tanggal_lahir":       kasir.TanggalLahir,
			"jenis_kelamin":       kasir.JenisKelamin,
			"alamat":              kasir.Alamat,
			"pendidikan_terakhir": kasir.PendidikanTerakhir,
			"nik":                 kasir.Nik,
			"id_user":             kasir.User.IdUser,
			"username":            kasir.User.Username,
			"email":               kasir.User.Email,
			"nama_lengkap":        kasir.User.NamaLengkap,
			"foto_profil":         fotoProfil,
		},
	})
}

// UpdateKasir memperbarui data kasir.
// Dipakai oleh: admin (PUT /admin/user/kasir/:id_kasir)
// Auth: Wajib login, role admin
func UpdateKasir(c *gin.Context) {
	idKasir := c.Param("id_kasir")

	var input struct {
		NamaLengkap        string `json:"nama_lengkap"`
		Email              string `json:"email"`
		NoTelp             string `json:"no_telp"`
		TempatLahir        string `json:"tempat_lahir"`
		TanggalLahir       string `json:"tanggal_lahir"`
		JenisKelamin       string `json:"jenis_kelamin"`
		Alamat             string `json:"alamat"`
		PendidikanTerakhir string `json:"pendidikan_terakhir"`
		Nik                string `json:"nik"`
		FotoProfil         string `json:"foto_profil"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "error",
			"message": "Validasi gagal",
			"error":   gin.H{"code": "VAL_001", "detail": "Input tidak valid"},
		})
		return
	}

	var kasir models.Kasir
	if err := config.DB.Preload("User").Where("id_kasir = ?", idKasir).First(&kasir).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data kasir tidak ditemukan",
			"error":   gin.H{"code": "DATA_004", "detail": "Kasir tidak ditemukan di database"},
		})
		return
	}

	if input.Email != "" && input.Email != kasir.User.Email {
		var existingUser models.User
		if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{
				"status":  "error",
				"message": "Email sudah terdaftar",
				"error":   gin.H{"code": "CONF_002", "detail": "Email telah digunakan"},
			})
			return
		}
		kasir.User.Email = input.Email
	}

	if input.NamaLengkap != "" {
		kasir.User.NamaLengkap = input.NamaLengkap
	}
	if input.FotoProfil != "" {
		kasir.User.FotoProfil = input.FotoProfil
	}
	if input.NoTelp != "" {
		kasir.NoTelp = input.NoTelp
	}
	if input.TempatLahir != "" {
		kasir.TempatLahir = input.TempatLahir
	}
	if input.TanggalLahir != "" {
		tglLahir, _ := time.Parse("2006-01-02", input.TanggalLahir)
		kasir.TanggalLahir = tglLahir
	}
	if input.JenisKelamin != "" {
		kasir.JenisKelamin = input.JenisKelamin
	}
	if input.Alamat != "" {
		kasir.Alamat = input.Alamat
	}
	if input.PendidikanTerakhir != "" {
		kasir.PendidikanTerakhir = input.PendidikanTerakhir
	}
	if input.Nik != "" {
		kasir.Nik = input.Nik
	}

	tx := config.DB.Begin()

	if err := tx.Save(&kasir).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memperbarui data kasir",
			"error":   gin.H{"code": "SERVER_001", "detail": err.Error()},
		})
		return
	}

	if err := tx.Save(&kasir.User).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memperbarui data user",
			"error":   gin.H{"code": "SERVER_001", "detail": err.Error()},
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data kasir berhasil diperbarui",
	})
}

// HapusKasir menghapus akun kasir.
// Dipakai oleh: admin (DELETE /admin/user/kasir/:id_kasir)
// Auth: Wajib login, role admin
func HapusKasir(c *gin.Context) {
	idKasir := c.Param("id_kasir")

	var kasir models.Kasir
	if err := config.DB.Where("id_kasir = ?", idKasir).First(&kasir).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data kasir tidak ditemukan",
			"error":   gin.H{"code": "DATA_004", "detail": "Kasir tidak ditemukan"},
		})
		return
	}

	tx := config.DB.Begin()

	// 1. Hapus kasir
	if err := tx.Where("id_kasir = ?", kasir.IdKasir).Delete(&models.Kasir{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menghapus kasir",
			"error":   gin.H{"code": "SERVER_001", "detail": err.Error()},
		})
		return
	}

	// 2. Hapus user
	if err := tx.Where("id_user = ?", kasir.UserId).Delete(&models.User{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menghapus user",
			"error":   gin.H{"code": "SERVER_001", "detail": err.Error()},
		})
		return
	}

	// 3. Revoke semua refresh token
	now := time.Now()
	if err := tx.Model(&models.RefreshToken{}).
		Where("id_user = ? AND revoked_at IS NULL", kasir.UserId).
		Update("revoked_at", &now).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal revoke token kasir",
			"error":   gin.H{"code": "SERVER_001", "detail": err.Error()},
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Kasir berhasil dihapus",
	})
}

// GetDaftarKaryawan mengambil semua karyawan (kasir + kurir).
// Dipakai oleh: admin (GET /admin/user/karyawan)
// Auth: Wajib login, role admin
func GetDaftarKaryawan(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Daftar karyawan belum diimplementasi sepenuhnya",
		"data": []gin.H{},
	})
}
