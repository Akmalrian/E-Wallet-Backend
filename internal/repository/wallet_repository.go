package repository

import (
	"context"
	"ewallet-backend/internal/dto"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WalletRepository struct {
	db *pgxpool.Pool
}

func NewWalletRepository(db *pgxpool.Pool) *WalletRepository {
	return &WalletRepository{db: db}
}

func (w *WalletRepository) GetDashboardInfo(ctx context.Context, userID int) (dto.DashboardResponse, error) {
	sql := `SELECT
  w.balance,

  (
    SELECT COALESCE(SUM(t.amount), 0)
    FROM transactions t
    WHERE t.user_id = $1
      AND t.type    = 'topup'
      AND t.status  = 'success'
  ) AS total_income,

  (
    SELECT COALESCE(SUM(t.amount), 0)
    FROM transactions t
    WHERE t.user_id = $1
      AND t.type    = 'transfer'
      AND t.status  = 'success'
  ) AS total_expense
   
	FROM wallet w
	WHERE w.user_id = $1;`

	var dashboard dto.DashboardResponse

	err := w.db.QueryRow(ctx, sql, userID).Scan(
		&dashboard.Balance,
		&dashboard.TotalIncome,
		&dashboard.TotalExpense,
	)
	if err != nil {
		return dto.DashboardResponse{}, err
	}
	return dashboard, nil

}
