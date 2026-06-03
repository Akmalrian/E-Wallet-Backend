package controller

import (
	"ewallet-backend/internal/dto"
	"ewallet-backend/internal/service"
	"ewallet-backend/pkg"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type TransactionController struct {
	transactionService *service.TransactionService
}

func NewTransactionController(transactionService *service.TransactionService) *TransactionController {
	return &TransactionController{transactionService: transactionService}
}

// CreateTopup godoc
//
//	@Summary		Top up saldo
//	@Description	Menambah saldo wallet user
//	@Tags			Transactions
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.TopupBody				true	"Topup Request"
//	@Success		201		{object}	dto.SwaggerTopupResponse
//	@Failure		400		{object}	pkg.ErrorResponse
//	@Failure		401		{object}	pkg.ErrorResponse
//	@Failure		500		{object}	pkg.ErrorResponse
//	@Security		BearerAuth
//	@Router			/transactions/topup [post]
func (t *TransactionController) CreateTopup(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	var body dto.TopupBody
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println("Error:", err.Error())
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{
			Message: "Bad Request",
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	result, err := t.transactionService.CreateTopup(ctx.Request.Context(), claims.Id, body)
	if err != nil {
		log.Println("Error:", err.Error())
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{
			Message: "Bad Request",
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, pkg.Response[dto.TopupResponse]{
		Message: "Top up successful",
		Success: true,
		Data:    result,
	})
}

// CreateTransfer godoc
//
//	@Summary		Transfer saldo
//	@Description	Transfer saldo ke user lain dengan verifikasi PIN
//	@Tags			Transactions
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.TransferBody			true	"Transfer Request"
//	@Success		201		{object}	dto.SwaggerTransferResponse
//	@Failure		400		{object}	pkg.ErrorResponse
//	@Failure		401		{object}	pkg.ErrorResponse
//	@Failure		500		{object}	pkg.ErrorResponse
//	@Security		BearerAuth
//	@Router			/transactions/transfer [post]
func (t *TransactionController) CreateTransfer(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	var body dto.TransferBody
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println("Error:", err.Error())
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{
			Message: "Bad Request",
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	result, err := t.transactionService.CreateTransfer(ctx.Request.Context(), claims.Id, body)
	if err != nil {
		log.Println("Error:", err.Error())

		statusCode := http.StatusBadRequest
		if err.Error() == "invalid pin" {
			statusCode = http.StatusUnauthorized
		}

		ctx.JSON(statusCode, pkg.ErrorResponse{
			Message: "Failed",
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, pkg.Response[dto.TransferResponse]{
		Message: "Transfer successful",
		Success: true,
		Data:    result,
	})
}

// GetHistory godoc
//
//	@Summary		History transaksi
//	@Description	Mengambil riwayat transaksi user dengan filter dan pagination
//	@Tags			Transactions
//	@Produce		json
//	@Param			type	query		string	false	"Filter tipe: topup / transfer"
//	@Param			page	query		int		false	"Halaman (default: 1)"
//	@Param			limit	query		int		false	"Jumlah data (default: 10)"
//	@Success		200		{object}	dto.SwaggerHistoryResponse
//	@Failure		401		{object}	pkg.ErrorResponse
//	@Failure		500		{object}	pkg.ErrorResponse
//	@Security		BearerAuth
//	@Router			/transactions/history [get]
func (t *TransactionController) GetHistory(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	// Query params
	transType := ctx.DefaultQuery("type", "")

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	result, err := t.transactionService.GetHistory(
		ctx.Request.Context(),
		claims.Id,
		transType,
		page,
		limit,
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

	ctx.JSON(http.StatusOK, pkg.Response[dto.HistoryListResponse]{
		Message: "OK",
		Success: true,
		Data:    result,
	})
}
