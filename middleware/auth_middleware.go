package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the custom claims in our JWT
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	PublicID string `json:"public_id"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GetTokenFromRequest extracts token from Authorization header or cookie
func GetTokenFromRequest(c *gin.Context) (string, error) {
	// Try Authorization header first (for Flutter)
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer "), nil
	}

	// Fallback to cookie (for Next.js)
	cookie, err := c.Cookie("access_token")
	if err == nil && cookie != "" {
		return cookie, nil
	}

	return "", errors.New("token tidak ditemukan di header maupun cookie")
}

// AuthMiddleware validates the JWT token and sets user context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := GetTokenFromRequest(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Token tidak valid atau sudah expired",
				"error": gin.H{
					"code":   "AUTH_001",
					"detail": err.Error(),
				},
			})
			c.Abort()
			return
		}

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "rahasia_dapur_mantra"
		}

		claims := &JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Token tidak valid atau sudah expired",
				"error": gin.H{
					"code":   "AUTH_001",
					"detail": err.Error(),
				},
			})
			c.Abort()
			return
		}

		// Save to context for next handlers
		c.Set("user_id", claims.UserID)
		c.Set("public_id", claims.PublicID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
