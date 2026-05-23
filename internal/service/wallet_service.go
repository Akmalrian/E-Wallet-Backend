package service

import (
	"context"
	"ewallet-backend/internal/dto"
	"ewallet-backend/internal/repository"
)

type WalletService struct {
	walletRepo *repository.WalletRepository
}

func NewWalletService(walletRepo *repository.WalletRepository) *WalletService {
	return &WalletService{walletRepo: walletRepo}
}

func (w *WalletService) GetDashboardInfo(ctx context.Context, userID int) (dto.DashboardResponse, error) {
	return w.walletRepo.GetDashboardInfo(ctx, userID)
}
