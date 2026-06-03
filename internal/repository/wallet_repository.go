package repository

import (
	"context"
	"ewallet-backend/internal/dto"
	"fmt"
	"time"

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
      AND t.type    = 'receive'
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

// GetGraphData — ambil data graph berdasarkan filter
func (w *WalletRepository) GetGraphData(
	ctx context.Context,
	userID int,
	graphType string, // "income", "expense", "both"
	startDate time.Time,
	endDate time.Time,
) ([]dto.GraphPoint, error) {

	// Build query berdasarkan graphType
	// Selalu group by hari agar bisa lihat trend per hari
	var incomeExpr, expenseExpr string

	switch graphType {
	case "income":
		// Hanya tampilkan income, expense selalu 0
		incomeExpr = `SUM(CASE WHEN type = 'topup'    AND status = 'success' THEN amount ELSE 0 END)`
		expenseExpr = `0`
	case "expense":
		// Hanya tampilkan expense, income selalu 0
		incomeExpr = `0`
		expenseExpr = `SUM(CASE WHEN type = 'transfer' AND status = 'success' THEN amount ELSE 0 END)`
	default:
		// "both" → tampilkan keduanya
		incomeExpr = `SUM(CASE WHEN type = 'topup'    AND status = 'success' THEN amount ELSE 0 END)`
		expenseExpr = `SUM(CASE WHEN type = 'transfer' AND status = 'success' THEN amount ELSE 0 END)`
	}

	sql := fmt.Sprintf(`
		SELECT
		  TO_CHAR(DATE_TRUNC('day', created_at), 'DD Mon') AS label,
		  %s AS income,
		  %s AS expense
		FROM transactions
		WHERE user_id    = $1
		  AND created_at >= $2
		  AND created_at <= $3
		GROUP BY DATE_TRUNC('day', created_at)
		ORDER BY DATE_TRUNC('day', created_at) ASC
	`, incomeExpr, expenseExpr)

	rows, err := w.db.Query(ctx, sql, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []dto.GraphPoint
	for rows.Next() {
		var point dto.GraphPoint
		if err := rows.Scan(
			&point.Label,
			&point.Income,
			&point.Expense,
		); err != nil {
			return nil, err
		}
		points = append(points, point)
	}

	return points, nil
}

// CheckWalletExists — cek apakah wallet ada berdasarkan wallet id
func (t *TransactionRepository) CheckWalletExists(
	ctx context.Context,
	walletID int,
	balance *float64,
) error {
	return t.db.QueryRow(ctx,
		"SELECT balance FROM wallet WHERE id = $1",
		walletID,
	).Scan(balance)
}
