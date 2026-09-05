package repository

import (
	"medicity/database"
	"medicity/internal/dto"
	//"medicity/internal/models"
	//"gorm.io/gorm"
)

type PrescriptionRepository interface {
	//CreatePrescription(tx *gorm.DB, prescription *models.Prescription) error
	GetRecentPrescriptionsByPatientID(patientID uint) ([]dto.RecentPrescription, error)
}

type prescriptionRepository struct{}

func NewPrescriptionRepository() PrescriptionRepository {
	return &prescriptionRepository{}
}
func (r *prescriptionRepository) GetRecentPrescriptionsByPatientID(patientID uint) ([]dto.RecentPrescription, error) {

	var prescriptions []dto.RecentPrescription

	err := database.DB.
		Table("prescriptions").
		Select(`
			prescriptions.prescription_id,
			prescriptions.created_at,
			CONCAT(doctors.first_name, ' ', doctors.last_name) AS doctor_name,
			doctors.professional_role AS professional_role
		`).
		Joins(`
			JOIN appointments
			ON appointments.appointment_id = prescriptions.appointment_id
		`).
		Joins(`
			JOIN doctors
			ON doctors.doctor_id = appointments.doctor_id
		`).
		Where("appointments.patient_id = ?", patientID).
		Order("prescriptions.created_at DESC").
		Limit(5).
		Scan(&prescriptions).Error

	if err != nil {
		return nil, err
	}

	return prescriptions, nil
}