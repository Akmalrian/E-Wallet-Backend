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
}

func NewTransactionService(
	transactionRepo *repository.TransactionRepository,
	userRepo *repository.UserRepository,
) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
	}
}

// CreateTopup — proses topup
func (t *TransactionService) CreateTopup(
	ctx context.Context,
	userID int,
	body dto.TopupBody,
) (dto.TopupResponse, error) {

	// Ambil wallet user
	walletID, _, err := t.transactionRepo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return dto.TopupResponse{}, errors.New("wallet not found")
	}

	// Proses topup dengan DB transaction
	return t.transactionRepo.CreateTopup(ctx, userID, walletID, body)
}

// CreateTransfer — proses transfer dengan validasi PIN
func (t *TransactionService) CreateTransfer(
	ctx context.Context,
	senderUserID int,
	body dto.TransferBody,
) (dto.TransferResponse, error) {

	// verifikasi PIN dengan hash
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

	// Ambil wallet pengirim
	senderWalletID, _, err := t.transactionRepo.GetWalletByUserID(ctx, senderUserID)
	if err != nil {
		return dto.TransferResponse{}, errors.New("sender wallet not found")
	}

	// Cek receiver wallet ada
	_, _, err = t.transactionRepo.GetWalletByUserID(ctx, 0)
	// langsung cek via wallet id
	var receiverBalance float64
	err = t.transactionRepo.CheckWalletExists(ctx, body.ReceiverWalletId, &receiverBalance)
	if err != nil {
		return dto.TransferResponse{}, errors.New("receiver wallet not found")
	}

	// Proses transfer dengan DB transaction
	return t.transactionRepo.CreateTransfer(ctx, senderUserID, senderWalletID, body)
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
