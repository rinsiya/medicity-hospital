package service

import (
	"medicity/internal/dto"
	//"medicity/internal/models"
	"medicity/internal/repository"
)

type DepartmentService interface {
	GetAllDepartments() ([]dto.DepartmentResponse, error)
}

type departmentService struct {
	departmentRepo repository.DepartmentRepository
}

func NewDepartmentService(
	departmentRepo repository.DepartmentRepository,
) DepartmentService {
	return &departmentService{
		departmentRepo: departmentRepo,
	}
}

func (s *departmentService) GetAllDepartments() ([]dto.DepartmentResponse, error) {

	return s.departmentRepo.GetAllDepartments()
}