package service

import (
	"admin-panel/internal/domain"
	"admin-panel/internal/repository"
	"admin-panel/pkg/lib/utils"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type AdminAuthService struct {
	AdminAuthRepository repository.AdminAuthRepository
}

func NewAdminAuthService(adminAuthRepository repository.AdminAuthRepository) *AdminAuthService {
	return &AdminAuthService{AdminAuthRepository: adminAuthRepository}
}

func (s *AdminAuthService) LoginAdmin(username, password string) (string, string, error) {
	admin, err := s.AdminAuthRepository.GetAdminByUsername(username)
	if err != nil {
		slog.Error("Error getting admin by username:", utils.Err(err))
		return "", "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
	if err != nil {
		slog.Error("Error comparing passwords:", utils.Err(err))
		return "", "", domain.ErrInvalidCredentials
	}

	accessToken, refreshToken, err := s.AdminAuthRepository.GenerateTokenPair(admin)
	if err != nil {
		slog.Error("Error generating token pair:", utils.Err(err))
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AdminAuthService) RefreshTokens(refreshToken string) (string, string, error) {
	claims, err := s.AdminAuthRepository.ValidateRefreshToken(refreshToken)
	if err != nil {
		slog.Error("Error validating refresh token:", utils.Err(err))
		return "", "", err
	}

	adminIDFloat, ok := claims["adminID"].(float64)
	if !ok {
		slog.Error("AdminID not found or not a number in refresh token claims")
		return "", "", domain.ErrInvalidRefreshToken
	}

	// Convert adminID to int
	adminID := int(adminIDFloat)

	admin, err := s.AdminAuthRepository.GetAdminByID(adminID)
	if err != nil {
		slog.Error("Error getting admin by ID:", utils.Err(err))
		return "", "", err
	}

	newAccessToken, newRefreshToken, err := s.AdminAuthRepository.GenerateTokenPair(admin)
	if err != nil {
		slog.Error("Error generating token pair:", utils.Err(err))
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *AdminAuthService) LogoutAdmin(refreshToken string) error {
	err := s.AdminAuthRepository.DeleteRefreshToken(refreshToken)
	if err != nil {
		slog.Error("Error deleting refresh token during logout:", utils.Err(err))
		return err
	}

	return nil
}
