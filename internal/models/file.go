package models

import "time"

type FileCategory string

const (
	FileCategoryDoctorCertificate FileCategory = "doctor_certificate"
	FileCategoryMedicalReport     FileCategory = "medical_report"
	FileCategoryIdentityProof     FileCategory = "identity_proof"
	FileCategoryProfilePhoto      FileCategory = "profile_photo"
	FileCategoryOther             FileCategory = "other"
)

type File struct {
	FileID uint `gorm:"primaryKey;column:file_id"`

	UserID uint `gorm:"not null;index"`

	Category FileCategory `gorm:"type:varchar(30);not null"`

	PublicID     string `gorm:"not null;uniqueIndex"`
	SecureURL    string `gorm:"type:text"`
	FileName     string `gorm:"size:255;not null"`
	FileType     string `gorm:"size:100;not null"`
	ResourceType string

	Remarks string `gorm:"type:text"`

	UploadedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time
}