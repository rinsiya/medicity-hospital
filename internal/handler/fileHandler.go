package handler

import (
	"fmt"
	"medicity/internal/models"
	"medicity/internal/service"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FileHandler interface {
	GetRecentMedicalFiles(c *gin.Context)
	UploadFile(c *gin.Context)
}

type fileHandler struct {
	fileService service.FileService
}

func NewFileHandler(fileService service.FileService) FileHandler {
	return &fileHandler{
		fileService: fileService,
	}
}

func (h *fileHandler) GetRecentMedicalFiles(c *gin.Context) {

	userID := c.GetUint("UserID")

	files, err := h.fileService.GetRecentFilesByUserID(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to load medical reports",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    files,
	})
}

func (h *fileHandler) UploadFile(c *gin.Context) {

	// Get logged-in user's ID
	userID := c.GetUint("UserID")

	// Get category from hidden form input
	category := models.FileCategory(c.PostForm("category"))

	// Validate category
	if !isValidCategory(category) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid file category",
		})
		return
	}

	// Get uploaded file
	file, header, err := c.Request.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "File is required",
		})
		return
	}

	defer file.Close()

	// Validate file according to category
	if err := validateFile(category, header); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Send file + category to service
	result, err := h.fileService.UploadFile(
		c.Request.Context(),
		userID,
		category,
		file,
		header,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to upload file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"file":    result,
	})
}

func isValidCategory(category models.FileCategory) bool {

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

	case models.FileCategoryMedicalReport:

		allowedTypes := map[string]bool{
			"image/jpeg":      true,
			"image/png":       true,
			"image/webp":      true,
			"application/pdf": true,
		}

		if !allowedTypes[contentType] {
			return fmt.Errorf(
				"medical report must be an image or PDF",
			)
		}

	case models.FileCategoryDoctorCertificate:

		allowedTypes := map[string]bool{
			"image/jpeg":      true,
			"image/png":       true,
			"image/webp":      true,
			"application/pdf": true,
		}

		if !allowedTypes[contentType] {
			return fmt.Errorf(
				"doctor certificate must be an image or PDF",
			)
		}

	case models.FileCategoryIdentityProof:

		allowedTypes := map[string]bool{
			"image/jpeg":      true,
			"image/png":       true,
			"image/webp":      true,
			"application/pdf": true,
		}

		if !allowedTypes[contentType] {
			return fmt.Errorf(
				"identity proof must be an image or PDF",
			)
		}

	case models.FileCategoryOther:

		allowedTypes := map[string]bool{
			"image/jpeg":      true,
			"image/png":       true,
			"image/webp":      true,
			"application/pdf": true,
		}

		if !allowedTypes[contentType] {
			return fmt.Errorf(
				"unsupported file type",
			)
		}

	default:
		return fmt.Errorf("invalid file category")
	}

	return nil
}