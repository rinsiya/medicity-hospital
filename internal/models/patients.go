package models

import "time"

type Patient struct {
	PatientID uint `gorm:"primaryKey;column:patient_id"`

	UserID uint `gorm:"not null;uniqueIndex"`

	FirstName string `gorm:"size:50;not null"`

	LastName string `gorm:"size:50;not null"`

	Gender string `gorm:"size:10;"`

	DOB time.Time `gorm:"type:date;"`

	Height float32 `gorm:"type:float;"`

	Weight float32 `gorm:"type:float;"`

	ProfilePhotoID *uint

	ProfilePhoto *File `gorm:"foreignKey:ProfilePhotoID;references:FileID"`

	CreatedAt time.Time `gorm:"autoCreateTime"`

	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
