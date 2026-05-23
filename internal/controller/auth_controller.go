package controller

import (
	"ewallet-backend/internal/dto"
	"ewallet-backend/internal/service"
	"ewallet-backend/pkg"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Register POST /auth/register
func (a *AuthController) Register(ctx *gin.Context) {
	var body dto.RegisterBody
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println("Error:", err.Error())
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{
			Message: "Bad Request",
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := a.authService.Register(ctx.Request.Context(), body); err != nil {
		log.Println("Error:", err.Error())
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{
			Message: "Bad Request",
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, pkg.BaseResponse{
		Message: "Registration successful",
		Success: true,
	})
}

// Login POST /auth/login
func (a *AuthController) Login(ctx *gin.Context) {
	var body dto.LoginBody
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println("Error:", err.Error())
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{
			Message: "Bad Request",
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	result, err := a.authService.Login(ctx.Request.Context(), body)
	if err != nil {
		log.Println("Error:", err.Error())
		ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse{
			Message: "Unauthorized",
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, pkg.Response[dto.LoginResponse]{
		Message: "Login successful",
		Success: true,
		Data:    result,
	})
}

// Logout DELETE /auth/logout
func (a *AuthController) Logout(ctx *gin.Context) {
	// Ambil token dari context (diset oleh AuthMiddleware)
	token, _ := ctx.Get("token")

	if err := a.authService.Logout(ctx.Request.Context(), token.(string)); err != nil {
		log.Println("Error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{
			Message: "Internal Error",
			Success: false,
			Error:   "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, pkg.BaseResponse{
		Message: "Logout successful",
		Success: true,
	})
}
