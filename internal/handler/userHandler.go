package handler

import (
	"errors"
	"medicity/internal/dto"
	"medicity/internal/models"
	"medicity/internal/service"
	"medicity/logger"
	"medicity/pkg/utils"
	"net/http"
	//"net/url"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	Service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		Service: service,
	}
}

func (h *UserHandler) Home(c *gin.Context) {
	c.HTML(200, "home.html", nil)
}

func (h *UserHandler) PatientLogin(c *gin.Context) {
	c.HTML(200, "patientLogin.html", nil)
}

func (h *UserHandler) DoctorLogin(c *gin.Context) {
	c.HTML(200, "doctorLogin.html", nil)
}

func (h *UserHandler) AdminLogin(c *gin.Context) {
	c.HTML(200, "adminLogin.html", nil)
}
func (h *UserHandler) VerifyPasswordResetOTP(c *gin.Context) {
	c.HTML(http.StatusOK, "forgot-password.html", gin.H{
		"role": c.Param("role"),
	})
}
func (h *UserHandler) ForgotPasswordForm(c *gin.Context) {
	c.HTML(http.StatusOK, "forgot-password.html", gin.H{
		"role": c.Param("role"),
	})
}

func (h *UserHandler) UpdateForgotPassword(c *gin.Context) {
session :=sessions.Default(c)
	phone, ok  := session.Get("reset_phone").(string)
    role,ok := session.Get("reset_role").(string)
if !ok || phone == "" || role == "" {
		logger.Log.Error("Reset phone number or role not found in session")
		c.HTML(http.StatusBadRequest, "forgot-password.html", gin.H{
			"role": role,
			"error": "unexpected error occurred. Please start the password reset process again.",
		})
		return
	}
	input := dto.ChangePasswordInput{}
	if err := c.ShouldBind(&input); err != nil {
		logger.Log.Error("Validation error", zap.Error(err))

		c.HTML(http.StatusBadRequest, "change-forgot-password.html", gin.H{
			"error": "Invalid input",
			"role":  role,
		})
		return
	}
	err := h.Service.UpdatePassword(phone,role, input.Password)
	if err != nil {
		logger.Log.Error("Failed to update password", zap.Error(err))

		c.HTML(http.StatusInternalServerError, "change-forgot-password.html", gin.H{
			"error": "Failed to update password. Please try again.",
			"role":  role,
		})
		return
	}
	  //  session = sessions.Default(c)
    session.Clear()

    // Save the cleared session
    if err := session.Save(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "message": "Password changed, but failed to clear session",
        })
        return
    }
	logger.Log.Info("Password updated successfully", zap.String("phone", phone), zap.String("role", role))
c.HTML(http.StatusOK, role+"Login.html", gin.H{
			"role": role,
			"message": "Password updated successfully. Please log in with your new password.",
		})
}
func (h *UserHandler) ForgotPasswordOTPVerificationRequest(c *gin.Context) {
	session := sessions.Default(c)
phoneValue := session.Get("reset_phone")
	phone, ok := phoneValue.(string)
		role := session.Get("reset_role").(string)
if !ok || phone == "" {
		logger.Log.Error("Reset phone number not found in session")
		c.HTML(http.StatusBadRequest, "forgot-password.html", gin.H{
			"role": role,
			"error": "unexpected error occurred. Please start the password reset process again.",
		})
		return
	}
	logger.Log.Info(
		"Password reset OTP verification started",
		zap.String("role", role),
		zap.String("phone", phone),
	)

otp := c.PostForm("otp1") +
		c.PostForm("otp2") +
		c.PostForm("otp3") +
		c.PostForm("otp4") +
		c.PostForm("otp5") +
		c.PostForm("otp6")

	// Validate OTP length
	if len(otp) != 6 {

		c.HTML(http.StatusBadRequest, "change-password-verification.html", gin.H{
			"phone": phone,
			"role":  role,
			"error": "Please enter a valid 6-digit OTP",
		})

		return
	}

	// Verify OTP and create user
	otpExpiresAt, err := h.Service.VerifyPasswordResetOTP(phone, otp)

	if err != nil {

		logger.Log.Error(
			"OTP verification failed",
			zap.String("phone", phone),
			zap.String("role", role),
			zap.Error(err),
		)

		data := gin.H{
			"phone": phone,
			"role":  role,
			"error": err.Error(),
		}

		// Only send expiry timestamp if available
		if !otpExpiresAt.IsZero() {
			data["otpExpiresAt"] = otpExpiresAt.UnixMilli()
		}

		c.HTML(http.StatusBadRequest, "change-password-verification.html", data)

		return
	}
	
	logger.Log.Info("OTP verified successfully",
		zap.String("phone", phone),
		zap.String("role", string(role)),
	)
		c.HTML(http.StatusSeeOther, "change-forgot-password.html", gin.H{
			"message": "OTP verified successfully. Please enter your new password.",
		"role":  role,
		})
	
}
func (h *UserHandler) ForgotPasswordOTPVerification(c *gin.Context) {
	session := sessions.Default(c)
	phoneValue := session.Get("reset_phone")
	phone, ok := phoneValue.(string)
	role := session.Get("reset_role").(string)

	if !ok || phone == "" {
		logger.Log.Error("Reset phone number not found in session")
		c.HTML(http.StatusBadRequest, "forgot-password.html", gin.H{
			"role": role,
			"error": "unexpected error occurred.",
		})
		return
	}

	
	logger.Log.Info(
		"OTP verification page requested for password reset",
		zap.String("phone", phone),
	)

		passwordResetUser, err := h.Service.FindPasswordResetUserByPhone(phone)

	if err != nil {

		logger.Log.Error(
			"Failed to find user",
			zap.String("phone", phone),
			zap.Error(err),
		)

		c.HTML(http.StatusInternalServerError, "change-password-verification.html", gin.H{
			"phone": phone,
			"role":  role,
			"error": "Unable to load OTP verification page.",
		})

		return
	}
	
	logger.Log.Info(
		"load OTP verification page",
		zap.String("role", role),
		zap.String("phone", phone),
	)

	c.HTML(http.StatusOK, "change-password-verification.html", gin.H{
		"role": role,
		"phone": phone,
		"otpExpiresAt": passwordResetUser.OTPExpiresAt.UnixMilli(),

	})

}
func (h *UserHandler) ForgotPasswordPhoneVerification(c *gin.Context) {
	role := c.Param("role")

	switch role {
	case "patient", "doctor", "admin":
		// valid role
	default:
		c.HTML(http.StatusBadRequest, "forgot-password.html", gin.H{
			"error": "Invalid user",
			"role":  role,
		})
		return
	}
	input := dto.ForgotPasswordPhone{}
	if err := c.ShouldBind(&input); err != nil {
		logger.Log.Error("Validation error", zap.Error(err))

		c.HTML(http.StatusBadRequest, "forgot-password.html", gin.H{
			"error": "Invalid phone number",
			"role":  role,
		})
		return
	}

	err := h.Service.UserExistByPhoneAndRole(input.Phone, role)

	if errors.Is(err, service.ErrUserNotFound) {

		logger.Log.Error("User not found", zap.Error(err))

		c.HTML(http.StatusNotFound, "forgot-password.html", gin.H{
			"error": "User not found",
			"role":  role,
		})
		return
	}

	if errors.Is(err, service.ErrFailedPasswordReset) {

		logger.Log.Error("Forgot password verification failed", zap.Error(err))
		c.HTML(http.StatusUnauthorized, "forgot-password.html", gin.H{
			"error": "Forgot password verification failed, try again later",
			"role":  role,
		})
		return
	}

	if err != nil {
		logger.Log.Error("Forgot password verification failed", zap.Error(err))
		c.HTML(http.StatusUnauthorized, "forgot-password.html", gin.H{
			"role":  role,
			"error": "Forgot password verification failed, try again later",
		})
		return

	}
	logger.Log.Info("Start session with reset password phone as session variable", zap.String("phone", input.Phone))
	session := sessions.Default(c)
	session.Set("reset_phone", input.Phone)
	session.Set("reset_role", role)

	if err := session.Save(); err != nil {
		logger.Log.Error("Failed to save reset phone to session", zap.Error(err))

		c.HTML(http.StatusInternalServerError, "forgot-password.html", gin.H{
			"role": role,

			"error": "Something went wrong. Please try again.",
		})
		return
	}
	c.Redirect(http.StatusSeeOther, "/"+role+"/forgot-password/otp-verification")

}

func (h *UserHandler) ResendOTP(c *gin.Context) {

session := sessions.Default(c)
	phoneValue := session.Get("reset_phone")
	phone, ok := phoneValue.(string)
	role := session.Get("reset_role").(string)
	
	if !ok || phone == "" || role == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid phone number",
		})
		return
	}
	logger.Log.Info(
		"OTP resend requested",
		zap.String("phone", phone),
	)

	// Generate new OTP
	err := h.Service.ResendOTP(phone)

	if err != nil {

		logger.Log.Error(
			"Failed to resend OTP",
			zap.String("phone", phone),
			zap.Error(err),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to resend OTP",
		})

		return
	}

	// Get the new expiry time
	passwordResetUser, err := h.Service.FindPasswordResetUserByPhone(phone)

	if err != nil || passwordResetUser == nil {

		logger.Log.Error(
			"Failed to get new OTP expiry",
			zap.String("phone", phone),
			zap.Error(err),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get OTP expiry",
		})

		return
	}

	c.HTML(http.StatusOK, "change-password-verification.html", gin.H{
		"phone":        phone,
		"message":      "OTP resent successfully. Check your SMS and enter the OTP carefully.",
		"otpExpiresAt": passwordResetUser.OTPExpiresAt.UnixMilli(),
	})
}

func (h *UserHandler) Login(c *gin.Context) {

	role := c.Param("role")

	switch role {
	case "patient", "doctor", "admin":
		// valid role
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid role",
		})
		return
	}

	input := dto.LoginInput{}
	if err := c.ShouldBind(&input); err != nil {
		logger.Log.Error("Validation error", zap.Error(err))

		c.HTML(http.StatusBadRequest, role+"Login.html", gin.H{
			"error": "Invalid input",
			"role":  role,
		})
		return
	}

	user, err := h.Service.Login(input, role)
	if err != nil {
		logger.Log.Error("Login failed", zap.Error(err))

		c.HTML(http.StatusUnauthorized,role+"Login.html", gin.H{
			"error": "Invalid credentials",
			"role":  role,
		})
		return
	}

patient_id := h.Service.GetPatientIDByUserID(user.UserID)

	token, err := utils.GenerateJWT(user.UserID,patient_id, string(user.Role))

	if err != nil {
		logger.Log.Error("Token generation failed", zap.Error(err))
		c.HTML(http.StatusUnauthorized, role+"Login.html", gin.H{
			"error": "something went wrong . please try again",
			"role":  role,
		})
		return
	}

	// Store JWT in HttpOnly cookie
	c.SetCookie(
		"access_token", // cookie name
		token,          // token
		86400,          // 24 hours
		"/",            // available throughout the application
		"",             // domain
		false,          // secure - set true when using HTTPS
		true,           // HttpOnly
	)

	switch user.Role {

	case models.RoleAdmin:
		c.Redirect(http.StatusSeeOther, "/admin/dashboard")

	case models.RolePatient:
		c.Redirect(http.StatusSeeOther, "/patient/home")

	case models.RoleDoctor:
		c.Redirect(http.StatusSeeOther, "/doctor/dashboard")

	default:
		logger.Log.Warn(
			"Unknown user role",
			zap.String("role", string(user.Role)),
		)

		c.HTML(http.StatusForbidden, role+"Login.html", gin.H{
			"error": "Invalid user role",
			"role":  role,
		})
	}
}
