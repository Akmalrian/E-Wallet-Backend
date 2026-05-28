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

	// ✅ Mulai DB transaction
	tx, err := t.db.Begin(ctx)
	if err != nil {
		return dto.TopupResponse{}, err
	}
	// ✅ Defer rollback — otomatis rollback jika ada error
	defer tx.Rollback(ctx)

	// Step 1: Insert ke tabel transactions
	var transactionId int
	err = tx.QueryRow(ctx, `
		INSERT INTO transactions (user_id, type, amount, status)
		VALUES ($1, 'topup', $2, 'success')
		RETURNING id
	`, userID, body.OrderAmount).Scan(&transactionId)
	if err != nil {
		return dto.TopupResponse{}, err
	}

	// Step 2: Insert ke tabel topup_details
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

	// Step 3: Tambah balance wallet
	_, err = tx.Exec(ctx, `
		UPDATE wallet
		SET balance    = balance + $1,
		    updated_at = NOW()
		WHERE id = $2
	`, body.OrderAmount, walletID)
	if err != nil {
		return dto.TopupResponse{}, err
	}

	// ✅ Commit — semua query berhasil, simpan perubahan
	if err := tx.Commit(ctx); err != nil {
		return dto.TopupResponse{}, err
	}

	// Ambil nama payment method
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

// CreateTransfer — buat transaksi transfer dengan DB transaction
func (t *TransactionRepository) CreateTransfer(
	ctx context.Context,
	senderUserID int,
	senderWalletID int,
	body dto.TransferBody,
) (dto.TransferResponse, error) {

	// ✅ Mulai DB transaction
	tx, err := t.db.Begin(ctx)
	if err != nil {
		return dto.TransferResponse{}, err
	}
	defer tx.Rollback(ctx)

	// Step 1: Cek dan lock saldo pengirim
	// FOR UPDATE → lock row agar tidak diubah proses lain
	var balance float64
	err = tx.QueryRow(ctx, `
		SELECT balance FROM wallet
		WHERE id = $1
		FOR UPDATE
	`, senderWalletID).Scan(&balance)
	if err != nil {
		return dto.TransferResponse{}, err
	}

	// Step 2: Validasi saldo mencukupi
	if balance < body.Amount {
		return dto.TransferResponse{}, fmt.Errorf("insufficient balance")
	}

	// Step 3: Insert ke tabel transactions
	var transactionId int
	err = tx.QueryRow(ctx, `
		INSERT INTO transactions (user_id, type, amount, status)
		VALUES ($1, 'transfer', $2, 'success')
		RETURNING id
	`, senderUserID, body.Amount).Scan(&transactionId)
	if err != nil {
		return dto.TransferResponse{}, err
	}

	// Step 4: Insert ke tabel transfer_details
	_, err = tx.Exec(ctx, `
		INSERT INTO transfer_details (
			transaction_id, sender_wallet_id,
			receiver_wallet_id, amount, notes
		) VALUES ($1, $2, $3, $4, $5)
	`, transactionId, senderWalletID,
		body.ReceiverWalletId, body.Amount, body.Notes)
	if err != nil {
		return dto.TransferResponse{}, err
	}

	// Step 5: Kurangi saldo pengirim
	_, err = tx.Exec(ctx, `
		UPDATE wallet
		SET balance    = balance - $1,
		    updated_at = NOW()
		WHERE id = $2
	`, body.Amount, senderWalletID)
	if err != nil {
		return dto.TransferResponse{}, err
	}

	// Step 6: Tambah saldo penerima
	_, err = tx.Exec(ctx, `
		UPDATE wallet
		SET balance    = balance + $1,
		    updated_at = NOW()
		WHERE id = $2
	`, body.Amount, body.ReceiverWalletId)
	if err != nil {
		return dto.TransferResponse{}, err
	}

	// ✅ Commit — semua berhasil
	if err := tx.Commit(ctx); err != nil {
		return dto.TransferResponse{}, err
	}

	return dto.TransferResponse{
		TransactionId:    transactionId,
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

	// Build query dinamis berdasarkan filter type
	query := `
		SELECT tr.id, tr.type, tr.amount, tr.status, COALESCE(td.notes, 'Isi Saldo'), tr.created_at
		FROM transactions tr
		LEFT JOIN transfer_details td ON td.transaction_id = tr.id
		WHERE tr.user_id = $1
	`
	args := []interface{}{userID}
	argIndex := 2

	// Tambah filter type jika ada
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
		if err := rows.Scan(
			&trx.Id,
			&trx.Type,
			&trx.Amount,
			&trx.Status,
			&trx.Notes,
			&trx.CreatedAt,
		); err != nil {
			return nil, err
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
