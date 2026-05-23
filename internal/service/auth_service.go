package service

import (
	"context"
	"errors"
	"ewallet-backend/internal/dto"
	"ewallet-backend/internal/model"
	"ewallet-backend/internal/repository"
	"ewallet-backend/pkg"

	"github.com/jackc/pgx/v5"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	tokenRepo *repository.TokenRepository
}

func NewAuthService(
	userRepo *repository.UserRepository,
	tokenRepo *repository.TokenRepository,
) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

// Register — proses registrasi
func (a *AuthService) Register(ctx context.Context, body dto.RegisterBody) error {
	// Cek email sudah terdaftar
	_, err := a.userRepo.FindByEmail(ctx, body.Email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	if err == nil {
		return errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := pkg.HashPassword(body.Password)
	if err != nil {
		return err
	}

	// Buat user
	userID, err := a.userRepo.Create(ctx, body.Email, hashedPassword)
	if err != nil {
		return err
	}

	// Buat wallet
	return a.userRepo.CreateWallet(ctx, userID)
}

// Login — proses login
func (a *AuthService) Login(ctx context.Context, body dto.LoginBody) (dto.LoginResponse, error) {
	// Cari user
	user, err := a.userRepo.FindByEmail(ctx, body.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.LoginResponse{}, errors.New("invalid email or password")
		}
		return dto.LoginResponse{}, err
	}

	// Verifikasi password
	if !pkg.VerifyPassword(body.Password, user.Password) {
		return dto.LoginResponse{}, errors.New("invalid email or password")
	}

	// Generate token
	token, expiredAt, err := pkg.GenerateToken(user.Id, user.Email)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	// Simpan token ke database
	if err := a.tokenRepo.Add(ctx, user.Id, token, expiredAt); err != nil {
		return dto.LoginResponse{}, err
	}

	return dto.LoginResponse{
		Token: token,
		User:  toUserResponse(user),
	}, nil
}

// Logout — hapus token dari database
func (a *AuthService) Logout(ctx context.Context, token string) error {
	return a.tokenRepo.Remove(ctx, token)
}

// toUserResponse — konversi model ke dto
func toUserResponse(user model.User) dto.UserResponse {
	return dto.UserResponse{
		Id:          user.Id,
		Email:       user.Email,
		Fullname:    user.Fullname,
		PhotoPath:   user.PhotoPath,
		PhoneNumber: user.PhoneNumber,
		CreatedAt:   user.CreatedAt,
	}
}
