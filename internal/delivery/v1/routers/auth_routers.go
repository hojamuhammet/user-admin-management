package routers

import (
	"admin-panel/internal/delivery/v1/handlers"
	repository "admin-panel/internal/repository/interfaces"
	service "admin-panel/internal/service/interfaces"

	"github.com/go-chi/chi/v5"
)

func SetupAuthRoutes(AuthRepository repository.AuthRepository, AuthService service.AuthService, authRouter *chi.Mux) {
	authHandler := handlers.AuthHandler{
		AuthRepository: AuthRepository,
		AuthService:    AuthService,
		Router:         authRouter,
	}

	authRouter.Post("/login", authHandler.LoginHandler)
	authRouter.Post("/refresh", authHandler.RefreshTokensHandler)
	authRouter.Post("/logout", authHandler.LogoutHandler)
}
