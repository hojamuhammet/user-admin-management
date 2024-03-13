package handlers_test

import (
	"admin-panel/internal/delivery/v1/handlers"
	"admin-panel/internal/domain"
	mocks "admin-panel/internal/mocks/service"
	"admin-panel/pkg/lib/errors"
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllAdminsHandler(t *testing.T) {
	testCases := []struct {
		name            string
		page            string
		pageSize        string
		mockReturn      *domain.AdminsList
		mockReturnErr   error
		mockCountReturn int
		mockCountErr    error
		expectedStatus  int
	}{
		{
			name:            "Success",
			page:            "1",
			pageSize:        "8",
			mockReturn:      &domain.AdminsList{},
			mockReturnErr:   nil,
			mockCountReturn: 10,
			mockCountErr:    nil,
			expectedStatus:  http.StatusOK,
		},
		{
			name:            "Invalid Page Number",
			page:            "0",
			pageSize:        "8",
			mockReturn:      nil,
			mockReturnErr:   nil,
			mockCountReturn: 10,
			mockCountErr:    nil,
			expectedStatus:  http.StatusOK, // Default page is 1
		},
		{
			name:            "Invalid Page Size",
			page:            "1",
			pageSize:        "0",
			mockReturn:      nil,
			mockReturnErr:   nil,
			mockCountReturn: 10,
			mockCountErr:    nil,
			expectedStatus:  http.StatusOK, // Default page size is 8
		},
		{
			name:            "Error Getting Admins",
			page:            "1",
			pageSize:        "8",
			mockReturn:      nil,
			mockReturnErr:   errors.ErrGettingAdmins,
			mockCountReturn: 10,
			mockCountErr:    nil,
			expectedStatus:  http.StatusInternalServerError,
		},
		{
			name:            "Error Getting Total Admins Count",
			page:            "1",
			pageSize:        "8",
			mockReturn:      &domain.AdminsList{},
			mockReturnErr:   nil,
			mockCountReturn: 0,
			mockCountErr:    errors.ErrGettingTotalAdminCount,
			expectedStatus:  http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockAdminService := new(mocks.MockAdminService)
			router := chi.NewRouter()
			handler := handlers.AdminHandler{
				AdminService: mockAdminService,
			}

			mockAdminService.On("GetAllAdmins", mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(tc.mockReturn, tc.mockReturnErr)
			mockAdminService.On("GetTotalAdminsCount").Return(tc.mockCountReturn, tc.mockCountErr)

			rr := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/api/admin?page=%s&pageSize=%s", tc.page, tc.pageSize), nil)

			router.Get("/api/admin", handler.GetAllAdminsHandler)
			router.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
		})
	}
}

func TestGetAdminByID(t *testing.T) {
	testCases := []struct {
		name           string
		mockReturn     *domain.GetAdminResponse
		mockReturnErr  error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",
			mockReturn: &domain.GetAdminResponse{
				ID:       1,
				Username: "Admin1",
				Role:     "Admin",
			},
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"username":"Admin1","role":"Admin"}`,
		},
		{
			name:           "Admin Not Found",
			mockReturn:     nil,
			mockReturnErr:  errors.ErrAdminNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"status":404,"message":"Admin not found"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockAdminService := new(mocks.MockAdminService)
			router := chi.NewRouter()
			handler := handlers.AdminHandler{
				AdminService: mockAdminService,
			}

			mockAdminService.On("GetAdminByID", mock.AnythingOfType("int32")).Return(tc.mockReturn, tc.mockReturnErr)

			rr := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/admin/1", nil)

			router.Get("/api/admin/{id}", handler.GetAdminByID)
			router.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			assert.Equal(t, tc.expectedBody, strings.TrimSpace(rr.Body.String()))
		})
	}
}

func TestCreateAdminHandler(t *testing.T) {
	testCases := []struct {
		name           string
		requestBody    string
		mockReturn     *domain.CreateAdminResponse
		mockReturnErr  error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "Success",
			requestBody: `{"username":"Admin1","password":"password","role":"Admin"}`,
			mockReturn: &domain.CreateAdminResponse{
				ID:       1,
				Username: "Admin1",
				Role:     "Admin",
			},
			mockReturnErr:  nil,
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":1,"username":"Admin1","role":"Admin"}`,
		},
		{
			name:           "Admin Already Exists",
			requestBody:    `{"username":"Admin1","password":"password","role":"Admin"}`,
			mockReturn:     nil,
			mockReturnErr:  errors.ErrAdminAlreadyExists,
			expectedStatus: http.StatusConflict,
			expectedBody:   `{"status":409,"message":"Admin with the same username already exists"}`,
		},
		{
			name:           "Missing Fields",
			requestBody:    `{"username":"","password":"password","role":"Admin"}`,
			mockReturn:     nil,
			mockReturnErr:  nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"status":400,"message":"Username, password, and role are required fields"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockAdminService := new(mocks.MockAdminService)
			router := chi.NewRouter()
			handler := handlers.AdminHandler{
				AdminService: mockAdminService,
			}

			mockAdminService.On("CreateAdmin", mock.AnythingOfType("*domain.CreateAdminRequest")).Return(tc.mockReturn, tc.mockReturnErr)

			rr := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/admin", bytes.NewBuffer([]byte(tc.requestBody)))

			router.Post("/api/admin", handler.CreateAdminHandler)
			router.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			assert.Equal(t, tc.expectedBody, strings.TrimSpace(rr.Body.String()))
		})
	}
}

func TestUpdateAdminHandler(t *testing.T) {
	testCases := []struct {
		name           string
		requestBody    string
		mockReturn     *domain.UpdateAdminResponse
		mockReturnErr  error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "Success",
			requestBody: `{"username":"Admin1","password":"password","role":"Admin"}`,
			mockReturn: &domain.UpdateAdminResponse{
				ID:       1,
				Username: "Admin1",
				Role:     "Admin",
			},
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"username":"Admin1","role":"Admin"}`,
		},
		{
			name:           "Admin Not Found",
			requestBody:    `{"username":"Admin1","password":"password","role":"Admin"}`,
			mockReturn:     nil,
			mockReturnErr:  errors.ErrAdminNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"status":404,"message":"Admin not found"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockAdminService := new(mocks.MockAdminService)
			router := chi.NewRouter()
			handler := handlers.AdminHandler{
				AdminService: mockAdminService,
			}

			mockAdminService.On("UpdateAdmin", mock.AnythingOfType("int32"), mock.AnythingOfType("*domain.UpdateAdminRequest")).Return(tc.mockReturn, tc.mockReturnErr)

			rr := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", "/api/admin/1", bytes.NewBuffer([]byte(tc.requestBody)))

			router.Put("/api/admin/{id}", handler.UpdateAdminHandler)
			router.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			assert.Equal(t, tc.expectedBody, strings.TrimSpace(rr.Body.String()))
		})
	}
}

func TestDeleteAdminHandler(t *testing.T) {
	testCases := []struct {
		name           string
		mockReturnErr  error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Success",
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":200,"message":"Admin deleted successfully"}`,
		},
		{
			name:           "Admin Not Found",
			mockReturnErr:  errors.ErrAdminNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"status":404,"message":"Admin not found"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockAdminService := new(mocks.MockAdminService)
			router := chi.NewRouter()
			handler := handlers.AdminHandler{
				AdminService: mockAdminService,
			}

			mockAdminService.On("DeleteAdmin", mock.AnythingOfType("int32")).Return(tc.mockReturnErr)

			rr := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/api/admin/1", nil)

			router.Delete("/api/admin/{id}", handler.DeleteAdminHandler)
			router.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			assert.Equal(t, tc.expectedBody, strings.TrimSpace(rr.Body.String()))
		})
	}
}

func TestSearchAdminsHandler(t *testing.T) {
	testCases := []struct {
		name            string
		query           string
		page            string
		pageSize        string
		mockReturn      *domain.AdminsList
		mockReturnErr   error
		mockCountReturn int
		mockCountErr    error
		expectedStatus  int
	}{
		{
			name:            "Success",
			query:           "Admin1",
			page:            "1",
			pageSize:        "8",
			mockReturn:      &domain.AdminsList{},
			mockReturnErr:   nil,
			mockCountReturn: 10,
			mockCountErr:    nil,
			expectedStatus:  http.StatusOK,
		},
		{
			name:            "Invalid Page Number",
			query:           "Admin1",
			page:            "0",
			pageSize:        "8",
			mockReturn:      nil,
			mockReturnErr:   nil,
			mockCountReturn: 10,
			mockCountErr:    nil,
			expectedStatus:  http.StatusOK, // Default page is 1
		},
		{
			name:            "Invalid Page Size",
			query:           "Admin1",
			page:            "1",
			pageSize:        "0",
			mockReturn:      nil,
			mockReturnErr:   nil,
			mockCountReturn: 10,
			mockCountErr:    nil,
			expectedStatus:  http.StatusOK, // Default page size is 8
		},
		{
			name:            "Error Getting Total Admins Count",
			query:           "Admin1",
			page:            "1",
			pageSize:        "8",
			mockReturn:      &domain.AdminsList{},
			mockReturnErr:   nil,
			mockCountReturn: 0,
			mockCountErr:    errors.ErrGettingTotalAdminCount,
			expectedStatus:  http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockAdminService := new(mocks.MockAdminService)
			router := chi.NewRouter()
			handler := handlers.AdminHandler{
				AdminService: mockAdminService,
			}

			mockAdminService.On("SearchAdmins", tc.query, mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(tc.mockReturn, tc.mockReturnErr)
			mockAdminService.On("GetTotalAdminsCount").Return(tc.mockCountReturn, tc.mockCountErr)

			rr := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/api/admin/search?query=%s&page=%s&pageSize=%s", tc.query, tc.page, tc.pageSize), nil)

			router.Get("/api/admin/search", handler.SearchAdminsHandler)
			router.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
		})
	}
}
