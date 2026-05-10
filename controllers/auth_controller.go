package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Login berhasil",
		"data": gin.H{
			"access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.mock_token",
			"token_type":   "Bearer",
			"expires_in":   900,
			"user": gin.H{
				"id_user":      1,
				"username":     "john_doe",
				"email":        "john@email.com",
				"nama_lengkap": "John Doe",
				"role":         "customer",
				"profile_id":   10,
			},
		},
	})
}

// RefreshToken handles token refresh
func RefreshToken(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Token berhasil diperbarui",
		"data": gin.H{
			"access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.mock_new_token",
			"expires_in":   900,
		},
	})
}

// Logout handles user logout
func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Logout berhasil",
	})
}

// RegisterCustomer handles customer registration
func RegisterCustomer(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Registrasi berhasil",
		"data": gin.H{
			"id_user":      1,
			"username":     "john_doe",
			"email":        "john@email.com",
			"nama_lengkap": "John Doe",
			"no_telp":      "081234567890",
			"role":         "customer",
		},
	})
}

// ChangePassword handles password change
func ChangePassword(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Password berhasil diubah",
	})
}
