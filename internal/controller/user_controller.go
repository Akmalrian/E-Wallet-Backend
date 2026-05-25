package controller

import (
	"errors"
	"ewallet-backend/internal/dto"
	"ewallet-backend/internal/service"
	"ewallet-backend/pkg"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

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

// UpdateProfile godoc
//
//	@Summary		Update profile user
//	@Description	Update profile dengan atau tanpa upload foto
//	@Tags			Users
//	@Accept			mpfd
//	@Produce		json
//	@Param			fullname		formData	string	false	"Nama lengkap"
//	@Param			phone_number	formData	string	false	"Nomor telepon"
//	@Param			photo			formData	file	false	"Foto profile (opsional)"
//	@Success		200				{object}	pkg.BaseResponse
//	@Failure		400				{object}	pkg.ErrorResponse
//	@Failure		401				{object}	pkg.ErrorResponse
//	@Failure		500				{object}	pkg.ErrorResponse
//	@Security		BearerAuth
//	@Router			/users/profile [patch]
func (u *UserController) UpdateProfile(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	// Bind form data (bukan JSON)
	var body dto.UpdateProfileBody
	if err := ctx.ShouldBind(&body); err != nil {
		log.Println("Error:", err.Error())
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{
			Message: "Bad Request",
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Handle file upload
	photoPath := ""
	file, err := ctx.FormFile("photo")

	if err == nil {
		// ✅ Ada file yang diupload
		// Validasi tipe file
		allowedTypes := map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
			"image/jpg":  true,
		}

		// Buka file untuk cek content type
		openedFile, err := file.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{
				Message: "Internal Error",
				Success: false,
				Error:   "failed to open file",
			})
			return
		}
		defer openedFile.Close()

		// Baca 512 byte pertama untuk deteksi content type
		buffer := make([]byte, 512)
		openedFile.Read(buffer)
		contentType := http.DetectContentType(buffer)

		if !allowedTypes[contentType] {
			ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{
				Message: "Bad Request",
				Success: false,
				Error:   "only jpg and png files are allowed",
			})
			return
		}

		// Validasi ukuran file (max 2MB)
		if file.Size > 2*1024*1024 {
			ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{
				Message: "Bad Request",
				Success: false,
				Error:   "file size must be less than 2MB",
			})
			return
		}

		// Buat nama file unik pakai timestamp
		ext := filepath.Ext(file.Filename)
		fileName := fmt.Sprintf("photo_%d_%d%s",
			claims.Id,
			time.Now().Unix(),
			ext,
		)

		// Pastikan folder public/uploads ada
		uploadDir := "public/uploads"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{
				Message: "Internal Error",
				Success: false,
				Error:   "failed to create upload directory",
			})
			return
		}

		// Simpan file
		savePath := fmt.Sprintf("%s/%s", uploadDir, fileName)
		if err := ctx.SaveUploadedFile(file, savePath); err != nil {
			log.Println("Error saving file:", err.Error())
			ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{
				Message: "Internal Error",
				Success: false,
				Error:   "failed to save file",
			})
			return
		}

		// Path yang disimpan ke database
		photoPath = fmt.Sprintf("/uploads/%s", fileName)

	} else if err != http.ErrMissingFile {
		// ❌ Ada error selain "file tidak ada"
		log.Println("Error:", err.Error())
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{
			Message: "Bad Request",
			Success: false,
			Error:   "invalid file",
		})
		return
	}
	// ✅ Jika err == http.ErrMissingFile → tidak ada file → photoPath tetap ""
	// UpdateProfile di repository akan pakai foto lama

	if err := u.userService.UpdateProfile(ctx.Request.Context(), claims.Id, body, photoPath); err != nil {
		log.Println("Error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{
			Message: "Internal Error",
			Success: false,
			Error:   "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, pkg.BaseResponse{
		Message: "Profile updated successfully",
		Success: true,
	})
}

// UpdatePassword PATCH /users/password
//
//	@Summary		Update password user
//	@Description	Memverifikasi Password user lama sebelum mengganti password baru
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.UpdatePasswordBody	true	"Update Password Request"
//	@Success		200		{object}	pkg.BaseResponse
//	@Failure		400		{object}	pkg.ErrorResponse
//	@Failure		401		{object}	pkg.ErrorResponse
//	@Security		BearerAuth
//	@Router			/users/password [patch]
func (u *UserController) UpdatePassword(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	var body dto.UpdatePasswordBody
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println("Error:", err.Error())
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{
			Message: "Bad Request",
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := u.userService.UpdatePassword(ctx.Request.Context(), claims.Id, body); err != nil {
		log.Println("Error:", err.Error())

		statusCode := http.StatusInternalServerError
		if err.Error() == "old password is incorrect" ||
			err.Error() == "new password must be different from old password" {
			statusCode = http.StatusBadRequest
		}

		ctx.JSON(statusCode, pkg.ErrorResponse{
			Message: "Failed",
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, pkg.BaseResponse{
		Message: "Password updated successfully",
		Success: true,
	})
}

// UpdatePin PATCH /users/pin
//
//	@Summary		Update PIN user
//	@Description	Mengubah PIN lama dengan varifikasi PIN yang baru
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.UpdatePinBody	true	"Update PIN Request"
//	@Success		200		{object}	pkg.BaseResponse
//	@Failure		400		{object}	pkg.ErrorResponse
//	@Failure		401		{object}	pkg.ErrorResponse
//	@Security		BearerAuth
//	@Router			/users/pin [patch]
func (u *UserController) UpdatePin(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	var body dto.UpdatePinBody
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println("Error:", err.Error())
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{
			Message: "Bad Request",
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := u.userService.UpdatePin(ctx.Request.Context(), claims.Id, body); err != nil {
		log.Println("Error:", err.Error())

		statusCode := http.StatusInternalServerError
		if err.Error() == "old pin is incorrect" ||
			err.Error() == "new pin must be different from old pin" ||
			err.Error() == "pin has not been set" {
			statusCode = http.StatusBadRequest
		}

		ctx.JSON(statusCode, pkg.ErrorResponse{
			Message: "Failed",
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, pkg.BaseResponse{
		Message: "PIN updated successfully",
		Success: true,
	})
}
