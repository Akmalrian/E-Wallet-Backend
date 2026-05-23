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

// UserToken — representasi tabel user_tokens di database
type UserToken struct {
	Id        int
	UserId    int
	Token     string
	ExpiredAt time.Time
	CreatedAt time.Time
}
