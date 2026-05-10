package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"backend-mantra/config"
	"backend-mantra/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Login handles user authentication
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /login [post]
func Login(c *gin.Context) {
	var req models.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Input tidak valid",
		})
		return
	}

	var user models.User
	// Cari user berdasarkan username atau email
	var query = config.DB.Preload("Role")
	if req.Username != "" {
		query = query.Where("username = ?", req.Username)
	} else if req.Email != "" {
		query = query.Where("email = ?", req.Email)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Username atau Email harus diisi",
		})
		return
	}

	if err := query.First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Username/Email atau password salah",
		})
		return
	}

	// Cek password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Username/Email atau password salah",
		})
		return
	}

	// Cari profile_id berdasarkan role
	var profileID uint
	roleName := user.Role.NamaRole

	switch roleName {
	case "Customer":
		var customer models.Customer
		if err := config.DB.Where("id_user = ?", user.IdUser).First(&customer).Error; err == nil {
			profileID = customer.IdCustomer
		}
	case "Kasir":
		var kasir models.Kasir
		if err := config.DB.Where("id_user = ?", user.IdUser).First(&kasir).Error; err == nil {
			profileID = kasir.IdKasir
		}
	case "Kurir":
		var kurir models.Kurir
		if err := config.DB.Where("id_user = ?", user.IdUser).First(&kurir).Error; err == nil {
			profileID = kurir.IdKurir
		}
	}

	// Mock token for now since JWT is not fully implemented yet in the project
	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.mock_token"

	// Generate random refresh token
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal generate refresh token",
		})
		return
	}
	refreshTokenStr := hex.EncodeToString(b)

	// Simpan refresh token ke database
	refreshToken := models.RefreshToken{
		Token:     refreshTokenStr,
		UserID:    user.IdUser,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // Berlaku 7 hari
		CreatedAt: time.Now(),
	}

	if err := config.DB.Create(&refreshToken).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menyimpan refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Login berhasil",
		"data": gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshTokenStr,
			"token_type":    "Bearer",
			"expires_in":    900,
			"user": gin.H{
				"id_user":      user.IdUser,
				"username":     user.Username,
				"email":        user.Email,
				"nama_lengkap": user.NamaLengkap,
				"role":         roleName,
				"profile_id":   profileID,
			},
		},
	})
}

// RefreshToken handles token refresh
func RefreshToken(c *gin.Context) {
	var req models.RefreshToken
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Input tidak valid",
		})
		return
	}

	var storedToken models.RefreshToken
	// Cari token di database
	if err := config.DB.Where("token = ?", req.Token).First(&storedToken).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Token tidak valid atau sudah kedaluwarsa",
		})
		return
	}

	// Cek apakah sudah expired
	if storedToken.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Token sudah kedaluwarsa",
		})
		return
	}

	// Cek apakah sudah direvoke
	if storedToken.RevokedAt != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Token sudah tidak berlaku",
		})
		return
	}

	// Jika valid, generate access token baru
	// (Karena belum ada JWT helper, kita pakai mock token baru)
	newAccessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.mock_new_token"

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Token berhasil diperbarui",
		"data": gin.H{
			"access_token": newAccessToken,
			"expires_in":   900,
		},
	})
}

// Logout handles user logout
func Logout(c *gin.Context) {
	var req models.RefreshToken
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Input tidak valid",
		})
		return
	}

	now := time.Now()
	// Update RevokedAt untuk token yang bersangkutan
	result := config.DB.Model(&models.RefreshToken{}).
		Where("token = ?", req.Token).
		Update("revoked_at", &now)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal melakukan logout",
		})
		return
	}

	// Jika tidak ada baris yang terupdate, berarti token tidak ditemukan
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Token tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Logout berhasil",
	})
}

// RegisterCustomer handles customer registration
func RegisterCustomer(c *gin.Context) {
	var req models.Customer
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Input tidak valid",
		})
		return
	}

	// Cek apakah email atau username sudah terdaftar
	var existingUser models.User
	if err := config.DB.Where("email = ? OR username = ?", req.User.Email, req.User.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"status":  "error",
			"message": "Username atau Email sudah terdaftar",
		})
		return
	}

	// Cari Role "Customer"
	var role models.Role
	if err := config.DB.Where("nama_role = ?", "Customer").First(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Role Customer tidak ditemukan di sistem",
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.User.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memproses password",
		})
		return
	}

	// Mulai transaksi
	tx := config.DB.Begin()

	req.User.Password = string(hashedPassword)
	req.User.RoleID = role.IdRole

	if err := tx.Create(&req.User).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal membuat user baru",
		})
		return
	}

	req.UserId = req.User.IdUser

	if err := tx.Create(&req).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal membuat data customer",
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Registrasi berhasil",
		"data": gin.H{
			"id_user":      req.User.IdUser,
			"username":     req.User.Username,
			"email":        req.User.Email,
			"nama_lengkap": req.User.NamaLengkap,
			"no_telp":      req.NoTelp,
			"role":         "Customer",
		},
	})
}

// ChangePassword handles password change
func ChangePassword(c *gin.Context) {
	// Struct lokal agar tidak mengotori models
	type ChangePasswordInput struct {
		PasswordLama string `json:"password_lama" binding:"required"`
		PasswordBaru string `json:"password_baru" binding:"required"`
	}

	var req ChangePasswordInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Input tidak valid",
		})
		return
	}

	// TODO: Ambil ID User dari context JWT jika middleware sudah aktif
	userID := uint(1)

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "User tidak ditemukan",
		})
		return
	}

	// Cek password lama
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.PasswordLama)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Password lama salah",
		})
		return
	}

	// Hash password baru
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.PasswordBaru), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memproses password baru",
		})
		return
	}

	// Update password
	user.Password = string(hashedPassword)
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengubah password",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Password berhasil diubah",
	})
}
