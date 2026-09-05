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
	c.HTML(http.StatusOK, "doctorDashboard.html", nil)
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


