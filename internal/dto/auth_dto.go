package dto

import "time"

// ── Request ──────────────────────────────────

type RegisterBody struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginBody struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// ── Response ─────────────────────────────────

type UserResponse struct {
	Id          int       `json:"id"`
	Email       string    `json:"email"`
	Fullname    *string   `json:"fullname"`
	PhotoPath   *string   `json:"photo_path"`
	PhoneNumber *string   `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
