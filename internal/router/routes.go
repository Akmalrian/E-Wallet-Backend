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
	tokenRepo := repository.NewTokenRepository(rdb)
	walletRepo := repository.NewWalletRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	// ── Service ──
	authService := service.NewAuthService(userRepo, tokenRepo)
	userService := service.NewUserService(userRepo, rdb)
	walletService := service.NewWalletService(walletRepo, rdb)
	transactionService := service.NewTransactionService(
		transactionRepo,
		userRepo,
		walletRepo,
		rdb,
	)

	// ── Controller ──
	authCtrl := controller.NewAuthController(authService)
	userCtrl := controller.NewUserController(userService)
	walletCtrl := controller.NewWalletController(walletService)
	transactionCtrl := controller.NewTransactionController(transactionService)

	// ── Global Middleware ──
	app.Use(middleware.CORSMiddleware)

	// ── Setup Routes per grup ──
	SetupAuthRoutes(app, authCtrl, tokenRepo)
	SetupUserRoutes(app, userCtrl, walletCtrl, tokenRepo)
	SetupTransactionRoutes(app, transactionCtrl, tokenRepo)
}
