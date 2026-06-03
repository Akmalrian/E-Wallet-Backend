package service

import (
	"context"
	"errors"
	"ewallet-backend/internal/dto"
	"ewallet-backend/internal/repository"
	"ewallet-backend/pkg"
)

type TransactionService struct {
	transactionRepo *repository.TransactionRepository
	userRepo        *repository.UserRepository
	walletService   *WalletService // ← tambah untuk invalidate cache
}

func NewTransactionService(
	transactionRepo *repository.TransactionRepository,
	userRepo *repository.UserRepository,
	walletService *WalletService, // ← tambah
) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
		walletService:   walletService,
	}
}

// CreateTopup — setelah berhasil, invalidate cache dashboard
func (t *TransactionService) CreateTopup(
	ctx context.Context,
	userID int,
	body dto.TopupBody,
) (dto.TopupResponse, error) {
	walletID, _, err := t.transactionRepo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return dto.TopupResponse{}, errors.New("wallet not found")
	}

	result, err := t.transactionRepo.CreateTopup(ctx, userID, walletID, body)
	if err != nil {
		return dto.TopupResponse{}, err
	}

	// Invalidate cache dashboard setelah topup
	t.walletService.InvalidateDashboardCache(ctx, userID)

	return result, nil
}

// CreateTransfer — setelah berhasil, invalidate cache dashboard kedua user
func (t *TransactionService) CreateTransfer(
	ctx context.Context,
	senderUserID int,
	body dto.TransferBody,
) (dto.TransferResponse, error) {
	pin, err := t.userRepo.FindPinByID(ctx, senderUserID)
	if err != nil {
		return dto.TransferResponse{}, err
	}
	if pin == "" {
		return dto.TransferResponse{}, errors.New("pin has not been set")
	}
	if !pkg.VerifyPassword(body.Pin, pin) {
		return dto.TransferResponse{}, errors.New("invalid pin")
	}

	senderWalletID, _, err := t.transactionRepo.GetWalletByUserID(ctx, senderUserID)
	if err != nil {
		return dto.TransferResponse{}, errors.New("sender wallet not found")
	}

	var receiverBalance float64
	err = t.transactionRepo.CheckWalletExists(ctx, body.ReceiverWalletId, &receiverBalance)
	if err != nil {
		return dto.TransferResponse{}, errors.New("receiver wallet not found")
	}

	result, err := t.transactionRepo.CreateTransfer(ctx, senderUserID, senderWalletID, body)
	if err != nil {
		return dto.TransferResponse{}, err
	}

	// Invalidate cache dashboard pengirim
	t.walletService.InvalidateDashboardCache(ctx, senderUserID)

	return result, nil
}

// GetHistory — ambil history dengan pagination
func (t *TransactionService) GetHistory(
	ctx context.Context,
	userID int,
	transType string,
	page int,
	limit int,
) (dto.HistoryListResponse, error) {

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	transactions, err := t.transactionRepo.GetHistory(ctx, userID, transType, limit, offset)
	if err != nil {
		return dto.HistoryListResponse{}, err
	}

	if transactions == nil {
		transactions = []dto.HistoryResponse{}
	}

	total, err := t.transactionRepo.CountHistory(ctx, userID, transType)
	if err != nil {
		return dto.HistoryListResponse{}, err
	}

	totalPages := (total + limit - 1) / limit

	return dto.HistoryListResponse{
		Transactions: transactions,
		Meta: dto.PaginationMeta{
			CurrentPage: page,
			TotalPages:  totalPages,
			TotalData:   total,
			Limit:       limit,
		},
	}, nil
}
