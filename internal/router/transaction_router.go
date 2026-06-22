package router

import (
	"ewallet-backend/internal/controller"
	"ewallet-backend/internal/middleware"
	"ewallet-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

func SetupTransactionRoutes(app *gin.Engine, transactionCtrl *controller.TransactionController, tokenRepo *repository.TokenRepository) {

	transactions := app.Group("/transactions")
	transactions.Use(middleware.AuthMiddleware(tokenRepo))

	transactions.POST("/topup", transactionCtrl.CreateTopup)
	transactions.POST("/transfer", transactionCtrl.CreateTransfer)
	transactions.GET("/history", transactionCtrl.GetHistory)
}
