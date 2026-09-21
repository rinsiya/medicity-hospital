package service

import (
	"context"
	"fmt"
	"mime/multipart"

	"medicity/internal/models"
	"medicity/internal/repository"
	"medicity/pkg/storage"

	"github.com/google/uuid"
)

type FileService interface {
	GetRecentFilesByUserID(
		userID uint,
	) ([]models.File, error)

	UploadFile(
		ctx context.Context,
		userID uint,
		category models.FileCategory,
		file multipart.File,
		header *multipart.FileHeader,
	) (*models.File, error)
}

type fileService struct {
	fileRepo   repository.FileRepository
	cloudinary *storage.CloudinaryStorage
}

func NewFileService(
	fileRepo repository.FileRepository,
	cloudinary *storage.CloudinaryStorage,
) FileService {

	return &fileService{
		fileRepo:   fileRepo,
		cloudinary: cloudinary,
	}
}

func (s *fileService) GetRecentFilesByUserID(
	userID uint,
) ([]models.File, error) {

	return s.fileRepo.GetRecentFilesByUserID(userID)
}

func (s *fileService) UploadFile(
	ctx context.Context,
	userID uint,
	category models.FileCategory,
	file multipart.File,
	header *multipart.FileHeader,
) (*models.File, error) {

	// Validate category
	if !isValidFileCategory(category) {
		return nil, fmt.Errorf("invalid file category")
	}

	// Validate file
	if err := validateFile(category, header); err != nil {
		return nil, err
	}

	// Generate unique ID
	fileID := uuid.New().String()

	// Generate Cloudinary public ID
	publicID := generatePublicID(
		category,
		userID,
		fileID,
	)

	if publicID == "" {
		return nil, fmt.Errorf("failed to generate file public ID")
	}

	// Upload to Cloudinary

result, err := s.cloudinary.Upload(
	ctx,
	file,
	publicID,
)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// Create database record
	dbFile := &models.File{
		UserID:       userID,
		Category:     category,
		PublicID:     result.PublicID,
		SecureURL:    result.SecureURL,
		FileName:     header.Filename,
		FileType:     header.Header.Get("Content-Type"),
		ResourceType: result.ResourceType,
	}

	// Save metadata in database
	if err := s.fileRepo.Create(dbFile); err != nil {
		return nil, fmt.Errorf("failed to save file metadata: %w", err)
	}

	return dbFile, nil
}

func generatePublicID(
	category models.FileCategory,
	userID uint,
	fileID string,
) string {

	switch category {

	case models.FileCategoryProfilePhoto:

		return fmt.Sprintf(
			"medicity/profile_photos/user_%d_%s",
			userID,
			fileID,
		)

	case models.FileCategoryMedicalReport:

		return fmt.Sprintf(
			"medicity/medical_reports/patient_%d_%s",
			userID,
			fileID,
		)

	case models.FileCategoryDoctorCertificate:

		return fmt.Sprintf(
			"medicity/doctor_certificates/doctor_%d_%s",
			userID,
			fileID,
		)

	case models.FileCategoryIdentityProof:

		return fmt.Sprintf(
			"medicity/identity_proofs/user_%d_%s",
			userID,
			fileID,
		)

	case models.FileCategoryOther:

		return fmt.Sprintf(
			"medicity/other/user_%d_%s",
			userID,
			fileID,
		)

	default:
		return ""
	}
}