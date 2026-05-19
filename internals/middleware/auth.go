package middleware

import (
	"ewallet-backend/internals/repositories"
	"ewallet-backend/pkg"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware — validasi JWT token di setiap request protected
func AuthMiddleware(blacklist *repositories.TokenBlacklist) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Format harus "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Invalid format. Use: Bearer <token>",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Cek apakah token sudah di blacklist (sudah logout)
		if blacklist.IsBlacklisted(tokenString) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Token has been logged out",
			})
			c.Abort()
			return
		}

		// Parse & validasi token
		claims, err := pkg.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Simpan data ke context untuk dipakai controller
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("token", tokenString) // simpan token untuk keperluan logout

		c.Next()
	}
}
