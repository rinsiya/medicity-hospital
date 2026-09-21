package service

import (
	"fmt"
	"medicity/internal/dto"
	"medicity/internal/models"
	"medicity/internal/repository"
	"strconv"
	"strings"
	"time"
)

type PatientService interface {
	GetPatientByID(patientID, userID uint) (*dto.PatientProfile, error)
	IsProfileComplete(patientID uint) (bool, error)
	FindByID(patientID uint) (*models.Patient, error)

	CompleteProfileRequest(patientID uint,userID uint,dob string,gender string,weight string,height string) error
}

type patientService struct {
	patientRepo repository.PatientRepository
	fileRepo repository.FileRepository
}

func NewPatientService(patientRepo repository.PatientRepository,
	fileRepo repository.FileRepository,
	) PatientService {
	return &patientService{
		patientRepo: patientRepo,
		fileRepo: fileRepo,
	}
}

func (s *patientService) CompleteProfileRequest(
	patientID uint,
	userID uint,
	dob string,
	gender string,
	weight string,
	height string,
) error {

	// Parse date of birth
	dateOfBirth, err := time.Parse("2006-01-02", dob)
	if err != nil {
		return fmt.Errorf("invalid date of birth")
	}

	// Validate gender
	gender = strings.TrimSpace(gender)

	if !strings.EqualFold(gender, "Male") &&
		!strings.EqualFold(gender, "Female") &&
		!strings.EqualFold(gender, "Other") {
		return fmt.Errorf("invalid gender")
	}

	// Parse weight
	var patientWeight *float64

	if weight != "" {
		value, err := strconv.ParseFloat(weight, 64)
		if err != nil {
			return fmt.Errorf("invalid weight")
		}

		if value <= 0 {
			return fmt.Errorf("weight must be greater than zero")
		}

		patientWeight = &value
	}

	// Parse height
	var patientHeight *float64

	if height != "" {
		value, err := strconv.ParseFloat(height, 64)
		if err != nil {
			return fmt.Errorf("invalid height")
		}

		if value <= 0 {
			return fmt.Errorf("height must be greater than zero")
		}

		patientHeight = &value
	}

	// Save patient profile.
	// Profile photo is handled by FileService.
	return s.patientRepo.CompleteProfileRequest(
		patientID,
		userID,
		dateOfBirth,
		gender,
		patientWeight,
		patientHeight,
	)
}

func (s *patientService) IsProfileComplete(patientID uint) (bool, error) {
	return s.patientRepo.IsProfileComplete(patientID)
}

func (s *patientService) GetPatientByID(patientID, userID uint) (*dto.PatientProfile, error) {
	patient, err := s.patientRepo.GetPatientByID(patientID, userID)
	if err != nil {
		return nil, err
	}
	profile := &dto.PatientProfile{
		PatientName: patient.FirstName + " " + patient.LastName,
		Gender:      patient.Gender,
		DOB:         patient.DOB,
		Height:      patient.Height,
		Weight:      patient.Weight,
		ProfilePhoto: "",
	}
	file, err := s.fileRepo.GetProfilePhotoByUserID(userID)

	if err == nil && file != nil {
		profile.ProfilePhoto = file.SecureURL
	}

	return profile, nil
	
}

func (s *patientService) FindByID(patientID uint) (*models.Patient, error) {
	patient, err := s.patientRepo.FindByID(patientID)
	if err != nil {
		return nil, err
	}

	return patient, nil
}