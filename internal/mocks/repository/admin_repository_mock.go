package mocks

import (
	"admin-panel/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockAdminRepository struct {
	mock.Mock
}

func (m *MockAdminRepository) GetAllAdmins(page, pageSize int) (*domain.AdminsList, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).(*domain.AdminsList), args.Error(1)
}

func (m *MockAdminRepository) GetTotalAdminsCount() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockAdminRepository) GetAdminByID(id int32) (*domain.GetAdminResponse, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.GetAdminResponse), args.Error(1)
}

func (m *MockAdminRepository) CreateAdmin(request *domain.CreateAdminRequest) (*domain.CreateAdminResponse, error) {
	args := m.Called(request)
	return args.Get(0).(*domain.CreateAdminResponse), args.Error(1)
}

func (m *MockAdminRepository) UpdateAdmin(id int32, request *domain.UpdateAdminRequest) (*domain.UpdateAdminResponse, error) {
	args := m.Called(id, request)
	return args.Get(0).(*domain.UpdateAdminResponse), args.Error(1)
}

func (m *MockAdminRepository) DeleteAdmin(id int32) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAdminRepository) BlockAdmin(id int32) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAdminRepository) UnblockAdmin(id int32) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAdminRepository) SearchAdmins(query string, page, pageSize int) (*domain.AdminsList, error) {
	args := m.Called(query, page, pageSize)
	return args.Get(0).(*domain.AdminsList), args.Error(1)
}
