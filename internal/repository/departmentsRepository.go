package repository

import (
	"medicity/database"
	"medicity/internal/models"
	"medicity/internal/dto"
)

type DepartmentRepository interface {
	GetAllDepartments() ([]dto.DepartmentResponse, error)
}

type departmentRepository struct{}

func NewDepartmentRepository() DepartmentRepository {
	return &departmentRepository{}
}

func (r *departmentRepository) GetAllDepartments() ([]dto.DepartmentResponse, error) {

	var departments []models.Department
	var departmentResponse []dto.DepartmentResponse

	err := database.DB.
		Where("status = ?", models.DepartmentActive).
		Order("department_name ASC").
		Find(&departments).Error

	if err != nil {
		return nil, err
	}
	for _, department := range departments {
		departmentResponse = append(departmentResponse, dto.DepartmentResponse{
			DepartmentID:   department.DepartmentID,
			DepartmentName: department.DepartmentName,
		})
	}

	return departmentResponse, nil
}