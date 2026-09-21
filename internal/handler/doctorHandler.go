package handler

import (
	//"medicity/internal/dto"
	"medicity/internal/service"
	"medicity/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)


type doctorHandler struct {
	doctorService service.DoctorService
}

func NewDoctorHandler(doctorService service.DoctorService) *doctorHandler {
	return &doctorHandler{
		doctorService: doctorService,
	}
}


func (h *doctorHandler) DoctorSignup(c *gin.Context) {
	c.HTML(http.StatusOK, "doctorSignup.html", nil)
}

func (h *doctorHandler) Dashboard(c *gin.Context) {
//userID := c.GetUint("UserID")
doctorID := c.GetUint("DoctorID")
	is_profile_complete := h.doctorService.IsProfileComplete(doctorID)
	if is_profile_complete{
	c.HTML(http.StatusOK, "doctorForVerification.html", nil)

	}else{
	c.HTML(http.StatusOK, "doctorCompleteProfile.html", nil)

	}

}

func (h *doctorHandler) ChangePassword(c *gin.Context) {
	c.HTML(http.StatusOK, "change-password.html", nil)}
// logout function
func (h *doctorHandler) Logout(c *gin.Context) {
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

	c.Redirect(http.StatusSeeOther, "/doctor/login")
}

func (h *doctorHandler) UpdatePassword(c *gin.Context) {
	
}
func (h *doctorHandler) CompleteProfile(c *gin.Context) {

	c.HTML(http.StatusOK, "doctorCompleteProfile.html", nil)
}

func (h *doctorHandler) VerificationPending(c *gin.Context) {
	c.HTML(http.StatusOK, "doctorWaitingForVerification.html", nil)
}

func (h *doctorHandler) ProfileRejected(c *gin.Context) {
	c.HTML(http.StatusOK, "profile-rejected.html", nil)
}




