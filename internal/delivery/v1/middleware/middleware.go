package middleware

import (
	"admin-panel/internal/config"
	"admin-panel/pkg/lib/errors"
	"admin-panel/pkg/lib/status"
	"admin-panel/pkg/lib/utils"
	"context"
	"fmt"
	"net/http"
	"strings"

	"log/slog"

	"github.com/dgrijalva/jwt-go"
)

// contextKey is a custom type for the context key used to store the JWT claims.
type contextKey string

const (
	// tokenKey is the context key for storing the JWT claims in the context.
	tokenKey contextKey = "token"
)

func AuthMiddleware(cfg *config.Config, allowedRoles []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := extractTokenFromHeader(r)
			if tokenString == "" {
				utils.RespondWithErrorJSON(w, status.Unauthorized, errors.AuthorizationTokenNotProvided)
				return
			}

			claims, err := validateToken(tokenString, cfg, isRefreshToken(r))
			if err != nil {
				utils.RespondWithErrorJSON(w, status.Unauthorized, fmt.Sprintf("Invalid authorization token: %v", err))
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, tokenKey, claims)

			if hasRequiredRole(claims["role"].(string), []string{"super_admin"}) {
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			if !hasRequiredRole(claims["role"].(string), allowedRoles) {
				utils.RespondWithErrorJSON(w, status.Forbidden, errors.InsufficientPermission)
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func validateToken(tokenString string, cfg *config.Config, isRefreshToken bool) (jwt.MapClaims, error) {
	var secretKey string

	if isRefreshToken {
		secretKey = cfg.JWT.RefreshSecretKey
	} else {
		secretKey = cfg.JWT.AccessSecretKey
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		slog.Error("Token validation error: %v", utils.Err(err))
		return nil, fmt.Errorf("token validation error: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims == nil {
		slog.Error("Invalid token claims")
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func isRefreshToken(r *http.Request) bool {
	return strings.Contains(r.URL.Path, "/refresh")
}

func extractTokenFromHeader(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		slog.Error("Authorization header not found")
		return ""
	}

	if !strings.HasPrefix(bearerToken, "Bearer ") {
		slog.Error("Invalid Authorization header format")
		return ""
	}

	return strings.TrimPrefix(bearerToken, "Bearer ")
}

func hasRequiredRole(adminRole string, allowedRoles []string) bool {
	for _, allowedRole := range allowedRoles {
		if adminRole == allowedRole {
			return true
		}
	}
	return false
}
