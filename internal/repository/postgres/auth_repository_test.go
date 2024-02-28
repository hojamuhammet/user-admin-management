package repository_test

import (
	"admin-panel/internal/domain"
	mocks "admin-panel/internal/mocks/repository"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateTokenPair(t *testing.T) {
	tests := []struct {
		name          string
		admin         *domain.Admin
		expectedError error
	}{
		{
			name: "Success",
			admin: &domain.Admin{
				ID:   1,
				Role: "admin",
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockAuthRepository)
			mockRepo.On("GenerateTokenPair", tt.admin).Return("access_token", "refresh_token", tt.expectedError)

			accessToken, refreshToken, err := mockRepo.GenerateTokenPair(tt.admin)

			if tt.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "access_token", accessToken)
				assert.Equal(t, "refresh_token", refreshToken)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestValidateRefreshToken(t *testing.T) {
	tests := []struct {
		name          string
		refreshToken  string
		mockReturn    map[string]interface{}
		expectedError error
	}{
		{
			name:          "Valid Refresh Token",
			refreshToken:  "valid_refresh_token",
			mockReturn:    map[string]interface{}{"adminID": "123"},
			expectedError: nil,
		},
		{
			name:          "Invalid Refresh Token",
			refreshToken:  "invalid_refresh_token",
			mockReturn:    nil,
			expectedError: errors.New("error parsing refresh token"),
		},
		{
			name:          "Missing AdminID Claim",
			refreshToken:  "missing_adminID_claim_refresh_token",
			mockReturn:    nil,
			expectedError: errors.New("adminID claim not found in refresh token"),
		},
		{
			name:          "Refresh Token Not Found in Database",
			refreshToken:  "not_in_db_refresh_token",
			mockReturn:    nil,
			expectedError: errors.New("refresh token not found in the database"),
		},
		{
			name:          "Error Validating Refresh Token",
			refreshToken:  "error_validating_refresh_token",
			mockReturn:    nil,
			expectedError: errors.New("error validating refresh token"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockAuthRepository)
			mockRepo.On("ValidateRefreshToken", tt.refreshToken).Return(tt.mockReturn, tt.expectedError)

			claims, err := mockRepo.ValidateRefreshToken(tt.refreshToken)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "123", claims["adminID"])
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteRefreshToken(t *testing.T) {
	tests := []struct {
		name          string
		refreshToken  string
		expectedError error
	}{
		{
			name:          "Success",
			refreshToken:  "valid_refresh_token",
			expectedError: nil,
		},
		{
			name:          "Refresh Token Not Found",
			refreshToken:  "not_in_db_refresh_token",
			expectedError: errors.New("refresh token not found in the database"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockAuthRepository)
			mockRepo.On("DeleteRefreshToken", tt.refreshToken).Return(tt.expectedError)

			err := mockRepo.DeleteRefreshToken(tt.refreshToken)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetAdminByUsername(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		mockReturn    *domain.Admin
		expectedError error
	}{
		{
			name:     "Success",
			username: "valid_username",
			mockReturn: &domain.Admin{
				ID:       1,
				Username: "valid_username",
				Password: "hashed_password",
				Role:     "admin",
			},
			expectedError: nil,
		},
		{
			name:          "Admin Not Found",
			username:      "invalid_username",
			mockReturn:    nil,
			expectedError: errors.New("admin not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockAuthRepository)
			mockRepo.On("GetAdminByUsername", tt.username).Return(tt.mockReturn, tt.expectedError)

			admin, err := mockRepo.GetAdminByUsername(tt.username)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockReturn, admin)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
