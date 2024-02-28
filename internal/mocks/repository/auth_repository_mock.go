package mocks

import (
	"admin-panel/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) GenerateTokenPair(admin *domain.Admin) (string, string, error) {
	args := m.Called(admin)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockAuthRepository) ValidateRefreshToken(refreshToken string) (map[string]interface{}, error) {
	args := m.Called(refreshToken)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockAuthRepository) DeleteRefreshToken(refreshToken string) error {
	args := m.Called(refreshToken)
	return args.Error(0)
}

func (m *MockAuthRepository) GetAdminByUsername(username string) (*domain.Admin, error) {
	args := m.Called(username)
	return args.Get(0).(*domain.Admin), args.Error(1)
}
