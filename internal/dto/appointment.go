package dto

import "time"

type RecentAppointment struct {
	AppointmentID    uint      `json:"appointment_id"`
	DateTime         time.Time `json:"date_time"`
	DoctorName       string    `json:"doctor_name"`
	ProfessionalRole string    `json:"professional_role"`
	Status           string    `json:"status"`
}

type PatientAppointment struct {
	AppointmentID    uint      `json:"appointment_id"`
	DateTime         time.Time `json:"date_time"`
	DoctorID         uint      `json:"doctor_id"`
	DoctorName       string    `json:"doctor_name"`
	ProfessionalRole string    `json:"professional_role"`
	DepartmentName   string    `json:"department_name"`
	Status           string    `json:"status"`
}
