package service

import (
	//"errors"
	"fmt"
	"medicity/internal/models"
	"medicity/internal/repository"
	"strconv"
	"strings"
	"time"
)

type PatientService interface {
	GetPatientByID(patientID,userID uint) (*models.Patient, error)
	IsProfileComplete(patientID uint) (bool, error)
	CompleteProfileRequest(patientID uint,userID uint, dob string, gender string, weight string, height string,photoPath string,fileName string,fileType string) error
}

type patientService struct {
	patientRepo        repository.PatientRepository
}

func NewPatientService(patientRepo repository.PatientRepository) PatientService {
	return &patientService{
		patientRepo: patientRepo,
	}

}

func (s *patientService) CompleteProfileRequest(patientID uint,userID uint,	dob string,gender string,weight string,height string,photoPath string,fileName string,fileType string) error {

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

	// Create file record only when photo is uploaded
	var fileRecord *models.File

	if photoPath != "" {

		fileRecord = &models.File{
			UserID:      userID,
			Category:    models.FileCategory("profile_photo"),
			StoragePath: photoPath,
			FileName:    fileName,
			FileType:    fileType,
			Remarks:     "Patient profile photo",
		}
	}

	// Save patient profile and file information
	return s.patientRepo.CompleteProfileRequest(
		patientID,
		userID,
		dateOfBirth,
		gender,
		patientWeight,
		patientHeight,
		fileRecord,
	)
}
func (s *patientService) IsProfileComplete(patientID uint) (bool, error) {
	return s.patientRepo.IsProfileComplete(patientID)
}
func (s *patientService) GetPatientByID(patientID ,userID uint) (*models.Patient, error) {
	patient, err := s.patientRepo.GetPatientByID(patientID,userID)
	if err != nil {
		return nil, err
	}
	return patient, nil
}