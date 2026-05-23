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
