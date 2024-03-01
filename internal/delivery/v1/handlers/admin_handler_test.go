package handlers_test

import (
	"admin-panel/internal/delivery/v1/handlers"
	"admin-panel/internal/domain"
	mocks "admin-panel/internal/mocks/service"
	"admin-panel/pkg/lib/errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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
