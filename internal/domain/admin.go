package domain

import (
	"time"
)

type AdminsList struct {
	Admins []GetAdminResponse `json:"admins"`
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

type AdminListResponse struct {
	Admins      *AdminsList `json:"admins"`
	CurrentPage int         `json:"currentPage"`
	PrevPage    int         `json:"previousPage"`
	NextPage    int         `json:"nextPage"`
	FirstPage   int         `json:"firstPage"`
	LastPage    int         `json:"lastPage"`
}

type CommonAdminRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type CommonAdminResponse struct {
	ID       int32  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type GetAdminResponse CommonAdminResponse

type CreateAdminRequest CommonAdminRequest

type CreateAdminResponse CommonAdminResponse

type UpdateAdminRequest CommonAdminRequest

type UpdateAdminResponse CommonAdminResponse
