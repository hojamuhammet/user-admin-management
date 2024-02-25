package service

import (
	"admin-panel/internal/domain"
	repository "admin-panel/internal/repository/interfaces"
	service "admin-panel/internal/service/interfaces"
)

type AdminService struct {
	AdminRepository repository.AdminRepository
}

func NewAdminService(adminRepository repository.AdminRepository) *AdminService {
	return &AdminService{AdminRepository: adminRepository}
}

func (s *AdminService) GetAllAdmins(page, pageSize int) (*domain.AdminsList, error) {
	return s.AdminRepository.GetAllAdmins(page, pageSize)
}

func (s *AdminService) GetTotalAdminsCount() (int, error) {
	return s.AdminRepository.GetTotalAdminsCount()
}

func (s *AdminService) GetAdminByID(id int32) (*domain.CommonAdminResponse, error) {
	return s.AdminRepository.GetAdminByID(id)
}

func (s *AdminService) CreateAdmin(request *domain.CreateAdminRequest) (*domain.CommonAdminResponse, error) {
	return s.AdminRepository.CreateAdmin(request)
}

func (s *AdminService) UpdateAdmin(id int32, request *domain.UpdateAdminRequest) (*domain.CommonAdminResponse, error) {
	return s.AdminRepository.UpdateAdmin(id, request)
}

func (s *AdminService) DeleteAdmin(id int32) error {
	return s.AdminRepository.DeleteAdmin(id)
}

func (s *AdminService) SearchAdmins(query string, page, pageSize int) (*domain.AdminsList, error) {
	return s.AdminRepository.SearchAdmins(query, page, pageSize)
}

var _ service.AdminService = &AdminService{}
