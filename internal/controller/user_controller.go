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

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService: userService}
}

// GetProfile GET /users/profile
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
