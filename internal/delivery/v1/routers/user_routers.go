package routers

import (
	"admin-panel/internal/delivery/v1/handlers"
	repository "admin-panel/internal/repository/postgres"
	"admin-panel/internal/service"

	"github.com/go-chi/chi/v5"
)

func SetupUserRoutes(userRouter *chi.Mux, userService *service.UserService, userRepository *repository.PostgresUserRepository) {
	userHandler := handlers.NewUserHandler(userService, userRepository, userRouter)

	userRouter.Get("/", userHandler.GetAllUsersHandler)
	userRouter.Get("/{id}", userHandler.GetUserByIDHandler)
	userRouter.Post("/", userHandler.CreateUserHandler)
	userRouter.Put("/{id}", userHandler.UpdateUserHandler)
	userRouter.Delete("/{id}", userHandler.DeleteUserHandler)
	userRouter.Post("/{id}/block", userHandler.BlockUserHandler)
	userRouter.Post("/{id}/unblock", userHandler.UnblockUserHandler)
	userRouter.Get("/search", userHandler.SearchUsersHandler)
}
