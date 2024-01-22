package handlers

import (
	"admin-panel/internal/domain"
	"admin-panel/internal/service"
	"admin-panel/pkg/lib/errors"
	"admin-panel/pkg/lib/status"
	"admin-panel/pkg/lib/utils"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type AdminAuthHandler struct {
	AdminAuthService service.AdminAuthService
	Router           *chi.Mux
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type StatusMessage struct {
	Status  int    `json:"code"`
	Message string `json:"message"`
}

func (h *AdminAuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginRequest LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		slog.Error("Error decoding login request:", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.BadRequest, errors.InvalidRequestFormat)
		return
	}

	accessToken, refreshToken, err := h.AdminAuthService.LoginAdmin(loginRequest.Username, loginRequest.Password)
	if err != nil {
		switch err {
		case domain.ErrAdminNotFound:
			utils.RespondWithErrorJSON(w, status.NotFound, errors.AdminNotFound)
		default:
			slog.Error("Error during login:", utils.Err(err))
			utils.RespondWithErrorJSON(w, status.Unauthorized, errors.InvalidCredentials)
		}
		return
	}

	loginResponse := LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	utils.RespondWithJSON(w, status.OK, loginResponse)
}

func (h *AdminAuthHandler) RefreshTokensHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken := extractTokenFromHeader(r)
	if refreshToken == "" {
		slog.Error("Refresh token is not provided")
		utils.RespondWithErrorJSON(w, status.Unauthorized, errors.RefreshTokenNotProvided)
		return
	}

	newAccessToken, newRefreshToken, err := h.AdminAuthService.RefreshTokens(refreshToken)
	if err != nil {
		slog.Error("Error refreshing tokens:", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.Unauthorized, errors.InvalidRefreshToken)
		return
	}

	utils.RespondWithJSON(w, status.OK, map[string]string{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

func (h *AdminAuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	var requestData map[string]string
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		utils.RespondWithErrorJSON(w, status.BadRequest, errors.InvalidRequestFormat)
		return
	}

	refreshToken := requestData["refresh_token"]
	if refreshToken == "" {
		utils.RespondWithErrorJSON(w, status.BadRequest, errors.RefreshTokenNotProvided)
		return
	}

	err := h.AdminAuthService.LogoutAdmin(refreshToken)
	if err != nil {
		utils.RespondWithErrorJSON(w, status.InternalServerError, errors.InternalServerError)
		return
	}

	response := StatusMessage{
		Status:  status.OK,
		Message: "Logout successful",
	}

	utils.RespondWithJSON(w, status.OK, response)
}

func extractTokenFromHeader(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		slog.Error("Authorization header not found")
		return ""
	}

	return strings.TrimPrefix(bearerToken, "Bearer ")
}
