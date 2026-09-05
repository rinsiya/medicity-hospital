package repository

import (
	"errors"
	"medicity/database"
	"medicity/internal/models"

	"gorm.io/gorm"
)
type PasswordResetRepository interface {

	Create(passwordReset *models.PasswordReset) error
FindByPhone(phone string) (*models.PasswordReset, error)
	Update(passwordReset *models.PasswordReset) error
Delete(passwordResetID uint) error
}
type passwordResetRepository struct {}

func NewPasswordResetRepository() PasswordResetRepository {
	return &passwordResetRepository{}
}

func (r *passwordResetRepository) Create(passwordReset *models.PasswordReset) error {
	return database.DB.Create(passwordReset).Error
}

func (r *passwordResetRepository) FindByPhone(phone string) (*models.PasswordReset, error) {

	var passwordReset models.PasswordReset

	err := database.DB.Where("phone = ?", phone).First(&passwordReset).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &passwordReset, nil
}

func (r *passwordResetRepository) Update(passwordReset *models.PasswordReset) error {
	return database.DB.Save(passwordReset).Error
}
		func (r *passwordResetRepository) Delete(passwordResetID uint) error {
			return database.DB.Delete(&models.PasswordReset{}, passwordResetID).Error
		}




