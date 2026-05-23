package service

import (
	"context"
	"ewallet-backend/internal/dto"
	"ewallet-backend/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// GetProfile — ambil profile user berdasarkan id dari token
func (u *UserService) GetProfile(ctx context.Context, id int) (dto.GetProfileResponse, error) {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return dto.GetProfileResponse{}, err
	}

	return dto.GetProfileResponse{
		Id:          user.Id,
		Email:       user.Email,
		Fullname:    user.Fullname,
		PhotoPath:   user.PhotoPath,
		PhoneNumber: user.PhoneNumber,
		CreatedAt:   user.CreatedAt,
	}, nil
}
