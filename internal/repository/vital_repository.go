package repository

import (
	"encoding/json"
	"fmt"
	"medicity/database"
	"medicity/internal/models"
	"strconv"
	"strings"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)
type Vital struct {
    Name  string      `json:"name"`
    Value interface{} `json:"value"`
    Unit  string      `json:"unit"`
}
type VitalRepository interface {
	GetByPatientID(patientID uint) (*models.VitalData, error)
	SavePatientVitals(patientID uint, vitals interface{}) error
}

type vitalRepository struct {
}

func NewVitalRepository() VitalRepository {
	return &vitalRepository{}
}

func (r *vitalRepository) GetByPatientID(patientID uint) (*models.VitalData, error) {

	var vitalData models.VitalData

	err := database.DB.Where("patient_id = ?", patientID).Order("created_at DESC").First(&vitalData).Error

	if err != nil {
		return nil, err
	}

	return &vitalData, nil
}


func (r *vitalRepository) SavePatientVitals(patientID uint, vitals interface{}) error {

	// Convert vitals into JSON
	vitalJSON, err := json.Marshal(vitals)
	if err != nil {
		return err
	}

	// Convert JSON into Vital list
	var vitalList []Vital

	if err := json.Unmarshal(vitalJSON, &vitalList); err != nil {
		return err
	}

	// Store weight if it is present
	var weight *float64

	for _, vital := range vitalList {

		if strings.EqualFold(vital.Name, "Weight") {

		
				parsedWeight, err := strconv.ParseFloat(vital.Value.(string), 64)
				if err != nil {
					return fmt.Errorf("invalid weight value: %s", vital.Value)
				}

				weight = &parsedWeight

		
	}}

	// Start transaction
	return database.DB.Transaction(func(tx *gorm.DB) error {

		// Create NEW vital history record
		vitalData := models.VitalData{
			PatientID: patientID,
			Vitals:    datatypes.JSON(vitalJSON),
		}

		if err := tx.Create(&vitalData).Error; err != nil {
			return err
		}

		// Update patient's current weight if Weight was provided
		if weight != nil {

			result := tx.
				Model(&models.Patient{}).
				Where("patient_id = ?", patientID).
				Update("weight", *weight)

			if result.Error != nil {
				return result.Error
			}

			if result.RowsAffected == 0 {
				return fmt.Errorf("patient not found")
			}
		}

		return nil
	})
}