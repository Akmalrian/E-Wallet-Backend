package controller

import (
	"errors"
	"ewallet-backend/internal/dto"
	"ewallet-backend/internal/service"
	"ewallet-backend/pkg"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type WalletController struct {
	walletService *service.WalletService
}

func NewWalletController(walletService *service.WalletService) *WalletController {
	return &WalletController{walletService: walletService}
}

func (w *WalletController) GetDashboardInfo(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	dashboard, err := w.walletService.GetDashboardInfo(ctx.Request.Context(), claims.Id)
	if err != nil {
		// Cek apakah user tidak ditemukan
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("Error:", err.Error())
			ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{
				Message: "Not Found",
				Success: false,
				Error:   "wallet not found",
			})
			return
		}

		log.Println("Error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{
			Message: "Internal Error",
			Success: false,
			Error:   "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, pkg.Response[dto.DashboardResponse]{
		Message: "OK",
		Success: true,
		Data:    dashboard,
	})
}
