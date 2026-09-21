package service

import (
	"medicity/internal/dto"
	"medicity/internal/repository"

	"go.uber.org/zap"
	"medicity/logger"
)

type AppointmentService interface {
	GetRecentAppointments(patientID uint) ([]dto.RecentAppointment, error)

	GetPatientAppointments(
		patientID uint,
		page int,
		pageSize int,
	) ([]dto.PatientAppointment, int64, error)
}

type appointmentService struct {
	appointmentRepo repository.AppointmentRepository
}

func NewAppointmentService(
	appointmentRepo repository.AppointmentRepository,
) AppointmentService {

	return &appointmentService{
		appointmentRepo: appointmentRepo,
	}
}

func (s *appointmentService) GetRecentAppointments(
	patientID uint,
) ([]dto.RecentAppointment, error) {

	return s.appointmentRepo.GetLastFiveByPatientID(patientID)
}

func (s *appointmentService) GetPatientAppointments(
	patientID uint,
	page int,
	pageSize int,
) ([]dto.PatientAppointment, int64, error) {

	// Safety check
	if page < 1 {
		page = 1
	}

	// We want exactly 10 appointments per page
	if pageSize <= 0 {
		pageSize = 10
	}

	appointments, total, err :=
		s.appointmentRepo.GetAppointmentsByPatientID(
			patientID,
			page,
			pageSize,
		)

	if err != nil {

		logger.Log.Error(
			"Failed to fetch patient appointments",
			zap.Uint("patientID", patientID),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.Error(err),
		)

		return nil, 0, err
	}

	return appointments, total, nil
}