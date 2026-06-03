package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenRepository struct {
	rdb *redis.Client
}

func NewTokenRepository(rdb *redis.Client) *TokenRepository {
	return &TokenRepository{rdb: rdb}
}

func tokenKey(token string) string {
	return fmt.Sprintf("token:%s", token)
}

// Add — simpan token ke Redis saat login
// TTL = durasi expired token (otomatis terhapus saat expired)
func (t *TokenRepository) Add(ctx context.Context, userID int, token string, expiredAt time.Time) error {
	ttl := time.Until(expiredAt)

	return t.rdb.Set(ctx,
		tokenKey(token),
		userID,
		ttl,
	).Err()
}

// Remove — hapus token dari Redis saat logout
func (t *TokenRepository) Remove(ctx context.Context, token string) error {
	return t.rdb.Del(ctx, tokenKey(token)).Err()
}

// IsWhitelisted — cek apakah token ada di Redis
func (t *TokenRepository) IsWhitelisted(ctx context.Context, token string) (bool, error) {
	result, err := t.rdb.Exists(ctx, tokenKey(token)).Result()
	if err != nil {
		return false, err
	}
	// result = 1 jika ada, 0 jika tidak ada
	return result == 1, nil
}

// GetUserID — ambil userID dari token di Redis
func (t *TokenRepository) GetUserID(ctx context.Context, token string) (int, error) {
	val, err := t.rdb.Get(ctx, tokenKey(token)).Int()
	if err != nil {
		return 0, err
	}
	return val, nil
}
