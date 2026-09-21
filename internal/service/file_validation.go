package service

import (
	"fmt"
	"mime/multipart"

	"medicity/internal/models"
)

func isValidFileCategory(category models.FileCategory) bool {

	switch category {

	case models.FileCategoryProfilePhoto,
		models.FileCategoryMedicalReport,
		models.FileCategoryDoctorCertificate,
		models.FileCategoryIdentityProof,
		models.FileCategoryOther:

		return true

	default:
		return false
	}
}

func validateFile(
	category models.FileCategory,
	header *multipart.FileHeader,
) error {

	const maxSize = 10 << 20 // 10 MB

	if header.Size > maxSize {
		return fmt.Errorf("file size must not exceed 10 MB")
	}

	contentType := header.Header.Get("Content-Type")

	switch category {

	case models.FileCategoryProfilePhoto:

		allowedTypes := map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
			"image/webp": true,
		}

		if !allowedTypes[contentType] {
			return fmt.Errorf(
				"profile photo must be JPEG, PNG, or WEBP",
			)
		}

	case models.FileCategoryMedicalReport,
		models.FileCategoryDoctorCertificate,
		models.FileCategoryIdentityProof,
		models.FileCategoryOther:

		allowedTypes := map[string]bool{
			"image/jpeg":      true,
			"image/png":       true,
			"image/webp":      true,
			"application/pdf": true,
		}

		if !allowedTypes[contentType] {
			return fmt.Errorf(
				"file must be an image or PDF",
			)
		}

	default:
		return fmt.Errorf("invalid file category")
	}

	return nil
}