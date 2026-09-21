package handler

import (
	"medicity/internal/service"
	"math"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

type AppointmentHandler interface {
	GetRecentAppointments(c *gin.Context)
	GetPatientAppointments(c *gin.Context)
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
patientID, exists := c.Get("RoleID")
	if !exists{
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"success": false,
				"message": "User not authenticated",
			},
		)
		return
	}


		appointments, err := h.appointmentService.GetRecentAppointments(uint(patientID.(uint)))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch recent appointments",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"appointments": appointments,
	})


if patientID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Patient not authenticated",
		})
		return
	}
}
func (h *appointmentHandler) GetPatientAppointments(c *gin.Context) {

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

	appointments, total, err := h.appointmentService.GetPatientAppointments(
		uint(patientID),
		page,
		pageSize,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch appointments",
		})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"appointments": appointments,
		"pagination": gin.H{
			"current_page": page,
			"page_size":    pageSize,
			"total_items":  total,
			"total_pages":  totalPages,
		},
	})
}