package services

import (
	"errors"
	"ewallet-backend/internals/dto"
	"ewallet-backend/internals/repositories"
)

type UserService struct {
	UserRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{UserRepo: userRepo}
}

// GetProfile — ambil profile user berdasarkan id dari token
func (s *UserService) GetProfile(userID int) (*dto.GetProfileResponse, error) {
	user, err := s.UserRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return &dto.GetProfileResponse{
		ID:          user.ID,
		Email:       user.Email,
		Fullname:    user.Fullname,
		PhotoPath:   user.PhotoPath,
		PhoneNumber: user.PhoneNumber,
		CreatedAt:   user.CreatedAt,
	}, nil
}
