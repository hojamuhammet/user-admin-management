package domain

import (
	"errors"
	"time"
)

type Admin struct {
	ID       int32  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type RefreshToken struct {
    ID             int    `json:"id"`
    AdminID        int    `json:"admin_id"`
    Token          string `json:"token"`
    ExpirationTime time.Time   `json:"expiration_time"`
    CreatedAt      time.Time   `json:"created_at"`
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

var (
	ErrAdminNotFound = errors.New("admin not found")
)