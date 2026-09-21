package repository

import (
	"medicity/database"
	"medicity/internal/dto"
	"medicity/internal/models"
)

type FileRepository interface {
	Create(file *models.File) error

	FindByUserID(userID uint) ([]models.File, error)

	GetRecentFilesByUserID(userID uint) ([]models.File, error)

	GetFilesByUserID(
		userID uint,
		page int,
		pageSize int,
	) ([]dto.PatientFiles, int64, error)

	GetProfilePhotoByUserID(userID uint) (*models.File, error)
}

type fileRepository struct{}

func NewFileRepository() FileRepository {
	return &fileRepository{}
}
func (r *fileRepository) GetProfilePhotoByUserID(
	userID uint,
) (*models.File, error) {

	var file models.File

	err := database.DB.
		Where(
			"user_id = ? AND category = ?",
			userID,
			models.FileCategoryProfilePhoto,
		).
		Order("uploaded_at DESC").
		First(&file).Error

	if err != nil {
		return nil, err
	}

	return &file, nil
}
func (r *fileRepository) GetRecentFilesByUserID(
	userID uint,
) ([]models.File, error) {

	var files []models.File

	err := database.DB.
		Where(
			"user_id = ? AND category = ?",
			userID,
			models.FileCategoryMedicalReport,
		).Order("uploaded_at DESC").
		Limit(5).
		Find(&files).Error

	if err != nil {
		return nil, err
	}

	return files, nil
}

func (r *fileRepository) Create(
	file *models.File,
) error {

	return database.DB.Create(file).Error
}

func (r *fileRepository) FindByUserID(
	userID uint,
) ([]models.File, error) {

	var files []models.File

	err := database.DB.
		Where("user_id = ?", userID).
		Order("uploaded_at DESC").
		Find(&files).Error

	return files, err
}

func (r *fileRepository) GetFilesByUserID(
	userID uint,
	page int,
	pageSize int,
) ([]dto.PatientFiles, int64, error) {

	var files []dto.PatientFiles
	var total int64

	offset := (page - 1) * pageSize

	// Count files belonging to this user
	err := database.DB.
		Model(&models.File{}).
		Where("user_id = ?", userID).
		Count(&total).Error

	if err != nil {
		return nil, 0, err
	}

	// Fetch paginated files
	err = database.DB.
		Table("files").
		Select(`
			files.uploaded_at,
			files.file_name,
			files.category,
			files.file_type
		`).
		Where("files.user_id = ?", userID).
		Order("files.uploaded_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&files).Error

	if err != nil {
		return nil, 0, err
	}

	return files, total, nil
}