package handler

import (
	"math"
	"medicity/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PrescriptionHandler interface {
	GetRecentPrescriptions(c *gin.Context)
	GetPatientPrescriptions(c *gin.Context)
}

type prescriptionHandler struct {
	prescriptionService service.PrescriptionService
}

func NewPrescriptionHandler(prescriptionService service.PrescriptionService) PrescriptionHandler{
	return &prescriptionHandler{
		prescriptionService: prescriptionService,
	}
}

func (h *prescriptionHandler) GetRecentPrescriptions(c *gin.Context) {

	patientID := c.GetUint("RoleID")

	prescriptions, err :=
		h.prescriptionService.GetRecentPrescriptions(patientID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to load recent prescriptions",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    prescriptions,
	})
}

func (h *prescriptionHandler) GetPatientPrescriptions(c *gin.Context){

	patientID := c.GetUint("RoleID")
if patientID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Patient not authenticated",
		})
		return
	}
	page := 1
	if pageParam := c.Query("page"); pageParam != "" {
		parsedPage, err := strconv.Atoi(pageParam)
		if err != nil || parsedPage < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid page number",
			})
			return
		}
		page = parsedPage
	}

	const pageSize = 10

	prescriptions, total, err := h.prescriptionService.GetPatientPrescriptions(
		uint(patientID),
		page,
		pageSize,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch prescriptions",
		})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"prescriptions": prescriptions,
		"pagination": gin.H{
			"current_page": page,
			"page_size":    pageSize,
			"total_items":  total,
			"total_pages":  totalPages,
		},
	})

}