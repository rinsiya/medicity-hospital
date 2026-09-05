package service

import (
	"medicity/internal/models"
	"medicity/internal/repository"

	//"gorm.io/gorm"
)

type FileService interface {
	//CreateFile(tx *gorm.DB, file *models.File) error
	//GetFilesByPatientID(patientID uint) ([]models.File, error)
	GetRecentFilesByUserID(userID uint) ([]models.File, error)
}
type fileService struct {
	fileRepo repository.FileRepository
}

func NewFileService(fileRepo repository.FileRepository) FileService {
	return &fileService{
		fileRepo: fileRepo,
	}
}

func (s *fileService) GetRecentFilesByUserID(userID uint) ([]models.File, error) {

	return s.fileRepo.GetRecentFilesByUserID(userID)

}