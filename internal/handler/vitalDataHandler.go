package handler

import (
	"medicity/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"medicity/internal/dto"
)

type VitalHandler struct {
	patientService service.PatientService
	vitalService   service.VitalService
}

func NewVitalHandler(
	patientService service.PatientService,
	vitalService service.VitalService,
) *VitalHandler {

	return &VitalHandler{
		patientService: patientService,
		vitalService:   vitalService,
	}
}

func (h *VitalHandler) GetPatientVitals(c *gin.Context) {

	//userIDValue, userExists :=c.Get("userID")
	PatientID := c.GetUint("RoleID")


   vitalData, err :=h.vitalService.GetPatientVitals(PatientID)

	if err != nil {

		c.JSON(
			http.StatusNotFound,
			gin.H{
				"success": true,
				"data":    []interface{}{},
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"success": true,
			"data":    vitalData.Vitals,
		},
	)
}
func (h *VitalHandler) SavePatientVitals(c *gin.Context) {

	PatientID:= c.GetUint("RoleID")
	var request dto.SavePatientVitalsRequest

	if err := c.ShouldBindJSON(&request); err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"success": false,
				"message": "Invalid vital data",
			},
		)

		return
	}

	if len(request.Vitals) == 0 {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"success": false,
				"message": "At least one vital is required",
			},
		)

		return
	}

	err := h.vitalService.SavePatientVitals(
		PatientID,
		request.Vitals,
	)

	if err != nil {

		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"success": false,
				"message": "Failed to save vital data",
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"success": true,
			"message": "Vital data saved successfully",
		},
	)
}