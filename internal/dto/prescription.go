package dto

import "time"

type RecentPrescription struct {
	PrescriptionID uint      `json:"prescription_id"`
	CreatedAt      time.Time `json:"created_at"`
	DoctorName     string    `json:"doctor_name"`
	ProfessionalRole string    `json:"professional_role"`
}