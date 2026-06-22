package middleware

import (
	"ewallet-backend/internal/repository"
	"ewallet-backend/pkg"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Hanya validasi JWT signature dan expiry
func validateToken(tokenString string) (*pkg.Claims, error) {
	claims, err := pkg.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}
	return &claims, nil
}

// Hanya cek whitelist di Redis
func checkWhitelist(ctx *gin.Context, tokenRepo *repository.TokenRepository, tokenString string) bool {
	isWhitelisted, err := tokenRepo.IsWhitelisted(ctx.Request.Context(), tokenString)
	if err != nil || !isWhitelisted {
		return false
	}
	return true
}

// AuthMiddleware
func AuthMiddleware(tokenRepo *repository.TokenRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Ambil token dari header
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

		// Validasi token dulu (signature + expiry)
		claims, err := validateToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse{
				Message: "Unauthorized",
				Success: false,
				Error:   "invalid or expired token",
			})
			ctx.Abort()
			return
		}

		// Baru cek whitelist
		if !checkWhitelist(ctx, tokenRepo, tokenString) {
			ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse{
				Message: "Unauthorized",
				Success: false,
				Error:   "token is not active. please login again",
			})
			ctx.Abort()
			return
		}

		// Simpan claims ke context untuk dipakai controller
		ctx.Set("claims", claims)
		ctx.Set("token", tokenString)
		ctx.Next()
	}
}
