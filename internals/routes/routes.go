package routes

import (
	"ewallet-backend/internals/controllers"
	"ewallet-backend/internals/middleware"
	"ewallet-backend/internals/repositories"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	blacklist *repositories.TokenBlacklist,
	authCtrl *controllers.AuthController,
	userCtrl *controllers.UserController,
) {
	r.Use(middleware.CORSMiddleware())

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "E-Wallet API is running",
			"version": "1.0.0",
		})
	})

	// ── Public (tidak perlu token)
	auth := r.Group("/auth")

	auth.POST("/register", authCtrl.Register)
	auth.POST("/login", authCtrl.Login)

	// ── Protected (butuh token)
	protected := r.Group("/users")
	protected.Use(middleware.AuthMiddleware(blacklist))

	protected.POST("/logout", authCtrl.Logout)

	protected.GET("/profile", userCtrl.GetProfile)

}
