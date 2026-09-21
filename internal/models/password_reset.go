package models

import "time"

type PasswordReset struct {
	PasswordResetID uint      `gorm:"primaryKey;column:password_reset_id"`
	Phone           string    `gorm:"size:15;not null;index"`
	SessionID        string    `gorm:"size:255;not null"`
	OTPExpiresAt    time.Time `gorm:"not null"`
	Verified        bool      `gorm:"not null;default:false"`

	CreatedAt       time.Time `gorm:"autoCreateTime"`
}