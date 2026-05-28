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

	var cached dto.DashboardResponse
	if hit := cache.Get(ctx, w.rdb, cacheKey, &cached); hit {
		log.Println("dashboard: cache hit")
		return cached, nil
	}

	log.Println("dashboard: cache miss")
	result, err := w.walletRepo.GetDashboardInfo(ctx, userID)
	if err != nil {
		return dto.DashboardResponse{}, err
	}

	cache.Set(ctx, w.rdb, cacheKey, result, 1*time.Minute)
	return result, nil
}

// GetGraphData — ambil data graph dengan validasi filter
func (w *WalletService) GetGraphData(
	ctx context.Context,
	userID int,
	graphType string,
	startDateStr string,
	endDateStr string,
) (dto.GraphResponse, error) {

	// Validasi graphType
	// Default: "both" jika kosong atau tidak valid
	validTypes := map[string]bool{
		"income":  true,
		"expense": true,
		"both":    true,
	}
	if !validTypes[graphType] {
		graphType = "both"
	}

	// Parse startDate
	// Default: 7 hari yang lalu jika tidak diisi
	var startDate time.Time
	var err error
	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			startDate = time.Now().AddDate(0, 0, -7)
		}
	} else {
		startDate = time.Now().AddDate(0, 0, -7)
	}

	// Parse endDate
	// Default: hari ini jika tidak diisi
	var endDate time.Time
	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			endDate = time.Now()
		}
		// Set ke akhir hari (23:59:59)
		endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	} else {
		endDate = time.Now()
	}

	// Validasi startDate tidak boleh lebih dari endDate
	if startDate.After(endDate) {
		startDate = endDate.AddDate(0, 0, -7)
	}

	points, err := w.walletRepo.GetGraphData(ctx, userID, graphType, startDate, endDate)
	if err != nil {
		return dto.GraphResponse{}, err
	}

	// Jika tidak ada data, kembalikan slice kosong
	if points == nil {
		points = []dto.GraphPoint{}
	}

	return dto.GraphResponse{
		Points: points,
		Filter: dto.GraphFilter{
			Type:      graphType,
			StartDate: startDate.Format("2006-01-02"),
			EndDate:   endDate.Format("2006-01-02"),
		},
	}, nil
}

// InvalidateDashboardCache — hapus cache dashboard
func (w *WalletService) InvalidateDashboardCache(ctx context.Context, userID int) {
	cacheKey := fmt.Sprintf("user:dashboard:%d", userID)
	cache.Delete(ctx, w.rdb, cacheKey)
}
