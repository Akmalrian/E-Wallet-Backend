package controllers

import (
	"ewallet-backend/internals/dto"
	"ewallet-backend/internals/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	AuthService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{AuthService: authService}
}

// Register POST /auth/register
func (ctrl *AuthController) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	if err := ctrl.AuthService.Register(req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse{
		Status:  "success",
		Message: "Registration successful",
	})
}

// Login POST /auth/login
func (ctrl *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	response, err := ctrl.AuthService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Status:  "success",
		Message: "Login successful",
		Data:    response,
	})
}

// Logout POST /auth/logout
func (ctrl *AuthController) Logout(c *gin.Context) {
	// Ambil token dari context (diset oleh AuthMiddleware)
	token, _ := c.Get("token")

	// Masukkan token ke blacklist
	ctrl.AuthService.Logout(token.(string))

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Status:  "success",
		Message: "Logout successful",
	})
}
