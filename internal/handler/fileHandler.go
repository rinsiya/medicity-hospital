package handler

import (
	"medicity/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)
type FileHandler interface{
GetRecentMedicalFiles(c *gin.Context)
}

type fileHandler struct{
fileService service.FileService
}

func NewFileHandler(fileService service.FileService) FileHandler{
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