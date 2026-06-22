package router

import (
	"ewallet-backend/internal/controller"
	"ewallet-backend/internal/middleware"
	"ewallet-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(app *gin.Engine, authCtrl *controller.AuthController, tokenRepo *repository.TokenRepository) {
	// Public routes
	auth := app.Group("/auth")
	auth.POST("/register", authCtrl.Register)
	auth.POST("/login", authCtrl.Login)
	auth.POST("/forgot-password", authCtrl.ForgotPassword)
	auth.POST("/reset-password", authCtrl.ResetPassword)

	// Protected routes
	authProtected := app.Group("/auth")
	authProtected.Use(middleware.AuthMiddleware(tokenRepo))
	authProtected.POST("/enter-pin", authCtrl.EnterPin)
	authProtected.DELETE("/logout", authCtrl.Logout)
}
