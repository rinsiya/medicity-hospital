package models

import "time"

type DoctorProfile struct {
	ProfileID uint `gorm:"primaryKey;column:profile_id"`

	DoctorID uint `gorm:"not null;uniqueIndex"`

	HospitalAddress string `gorm:"size:255"`

	
	Experience uint

	ProfessionalSummary string `gorm:"type:text"`

	Biography string `gorm:"type:text"`

	ProfilePhotoID *uint

	//ProfilePhoto *File `gorm:"foreignKey:ProfilePhotoID;references:FileID"`

	CreatedAt time.Time `gorm:"autoCreateTime"`

	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	//Doctor Doctor `gorm:"foreignKey:DoctorID;references:DoctorID"`
}