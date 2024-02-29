package handlers

import (
	mocks "admin-panel/internal/mocks/service"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginHandler(t *testing.T) {
	mockAuthService := new(mocks.MockAuthService)
	handler := AuthHandler{
		AuthService: mockAuthService,
	}

	testCases := []struct {
		name           string
		username       string
		password       string
		mockReturn     []interface{}
		expectedStatus int
	}{
		{
			name:           "Successful login",
			username:       "admin",
			password:       "password",
			mockReturn:     []interface{}{"access_token", "refresh_token", nil},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid credentials",
			username:       "admin",
			password:       "wrong_password",
			mockReturn:     []interface{}{"", "", errors.New("invalid credentials")},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			loginRequest := LoginRequest{
				Username: tc.username,
				Password: tc.password,
			}
			requestBody, _ := json.Marshal(loginRequest)
			req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
			rr := httptest.NewRecorder()

			mockAuthService.On("LoginAdmin", tc.username, tc.password).Return(tc.mockReturn...)

			handler.LoginHandler(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			mockAuthService.AssertExpectations(t)
		})
	}
}

func TestRefreshTokensHandler(t *testing.T) {
	mockAuthService := new(mocks.MockAuthService)
	handler := AuthHandler{
		AuthService: mockAuthService,
	}

	testCases := []struct {
		name           string
		refreshToken   string
		mockReturn     []interface{}
		expectedStatus int
	}{
		{
			name:           "Successful token refresh",
			refreshToken:   "valid_refresh_token",
			mockReturn:     []interface{}{"new_access_token", "new_refresh_token", nil},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid refresh token",
			refreshToken:   "invalid_refresh_token",
			mockReturn:     []interface{}{"", "", errors.New("invalid refresh token")},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Expired refresh token",
			refreshToken:   "Bearer expired_refresh_token",
			mockReturn:     []interface{}{"", "", errors.New("expired refresh token")},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/refresh", nil)
			req.Header.Add("Authorization", tc.refreshToken)
			rr := httptest.NewRecorder()

			mockAuthService.On("RefreshTokens", strings.TrimPrefix(tc.refreshToken, "Bearer ")).Return(tc.mockReturn...)

			handler.RefreshTokensHandler(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			mockAuthService.AssertExpectations(t)
		})
	}
}
