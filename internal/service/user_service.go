package service

import (
	"context"
	"errors"
	"ewallet-backend/internal/cache"
	"ewallet-backend/internal/dto"
	"ewallet-backend/internal/repository"
	"ewallet-backend/pkg"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type UserService struct {
	userRepo *repository.UserRepository
	rdb      *redis.Client
}

func NewUserService(userRepo *repository.UserRepository, rdb *redis.Client) *UserService {
	return &UserService{
		userRepo: userRepo,
		rdb:      rdb,
	}
}

// GetProfile — ambil profile dengan cache
func (u *UserService) GetProfile(ctx context.Context, id int) (dto.GetProfileResponse, error) {
	cacheKey := fmt.Sprintf("user:profile:%d", id)

	var cached dto.GetProfileResponse
	if hit := cache.Get(ctx, u.rdb, cacheKey, &cached); hit {
		log.Println("profile: cache hit")
		return cached, nil
	}

	log.Println("profile: cache miss")
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return dto.GetProfileResponse{}, err
	}

	result := dto.GetProfileResponse{
		Id:          user.Id,
		Email:       user.Email,
		Fullname:    user.Fullname,
		PhotoPath:   user.PhotoPath,
		PhoneNumber: user.PhoneNumber,
		CreatedAt:   user.CreatedAt,
	}

	cache.Set(ctx, u.rdb, cacheKey, result, 5*time.Minute)

	return result, nil
}

// UpdateProfile — update profile dengan file upload
func (u *UserService) UpdateProfile(
	ctx context.Context,
	id int,
	body dto.UpdateProfileBody,
	photoPath string,
) error {
	_, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("user not found")
		}
		return err
	}

	if err := u.userRepo.UpdateProfile(ctx, id, body.Fullname, body.PhoneNumber, photoPath); err != nil {
		return err
	}

	// Hapus cache
	cacheKey := fmt.Sprintf("user:profile:%d", id)
	cache.Delete(ctx, u.rdb, cacheKey)

	return nil
}

// FindReceivers — cari penerima dengan search dan pagination
func (u *UserService) FindReceivers(
	ctx context.Context,
	currentUserID int,
	search string,
	page int,
	limit int,
) (dto.ReceiverListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 7
	}

	offset := (page - 1) * limit

	receivers, err := u.userRepo.FindReceivers(ctx, currentUserID, search, limit, offset)
	if err != nil {
		return dto.ReceiverListResponse{}, err
	}

	if receivers == nil {
		receivers = []dto.ReceiverResponse{}
	}

	total, err := u.userRepo.CountReceivers(ctx, currentUserID, search)
	if err != nil {
		return dto.ReceiverListResponse{}, err
	}

	totalPages := (total + limit - 1) / limit

	return dto.ReceiverListResponse{
		Receivers: receivers,
		Meta: dto.PaginationMeta{
			CurrentPage: page,
			TotalPages:  totalPages,
			TotalData:   total,
			Limit:       limit,
		},
	}, nil
}

// CheckPin — verifikasi PIN user
func (u *UserService) CheckPin(ctx context.Context, userID int, body dto.CheckPinBody) error {
	pin, err := u.userRepo.FindPinByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("user not found")
		}
		return err
	}

	if pin == "" {
		return errors.New("pin has not been set")
	}

	if !pkg.VerifyPassword(body.Pin, pin) {
		return errors.New("invalid pin")
	}
	return nil
}

// UpdatePassword — update password user
func (u *UserService) UpdatePassword(ctx context.Context, id int, body dto.UpdatePasswordBody) error {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("user not found")
		}
		return err
	}

	if !pkg.VerifyPassword(body.OldPassword, user.Password) {
		return errors.New("old password is incorrect")
	}

	if body.OldPassword == body.NewPassword {
		return errors.New("new password must be different from old password")
	}

	hashedPassword, err := pkg.HashPassword(body.NewPassword)
	if err != nil {
		return err
	}

	return u.userRepo.UpdatePassword(ctx, id, hashedPassword)
}

// UpdatePin — update PIN user
func (u *UserService) UpdatePin(ctx context.Context, id int, body dto.UpdatePinBody) error {

	currentPin, err := u.userRepo.FindPinByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("user not found")
		}
		return err
	}

	if currentPin != "" {
		if !pkg.VerifyPassword(body.OldPin, currentPin) {
			return errors.New("old pin is incorrect")
		}

		if !pkg.VerifyPassword(body.NewPin, currentPin) {
			return errors.New("new pin must be different from old pin")
		}

	}
	hashedPin, err := pkg.HashPassword(body.NewPin)
	if err != nil {
		return err
	}
	return u.userRepo.UpdatePin(ctx, id, hashedPin)
}

// ForgotPassword — step 1: verifikasi email
func (a *AuthService) ForgotPassword(ctx context.Context, body dto.ForgotPasswordBody) error {
	// Cek apakah email terdaftar
	_, err := a.userRepo.FindByEmail(ctx, body.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("email not found")
		}
		return err
	}

	// Email valid → bisa lanjut ke step 2
	return nil
}

// ResetPassword — step 2: ganti password baru
func (a *AuthService) ResetPassword(ctx context.Context, body dto.ResetPasswordBody) error {
	// Cek ulang apakah email masih valid
	_, err := a.userRepo.FindByEmail(ctx, body.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("email not found")
		}
		return err
	}

	// Hash password baru
	hashedPassword, err := pkg.HashPassword(body.NewPassword)
	if err != nil {
		return err
	}

	// Update password di database
	return a.userRepo.UpdatePasswordByEmail(ctx, body.Email, hashedPassword)
}
