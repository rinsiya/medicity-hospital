package models


type VerificationStatus string

const (
	PendingVerification VerificationStatus = "pending"
	Verified            VerificationStatus = "verified"
	Rejected            VerificationStatus = "rejected"
)

type Doctor struct {
	DoctorID uint `gorm:"primaryKey;column:doctor_id"`

	UserID uint `gorm:"not null;uniqueIndex"`

	FirstName string `gorm:"size:50;not null"`

	LastName string `gorm:"size:50;not null"`

	ProfessionalRole string `gorm:"size:100"`

	DepartmentID *uint `gorm:"index"`

	ConsultationFee int 

	VerificationStatus VerificationStatus `gorm:"type:varchar(20);not null;default:'pending'"`

	}