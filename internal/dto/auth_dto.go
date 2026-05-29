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

// ForgotPasswordBody — step 1: verifikasi email
type ForgotPasswordBody struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordBody — step 2: ganti password
type ResetPasswordBody struct {
	Email       string `json:"email"        binding:"required,email"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// EnterPinBody — request body untuk set PIN pertama kali
type EnterPinBody struct {
	Pin string `json:"pin" binding:"required,len=6"`
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
