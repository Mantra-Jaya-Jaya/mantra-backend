package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// RoleMiddleware checks if the user has the required role
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("role")
		
		if userRole == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": "Akses ditolak",
				"error": gin.H{
					"code":   "AUTH_002",
					"detail": "Role tidak ditemukan pada token",
				},
			})
			c.Abort()
			return
		}

		isAllowed := false
		for _, role := range allowedRoles {
			if strings.EqualFold(userRole, role) {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": "Akses ditolak",
				"error": gin.H{
					"code":   "AUTH_003",
					"detail": "Anda tidak memiliki izin untuk mengakses resource ini",
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
