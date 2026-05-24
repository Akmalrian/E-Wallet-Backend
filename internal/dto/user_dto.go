package dto

import "time"

// GetProfileResponse — response data profile user
type GetProfileResponse struct {
	Id          int       `json:"id"`
	Email       string    `json:"email"`
	Fullname    *string   `json:"fullname"`
	PhotoPath   *string   `json:"photo_path"`
	PhoneNumber *string   `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
}

type ReceiverResponse struct {
	Id          int     `json:"id"`
	Email       string  `json:"email"`
	Fullname    *string `json:"fullname"`
	PhotoPath   *string `json:"photo_path"`
	PhoneNumber *string `json:"phone_number"`
	WalletId    int     `json:"wallet_id"`
}

type PaginationMeta struct {
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
	TotalData   int `json:"total_data"`
	Limit       int `json:"limit"`
}

type ReceiverListResponse struct {
	Receivers []ReceiverResponse `json:"receivers"`
	Meta      PaginationMeta     `json:"meta"`
}

type CheckPinBody struct {
	Pin string `json:"pin" binding:"required,len=6,numeric"`
}
