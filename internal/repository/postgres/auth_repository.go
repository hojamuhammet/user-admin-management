package repository

import (
	"admin-panel/internal/domain"
	"admin-panel/pkg/lib/utils"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"admin-panel/internal/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type PostgresAdminAuthRepository struct {
	DB        *sql.DB
	JWTConfig config.JWT
}

func NewPostgresAdminAuthRepository(db *sql.DB, jwtConfig config.JWT) *PostgresAdminAuthRepository {
	return &PostgresAdminAuthRepository{DB: db, JWTConfig: jwtConfig}
}

const (
	accessTokenExpiration  = 30 * time.Minute
	refreshTokenExpiration = 7 * 24 * time.Hour
)

func (r *PostgresAdminAuthRepository) GenerateTokenPair(admin *domain.Admin) (string, string, error) {
	// Generate a new access token
	accessToken, err := r.generateAccessToken(admin)
	if err != nil {
		return "", "", err
	}

	// Generate a new refresh token
	refreshToken, err := r.generateRefreshToken(admin)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (r *PostgresAdminAuthRepository) ValidateRefreshToken(refreshToken string) (map[string]interface{}, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(refreshToken, jwt.MapClaims{})
	if err != nil {
		slog.Error("Error parsing refresh token: %v", err)
		return nil, fmt.Errorf("error parsing refresh token: %v", err)
	}

	adminIDClaim, ok := token.Claims.(jwt.MapClaims)["adminID"]
	if !ok {
		slog.Error("AdminID claim not found in refresh token")
		return nil, fmt.Errorf("adminID claim not found in refresh token")
	}

	query := `
        SELECT 1
        FROM admins
        WHERE refresh_token = $1 AND id = $2
    `

	var exists bool
	err = r.DB.QueryRow(query, refreshToken, adminIDClaim).Scan(&exists)
	if err != nil {
		slog.Error("Error checking refresh token existence in database: %v", utils.Err(err))
		return nil, fmt.Errorf("error checking refresh token existence in database: %v", err)
	}

	if !exists {
		slog.Error("Refresh token not found in the database")
		return nil, fmt.Errorf("refresh token not found in the database")
	}

	claims, err := r.validateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (r *PostgresAdminAuthRepository) DeleteRefreshToken(refreshToken string) error {
	query := `
        UPDATE admins
        SET refresh_token = NULL,
            refresh_token_created_at = NULL,
            refresh_token_expiration_time = NULL
        WHERE refresh_token = $1
    `

	_, err := r.DB.Exec(query, refreshToken)
	if err != nil {
		slog.Error("Error deleting refresh token: %v", utils.Err(err))
		return err
	}

	return nil
}

func (r *PostgresAdminAuthRepository) GetAdminByUsername(username string) (*domain.Admin, error) {
	query := `
		SELECT id, username, password, role
		FROM admins
		WHERE username = $1
		LIMIT 1
	`

	row := r.DB.QueryRow(query, username)

	var admin domain.Admin

	err := row.Scan(&admin.ID, &admin.Username, &admin.Password, &admin.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrAdminNotFound
		}

		slog.Error("Error getting admin by username: %v", err)
		return nil, err
	}

	return &admin, nil
}

func (r *PostgresAdminAuthRepository) GetAdminByID(adminID int) (*domain.Admin, error) {
	query := `
		SELECT id, username, password, role
		FROM admins
		WHERE id = $1
		LIMIT 1
	`

	row := r.DB.QueryRow(query, adminID)

	var admin domain.Admin

	err := row.Scan(&admin.ID, &admin.Username, &admin.Password, &admin.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("Admin not found")
			return nil, domain.ErrAdminNotFound
		}

		slog.Error("Error getting admin by ID: %v", err)
		return nil, err
	}

	return &admin, nil
}

func (r *PostgresAdminAuthRepository) generateAccessToken(admin *domain.Admin) (string, error) {
	claims := jwt.MapClaims{
		"id":   admin.ID,
		"role": admin.Role,
		"exp":  time.Now().Add(accessTokenExpiration).Unix(), // Token expiration time
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(r.JWTConfig.AccessSecretKey))
	if err != nil {
		slog.Error("Error generating access token: %v", utils.Err(err))
		return "", err
	}

	return tokenString, nil
}

func (r *PostgresAdminAuthRepository) generateRefreshToken(admin *domain.Admin) (string, error) {
	refreshTokenID := uuid.New().String()

	refreshClaims := jwt.MapClaims{
		"id":      refreshTokenID,
		"adminID": admin.ID,
		"exp":     time.Now().Add(refreshTokenExpiration).Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	refreshTokenString, err := refreshToken.SignedString([]byte(r.JWTConfig.RefreshSecretKey))
	if err != nil {
		return "", err
	}

	query := `
        UPDATE admins
        SET refresh_token = $1,
            refresh_token_created_at = CURRENT_TIMESTAMP,
            refresh_token_expiration_time = TO_TIMESTAMP($2)
        WHERE id = $3
    `

	_, err = r.DB.Exec(query, refreshTokenString, refreshClaims["exp"].(int64), admin.ID)
	if err != nil {
		slog.Error("Failed to update refresh token in database", utils.Err(err))
		return "", err
	}

	return refreshTokenString, nil
}

func (r *PostgresAdminAuthRepository) validateRefreshToken(refreshToken string) (map[string]interface{}, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(r.JWTConfig.RefreshSecretKey), nil
	})

	if err != nil || !token.Valid {
		slog.Error("Refresh token validation error: %v !BADKEY=\"%s\"", err, r.JWTConfig.RefreshSecretKey)
		return nil, fmt.Errorf("refresh token validation error: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims == nil {
		slog.Error("Invalid refresh token claims")
		return nil, fmt.Errorf("invalid refresh token claims")
	}

	return claims, nil
}
