package service

import "admin-panel/internal/domain"

type AdminService interface {
	GetAllAdmins(page, pageSize int) (*domain.AdminsList, error)
	GetTotalAdminsCount() (int, error)
	GetAdminByID(id int32) (*domain.GetAdminResponse, error)
	CreateAdmin(request *domain.CreateAdminRequest) (*domain.CreateAdminResponse, error)
	UpdateAdmin(id int32, request *domain.UpdateAdminRequest) (*domain.UpdateAdminResponse, error)
	DeleteAdmin(id int32) error
	SearchAdmins(query string, page, pageSize int) (*domain.AdminsList, error)
}
