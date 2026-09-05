package handler

import (
	//"medicity/internal/models"
	"medicity/internal/service"
	//"net/http"

	"github.com/gin-gonic/gin"
)

type AppointmentHandler interface {
	GetRecentAppointments(c *gin.Context)
}

type appointmentHandler struct {
	appointmentService service.AppointmentService
	patientService     service.PatientService
}

func NewAppointmentHandler(appointmentService service.AppointmentService, patientService service.PatientService) AppointmentHandler {
	return &appointmentHandler{
		appointmentService: appointmentService,
		patientService:     patientService,
	}
}

func (h *appointmentHandler) GetRecentAppointments(c *gin.Context) {

	//userIDValue, userExists := c.Get("userID")
PatientID := c.GetUint("RoleID")
	// if !exists || !userExists {
	// 	c.JSON(
	// 		http.StatusUnauthorized,
	// 		gin.H{
	// 			"success": false,
	// 			"message": "User not authenticated",
	// 		},
	// 	)
	// 	return
	// }


	appointments, err := h.appointmentService.GetRecentAppointments(uint(PatientID))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"appointments": appointments})
}