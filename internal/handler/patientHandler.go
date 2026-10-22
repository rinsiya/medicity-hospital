package handler

import (
	"errors"
	"medicity/internal/dto"
	"medicity/internal/models"
	"medicity/internal/service"
	"medicity/logger"
	"medicity/pkg/utils"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type patientHandler struct {
	patientService service.PatientService
	Service        service.UserService
	fileService    service.FileService
}

func NewPatientHandler(
	patientService service.PatientService,
	service service.UserService,
	fileService service.FileService,
) *patientHandler {
	return &patientHandler{
		patientService: patientService,
		Service:        service,
		fileService:    fileService,
	}
}

func (h *patientHandler) CompleteAddress(c *gin.Context) {
	patientID := c.GetUint("RoleID")
	if patientID == 0 {
		logger.Log.Error(
			"Unauthorized access to complete address",
			zap.String("role", "patient"),
			zap.Uint("patientID", patientID),
		)

		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
		})
		return
	}

	editMode := c.Query("edit") == "true"

	c.HTML(http.StatusOK, "completeAddress.html", gin.H{
		"editMode": editMode,
	})
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
	userID := c.GetUint("UserID")

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
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Date of birth is required",
		})
		return
	}

	if gender == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Gender is required",
		})
		return
	}

	var photoFile multipart.File
	var photoHeader *multipart.FileHeader

	file, header, err := c.Request.FormFile("photo")

	if err == nil {
		photoFile = file
		photoHeader = header

		defer photoFile.Close()
	} else if !errors.Is(err, http.ErrMissingFile) {
		logger.Log.Error(
			"Failed to read profile photo",
			zap.Uint("patientID", patientID),
			zap.Error(err),
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to read profile photo",
		})
		return
	}

	err = h.patientService.CompleteProfileRequest(
		patientID,
		userID,
		dob,
		gender,
		weight,
		height,
	)

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

	// Upload profile photo to Cloudinary if provided

	if photoFile != nil {

		_, err := h.fileService.UploadFile(
			c.Request.Context(),
			userID,
			models.FileCategoryProfilePhoto,
			photoFile,
			photoHeader,
		)

		if err != nil {

			logger.Log.Error(
				"Failed to upload profile photo to Cloudinary",
				zap.String("role", "patient"),
				zap.Uint("patientID", patientID),
				zap.Uint("userID", userID),
				zap.Error(err),
			)

			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Profile completed, but failed to upload profile photo",
			})
			return
		}

		logger.Log.Info(
			"Profile photo uploaded successfully",
			zap.String("role", "patient"),
			zap.Uint("patientID", patientID),
			zap.Uint("userID", userID),
		)
	}

	logger.Log.Info(
		"Patient profile completed successfully",
		zap.String("role", "patient"),
		zap.Uint("patientID", patientID),
		zap.Uint("userID", userID),
	)

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

	logger.Log.Info(
		"Patient home request started",
		zap.String("role", "patient"),
	)

	userIDValue, userExists := c.Get("UserID")
	patientIDValue, patientExists := c.Get("RoleID")

	if !userExists || !patientExists {
		logger.Log.Warn(
			"Patient home access denied: authentication context missing",
			zap.String("role", "patient"),
			zap.Bool("userIDExists", userExists),
			zap.Bool("patientIDExists", patientExists),
		)

		c.HTML(http.StatusUnauthorized, "patientLogin.html", gin.H{
			"success": false,
			"error":   "User not authenticated",
		})
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok || userID == 0 {
		logger.Log.Warn(
			"Invalid user ID in authentication context",
			zap.String("role", "patient"),
		)

		c.HTML(http.StatusUnauthorized, "patientLogin.html", gin.H{
			"success": false,
			"error":   "User not authenticated",
		})
		return
	}

	patientID, ok := patientIDValue.(uint)
	if !ok || patientID == 0 {
		logger.Log.Warn(
			"Invalid patient ID in authentication context",
			zap.Uint("userID", userID),
			zap.String("role", "patient"),  
		)
  
		c.HTML(http.StatusUnauthorized, "patientLogin.html", gin.H{
			"success": false,
			"error":   "Patient not authenticated",
		})
		return
	}

	logger.Log.Info(
		"Loading patient home data",
		zap.String("role", "patient"),
		zap.Uint("userID", userID),
		zap.Uint("patientID", patientID),
	)

	// Get patient details
	patient, err := h.patientService.GetPatientByID(patientID, userID)

	if err != nil {

		logger.Log.Error(
			"Failed to load patient details",
			zap.String("role", "patient"),
			zap.Uint("userID", userID),
			zap.Uint("patientID", patientID),
			zap.Error(err),
		)

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

	logger.Log.Info(
		"Patient details loaded successfully",
		zap.String("role", "patient"),
		zap.Uint("userID", userID),
		zap.Uint("patientID", patientID),
	)

 	// Get user details
	user, err := h.Service.GetUserByID(userID)

	if err != nil {

		logger.Log.Error(
			"Failed to load user details",
			zap.String("role", "patient"),
			zap.Uint("userID", userID),
			zap.Uint("patientID", patientID),
			zap.Error(err),
		)

		if errors.Is(err, service.ErrUserNotFound) {

			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "User not found",
			})
			return
		}

		c.HTML(http.StatusInternalServerError, "patient-home.html", gin.H{
			"error": "Internal server error",
		})
		return
	}

	logger.Log.Info(
		"User contact details loaded successfully",
		zap.String("role", "patient"),
		zap.Uint("userID", userID),
		zap.Uint("patientID", patientID),
	)

	logger.Log.Info(
		"Patient home loaded successfully",
		zap.String("role", "patient"),
		zap.Uint("userID", userID),
		zap.Uint("patientID", patientID),
	)

	c.HTML(http.StatusOK, "patient-home.html", gin.H{
		"patient": patient,
		"user":    user,
	})
}

// func (h *patientHandler) GetMiniProfile(c *gin.Context) {

//     userIDValue, userExists := c.Get("UserID")
//     patientIDValue, exists := c.Get("RoleID")
//     if !exists || !userExists {
//         c.JSON(http.StatusUnauthorized, gin.H{
//             "success": false,
//             "message": "User not authenticated",
//         })
//         return
//     }

//     userID := userIDValue.(uint)
//     patientID := patientIDValue.(uint)

//     patient, err := h.patientService.GetPatientMiniProfile(userID,patientID)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{
//             "success": false,
//             "message": "Failed to load patient profile",
//         })
//         return
//     }

//	    c.JSON(http.StatusOK, gin.H{
//	        "success": true,
//	        "patient": patient,
//	    })
//	}
func (h *patientHandler) Logout(c *gin.Context) {
	utils.Logout(c, "/patient/login")
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

func (h *patientHandler) GetAddress(c *gin.Context) {

	logger.Log.Info(
		"Patient address request started",
		zap.String("role", "patient"),
	)

	patientIDValue, exists := c.Get("RoleID")

	if !exists {
		logger.Log.Warn(
			"Patient address access denied: patient ID missing",
			zap.String("role", "patient"),
		)

		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Patient not authenticated",
		})
		return
	}

	patientID, ok := patientIDValue.(uint)

	if !ok || patientID == 0 {
		logger.Log.Warn(
			"Patient address access denied: invalid patient ID",
			zap.String("role", "patient"),
		)

		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Patient not authenticated",
		})
		return
	}

	logger.Log.Info(
		"Loading patient address",
		zap.String("role", "patient"),
		zap.Uint("patientID", patientID),
	)

	address, err := h.patientService.GetAddressByPatientID(patientID)

	if err != nil {

		// No address is not an error
		if errors.Is(err, gorm.ErrRecordNotFound) {

			logger.Log.Info(
				"Patient address not found",
				zap.String("role", "patient"),
				zap.Uint("patientID", patientID),
			)

			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"exists":  false,
				"address": nil,
			})
			return
		}

		logger.Log.Error(
			"Failed to load patient address",
			zap.String("role", "patient"),
			zap.Uint("patientID", patientID),
			zap.Error(err),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to load address",
		})
		return
	}

	logger.Log.Info(
		"Patient address loaded successfully",
		zap.String("role", "patient"),
		zap.Uint("patientID", patientID),
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"exists":  true,
		"address": address,
	})
}
