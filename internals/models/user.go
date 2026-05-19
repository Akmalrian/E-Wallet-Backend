package models

import "time"

type User struct {
	ID          int
	Email       string
	Password    string
	Pin         *string
	Fullname    *string
	PhotoPath   *string
	PhoneNumber *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
