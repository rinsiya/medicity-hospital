package handler

import (
	"encoding/json"
	"medicity/internal/dto"
	"medicity/internal/models"
	"medicity/internal/service"
	"medicity/logger"
	"medicity/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type doctorHandler struct {
	doctorService service.DoctorService
	fileService   service.FileService
}

func NewDoctorHandler(
	doctorService service.DoctorService,
	fileService service.FileService,
) *doctorHandler {

	logger.Log.Info("Doctor handler initialized")

	return &doctorHandler{
		doctorService: doctorService,
		fileService:   fileService,
	}
}

// --------------------------------------------------
// Doctor Signup
// --------------------------------------------------

func (h *doctorHandler) DoctorSignup(c *gin.Context) {

	logger.Log.Info("Doctor signup page requested")

	c.HTML(
		http.StatusOK,
		"doctorSignup.html",
		nil,
	)
}

// --------------------------------------------------
// Doctor Dashboard
// --------------------------------------------------

func (h *doctorHandler) Dashboard(c *gin.Context) {

	doctorID := c.GetUint("RoleID")

	logger.Log.Info(
		"Doctor dashboard requested",
		zap.Uint("doctor_id", doctorID),
	)

	isProfileComplete := h.doctorService.IsProfileComplete(doctorID)

	logger.Log.Info(
		"Doctor profile completion status checked",
		zap.Uint("doctor_id", doctorID),
		zap.Bool("profile_complete", isProfileComplete),
	)

	if isProfileComplete {

		logger.Log.Info(
			"Showing doctor verification page",
			zap.Uint("doctor_id", doctorID),
		)

		c.HTML(
			http.StatusOK,
			"doctorForVerification.html",
			nil,
		)

	} else {

		logger.Log.Info(
			"Showing doctor profile completion page",
			zap.Uint("doctor_id", doctorID),
		)

		c.HTML(
			http.StatusOK,
			"doctorCompleteProfile.html",
			nil,
		)
	}
}

// --------------------------------------------------
// Change Password
// --------------------------------------------------

func (h *doctorHandler) ChangePassword(c *gin.Context) {

	logger.Log.Info("Doctor change password page requested")

	c.HTML(
		http.StatusOK,
		"changePassword.html",
		nil,
	)
}

// --------------------------------------------------
// Logout
// --------------------------------------------------

func (h *doctorHandler) Logout(c *gin.Context) {

	logger.Log.Info("Doctor logout requested")

	utils.Logout(c, "/doctor/login")

	logger.Log.Info("Doctor logged out successfully")
}

// --------------------------------------------------
// Update Password
// --------------------------------------------------

func (h *doctorHandler) UpdatePassword(c *gin.Context) {

	logger.Log.Info("Doctor password update requested")

	// Password update logic will be implemented here.
}

// --------------------------------------------------
// Complete Profile
// --------------------------------------------------

func (h *doctorHandler) CompleteProfile(c *gin.Context) {

	logger.Log.Info("Doctor complete profile page requested")

	c.HTML(
		http.StatusOK,
		"doctorCompleteProfile.html",
		nil,
	)
}

// --------------------------------------------------
// Verification Pending
// --------------------------------------------------

func (h *doctorHandler) VerificationPending(c *gin.Context) {

	logger.Log.Info("Doctor verification pending page requested")

	c.HTML(
		http.StatusOK,
		"doctorWaitingForVerification.html",
		nil,
	)
}

// --------------------------------------------------
// Profile Rejected
// --------------------------------------------------

func (h *doctorHandler) ProfileRejected(c *gin.Context) {

	logger.Log.Info("Doctor profile rejected page requested")

	c.HTML(
		http.StatusOK,
		"profile-rejected.html",
		nil,
	)
}

// --------------------------------------------------
// Submit Doctor Verification
// --------------------------------------------------

func (h *doctorHandler) SubmitVerification(c *gin.Context) {

	logger.Log.Info("Doctor profile submission started")

	// Get user ID

	userIDValue, exists := c.Get("UserID")

	if !exists {

		logger.Log.Warn(
			"Doctor profile submission failed: user ID not found",
		)

		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
		})

		return
	}

	userID, ok := userIDValue.(uint)

	if !ok {

		logger.Log.Warn(
			"Doctor profile submission failed: invalid user ID type",
		)

		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})

		return
	}

	logger.Log.Info(
		"Doctor user identified",
		zap.Uint("user_id", userID),
	)

	var input dto.DoctorProfileSubmission
	var messages []string

	// Parse profile data

	input.ProfessionalRole = c.PostForm("professional_role")
	input.HospitalAddress = c.PostForm("hospital_address")
	input.Biography = c.PostForm("biography")

	// Department

	departmentID, err := strconv.ParseUint(
		c.PostForm("department_id"),
		10,
		64,
	)

	if err != nil {

		logger.Log.Warn(
			"Invalid department ID",
			zap.Uint("user_id", userID),
			zap.String("error", err.Error()),
		)

		messages = append(
			messages,
			"Invalid department ID",
		)

	} else {

		input.DepartmentID = uint(departmentID)

		logger.Log.Info(
			"Department ID parsed",
			zap.Uint("user_id", userID),
			zap.Uint("department_id", input.DepartmentID),
		)
	}

	// Experience

	experience, err := strconv.ParseUint(
		c.PostForm("experience"),
		10,
		64,
	)

	if err != nil {

		logger.Log.Warn(
			"Invalid years of experience",
			zap.Uint("user_id", userID),
			zap.String("error", err.Error()),
		)

		messages = append(
			messages,
			"Invalid years of experience",
		)

	} else {

		input.Experience = uint(experience)

		logger.Log.Info(
			"Doctor experience parsed",
			zap.Uint("user_id", userID),
			zap.Uint("experience", input.Experience),
		)
	}

	// Qualifications

	qualificationsJSON := c.PostForm("qualifications")

	if qualificationsJSON != "" {

		err = json.Unmarshal(
			[]byte(qualificationsJSON),
			&input.Qualifications,
		)

		if err != nil {

			logger.Log.Warn(
				"Invalid qualifications data",
				zap.Uint("user_id", userID),
				zap.String("error", err.Error()),
			)

			messages = append(
				messages,
				"Invalid qualifications data",
			)

		} else {

			logger.Log.Info(
				"Doctor qualifications parsed",
				zap.Uint("user_id", userID),
				zap.Int("qualification_count", len(input.Qualifications)),
			)
		}
	}

	// Experiences

	experiencesJSON := c.PostForm("experiences")

	if experiencesJSON != "" {

		err = json.Unmarshal(
			[]byte(experiencesJSON),
			&input.Experiences,
		)

		if err != nil {

			logger.Log.Warn(
				"Invalid professional experience data",
				zap.Uint("user_id", userID),
				zap.String("error", err.Error()),
			)

			messages = append(
				messages,
				"Invalid professional experience data",
			)

		} else {

			logger.Log.Info(
				"Doctor professional experiences parsed",
				zap.Uint("user_id", userID),
				zap.Int("experience_count", len(input.Experiences)),
			)
		}
	}

	// Update doctor profile

	logger.Log.Info(
		"Updating doctor profile",
		zap.Uint("user_id", userID),
	)

	err = h.doctorService.SubmitDoctorProfile(
		userID,
		input,
	)

	if err != nil {

		logger.Log.Error(
			"Doctor profile update failed",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)

		messages = append(
			messages,
			"Profile update failed: "+err.Error(),
		)

	} else {

		logger.Log.Info(
			"Doctor profile updated successfully",
			zap.Uint("user_id", userID),
		)
	}

	// Profile photo

	var profilePhotoResult interface{}

	file, header, err := c.Request.FormFile("file")

	if err != nil {

		logger.Log.Info(
			"Doctor profile photo not uploaded",
			zap.Uint("user_id", userID),
		)

		messages = append(
			messages,
			"Profile photo was not uploaded",
		)

	} else {

		defer file.Close()

		category := models.FileCategoryProfilePhoto

		logger.Log.Info(
			"Doctor profile photo received",
			zap.Uint("user_id", userID),
			zap.String("filename", header.Filename),
		)

		// Validate profile photo

		if err := validateFile(category, header); err != nil {

			logger.Log.Warn(
				"Doctor profile photo validation failed",
				zap.Uint("user_id", userID),
				zap.String("filename", header.Filename),
				zap.Error(err),
			)

			messages = append(
				messages,
				"Profile photo: "+err.Error(),
			)

		} else {

			// Upload profile photo

			result, err := h.fileService.UploadFile(
				c.Request.Context(),
				userID,
				category,
				file,
				header,
			)

			if err != nil {

				logger.Log.Error(
					"Doctor profile photo upload failed",
					zap.Uint("user_id", userID),
					zap.String("filename", header.Filename),
					zap.Error(err),
				)

				messages = append(
					messages,
					"Profile photo upload failed: "+err.Error(),
				)

			} else {

				logger.Log.Info(
					"Doctor profile photo uploaded successfully",
					zap.Uint("user_id", userID),
					zap.String("filename", header.Filename),
				)

				profilePhotoResult = result
			}
		}
	}

	// Verification documents

	var uploadedDocuments []interface{}

	multipartForm, err := c.MultipartForm()

	if err != nil {

		logger.Log.Warn(
			"Could not read multipart form",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)

		messages = append(
			messages,
			"Could not read uploaded documents",
		)

	} else {

		documentFiles := multipartForm.File["documents"]
		documentCategories := c.PostFormArray("document_categories")
		documentRemarks := c.PostFormArray("document_remarks")

		logger.Log.Info(
			"Doctor verification documents received",
			zap.Uint("user_id", userID),
			zap.Int("document_count", len(documentFiles)),
		)

		// No documents is not an error.
		// The doctor can upload them later.

		if len(documentFiles) > 0 {

			// Validate document information

			if len(documentCategories) != len(documentFiles) ||
				len(documentRemarks) != len(documentFiles) {

				logger.Log.Warn(
					"Doctor document information is incomplete",
					zap.Uint("user_id", userID),
					zap.Int("document_count", len(documentFiles)),
					zap.Int("category_count", len(documentCategories)),
					zap.Int("remarks_count", len(documentRemarks)),
				)

				messages = append(
					messages,
					"Some document information is incomplete",
				)

			} else {

				for i, documentHeader := range documentFiles {

					category := models.FileCategory(
						documentCategories[i],
					)

					remarks := documentRemarks[i]

					// Validate category

					switch category {

					case models.FileCategoryDoctorCertificate,
						models.FileCategoryIdentityProof,
						models.FileCategoryExperienceProof:

					default:

						logger.Log.Warn(
							"Invalid doctor document category",
							zap.Uint("user_id", userID),
							zap.String("filename", documentHeader.Filename),
							zap.String("category", string(category)),
						)

						messages = append(
							messages,
							documentHeader.Filename+
								": invalid document category",
						)

						continue
					}

					// Validate file

					if err := validateFile(
						category,
						documentHeader,
					); err != nil {

						logger.Log.Warn(
							"Doctor document validation failed",
							zap.Uint("user_id", userID),
							zap.String("filename", documentHeader.Filename),
							zap.String("category", string(category)),
							zap.Error(err),
						)

						messages = append(
							messages,
							documentHeader.Filename+
								": "+err.Error(),
						)

						continue
					}

					// Open file

					documentFile, err :=
						documentHeader.Open()

					if err != nil {

						logger.Log.Error(
							"Unable to open doctor document",
							zap.Uint("user_id", userID),
							zap.String("filename", documentHeader.Filename),
							zap.Error(err),
						)

						messages = append(
							messages,
							documentHeader.Filename+
								": unable to open file",
						)

						continue
					}

					// Upload document

					logger.Log.Info(
						"Uploading doctor verification document",
						zap.Uint("user_id", userID),
						zap.String("filename", documentHeader.Filename),
						zap.String("category", string(category)),
					)

					result, err := h.fileService.UploadFile(
						c.Request.Context(),
						userID,
						category,
						documentFile,
						documentHeader,
					)

					documentFile.Close()

					if err != nil {

						logger.Log.Error(
							"Doctor document upload failed",
							zap.Uint("user_id", userID),
							zap.String("filename", documentHeader.Filename),
							zap.String("category", string(category)),
							zap.Error(err),
						)

						messages = append(
							messages,
							documentHeader.Filename+
								": upload failed - "+
								err.Error(),
						)

						continue
					}

					// Successfully uploaded

					logger.Log.Info(
						"Doctor verification document uploaded successfully",
						zap.Uint("user_id", userID),
						zap.String("filename", documentHeader.Filename),
						zap.String("category", string(category)),
					)

					uploadedDocuments = append(
						uploadedDocuments,
						gin.H{
							"file":    result,
							"remarks": remarks,
						},
					)
				}
			}
		}
	}

	// Final response

	response := gin.H{
		"success": true,
		"message": "Profile update completed",
	}

	if profilePhotoResult != nil {

		response["profile_photo"] = profilePhotoResult
	}

	if len(uploadedDocuments) > 0 {

		response["documents"] = uploadedDocuments
	}

	if len(messages) > 0 {

		response["messages"] = messages
	}

	logger.Log.Info(
		"Doctor profile submission completed",
		zap.Uint("user_id", userID),
		zap.Bool("profile_photo_uploaded", profilePhotoResult != nil),
		zap.Int("documents_uploaded", len(uploadedDocuments)),
		zap.Int("message_count", len(messages)),
	)

	c.JSON(
		http.StatusOK,
		response,
	)
}

// --------------------------------------------------
// Get Doctor Basic Profile
// --------------------------------------------------

func (h *doctorHandler) GetDoctorBasicProfile(c *gin.Context) {

	logger.Log.Info(
		"Doctor basic profile requested",
	)

	userIDValue, exists := c.Get("UserID")

	if !exists {

		logger.Log.Warn(
			"Doctor basic profile request failed: user ID not found",
		)

		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
		})

		return
	}

	userID, ok := userIDValue.(uint)

	if !ok {

		logger.Log.Warn(
			"Doctor basic profile request failed: invalid user ID type",
		)

		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})

		return
	}

	logger.Log.Info(
		"Fetching doctor basic profile",
		zap.Uint("user_id", userID),
	)

	profile, err := h.doctorService.GetDoctorBasicProfile(userID)

	if err != nil {

		logger.Log.Error(
			"Failed to fetch doctor basic profile",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})

		return
	}

	logger.Log.Info(
		"Doctor basic profile fetched successfully",
		zap.Uint("user_id", userID),
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"profile": profile,
	})
}
