package service

import (
	"medicity/internal/dto"
	//"medicity/internal/models"
	"medicity/internal/repository"
	// "gorm.io/gorm"
)

type AppointmentService interface {
	//CreateAppointment(tx *gorm.DB, appointment *models.Appointment) error
	//GetAppointmentsByPatientID(patientID uint) ([]models.Appointment, error)
	GetRecentAppointments(patientID uint) ([]dto.RecentAppointment, error)
}

type appointmentService struct {
	appointmentRepo repository.AppointmentRepository
}

func NewAppointmentService(appointmentRepo repository.AppointmentRepository) AppointmentService {
	return &appointmentService{
		appointmentRepo: appointmentRepo,
	}
}

// func (s *appointmentService) CreateAppointment(tx *gorm.DB, appointment *models.Appointment) error
// {}



func (s *appointmentService) GetRecentAppointments(patientID uint) ([]dto.RecentAppointment, error) {

	return s.appointmentRepo.GetLastFiveByPatientID(patientID)

}