package dto

import "time"

// GetProfileResponse
type GetProfileResponse struct {
	ID          int       `json:"id"`
	Email       string    `json:"email"`
	Fullname    *string   `json:"fullname"`
	PhotoPath   *string   `json:"photo_path"`
	PhoneNumber *string   `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
}
