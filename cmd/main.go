package main

import (
	"admin-panel/internal/config"
	"admin-panel/internal/delivery/v1/middleware"
	"admin-panel/internal/delivery/v1/routers"
	repository "admin-panel/internal/repository/postgres"
	"admin-panel/internal/service"
	"admin-panel/pkg/database"
	utils "admin-panel/pkg/lib/utils"
	"admin-panel/pkg/logger"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "admin-panel/docs"

	"github.com/go-chi/chi/v5"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Admin Panel API
// @version 1.0
// @description API for Admin Panel.
// @description This API provides endpoints to manage users, administrators, and authentication in the admin panel.
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey jwt
// @in header
// @name Authorization
func main() {
	cfg := config.LoadConfig()

	log := logger.SetupLogger(cfg.Env)

	slog.Info("Starting the server...", slog.String("env", cfg.Env))
	slog.Debug("Debug messages are enabled") // If env is set to prod, debug messages are going to be disabled

	db, err := database.InitDB(cfg)
	if err != nil {
		slog.Error("Failed to init database:", utils.Err(err))
		os.Exit(1)
	}
	defer db.Close()

	mainRouter := chi.NewRouter()

	authMiddlewareForAdmin := middleware.AuthMiddleware(cfg, []string{"admin"})
	authMiddlewareForSuperAdmin := middleware.AuthMiddleware(cfg, []string{"super_admin"})

	// Authentication routes
	authRouter := chi.NewRouter()
	mainRouter.Route("/auth", func(r chi.Router) {
		r.Mount("/", authRouter)
	})

	adminAuthRepository := repository.NewPostgresAdminAuthRepository(db.GetDB(), cfg.JWT)
	adminAuthService := service.NewAdminAuthService(adminAuthRepository)
	routers.SetupAuthRoutes(authRouter, adminAuthService)

	// Admin routes
	adminRouter := chi.NewRouter()
	adminRouter.Use(authMiddlewareForSuperAdmin) // Apply auth middleware to admin routes
	mainRouter.Route("/api/admin", func(r chi.Router) {
		r.Mount("/", adminRouter)
	})

	adminRepository := repository.NewPostgresAdminRepository(db.GetDB())
	adminService := service.NewAdminService(adminRepository)
	routers.SetupAdminRoutes(adminRouter, adminService)

	// User routes
	userRouter := chi.NewRouter()
	userRouter.Use(authMiddlewareForAdmin)
	mainRouter.Route("/api/user", func(r chi.Router) {
		r.Mount("/", userRouter)
	})

	userRepository := repository.NewPostgresUserRepository(db.GetDB())
	userService := service.NewUserService(userRepository)
	routers.SetupUserRoutes(userRouter, userService)

	mainRouter.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	// Handling graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		log.Info("Shutting down the server gracefully...")

		if err := db.Close(); err != nil {
			slog.Error("Error closing database:", utils.Err(err))
		}
		os.Exit(0)
	}()

	err = http.ListenAndServe(cfg.Address, mainRouter)
	if err != nil {
		slog.Error("Server failed to start:", utils.Err(err))
	}
}
