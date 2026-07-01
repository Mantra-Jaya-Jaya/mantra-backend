package auth

import (
	"backend-mantra/models"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// NormalizeRoleName membuat representasi role yang konsisten untuk JWT, context, dan handler.
func NormalizeRoleName(role string) string {
	return strings.ToLower(strings.TrimSpace(role))
}

// IsSecureCookie returns true only in production (HTTPS), false for local dev (HTTP)
func IsSecureCookie(c *gin.Context) bool {
	host := c.Request.Host
	if strings.HasPrefix(host, "localhost:") || strings.HasPrefix(host, "127.0.0.1:") || host == "localhost" || host == "127.0.0.1" {
		return false
	}
	return os.Getenv("MIDTRANS_ENVIRONMENT") == "production"
}

// GenerateJWT generates a real JWT token for a user
func GenerateJWT(userID uint, publicID string, role string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "rahasia_dapur_mantra" // Fallback
	}

	normalizedRole := NormalizeRoleName(role)

	claims := jwt.MapClaims{
		"user_id":   userID,
		"public_id": publicID,
		"role":      normalizedRole,
		"exp":       time.Now().Add(30 * time.Minute).Unix(), // 30 menit
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// RespondWithSuccess handles the different response formats for Flutter and Next.js
func RespondWithSuccess(c *gin.Context, clientType string, user models.User, roleName string, profileID uint, accessToken string, refreshToken string) {
	lowerRole := NormalizeRoleName(roleName)

	if clientType == "nextjs" {
		isSecure := IsSecureCookie(c)
		c.SetCookie("access_token", accessToken, 1800, "/", "", isSecure, true)
		c.SetCookie("refresh_token", refreshToken, 604800, "/", "", isSecure, true)

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Login berhasil",
			"data": gin.H{
				"user": gin.H{
					"id_user":      user.IdUser,
					"username":     user.Username,
					"nama_lengkap": user.NamaLengkap,
					"role":         lowerRole,
				},
			},
		})
	} else {
		// Default ke Flutter
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Login berhasil",
			"data": gin.H{
				"access_token":  accessToken,
				"refresh_token": refreshToken,
				"token_type":    "Bearer",
				"expires_in":    1800,
				"user": gin.H{
					"id_user":      user.IdUser,
					"username":     user.Username,
					"email":        user.Email,
					"nama_lengkap": user.NamaLengkap,
					"role":         lowerRole,
					"profile_id":   profileID,
				},
			},
		})
	}
}

// RespondWithError handles standard error responses according to the contract
func RespondWithError(c *gin.Context, status int, message string, code string, detail string) {
	c.JSON(status, gin.H{
		"status":  "error",
		"message": message,
		"error": gin.H{
			"code":   code,
			"detail": detail,
		},
	})
}

// IsOwnerOrAdmin checks if the current user is an admin or the owner of the data
func IsOwnerOrAdmin(c *gin.Context, dataOwnerPublicID string) bool {
	userRole := c.GetString("role")
	if strings.EqualFold(userRole, "admin") {
		return true
	}

	userPublicID := c.GetString("public_id")
	if userPublicID == "" || dataOwnerPublicID == "" {
		return false
	}

	return userPublicID == dataOwnerPublicID
}
