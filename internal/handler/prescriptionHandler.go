package handler

import (
	"medicity/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PrescriptionHandler interface {
	GetRecentPrescriptions(c *gin.Context)
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