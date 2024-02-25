package routers

import (
	"admin-panel/internal/delivery/v1/handlers"
	repository "admin-panel/internal/repository/interfaces"
	service "admin-panel/internal/service/interfaces"

	"github.com/go-chi/chi/v5"
)

func SetupAdminRoutes(adminRepository repository.AdminRepository, adminService service.AdminService, adminRouter *chi.Mux) {
	adminHandler := handlers.NewAdminHandler(adminRepository, adminService, adminRouter)

	adminRouter.Get("/", adminHandler.GetAllAdminsHandler)
	adminRouter.Get("/{id}", adminHandler.GetAdminByID)
	adminRouter.Post("/", adminHandler.CreateAdminHandler)
	adminRouter.Put("/{id}", adminHandler.UpdateAdminHandler)
	adminRouter.Delete("/{id}", adminHandler.DeleteAdminHandler)
	adminRouter.Get("/search", adminHandler.SearchAdminsHandler)
}
