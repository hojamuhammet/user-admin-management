package routers

import (
	"admin-panel/internal/delivery/v1/handlers"
	repository "admin-panel/internal/repository/interfaces"
	service "admin-panel/internal/service/interfaces"

	"github.com/go-chi/chi/v5"
)

func SetupUserRoutes(userRepository repository.UserRepository, userService service.UserService, userRouter *chi.Mux) {
	userHandler := handlers.NewUserHandler(userRepository, userService, userRouter)

	userRouter.Get("/", userHandler.GetAllUsersHandler)
	userRouter.Get("/{id}", userHandler.GetUserByIDHandler)
	userRouter.Post("/", userHandler.CreateUserHandler)
	userRouter.Put("/{id}", userHandler.UpdateUserHandler)
	userRouter.Delete("/{id}", userHandler.DeleteUserHandler)
	userRouter.Post("/{id}/block", userHandler.BlockUserHandler)
	userRouter.Post("/{id}/unblock", userHandler.UnblockUserHandler)
	userRouter.Get("/search", userHandler.SearchUsersHandler)
}
