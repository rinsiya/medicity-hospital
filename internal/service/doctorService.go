package service

import (
	"errors"
	"medicity/internal/dto"
	"medicity/internal/models"
	"medicity/internal/repository"
	
)

type DoctorService interface {

    IsProfileComplete(doctorID uint) bool

    GetDoctorByUserID(userID uint) *models.Doctor

    GetDoctorBasicProfile(userID uint) (*dto.DoctorBasicProfileResponse, error)

    SubmitDoctorProfile(userID uint,input dto.DoctorProfileSubmission) error
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
func (s *doctorService) SubmitDoctorProfile(userID uint,input dto.DoctorProfileSubmission) error {

    doctor, err := s.doctorRepo.GetDoctorByUserID(userID)

    if err != nil {
        return err
    }

    if doctor == nil {
        return errors.New("doctor not found")
    }

    return s.doctorRepo.SubmitDoctorProfile(
        doctor.DoctorID,
        input,
    )
}

func (s *doctorService) GetDoctorBasicProfile(
    userID uint,
) (*dto.DoctorBasicProfileResponse, error) {

    return s.doctorRepo.GetDoctorBasicProfile(userID)
}