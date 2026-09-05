package repository

import (
	"medicity/database"
	"medicity/internal/models"

	//"gorm.io/gorm"
)

type FileRepository interface {
	//CreateAppointment(tx *gorm.DB, appointment *models.Appointment) error
	GetRecentFilesByUserID(userID uint) ([]models.File, error)
}

type fileRepository struct{}

func NewFileRepository() FileRepository{
	return &fileRepository{}
}

func (r *fileRepository) GetRecentFilesByUserID(userID uint) ([]models.File, error) {

	var files []models.File

	err := database.DB.
		Where("user_id = ?", userID).
		Order("uploaded_at DESC").
		Limit(5).
		Find(&files).Error

	if err != nil {
		return nil, err
	}

	return files, nil
}