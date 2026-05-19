package main

import (
	"ewallet-backend/config"
	"ewallet-backend/internals/controllers"
	"ewallet-backend/internals/repositories"
	"ewallet-backend/internals/routes"
	"ewallet-backend/internals/services"

	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load .env:", err)
	}

	config.ConnectDatabase()
	defer config.DB.Close()

	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// ── Blacklist (untuk logout) ──
	blacklist := repositories.NewTokenBlacklist()

	// ── Repository ──
	userRepo := repositories.NewUserRepository(config.DB)

	// ── Service ──
	authService := services.NewAuthService(userRepo, blacklist)
	userService := services.NewUserService(userRepo)

	// ── Controller ──
	authCtrl := controllers.NewAuthController(authService)
	userCtrl := controllers.NewUserController(userService)

	// ── Routes ──
	routes.SetupRoutes(r, blacklist, authCtrl, userCtrl)

	port := os.Getenv("APP_PORT")
	log.Printf("Server running on port %s", port)
	r.Run(":" + port)
}
