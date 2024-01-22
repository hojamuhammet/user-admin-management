package routers

import (
	"admin-panel/internal/delivery/v1/handlers"
	"admin-panel/internal/service"

	"github.com/go-chi/chi/v5"
)

func SetupAdminRoutes(adminRouter *chi.Mux, adminService *service.AdminService) {
	adminHandler := handlers.AdminHandler{
		AdminService: adminService,
		Router:       adminRouter,
	}

	adminRouter.Get("/", adminHandler.GetAllAdminsHandler)
	adminRouter.Get("/{id}", adminHandler.GetAdminByID)
	adminRouter.Post("/", adminHandler.CreateAdminHandler)
	adminRouter.Put("/{id}", adminHandler.UpdateAdminHandler)
	adminRouter.Delete("/{id}", adminHandler.DeleteAdminHandler)
	adminRouter.Get("/search", adminHandler.SearchAdminsHandler)
}
