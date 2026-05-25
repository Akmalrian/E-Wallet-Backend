package router

import (
	"ewallet-backend/internal/controller"
	"ewallet-backend/internal/middleware"
	"ewallet-backend/internal/repository"
	"ewallet-backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitRouter(app *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	app.Static("/uploads", "./public/uploads")
	// ── Repository ──
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	walletRepo := repository.NewWalletRepository(db)

	// ── Service ──
	authService := service.NewAuthService(userRepo, tokenRepo)
	userService := service.NewUserService(userRepo, rdb)
	walletService := service.NewWalletService(walletRepo, rdb)

	// ── Controller ──
	authCtrl := controller.NewAuthController(authService)
	userCtrl := controller.NewUserController(userService)
	walletCtrl := controller.NewWalletController(walletService)

	// ── Global Middleware ──
	app.Use(middleware.CORSMiddleware)

	// ── Routes ──

	auth := app.Group("/auth")

	auth.POST("/register", authCtrl.Register)
	auth.POST("/login", authCtrl.Login)

	protected := app.Group("/")
	protected.Use(middleware.AuthMiddleware(tokenRepo))
	protected.DELETE("/auth/logout", authCtrl.Logout)
	protected.GET("/users/profile", userCtrl.GetProfile)
	protected.GET("/users/dashboard", walletCtrl.GetDashboardInfo)
	protected.GET("/users/receiver", userCtrl.FindReceivers)
	protected.PATCH("/users/profile", userCtrl.UpdateProfile)
	protected.PATCH("/users/password", userCtrl.UpdatePassword)
	protected.PATCH("/users/pin", userCtrl.UpdatePin)
	protected.POST("/users/pin/check", userCtrl.CheckPin)
}
