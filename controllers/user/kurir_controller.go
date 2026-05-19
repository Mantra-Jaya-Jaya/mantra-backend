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

// GetProfilKurir mengambil data profil kurir yang sedang login.
// Dipakai oleh: kurir (GET /kurir/profil)
// Auth: Wajib login, role kurir
func GetProfilKurir(c *gin.Context) {
	uid, exists := getUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User belum login",
			"error":   gin.H{"code": "AUTH_001", "detail": "Token tidak valid"},
		})
		return
	}

	var kurir models.Kurir
	if err := config.DB.Preload("User").Where("id_user = ?", uid).First(&kurir).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data kurir tidak ditemukan",
			"error":   gin.H{"code": "DATA_004", "detail": "Kurir tidak ditemukan di database"},
		})
		return
	}

	fotoProfil := kurir.User.FotoProfil
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
			"id_kurir":            kurir.IdKurir,
			"public_id":           kurir.PublicId,
			"no_telp":             kurir.NoTelp,
			"tempat_lahir":        kurir.TempatLahir,
			"tanggal_lahir":       kurir.TanggalLahir,
			"jenis_kelamin":       kurir.JenisKelamin,
			"alamat":              kurir.Alamat,
			"pendidikan_terakhir": kurir.PendidikanTerakhir,
			"nik":                 kurir.Nik,
			"id_user":             kurir.User.IdUser,
			"user_public_id":      kurir.User.PublicId,
			"username":            kurir.User.Username,
			"email":               kurir.User.Email,
			"nama_lengkap":        kurir.User.NamaLengkap,
			"foto_profil":         fotoProfil,
		},
	})
}

// GetDaftarKurir mengambil daftar semua kurir.
// Dipakai oleh: admin (GET /admin/user/kurir)
// Auth: Wajib login, role admin
func GetDaftarKurir(c *gin.Context) {
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

	var kurirs []models.Kurir
	var total int64

	query := config.DB.Model(&models.Kurir{}).Joins("User")

	if search != "" {
		query = query.Where("User.nama_lengkap ILIKE ? OR User.username ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&total)

	if err := query.Preload("User").Offset(offset).Limit(limit).Find(&kurirs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data",
			"error":   gin.H{"code": "SERVER_001", "detail": err.Error()},
		})
		return
	}

	var response []gin.H
	baseURL := os.Getenv("BASE_URL")

	for _, k := range kurirs {
		fotoProfil := k.User.FotoProfil
		if fotoProfil != "" && !strings.HasPrefix(fotoProfil, "http") && baseURL != "" {
			fotoProfil = strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(fotoProfil, "/")
		}

		response = append(response, gin.H{
			"id_kurir":     k.IdKurir,
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
		"message": "Daftar kurir berhasil diambil",
		"data":    response,
		"meta": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// TambahKurir membuat akun kurir baru.
// Dipakai oleh: admin (POST /admin/user/kurir)
// Auth: Wajib login, role admin
func TambahKurir(c *gin.Context) {
	type TambahKurirInput struct {
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

	var input TambahKurirInput
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
	if err := config.DB.Where("nama_role = ?", "Kurir").First(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Role tidak ditemukan",
			"error":   gin.H{"code": "SERVER_001", "detail": "Role Kurir tidak ada di database"},
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

	newKurir := models.Kurir{
		NoTelp:             input.NoTelp,
		TempatLahir:        input.TempatLahir,
		TanggalLahir:       tglLahir,
		JenisKelamin:       input.JenisKelamin,
		Alamat:             input.Alamat,
		PendidikanTerakhir: input.PendidikanTerakhir,
		Nik:                input.Nik,
		UserId:             newUser.IdUser,
	}

	if err := tx.Create(&newKurir).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal membuat data kurir",
			"error":   gin.H{"code": "SERVER_001", "detail": err.Error()},
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Kurir berhasil ditambahkan",
		"data": gin.H{
			"id_kurir":     newKurir.IdKurir,
			"nama_lengkap": newUser.NamaLengkap,
		},
	})
}

// GetDetailKurir mengambil detail profil satu kurir.
// Dipakai oleh: admin (GET /admin/user/kurir/:id_kurir)
// Auth: Wajib login, role admin
func GetDetailKurir(c *gin.Context) {
	idKurir := c.Param("id_kurir")

	var kurir models.Kurir
	if err := config.DB.Preload("User").Where("public_id = ?", idKurir).First(&kurir).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Detail kurir tidak ditemukan",
			"error":   gin.H{"code": "DATA_004", "detail": "Kurir tidak ada di database"},
		})
		return
	}

	fotoProfil := kurir.User.FotoProfil
	if fotoProfil != "" && !strings.HasPrefix(fotoProfil, "http") {
		baseURL := os.Getenv("BASE_URL")
		if baseURL != "" {
			fotoProfil = strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(fotoProfil, "/")
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Detail kurir berhasil diambil",
		"data": gin.H{
			"id_kurir":            kurir.PublicId,
			"public_id":           kurir.PublicId,
			"no_telp":             kurir.NoTelp,
			"tempat_lahir":        kurir.TempatLahir,
			"tanggal_lahir":       kurir.TanggalLahir,
			"jenis_kelamin":       kurir.JenisKelamin,
			"alamat":              kurir.Alamat,
			"pendidikan_terakhir": kurir.PendidikanTerakhir,
			"nik":                 kurir.Nik,
			"id_user":             kurir.User.IdUser,
			"username":            kurir.User.Username,
			"email":               kurir.User.Email,
			"nama_lengkap":        kurir.User.NamaLengkap,
			"foto_profil":         fotoProfil,
		},
	})
}

// UpdateKurir memperbarui data kurir.
// Dipakai oleh: admin (PUT /admin/user/kurir/:id_kurir)
// Auth: Wajib login, role admin
func UpdateKurir(c *gin.Context) {
	idKurir := c.Param("id_kurir")

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

	var kurir models.Kurir
	if err := config.DB.Preload("User").Where("public_id = ?", idKurir).First(&kurir).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data kurir tidak ditemukan",
			"error":   gin.H{"code": "DATA_004", "detail": "Kurir tidak ditemukan di database"},
		})
		return
	}

	if input.Email != "" && input.Email != kurir.User.Email {
		var existingUser models.User
		if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{
				"status":  "error",
				"message": "Email sudah terdaftar",
				"error":   gin.H{"code": "CONF_002", "detail": "Email telah digunakan"},
			})
			return
		}
		kurir.User.Email = input.Email
	}

	if input.NamaLengkap != "" {
		kurir.User.NamaLengkap = input.NamaLengkap
	}
	if input.FotoProfil != "" {
		kurir.User.FotoProfil = input.FotoProfil
	}
	if input.NoTelp != "" {
		kurir.NoTelp = input.NoTelp
	}
	if input.TempatLahir != "" {
		kurir.TempatLahir = input.TempatLahir
	}
	if input.TanggalLahir != "" {
		tglLahir, _ := time.Parse("2006-01-02", input.TanggalLahir)
		kurir.TanggalLahir = tglLahir
	}
	if input.JenisKelamin != "" {
		kurir.JenisKelamin = input.JenisKelamin
	}
	if input.Alamat != "" {
		kurir.Alamat = input.Alamat
	}
	if input.PendidikanTerakhir != "" {
		kurir.PendidikanTerakhir = input.PendidikanTerakhir
	}
	if input.Nik != "" {
		kurir.Nik = input.Nik
	}

	tx := config.DB.Begin()

	if err := tx.Save(&kurir).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memperbarui data kurir",
			"error":   gin.H{"code": "SERVER_001", "detail": err.Error()},
		})
		return
	}

	if err := tx.Save(&kurir.User).Error; err != nil {
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
		"message": "Data kurir berhasil diperbarui",
	})
}

// HapusKurir menghapus akun kurir.
// Dipakai oleh: admin (DELETE /admin/user/kurir/:id_kurir)
// Auth: Wajib login, role admin
func HapusKurir(c *gin.Context) {
	idKurir := c.Param("id_kurir")

	var kurir models.Kurir
	if err := config.DB.Where("public_id = ?", idKurir).First(&kurir).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data kurir tidak ditemukan",
			"error":   gin.H{"code": "DATA_004", "detail": "Kurir tidak ditemukan"},
		})
		return
	}

	tx := config.DB.Begin()

	// 1. Hapus kurir
	if err := tx.Where("id_kurir = ?", kurir.IdKurir).Delete(&models.Kurir{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menghapus kurir",
			"error":   gin.H{"code": "SERVER_001", "detail": err.Error()},
		})
		return
	}

	// 2. Hapus user
	if err := tx.Where("id_user = ?", kurir.UserId).Delete(&models.User{}).Error; err != nil {
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
		Where("id_user = ? AND revoked_at IS NULL", kurir.UserId).
		Update("revoked_at", &now).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal revoke token kurir",
			"error":   gin.H{"code": "SERVER_001", "detail": err.Error()},
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Kurir berhasil dihapus",
	})
}
