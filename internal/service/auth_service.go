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
// Login — tambah HasPin di response
func (a *AuthService) Login(ctx context.Context, body dto.LoginBody) (dto.LoginResponse, error) {
	user, err := a.userRepo.FindByEmail(ctx, body.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.LoginResponse{}, errors.New("invalid email or password")
		}
		return dto.LoginResponse{}, err
	}

	if !pkg.VerifyPassword(body.Password, user.Password) {
		return dto.LoginResponse{}, errors.New("invalid email or password")
	}

	// Generate token
	token, expiredAt, err := pkg.GenerateToken(user.Id, user.Email)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	// Simpan token ke Redis
	if err := a.tokenRepo.Add(ctx, user.Id, token, expiredAt); err != nil {
		return dto.LoginResponse{}, err
	}

	// ✅ Cek apakah user sudah punya PIN
	hasPin := user.Pin != nil && *user.Pin != ""

	return dto.LoginResponse{
		Token:  token,
		HasPin: hasPin, // ← false jika PIN kosong/null
		User:   toUserResponse(user),
	}, nil
}

// EnterPin — set PIN pertama kali setelah login
func (a *AuthService) EnterPin(ctx context.Context, userID int, body dto.EnterPinBody) error {
	// Ambil PIN saat ini
	pin, err := a.userRepo.FindPinByID(ctx, userID)
	if err != nil {
		return err
	}

	// Jika sudah punya PIN → tidak boleh pakai endpoint ini
	// Gunakan PATCH /users/pin untuk ganti PIN
	if pin != "" {
		return errors.New("pin already set. use change pin to update")
	}

	// Hash PIN
	hashedPin, err := pkg.HashPassword(body.Pin)
	if err != nil {
		return err
	}

	// Simpan PIN ke database
	return a.userRepo.UpdatePin(ctx, userID, hashedPin)
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
