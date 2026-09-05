package repository

import (
	"medicity/database"
	"medicity/internal/dto"
	"medicity/internal/models"

	"gorm.io/gorm"
)

type AppointmentRepository interface {
	CreateAppointment(tx *gorm.DB, appointment *models.Appointment) error
	GetAppointmentsByPatientID(patientID uint) ([]models.Appointment, error)
	GetLastFiveByPatientID(patientID uint) ([]dto.RecentAppointment, error)
}

type appointmentRepository struct{}

func NewAppointmentRepository() AppointmentRepository {
	return &appointmentRepository{}
}

func (r *appointmentRepository) CreateAppointment(tx *gorm.DB, appointment *models.Appointment) error {
	return tx.Create(appointment).Error
}

func (r *appointmentRepository) GetAppointmentsByPatientID(patientID uint) ([]models.Appointment, error) {
	var appointments []models.Appointment

	err := database.DB.Where("patient_id = ?", patientID).Find(&appointments).Error
	if err != nil {
		return nil, err
	}

	return appointments, nil
}		
func (r *appointmentRepository) GetLastFiveByPatientID(patientID uint) ([]dto.RecentAppointment, error) {

	var appointments []dto.RecentAppointment

err := database.DB.
Table("appointments").
Select(`
appointments.appointment_id,
appointments.date_time,
CONCAT(doctors.first_name, ' ', doctors.last_name) AS doctor_name,
doctors.professional_role AS professional_role,
appointments.status AS status
`).
	Joins(`
		JOIN doctors
		ON doctors.doctor_id = appointments.doctor_id
	`).
	Where("patient_id = ?", patientID).
	Order("date_time DESC").
	Limit(5).
	Find(&appointments).Error
	if err != nil {
		return nil, err
	}

	return appointments, nil
}


