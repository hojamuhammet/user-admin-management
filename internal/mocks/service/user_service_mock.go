package mocks

import (
	"admin-panel/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetAllUsers(page, pageSize int) (*domain.UsersList, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).(*domain.UsersList), args.Error(1)
}

func (m *MockUserService) GetTotalUsersCount() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockUserService) GetUserByID(id int32) (*domain.GetUserResponse, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.GetUserResponse), args.Error(1)
}

func (m *MockUserService) CreateUser(request *domain.CreateUserRequest) (*domain.CreateUserResponse, error) {
	args := m.Called(request)
	return args.Get(0).(*domain.CreateUserResponse), args.Error(1)
}

func (m *MockUserService) UpdateUser(id int32, request *domain.UpdateUserRequest) (*domain.UpdateUserResponse, error) {
	args := m.Called(id, request)
	return args.Get(0).(*domain.UpdateUserResponse), args.Error(1)
}

func (m *MockUserService) DeleteUser(id int32) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserService) BlockUser(id int32) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserService) UnblockUser(id int32) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserService) SearchUsers(query string, page, pageSize int) (*domain.UsersList, error) {
	args := m.Called(query, page, pageSize)
	return args.Get(0).(*domain.UsersList), args.Error(1)
}
