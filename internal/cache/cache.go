package cache

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// Set — simpan data ke Redis
func Set(ctx context.Context, rdb *redis.Client, key string, value any, ttl time.Duration) {
	data, err := json.Marshal(value)
	if err != nil {
		log.Println("Cache marshal error:", err.Error())
		return
	}

	if err := rdb.Set(ctx, key, string(data), ttl).Err(); err != nil {
		log.Println("Cache set error:", err.Error())
	}
}

// Get — ambil data dari Redis
func Get(ctx context.Context, rdb *redis.Client, key string, dest any) bool {
	data, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			log.Println("Cache get error:", err.Error())
		}
		return false
	}

	if err := json.Unmarshal([]byte(data), dest); err != nil {
		log.Println("Cache unmarshal error:", err.Error())
		return false
	}

	return true
}

// Delete — hapus data dari Redis
func Delete(ctx context.Context, rdb *redis.Client, key string) {
	if err := rdb.Del(ctx, key).Err(); err != nil {
		log.Println("Cache delete error:", err.Error())
	}
}
