package mocks

import (
	"admin-panel/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetAllUsers(page, pageSize int) (*domain.UsersList, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).(*domain.UsersList), args.Error(1)
}

func (m *MockUserRepository) GetTotalUsersCount() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(id int32) (*domain.GetUserResponse, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.GetUserResponse), args.Error(1)
}

func (m *MockUserRepository) CreateUser(request *domain.CreateUserRequest) (*domain.CreateUserResponse, error) {
	args := m.Called(request)
	return args.Get(0).(*domain.CreateUserResponse), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(id int32, request *domain.UpdateUserRequest) (*domain.UpdateUserResponse, error) {
	args := m.Called(id, request)
	return args.Get(0).(*domain.UpdateUserResponse), args.Error(1)
}

func (m *MockUserRepository) DeleteUser(id int32) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) BlockUser(id int32) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) UnblockUser(id int32) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) SearchUsers(query string, page, pageSize int) (*domain.UsersList, error) {
	args := m.Called(query, page, pageSize)
	return args.Get(0).(*domain.UsersList), args.Error(1)
}
