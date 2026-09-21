package service

import (
	"medicity/internal/dto"
	"medicity/internal/repository"
	"medicity/logger"

	"go.uber.org/zap"
)

type PrescriptionService interface {
    GetPatientPrescriptions(patientID uint,page int,pageSize int) ([]dto.PatientPrescription, int64, error)
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


func (s *prescriptionService) GetPatientPrescriptions(patientID uint,page int,pageSize int) ([]dto.PatientPrescription, int64, error) {

	// Safety check
	if page < 1 {
		page = 1
	}

	// We want exactly 10 prescriptions per page
	if pageSize <= 0 {
		pageSize = 10
	}

	prescriptions, total, err :=
		s.prescriptionRepo.GetPrescriptionsByPatientID(patientID,page,pageSize)

	if err != nil {

		logger.Log.Error(
			"Failed to fetch patient prescriptions",
			zap.Uint("patientID", patientID),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.Error(err),
		)

		return nil, 0, err
	}

	return prescriptions, total, nil
}
