package controllers

import (
	"ewallet-backend/internals/dto"
	"ewallet-backend/internals/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{UserService: userService}
}

// GetProfile GET /users/profile
func (ctrl *UserController) GetProfile(c *gin.Context) {
	// Ambil user_id dari context (diset oleh AuthMiddleware)
	userID := c.GetInt("user_id")

	profile, err := ctrl.UserService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Status:  "success",
		Message: "Profile retrieved successfully",
		Data:    profile,
	})
}
