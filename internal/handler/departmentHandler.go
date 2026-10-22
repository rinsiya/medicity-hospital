package handler

import (
	"net/http"

	"medicity/internal/service"

	"github.com/gin-gonic/gin"
)

type departmentHandler struct {
	departmentService service.DepartmentService
}

func NewDepartmentHandler(
	departmentService service.DepartmentService,
) *departmentHandler {
	return &departmentHandler{
		departmentService: departmentService,
	}
}

func (h *departmentHandler) GetDepartments(c *gin.Context) {

	departments, err := h.departmentService.GetAllDepartments()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch departments",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"departments": departments,
	})
}
