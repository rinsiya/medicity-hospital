package models

type ProfessionalExperience struct {
    ExperienceID uint   `gorm:"primaryKey"`
    DoctorID     uint   `gorm:"not null;index"`
    Details      string `gorm:"type:text;not null"`
}