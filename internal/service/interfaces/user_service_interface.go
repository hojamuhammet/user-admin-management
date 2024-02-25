package service

import "admin-panel/internal/domain"

type UserService interface {
	GetAllUsers(page, pageSize int) (*domain.UsersList, error)
	GetTotalUsersCount() (int, error)
	GetUserByID(id int32) (*domain.GetUserResponse, error)
	CreateUser(request *domain.CreateUserRequest) (*domain.CreateUserResponse, error)
	UpdateUser(id int32, request *domain.UpdateUserRequest) (*domain.UpdateUserResponse, error)
	DeleteUser(id int32) error
	BlockUser(id int32) error
	UnblockUser(id int32) error
	SearchUsers(query string, page, pageSize int) (*domain.UsersList, error)
}
