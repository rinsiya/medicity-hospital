package models

import "time"

type PendingUserSignup struct {
	PendingID uint `gorm:"primaryKey;column:pending_id"`
	Role UserRole `gorm:"type:varchar(20);not null"`

	FirstName string    `gorm:"size:50;not null"`
	LastName  string    `gorm:"size:50;not null"`
	
	Email string `gorm:"size:150;not null;index"`
	Phone string `gorm:"size:20;not null;index"`

	Password string `gorm:"size:255;not null"`

	SessionID string `gorm:"size:255;not null"`

	OTPExpiresAt time.Time `gorm:"not null"`

	CreatedAt time.Time
}