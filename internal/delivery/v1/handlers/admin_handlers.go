package handlers

import (
	"admin-panel/internal/domain"
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

type AdminHandler struct {
	AdminService *service.AdminService
	Router       chi.Router
}

// @Summary Get all admins
// @Description Retrieves a list of all administrators with pagination.
// @Tags admins
// @Accept json
// @Produce json
// @Security jwt
// @Param page query int false "Page number (default: 1)"
// @Param pageSize query int false "Page size (default: 8)"
// @Success 200 {object} domain.AdminListResponse
// @Failure 500 {string} string
// @Router /api/admin [get]
func (h *AdminHandler) GetAllAdminsHandler(w http.ResponseWriter, r *http.Request) {
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

	nextPage := page + 1

	admins, err := h.AdminService.GetAllAdmins(page, pageSize)
	if err != nil {
		slog.Error("Error getting admins: ", utils.Err(err))
		http.Error(w, errors.InternalServerError, status.InternalServerError)
		return
	}

	response := domain.AdminListResponse{
		Admins:      admins,
		CurrentPage: page,
		PrevPage:    previousPage,
		NextPage:    nextPage,
	}

	utils.RespondWithJSON(w, status.OK, response)
}

// @Summary Get admin by ID
// @Description Retrieves an administrator by ID.
// @Tags admins
// @Accept json
// @Produce json
// @Security jwt
// @Param id path int true "Admin ID"
// @Success 200 {object} domain.Admin
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/admin/{id} [get]
func (h *AdminHandler) GetAdminByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondWithErrorJSON(w, status.BadRequest, errors.InvalidID)
		return
	}

	admin, err := h.AdminService.GetAdminByID(int32(id))
	if err != nil {
		if err.Error() == "admin not found" {
			utils.RespondWithErrorJSON(w, status.NotFound, errors.AdminNotFound)
			return
		}

		slog.Error("Error retrieving admin: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, "Error retrieving admin")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(admin)
}

// @Summary Create admin
// @Description Creates a new administrator with the provided details.
// @Tags admins
// @Accept json
// @Produce json
// @Security jwt
// @Param admin body domain.CreateAdminRequest true "Admin data"
// @Success 200 {object} domain.Admin
// @Failure 400 {string} string
// @Failure 409 {string} string
// @Failure 500 {string} string
// @Router /api/admin [post]
func (h *AdminHandler) CreateAdminHandler(w http.ResponseWriter, r *http.Request) {
	var admin domain.CreateAdminRequest
	err := json.NewDecoder(r.Body).Decode(&admin)
	if err != nil {
		utils.RespondWithErrorJSON(w, status.BadRequest, "Invalid request body")
		return
	}

	if admin.Username == "" || admin.Password == "" || admin.Role == "" {
		utils.RespondWithErrorJSON(w, status.BadRequest, "Username, password, and role are required fields")
		return
	}

	createdAdmin, err := h.AdminService.CreateAdmin(&admin)
	if err != nil {
		switch err {
		case domain.ErrAdminAlreadyExists:
			utils.RespondWithErrorJSON(w, status.Conflict, "Admin with the same username already exists")
		default:
			utils.RespondWithErrorJSON(w, status.InternalServerError, fmt.Sprintf("Error creating admin: %v", err))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdAdmin)
}

// @Summary Update admin
// @Description Updates an existing administrator with the provided data.
// @Tags admins
// @Accept json
// @Produce json
// @Security jwt
// @Param id path int true "Admin ID"
// @Param admin body domain.UpdateAdminRequest true "Updated admin data"
// @Success 200 {object} domain.Admin
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/admin/{id} [put]
func (h *AdminHandler) UpdateAdminHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondWithErrorJSON(w, status.BadRequest, errors.InvalidID)
		return
	}

	var updateAdminRequest domain.UpdateAdminRequest

	err = json.NewDecoder(r.Body).Decode(&updateAdminRequest)
	if err != nil {
		utils.RespondWithErrorJSON(w, status.BadRequest, errors.InvalidRequestBody)
		return
	}

	admin, err := h.AdminService.UpdateAdmin(int32(id), &updateAdminRequest)
	if err != nil {
		if err == domain.ErrAdminNotFound {
			utils.RespondWithErrorJSON(w, status.NotFound, errors.AdminNotFound)
			return
		}

		slog.Error("Error updating admin: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, fmt.Sprintf("error updating admin: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status.OK)
	json.NewEncoder(w).Encode(admin)
}

// @Summary Delete admin
// @Description Deletes an administrator by their unique ID.
// @Tags admins
// @Accept json
// @Produce json
// @Security jwt
// @Param id path int true "Admin ID"
// @Success 200 {object} StatusMessage
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/admin/{id} [delete]
func (h *AdminHandler) DeleteAdminHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondWithErrorJSON(w, status.BadRequest, errors.InvalidID)
		return
	}

	if err := h.AdminService.DeleteAdmin(int32(id)); err != nil {
		if err == domain.ErrAdminNotFound {
			utils.RespondWithErrorJSON(w, status.NotFound, errors.AdminNotFound)
			return
		}

		slog.Error("Error deleting admin: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, fmt.Sprintf("Error deleting admin: %s", err))
		return
	}

	utils.RespondWithJSON(w, status.OK, StatusMessage{
		Status:  status.OK,
		Message: "Admin deleted successfully",
	})
}

// @Summary Search admins
// @Description Search administrators by query with pagination
// @Tags admins
// @Accept json
// @Produce json
// @Security jwt
// @Param query query string true "Search query"
// @Param page query int false "Page number (default: 1)"
// @Param pageSize query int false "Page size (default: 8)"
// @Success 200 {object} domain.AdminListResponse
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /api/admin/search [get]
func (h *AdminHandler) SearchAdminsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		utils.RespondWithErrorJSON(w, status.BadRequest, "Search query is required")
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

	admins, err := h.AdminService.SearchAdmins(query, page, pageSize)
	if err != nil {
		slog.Error("Error searching admins: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errors.InternalServerError)
		return
	}

	previousPage := page - 1
	if previousPage < 1 {
		previousPage = 1
	}

	nextPage := page + 1

	response := domain.AdminListResponse{
		Admins:      admins,
		CurrentPage: page,
		PrevPage:    previousPage,
		NextPage:    nextPage,
	}

	utils.RespondWithJSON(w, status.OK, response)
}
