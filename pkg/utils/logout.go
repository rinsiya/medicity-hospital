package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context, loginURL string) {

	// Clear JWT cookie
	c.SetCookie(
		"access_token",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	c.Redirect(http.StatusSeeOther, loginURL)
}