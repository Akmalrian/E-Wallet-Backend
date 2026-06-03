package model

import "time"

// User — representasi tabel users di database
type User struct {
	Id          int
	Email       string
	Password    string
	Pin         *string
	Fullname    *string
	PhotoPath   *string
	PhoneNumber *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
