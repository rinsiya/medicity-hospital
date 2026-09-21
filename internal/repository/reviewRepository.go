package repository

import (
	"medicity/database"
	"medicity/internal/models"


)

// ReviewRepository defines the interface for database operations on doctor_reviews table
type ReviewRepository interface {
	GetAllReviews() ([]models.DoctorReview, error)
	GetReviewsByDoctorID(doctorID uint) ([]models.DoctorReview, error)
	CreateReview(review *models.DoctorReview) error
}

type reviewRepository struct {}

// NewReviewRepository creates a new instance of ReviewRepository
func NewReviewRepository() ReviewRepository {
	return &reviewRepository{}
}

// GetAllReviews fetches all review records from 'doctor_reviews' table ordered by newest first
func (r *reviewRepository) GetAllReviews() ([]models.DoctorReview, error) {
	var reviews []models.DoctorReview
	err := database.DB.Order("created_at desc").Find(&reviews).Error
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

// GetReviewsByDoctorID fetches reviews for a specific doctor
func (r *reviewRepository) GetReviewsByDoctorID(doctorID uint) ([]models.DoctorReview, error) {
	var reviews []models.DoctorReview
	err := database.DB.Where("doctor_id = ?", doctorID).Order("created_at desc").Find(&reviews).Error
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

// CreateReview inserts a new review record into 'doctor_reviews'
func (r *reviewRepository) CreateReview(review *models.DoctorReview) error {
	return database.DB.Create(review).Error
}
