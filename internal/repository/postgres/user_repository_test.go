package repository_test

import (
	"admin-panel/internal/domain"
	mocks "admin-panel/internal/mocks/repository"
	repository "admin-panel/internal/repository/postgres"
	errors "admin-panel/pkg/lib/errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetAllUsers(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.NewPostgresUserRepository(db)

	testCases := []struct {
		name           string
		page           int
		limit          int
		mockUsers      []domain.CommonUserResponse
		expectedLength int
	}{
		{
			name:  "Success - Users exist",
			page:  1,
			limit: 10,
			mockUsers: []domain.CommonUserResponse{
				{
					ID:               1,
					FirstName:        "John",
					LastName:         "Doe",
					PhoneNumber:      "1234567890",
					Blocked:          false,
					Gender:           "Male",
					RegistrationDate: time.Now(),
					DateOfBirth:      time.Now(),
					Location:         "Location1",
					Email:            "john.doe@example.com",
					ProfilePhotoURL:  "https://example.com/john.jpg",
				},
			},
			expectedLength: 1,
		},
		{
			name:           "Failure - No users exist",
			page:           1,
			limit:          10,
			mockUsers:      []domain.CommonUserResponse{},
			expectedLength: 0,
		},
		{
			name:           "Failure - Requested page number exceeds total pages",
			page:           3,
			limit:          10,
			mockUsers:      []domain.CommonUserResponse{},
			expectedLength: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query := `SELECT id, first_name, last_name, phone_number, blocked, registration_date, gender, date_of_birth, location, email, profile_photo_url FROM users ORDER BY id LIMIT \$1 OFFSET \$2`

			rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "phone_number", "blocked", "registration_date", "gender", "date_of_birth", "location", "email", "profile_photo_url"})
			for _, user := range tc.mockUsers {
				rows.AddRow(user.ID, user.FirstName, user.LastName, user.PhoneNumber, user.Blocked, user.RegistrationDate, user.Gender, user.DateOfBirth, user.Location, user.Email, user.ProfilePhotoURL)
			}
			mock.ExpectPrepare(query)
			mock.ExpectQuery(query).WithArgs(tc.limit, (tc.page-1)*tc.limit).WillReturnRows(rows)

			users, _ := repo.GetAllUsers(tc.page, tc.limit)

			assert.Equal(t, tc.expectedLength, len(users.Users))
		})
	}
}

func TestGetTotalUsersCount(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.NewPostgresUserRepository(db)

	testCases := []struct {
		name               string
		mockTotalUsers     int
		expectedTotalUsers int
	}{
		{
			name:               "Database has users",
			mockTotalUsers:     10,
			expectedTotalUsers: 10,
		},
		{
			name:               "Database is empty",
			mockTotalUsers:     0,
			expectedTotalUsers: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query := "SELECT COUNT\\(\\*\\) FROM users"
			rows := sqlmock.NewRows([]string{"count"}).AddRow(tc.mockTotalUsers)
			mock.ExpectQuery(query).WillReturnRows(rows)

			totalUsers, err := repo.GetTotalUsersCount()
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedTotalUsers, totalUsers)
		})
	}
}

func TestGetUserByID(t *testing.T) {
	testCases := []struct {
		name           string
		id             int32
		mockReturnUser *domain.GetUserResponse
		mockReturnErr  error
		expectedErr    error
	}{
		{
			name: "Success",
			id:   1,
			mockReturnUser: &domain.GetUserResponse{
				ID:              1,
				FirstName:       "Kemal",
				LastName:        "Atdayew",
				PhoneNumber:     "+99362008971",
				Gender:          "Male",
				Location:        "Ashgabat",
				Email:           "atdayewkemal@gmail.com",
				ProfilePhotoURL: "https://example.com/profile.jpg",
			},
			mockReturnErr: nil,
			expectedErr:   nil,
		},
		{
			name:           "User Not Found",
			id:             2,
			mockReturnUser: nil,
			mockReturnErr:  errors.ErrUserNotFound,
			expectedErr:    errors.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUserRepository := new(mocks.MockUserRepository)

			mockUserRepository.ExpectedCalls = nil

			mockUserRepository.On("GetUserByID", tc.id).Return(tc.mockReturnUser, tc.mockReturnErr)

			user, err := mockUserRepository.GetUserByID(tc.id)

			if err != nil {
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.Equal(t, tc.mockReturnUser, user)
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	testCases := []struct {
		name           string
		request        *domain.CreateUserRequest
		mockReturnUser *domain.CreateUserResponse
		mockReturnErr  error
		expectedErr    error
	}{
		{
			name: "Success",
			request: &domain.CreateUserRequest{
				FirstName:       "Kemal",
				LastName:        "Atdayew",
				PhoneNumber:     "+99362008971",
				Gender:          "Male",
				Location:        "Ashgabat",
				Email:           "atdayewkemal@gmail.com",
				ProfilePhotoURL: "https://example.com/profile.jpg",
			},
			mockReturnUser: &domain.CreateUserResponse{
				ID:              1,
				FirstName:       "Kemal",
				LastName:        "Atdayew",
				PhoneNumber:     "+99362008971",
				Gender:          "Male",
				Location:        "Ashgabat",
				Email:           "atdayewkemal@gmail.com",
				ProfilePhotoURL: "https://example.com/profile.jpg",
			},
			mockReturnErr: nil,
			expectedErr:   nil,
		},
		{
			name: "Invalid Phone Number",
			request: &domain.CreateUserRequest{
				FirstName:       "Kemal",
				LastName:        "Atdayew",
				PhoneNumber:     "invalid",
				Gender:          "Male",
				Location:        "Ashgabat",
				Email:           "atdayewkemal@gmail.com",
				ProfilePhotoURL: "https://example.com/profile.jpg",
			},
			mockReturnUser: nil,
			mockReturnErr:  errors.ErrInvalidPhoneNumber,
			expectedErr:    errors.ErrInvalidPhoneNumber,
		},
		{
			name: "Email In Use",
			request: &domain.CreateUserRequest{
				FirstName:       "Kemal",
				LastName:        "Atdayew",
				PhoneNumber:     "+99362008971",
				Gender:          "Male",
				Location:        "Ashgabat",
				Email:           "atdayewkemal@gmail.com",
				ProfilePhotoURL: "https://example.com/profile.jpg",
			},
			mockReturnUser: nil,
			mockReturnErr:  errors.ErrEmailInUse,
			expectedErr:    errors.ErrEmailInUse,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUserRepository := new(mocks.MockUserRepository)

			mockUserRepository.On("CreateUser", tc.request).Return(tc.mockReturnUser, tc.mockReturnErr)

			user, err := mockUserRepository.CreateUser(tc.request)

			if err != nil {
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.Equal(t, tc.mockReturnUser, user)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	testCases := []struct {
		name           string
		id             int32
		request        *domain.UpdateUserRequest
		mockReturnUser *domain.UpdateUserResponse
		mockReturnErr  error
		expectedErr    error
	}{
		{
			name: "Success",
			id:   1,
			request: &domain.UpdateUserRequest{
				FirstName:       "Kemal",
				LastName:        "Atdayew",
				Gender:          "Male",
				Location:        "Ashgabat",
				Email:           "atdayewkemal@gmail.com",
				ProfilePhotoURL: "https://example.com/profile.jpg",
			},
			mockReturnUser: &domain.UpdateUserResponse{
				ID:              1,
				FirstName:       "Kemal",
				LastName:        "Atdayew",
				PhoneNumber:     "+99362008971",
				Gender:          "Male",
				Location:        "Ashgabat",
				Email:           "atdayewkemal@gmail.com",
				ProfilePhotoURL: "https://example.com/profile.jpg",
			},
			mockReturnErr: nil,
			expectedErr:   nil,
		},
		{
			name: "User Not Found",
			id:   2,
			request: &domain.UpdateUserRequest{
				FirstName:       "Kemal",
				LastName:        "Atdayew",
				Gender:          "Male",
				Location:        "Ashgabat",
				Email:           "atdayewkemal@gmail.com",
				ProfilePhotoURL: "https://example.com/profile.jpg",
			},
			mockReturnUser: nil,
			mockReturnErr:  errors.ErrUserNotFound,
			expectedErr:    errors.ErrUserNotFound,
		},
		{
			name: "Email In Use",
			id:   3,
			request: &domain.UpdateUserRequest{
				FirstName:       "Kemal",
				LastName:        "Atdayew",
				Gender:          "Male",
				Location:        "Ashgabat",
				Email:           "atdayewkemal@gmail.com",
				ProfilePhotoURL: "https://example.com/profile.jpg",
			},
			mockReturnUser: nil,
			mockReturnErr:  errors.ErrEmailInUse,
			expectedErr:    errors.ErrEmailInUse,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUserRepository := new(mocks.MockUserRepository)
			mockUserRepository.On("UpdateUser", tc.id, tc.request).Return(tc.mockReturnUser, tc.mockReturnErr)
			user, err := mockUserRepository.UpdateUser(tc.id, tc.request)
			if err != nil {
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.Equal(t, tc.mockReturnUser, user)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
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
			name:          "User Not Found",
			id:            2,
			mockReturnErr: errors.ErrUserNotFound,
			expectedErr:   errors.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUserRepository := new(mocks.MockUserRepository)

			mockUserRepository.On("DeleteUser", tc.id).Return(tc.mockReturnErr)

			err := mockUserRepository.DeleteUser(tc.id)

			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestBlockUser(t *testing.T) {
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
			name:          "User Not Found",
			id:            2,
			mockReturnErr: errors.ErrUserNotFound,
			expectedErr:   errors.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUserRepository := new(mocks.MockUserRepository)

			mockUserRepository.On("BlockUser", tc.id).Return(tc.mockReturnErr)

			err := mockUserRepository.BlockUser(tc.id)

			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestUnblockUser(t *testing.T) {
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
			name:          "User Not Found",
			id:            2,
			mockReturnErr: errors.ErrUserNotFound,
			expectedErr:   errors.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUserRepository := new(mocks.MockUserRepository)

			mockUserRepository.On("UnblockUser", tc.id).Return(tc.mockReturnErr)

			err := mockUserRepository.UnblockUser(tc.id)

			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestSearchUsers(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.NewPostgresUserRepository(db)

	testCases := []struct {
		name              string
		query             string
		page              int
		pageSize          int
		mockUsers         []domain.CommonUserResponse
		expectedUserCount int
	}{
		{
			name:     "Success - Users found",
			query:    "John",
			page:     1,
			pageSize: 10,
			mockUsers: []domain.CommonUserResponse{
				{
					ID:               1,
					FirstName:        "John",
					LastName:         "Doe",
					PhoneNumber:      "1234567890",
					Blocked:          false,
					Gender:           "Male",
					RegistrationDate: time.Now(),
					DateOfBirth:      time.Now(),
					Location:         "Location1",
					Email:            "john.doe@example.com",
					ProfilePhotoURL:  "https://example.com/john.jpg",
				},
			},
			expectedUserCount: 1,
		},
		{
			name:              "Failure - No users found",
			query:             "Nonexistent",
			page:              1,
			pageSize:          10,
			mockUsers:         []domain.CommonUserResponse{},
			expectedUserCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			searchQuery := `
				SELECT id, first_name, last_name, phone_number, blocked,
				registration_date, gender, date_of_birth, location,
				email, profile_photo_url
				FROM users
				WHERE first_name ILIKE \$1 OR last_name ILIKE \$1 OR phone_number ILIKE \$1 OR email ILIKE \$1
				ORDER BY id
				LIMIT \$2 OFFSET \$3
			`
			rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "phone_number", "blocked", "registration_date", "gender", "date_of_birth", "location", "email", "profile_photo_url"})
			for _, user := range tc.mockUsers {
				rows.AddRow(user.ID, user.FirstName, user.LastName, user.PhoneNumber, user.Blocked, user.RegistrationDate, user.Gender, user.DateOfBirth, user.Location, user.Email, user.ProfilePhotoURL)
			}
			mock.ExpectPrepare(searchQuery)
			mock.ExpectQuery(searchQuery).WithArgs("%"+tc.query+"%", tc.pageSize, (tc.page-1)*tc.pageSize).WillReturnRows(rows)

			users, _ := repo.SearchUsers(tc.query, tc.page, tc.pageSize)

			assert.Equal(t, tc.expectedUserCount, len(users.Users))
		})
	}
}
