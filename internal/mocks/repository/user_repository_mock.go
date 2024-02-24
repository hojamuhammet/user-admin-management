package mocks

import (
	"admin-panel/internal/domain"
	"errors"
)

type UserRepositoryMock struct {
	GetAllUsersFunc func(page, pageSize int) (*domain.UsersList, error)
	GetUserByIDFunc func(id int32) (*domain.GetUserResponse, error)
	CreateUserFunc  func(request *domain.CreateUserRequest) (*domain.CreateUserResponse, error)
	UpdateUserFunc  func(id int32, request *domain.UpdateUserRequest) (*domain.UpdateUserResponse, error)
	DeleteUserFunc  func(id int32) error
	BlockUserFunc   func(id int32) error
	UnblockUserFunc func(id int32) error
	SearchUsersFunc func(query string, page, pageSize int) (*domain.UsersList, error)
}

func NewUserRepositoryMock() *UserRepositoryMock {
	return &UserRepositoryMock{}
}

func (m *UserRepositoryMock) GetAllUsers(page, pageSize int) (*domain.UsersList, error) {
	if m.GetAllUsersFunc != nil {
		return m.GetAllUsersFunc(page, pageSize)
	}
	return nil, errors.New("GetAllUsers is not implemented in the mock")
}

func (m *UserRepositoryMock) GetUserByID(id int32) (*domain.GetUserResponse, error) {
	if m.GetUserByIDFunc != nil {
		return m.GetUserByIDFunc(id)
	}
	return nil, errors.New("GetUserByID is not implemented in the mock")
}

func (m *UserRepositoryMock) CreateUser(request *domain.CreateUserRequest) (*domain.CreateUserResponse, error) {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(request)
	}
	return nil, errors.New("CreateUser is not implemented in the mock")
}

func (m *UserRepositoryMock) UpdateUser(id int32, request *domain.UpdateUserRequest) (*domain.UpdateUserResponse, error) {
	if m.UpdateUserFunc != nil {
		return m.UpdateUserFunc(id, request)
	}
	return nil, errors.New("UpdateUser is not implemented in the mock")
}

func (m *UserRepositoryMock) DeleteUser(id int32) error {
	if m.DeleteUserFunc != nil {
		return m.DeleteUserFunc(id)
	}
	return errors.New("DeleteUser is not implemented in the mock")
}

func (m *UserRepositoryMock) BlockUser(id int32) error {
	if m.BlockUserFunc != nil {
		return m.BlockUserFunc(id)
	}
	return nil
}

func (m *UserRepositoryMock) UnblockUser(id int32) error {
	if m.UnblockUserFunc != nil {
		return m.UnblockUserFunc(id)
	}
	return nil
}

func (m *UserRepositoryMock) SearchUsers(query string, page, pageSize int) (*domain.UsersList, error) {
	if m.SearchUsersFunc != nil {
		return m.SearchUsersFunc(query, page, pageSize)
	}
	return nil, errors.New("SearchUsers is not implemented in the mock")
}
