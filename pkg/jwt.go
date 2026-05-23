package pkg

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims — data yang disimpan di dalam token
type Claims struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// GetExpiredDuration — ambil durasi expired dari env
func GetExpiredDuration() time.Duration {
	expiredHours, _ := strconv.Atoi(os.Getenv("JWT_EXPIRED_HOURS"))
	if expiredHours == 0 {
		expiredHours = 24
	}
	return time.Duration(expiredHours) * time.Hour
}

// GenerateToken — buat JWT token baru
func GenerateToken(id int, email string) (string, time.Time, error) {
	duration := GetExpiredDuration()
	expiredAt := time.Now().Add(duration)

	claims := Claims{
		Id:    id,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", time.Time{}, err
	}

	return signed, expiredAt, nil
}

// ParseToken — parse dan validasi token
func ParseToken(tokenString string) (Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		},
	)
	if err != nil {
		return Claims{}, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return Claims{}, errors.New("invalid token")
	}

	return *claims, nil
}
