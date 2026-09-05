package repository

import (
	"fmt"
	"medicity/database"
	"medicity/internal/models"
	"time"

	"gorm.io/gorm"
)

type PatientRepository interface {
	//FindByEmailOrPhone(username string) (*models.User, error)
	FindByID(id uint) (*models.Patient, error)
	CreatePatient(tx *gorm.DB, patient *models.Patient) error
	EmailExists(email string) (bool, error)
	PhoneExists(phone string) (bool, error)
	IsProfileComplete(patientID uint) (bool, error)
	GetPatientByID(PatientID,UserID uint)(*models.Patient, error)
	CompleteProfileRequest(patientID uint, userID uint, dob time.Time, gender string, weight *float64, height *float64, fileRecord *models.File) error
}

type patientRepository struct{}

func NewPatientRepository() PatientRepository {
	return &patientRepository{}
}

func (r *patientRepository) CompleteProfileRequest(patientID uint,userID uint,
	dob time.Time,
	gender string,
	weight *float64,
	height *float64,
	fileRecord *models.File,
) error {

	return database.DB.Transaction(func(tx *gorm.DB) error {

		// Store Profile Photo ID
		var profilePhotoID *uint

		// Create file record if photo exists
		if fileRecord != nil {

			if err := tx.Create(fileRecord).Error; err != nil {
				return err
			}

			profilePhotoID = &fileRecord.FileID
		}

		// Patient fields to update
		updates := map[string]interface{}{
			"dob":    dob,
			"gender": gender,
		}

		if weight != nil {
			updates["weight"] = *weight
		}

		if height != nil {
			updates["height"] = *height
		}

		if profilePhotoID != nil {
			updates["profile_photo_id"] = *profilePhotoID
		}

		// Update patient
		result := tx.
			Model(&models.Patient{}).
			Where("patient_id = ? AND user_id = ?", patientID, userID).
			Updates(updates)

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("patient not found")
		}

		return nil
	})
}
func (r *patientRepository) EmailExists(email string) (bool, error) {
	var count int64

	err := database.DB.
		Model(&models.User{}).
		Where("email = ?", email).
		Count(&count).
		Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *patientRepository) PhoneExists(phone string) (bool, error) {
	var count int64

	err := database.DB.
		Model(&models.User{}).
		Where("phone = ?", phone).
		Count(&count).
		Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *patientRepository) CreatePatient(tx *gorm.DB, patient *models.Patient) error {
	// Create patient
	return tx.Create(patient).Error

}
func (r *patientRepository) FindByID(UserId uint) (*models.Patient, error) {
	var patient models.Patient
	err := database.DB.Where("user_id = ?", UserId).First(&patient).Error
	if err != nil {
		return nil, err
	}
	return &patient, nil
}
func (r *patientRepository) GetPatientByID(patientID uint,userID uint) (*models.Patient, error) {

	var patient models.Patient

	err := database.DB.
		Preload("ProfilePhoto").
		Where(
			"patient_id = ? AND user_id = ?",
			patientID,
			userID,
		).
		First(&patient).Error

	if err != nil {
		return nil, err
	}

	return &patient, nil
}

func (r *patientRepository) IsProfileComplete(patientID uint) (bool, error) {

	var patient models.Patient

	err := database.DB.
		Select("first_name", "last_name", "gender", "dob").
		Where("patient_id = ?", patientID).
		First(&patient).Error

	if err != nil {
		return false, err
	}

	if patient.Gender == "" ||
		patient.DOB.IsZero() {

		return false, nil
	}

	return true, nil
}
