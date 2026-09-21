package handler

import (
	"errors"
	"medicity/internal/dto"
	"medicity/internal/models"
	"medicity/internal/service"
	"medicity/logger"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

// func (h *patientHandler) GetMiniProfile(c *gin.Context) {

//     userIDValue, userExists := c.Get("userID")
//     patientIDValue, exists := c.Get("roleID")
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

//     c.JSON(http.StatusOK, gin.H{
//         "success": true,
//         "patient": patient,
//     })
// }
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
