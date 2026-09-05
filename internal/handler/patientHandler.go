package handler

import (
	"errors"
	"fmt"
	"medicity/internal/dto"
	"medicity/internal/service"
	"medicity/logger"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type patientHandler struct {
	patientService service.PatientService
	Service service.UserService
}

func NewPatientHandler(patientService service.PatientService,service service.UserService) *patientHandler {
	return &patientHandler{
		patientService: patientService,
	
		Service: service,
	}
}

func (h *patientHandler) PatientSignup(c *gin.Context) {
	c.HTML(http.StatusOK, "patientSignup.html", nil)
}

func (h *patientHandler) ChangePassword(c *gin.Context) {
	c.HTML(http.StatusOK, "change-password.html", gin.H{
		"role": "patient",
	})
}

func (h *patientHandler) CompleteProfile(c *gin.Context) {
	c.HTML(http.StatusOK, "complete-profile.html", nil)
}

func (h *patientHandler) CompleteProfileRequest(c *gin.Context) {

	patientID := c.GetUint("RoleID")
	userID := c.GetUint("userID")

	if patientID == 0 || userID == 0 {
		logger.Log.Error(
			"Unauthorized access to complete profile",
			zap.String("role", "patient"),
			zap.Uint("patientID", patientID),
			zap.Uint("userID", userID),
		)
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
		})
		return
	}

	// Get form values
	dob := c.PostForm("dob")
	gender := c.PostForm("gender")
	weight := c.PostForm("weight")
	height := c.PostForm("height")

	// Validate required fields
	if dob == "" {
		logger.Log.Error(
			"Date of birth is required",
			zap.String("role", "patient"),
			zap.Uint("patientID", patientID),
			zap.Uint("userID", userID),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Date of birth is required",
		})
		return
	}

	if gender == "" {
		logger.Log.Error(
			"Gender is required",
			zap.String("role", "patient"),
			zap.Uint("patientID", patientID),
			zap.Uint("userID", userID),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Gender is required",
		})
		return
	}

	// Profile photo information
	var photoPath string
	var fileName string
	var fileType string

	// Get uploaded photo
	file, err := c.FormFile("photo")

	if err == nil {

		// Validate file type
		fileType = file.Header.Get("Content-Type")

		if fileType != "image/jpeg" &&
			fileType != "image/png" &&
			fileType != "image/webp" {

			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Only JPG, PNG and WEBP images are allowed",
			})
			return
		}

		// Get file extension
		extension := filepath.Ext(file.Filename)

		// Generate unique filename
		fileName = fmt.Sprintf("patient_%d_%d%s", patientID, time.Now().UnixNano(), extension)

		photoPath = filepath.Join("uploads", "patients", fileName)

		// Create upload directory
		err = os.MkdirAll(filepath.Dir(photoPath), 0755)

		if err != nil {

			logger.Log.Error(
				"Failed to create upload directory",
				zap.String("role", "patient"),
				zap.Uint("patientID", patientID),
				zap.String("photoPath", photoPath),
				zap.Error(err),
			)

			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to create upload directory",
			})
			return
		}

		// Save physical file

	}

	// Call service
	err = h.patientService.CompleteProfileRequest(patientID, userID, dob, gender, weight, height, photoPath, fileName, fileType)

	if err != nil {

		logger.Log.Error(
			"Failed to complete patient profile",
			zap.String("role", "patient"),
			zap.Uint("patientID", patientID),
			zap.Uint("userID", userID),
			zap.Error(err),
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	logger.Log.Info(
		"Patient profile completed successfully",
		zap.String("role", "patient"),
		zap.Uint("patientID", patientID),
		zap.Uint("userID", userID),
	)
	err = c.SaveUploadedFile(file, photoPath)

	if err != nil {

		logger.Log.Error(
			"Failed to save profile photo",
			zap.String("role", "patient"),
			zap.Uint("patientID", patientID),
			zap.String("photoPath", photoPath),
			zap.Error(err),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to save profile photo",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  "Profile completed successfully",
		"redirect": "/patient/home",
	})
}

func (h *patientHandler) CheckProfileComplete(c *gin.Context) {

	patientID := c.GetUint("RoleID")
	complete, err := h.patientService.IsProfileComplete(patientID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to check profile",
		})
		return
	}

	if !complete {
		c.JSON(http.StatusOK, gin.H{
			"success":  true,
			"redirect": "/patient/complete-profile",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"redirect": "/patient/schedule-appointment",
	})
}
func (h *patientHandler) Home(c *gin.Context) {

	userID:= c.GetUint("userID")
	patientID:= c.GetUint("RoleID")

	patient, err := h.patientService.GetPatientByID(patientID, userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {

			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Patient not found",
			})
			return
		}
		c.HTML(http.StatusInternalServerError, "patient-home.html", gin.H{
			"error": "Internal server error",
		})
		return
	}
	user,err:=h.Service.GetUserByID(userID)
if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {

			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Patient not found",
			})
			return
		}
		c.HTML(http.StatusInternalServerError, "patient-home.html", gin.H{
			"error": "Internal server error",
		})
		return
	}
	
	c.HTML(http.StatusOK, "patient-home.html", gin.H{
		"patient": patient,
		"user":user,
	})
}


// logout function to clear the JWT cookie and redirect to login page
func (h *patientHandler) Logout(c *gin.Context) {
	// Clear the JWT cookie
	c.SetCookie(
		"access_token",
		"",
		-1, // expire immediately
		"/",
		"",
		false,
		true,
	)

	logger.Log.Info("User logged out")

	c.Redirect(http.StatusSeeOther, "/patient/login")
}

func (h *patientHandler) UpdatePassword(c *gin.Context) {

	var input dto.ChangePasswordInput

	logger.Log.Info("Request to change forgot password",
		zap.String("role", "patient"))
	if err := c.ShouldBind(&input); err != nil {
		logger.Log.Error(
			"password validation failed",
			zap.String("role", "patient"),
			zap.Error(err),
		)
		c.HTML(http.StatusBadRequest, "/patient/change-password.html", gin.H{
			"error": err.Error(),
		})
		if input.Password != input.ConfirmPassword {
			err = errors.New("Password mismatch")
			logger.Log.Error(
				"password mismatch",
				zap.String("role", "patient"),
				zap.Error(err),
			)
			c.HTML(http.StatusBadRequest, "/patient/change-password.html", gin.H{
				"error": err.Error(),
			})
		}
	}

}
