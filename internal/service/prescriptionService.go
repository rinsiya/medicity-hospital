package service

import (
	"medicity/internal/dto"
	"medicity/internal/repository"
)

type PrescriptionService interface {
	GetRecentPrescriptions(patientID uint) ([]dto.RecentPrescription, error)
}

type prescriptionService struct {
	prescriptionRepo repository.PrescriptionRepository
}

func NewPrescriptionService(prescriptionRepo repository.PrescriptionRepository) PrescriptionService{
	return &prescriptionService{
		prescriptionRepo: prescriptionRepo,
	}
}

func (s *prescriptionService) GetRecentPrescriptions(patientID uint) ([]dto.RecentPrescription, error) {
	return s.prescriptionRepo.GetRecentPrescriptionsByPatientID(patientID)
}
