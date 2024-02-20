package repository_test

import (
	"admin-panel/internal/domain"
	"admin-panel/internal/domain/mocks"
	"admin-panel/internal/service"
	"errors"
	"testing"
	"time"
)

func TestGetUserByID(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock()

	registrationDate, _ := time.Parse(time.RFC3339, "2002-06-08T00:00:00Z")
	dateOfBirth, _ := time.Parse(time.RFC3339, "2024-02-16T22:48:15Z")

	expectedUser := &domain.GetUserResponse{
		ID:               66,
		FirstName:        "",
		LastName:         "",
		PhoneNumber:      "+99376065810",
		Blocked:          false,
		Gender:           "",
		RegistrationDate: registrationDate,
		DateOfBirth:      dateOfBirth,
		Location:         "",
		Email:            "",
		ProfilePhotoURL:  "",
	}
	mockRepo.GetUserByIDFunc = func(id int32) (*domain.GetUserResponse, error) {
		if id == expectedUser.ID {
			return expectedUser, nil
		}
		return nil, errors.New("user not found")
	}

	userService := service.NewUserService(mockRepo)

	user, err := userService.GetUserByID(66)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if user == nil {
		t.Error("Expected user, got nil")
	} else if user.ID != expectedUser.ID {
		t.Errorf("Expected user ID to be %d, got %d", expectedUser.ID, user.ID)
	}
}

func TestCreateUser(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock()

	testCases := []struct {
		Name             string
		CreateUserFields map[string]interface{}
		ExpectedResponse *domain.CreateUserResponse
		ExpectedError    error
	}{
		{
			Name: "CreateUser with only first name",
			CreateUserFields: map[string]interface{}{
				"FirstName":   "John",
				"PhoneNumber": "+99362008971",
			},
			ExpectedResponse: &domain.CreateUserResponse{ID: 1},
			ExpectedError:    nil,
		},
		{
			Name: "CreateUser with only last name",
			CreateUserFields: map[string]interface{}{
				"LastName":    "Doe",
				"PhoneNumber": "+99362008971",
			},
			ExpectedResponse: &domain.CreateUserResponse{ID: 1},
			ExpectedError:    nil,
		},
		{
			Name: "CreateUser with only gender",
			CreateUserFields: map[string]interface{}{
				"Gender":      "Male",
				"PhoneNumber": "+99362008971",
			},
			ExpectedResponse: &domain.CreateUserResponse{ID: 1},
			ExpectedError:    nil,
		},
		{
			Name: "CreateUser with only date of birth",
			CreateUserFields: map[string]interface{}{
				"DateOfBirth": time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
				"PhoneNumber": "+99362008971",
			},
			ExpectedResponse: &domain.CreateUserResponse{ID: 1},
			ExpectedError:    nil,
		},
		{
			Name: "CreateUser with only location",
			CreateUserFields: map[string]interface{}{
				"Location":    "Ashgabat",
				"PhoneNumber": "+99362008971",
			},
			ExpectedResponse: &domain.CreateUserResponse{ID: 1},
			ExpectedError:    nil,
		},
		{
			Name: "CreateUser with only email",
			CreateUserFields: map[string]interface{}{
				"Email":       "asdasdasd@gmail.com",
				"PhoneNumber": "+99362008971",
			},
			ExpectedResponse: &domain.CreateUserResponse{ID: 1},
			ExpectedError:    nil,
		},
		{
			Name: "CreateUser with only profile photo url",
			CreateUserFields: map[string]interface{}{
				"ProfilePhotoURL": "my_photo.jpeg",
				"PhoneNumber":     "+99362008971",
			},
			ExpectedResponse: &domain.CreateUserResponse{ID: 1},
			ExpectedError:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			mockRepo.CreateUserFunc = func(request *domain.CreateUserRequest) (*domain.CreateUserResponse, error) {
				for field, value := range tc.CreateUserFields {
					if val, ok := value.(string); ok {
						switch field {
						case "FirstName":
							request.FirstName = val
						case "LastName":
							request.LastName = val
						case "Gender":
							request.Gender = val
						case "Location":
							request.Location = val
						case "Email":
							request.Email = val
						case "ProfilePhotoURL":
							request.ProfilePhotoURL = val
						}
					}
					if field == "DateOfBirth" {
						if val, ok := value.(time.Time); ok {
							request.DateOfBirth = val
						}
					}
				}
				return tc.ExpectedResponse, nil
			}

			userService := service.NewUserService(mockRepo)

			createUserRequest := &domain.CreateUserRequest{
				PhoneNumber: tc.CreateUserFields["PhoneNumber"].(string),
			}

			response, err := userService.CreateUser(createUserRequest)

			if err != nil {
				if tc.ExpectedError == nil || err.Error() != tc.ExpectedError.Error() {
					t.Errorf("Unexpected error: %v", err)
				}
			} else {
				if response == nil {
					t.Error("Expected response, got nil")
				}
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock()

	testCases := []struct {
		Name             string
		UpdateUserFields map[string]interface{}
		ExpectedResponse *domain.UpdateUserResponse
		ExpectedError    error
	}{
		{
			Name: "UpdateUser with only first name",
			UpdateUserFields: map[string]interface{}{
				"FirstName": "John",
			},
			ExpectedResponse: &domain.UpdateUserResponse{ID: 1},
			ExpectedError:    nil,
		},
		{
			Name: "UpdateUser with only last name",
			UpdateUserFields: map[string]interface{}{
				"LastName": "Doe",
			},
			ExpectedResponse: &domain.UpdateUserResponse{ID: 1},
			ExpectedError:    nil,
		},
		{
			Name: "UpdateUser with only gender",
			UpdateUserFields: map[string]interface{}{
				"Gender": "Male",
			},
			ExpectedResponse: &domain.UpdateUserResponse{ID: 1},
			ExpectedError:    nil,
		},
		{
			Name: "UpdateUser with only date of birth",
			UpdateUserFields: map[string]interface{}{
				"DateOfBirth": time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			ExpectedResponse: &domain.UpdateUserResponse{ID: 1},
			ExpectedError:    nil,
		},
		{
			Name: "UpdateUser with only location",
			UpdateUserFields: map[string]interface{}{
				"Location": "Ashgabat",
			},
			ExpectedResponse: &domain.UpdateUserResponse{ID: 1},
			ExpectedError:    nil,
		},
		{
			Name: "UpdateUser with only email",
			UpdateUserFields: map[string]interface{}{
				"Email": "asdasdasd@gmail.com",
			},
			ExpectedResponse: &domain.UpdateUserResponse{ID: 1},
			ExpectedError:    nil,
		},
		{
			Name: "UpdateUser with only profile photo URL",
			UpdateUserFields: map[string]interface{}{
				"ProfilePhotoURL": "my_photo.jpeg",
			},
			ExpectedResponse: &domain.UpdateUserResponse{ID: 1},
			ExpectedError:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			mockRepo.UpdateUserFunc = func(userID int32, request *domain.UpdateUserRequest) (*domain.UpdateUserResponse, error) {
				for field, value := range tc.UpdateUserFields {
					if val, ok := value.(string); ok {
						switch field {
						case "FirstName":
							request.FirstName = val
						case "LastName":
							request.LastName = val
						case "Gender":
							request.Gender = val
						case "Location":
							request.Location = val
						case "Email":
							request.Email = val
						case "ProfilePhotoURL":
							request.ProfilePhotoURL = val
						}
					}
					if field == "DateOfBirth" {
						if val, ok := value.(time.Time); ok {
							request.DateOfBirth = val
						}
					}
				}
				return tc.ExpectedResponse, nil
			}

			userService := service.NewUserService(mockRepo)

			updateUserRequest := &domain.UpdateUserRequest{}
			for field, value := range tc.UpdateUserFields {
				if val, ok := value.(string); ok {
					switch field {
					case "FirstName":
						updateUserRequest.FirstName = val
					case "LastName":
						updateUserRequest.LastName = val
					case "Gender":
						updateUserRequest.Gender = val
					case "Location":
						updateUserRequest.Location = val
					case "Email":
						updateUserRequest.Email = val
					case "ProfilePhotoURL":
						updateUserRequest.ProfilePhotoURL = val
					}
				}
			}

			userID := int32(1)

			response, err := userService.UpdateUser(userID, updateUserRequest)

			if err != nil {
				if tc.ExpectedError == nil || err.Error() != tc.ExpectedError.Error() {
					t.Errorf("Test case '%s' failed: unexpected error: %v", tc.Name, err)
				}
			} else {
				if response == nil {
					t.Errorf("Test case '%s' failed: expected non-nil response, got nil", tc.Name)
				}
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock()

	userID := int32(38)
	mockRepo.DeleteUserFunc = func(id int32) error {
		if id == userID {
			return nil
		}
		return errors.New("user not found")
	}

	userService := service.NewUserService(mockRepo)

	err := userService.DeleteUser(userID)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestBlockUser(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock()

	userID := int32(38)
	mockRepo.BlockUserFunc = func(id int32) error {
		if id == userID {
			return nil
		}
		return errors.New("user not found")
	}

	userService := service.NewUserService(mockRepo)

	err := userService.BlockUser(userID)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestUnblockUser(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock()

	userID := int32(38)
	mockRepo.UnblockUserFunc = func(id int32) error {
		if id == userID {
			return nil
		}
		return errors.New("user not found")
	}

	userService := service.NewUserService(mockRepo)

	err := userService.UnblockUser(userID)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
