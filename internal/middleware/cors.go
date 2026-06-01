package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware — handle CORS & preflight request
func CORSMiddleware(ctx *gin.Context) {
	allowedOrigins := []string{
		"http://localhost:5173", // ← Vite default port
		"http://127.0.0.1:5173",
	}

	currentOrigin := ctx.GetHeader("Origin")
	if slices.Contains(allowedOrigins, currentOrigin) {
		ctx.Header("Access-Control-Allow-Origin", currentOrigin)
	}

	allowedHeaders := []string{"Content-Type", "Authorization"}
	ctx.Header("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ", "))

	allowedMethods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodOptions,
	}
	ctx.Header("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))

	// Handle preflight request
	if ctx.Request.Method == http.MethodOptions {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}

	ctx.Next()
}
