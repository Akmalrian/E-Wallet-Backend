package middleware

import (
	"context"
	"ewallet-backend/internal/repository"
	"ewallet-backend/pkg"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware — validasi JWT token
func AuthMiddleware(tokenRepo *repository.TokenRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse{
				Message: "Unauthorized",
				Success: false,
				Error:   "authorization header is required",
			})
			ctx.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse{
				Message: "Unauthorized",
				Success: false,
				Error:   "invalid format. use: Bearer <token>",
			})
			ctx.Abort()
			return
		}

		tokenString := parts[1]

		// Cek whitelist di database
		isValid, err := tokenRepo.IsWhitelisted(context.Background(), tokenString)
		if err != nil || !isValid {
			ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse{
				Message: "Unauthorized",
				Success: false,
				Error:   "token is not active. please login first",
			})
			ctx.Abort()
			return
		}

		// Parse token
		claims, err := pkg.ParseToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse{
				Message: "Unauthorized",
				Success: false,
				Error:   "invalid or expired token",
			})
			ctx.Abort()
			return
		}

		// Simpan claims ke context — sama persis dengan referensi
		ctx.Set("claims", claims)
		ctx.Set("token", tokenString)

		ctx.Next()
	}
}
