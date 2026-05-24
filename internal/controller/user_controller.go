package controller

import (
	"errors"
	"ewallet-backend/internal/dto"
	"ewallet-backend/internal/service"
	"ewallet-backend/pkg"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jackc/pgx/v5"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService: userService}
}

// GetProfile godoc
//
//	@Summary		Get profile user
//	@Description	Mengambil data profile user yang sedang login
//	@Tags			Users
//	@Produce		json
//	@Success		200	{object}	dto.SwaggerProfileResponse	"Profile berhasil diambil"
//	@Failure		401	{object}	pkg.ErrorResponse
//	@Failure		404	{object}	pkg.ErrorResponse
//	@Failure		500	{object}	pkg.ErrorResponse
//	@Security		BearerAuth
//	@Router			/users/profile [get]
func (u *UserController) GetProfile(ctx *gin.Context) {
	// Ambil claims dari context yang diset AuthMiddleware
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	user, err := u.userService.GetProfile(ctx.Request.Context(), claims.Id)
	if err != nil {
		// Cek apakah user tidak ditemukan
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("Error:", err.Error())
			ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{
				Message: "Not Found",
				Success: false,
				Error:   "user not found",
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

	ctx.JSON(http.StatusOK, pkg.Response[dto.GetProfileResponse]{
		Message: "OK",
		Success: true,
		Data:    user,
	})
}

// FindReceivers GET /users/receivers?search=&page=&limit=
// FindReceivers godoc
//
//	@Summary		Cari penerima transfer
//	@Description	Mencari user lain sebagai penerima transfer dengan search dan pagination
//	@Tags			Users
//	@Produce		json
//	@Param			search	query		string	false	"Cari berdasarkan nama/email/phone"
//	@Param			page	query		int		false	"Halaman (default: 1)"
//	@Param			limit	query		int		false	"Jumlah data per halaman (default: 7)"
//	@Success		200		{object}	dto.SwaggerReceiversResponse	"Berhasil"
//	@Failure		401		{object}	pkg.ErrorResponse
//	@Failure		500		{object}	pkg.ErrorResponse
//	@Security		BearerAuth
//	@Router			/users/receivers [get]
func (u *UserController) FindReceivers(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	search := ctx.DefaultQuery("search", "")

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	result, err := u.userService.FindReceivers(
		ctx.Request.Context(),
		claims.Id,
		search,
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

	ctx.JSON(http.StatusOK, pkg.Response[dto.ReceiverListResponse]{
		Message: "OK",
		Success: true,
		Data:    result,
	})
}

// CheckPin godoc
//
//	@Summary		Check PIN user
//	@Description	Memverifikasi PIN user sebelum melakukan transaksi
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.CheckPinBody	true	"Check PIN Request"
//	@Success		200		{object}	pkg.BaseResponse
//	@Failure		400		{object}	pkg.ErrorResponse
//	@Failure		401		{object}	pkg.ErrorResponse
//	@Security		BearerAuth
//	@Router			/users/pin/check [post]
func (u *UserController) CheckPin(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	var body dto.CheckPinBody
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println("Error:", err.Error())
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{
			Message: "Bad Request",
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := u.userService.CheckPin(ctx.Request.Context(), claims.Id, body); err != nil {
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

	ctx.JSON(http.StatusOK, pkg.BaseResponse{
		Message: "Pin Is Correct",
		Success: true,
	})
}
