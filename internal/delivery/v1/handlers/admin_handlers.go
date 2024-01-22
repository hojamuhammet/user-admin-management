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
	"strings"

	"github.com/go-chi/chi/v5"
)

type AdminHandler struct {
	AdminService *service.AdminService
	Router       chi.Router
}

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

	response := struct {
		Admins      *domain.AdminsList `json:"admins"`
		CurrentPage int                `json:"currentPage"`
		PrevPage    int                `json:"previousPage"`
		NextPage    int                `json:"nextPage"`
	}{
		Admins:      admins,
		CurrentPage: page,
		PrevPage:    previousPage,
		NextPage:    nextPage,
	}

	utils.RespondWithJSON(w, status.OK, response)
}

func (h *AdminHandler) GetAdminByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondWithErrorJSON(w, status.BadRequest, "Invalid ID")
		return
	}

	admin, err := h.AdminService.GetAdminByID(int32(id))
	if err != nil {
		slog.Error("Error retrieving admin: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, "Error retrieving admin")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(admin)
}

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

	updateAdminRequest.ID = int32(id)

	admin, err := h.AdminService.UpdateAdmin(&updateAdminRequest)
	if err != nil {
		slog.Error("Error updating admin: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, fmt.Sprintf("error updating admin: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status.OK)
	json.NewEncoder(w).Encode(admin)
}

func (h *AdminHandler) DeleteAdminHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondWithErrorJSON(w, status.BadRequest, errors.InvalidID)
		return
	}

	if err := h.AdminService.DeleteAdmin(int32(id)); err != nil {
		slog.Error("Error deleting admin: ", utils.Err(err))

		if strings.Contains(err.Error(), "not found") {
			errorMessage := fmt.Sprintf("Admin with ID %d not found", id)
			utils.RespondWithErrorJSON(w, status.NotFound, errorMessage)
			return
		}

		utils.RespondWithErrorJSON(w, status.InternalServerError, fmt.Sprintf("Error deleting admin: %s", err))
		return
	}

	utils.RespondWithJSON(w, status.OK, StatusMessage{
		Status:  status.OK,
		Message: "Admin deleted successfully",
	})
}

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

	response := struct {
		Admins      *domain.AdminsList `json:"admins"`
		CurrentPage int                `json:"currentPage"`
		PrevPage    int                `json:"previousPage"`
		NextPage    int                `json:"nextPage"`
	}{
		Admins:      admins,
		CurrentPage: page,
		PrevPage:    previousPage,
		NextPage:    nextPage,
	}

	utils.RespondWithJSON(w, status.OK, response)
}
