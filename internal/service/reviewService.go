package service

import (
	"errors"
	"medicity/internal/models"
	"medicity/internal/repository"
)

// ReviewService defines the interface for review business logic
type ReviewService interface {
	GetAllReviews() ([]models.DoctorReview, error)
	GetReviewsByDoctorID(doctorID uint) ([]models.DoctorReview, error)
	CreateReview(review *models.DoctorReview) error
	GetAverageRating() (float64, error)
}

type reviewService struct {
	reviewRepository repository.ReviewRepository
}

// NewReviewService initializes a new ReviewService instance
func NewReviewService(reviewRepository repository.ReviewRepository) ReviewService {
	return &reviewService{reviewRepository: reviewRepository}
}

// GetAllReviews handles fetching all doctor reviews
func (s *reviewService) GetAllReviews() ([]models.DoctorReview, error) {
	return s.reviewRepository.GetAllReviews()
}

// GetReviewsByDoctorID retrieves reviews for a specified doctor ID
func (s *reviewService) GetReviewsByDoctorID(doctorID uint) ([]models.DoctorReview, error) {
	if doctorID == 0 {
		return nil, errors.New("invalid doctor ID")
	}
	return s.reviewRepository.GetReviewsByDoctorID(doctorID)
}

// CreateReview validates and delegates new review creation
func (s *reviewService) CreateReview(review *models.DoctorReview) error {
	if review.Rating < 1 || review.Rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}
	if len(review.Review) == 0 {
		return errors.New("review text cannot be empty")
	}
	return s.reviewRepository.CreateReview(review)
}

// GetAverageRating calculates the overall rating average from retrieved review
func (s *reviewService) GetAverageRating() (float64, error) {
	reviews, err := s.reviewRepository.GetAllReviews()
	if err != nil || len(reviews) == 0 {
		return 5.0, err
	}
	var total uint
	for _, r := range reviews {
		total += uint(r.Rating)
	}
	return float64(total) / float64(len(reviews)), nil
}
