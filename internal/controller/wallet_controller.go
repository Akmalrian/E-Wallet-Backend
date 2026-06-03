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

// GetDashboard godoc
//
//	@Summary		Get dashboard info
//	@Description	Mengambil informasi dashboard (balance, total income, total expense)
//	@Tags			Dashboard
//	@Produce		json
//	@Success		200	{object}	dto.SwaggerDashboardResponse	"Berhasil"
//	@Failure		401	{object}	pkg.ErrorResponse
//	@Failure		404	{object}	pkg.ErrorResponse
//	@Failure		500	{object}	pkg.ErrorResponse
//	@Security		BearerAuth
//	@Router			/users/dashboard [get]
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

// GetGraphData godoc
//
//	@Summary		Get dashboard graph
//	@Description	Mengambil data graph transaksi dengan filter type dan date range
//	@Tags			Dashboard
//	@Produce		json
//	@Param			type		query		string	false	"Filter: income / expense / both (default: both)"
//	@Param			start_date	query		string	false	"Tanggal mulai format: 2024-01-01 (default: 7 hari lalu)"
//	@Param			end_date	query		string	false	"Tanggal akhir format: 2024-01-31 (default: hari ini)"
//	@Success		200			{object}	dto.SwaggerGraphResponse
//	@Failure		401			{object}	pkg.ErrorResponse
//	@Failure		500			{object}	pkg.ErrorResponse
//	@Security		BearerAuth
//	@Router			/users/dashboard/graph [get]
func (w *WalletController) GetGraphData(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	// Ambil query params
	graphType := ctx.DefaultQuery("type", "both")
	startDate := ctx.DefaultQuery("start_date", "")
	endDate := ctx.DefaultQuery("end_date", "")

	result, err := w.walletService.GetGraphData(
		ctx.Request.Context(),
		claims.Id,
		graphType,
		startDate,
		endDate,
	)
	if err != nil {
		log.Println("Error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{
			Message: "Internal Error",
			Success: false,
			Error:   "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, pkg.Response[dto.GraphResponse]{
		Message: "OK",
		Success: true,
		Data:    result,
	})
}
