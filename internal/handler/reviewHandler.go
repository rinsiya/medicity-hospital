package handler

import (
	"net/http"
	"strconv"

	"medicity/internal/models"
	"medicity/internal/service"

	"github.com/gin-gonic/gin"
)

type ReviewHandler interface {
	GetReviews(c *gin.Context)
	CreateReview(c *gin.Context)
}

// ReviewHandler manages HTTP request routing for doctor reviews
type reviewHandler struct {
	reviewService service.ReviewService
}

// NewReviewHandler returns a new instance of ReviewHandler
func NewReviewHandler(reviewService service.ReviewService) ReviewHandler {
	return &reviewHandler{reviewService: reviewService}
}

// GetReviews handles GET /api/reviews
// Responds with c.JSON containing the array of DoctorReview records from PostgreSQL 'doctor_reviews' table
func (h *reviewHandler) GetReviews(c *gin.Context) {
	doctorIDStr := c.Query("doctor_id")

	var reviews []models.DoctorReview
	var err error

	if doctorIDStr != "" {
		doctorID, pErr := strconv.ParseUint(doctorIDStr, 10, 32)
		if pErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid doctor ID format",
			})
			return
		}
		reviews, err = h.reviewService.GetReviewsByDoctorID(uint(doctorID))
	} else {
		reviews, err = h.reviewService.GetAllReviews()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to fetch review reviews: " + err.Error(),
		})
		return
	}

	// Calculate average rating
	avgRating, _ := h.reviewService.GetAverageRating()

	// Return c.JSON with array of reviews retrieved from doctor_reviews table
	c.JSON(http.StatusOK, gin.H{
		"status":         "success",
		"count":          len(reviews),
		"average_rating": avgRating,
		"data":           reviews,
	})
}

// CreateReview handles POST /api/reviews
// Accepts new review details, stores in PostgreSQL database, and responds with c.JSON
func (h *reviewHandler) CreateReview(c *gin.Context) {
	var input models.DoctorReview
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data: " + err.Error(),
		})
		return
	}

	if err := h.reviewService.CreateReview(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Review submitted successfully",
		"data":    input,
	})
}
