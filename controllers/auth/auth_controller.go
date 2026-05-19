package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/mail"
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
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginInput
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Input tidak valid", "VAL_001", "Input tidak memenuhi aturan validasi")
		return
	}

	fmt.Printf("[DEBUG] Login Attempt - Username: '%s', Password: '%s'\n", req.Username, req.Password)

	var user models.User
	// Cari user berdasarkan username atau email
	if err := config.DB.Preload("Role").Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error; err != nil {
		fmt.Printf("[DEBUG] User not found: %v\n", err)
		RespondWithError(c, http.StatusUnauthorized, "Username/Email atau password salah", "AUTH_001", "Credential tidak valid")
		return
	}

	fmt.Printf("[DEBUG] User Found - Stored Hash: '%s'\n", user.Password)

	// Cek password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		fmt.Printf("[DEBUG] Password Mismatch: %v\n", err)
		RespondWithError(c, http.StatusUnauthorized, "Username/Email atau password salah", "AUTH_001", "Credential tidak valid")
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
		RespondWithError(c, http.StatusInternalServerError, "Gagal generate token", "SERVER_001", err.Error())
		return
	}

	// Generate random refresh token
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Gagal generate refresh token", "SERVER_001", err.Error())
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
		RespondWithError(c, http.StatusInternalServerError, "Gagal menyimpan refresh token", "SERVER_001", err.Error())
		return
	}

	// Deteksi client type berdasarkan header khusus (lebih andal dari Authorization)
	clientType := "nextjs"
	if c.GetHeader("X-Client-Type") == "flutter" {
		clientType = "flutter"
	}

	// Kirim response sukses sesuai client
	RespondWithSuccess(c, clientType, user, roleName, profileID, accessToken, refreshTokenStr)
}

// RefreshTokenInput is the DTO for refresh token request
type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshToken handles token refresh
func RefreshToken(c *gin.Context) {
	var req RefreshTokenInput
	_ = c.ShouldBindJSON(&req)

	tokenStr := req.RefreshToken
	clientType := "flutter"

	if c.GetHeader("X-Client-Type") == "nextjs" {
		clientType = "nextjs"
	}

	if tokenStr == "" || tokenStr == "token_kosong" {
		cookieToken, err := c.Cookie("refresh_token")
		if err == nil && cookieToken != "" {
			tokenStr = cookieToken
			clientType = "nextjs"
		}
	}

	if tokenStr == "" {
		RespondWithError(c, http.StatusBadRequest, "Input tidak valid", "VAL_001", "Refresh token tidak ditemukan")
		return
	}

	var storedToken models.RefreshToken
	// Cari token di database
	if err := config.DB.Where("token = ? AND expires_at > ? AND revoked_at IS NULL", tokenStr, time.Now()).First(&storedToken).Error; err != nil {
		RespondWithError(c, http.StatusUnauthorized, "Sesi Anda telah berakhir, silakan login kembali", "AUTH_003", "Refresh token tidak valid atau expired")
		return
	}

	// Cari user dan role
	var user models.User
	if err := config.DB.Preload("Role").First(&user, storedToken.UserID).Error; err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Gagal memproses token", "SERVER_001", "User tidak ditemukan")
		return
	}

	// Jika valid, generate access token baru
	newAccessToken, err := GenerateJWT(user.IdUser, user.PublicId.String(), user.Role.NamaRole)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Gagal generate token baru", "SERVER_001", err.Error())
		return
	}

	if clientType == "nextjs" {
		c.SetCookie("access_token", newAccessToken, 900, "/", "", true, true)
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Token berhasil diperbarui",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Token berhasil diperbarui",
			"data": gin.H{
				"access_token": newAccessToken,
				"expires_in":   900,
			},
		})
	}
}

// LogoutInput is the DTO for logout request
type LogoutInput struct {
	RefreshToken string `json:"refresh_token"`
}

// Logout handles user logout
func Logout(c *gin.Context) {
	var req LogoutInput
	_ = c.ShouldBindJSON(&req)

	tokenStr := req.RefreshToken
	clientType := "flutter"

	if c.GetHeader("X-Client-Type") == "nextjs" {
		clientType = "nextjs"
	}

	if tokenStr == "" || tokenStr == "token_kosong" {
		cookieToken, err := c.Cookie("refresh_token")
		if err == nil && cookieToken != "" {
			tokenStr = cookieToken
			clientType = "nextjs"
		}
	}

	if tokenStr == "" && clientType != "nextjs" {
		RespondWithError(c, http.StatusBadRequest, "Input tidak valid", "VAL_001", "Refresh token tidak ditemukan")
		return
	}

	// Ambil ID User dari context JWT middleware
	userID, exists := c.Get("user_id")
	if !exists {
		RespondWithError(c, http.StatusUnauthorized, "User belum login", "AUTH_001", "Token tidak valid")
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		if uidFloat, ok := userID.(float64); ok {
			uid = uint(uidFloat)
		} else {
			RespondWithError(c, http.StatusInternalServerError, "Kesalahan sistem", "SERVER_001", "Format ID user tidak valid")
			return
		}
	}

	if tokenStr != "" {
		now := time.Now()
		// Update RevokedAt untuk token yang bersangkutan & milik user tsb
		result := config.DB.Model(&models.RefreshToken{}).
			Where("token = ? AND id_user = ?", tokenStr, uid).
			Update("revoked_at", &now)

		if result.Error != nil {
			RespondWithError(c, http.StatusInternalServerError, "Gagal melakukan logout", "SERVER_001", result.Error.Error())
			return
		}
	}

	if clientType == "nextjs" {
		c.SetCookie("access_token", "", -1, "/", "", true, true)
		c.SetCookie("refresh_token", "", -1, "/", "", true, true)
		// Juga bersihkan path lama jika user masih menyimpannya
		c.SetCookie("refresh_token", "", -1, "/api/v1/auth/refresh", "", true, true)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Logout berhasil",
	})
}

// RegisterCustomerInput is the DTO for register request
type RegisterCustomerInput struct {
	Username           string `json:"username" binding:"required"`
	Email              string `json:"email" binding:"required"`
	Password           string `json:"password" binding:"required"`
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

	if len(req.Password) < 8 {
		RespondWithError(c, http.StatusUnprocessableEntity, "Validasi gagal", "VAL_001", "Password minimal 8 karakter")
		return
	}

	if req.Password != req.KonfirmasiPassword {
		RespondWithError(c, http.StatusUnprocessableEntity, "Validasi gagal", "VAL_002", "Konfirmasi password tidak cocok")
		return
	}

	if _, err := mail.ParseAddress(req.Email); err != nil {
		RespondWithError(c, http.StatusUnprocessableEntity, "Validasi gagal", "VAL_003", "Format email tidak valid")
		return
	}

	// Cek duplikasi username
	var existingUser models.User
	if err := config.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		RespondWithError(c, http.StatusConflict, "Username sudah terdaftar", "CONF_001", "Username telah digunakan")
		return
	}

	// Cek duplikasi email
	if err := config.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		RespondWithError(c, http.StatusConflict, "Email sudah terdaftar", "CONF_002", "Email telah digunakan")
		return
	}

	// Cari Role "Customer"
	var role models.Role
	if err := config.DB.Where("nama_role = ?", "Customer").First(&role).Error; err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Role Customer tidak ditemukan di sistem", "SERVER_001", "Role tidak ditemukan")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Gagal memproses password", "SERVER_001", err.Error())
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
		RespondWithError(c, http.StatusInternalServerError, "Gagal membuat user baru", "SERVER_001", err.Error())
		return
	}

	newCustomer := models.Customer{
		NoTelp: req.NoTelp,
		UserId: newUser.IdUser,
	}

	if err := tx.Create(&newCustomer).Error; err != nil {
		tx.Rollback()
		RespondWithError(c, http.StatusInternalServerError, "Gagal membuat data customer", "SERVER_001", err.Error())
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
			RespondWithError(c, http.StatusInternalServerError, "Kesalahan sistem", "SERVER_001", "Format ID user tidak valid")
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.PasswordBaru), 12)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Gagal memproses password baru", "SERVER_001", err.Error())
		return
	}

	// Mulai transaksi
	tx := config.DB.Begin()

	// Update password
	user.Password = string(hashedPassword)
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		RespondWithError(c, http.StatusInternalServerError, "Gagal mengubah password", "SERVER_001", err.Error())
		return
	}

	// Revoke seluruh refresh token user (logout semua device)
	now := time.Now()
	if err := tx.Model(&models.RefreshToken{}).
		Where("id_user = ? AND revoked_at IS NULL", uid).
		Update("revoked_at", &now).Error; err != nil {
		tx.Rollback()
		RespondWithError(c, http.StatusInternalServerError, "Gagal me-revoke sesi lama", "SERVER_001", err.Error())
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Password berhasil diubah",
	})
}
