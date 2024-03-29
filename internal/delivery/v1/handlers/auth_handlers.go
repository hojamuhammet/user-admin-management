package handlers

import (
	repository "admin-panel/internal/repository/interfaces"
	service "admin-panel/internal/service/interfaces"
	"admin-panel/pkg/lib/errors"
	"admin-panel/pkg/lib/status"
	"admin-panel/pkg/lib/utils"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	AuthRepository repository.AuthRepository
	AuthService    service.AuthService
	Router         *chi.Mux
}

func NewAuthHandler(repository repository.AuthRepository, service service.AuthService, router *chi.Mux) *AuthHandler {
	return &AuthHandler{
		AuthRepository: repository,
		AuthService:    service,
		Router:         router,
	}
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
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// @Summary Admin Login
// @Description Logs in an admin and returns access and refresh tokens.
// @Tags auth
// @Accept json
// @Produce json
// @Param loginRequest body LoginRequest true "Admin login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} StatusMessage
// @Failure 401 {object} StatusMessage
// @Router /auth/login [post]
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginRequest LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		slog.Error("Error decoding login request:", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.BadRequest, errors.InvalidRequestFormat)
		return
	}

	accessToken, refreshToken, err := h.AuthService.LoginAdmin(loginRequest.Username, loginRequest.Password)
	if err != nil {
		if err == errors.ErrAdminNotFound {
			utils.RespondWithErrorJSON(w, status.NotFound, errors.AdminNotFound)
		} else {
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

// @Summary Refresh Tokens
// @Description Provide with your refresh token in header to make new refresh and access token pair.
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Refresh token to renew access and refresh tokens" default(Bearer your_refresh_token)
// @Success 200 {object} map[string]string
// @Failure 401 {object} StatusMessage
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshTokensHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken := extractTokenFromHeader(r)
	if refreshToken == "" {
		slog.Error("Refresh token is not provided")
		utils.RespondWithErrorJSON(w, status.Unauthorized, errors.RefreshTokenNotProvided)
		return
	}

	newAccessToken, newRefreshToken, err := h.AuthService.RefreshTokens(refreshToken)
	if err != nil {
		slog.Error("Error refreshing tokens:", utils.Err(err))
		if err == errors.ErrRefreshTokenExpired {
			utils.RespondWithErrorJSON(w, status.Unauthorized, errors.RefreshTokenExpired)
		} else if err == errors.ErrRefreshNotFoundInDB {
			utils.RespondWithErrorJSON(w, status.Unauthorized, errors.RefreshNotFoundInDB)
		} else {
			utils.RespondWithErrorJSON(w, status.Unauthorized, err.Error())
		}
		return
	}

	utils.RespondWithJSON(w, status.OK, map[string]string{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

// @Summary Admin Logout
// @Description Provide your refresh token in body of request to log out an admin by invalidating the provided refresh token.
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param refresh_token body string true "Refresh token to be invalidated" example:"your_refresh_token"
// @Success 200 {object} StatusMessage "Logout successful"
// @Failure 400 {object} StatusMessage "Invalid request format"
// @Failure 401 {object} StatusMessage "Refresh token not provided"
// @Failure 500 {object} StatusMessage "Internal server error"
// @Router /auth/logout [post]
func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
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

	err := h.AuthService.LogoutAdmin(refreshToken)
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
