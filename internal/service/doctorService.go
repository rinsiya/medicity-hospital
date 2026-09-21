package service

import (
	"medicity/internal/models"
	"medicity/internal/repository"
)

type DoctorService interface {
	IsProfileComplete(doctorID uint) bool
	GetDoctorByUserID(userID uint) *models.Doctor
}

type doctorService struct {
	doctorRepo repository.DoctorRepository
}

func NewDoctorService(doctorRepo repository.DoctorRepository) DoctorService {
	return &doctorService{
		doctorRepo: doctorRepo,
	}
}

func (s *doctorService) IsProfileComplete(doctorID uint) bool {
	isComplete, err := s.doctorRepo.IsProfileComplete(doctorID)

	if err != nil {
		return false
	}

	return isComplete;
}

func (s *doctorService) GetDoctorByUserID(userID uint) *models.Doctor {
	doctor, err := s.doctorRepo.GetDoctorByUserID(userID)

	if err != nil {
		return nil
	}

	return doctor
}