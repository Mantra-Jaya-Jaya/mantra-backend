package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// OwnershipMiddleware checks if the requested resource's public ID in the URL matches the user's public ID.
// Cocok digunakan untuk rute seperti GET /users/:public_id atau PUT /customers/:public_id
// Note: Role "admin" akan selalu diizinkan melewati pengecekan ini.
func OwnershipMiddleware(paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("role")

		// Admin selalu memiliki akses (Bypass)
		if strings.EqualFold(userRole, "admin") {
			c.Next()
			return
		}

		userPublicID := c.GetString("public_id")
		resourcePublicID := c.Param(paramName)

		// Jika param kosong, atau userPublicId kosong, atau tidak cocok
		if userPublicID == "" || resourcePublicID == "" || userPublicID != resourcePublicID {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": "Akses ditolak",
				"error": gin.H{
					"code":   "AUTH_004",
					"detail": "Anda tidak memiliki kepemilikan atas data ini",
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
