package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenRepository struct {
	db *pgxpool.Pool
}

func NewTokenRepository(db *pgxpool.Pool) *TokenRepository {
	return &TokenRepository{db: db}
}

// Add — simpan token ke database saat login
func (t *TokenRepository) Add(ctx context.Context, userID int, token string, expiredAt time.Time) error {
	sql := `
		INSERT INTO user_tokens (user_id, token, expired_at)
		VALUES ($1, $2, $3)
	`
	_, err := t.db.Exec(ctx, sql, userID, token, expiredAt)
	return err
}

// Remove — hapus token saat logout
func (t *TokenRepository) Remove(ctx context.Context, token string) error {
	sql := `DELETE FROM user_tokens WHERE token = $1`
	_, err := t.db.Exec(ctx, sql, token)
	return err
}

// IsWhitelisted — cek apakah token aktif di database
func (t *TokenRepository) IsWhitelisted(ctx context.Context, token string) (bool, error) {
	sql := `
		SELECT COUNT(*)
		FROM user_tokens
		WHERE token      = $1
		  AND expired_at > NOW()
	`
	var count int
	err := t.db.QueryRow(ctx, sql, token).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
