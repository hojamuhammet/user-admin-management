package service_test

import (
	"admin-panel/internal/domain"
	mocks "admin-panel/internal/mocks/repository"
	"admin-panel/internal/service"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginAdmin(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)

	testCases := []struct {
		name          string
		username      string
		password      string
		mockRepo      func() *mocks.MockAuthRepository
		expectedError error
	}{
		{
			name:     "Successful Login",
			username: "testuser",
			password: "testpass",
			mockRepo: func() *mocks.MockAuthRepository {
				mockRepo := new(mocks.MockAuthRepository)
				mockRepo.On("GetAdminByUsername", "testuser").Return(&domain.Admin{Username: "testuser", Password: string(hashedPassword)}, nil) // password is "testpass"
				mockRepo.On("GenerateTokenPair", mock.AnythingOfType("*domain.Admin")).Return("mockAccessToken", "mockRefreshToken", nil)
				return mockRepo
			},
			expectedError: nil,
		},
		{
			name:     "Invalid Password",
			username: "testuser",
			password: "wrongpass",
			mockRepo: func() *mocks.MockAuthRepository {
				mockRepo := new(mocks.MockAuthRepository)
				mockRepo.On("GetAdminByUsername", "testuser").Return(&domain.Admin{Username: "testuser", Password: string(hashedPassword)}, nil) // password is "testpass"
				return mockRepo
			},
			expectedError: errors.New("invalid credentials"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := service.NewAuthService(tc.mockRepo())

			_, _, err := s.LoginAdmin(tc.username, tc.password)

			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestRefreshTokens(t *testing.T) {
	testCases := []struct {
		name          string
		refreshToken  string
		mockRepo      func() *mocks.MockAuthRepository
		expectedError error
	}{
		{
			name:         "Successful Refresh",
			refreshToken: "validRefreshToken",
			mockRepo: func() *mocks.MockAuthRepository {
				mockRepo := new(mocks.MockAuthRepository)
				mockRepo.On("ValidateRefreshToken", "validRefreshToken").Return(map[string]interface{}{"adminID": float64(1)}, nil)
				mockRepo.On("GetAdminByID", 1).Return(&domain.Admin{ID: 1, Username: "testuser"}, nil)
				mockRepo.On("GenerateTokenPair", mock.AnythingOfType("*domain.Admin")).Return("newAccessToken", "newRefreshToken", nil)
				return mockRepo
			},
			expectedError: nil,
		},
		{
			name:         "Invalid Refresh Token",
			refreshToken: "invalidRefreshToken",
			mockRepo: func() *mocks.MockAuthRepository {
				mockRepo := new(mocks.MockAuthRepository)
				mockRepo.On("ValidateRefreshToken", "invalidRefreshToken").Return(make(map[string]interface{}), errors.New("invalid refresh token"))
				return mockRepo
			},
			expectedError: errors.New("invalid refresh token"),
		},
		{
			name:         "AdminID Claim Not Found",
			refreshToken: "validRefreshToken",
			mockRepo: func() *mocks.MockAuthRepository {
				mockRepo := new(mocks.MockAuthRepository)
				mockRepo.On("ValidateRefreshToken", "validRefreshToken").Return(make(map[string]interface{}), nil)
				return mockRepo
			},
			expectedError: errors.New("invalid refresh token"),
		},
		{
			name:         "GetAdminByID Returns Error",
			refreshToken: "validRefreshToken",
			mockRepo: func() *mocks.MockAuthRepository {
				mockRepo := new(mocks.MockAuthRepository)
				mockRepo.On("ValidateRefreshToken", "validRefreshToken").Return(map[string]interface{}{"adminID": float64(1)}, nil)
				mockRepo.On("GetAdminByID", 1).Return(&domain.Admin{}, errors.New("admin not found"))
				return mockRepo
			},
			expectedError: errors.New("admin not found"),
		},
		{
			name:         "GenerateTokenPair Returns Error",
			refreshToken: "validRefreshToken",
			mockRepo: func() *mocks.MockAuthRepository {
				mockRepo := new(mocks.MockAuthRepository)
				mockRepo.On("ValidateRefreshToken", "validRefreshToken").Return(map[string]interface{}{"adminID": float64(1)}, nil)
				mockRepo.On("GetAdminByID", 1).Return(&domain.Admin{ID: 1, Username: "testuser"}, nil)
				mockRepo.On("GenerateTokenPair", mock.AnythingOfType("*domain.Admin")).Return("", "", errors.New("error generating token pair"))
				return mockRepo
			},
			expectedError: errors.New("error generating token pair"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := service.NewAuthService(tc.mockRepo())

			newAccessToken, newRefreshToken, err := s.RefreshTokens(tc.refreshToken)

			assert.Equal(t, tc.expectedError, err)
			if err == nil {
				assert.Equal(t, "newAccessToken", newAccessToken)
				assert.Equal(t, "newRefreshToken", newRefreshToken)
			}
		})
	}
}

func TestLogoutAdmin(t *testing.T) {
	testCases := []struct {
		name          string
		refreshToken  string
		mockRepo      func() *mocks.MockAuthRepository
		expectedError error
	}{
		{
			name:         "Successful Logout",
			refreshToken: "validRefreshToken",
			mockRepo: func() *mocks.MockAuthRepository {
				mockRepo := new(mocks.MockAuthRepository)
				mockRepo.On("DeleteRefreshToken", "validRefreshToken").Return(nil)
				return mockRepo
			},
			expectedError: nil,
		},
		{
			name:         "Failed Logout",
			refreshToken: "invalidRefreshToken",
			mockRepo: func() *mocks.MockAuthRepository {
				mockRepo := new(mocks.MockAuthRepository)
				mockRepo.On("DeleteRefreshToken", "invalidRefreshToken").Return(errors.New("error deleting refresh token"))
				return mockRepo
			},
			expectedError: errors.New("error deleting refresh token"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := service.NewAuthService(tc.mockRepo())

			err := s.LogoutAdmin(tc.refreshToken)

			// Assertions
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
