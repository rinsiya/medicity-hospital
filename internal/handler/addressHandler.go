package handler

import (
	"net/http"
	"strings"

	"medicity/internal/service"

	"github.com/gin-gonic/gin"
)

type AddressHandler interface {
	AddAddress(c *gin.Context)
}

type addressHandler struct {
	addressService service.AddressService
	patientService service.PatientService
}

func NewAddressHandler(
	addressService service.AddressService,
	patientService service.PatientService,
) AddressHandler {
	return &addressHandler{
		addressService: addressService,
		patientService: patientService,
	}
}

func (h *addressHandler) AddAddress(c *gin.Context) {

	// Get user ID from JWT middleware
	patientIDValue, exists := c.Get("roleID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	patientID, ok := patientIDValue.(uint)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})
		return
	}

	// Find patient using logged-in user's ID
	patient, err := h.patientService.FindByID(patientID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Patient not found",
		})
		return
	}

	address := strings.TrimSpace(c.PostForm("address"))
	place := strings.TrimSpace(c.PostForm("place"))
	country := strings.TrimSpace(c.PostForm("country"))

	// Since address is optional, allow all fields to be empty.
	if address == "" && place == "" && country == "" {
		c.JSON(http.StatusOK, gin.H{
			"success":  true,
			"message":  "Address skipped successfully",
			"redirect": "/patient/dashboard",
		})
		return
	}

	err = h.addressService.AddOrUpdateAddress(
		patient.PatientID,
		address,
		place,
		country,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to save address",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  "Address saved successfully",
		"redirect": "/patient/home",
	})
}