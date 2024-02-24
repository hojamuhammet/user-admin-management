package handlers_test

import (
	"admin-panel/internal/delivery/v1/handlers"
	"admin-panel/internal/domain"
	mocks "admin-panel/internal/mocks/service"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var dateOfBirth = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)

func TestGetUserByIDHandler(t *testing.T) {
	testCases := []struct {
		name           string
		mockReturnUser *domain.GetUserResponse
		mockReturnErr  error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",
			mockReturnUser: &domain.GetUserResponse{
				ID:              1,
				FirstName:       "Kemal",
				LastName:        "Atdayew",
				PhoneNumber:     "+99362008971",
				Gender:          "Male",
				DateOfBirth:     dateOfBirth,
				Location:        "Ashgabat",
				Email:           "atdayewkemal@gmail.com",
				ProfilePhotoURL: "https://example.com/profile.jpg"},
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"first_name":"Kemal","last_name":"Atdayew","phone_number":"+99362008971","blocked":false,"gender":"Male","registration_date":"0001-01-01T00:00:00Z","date_of_birth":"2000-01-01T00:00:00Z","location":"Ashgabat","email":"atdayewkemal@gmail.com","profile_photo_url":"https://example.com/profile.jpg"}`,
		},
		{
			name:           "NotFound",
			mockReturnUser: nil,
			mockReturnErr:  domain.ErrUserNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"status":404,"message":"User not found"}`,
		},
		{
			name:           "InternalServerError",
			mockReturnUser: nil,
			mockReturnErr:  errors.New("internal server error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"status":500,"message":"Error retrieving user"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUserService := new(mocks.MockUserService)
			router := chi.NewRouter()
			handler := handlers.NewUserHandler(mockUserService, nil, router)

			mockUserService.On("GetUserByID", mock.AnythingOfType("int32")).Return(tc.mockReturnUser, tc.mockReturnErr)

			rr := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/user/1", nil)

			router.Get("/api/user/{id}", handler.GetUserByIDHandler)
			router.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			assert.Equal(t, tc.expectedBody, strings.TrimSpace(rr.Body.String()))
		})
	}
}

func TestCreateUserHandler(t *testing.T) {
	tests := []struct {
		name            string
		requestBody     *domain.CreateUserRequest
		mockUserService func() *mocks.MockUserService
		expectedStatus  int
		expectedBody    string
	}{
		{
			name: "successful request",
			requestBody: &domain.CreateUserRequest{
				FirstName:       "Kemal",
				LastName:        "Atdayew",
				PhoneNumber:     "+99362008971",
				Gender:          "Male",
				DateOfBirth:     dateOfBirth,
				Location:        "Ashgabat",
				Email:           "atdayewkemal@gmail.com",
				ProfilePhotoURL: "https://example.com/profile.jpg",
			},
			mockUserService: func() *mocks.MockUserService {
				userService := &mocks.MockUserService{}
				userService.On("CreateUser", mock.Anything).Return(&domain.CreateUserResponse{
					ID:              1,
					FirstName:       "Kemal",
					LastName:        "Atdayew",
					PhoneNumber:     "+99362008971",
					Gender:          "Male",
					DateOfBirth:     dateOfBirth,
					Location:        "Ashgabat",
					Email:           "atdayewkemal@gmail.com",
					ProfilePhotoURL: "https://example.com/profile.jpg"}, nil)
				return userService
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":1,"first_name":"Kemal","last_name":"Atdayew","phone_number":"+99362008971","blocked":false,"gender":"Male","registration_date":"0001-01-01T00:00:00Z","date_of_birth":"2000-01-01T00:00:00Z","location":"Ashgabat","email":"atdayewkemal@gmail.com","profile_photo_url":"https://example.com/profile.jpg"}`,
		},
		{
			name: "phone number already in use",
			requestBody: &domain.CreateUserRequest{
				FirstName:       "Kemal",
				LastName:        "Atdayew",
				PhoneNumber:     "+99362008971",
				Gender:          "Male",
				DateOfBirth:     dateOfBirth,
				Location:        "Ashgabat",
				Email:           "atdayewkemal@gmail.com",
				ProfilePhotoURL: "https://example.com/profile.jpg",
			},
			mockUserService: func() *mocks.MockUserService {
				userService := &mocks.MockUserService{}
				userService.On("CreateUser", mock.Anything).Return(&domain.CreateUserResponse{}, domain.ErrPhoneNumberInUse)
				return userService
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   `{"status":409,"message":"Phone number already in use"}`,
		},
		{
			name: "email already in use",
			requestBody: &domain.CreateUserRequest{
				FirstName:       "Kemal",
				LastName:        "Atdayew",
				PhoneNumber:     "+99362008971",
				Gender:          "Male",
				DateOfBirth:     dateOfBirth,
				Location:        "Ashgabat",
				Email:           "atdayewkemal@gmail.com",
				ProfilePhotoURL: "https://example.com/profile.jpg",
			},
			mockUserService: func() *mocks.MockUserService {
				userService := &mocks.MockUserService{}
				userService.On("CreateUser", mock.Anything).Return(&domain.CreateUserResponse{}, domain.ErrEmailInUse)
				return userService
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   `{"status":409,"message":"Email already in use"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/api/user", bytes.NewBuffer(body))
			rr := httptest.NewRecorder()
			router := chi.NewRouter()
			handler := handlers.NewUserHandler(tt.mockUserService(), nil, router)
			router.Post("/api/user", handler.CreateUserHandler)
			router.ServeHTTP(rr, req)
			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedBody, strings.TrimSpace(rr.Body.String()))
		})
	}
}

func TestUpdateUserHandler(t *testing.T) {
	tests := []struct {
		name            string
		requestBody     *domain.UpdateUserRequest
		mockUserService func() *mocks.MockUserService
		expectedStatus  int
		expectedBody    string
	}{
		{
			name: "successful request",
			requestBody: &domain.UpdateUserRequest{
				FirstName:       "Kemal",
				LastName:        "Atdayew",
				Gender:          "Male",
				DateOfBirth:     dateOfBirth,
				Location:        "Ashgabat",
				Email:           "atdayewkemal@gmail.com",
				ProfilePhotoURL: "https://example.com/profile.jpg",
			},
			mockUserService: func() *mocks.MockUserService {
				userService := &mocks.MockUserService{}
				userService.On("UpdateUser", mock.Anything, mock.Anything).Return(&domain.UpdateUserResponse{
					ID:              1,
					FirstName:       "Kemal",
					LastName:        "Atdayew",
					PhoneNumber:     "+99362008971",
					Gender:          "Male",
					DateOfBirth:     dateOfBirth,
					Location:        "Ashgabat",
					Email:           "atdayewkemal@gmail.com",
					ProfilePhotoURL: "https://example.com/profile.jpg"}, nil)
				return userService
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"first_name":"Kemal","last_name":"Atdayew","phone_number":"+99362008971","blocked":false,"gender":"Male","registration_date":"0001-01-01T00:00:00Z","date_of_birth":"2000-01-01T00:00:00Z","location":"Ashgabat","email":"atdayewkemal@gmail.com","profile_photo_url":"https://example.com/profile.jpg"}`,
		},
		{
			name: "user not found",
			requestBody: &domain.UpdateUserRequest{
				FirstName:       "Kemal",
				LastName:        "Atdayew",
				Gender:          "Male",
				DateOfBirth:     dateOfBirth,
				Location:        "Ashgabat",
				Email:           "atdayewkemal@gmail.com",
				ProfilePhotoURL: "https://example.com/profile.jpg",
			},
			mockUserService: func() *mocks.MockUserService {
				userService := &mocks.MockUserService{}
				userService.On("UpdateUser", mock.Anything, mock.Anything).Return(&domain.UpdateUserResponse{}, domain.ErrUserNotFound)
				return userService
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"status":404,"message":"User not found"}`,
		},
		{
			name: "email already in use",
			requestBody: &domain.UpdateUserRequest{
				FirstName:       "Kemal",
				LastName:        "Atdayew",
				Gender:          "Male",
				DateOfBirth:     dateOfBirth,
				Location:        "Ashgabat",
				Email:           "atdayewkemal@gmail.com",
				ProfilePhotoURL: "https://example.com/profile.jpg",
			},
			mockUserService: func() *mocks.MockUserService {
				userService := &mocks.MockUserService{}
				userService.On("UpdateUser", mock.Anything, mock.Anything).Return(&domain.UpdateUserResponse{}, domain.ErrEmailInUse)
				return userService
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   `{"status":409,"message":"Email already in use"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			requestBody, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest("PUT", "/api/user/1", bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			router := chi.NewRouter()

			mockUserService := tt.mockUserService()

			handler := &handlers.UserHandler{
				UserService: mockUserService,
			}

			router.Put("/api/user/{id}", handler.UpdateUserHandler)

			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if strings.TrimSpace(rr.Body.String()) != strings.TrimSpace(tt.expectedBody) {
				t.Errorf("handler returned unexpected body: got %v want %v", []byte(strings.TrimSpace(rr.Body.String())), []byte(strings.TrimSpace(tt.expectedBody)))
			}
		})
	}
}

func TestDeleteUserHandler(t *testing.T) {
	tests := []struct {
		name            string
		id              int
		mockUserService func() *mocks.MockUserService
		expectedStatus  int
		expectedBody    string
	}{
		{
			name: "successful request",
			id:   1,
			mockUserService: func() *mocks.MockUserService {
				userService := &mocks.MockUserService{}
				userService.On("DeleteUser", mock.Anything).Return(nil)
				return userService
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":200,"message":"User deleted successfully"}`,
		},
		{
			name: "user not found",
			id:   1,
			mockUserService: func() *mocks.MockUserService {
				userService := &mocks.MockUserService{}
				userService.On("DeleteUser", mock.Anything).Return(domain.ErrUserNotFound)
				return userService
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"status":404,"message":"User not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/user/%d", tt.id), nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			router := chi.NewRouter()

			mockUserService := tt.mockUserService()

			handler := &handlers.UserHandler{
				UserService: mockUserService,
			}

			router.Delete("/api/user/{id}", handler.DeleteUserHandler)

			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if strings.TrimSpace(rr.Body.String()) != strings.TrimSpace(tt.expectedBody) {
				t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestBlockUserHandler(t *testing.T) {
	tests := []struct {
		name            string
		id              int
		mockUserService func() *mocks.MockUserService
		expectedStatus  int
		expectedBody    string
	}{
		{
			name: "successful request",
			id:   1,
			mockUserService: func() *mocks.MockUserService {
				userService := &mocks.MockUserService{}
				userService.On("BlockUser", mock.Anything).Return(nil)
				return userService
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":200,"message":"User blocked successfully"}`,
		},
		{
			name: "user not found",
			id:   1,
			mockUserService: func() *mocks.MockUserService {
				userService := &mocks.MockUserService{}
				userService.On("BlockUser", mock.Anything).Return(domain.ErrUserNotFound)
				return userService
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"status":404,"message":"User not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("PUT", fmt.Sprintf("/users/%d/block", tt.id), nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			router := chi.NewRouter()

			mockUserService := tt.mockUserService()

			handler := &handlers.UserHandler{
				UserService: mockUserService,
			}

			router.Put("/users/{id}/block", handler.BlockUserHandler)

			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if strings.TrimSpace(rr.Body.String()) != strings.TrimSpace(tt.expectedBody) {
				t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestUnblockUserHandler(t *testing.T) {
	tests := []struct {
		name            string
		id              int
		mockUserService func() *mocks.MockUserService
		expectedStatus  int
		expectedBody    string
	}{
		{
			name: "successful request",
			id:   1,
			mockUserService: func() *mocks.MockUserService {
				userService := &mocks.MockUserService{}
				userService.On("UnblockUser", mock.Anything).Return(nil)
				return userService
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":200,"message":"User unblocked successfully"}`,
		},
		{
			name: "user not found",
			id:   1,
			mockUserService: func() *mocks.MockUserService {
				userService := &mocks.MockUserService{}
				userService.On("UnblockUser", mock.Anything).Return(domain.ErrUserNotFound)
				return userService
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"status":404,"message":"User not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("PUT", fmt.Sprintf("/api/user/%d/unblock", tt.id), nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			router := chi.NewRouter()

			mockUserService := tt.mockUserService()

			handlers := &handlers.UserHandler{
				UserService: mockUserService,
			}

			router.Put("/api/user/{id}/unblock", handlers.UnblockUserHandler)

			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if strings.TrimSpace(rr.Body.String()) != strings.TrimSpace(tt.expectedBody) {
				t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), tt.expectedBody)
			}
		})
	}
}
