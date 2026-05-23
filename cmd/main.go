package main

import (
	"ewallet-backend/internal/config"
	"ewallet-backend/internal/router"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading env. \ncause: %s", err.Error())
	}

	// Inisialisasi gin
	app := gin.Default()

	// Koneksi database
	db, err := config.ConnectPsql()
	if err != nil {
		log.Fatalf("DB connection error. \ncause: %s", err.Error())
	}
	defer db.Close()
	log.Println("DB Connected")

	// Setup router
	router.InitRouter(app, db)

	// Jalankan server
	app.Run(fmt.Sprintf("%s:%s",
		os.Getenv("APP_HOST"),
		os.Getenv("APP_PORT"),
	))
}
