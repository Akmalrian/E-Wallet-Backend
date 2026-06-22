package router

import (
	"ewallet-backend/internal/controller"
	"ewallet-backend/internal/middleware"
	"ewallet-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(
	app *gin.Engine,
	userCtrl *controller.UserController,
	walletCtrl *controller.WalletController,
	tokenRepo *repository.TokenRepository,
) {
	users := app.Group("/users")
	users.Use(middleware.AuthMiddleware(tokenRepo))

	users.GET("/profile", userCtrl.GetProfile)
	users.PATCH("/profile", userCtrl.UpdateProfile)
	users.PATCH("/password", userCtrl.UpdatePassword)
	users.PATCH("/pin", userCtrl.UpdatePin)
	users.POST("/pin/check", userCtrl.CheckPin)
	users.GET("/receiver", userCtrl.FindReceivers)
	users.GET("/dashboard", walletCtrl.GetDashboardInfo)
	users.GET("/dashboard/graph", walletCtrl.GetGraphData)
}
