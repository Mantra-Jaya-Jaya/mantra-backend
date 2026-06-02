package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"backend-mantra/controllers/auth"

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

		// Sliding Expiration: Renew token if less than 15 minutes left
		if claims.ExpiresAt != nil {
			timeRemaining := time.Until(claims.ExpiresAt.Time)
			if timeRemaining > 0 && timeRemaining < 15*time.Minute {
				newToken, errGenerate := auth.GenerateJWT(claims.UserID, claims.PublicID, claims.Role)
				if errGenerate == nil && newToken != "" {
					c.SetCookie("access_token", newToken, 1800, "/", "", true, true)
					c.Header("X-New-Access-Token", newToken)
					c.Header("Access-Control-Expose-Headers", "X-New-Access-Token")
				}
			}
		}

		// Save to context for next handlers
		c.Set("user_id", claims.UserID)
		c.Set("public_id", claims.PublicID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
