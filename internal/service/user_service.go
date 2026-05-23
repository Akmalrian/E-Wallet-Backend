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

// FindReceivers — cari penerima dengan search dan pagination
func (u *UserService) FindReceivers(
	ctx context.Context,
	currentUserID int,
	search string,
	page int,
	limit int,
) (dto.ReceiverListResponse, error) {

	// Validasi page dan limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 7
	}

	offset := (page - 1) * limit

	// Ambil data receivers
	receivers, err := u.userRepo.FindReceivers(ctx, currentUserID, search, limit, offset)
	if err != nil {
		return dto.ReceiverListResponse{}, err
	}

	// Jika tidak ada hasil, kembalikan slice kosong
	if receivers == nil {
		receivers = []dto.ReceiverResponse{}
	}

	// Hitung total data untuk pagination
	total, err := u.userRepo.CountReceivers(ctx, currentUserID, search)
	if err != nil {
		return dto.ReceiverListResponse{}, err
	}

	// Hitung total halaman
	// contoh: total=15, limit=7 → totalPages=3
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
