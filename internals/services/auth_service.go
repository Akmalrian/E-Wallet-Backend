package services

import (
	"errors"
	"ewallet-backend/internals/dto"
	"ewallet-backend/internals/models"
	"ewallet-backend/internals/repositories"
	"ewallet-backend/pkg"
)

type AuthService struct {
	UserRepo  *repositories.UserRepository
	Blacklist *repositories.TokenBlacklist
}

func NewAuthService(
	userRepo *repositories.UserRepository,
	blacklist *repositories.TokenBlacklist,
) *AuthService {
	return &AuthService{
		UserRepo:  userRepo,
		Blacklist: blacklist,
	}
}

// Register — proses registrasi
func (s *AuthService) Register(req dto.RegisterRequest) error {
	existing, err := s.UserRepo.FindByEmail(req.Email)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("email already registered")
	}

	hashedPassword, err := pkg.HashPassword(req.Password)
	if err != nil {
		return err
	}

	userID, err := s.UserRepo.Create(req.Email, hashedPassword)
	if err != nil {
		return err
	}

	return s.UserRepo.CreateWallet(userID)
}

// Login — proses login
func (s *AuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.UserRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	if !pkg.VerifyPassword(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	token, err := pkg.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User:  toUserResponse(user),
	}, nil
}

// Logout — tambahkan token ke blacklist
func (s *AuthService) Logout(token string) {
	s.Blacklist.Add(token)
}

// toUserResponse — konversi model ke DTO
func toUserResponse(user *models.User) dto.UserResponse {
	return dto.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		Fullname:    user.Fullname,
		PhotoPath:   user.PhotoPath,
		PhoneNumber: user.PhoneNumber,
		CreatedAt:   user.CreatedAt,
	}
}
