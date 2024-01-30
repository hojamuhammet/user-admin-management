package domain

import (
	"errors"
	"time"
)

type Date struct {
	Day   int32 `json:"day"`
	Month int32 `json:"month"`
	Year  int32 `json:"year"`
}

type UserID struct {
	ID int32 `json:"id"`
}

type UsersList struct {
	Users []CommonUserResponse `json:"users"`
}

type UsersListResponse struct {
	Users       *UsersList `json:"users"`
	CurrentPage int        `json:"currentPage"`
	PrevPage    int        `json:"previousPage"`
	NextPage    int        `json:"nextPage"`
}

// CommonUserResponse captures the common properties for GetUserResponse, CreateUserResponse and UpdateUserResponse
type CommonUserResponse struct {
	ID               int32     `json:"id"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	PhoneNumber      string    `json:"phone_number"`
	Blocked          bool      `json:"blocked"`
	Gender           string    `json:"gender"`
	RegistrationDate time.Time `json:"registration_date"`
	DateOfBirth      Date      `json:"date_of_birth"`
	Location         string    `json:"location"`
	Email            string    `json:"email"`
	ProfilePhotoURL  string    `json:"profile_photo_url"`
}

type GetUserResponse CommonUserResponse

type CreateUserRequest struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	PhoneNumber     string `json:"phone_number"`
	Gender          string `json:"gender"`
	DateOfBirth     Date   `json:"date_of_birth"`
	Location        string `json:"location"`
	Email           string `json:"email"`
	ProfilePhotoURL string `json:"profile_photo_url"`
}

type CreateUserResponse CommonUserResponse

type UpdateUserRequest struct {
	ID              int32  `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	PhoneNumber     string `json:"phone_number"`
	Gender          string `json:"gender"`
	DateOfBirth     Date   `json:"date_of_birth"`
	Location        string `json:"location"`
	Email           string `json:"email"`
	ProfilePhotoURL string `json:"profile_photo_url"`
}

type UpdateUserResponse CommonUserResponse

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)
