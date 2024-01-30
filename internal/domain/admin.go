package domain

import (
	"errors"
	"time"
)

type AdminsList struct {
	Admins []CommonAdminResponse `json:"admins"`
}

type Admin struct {
	ID           int32        `json:"id"`
	Username     string       `json:"username"`
	Password     string       `json:"password"`
	Role         string       `json:"role"`
	RefreshToken RefreshToken `json:"refresh_token"`
}

type RefreshToken struct {
	Token          string    `json:"token"`
	ExpirationTime time.Time `json:"expiration_time"`
	CreatedAt      time.Time `json:"created_at"`
}

type CreateAdminRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UpdateAdminRequest struct {
	ID       int32  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type AdminListResponse struct {
	Admins      *AdminsList `json:"admins"`
	CurrentPage int         `json:"currentPage"`
	PrevPage    int         `json:"previousPage"`
	NextPage    int         `json:"nextPage"`
}

type CommonAdminResponse struct {
	ID       int32  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrAdminNotFound        = errors.New("admin not found")
	ErrInvalidRefreshToken  = errors.New("invalid refresh token")
	ErrAdminAlreadyExists   = errors.New("admin already exists")
	ErrAdminCannotBeDeleted = errors.New("super admin cannot be deleted")
)
