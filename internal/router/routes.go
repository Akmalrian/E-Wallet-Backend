package router

import (
	"ewallet-backend/internal/controller"
	"ewallet-backend/internal/middleware"
	"ewallet-backend/internal/repository"
	"ewallet-backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitRouter(app *gin.Engine, db *pgxpool.Pool) {
	// ── Repository ──
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)

	// ── Service ──
	authService := service.NewAuthService(userRepo, tokenRepo)
	userService := service.NewUserService(userRepo) // ← tambah

	// ── Controller ──
	authCtrl := controller.NewAuthController(authService)
	userCtrl := controller.NewUserController(userService) // ← tambah

	// ── Global Middleware ──
	app.Use(middleware.CORSMiddleware)

	// ── Routes ──
	api := app.Group("/api/v1")

	// Public
	auth := api.Group("/auth")
	{
		auth.POST("/register", authCtrl.Register)
		auth.POST("/login", authCtrl.Login)
	}

	// Protected
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(tokenRepo))
	{
		protected.DELETE("/auth/logout", authCtrl.Logout)
		protected.GET("/users/profile", userCtrl.GetProfile) // ← tambah
	}
}
