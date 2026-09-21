package handler

import (
	"errors"
	"net/http"
	"net/url"

	"medicity/internal/dto"
	"medicity/internal/models"
	"medicity/internal/service"
	"medicity/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SignupHandler struct {
	signupService service.SignupService
	otpService service.OTPService
}

func NewSignupHandler(
	signupService service.SignupService,
	otpService service.OTPService,
) *SignupHandler {
	return &SignupHandler{
		signupService: signupService,
		otpService:    otpService,
	}
}
func (h *SignupHandler) Signup(c *gin.Context) {

	var input dto.SignupInput

	logger.Log.Info(
		"Signup request received",
		zap.String("role", c.Param("role")),
	)

	// Bind and validate form
	if err := c.ShouldBind(&input); err != nil {

		logger.Log.Warn(
			"Signup validation failed",
			zap.String("role", c.Param("role")),
			zap.Error(err),
		)

		c.HTML(http.StatusBadRequest, c.Param("role")+"Signup.html", gin.H{
			"error": err.Error(),
		})

		return
	}

	// Determine role
	var role models.UserRole

	switch c.Param("role") {

	case "patient":
		role = models.RolePatient

	case "doctor":
		role = models.RoleDoctor

	default:

		logger.Log.Warn(
			"Invalid signup role",
			zap.String("role", c.Param("role")),
			zap.String("email", input.Email),
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid signup role",
		})

		return
	}

	logger.Log.Info(
		"Signup attempt",
		zap.String("role", string(role)),
		zap.String("email", input.Email),
		zap.String("phone", input.Phone),
	)

	phone, err := h.signupService.Signup(&input, role)

	if err != nil {

		switch {

		case errors.Is(err, service.ErrEmailAndPhoneAlreadyExists):

			c.HTML(http.StatusConflict, string(role)+"Signup.html", gin.H{
				"error": "Email and phone already registered. Try to login",
			})
return
		case errors.Is(err, service.ErrEmailOrPhoneAlreadyExists):

			c.HTML(http.StatusConflict, string(role)+"Signup.html", gin.H{
				"error": "Email or Phone number already registered. try again",
			})
return
		case errors.Is(err, service.ErrPhonePendingVerification):
//user already ergistered but verification pending
logger.Log.Info("pending user to verify")
			
	c.Redirect(http.StatusSeeOther, "/"+string(role)+"/verify-otp?phone="+url.QueryEscape(phone))
return
		case errors.Is(err, service.ErrEmailAndPhoneMismatch):

			c.HTML(http.StatusConflict, string(role)+"Signup.html", gin.H{
				"error": "Email or Phone already registered. try again",
			})
return
		default:

			logger.Log.Error(
				"Signup failed",
				zap.String("role", string(role)),
				zap.String("email", input.Email),
				zap.Error(err),
			)

			c.HTML(http.StatusOK, string(role)+"Signup.html", gin.H{
				"error": err,
			})
return
		}
	}

	logger.Log.Info(
		"Signup successful. phone number verification required",
		zap.String("role", string(role)),
		zap.String("email", input.Email),
	)

	c.Redirect(http.StatusSeeOther, "/"+string(role)+"/verify-otp?phone="+url.QueryEscape(input.Phone))
}

func (h *SignupHandler) ShowOTPPage(c *gin.Context) {

	role := c.Param("role")
	phone := c.Query("phone")

	logger.Log.Info(
		"OTP verification page requested",
		zap.String("role", role),
		zap.String("phone", phone),
	)

	// Find pending signup

pendingUser, err := h.signupService.FindPendingUserByPhone(phone)

if err != nil || pendingUser == nil {
	logger.Log.Error(
		"Pending user not found",
		zap.String("phone", phone),
		zap.Error(err),
	)

	c.HTML(http.StatusBadRequest, "verify-otp.html", gin.H{
		"phone": phone,
		"role":  role,
		"error": "Signup session not found. Please signup again.",
	})

	return
}
		
	logger.Log.Info(
		"load OTP verification page",
		zap.String("role", role),
		zap.String("phone", phone),
	)
	c.HTML(http.StatusOK, "verify-otp.html", gin.H{
		"role":         role,
		"phone":        phone,
		"otpExpiresAt": pendingUser.OTPExpiresAt.UnixMilli(),
	})
}



func (h *SignupHandler) ResendOTP(c *gin.Context) {

	phone := c.Query("phone")

	logger.Log.Info(
		"OTP resend requested",
		zap.String("phone", phone),
	)

	// Generate new OTP
	err := h.signupService.ResendOTP(phone)

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
	pendingUser, err := h.signupService.FindPendingUserByPhone(phone)

	if err != nil || pendingUser == nil {

		logger.Log.Error(
			"Failed to get new OTP expiry",
			zap.String("phone", phone),
			zap.Error(err),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed",
		})

		return
	}

	c.HTML(http.StatusOK, "verify-otp.html", gin.H{
		"phone":        phone,
		"message":      "OTP resent successfully. Check your SMS and enter the OTP carefully.",
		"otpExpiresAt": pendingUser.OTPExpiresAt.UnixMilli(),
	})
}

func (h *SignupHandler) PatientVerificationSuccess(c *gin.Context) {

	c.HTML(http.StatusOK, "patientVerificationSuccess.html", gin.H{
		"message": "Your phone number has been verified successfully.",
	},
	)
}

func (h *SignupHandler) DoctorVerificationSuccess(c *gin.Context) {

	c.HTML(http.StatusOK, "doctorVerificationSuccess.html", gin.H{
		"message": "Your phone number has been verified successfully.",
	},
	)
}

func (h *SignupHandler) ChangePhone(c *gin.Context) {
	phone := c.Query("phone")
	c.HTML(http.StatusOK, "change-phone.html", gin.H{
		"phone": phone,
	},
	)

}

func (h *SignupHandler) UpdatePhone(c *gin.Context) {

	oldPhone := c.PostForm("OldPhone")
	newPhone := c.PostForm("NewPhone")

	role, phone, err := h.signupService.UpdatePhone(oldPhone, newPhone)
if err != nil {

    if phone != "" {
        c.Redirect(
            http.StatusSeeOther,
            "/"+string(role)+"/verify-otp?phone="+url.QueryEscape(phone)+
                "&error="+url.QueryEscape("Phone number updated. can't send OTP. Try resend OTP"),
        )
        return
    }

    if errors.Is(err, service.ErrPhoneAlreadyExists) {

        c.HTML(http.StatusBadRequest, "change-phone.html", gin.H{
            "phone": oldPhone,
            "error": "This phone number is already registered",
        })

        logger.Log.Error(
            "Failed to update phone number. Number already exists",
            zap.String("oldPhone", oldPhone),
            zap.String("newPhone", newPhone),
            zap.Error(err),
        )

        return
    }

    logger.Log.Error(
        "Failed to update phone number",
        zap.String("oldPhone", oldPhone),
        zap.String("newPhone", newPhone),
        zap.Error(err),
    )

    c.Redirect(
        http.StatusSeeOther,
        "/"+string(role)+"/change-phone?phone="+url.QueryEscape(oldPhone),
    )

    return
}
	// Phone updated and OTP sent successfully
	c.Redirect(
		http.StatusSeeOther,
		"/"+string(role)+"/verify-otp?phone="+url.QueryEscape(phone),
	)
}

func (h *SignupHandler) ValidateOTP(c *gin.Context) {

	phone := c.PostForm("phone")
	role := c.Param("role")

	otp := c.PostForm("otp1") +
		c.PostForm("otp2") +
		c.PostForm("otp3") +
		c.PostForm("otp4") +
		c.PostForm("otp5") +
		c.PostForm("otp6")

	// Validate OTP length
	if len(otp) != 6 {

		c.HTML(http.StatusBadRequest, "verify-otp.html", gin.H{
			"phone": phone,
			"role":  role,
			"error": "Please enter a valid 6-digit OTP",
		})

		return
	}

	// Find pending signup
	pendingUser, err := h.signupService.FindPendingUserByPhone(phone)

	if err != nil || pendingUser == nil {

		logger.Log.Error(
			"Pending signup not found",
			zap.String("phone", phone),
			zap.Error(err),
		)

		c.HTML(http.StatusBadRequest, "verify-otp.html", gin.H{
			"phone": phone,
			"role":  role,
			"error": "Signup session not found. Please signup again.",
		})

		return
	}

	// Verify OTP
	otpExpiresAt, err := h.otpService.VerifyOTP(
		pendingUser.SessionID,
		otp,
		pendingUser.OTPExpiresAt,
	)

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
			"error": "OTP verification failed",
		}

		if !otpExpiresAt.IsZero() {
			data["otpExpiresAt"] = otpExpiresAt.UnixMilli()
			data["error"] = "OTP expired"
		}

		c.HTML(http.StatusBadRequest, "verify-otp.html", data)

		return
	}

	// Create actual user
	user, err := h.signupService.CreateUser(phone)

	if err != nil {

		logger.Log.Error(
			"Failed to create user after OTP verification",
			zap.String("phone", phone),
			zap.Error(err),
		)

		c.HTML(http.StatusInternalServerError, "verify-otp.html", gin.H{
			"phone": phone,
			"role":  role,
			"error": "Unable to complete signup. Please try again.",
		})

		return
	}

	logger.Log.Info(
		"OTP verified successfully. User created",
		zap.String("phone", user.Phone),
		zap.String("role", string(user.Role)),
	)

	switch user.Role {

	case models.RolePatient:

		c.Redirect(
			http.StatusSeeOther,
			"/patient/verification-success",
		)

	case models.RoleDoctor:

		c.Redirect(
			http.StatusSeeOther,
			"/doctor/verification-success",
		)

	default:

		c.Redirect(
			http.StatusSeeOther,
			"/"+role+"/verification-success",
		)
	}

}