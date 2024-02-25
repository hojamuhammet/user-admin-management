package service

import (
	repository "admin-panel/internal/repository/interfaces"
	service "admin-panel/internal/service/interfaces"
	"admin-panel/pkg/lib/errors"
	"admin-panel/pkg/lib/utils"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	AuthRepository repository.AuthRepository
}

func NewAuthService(authRepository repository.AuthRepository) *AuthService {
	return &AuthService{AuthRepository: authRepository}
}

func (s *AuthService) LoginAdmin(username, password string) (string, string, error) {
	admin, err := s.AuthRepository.GetAdminByUsername(username)
	if err != nil {
		slog.Error("Error getting admin by username:", utils.Err(err))
		return "", "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
	if err != nil {
		slog.Error("Error comparing passwords:", utils.Err(err))
		return "", "", errors.ErrInvalidCredentials
	}

	accessToken, refreshToken, err := s.AuthRepository.GenerateTokenPair(admin)
	if err != nil {
		slog.Error("Error generating token pair:", utils.Err(err))
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) RefreshTokens(refreshToken string) (string, string, error) {
	claims, err := s.AuthRepository.ValidateRefreshToken(refreshToken)
	if err != nil {
		slog.Error("Error validating refresh token:", utils.Err(err))
		return "", "", err
	}

	adminIDFloat, ok := claims["adminID"].(float64)
	if !ok {
		slog.Error("AdminID not found or not a number in refresh token claims")
		return "", "", errors.ErrInvalidRefreshToken
	}

	adminID := int(adminIDFloat)

	admin, err := s.AuthRepository.GetAdminByID(adminID)
	if err != nil {
		slog.Error("Error getting admin by ID:", utils.Err(err))
		return "", "", err
	}

	newAccessToken, newRefreshToken, err := s.AuthRepository.GenerateTokenPair(admin)
	if err != nil {
		slog.Error("Error generating token pair:", utils.Err(err))
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *AuthService) LogoutAdmin(refreshToken string) error {
	err := s.AuthRepository.DeleteRefreshToken(refreshToken)
	if err != nil {
		slog.Error("Error deleting refresh token during logout:", utils.Err(err))
		return err
	}

	return nil
}

var _ service.AuthService = &AuthService{}
