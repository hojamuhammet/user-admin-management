package repository

import (
	"admin-panel/internal/domain"
)

type AuthRepository interface {
	GetAdminByUsername(username string) (*domain.Admin, error)
	GenerateTokenPair(admin *domain.Admin) (string, string, error)
	ValidateRefreshToken(refreshToken string) (map[string]interface{}, error)
	GetAdminByID(adminID int) (*domain.Admin, error)
	DeleteRefreshToken(refreshToken string) error
}
