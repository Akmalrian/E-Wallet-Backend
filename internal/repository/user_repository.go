package repository

import (
	"context"
	"ewallet-backend/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// FindByEmail — cari user berdasarkan email
func (u *UserRepository) FindByEmail(ctx context.Context, email string) (model.User, error) {
	sql := `
		SELECT id, email, password, pin, fullname, photo_path, phone_number, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	var user model.User
	err := u.db.QueryRow(ctx, sql, email).Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.Pin,
		&user.Fullname,
		&user.PhotoPath,
		&user.PhoneNumber,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

// FindByID — cari user berdasarkan id
func (u *UserRepository) FindByID(ctx context.Context, id int) (model.User, error) {
	sql := `
		SELECT id, email, password, pin, fullname, photo_path, phone_number, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var user model.User
	err := u.db.QueryRow(ctx, sql, id).Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.Pin,
		&user.Fullname,
		&user.PhotoPath,
		&user.PhoneNumber,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

// Create — buat user baru
func (u *UserRepository) Create(ctx context.Context, email, hashedPassword string) (int, error) {
	sql := `
		INSERT INTO users (email, password)
		VALUES ($1, $2)
		RETURNING id
	`
	var id int
	err := u.db.QueryRow(ctx, sql, email, hashedPassword).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// CreateWallet — buat wallet untuk user baru
func (u *UserRepository) CreateWallet(ctx context.Context, userID int) error {
	sql := `INSERT INTO wallet (user_id, balance) VALUES ($1, 0)`
	_, err := u.db.Exec(ctx, sql, userID)
	return err
}
