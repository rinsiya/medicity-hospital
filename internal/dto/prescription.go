package dto

import "time"

type RecentPrescription struct {
	PrescriptionID uint      `json:"prescription_id"`
	CreatedAt      time.Time `json:"created_at"`
	DoctorName     string    `json:"doctor_name"`
	ProfessionalRole string    `json:"professional_role"`
}

type PatientPrescription struct {
	AppointmentID    uint      `json:"appointment_id"`
	CreatedAt        time.Time `json:"created_at"`
	Diagnosis        string    `json:"diagnosis"`
	DoctorID         uint      `json:"doctor_id"`
	DoctorName       string    `json:"doctor_name"`
	ProfessionalRole string    `json:"professional_role"`
	DepartmentName   string    `json:"department_name"`
}