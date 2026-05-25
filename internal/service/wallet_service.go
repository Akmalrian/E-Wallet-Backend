package service

import (
	"context"
	"ewallet-backend/internal/cache"
	"ewallet-backend/internal/dto"
	"ewallet-backend/internal/repository"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type WalletService struct {
	walletRepo *repository.WalletRepository
	rdb        *redis.Client
}

func NewWalletService(walletRepo *repository.WalletRepository, rdb *redis.Client) *WalletService {
	return &WalletService{walletRepo: walletRepo, rdb: rdb}
}

// GetDashboardInfo — ambil dashboard dengan cache
func (w *WalletService) GetDashboardInfo(ctx context.Context, userID int) (dto.DashboardResponse, error) {
	cacheKey := fmt.Sprintf("user:dashboard:%d", userID)

	// Cek cache
	var cached dto.DashboardResponse
	if hit := cache.Get(ctx, w.rdb, cacheKey, &cached); hit {
		log.Println("dashboard: cache hit")
		return cached, nil
	}

	// Cache miss → query database
	log.Println("dashboard: cache miss")
	result, err := w.walletRepo.GetDashboardInfo(ctx, userID)
	if err != nil {
		return dto.DashboardResponse{}, err
	}

	// Simpan ke cache selama 1 menit
	// Dashboard lebih sering berubah jadi TTL lebih pendek
	cache.Set(ctx, w.rdb, cacheKey, result, 1*time.Minute)

	return result, nil
}

// InvalidateDashboardCache — hapus cache dashboard
// dipanggil setelah topup atau transfer
func (w *WalletService) InvalidateDashboardCache(ctx context.Context, userID int) {
	cacheKey := fmt.Sprintf("user:dashboard:%d", userID)
	cache.Delete(ctx, w.rdb, cacheKey)
}
