package repository

import (
	"context"
	"ewallet-backend/internal/dto"
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

// FindReceivers — cari penerima transfer dengan search dan pagination
func (u *UserRepository) FindReceivers(
	ctx context.Context,
	currentUserID int,
	search string,
	limit int,
	offset int,
) ([]dto.ReceiverResponse, error) {
	sql := `
		SELECT
		  u.id,
		  u.email,
		  u.fullname,
		  u.photo_path,
		  u.phone_number,
		  w.id AS wallet_id
		FROM users u
		JOIN wallet w ON w.user_id = u.id
		WHERE u.id != $1
		  AND (
		    u.fullname     ILIKE $2
		    OR u.phone_number ILIKE $2
		    OR u.email      ILIKE $2
		  )
		ORDER BY u.fullname ASC
		LIMIT  $3
		OFFSET $4
	`

	// $2 = search pattern, contoh: "%akmal%"
	searchPattern := "%" + search + "%"

	rows, err := u.db.Query(ctx, sql, currentUserID, searchPattern, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var receivers []dto.ReceiverResponse
	for rows.Next() {
		var receiver dto.ReceiverResponse
		if err := rows.Scan(
			&receiver.Id,
			&receiver.Email,
			&receiver.Fullname,
			&receiver.PhotoPath,
			&receiver.PhoneNumber,
			&receiver.WalletId,
		); err != nil {
			return nil, err
		}
		receivers = append(receivers, receiver)
	}

	return receivers, nil
}

// CountReceivers — hitung total penerima untuk pagination
func (u *UserRepository) CountReceivers(
	ctx context.Context,
	currentUserID int,
	search string,
) (int, error) {
	sql := `
		SELECT COUNT(*)
		FROM users u
		WHERE u.id != $1
		  AND (
		    u.fullname     ILIKE $2
		    OR u.phone_number ILIKE $2
		    OR u.email      ILIKE $2
		  )
	`

	searchPattern := "%" + search + "%"

	var total int
	err := u.db.QueryRow(ctx, sql, currentUserID, searchPattern).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

// Check Pin
func (u *UserRepository) FindPinByID(ctx context.Context, id int) (string, error) {
	sql := `SELECT COALESCE(pin, '') FROM users WHERE id = $1`

	var pin string
	err := u.db.QueryRow(ctx, sql, id).Scan(pin)
	if err != nil {
		return "", err
	}

	return pin, nil
}
