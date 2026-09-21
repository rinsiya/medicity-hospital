package repository

import (
	"medicity/database"
	"medicity/internal/dto"
	"medicity/internal/models"

	"gorm.io/gorm"
)

type AppointmentRepository interface {
	CreateAppointment(tx *gorm.DB, appointment *models.Appointment) error
	GetLastFiveByPatientID(patientID uint) ([]dto.RecentAppointment, error)
		GetAppointmentsByPatientID(patientID uint,page int,pageSize int,) ([]dto.PatientAppointment, int64, error)
}

type appointmentRepository struct{}

func NewAppointmentRepository() AppointmentRepository {
	return &appointmentRepository{}
}

func (r *appointmentRepository) CreateAppointment(tx *gorm.DB, appointment *models.Appointment) error {
	return tx.Create(appointment).Error
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
func (r *appointmentRepository) GetAppointmentsByPatientID(patientID uint, page int, pageSize int) ([]dto.PatientAppointment, int64, error) {

    var appointments []dto.PatientAppointment
    var total int64

    offset := (page - 1) * pageSize

    err := database.DB.
        Table("appointments").
        Where("appointments.patient_id = ?", patientID).
        Count(&total).Error

    if err != nil {
        return nil, 0, err
    }

    err = database.DB.
        Table("appointments").
        Select(`
            appointments.appointment_id,
            appointments.date_time,
            appointments.doctor_id,
            CONCAT(doctors.first_name, ' ', doctors.last_name) AS doctor_name,
            doctors.professional_role AS professional_role,
            departments.department_name AS department_name,
            appointments.status
        `).
        Joins(`
            JOIN doctors
            ON doctors.doctor_id = appointments.doctor_id
        `).
        Joins(`
            LEFT JOIN departments
            ON departments.department_id = doctors.department_id
        `).
        Where("appointments.patient_id = ?", patientID).
        Order("appointments.date_time DESC").
        Limit(pageSize).
        Offset(offset).
        Find(&appointments).Error

    if err != nil {
        return nil, 0, err
    }

    return appointments, total, nil
}