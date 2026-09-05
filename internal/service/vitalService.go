package service

import (
	"medicity/internal/models"
	"medicity/internal/repository"
)

type VitalService interface {
	GetPatientVitals(patientID uint) (*models.VitalData, error)
	SavePatientVitals(patientID uint, vitals interface{}) error
}

type vitalService struct {
	vitalRepo repository.VitalRepository
}

func NewVitalService(
	vitalRepo repository.VitalRepository,
) VitalService {

	return &vitalService{
		vitalRepo: vitalRepo,
	}
}

func (s *vitalService) GetPatientVitals(
	patientID uint,
) (*models.VitalData, error) {

	vitals, err :=
		s.vitalRepo.GetByPatientID(patientID)

	if err != nil {
		return nil, err
	}

	return vitals, nil
}

func (s *vitalService) SavePatientVitals(patientID uint,vitals interface{}) error {

	return s.vitalRepo.SavePatientVitals(patientID,vitals)
}