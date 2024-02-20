package handlers

import (
	"admin-panel/internal/domain"
	repository "admin-panel/internal/repository/postgres"
	"admin-panel/internal/service"
	"admin-panel/pkg/lib/errors"
	"admin-panel/pkg/lib/status"
	"admin-panel/pkg/lib/utils"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	UserService    *service.UserService
	UserRepository *repository.PostgresUserRepository
	Router         *chi.Mux
}

// @Summary Get all users
// @Description Retrieves a list of all users with pagination.
// @Tags users
// @Accept json
// @Produce json
// @Security jwt
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} domain.UsersListResponse "Success"
// @Failure 500 {string} string "Internal Server Error: " + errors.InternalServerError
// @Router /api/user [get]
func (h *UserHandler) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil || pageSize <= 0 {
		pageSize = 8 // Default page size
	}

	var previousPage int
	if page > 1 {
		previousPage = page - 1
	} else {
		previousPage = 1
	}

	users, err := h.UserService.GetAllUsers(page, pageSize)
	if err != nil {
		slog.Error("Error getting users: ", utils.Err(err))
		http.Error(w, errors.InternalServerError, status.InternalServerError)
		return
	}

	totalUsers, err := h.UserRepository.GetTotalUsersCount()
	if err != nil {
		slog.Error("Error getting total users count: ", utils.Err(err))
		http.Error(w, errors.InternalServerError, status.InternalServerError)
		return
	}

	firstPage := 1
	lastPage := (totalUsers + pageSize - 1) / pageSize

	nextPage := page + 1
	if page >= lastPage {
		nextPage = lastPage
	}

	response := domain.UsersListResponse{

		Users:       users,
		CurrentPage: page,
		PrevPage:    previousPage,
		NextPage:    nextPage,
		FirstPage:   firstPage,
		LastPage:    lastPage,
	}

	utils.RespondWithJSON(w, status.OK, response)
}

// @Summary Get user by ID
// @Description Retrieves a user by ID.
// @Tags users
// @Accept json
// @Produce json
// @Security jwt
// @Param id path int true "User ID"
// @Success 200 {object} domain.GetUserResponse "Success"
// @Failure 400 {string} string "Bad Request: " + errors.InvalidID
// @Failure 404 {string} string "User not found" + errors.UserNotFound
// @Failure 500 {string} string "Internal Server Error: " + errors.InternalServerError
// @Router /api/user/{id} [get]
func (h *UserHandler) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondWithErrorJSON(w, status.BadRequest, errors.InvalidID)
		return
	}

	user, err := h.UserService.GetUserByID(int32(id))
	if err != nil {
		if err.Error() == "user not found" {
			utils.RespondWithErrorJSON(w, status.NotFound, errors.UserNotFound)
			return
		}

		slog.Error("Error retrieving user: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, "Error retrieving user")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// @Summary Create user
// @Description Creates a new administrator with the provided details.
// @Tags users
// @Accept json
// @Produce json
// @Security jwt
// @Param request body domain.CreateUserRequest true "User creation request"
// @Success 201 {object} domain.CreateUserResponse "Created"
// @Failure 400 {string} string "Bad Request: " + errors.InvalidRequestBody
// @Failure 500 {string} string "Internal Server Error: " + errors.InternalServerError
// @Router /api/user [post]
func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var createUserRequest domain.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&createUserRequest)
	if err != nil {
		utils.RespondWithErrorJSON(w, status.BadRequest, errors.InvalidRequestBody)
		return
	}

	if !utils.IsValidPhoneNumber(createUserRequest.PhoneNumber) {
		utils.RespondWithErrorJSON(w, status.BadRequest, errors.InvalidPhoneNumberFormat)
		return
	}

	user, err := h.UserService.CreateUser(&createUserRequest)
	if err != nil {
		slog.Error("Error creating user: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, fmt.Sprintf("error creating user: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// @Summary Update user
// @Description Updates an existing user with the provided data.
// @Tags users
// @Accept json
// @Produce json
// @Security jwt
// @Param id path int true "User ID"
// @Param request body domain.UpdateUserRequest true "User update request"
// @Success 200 {object} domain.UpdateUserResponse "Updated"
// @Failure 400 {string} string "Bad Request: " + errors.InvalidID or errors.InvalidRequestBody
// @Failure 404 {string} string "User not found: " + errors.UserNotFound
// @Failure 500 {string} string "Internal Server Error: " + errors.InternalServerError
// @Router /api/user/{id} [put]
func (h *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondWithErrorJSON(w, status.BadRequest, errors.InvalidID)
		return
	}

	var updateUserRequest domain.UpdateUserRequest

	err = json.NewDecoder(r.Body).Decode(&updateUserRequest)
	if err != nil {
		utils.RespondWithErrorJSON(w, status.BadRequest, errors.InvalidRequestBody)
		return
	}

	user, err := h.UserService.UpdateUser(int32(id), &updateUserRequest)
	if err != nil {
		if err == domain.ErrUserNotFound {
			utils.RespondWithErrorJSON(w, status.NotFound, errors.UserNotFound)
			return
		}
		slog.Error("Error updating user: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, fmt.Sprintf("error updating user: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status.OK)
	json.NewEncoder(w).Encode(user)
}

// @Summary Delete user by ID
// @Description Deletes a user by their unique ID.
// @Tags users
// @Accept json
// @Produce json
// @Security jwt
// @Param id path int true "User ID"
// @Success 200 {object} StatusMessage "Deleted"
// @Failure 400 {string} string "Bad Request: " + errors.InvalidID
// @Failure 404 {string} string "User not found: " + errors.UserNotFound
// @Failure 500 {string} string "Internal Server Error: " + errors.InternalServerError
// @Router /api/user/{id} [delete]
func (h *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondWithErrorJSON(w, status.BadRequest, errors.InvalidID)
		return
	}

	if err := h.UserService.DeleteUser(int32(id)); err != nil {
		if err == domain.ErrUserNotFound {
			utils.RespondWithErrorJSON(w, status.NotFound, errors.UserNotFound)
			return
		}

		slog.Error("Error deleting user: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, fmt.Sprintf("Error deleting user: %s", err))
		return
	}

	utils.RespondWithJSON(w, status.OK, StatusMessage{
		Status:  status.OK,
		Message: "User deleted successfully",
	})
}

// @Summary Block user by ID
// @Description Blocks a user by their unique ID.
// @Tags users
// @Accept json
// @Produce json
// @Security jwt
// @Param id path int true "User ID"
// @Success 200 {object} StatusMessage "Blocked"
// @Failure 400 {string} string "Bad Request: " + errors.InvalidID
// @Failure 404 {string} string "User not found: " + errors.UserNotFound
// @Failure 500 {string} string "Internal Server Error: " + errors.InternalServerError
// @Router /api/user/{id}/block [post]
func (h *UserHandler) BlockUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondWithErrorJSON(w, status.BadRequest, errors.InvalidID)
		return
	}

	if err := h.UserService.BlockUser(int32(id)); err != nil {
		if err == domain.ErrUserNotFound {
			utils.RespondWithErrorJSON(w, status.NotFound, errors.UserNotFound)
			return
		}

		slog.Error("Error blocking user: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, fmt.Sprintf("Error blocking user: %s", err))
		return
	}

	utils.RespondWithJSON(w, status.OK, StatusMessage{
		Status:  status.OK,
		Message: "User blocked successfully",
	})
}

// @Summary Unblock user by ID
// @Description Unblocks a user by their unique ID.
// @Tags users
// @Accept json
// @Produce json
// @Security jwt
// @Param id path int true "User ID"
// @Success 200 {object} StatusMessage "Unblocked"
// @Failure 400 {string} string "Bad Request: " + errors.InvalidID
// @Failure 404 {string} string "User not found: " + errors.UserNotFound
// @Failure 500 {string} string "Internal Server Error: " + errors.InternalServerError
// @Router /api/user/{id}/unblock [post]
func (h *UserHandler) UnblockUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondWithErrorJSON(w, status.BadRequest, errors.InvalidID)
		return
	}

	if err := h.UserService.UnblockUser(int32(id)); err != nil {
		if err == domain.ErrUserNotFound {
			utils.RespondWithErrorJSON(w, status.NotFound, errors.UserNotFound)
			return
		}

		slog.Error("Error unblocking user: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, fmt.Sprintf("Error unblocking user: %s", err))
		return
	}

	utils.RespondWithJSON(w, status.OK, StatusMessage{
		Status:  status.OK,
		Message: "User unblocked successfully",
	})
}

// @Summary Search users
// @Description Search users by query with pagination
// @Tags users
// @Accept json
// @Produce json
// @Security jwt
// @Param query query string true "Search query"
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} domain.UsersListResponse "Success"
// @Failure 400 {string} string "Bad Request: " + errors.SearchQueryRequired
// @Failure 500 {string} string "Internal Server Error: " + errors.InternalServerError
// @Router /api/user/search [get]
func (h *UserHandler) SearchUsersHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		utils.RespondWithErrorJSON(w, status.BadRequest, errors.SearchQueryRequired)
		return
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil || pageSize <= 0 {
		pageSize = 8 // Default page size
	}

	users, err := h.UserService.SearchUsers(query, page, pageSize)
	if err != nil {
		slog.Error("Error searching users: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errors.InternalServerError)
		return
	}

	totalUsers, err := h.UserRepository.GetTotalUsersCount()
	if err != nil {
		slog.Error("Error getting total users count: ", utils.Err(err))
		http.Error(w, errors.InternalServerError, status.InternalServerError)
		return
	}

	firstPage := 1
	lastPage := (totalUsers + pageSize - 1) / pageSize

	var previousPage int
	if page > 1 {
		previousPage = page - 1
	} else {
		previousPage = 1
	}

	nextPage := page + 1
	if page >= lastPage {
		nextPage = lastPage
	}

	response := domain.UsersListResponse{

		Users:       users,
		CurrentPage: page,
		PrevPage:    previousPage,
		NextPage:    nextPage,
		FirstPage:   firstPage,
		LastPage:    lastPage,
	}

	utils.RespondWithJSON(w, status.OK, response)
}
