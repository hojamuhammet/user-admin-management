package mocks

import (
	"admin-panel/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockAdminService struct {
	mock.Mock
}

func (m *MockAdminService) GetAllAdmins(page, pageSize int) (*domain.AdminsList, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).(*domain.AdminsList), args.Error(1)
}

func (m *MockAdminService) GetTotalAdminsCount() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockAdminService) GetAdminByID(id int32) (*domain.GetAdminResponse, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.GetAdminResponse), args.Error(1)
}

func (m *MockAdminService) CreateAdmin(request *domain.CreateAdminRequest) (*domain.CreateAdminResponse, error) {
	args := m.Called(request)
	return args.Get(0).(*domain.CreateAdminResponse), args.Error(1)
}

func (m *MockAdminService) UpdateAdmin(id int32, request *domain.UpdateAdminRequest) (*domain.UpdateAdminResponse, error) {
	args := m.Called(id, request)
	return args.Get(0).(*domain.UpdateAdminResponse), args.Error(1)
}

func (m *MockAdminService) DeleteAdmin(id int32) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAdminService) BlockAdmin(id int32) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAdminService) UnblockAdmin(id int32) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAdminService) SearchAdmins(query string, page, pageSize int) (*domain.AdminsList, error) {
	args := m.Called(query, page, pageSize)
	return args.Get(0).(*domain.AdminsList), args.Error(1)
}
