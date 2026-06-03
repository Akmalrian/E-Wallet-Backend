package repository

import (
	"context"
	"ewallet-backend/internal/dto"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// CreateTopup — buat transaksi topup dengan DB transaction
func (t *TransactionRepository) CreateTopup(
	ctx context.Context,
	userID int,
	walletID int,
	body dto.TopupBody,
) (dto.TopupResponse, error) {

	// Hitung tax dan total
	taxAmount := body.OrderAmount * 0.1
	totalAmount := body.OrderAmount + taxAmount

	tx, err := t.db.Begin(ctx)
	if err != nil {
		return dto.TopupResponse{}, err
	}
	defer tx.Rollback(ctx)

	var transactionId int
	err = tx.QueryRow(ctx, `
		INSERT INTO transactions (user_id, type, amount, status)
		VALUES ($1, 'topup', $2, 'success')
		RETURNING id
	`, userID, body.OrderAmount).Scan(&transactionId)
	if err != nil {
		return dto.TopupResponse{}, err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO topup_details (
			transaction_id, wallet_id, payment_method_id,
			order_amount, tax_amount, delivery_fee, total_amount
		) VALUES ($1, $2, $3, $4, $5, 0, $6)
	`, transactionId, walletID, body.PaymentMethodId,
		body.OrderAmount, taxAmount, totalAmount)
	if err != nil {
		return dto.TopupResponse{}, err
	}

	_, err = tx.Exec(ctx, `
		UPDATE wallet
		SET balance    = balance + $1,
		    updated_at = NOW()
		WHERE id = $2
	`, body.OrderAmount, walletID)
	if err != nil {
		return dto.TopupResponse{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.TopupResponse{}, err
	}

	var paymentName string
	t.db.QueryRow(ctx,
		"SELECT payment_name FROM payment_methods WHERE id = $1",
		body.PaymentMethodId,
	).Scan(&paymentName)

	return dto.TopupResponse{
		TransactionId: transactionId,
		PaymentMethod: paymentName,
		OrderAmount:   body.OrderAmount,
		TaxAmount:     taxAmount,
		TotalAmount:   totalAmount,
		Status:        "success",
	}, nil
}

// CreateTransfer — insert transaksi untuk pengirim DAN penerima
func (t *TransactionRepository) CreateTransfer(
	ctx context.Context,
	senderUserID int,
	senderWalletID int,
	body dto.TransferBody,
) (dto.TransferResponse, error) {

	tx, err := t.db.Begin(ctx)
	if err != nil {
		return dto.TransferResponse{}, err
	}
	defer tx.Rollback(ctx)

	// Step 1: Cek dan lock saldo pengirim
	var balance float64
	err = tx.QueryRow(ctx, `
		SELECT balance FROM wallet
		WHERE id = $1
		FOR UPDATE
	`, senderWalletID).Scan(&balance)
	if err != nil {
		return dto.TransferResponse{}, err
	}

	if balance < body.Amount {
		return dto.TransferResponse{}, fmt.Errorf("insufficient balance")
	}

	// Step 2: Ambil user_id penerima dari wallet id
	var receiverUserID int
	err = tx.QueryRow(ctx, `
		SELECT user_id FROM wallet WHERE id = $1
	`, body.ReceiverWalletId).Scan(&receiverUserID)
	if err != nil {
		return dto.TransferResponse{}, fmt.Errorf("receiver wallet not found")
	}

	// Step 3: Insert transactions untuk PENGIRIM (type: transfer)
	var senderTransactionId int
	err = tx.QueryRow(ctx, `
		INSERT INTO transactions (user_id, type, amount, status)
		VALUES ($1, 'transfer', $2, 'success')
		RETURNING id
	`, senderUserID, body.Amount).Scan(&senderTransactionId)
	if err != nil {
		return dto.TransferResponse{}, err
	}

	// Step 4: Insert transactions untuk PENERIMA (type: receive)
	var receiverTransactionId int
	err = tx.QueryRow(ctx, `
		INSERT INTO transactions (user_id, type, amount, status)
		VALUES ($1, 'receive', $2, 'success')
		RETURNING id
	`, receiverUserID, body.Amount).Scan(&receiverTransactionId)
	if err != nil {
		return dto.TransferResponse{}, err
	}

	// Step 5: Insert transfer_details untuk transaksi pengirim
	_, err = tx.Exec(ctx, `
		INSERT INTO transfer_details (
			transaction_id, sender_wallet_id,
			receiver_wallet_id, amount, notes
		) VALUES ($1, $2, $3, $4, $5)
	`, senderTransactionId, senderWalletID,
		body.ReceiverWalletId, body.Amount, body.Notes)
	if err != nil {
		return dto.TransferResponse{}, err
	}

	// Step 6: Insert transfer_details untuk transaksi penerima
	_, err = tx.Exec(ctx, `
		INSERT INTO transfer_details (
			transaction_id, sender_wallet_id,
			receiver_wallet_id, amount, notes
		) VALUES ($1, $2, $3, $4, $5)
	`, receiverTransactionId, senderWalletID,
		body.ReceiverWalletId, body.Amount, body.Notes)
	if err != nil {
		return dto.TransferResponse{}, err
	}

	// Step 7: Kurangi saldo pengirim
	_, err = tx.Exec(ctx, `
		UPDATE wallet
		SET balance    = balance - $1,
		    updated_at = NOW()
		WHERE id = $2
	`, body.Amount, senderWalletID)
	if err != nil {
		return dto.TransferResponse{}, err
	}

	// Step 8: Tambah saldo penerima
	_, err = tx.Exec(ctx, `
		UPDATE wallet
		SET balance    = balance + $1,
		    updated_at = NOW()
		WHERE id = $2
	`, body.Amount, body.ReceiverWalletId)
	if err != nil {
		return dto.TransferResponse{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.TransferResponse{}, err
	}

	return dto.TransferResponse{
		TransactionId:    senderTransactionId,
		ReceiverWalletId: body.ReceiverWalletId,
		Amount:           body.Amount,
		Notes:            body.Notes,
		Status:           "success",
	}, nil
}

// GetHistory — ambil history transaksi dengan pagination
func (t *TransactionRepository) GetHistory(
	ctx context.Context,
	userID int,
	transType string,
	limit int,
	offset int,
) ([]dto.HistoryResponse, error) {

	query := `
    SELECT
      tr.id,
      tr.type,
      tr.amount,
      tr.status,
      td.notes,
      -- Sender info (untuk type receive)
      CASE WHEN tr.type = 'receive' THEN td.sender_wallet_id      ELSE NULL END,
      CASE WHEN tr.type = 'receive' THEN sender_user.fullname      ELSE NULL END,
      CASE WHEN tr.type = 'receive' THEN sender_user.phone_number  ELSE NULL END,
      CASE WHEN tr.type = 'receive' THEN sender_user.photo_path    ELSE NULL END, -- ← tambah
      -- Receiver info (untuk type transfer)
      CASE WHEN tr.type = 'transfer' THEN td.receiver_wallet_id     ELSE NULL END,
      CASE WHEN tr.type = 'transfer' THEN receiver_user.fullname     ELSE NULL END,
      CASE WHEN tr.type = 'transfer' THEN receiver_user.phone_number ELSE NULL END,
      CASE WHEN tr.type = 'transfer' THEN receiver_user.photo_path   ELSE NULL END, -- ← tambah
      tr.created_at
    FROM transactions tr
    LEFT JOIN transfer_details td      ON td.transaction_id   = tr.id
    LEFT JOIN wallet sender_wallet     ON sender_wallet.id    = td.sender_wallet_id
    LEFT JOIN users  sender_user       ON sender_user.id      = sender_wallet.user_id
    LEFT JOIN wallet receiver_wallet   ON receiver_wallet.id  = td.receiver_wallet_id
    LEFT JOIN users  receiver_user     ON receiver_user.id    = receiver_wallet.user_id
    WHERE tr.user_id = $1
`

	args := []interface{}{userID}
	argIndex := 2

	if transType != "" {
		query += fmt.Sprintf(" AND tr.type = $%d", argIndex)
		args = append(args, transType)
		argIndex++
	}

	query += fmt.Sprintf(
		" ORDER BY tr.created_at DESC LIMIT $%d OFFSET $%d",
		argIndex, argIndex+1,
	)
	args = append(args, limit, offset)

	rows, err := t.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []dto.HistoryResponse
	for rows.Next() {
		var trx dto.HistoryResponse

		var senderWalletId *int
		var senderFullname *string
		var senderPhone *string
		var senderPhoto *string
		var receiverWalletId *int
		var receiverFullname *string
		var receiverPhone *string
		var receiverPhoto *string

		if err := rows.Scan(
			&trx.Id,
			&trx.Type,
			&trx.Amount,
			&trx.Status,
			&trx.Notes,
			&senderWalletId,
			&senderFullname,
			&senderPhone,
			&senderPhoto,
			&receiverWalletId,
			&receiverFullname,
			&receiverPhone,
			&receiverPhoto,
			&trx.CreatedAt,
		); err != nil {
			return nil, err
		}

		// Isi SenderInfo untuk type receive
		if trx.Type == "receive" && senderWalletId != nil {
			trx.SenderInfo = &dto.SenderInfo{
				WalletId:    *senderWalletId,
				Fullname:    senderFullname,
				PhoneNumber: senderPhone,
				PhotoPath:   senderPhoto,
			}
		}

		// ✅ Isi ReceiverInfo untuk type transfer
		if trx.Type == "transfer" && receiverWalletId != nil {
			trx.ReceiverInfo = &dto.ReceiverInfo{
				WalletId:    *receiverWalletId,
				Fullname:    receiverFullname,
				PhoneNumber: receiverPhone,
				PhotoPath:   receiverPhoto,
			}
		}

		transactions = append(transactions, trx)
	}

	return transactions, nil
}

// CountHistory — hitung total transaksi untuk pagination
func (t *TransactionRepository) CountHistory(
	ctx context.Context,
	userID int,
	transType string,
) (int, error) {
	query := `SELECT COUNT(*) FROM transactions WHERE user_id = $1`
	args := []interface{}{userID}

	if transType != "" {
		query += " AND type = $2"
		args = append(args, transType)
	}

	var total int
	err := t.db.QueryRow(ctx, query, args...).Scan(&total)
	return total, err
}

// GetWalletByUserID — ambil wallet berdasarkan user id
func (t *TransactionRepository) GetWalletByUserID(ctx context.Context, userID int) (int, float64, error) {
	var walletID int
	var balance float64
	err := t.db.QueryRow(ctx, `
		SELECT id, balance FROM wallet WHERE user_id = $1
	`, userID).Scan(&walletID, &balance)
	return walletID, balance, err
}
