package auth

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
type LoginInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginInput
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Input tidak valid", "VAL_001", "Input tidak memenuhi aturan validasi")
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
		RespondWithError(c, http.StatusBadRequest, "Username atau Email harus diisi", "REQ_001", "Field username atau email kosong")
		return
	}

	if err := query.First(&user).Error; err != nil {
		RespondWithError(c, http.StatusUnauthorized, "Username atau password salah", "AUTH_001", "Credential tidak valid")
		return
	}

	// Cek password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		RespondWithError(c, http.StatusUnauthorized, "Username atau password salah", "AUTH_001", "Credential tidak valid")
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

	// Generate real JWT token
	accessToken, err := GenerateJWT(user.IdUser, user.PublicId.String(), roleName)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Gagal generate token", "SRV_001", err.Error())
		return
	}

	// Generate random refresh token
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Gagal generate refresh token", "SRV_002", err.Error())
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
		RespondWithError(c, http.StatusInternalServerError, "Gagal menyimpan refresh token", "SRV_003", err.Error())
		return
	}

	// Deteksi client type
	clientType := c.GetHeader("X-Client-Type")

	// Kirim response sukses sesuai client
	RespondWithSuccess(c, clientType, user, roleName, profileID, accessToken, refreshTokenStr)
}

// RefreshTokenInput is the DTO for refresh token request
type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshToken handles token refresh
func RefreshToken(c *gin.Context) {
	var req RefreshTokenInput
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Input tidak valid", "VAL_001", "Format request tidak sesuai")
		return
	}

	var storedToken models.RefreshToken
	// Cari token di database
	if err := config.DB.Where("token = ?", req.RefreshToken).First(&storedToken).Error; err != nil {
		RespondWithError(c, http.StatusUnauthorized, "Sesi Anda telah berakhir, silakan login kembali", "AUTH_003", "Refresh token tidak valid")
		return
	}

	// Cek apakah sudah expired
	if storedToken.ExpiresAt.Before(time.Now()) {
		RespondWithError(c, http.StatusUnauthorized, "Sesi Anda telah berakhir, silakan login kembali", "AUTH_003", "Refresh token expired")
		return
	}

	// Cek apakah sudah direvoke
	if storedToken.RevokedAt != nil {
		RespondWithError(c, http.StatusUnauthorized, "Sesi Anda telah berakhir, silakan login kembali", "AUTH_003", "Refresh token sudah direvoke")
		return
	}

	// Cari user dan role
	var user models.User
	if err := config.DB.Preload("Role").First(&user, storedToken.UserID).Error; err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Gagal memproses token", "SRV_001", "User tidak ditemukan")
		return
	}

	// Jika valid, generate access token baru
	newAccessToken, err := GenerateJWT(user.IdUser, user.PublicId.String(), user.Role.NamaRole)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Gagal generate token baru", "SRV_001", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Token berhasil diperbarui",
		"data": gin.H{
			"access_token": newAccessToken,
			"expires_in":   900,
		},
	})
}

// LogoutInput is the DTO for logout request
type LogoutInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Logout handles user logout
func Logout(c *gin.Context) {
	var req LogoutInput
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Input tidak valid", "VAL_001", "Format request tidak sesuai")
		return
	}

	now := time.Now()
	// Update RevokedAt untuk token yang bersangkutan
	result := config.DB.Model(&models.RefreshToken{}).
		Where("token = ?", req.RefreshToken).
		Update("revoked_at", &now)

	if result.Error != nil {
		RespondWithError(c, http.StatusInternalServerError, "Gagal melakukan logout", "SRV_001", result.Error.Error())
		return
	}

	// Jika tidak ada baris yang terupdate, berarti token tidak ditemukan
	if result.RowsAffected == 0 {
		RespondWithError(c, http.StatusUnauthorized, "Token tidak valid atau sudah expired", "AUTH_001", "Token tidak ditemukan atau sudah direvoke")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Logout berhasil",
	})
}

// RegisterCustomerInput is the DTO for register request
type RegisterCustomerInput struct {
	Username           string `json:"username" binding:"required"`
	Email              string `json:"email" binding:"required,email"`
	Password           string `json:"password" binding:"required,min=8"`
	KonfirmasiPassword string `json:"konfirmasi_password" binding:"required"`
	NamaLengkap        string `json:"nama_lengkap" binding:"required"`
	NoTelp             string `json:"no_telp" binding:"required"`
}

// RegisterCustomer handles customer registration
func RegisterCustomer(c *gin.Context) {
	var req RegisterCustomerInput
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusUnprocessableEntity, "Validasi gagal", "VAL_001", "Input tidak memenuhi aturan validasi")
		return
	}

	if req.Password != req.KonfirmasiPassword {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "error",
			"message": "Validasi gagal",
			"error": gin.H{
				"code":   "VAL_001",
				"detail": "Input tidak memenuhi aturan validasi",
			},
			"errors": gin.H{
				"konfirmasi_password": "Konfirmasi password tidak cocok",
			},
		})
		return
	}

	// Cek apakah email atau username sudah terdaftar
	var existingUser models.User
	if err := config.DB.Where("email = ? OR username = ?", req.Email, req.Username).First(&existingUser).Error; err == nil {
		RespondWithError(c, http.StatusConflict, "Username sudah terdaftar", "CONF_001", "Duplicate entry pada kolom username atau email")
		return
	}

	// Cari Role "Customer"
	var role models.Role
	if err := config.DB.Where("nama_role = ?", "Customer").First(&role).Error; err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Role Customer tidak ditemukan di sistem", "SRV_001", "Role tidak ditemukan")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Gagal memproses password", "SRV_001", err.Error())
		return
	}

	// Mulai transaksi
	tx := config.DB.Begin()

	newUser := models.User{
		Username:    req.Username,
		Email:       req.Email,
		Password:    string(hashedPassword),
		NamaLengkap: req.NamaLengkap,
		RoleID:      role.IdRole,
	}

	if err := tx.Create(&newUser).Error; err != nil {
		tx.Rollback()
		RespondWithError(c, http.StatusInternalServerError, "Gagal membuat user baru", "SRV_001", err.Error())
		return
	}

	newCustomer := models.Customer{
		NoTelp: req.NoTelp,
		UserId: newUser.IdUser,
	}

	if err := tx.Create(&newCustomer).Error; err != nil {
		tx.Rollback()
		RespondWithError(c, http.StatusInternalServerError, "Gagal membuat data customer", "SRV_001", err.Error())
		return
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Registrasi berhasil",
		"data": gin.H{
			"id_user":      newUser.IdUser,
			"username":     newUser.Username,
			"email":        newUser.Email,
			"nama_lengkap": newUser.NamaLengkap,
			"no_telp":      newCustomer.NoTelp,
			"role":         "customer",
		},
	})
}

// ChangePassword handles password change
func ChangePassword(c *gin.Context) {
	// Struct lokal agar tidak mengotori models
	type ChangePasswordInput struct {
		PasswordLama       string `json:"password_lama" binding:"required"`
		PasswordBaru       string `json:"password_baru" binding:"required"`
		KonfirmasiPassword string `json:"konfirmasi_password" binding:"required"`
	}

	var req ChangePasswordInput
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Input tidak valid", "VAL_001", "Format request tidak sesuai")
		return
	}

	if req.PasswordBaru != req.KonfirmasiPassword {
		RespondWithError(c, http.StatusUnprocessableEntity, "Konfirmasi password tidak cocok", "VAL_002", "password_baru dan konfirmasi_password tidak sama")
		return
	}

	// Ambil ID User dari context JWT middleware
	userID, exists := c.Get("user_id")
	if !exists {
		RespondWithError(c, http.StatusUnauthorized, "User belum login", "AUTH_001", "Token tidak valid")
		return
	}

	// Convert userID to uint
	uid, ok := userID.(uint)
	if !ok {
		// Sometimes numbers from context/json are float64 if not parsed explicitly as uint, 
		// but since we parse it in middleware as uint or float64 from JWT, we need to handle it.
		// In jwt-go, standard parsing returns float64 for numbers.
		if uidFloat, ok := userID.(float64); ok {
			uid = uint(uidFloat)
		} else {
			RespondWithError(c, http.StatusInternalServerError, "Kesalahan sistem", "SRV_001", "Format ID user tidak valid")
			return
		}
	}

	var user models.User
	if err := config.DB.First(&user, uid).Error; err != nil {
		RespondWithError(c, http.StatusNotFound, "User tidak ditemukan", "REQ_004", "Data user tidak ditemukan di database")
		return
	}

	// Cek password lama
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.PasswordLama)); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Password lama tidak sesuai", "REQ_003", "Password lama yang dimasukkan salah")
		return
	}

	// Hash password baru
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.PasswordBaru), bcrypt.DefaultCost)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Gagal memproses password baru", "SRV_001", err.Error())
		return
	}

	// Update password
	user.Password = string(hashedPassword)
	if err := config.DB.Save(&user).Error; err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Gagal mengubah password", "SRV_001", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Password berhasil diubah",
	})
}
