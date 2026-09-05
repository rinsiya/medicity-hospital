package models

import "time"

type AppointmentStatus string

const (
	AppointmentConfirmed AppointmentStatus = "confirmed"
	AppointmentCancelled AppointmentStatus = "cancelled"
	AppointmentCompleted AppointmentStatus = "completed"
	AppointmentMissed    AppointmentStatus = "missed"
)

type Appointment struct {
	AppointmentID uint `gorm:"primaryKey;column:appointment_id"`

	PatientID uint `gorm:"not null;index"`
	DoctorID  uint `gorm:"not null;index"`
	DateTime time.Time `gorm:"not null"`
	SlotID    uint `gorm:"not null;uniqueIndex"`

	// Store the consultation fee charged at the time of booking.
	ConsultationFee int `gorm:"not null"`

	Status AppointmentStatus `gorm:"type:varchar(20);not null;default:'confirmed'"`

	// Google Meet link (generated 10 minutes before the appointment).
	MeetingLink string `gorm:"size:255"`

	// Time until which the meeting link is valid.
	MeetingExpiry *time.Time

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	CompletedAt *time.Time
	
	// Patient Patient `gorm:"foreignKey:PatientID;references:PatientID"`
	// Doctor  Doctor  `gorm:"foreignKey:DoctorID;references:DoctorID"`
	// Slot    DoctorTimeSlot `gorm:"foreignKey:SlotID;references:SlotID"`

	// Payment               *Payment               `gorm:"foreignKey:AppointmentID"`
	// Prescription          *Prescription          `gorm:"foreignKey:AppointmentID"`
	// CancelledAppointment  *CancelledAppointment  `gorm:"foreignKey:AppointmentID"`
	// Refund                *Refund                `gorm:"foreignKey:AppointmentID"`
}