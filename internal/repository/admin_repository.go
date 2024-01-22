package repository

import "admin-panel/internal/domain"

type AdminRepository interface {
	GetAllAdmins(page, pageSize int) (*domain.AdminsList, error)
	GetAdminByID(id int32) (*domain.CommonAdminResponse, error)
	CreateAdmin(request *domain.CreateAdminRequest) (*domain.CommonAdminResponse, error)
	UpdateAdmin(request *domain.UpdateAdminRequest) (*domain.CommonAdminResponse, error)
	DeleteAdmin(id int32) error
	SearchAdmins(query string, page, pageSize int) (*domain.AdminsList, error)
}
