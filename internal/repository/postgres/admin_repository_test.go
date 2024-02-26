package repository_test

import (
	"admin-panel/internal/domain"
	mocks "admin-panel/internal/mocks/repository"
	repository "admin-panel/internal/repository/postgres"
	errors "admin-panel/pkg/lib/errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetAllAdmins(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.NewPostgresAdminRepository(db)

	testCases := []struct {
		name           string
		page           int
		pageSize       int
		mockAdmins     []domain.GetAdminResponse
		expectedLength int
	}{
		{
			name:     "Success - Admins exist",
			page:     1,
			pageSize: 10,
			mockAdmins: []domain.GetAdminResponse{
				{ID: 1, Username: "admin1", Role: "admin"},
				{ID: 2, Username: "admin2", Role: "admin"},
			},
			expectedLength: 2,
		},
		{
			name:           "Failure - No admins exist",
			page:           1,
			pageSize:       10,
			mockAdmins:     []domain.GetAdminResponse{},
			expectedLength: 0,
		},
		{
			name:           "Failure - Requested page number exceeds total pages",
			page:           3,
			pageSize:       10,
			mockAdmins:     []domain.GetAdminResponse{},
			expectedLength: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query := `SELECT id, username, role FROM admins ORDER BY id LIMIT \$1 OFFSET \$2`

			rows := sqlmock.NewRows([]string{"id", "username", "role"})
			for _, admin := range tc.mockAdmins {
				rows.AddRow(admin.ID, admin.Username, admin.Role)
			}
			mock.ExpectPrepare(query)
			mock.ExpectQuery(query).WithArgs(tc.pageSize, (tc.page-1)*tc.pageSize).WillReturnRows(rows)

			admins, _ := repo.GetAllAdmins(tc.page, tc.pageSize)

			assert.Equal(t, tc.expectedLength, len(admins.Admins))
		})
	}
}

func TestGetTotalAdminsCount(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.NewPostgresAdminRepository(db)

	testCases := []struct {
		name                string
		mockTotalAdmins     int
		expectedTotalAdmins int
	}{
		{
			name:                "Database has admins",
			mockTotalAdmins:     10,
			expectedTotalAdmins: 10,
		},
		{
			name:                "Database is empty",
			mockTotalAdmins:     0,
			expectedTotalAdmins: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query := "SELECT COUNT\\(\\*\\) FROM admins"
			rows := sqlmock.NewRows([]string{"count"}).AddRow(tc.mockTotalAdmins)
			mock.ExpectQuery(query).WillReturnRows(rows)

			totalAdmins, err := repo.GetTotalAdminsCount()
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedTotalAdmins, totalAdmins)
		})
	}
}

func TestGetAdminByID(t *testing.T) {
	testCases := []struct {
		name            string
		id              int32
		mockReturnAdmin *domain.GetAdminResponse
		mockReturnErr   error
		expectedErr     error
	}{
		{
			name: "Success",
			id:   1,
			mockReturnAdmin: &domain.GetAdminResponse{
				ID:       1,
				Username: "Kemal",
				Role:     "admin",
			},
			mockReturnErr: nil,
			expectedErr:   nil,
		},
		{
			name:            "Admin Not Found",
			id:              2,
			mockReturnAdmin: nil,
			mockReturnErr:   errors.ErrAdminNotFound,
			expectedErr:     errors.ErrAdminNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockAdminRepository := new(mocks.MockAdminRepository)

			mockAdminRepository.ExpectedCalls = nil

			mockAdminRepository.On("GetAdminByID", tc.id).Return(tc.mockReturnAdmin, tc.mockReturnErr)

			admin, err := mockAdminRepository.GetAdminByID(tc.id)

			if err != nil {
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.Equal(t, tc.mockReturnAdmin, admin)
			}
		})
	}
}

func TestCreateAdmin(t *testing.T) {
	testCases := []struct {
		name            string
		input           *domain.CreateAdminRequest
		mockReturnAdmin *domain.CreateAdminResponse
		mockReturnErr   error
		expectedErr     error
	}{
		{
			name: "Success",
			input: &domain.CreateAdminRequest{
				Username: "testuser",
				Password: "testpass",
				Role:     "user",
			},
			mockReturnAdmin: &domain.CreateAdminResponse{
				ID:       1,
				Username: "testuser",
				Role:     "user",
			},
			mockReturnErr: nil,
			expectedErr:   nil,
		},
		{
			name: "User Already Exists",
			input: &domain.CreateAdminRequest{
				Username: "existinguser",
				Password: "testpass",
				Role:     "user",
			},
			mockReturnAdmin: nil,
			mockReturnErr:   errors.ErrAdminAlreadyExists,
			expectedErr:     errors.ErrAdminAlreadyExists,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockAdminRepository := new(mocks.MockAdminRepository)

			mockAdminRepository.ExpectedCalls = nil

			mockAdminRepository.On("CreateAdmin", tc.input).Return(tc.mockReturnAdmin, tc.mockReturnErr)

			user, err := mockAdminRepository.CreateAdmin(tc.input)

			if err != nil {
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.Equal(t, tc.mockReturnAdmin, user)
			}
		})
	}
}

func TestUpdateAdmin(t *testing.T) {
	testCases := []struct {
		name            string
		id              int32
		input           *domain.UpdateAdminRequest
		mockReturnAdmin *domain.UpdateAdminResponse
		mockReturnErr   error
		expectedErr     error
	}{
		{
			name: "Success",
			id:   1,
			input: &domain.UpdateAdminRequest{
				Username: "updateduser",
				Password: "updatedpass",
				Role:     "admin",
			},
			mockReturnAdmin: &domain.UpdateAdminResponse{
				ID:       1,
				Username: "updateduser",
				Role:     "admin",
			},
			mockReturnErr: nil,
			expectedErr:   nil,
		},
		{
			name: "Admin Not Found",
			id:   2,
			input: &domain.UpdateAdminRequest{
				Username: "nonexistentuser",
				Password: "testpass",
				Role:     "admin",
			},
			mockReturnAdmin: nil,
			mockReturnErr:   errors.ErrAdminNotFound,
			expectedErr:     errors.ErrAdminNotFound,
		},
		{
			name: "Internal Server Error",
			id:   3,
			input: &domain.UpdateAdminRequest{
				Username: "servererror",
				Password: "testpass",
				Role:     "admin",
			},
			mockReturnAdmin: nil,
			mockReturnErr:   errors.ErrInternalServerError,
			expectedErr:     errors.ErrInternalServerError,
		},
		{
			name: "Empty Fields",
			id:   4,
			input: &domain.UpdateAdminRequest{
				Username: "",
				Password: "",
				Role:     "",
			},
			mockReturnAdmin: nil,
			mockReturnErr:   errors.ErrFillRequiredFields,
			expectedErr:     errors.ErrFillRequiredFields,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockAdminRepository := new(mocks.MockAdminRepository)

			mockAdminRepository.On("UpdateAdmin", tc.id, tc.input).Return(tc.mockReturnAdmin, tc.mockReturnErr)

			admin, err := mockAdminRepository.UpdateAdmin(tc.id, tc.input)

			if err != nil {
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.Equal(t, tc.mockReturnAdmin, admin)
			}
		})
	}
}

func TestDeleteAdmin(t *testing.T) {
	testCases := []struct {
		name          string
		id            int32
		mockReturnErr error
		expectedErr   error
	}{
		{
			name:          "Success",
			id:            1,
			mockReturnErr: nil,
			expectedErr:   nil,
		},
		{
			name:          "Admin Not Found",
			id:            2,
			mockReturnErr: errors.ErrAdminNotFound,
			expectedErr:   errors.ErrAdminNotFound,
		},
		{
			name:          "Internal Server Error",
			id:            3,
			mockReturnErr: errors.ErrInternalServerError,
			expectedErr:   errors.ErrInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockAdminRepository := new(mocks.MockAdminRepository)

			mockAdminRepository.On("DeleteAdmin", tc.id).Return(tc.mockReturnErr)

			err := mockAdminRepository.DeleteAdmin(tc.id)

			if err != nil {
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestSearchAdmins(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.NewPostgresAdminRepository(db)

	testCases := []struct {
		name           string
		query          string
		page           int
		pageSize       int
		mockAdmins     []domain.GetAdminResponse
		expectedLength int
	}{
		{
			name:     "Success - Admins exist",
			query:    "admin",
			page:     1,
			pageSize: 10,
			mockAdmins: []domain.GetAdminResponse{
				{ID: 1, Username: "admin1", Role: "admin"},
				{ID: 2, Username: "admin2", Role: "admin"},
			},
			expectedLength: 2,
		},
		{
			name:           "Failure - No admins exist",
			query:          "admin",
			page:           1,
			pageSize:       10,
			mockAdmins:     []domain.GetAdminResponse{},
			expectedLength: 0,
		},
		{
			name:           "Failure - Requested page number exceeds total pages",
			query:          "admin",
			page:           3,
			pageSize:       10,
			mockAdmins:     []domain.GetAdminResponse{},
			expectedLength: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			searchQuery := `SELECT id, username, role FROM admins WHERE username ILIKE \$1 OR role ILIKE \$1 ORDER BY id LIMIT \$2 OFFSET \$3`

			rows := sqlmock.NewRows([]string{"id", "username", "role"})
			for _, admin := range tc.mockAdmins {
				rows.AddRow(admin.ID, admin.Username, admin.Role)
			}
			mock.ExpectPrepare(searchQuery)
			mock.ExpectQuery(searchQuery).WithArgs("%"+tc.query+"%", tc.pageSize, (tc.page-1)*tc.pageSize).WillReturnRows(rows)

			admins, _ := repo.SearchAdmins(tc.query, tc.page, tc.pageSize)

			assert.Equal(t, tc.expectedLength, len(admins.Admins))
		})
	}
}
