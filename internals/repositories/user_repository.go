package repositories

import (
	"database/sql"
	"ewallet-backend/internals/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// FindByEmail — cari user berdasarkan email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, password, pin, fullname, photo_path, phone_number, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user models.User
	err := r.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Pin,
		&user.Fullname,
		&user.PhotoPath,
		&user.PhoneNumber,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // user tidak ditemukan
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Create — buat user baru
func (r *UserRepository) Create(email, hashedPassword string) (int, error) {
	var userID int
	query := `
		INSERT INTO users (email, password)
		VALUES ($1, $2)
		RETURNING id
	`
	err := r.DB.QueryRow(query, email, hashedPassword).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

// CreateWallet — buat wallet untuk user baru
func (r *UserRepository) CreateWallet(userID int) error {
	query := `INSERT INTO wallet (user_id, balance) VALUES ($1, 0)`
	_, err := r.DB.Exec(query, userID)
	return err
}

// FindByID — cari user berdasarkan id
func (r *UserRepository) FindByID(id int) (*models.User, error) {
	query := `
		SELECT id, email, password, pin, fullname, photo_path, phone_number, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := r.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Pin,
		&user.Fullname,
		&user.PhotoPath,
		&user.PhoneNumber,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}
