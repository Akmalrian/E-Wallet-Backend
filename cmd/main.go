package main

import (
	"ewallet-backend/docs"
	"ewallet-backend/internal/config"
	"ewallet-backend/internal/router"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title						Backend E-WALLET
// @version						1.0
// @description					Backend created by Koda using Gin

// @license.name				MIT

// @host						localhost:9000
// @BasePath					/

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Masukkan token dengan format: Bearer <token>

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading env. \ncause: %s", err.Error())
	}

	app := gin.Default()

	// Koneksi database
	db, err := config.ConnectPsql()
	if err != nil {
		log.Printf("DB connection error. \ncause: %s", err.Error())
	}
	defer db.Close()
	log.Println("DB Connected")

	// Koneksi Redis
	rdb, err := config.ConnectRedis()
	if err != nil {
		log.Printf("Redis connection error. \ncause: %s", err.Error())
	}
	defer rdb.Close()
	log.Println("Redis Connected")

	// Setup swagger
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s",
		os.Getenv("APP_HOST"),
		os.Getenv("APP_PORT"),
	)
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup router — pass rdb
	router.InitRouter(app, db, rdb)

	app.Run(fmt.Sprintf("%s:%s",
		os.Getenv("APP_HOST"),
		os.Getenv("APP_PORT"),
	))
}
