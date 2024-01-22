package routers

import (
	"admin-panel/internal/delivery/v1/handlers"
	"admin-panel/internal/service"

	"github.com/go-chi/chi/v5"
)

func SetupAuthRoutes(authRouter *chi.Mux, adminAuthService *service.AdminAuthService) {
	authHandler := handlers.AdminAuthHandler{
		AdminAuthService: *adminAuthService,
		Router:           authRouter,
	}

	authRouter.Post("/login", authHandler.LoginHandler)
	authRouter.Post("/refresh", authHandler.RefreshTokensHandler)
	authRouter.Post("/logout", authHandler.LogoutHandler)
}
